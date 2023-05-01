package logging

import (
	"strings"
	"testing"
	"time"
)

var testListSubscriber = NewListSubscriber()

func TestLogger_Info(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		factory *LoggerFactory
		args    args
	}{
		{
			name: "create logger and push some info messages",
			factory: NewLoggerFactory(
				NewLoggerConfig(
					WithTemplate("%v [%s] %s"),
					WithCallbacks(func() interface{} { return time.Now().Unix() }),
					WithSubscribers(testListSubscriber),
				),
			),
			args: args{
				message: "Some info message",
			},
		},
	}
	for _, tt := range tests {
		logger := tt.factory.GetLogger("Some very informative logger")
		t.Run(tt.name, func(t *testing.T) {
			logger.Info(tt.args.message)

			lastMessage := testListSubscriber.messages[len(testListSubscriber.messages)-1]
			if !strings.Contains(lastMessage, tt.args.message) {
				t.Errorf("Test message was not appended to subscriber")
			}
		})
	}
}
