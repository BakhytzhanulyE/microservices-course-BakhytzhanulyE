package kafka

// HeaderCarrier переносит контекст трассировки через заголовки Kafka-сообщения.
// Реализует propagation.TextMapCarrier из OpenTelemetry: пропагатор кладёт в него
// traceparent/tracestate на стороне продюсера и достаёт их на стороне консьюмера.
type HeaderCarrier map[string][]byte

// Get возвращает значение заголовка или пустую строку, если его нет.
func (c HeaderCarrier) Get(key string) string {
	return string(c[key])
}

// Set записывает заголовок.
func (c HeaderCarrier) Set(key, value string) {
	c[key] = []byte(value)
}

// Keys перечисляет имена заголовков.
func (c HeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}

	return keys
}
