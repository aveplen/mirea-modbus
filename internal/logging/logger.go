package logging

import (
	"fmt"
	"sync"
)

type Subscriber interface {
	Consume(message string)
}

type LoggerConfig struct {
	mu *sync.Mutex

	template    string
	callbacks   []func() interface{}
	subscribers []Subscriber
}

type LoggerConfigOption func(config *LoggerConfig)

func WithTemplate(template string) LoggerConfigOption {
	return func(config *LoggerConfig) {
		config.template = template
	}
}

func WithCallbacks(callbacks ...func() interface{}) LoggerConfigOption {
	return func(config *LoggerConfig) {
		if config.callbacks == nil {
			config.callbacks = make([]func() interface{}, 0, len(callbacks))
		}

		config.callbacks = append(config.callbacks, callbacks...)
	}
}

func WithSubscribers(subscribers ...Subscriber) LoggerConfigOption {
	return func(config *LoggerConfig) {
		if config.subscribers == nil {
			config.subscribers = make([]Subscriber, 0, len(subscribers))
		}

		config.subscribers = append(config.subscribers, subscribers...)
	}
}

func NewLoggerConfig(options ...LoggerConfigOption) *LoggerConfig {
	config := &LoggerConfig{
		mu: &sync.Mutex{},
	}

	for _, option := range options {
		option(config)
	}

	return config
}

type Logger struct {
	config *LoggerConfig
}

func (l *Logger) Info(message string) {
	results := make([]interface{}, 0, len(l.config.callbacks))
	for _, callback := range l.config.callbacks {
		results = append(results, callback())
	}

	results = append(results, "INFO")
	results = append(results, message)

	for _, sub := range l.config.subscribers {
		sub.Consume(fmt.Sprintf(l.config.template, results...))
	}
}

func (l *Logger) Infof(format string, values ...interface{}) {
	results := make([]interface{}, 0, len(l.config.callbacks))
	for _, callback := range l.config.callbacks {
		results = append(results, callback())
	}

	results = append(results, "INFO")
	results = append(results, fmt.Sprintf(format, values...))

	for _, sub := range l.config.subscribers {
		sub.Consume(fmt.Sprintf(l.config.template, results...))
	}
}

func (l *Logger) Debug(message string) {
	results := make([]interface{}, 0, len(l.config.callbacks))
	for _, callback := range l.config.callbacks {
		results = append(results, callback())
	}

	results = append(results, "DEBUG")
	results = append(results, message)

	for _, sub := range l.config.subscribers {
		sub.Consume(fmt.Sprintf(l.config.template, results...))
	}
}

func (l *Logger) Debugf(format string, values ...interface{}) {
	results := make([]interface{}, 0, len(l.config.callbacks))
	for _, callback := range l.config.callbacks {
		results = append(results, callback())
	}

	results = append(results, "DEBUG")
	results = append(results, fmt.Sprintf(format, values...))

	for _, sub := range l.config.subscribers {
		sub.Consume(fmt.Sprintf(l.config.template, results...))
	}
}

func (l *Logger) Error(err error) {
	results := make([]interface{}, 0, len(l.config.callbacks))
	for _, callback := range l.config.callbacks {
		results = append(results, callback())
	}

	results = append(results, "ERROR")
	results = append(results, err.Error())

	for _, sub := range l.config.subscribers {
		sub.Consume(fmt.Sprintf(l.config.template, results...))
	}
}

func (l *Logger) Errorf(format string, values ...interface{}) {
	results := make([]interface{}, 0, len(l.config.callbacks))
	for _, callback := range l.config.callbacks {
		results = append(results, callback())
	}

	results = append(results, "ERROR")
	results = append(results, fmt.Errorf(format, values...))

	for _, sub := range l.config.subscribers {
		sub.Consume(fmt.Sprintf(l.config.template, results...))
	}
}

func (l *Logger) AddSubscriber(subscriber Subscriber) {
	l.config.mu.Lock()
	defer l.config.mu.Unlock()
	l.config.subscribers = append(l.config.subscribers, subscriber)
}

type LoggerFactory struct {
	config *LoggerConfig
	cache  map[string]*Logger
}

func NewLoggerFactory(config *LoggerConfig) *LoggerFactory {
	loggerFactory := &LoggerFactory{
		cache:  make(map[string]*Logger),
		config: config,
	}
	return loggerFactory
}

func (f *LoggerFactory) GetLogger(name string) *Logger {
	_, ok := f.cache[name]

	if !ok {
		f.cache[name] = &Logger{
			config: f.config,
		}
	}

	return f.cache[name]
}

func (f *LoggerFactory) AddSubscriber(subscriber Subscriber) {
	f.config.mu.Lock()
	defer f.config.mu.Unlock()
	f.config.subscribers = append(f.config.subscribers, subscriber)
}
