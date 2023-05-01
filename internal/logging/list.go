package logging

type ListSubscriber struct {
	messages []string
}

func (l *ListSubscriber) Consume(message string) {
	l.messages = append(l.messages, message)
}

func NewListSubscriber() *ListSubscriber {
	listSubscriber := &ListSubscriber{
		messages: make([]string, 0, 20),
	}
	return listSubscriber
}
