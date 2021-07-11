package match

import (
	"github.com/ndphu/betmonitor/cache"
	"github.com/ndphu/betmonitor/config"
)

type Flip struct {
	User         config.User `json:"user"`
	From         string `json:"from"`
	To           string `json:"to"`
	MatchId      string `json:"matchId"`
	MatchDetails string `json:"matchDetails"`
}

func (flip *Flip) UpdateCache() error {
	return cache.Save("kmsbet:" + flip.MatchId + ":" + flip.User.Email, flip.To)
}
