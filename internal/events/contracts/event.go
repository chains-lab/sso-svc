package contracts

type Event struct {
	Topic        string
	EventType    string
	EventVersion int32
	Key          string
	Payload      interface{}
}
