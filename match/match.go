package match

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/ndphu/betmonitor/auth"
	"github.com/ndphu/betmonitor/cache"
	"github.com/ndphu/betmonitor/config"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Match struct {
	Key     string `json:"key"`
	Details string `json:"details"`
}

func (m *Match) GetBetList() ([]*Bet, error) {
	client := http.Client{}
	url := "http://kmsbet.appspot.com/matchBet?matchKey=" + m.Key
	log.Println("GetBetList:", "URL=", url)
	req, _ := http.NewRequest("GET", url, nil)

	auth.GetStore().SetCookie(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		payload, _ := ioutil.ReadAll(resp.Body)
		log.Println("GetBetList:", "Server return unexpected status code", resp.StatusCode, url)
		log.Println("GetBetList:", "Body", string(payload))
		return nil, errors.New("InvalidResponseStatus")
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("GetBetList:", "Fail to parse response body", err)
		return nil, err
	}

	bets := make([]*Bet, 0)

	doc.Find("ul > li").Each(func(i int, s *goquery.Selection) {
		re := regexp.MustCompile(`Người chơi (.*?) cá (.*?) lÃºc.*`)
		submatch := re.FindStringSubmatch(s.Text())
		if len(submatch) == 3 {
			b := &Bet{
				User: submatch[1],
				Bet:  submatch[2],
			}
			bets = append(bets, b)
		}
	})

	return bets, nil
}

func (m *Match) GetFlips() ([]*Flip, error) {
	bets, err := m.GetBetList()
	if err != nil {
		return nil, err
	}
	flips := make([]*Flip, 0)
	for _, bet := range bets {
		key := "kmsbet:" + m.Key + ":" + bet.User
		if exist, err := cache.Exists(key); err != nil {
			log.Println("BetFlip:", "Fail to check key exist", key, err)
			return nil, err
		} else if exist {
			// checking flip
			if vote, err := cache.Get(key); err != nil {
				log.Println("BetFlip:", "Fail to read key", key, err)
				return nil, err
			} else {
				if vote != bet.Bet {
					log.Println("BetFlip:", "User", bet.User, "flip the choice")
					flips = append(flips, &Flip{
						User:         config.FindUser(bet.User),
						From:         vote,
						To:           bet.Bet,
						MatchId:      m.Key,
						MatchDetails: m.Details,
					})
				}
			}
		} else {
			flips = append(flips, &Flip{
				User:         config.FindUser(bet.User),
				From:         "N/A",
				To:           bet.Bet,
				MatchId:      m.Key,
				MatchDetails: m.Details,
			})
		}
	}
	return flips, nil
}
