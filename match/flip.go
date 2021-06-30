package match

import "github.com/ndphu/betmonitor/cache"

type Flip struct {
	User         string `json:"user"`
	From         string `json:"from"`
	To           string `json:"to"`
	MatchId      string `json:"matchId"`
	MatchDetails string `json:"matchDetails"`
}

func (flip *Flip) UpdateCache() error {
	return cache.Save("kmsbet:" + flip.MatchId + ":" + flip.User, flip.To)
}
