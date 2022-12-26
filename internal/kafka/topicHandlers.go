package kafka

import (
	"sync"

	"github.com/google/uuid"
)

type topicHandlers struct {
	sync.RWMutex
	topics        map[string]chan struct{}       //topics
	handlers      map[string]HandleFx            //handlers
	topicHandlers map[string]map[string]struct{} //handlers indexed by topic
}

func newTopicHandlers() *topicHandlers {
	return &topicHandlers{
		topics:        make(map[string]chan struct{}),
		handlers:      make(map[string]HandleFx),
		topicHandlers: make(map[string]map[string]struct{}),
	}
}

func (t *topicHandlers) upsertTopic(topic string) (chan struct{}, bool) {
	t.Lock()
	defer t.Unlock()
	stopper, ok := t.topics[topic]
	if !ok {
		stopper = make(chan struct{})
		t.topics[topic] = stopper
		return stopper, true
	}
	return stopper, false
}

func (t *topicHandlers) readHandlers(topic string) map[string]HandleFx {
	t.RLock()
	defer t.RUnlock()
	handlers := make(map[string]HandleFx)
	for handlerId := range t.topicHandlers[topic] {
		handlers[handlerId] = t.handlers[handlerId]
	}
	return handlers
}

func (t *topicHandlers) writeHandler(topic string, handler HandleFx) (handlerId string) {
	t.Lock()
	defer t.Unlock()
	handlerId = uuid.Must(uuid.NewRandom()).String()
	t.handlers[handlerId] = handler
	if _, ok := t.topicHandlers[topic]; !ok {
		t.topicHandlers[topic] = make(map[string]struct{})
	}
	t.topicHandlers[topic][handlerId] = struct{}{}
	return
}

func (t *topicHandlers) deleteHandlers(topic string, handlerIds ...string) {
	t.Lock()
	defer t.Unlock()

	switch {
	default:
		//delete all handlers
		for handlerId := range t.topicHandlers[topic] {
			delete(t.handlers, handlerId)
			delete(t.topicHandlers[topic], handlerId)
		}
	case len(handlerIds) > 0:
		//delete specific handlers provided
		for _, handlerId := range handlerIds {
			delete(t.handlers, handlerId)
			delete(t.topicHandlers[topic], handlerId)
		}
	}
	if len(t.topicHandlers[topic]) <= 0 {
		if stopper, ok := t.topics[topic]; ok {
			delete(t.topics, topic)
			select {
			default:
				close(stopper)
			case <-stopper:
			}
		}
	}
}

func (t *topicHandlers) deleteTopics(topics ...string) {
	t.Lock()
	defer t.Unlock()

	switch {
	default:
		for topic, stopper := range t.topics {
			select {
			default:
				close(stopper)
			case <-stopper:
			}
			for handlerId := range t.topicHandlers[topic] {
				delete(t.topicHandlers[topic], handlerId)
				delete(t.handlers, handlerId)
			}
			delete(t.topicHandlers, topic)
		}
	case len(topics) > 0:
		for _, topic := range topics {
			stopper := t.topics[topic]
			select {
			default:
				close(stopper)
			case <-stopper:
			}
			for handlerId := range t.topicHandlers[topic] {
				delete(t.topicHandlers[topic], handlerId)
				delete(t.handlers, handlerId)
			}
			delete(t.topicHandlers, topic)
		}
	}
}
