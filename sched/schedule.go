package sched

import (
	"bufio"
	"encoding/json"
	"github.com/jasonlvhit/gocron"
	"github.com/ndphu/betmonitor/match"
	"github.com/ndphu/betmonitor/notifier"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func ScheduleJobs() {
	s := gocron.NewScheduler()
	s.Every(15).Second().Do(func() {
		getwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		nt := make([]*notifier.NotificationTarget, 0)
		payload, err := ioutil.ReadFile(path.Join(getwd, "targets.json"))
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(payload, &nt); err != nil {
			panic(err)
		}

		file, err := os.Open(path.Join(getwd, "matches.txt"))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {
			matchId := scanner.Text()
			_m := match.GetMatch(matchId)
			if _m != nil {
				func(m *match.Match) {
					if flips, err := m.GetFlips(); err != nil {
						log.Println("Scheduler:", m.Key, "Fail to notify flips by error", err)
						return
					} else {
						for _, flip := range flips {
							var notifyError error
							for _, n := range nt {
								skypeNotifier := notifier.NewSkypeNotifier(n.WorkerId, n.Target)
								if err := skypeNotifier.NotifyItem(flip); err != nil {
									notifyError = err
									break
								}
							}
							if notifyError == nil {
								// successfully sent notification
								if err := flip.UpdateCache(); err != nil {
									log.Println("Fail to update flip",err)
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
