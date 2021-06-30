package notifier

type NotificationTarget struct {
	WorkerId string `json:"workerId"`
	Target   string `json:"target"`
}
