package notification

import "context"

type PushoverService struct {
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Token    string `json:"token"`
}

func (s PushoverService) ServiceName() string { return "pushover" }

func (s PushoverService) SendNotification(ctx context.Context, n Notification) error {
	panic("not implemented")
}
