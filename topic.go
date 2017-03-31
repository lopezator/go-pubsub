package queue

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

// Topic errors
var (
	ErrAlreadyExistTopic        = errors.New("already exist topic")
	ErrAlreadyExistSubscription = errors.New("already exist subscription")
	ErrNotFoundTopic            = errors.New("not found topic")
	ErrNotHasSubscription       = errors.New("not has subscription")
	ErrInvalidTopic             = errors.New("invalid topic")
	ErrInvalidSubscription      = errors.New("invalid subscription")
)

// Global variable Topic map
var GlobalTopics *topics = newTopics()

type topics struct {
	store Datastore
	mu    sync.RWMutex
}

func newTopics() *topics {
	return &topics{
		store: NewMemory(),
	}
}

func (ts *topics) Get(key string) (*Topic, bool) {
	t := ts.store.Get(key)
	if v, ok := t.(*Topic); ok {
		return v, true
	}
	return nil, false
}

func (ts *topics) List() ([]*Topic, error) {
	values := ts.store.Dump()
	res := make([]*Topic, 0, len(values))
	for _, v := range values {
		if vt, ok := v.(*Topic); ok {
			res = append(res, vt)
		} else {
			return nil, ErrInvalidTopic
		}
	}
	return res, nil
}

func (ts *topics) Set(topic *Topic) error {
	return ts.store.Set(topic.name, topic)
}

func (ts *topics) Delete(key string) error {
	return ts.store.Delete(key)
}

// Topic object
type Topic struct {
	name string
	sub  Datastore
}

// Create topic, if not exist already topic name in GlobalTopics
func NewTopic(name string) (*Topic, error) {
	if _, ok := GlobalTopics.Get(name); ok {
		return nil, ErrAlreadyExistTopic
	}

	t := &Topic{
		name: name,
		sub:  NewMemory(),
	}
	GlobalTopics.Set(t)
	return t, nil
}

// Return topic object
func GetTopic(name string) (*Topic, error) {
	t, ok := GlobalTopics.Get(name)
	if !ok {
		return nil, ErrNotFoundTopic
	}

	return t, nil
}

// Return topic list
func ListTopic() ([]*Topic, error) {
	return GlobalTopics.List()
}

// Delete topic object at GlobalTopics
func (t *Topic) Delete() {
	GlobalTopics.Delete(t.name)
}

// Register subscription
func (t *Topic) AddSubscription(s Subscription) error {
	if _, ok := t.sub.Get(s.name).(Subscription); ok {
		return errors.Wrapf(ErrAlreadyExistSubscription, fmt.Sprintf("id=%s", s.name))
	}
	return t.sub.Set(s.name, s)
}

// Message store backend storage and delivery to Subscription
func (t *Topic) Publish(data []byte, attributes map[string]string) error {
	subList, err := t.GetSubscriptions()
	if err != nil {
		return errors.Wrap(err, "failed GetSubscriptions")
	}
	m := NewMessage(makeMessageID(), *t, data, attributes, subList)
	// TODO: save message to store

	for _, s := range subList {
		s.Subscribe(m)
	}
	return nil
}

func (t *Topic) GetSubscription(key string) (Subscription, error) {
	v := t.sub.Get(key)
	if s, ok := v.(Subscription); !ok {
		return Subscription{}, ErrNotHasSubscription
	} else {
		return s, nil
	}
}

// GetSubscriptions returns topic dependent Subscription list
func (t *Topic) GetSubscriptions() ([]Subscription, error) {
	dump := t.sub.Dump()
	subList := make([]Subscription, 0, len(dump))
	for _, v := range dump {
		if s, ok := v.(Subscription); ok {
			subList = append(subList, s)
		} else {
			return nil, ErrInvalidSubscription
		}
	}
	return subList, nil
}
