package match

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/ndphu/betmonitor/auth"
	"github.com/ndphu/betmonitor/cache"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type MatchList struct {
	Group string `json:"group"`
}

var matchCache = make(map[string]*Match, 0)

func NewMatchList(group string) *MatchList {
	return &MatchList{
		Group: group,
	}
}

func (ml *MatchList) List() ([]*Match, error) {
	req, _ := http.NewRequest("GET", "http://kmsbet.appspot.com/matchUser?batchKey="+ml.Group, nil)
	auth.GetStore().SetCookie(req)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("MatchList:", "Fail to retrieve match list", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("MatchList:", "Received unexpected status code", resp.StatusCode)
		return nil, errors.New("UnexpectedStatusCode")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("MatchList:", "Fail to parse response from server", err)
		return nil, err
	}

	matches := make([]*Match, 0)
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		anchor := s.Find("a").First()
		link := anchor.AttrOr("href", "")
		if strings.TrimSpace(link) == "" {
			return
		}
		//log.Println("MatchList:", "Found match", link)
		re := regexp.MustCompile(`/matchBet\?matchKey=(.*?)$`)
		submatch := re.FindStringSubmatch(link)
		if len(submatch) == 2 {
			match := &Match{
				Key:     submatch[1],
				Details: strings.TrimSpace(anchor.Text()),
			}
			matches = append(matches, match)
			saveMatch(match)
		}
	})

	return matches, nil
}

func InitMapDetails(groupId string) error {
	matches, err := NewMatchList(groupId).List()
	if err != nil {
		return err
	}
	for _, m := range matches {
		matchCache[m.Key] = m
	}
	return nil
}

func GetMatch(id string) (*Match) {
	if m, e := matchCache[id]; e {
		return m
	}
	if details, err := cache.Get("kmsbet:" + id + ":details"); err != nil {
		return nil
	} else {
		m := &Match{
			Key:     id,
			Details: details,
		}
		matchCache[id] = m
		return m
	}
}

func saveMatch(m *Match) error {
	return cache.Save("kmsbet:"+m.Key+":details", m.Details)
}
