package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ndphu/betmonitor/auth"
	"github.com/ndphu/betmonitor/cache"
	"github.com/ndphu/betmonitor/match"
	"github.com/ndphu/betmonitor/notifier"
	"github.com/ndphu/betmonitor/sched"
	"log"
	"os"
)

var groupId = os.Getenv("GROUP_ID")

type NotificationTarget struct {
	WorkerId string `json:"workerId"`
	Target   string `json:"target"`
}

func main() {
	store := auth.GetStore()
	store.Start()

	log.Println(store.Cookie())

	cache.Start()

	if groupId == "" {
		groupId = "aghzfmttc2JldHISCxIFQmF0Y2gYgICAgIDyiAoM"
	}

	if err := match.InitMapDetails(groupId); err != nil {
		log.Println("Fail to load match list")
		panic(err)
	}

	sched.ScheduleJobs()

	r := gin.Default()

	api := r.Group("/api")

	api.GET("/match/:matchId/bets", func(c *gin.Context) {
		matchId := c.Param("matchId")
		m := match.Match{Key: matchId}
		if bets, err := m.GetBetList(); err != nil {
			c.JSON(500, gin.H{"success": false, "error": err.Error()})
		} else {
			c.JSON(200, gin.H{"success": true, "bets": bets})
		}
	})

	api.GET("/matchList", func(c *gin.Context) {
		if matches, err := match.NewMatchList(groupId).List(); err != nil {
			c.JSON(500, gin.H{"success": false, "error": err.Error(),})
		} else {
			c.JSON(200, gin.H{"matches": matches, "success": true})
		}
	})

	api.GET("/match/:matchId/flips", func(c *gin.Context) {
		matchId := c.Param("matchId")
		m := match.GetMatch(matchId)
		if flips, err := m.GetFlips(); err != nil {
			c.JSON(500, gin.H{"success": false, "error": err.Error()})
		} else {
			c.JSON(200, gin.H{"success": true, "flips": flips})
		}
	})

	api.POST("/match/:matchId/notifyFlips", func(c *gin.Context) {
		matchId := c.Param("matchId")
		nt := make([]*NotificationTarget, 0)
		if err := c.ShouldBindJSON(&nt); err != nil {
			c.JSON(400, gin.H{"success": false, "error": err.Error()})
			return
		}
		m := match.GetMatch(matchId)
		if flips, err := m.GetFlips(); err != nil {
			log.Println("Fail to notify flips by error", err)
			c.JSON(500, gin.H{"success": false, "error": err.Error()})
		} else {
			for _, flip := range flips {
				skypeNotifier := notifier.NewSkypeNotifier()
				if err := skypeNotifier.NotifyItem(flip); err != nil {
					log.Printf("Fail to notify flip %v\n", flip)
				} else {
					// successfully sent notification
					if err := flip.UpdateCache(); err != nil {
						log.Println("Fail to update flip", err)
					}
				}
			}
			c.JSON(200, gin.H{"success": true, "flips": flips})
		}
	})

	r.Run(":1990")
}
