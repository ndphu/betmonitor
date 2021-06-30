package notifier

import
(
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ndphu/betmonitor/match"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Notifier interface {
	Id() string
	NotifyItem(item *match.Flip) error
	NotifyItems(item []*match.Flip) error
}

type SkypeNotifier struct {
	Receiver   string `json:"receiver"`
	MessageUrl string `json:"messageUrl"`
}

type SkypeMessage struct {
	Target string `json:"target"`
	Text   string `json:"text"`
}

type SkypeResponse struct {
	Success bool `json:"success"`
}

func NewSkypeNotifier(workerId, receiver string) Notifier {
	skypeHost := os.Getenv("SKYPE_HOST")
	skypePort := os.Getenv("SKYPE_PORT")

	return &SkypeNotifier{
		Receiver:   receiver,
		MessageUrl: fmt.Sprintf("http://%s:%s/api/skype/worker/%s/message/text", skypeHost, skypePort, workerId),
	}
}

func (sn *SkypeNotifier) NotifyItem(flip *match.Flip) (error) {
	sm := SkypeMessage{
		Target: sn.Receiver,
		Text:   flip.MatchDetails + "\n" + flip.User + "\nUpdated:\n" + flip.From + " ---> " + flip.To,
	}
	payload, err := json.Marshal(sm)
	if err != nil {
		log.Printf("fail to send skype message by error %v\n", err)
		return err
	}
	req, _ := http.NewRequest("POST", sn.MessageUrl, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("fail to send skype message by error %v\n", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("fail to send skype message by error: skypebot return status code: %d\n", resp.StatusCode)
		return errors.New("InvalidStatusCode")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	sr := SkypeResponse{}
	if err := json.Unmarshal(body, &sr); err != nil {
		return err
	}

	if sr.Success {
		return nil
	}

	log.Println("fail to send skype message by error:", sr.Success)
	return errors.New("FailToSendMessage")
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
