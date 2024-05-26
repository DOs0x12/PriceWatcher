package price

type SubscribingComm struct{}

func (SubscribingComm) SubscribeUser() string {
	return "Subscribed!"
}
