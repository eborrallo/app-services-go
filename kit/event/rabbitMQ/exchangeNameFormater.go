package rabbitMQ

func Retry(exchange string) string {
	return "retry-" + exchange
}
func DeadLetter(exchange string) string {
	return "dead_letter-" + exchange
}
