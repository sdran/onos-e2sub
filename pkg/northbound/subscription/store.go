// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package subscription

import (
	"context"
	"errors"
	"io"
	"time"

	_map "github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/gogo/protobuf/proto"
	subapi "github.com/onosproject/onos-e2sub/api/e2/subscription/v1beta1"
	"github.com/onosproject/onos-e2sub/pkg/config"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
)

// NewAtomixStore returns a new persistent Store
func NewAtomixStore() (Store, error) {
	ricConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	database, err := atomix.GetDatabase(ricConfig.Atomix, ricConfig.Atomix.GetDatabase(atomix.DatabaseTypeConsensus))
	if err != nil {
		return nil, err
	}

	subscriptions, err := database.GetMap(context.Background(), "subscriptions")
	if err != nil {
		return nil, err
	}

	return &atomixStore{
		subscriptions: subscriptions,
	}, nil
}

// NewLocalStore returns a new local end-point store
func NewLocalStore() (Store, error) {
	node, address := atomix.StartLocalNode()
	name := primitive.Name{
		Namespace: "local",
		Name:      "subscriptions",
	}

	session, err := primitive.NewSession(context.TODO(), primitive.Partition{ID: 1, Address: address})
	if err != nil {
		return nil, err
	}

	subscriptions, err := _map.New(context.Background(), name, []*primitive.Session{session})
	if err != nil {
		return nil, err
	}

	return &atomixStore{
		subscriptions: subscriptions,
		closer:        node.Stop,
	}, nil
}

// Store stores end-point registry information
type Store interface {
	io.Closer

	// Store stores a subscription in the store
	Store(point *subapi.Subscription) error

	// Delete deletes an subscription from the store
	Get(subapi.ID) (*subapi.Subscription, error)

	// Delete deletes an subscription from the store
	Delete(subapi.ID) error

	// List streams subscriptions to the given channel
	List(chan<- *subapi.Subscription) error

	// Watch streams subscription events to the given channel
	Watch(chan<- *Event, ...WatchOption) error
}

// WatchOption is a configuration option for Watch calls
type WatchOption interface {
	apply([]_map.WatchOption) []_map.WatchOption
}

// watchReplyOption is an option to replay events on watch
type watchReplayOption struct {
}

func (o watchReplayOption) apply(opts []_map.WatchOption) []_map.WatchOption {
	return append(opts, _map.WithReplay())
}

// WithReplay returns a WatchOption that replays past changes
func WithReplay() WatchOption {
	return watchReplayOption{}
}

// atomixStore is the implementation of the subscription Store
type atomixStore struct {
	subscriptions _map.Map
	closer        func() error
}

func (s *atomixStore) Store(endPoint *subapi.Subscription) error {
	if endPoint.ID == "" {
		return errors.New("ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	bytes, err := proto.Marshal(endPoint)
	if err != nil {
		return err
	}

	// Put the end-pPoint in the map using an optimistic lock if this is an update
	_, err = s.subscriptions.Put(ctx, string(endPoint.ID), bytes)

	if err != nil {
		return err
	}

	return err
}

func (s *atomixStore) Get(id subapi.ID) (*subapi.Subscription, error) {
	if id == "" {
		return nil, errors.New("ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	entry, err := s.subscriptions.Get(ctx, string(id))
	if err != nil {
		return nil, err
	}

	sub := &subapi.Subscription{}
	err = proto.Unmarshal(entry.Value, sub)
	return sub, err
}

func (s *atomixStore) Delete(id subapi.ID) error {
	if id == "" {
		return errors.New("ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := s.subscriptions.Remove(ctx, string(id))
	return err
}

func (s *atomixStore) List(ch chan<- *subapi.Subscription) error {
	mapCh := make(chan *_map.Entry)
	if err := s.subscriptions.Entries(context.Background(), mapCh); err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for entry := range mapCh {
			if endPoint, err := decodeObject(entry); err == nil {
				ch <- endPoint
			}
		}
	}()
	return nil
}

func (s *atomixStore) Watch(ch chan<- *Event, opts ...WatchOption) error {
	watchOpts := make([]_map.WatchOption, 0)
	for _, opt := range opts {
		watchOpts = opt.apply(watchOpts)
	}

	mapCh := make(chan *_map.Event)
	if err := s.subscriptions.Watch(context.Background(), mapCh, watchOpts...); err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for event := range mapCh {
			if endPoint, err := decodeObject(event.Entry); err == nil {
				ch <- &Event{
					Type:   EventType(event.Type),
					Object: endPoint,
				}
			}
		}
	}()
	return nil
}

func (s *atomixStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = s.subscriptions.Close(ctx)
	cancel()
	if s.closer != nil {
		return s.closer()
	}
	return nil
}

func decodeObject(entry *_map.Entry) (*subapi.Subscription, error) {
	endPoint := &subapi.Subscription{}
	if err := proto.Unmarshal(entry.Value, endPoint); err != nil {
		return nil, err
	}
	endPoint.ID = subapi.ID(entry.Key)
	return endPoint, nil
}

// EventType provides the type for a subscription event
type EventType string

const (
	// EventNone is no event
	EventNone EventType = ""
	// EventInserted is inserted
	EventInserted EventType = "inserted"
	// EventUpdated is updated
	EventUpdated EventType = "updated"
	// EventRemoved is removed
	EventRemoved EventType = "removed"
)

// Event is a store event for a subscription
type Event struct {
	Type   EventType
	Object *subapi.Subscription
}
