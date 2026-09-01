package kafka

import "time"

// Message — сообщение Kafka без привязки к конкретной библиотеке.
type Message struct {
	Headers        map[string][]byte
	Timestamp      time.Time
	BlockTimestamp time.Time

	Key       []byte
	Value     []byte
	Topic     string
	Partition int32
	Offset    int64
}
