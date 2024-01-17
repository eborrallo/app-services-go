package rabbitMQ

import (
	"app-services-go/kit/event"
	"reflect"
	"regexp"
	"strings"
)

type QueueFormatter struct {
	moduleName string
}

func NewQueueFormatter(moduleName string) *QueueFormatter {
	return &QueueFormatter{moduleName: moduleName}
}

func name(subscriber event.Subscriber) string {
	return reflect.TypeOf(subscriber).Name()
}
func regSplit(text string, delimeter string) string {
	re := regexp.MustCompile(delimeter)
	snakeCase := re.ReplaceAllStringFunc(text, func(match string) string {
		return match[:1] + "_" + match[1:]
	})
	return strings.ToLower(snakeCase)
}
func (b *QueueFormatter) Format(subscriber event.Subscriber) string {
	value := name(subscriber)
	fullName := regSplit(value, `([a-z0-9])([A-Z])`)

	return strings.ToLower(b.moduleName + "." + fullName)
}
func (b *QueueFormatter) formatRetry(subscriber event.Subscriber) string {
	return "retry." + b.Format(subscriber)
}
func (b *QueueFormatter) formatDeadLetter(subscriber event.Subscriber) string {
	return "dead_letter." + b.Format(subscriber)

}
