package notifier

import (
	"github.com/ndphu/betmonitor/match"
	"github.com/ndphu/betmonitor/utils"
	"log"
	"strings"
)

type Notifier interface {
	Id() string
	NotifyItem(item *match.Flip) error
	NotifyItems(item []*match.Flip) error
}

type SkypeNotifier struct {
}

func NewSkypeNotifier() Notifier {
	return &SkypeNotifier{
	}
}

func (sn *SkypeNotifier) NotifyItem(flip *match.Flip) (error) {
	user := flip.User.Alias
	if user == "" {
		user = flip.User.Email
	}
	msg := flip.MatchDetails + "\n" + strings.Repeat("-", len(flip.MatchDetails)) + "\n" + user + "\nUpdated:\n" + flip.From + " ---> " + flip.To
	log.Println("Flip message:\n", msg)
	return utils.Reply(msg)
}

func (sn *SkypeNotifier) NotifyItems(flips []*match.Flip) (error) {
	for _, flips := range flips {
		if err := sn.NotifyItem(flips); err != nil {
			return err
		}
	}
	return nil
}

func (sn *SkypeNotifier) Id() string {
	return "skype"
}
