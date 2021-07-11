package sched

import (
	"github.com/jasonlvhit/gocron"
	"github.com/ndphu/betmonitor/config"
	"github.com/ndphu/betmonitor/match"
	"github.com/ndphu/betmonitor/notifier"
	"log"
)

func ScheduleJobs() {
	s := gocron.NewScheduler()
	s.Every(15).Second().Do(func() {
		conf := config.GetConfig()

		for _, matchId := range conf.Matches {
			_m := match.GetMatch(matchId)
			if _m != nil {
				func(m *match.Match) {
					if flips, err := m.GetFlips(); err != nil {
						log.Println("Scheduler:", m.Key, "Fail to notify flips by error", err)
						return
					} else {
						for _, flip := range flips {
							skypeNotifier := notifier.NewSkypeNotifier()
							if err := skypeNotifier.NotifyItem(flip); err != nil {
								log.Printf("Fail to notify flip %v\n", flip)
							} else {
								// successfully sent notification
								log.Printf("Successfully notify flip %v\n", flip)
								if err := flip.UpdateCache(); err != nil {
									log.Println("Fail to update flip", err)
								}
							}

						}
					}
				}(_m)
			}
		}
	})
	s.Start()
}
