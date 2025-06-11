package reader

// TODO now is not used, but in future we can use it to read from kafka topics need to remade this package to use it
// for each topic each listener

//type TopicReader struct {
//	reader *kafka.Reader
//}
//
//func NewTopicReader(r *kafka.Reader) *TopicReader {
//	return &TopicReader{
//		reader: r,
//	}
//}
//
//func (r *TopicReader) ListenChan(ctx context.Context) <-chan bodies.InternalEvent {
//	out := make(chan bodies.InternalEvent)
//
//	go func() {
//		defer r.reader.Close()
//		defer close(out)
//
//		for {
//			m, err := r.reader.ReadMessage(ctx)
//			if err != nil {
//				if errors.Is(ctx.Err(), context.Canceled) {
//					//r.log.Info("Context canceled, stopping listener")
//					return
//				}
//				//r.log.WithError(err).Error("Error reading message")
//				continue
//			}
//
//			var ie bodies.InternalEvent
//			if err := json.Unmarshal(m.Value, &ie); err != nil {
//				//r.log.WithError(err).Error("Error unmarshalling InternalEvent")
//				continue
//			}
//
//			select {
//			case out <- ie:
//			case <-ctx.Done():
//				//r.log.Info("Context canceled while sending message to channel")
//				return
//			}
//		}
//	}()
//
//	return out
//}
