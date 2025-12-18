package contracts

//type Message struct {
//	Topic        string
//	EventType    string
//	EventVersion int32
//	Key          string
//	Payload      interface{}
//}
//
//type OutboxEvent struct {
//	ID           uuid.UUID `db:"id"`
//	Topic        string    `db:"topic"`
//	EventType    string    `db:"event_type"`
//	EventVersion int       `db:"event_version"`
//	Key          string    `db:"key"`
//	Payload      []byte    `db:"payload"`
//
//	Status      string     `db:"status"`
//	Attempts    int        `db:"attempts"`
//	NextRetryAt time.Time  `db:"next_retry_at"`
//	CreatedAt   time.Time  `db:"created_at"`
//	SentAt      *time.Time `db:"sent_at"`
//}
//
//func (e OutboxEvent) ToMessage() Message {
//	return Message{
//		Topic:        e.Topic,
//		EventType:    e.EventType,
//		EventVersion: int32(e.EventVersion),
//		Key:          e.Key,
//		Payload:      json.RawMessage(e.Payload),
//	}
//}
