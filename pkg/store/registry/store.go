// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package registry

import (
	"context"
	"errors"
	"io"
	"time"

	_map "github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/gogo/protobuf/proto"
	regapi "github.com/onosproject/onos-e2sub/api/e2/registry/v1beta1"
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

	endPoints, err := database.GetMap(context.Background(), "endPoints")
	if err != nil {
		return nil, err
	}

	return &atomixStore{
		endPoints: endPoints,
	}, nil
}

// NewLocalStore returns a new local end-point store
func NewLocalStore() (Store, error) {
	node, address := atomix.StartLocalNode()
	name := primitive.Name{
		Namespace: "local",
		Name:      "endPoints",
	}

	session, err := primitive.NewSession(context.TODO(), primitive.Partition{ID: 1, Address: address})
	if err != nil {
		return nil, err
	}

	endPoints, err := _map.New(context.Background(), name, []*primitive.Session{session})
	if err != nil {
		return nil, err
	}

	return &atomixStore{
		endPoints: endPoints,
		closer:    node.Stop,
	}, nil
}

// Store stores end-point registry information
type Store interface {
	io.Closer

	// Store stores an end-point in the store
	Store(ctx context.Context, point *regapi.TerminationEndPoint) error

	// Gets an end-point from the store
	Get(ctx context.Context, id regapi.ID) (*regapi.TerminationEndPoint, error)

	// Delete deletes an end-point from the store
	Delete(ctx context.Context, id regapi.ID) error

	// List streams end-points to the given channel
	List(ctx context.Context, ch chan<- *regapi.TerminationEndPoint) error

	// Watch streams end-point events to the given channel
	Watch(ctx context.Context, ch chan<- regapi.Event, opts ...WatchOption) error
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

// atomixStore is the implementation of the end-point Store
type atomixStore struct {
	endPoints _map.Map
	closer    func() error
}

func (s *atomixStore) Store(ctx context.Context, endPoint *regapi.TerminationEndPoint) error {
	if endPoint.ID == "" {
		return errors.New("ID cannot be empty")
	}

	bytes, err := proto.Marshal(endPoint)
	if err != nil {
		return err
	}

	// Put the end-point in the map using an optimistic lock if this is an update
	_, err = s.endPoints.Put(ctx, string(endPoint.ID), bytes)

	if err != nil {
		return err
	}

	return err
}

func (s *atomixStore) Get(ctx context.Context, id regapi.ID) (*regapi.TerminationEndPoint, error) {
	if id == "" {
		return nil, errors.New("ID cannot be empty")
	}
	entry, err := s.endPoints.Get(ctx, string(id))
	if err != nil {
		return nil, err
	}

	ep := &regapi.TerminationEndPoint{}
	err = proto.Unmarshal(entry.Value, ep)
	return ep, err
}

func (s *atomixStore) Delete(ctx context.Context, id regapi.ID) error {
	if id == "" {
		return errors.New("ID cannot be empty")
	}
	_, err := s.endPoints.Remove(ctx, string(id))
	return err
}

func (s *atomixStore) List(ctx context.Context, ch chan<- *regapi.TerminationEndPoint) error {
	mapCh := make(chan *_map.Entry)
	if err := s.endPoints.Entries(context.Background(), mapCh); err != nil {
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

func (s *atomixStore) Watch(ctx context.Context, ch chan<- regapi.Event, opts ...WatchOption) error {
	watchOpts := make([]_map.WatchOption, 0)
	for _, opt := range opts {
		watchOpts = opt.apply(watchOpts)
	}

	mapCh := make(chan *_map.Event)
	if err := s.endPoints.Watch(ctx, mapCh, watchOpts...); err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for event := range mapCh {
			if endPoint, err := decodeObject(event.Entry); err == nil {
				var eventType regapi.EventType
				switch event.Type {
				case _map.EventNone:
					eventType = regapi.EventType_NONE
				case _map.EventInserted:
					eventType = regapi.EventType_ADDED
				case _map.EventRemoved:
					eventType = regapi.EventType_REMOVED
				}
				ch <- regapi.Event{
					Type:     eventType,
					EndPoint: *endPoint,
				}
			}
		}
	}()
	return nil
}

func (s *atomixStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = s.endPoints.Close(ctx)
	cancel()
	if s.closer != nil {
		return s.closer()
	}
	return nil
}

func decodeObject(entry *_map.Entry) (*regapi.TerminationEndPoint, error) {
	endPoint := &regapi.TerminationEndPoint{}
	if err := proto.Unmarshal(entry.Value, endPoint); err != nil {
		return nil, err
	}
	endPoint.ID = regapi.ID(entry.Key)
	return endPoint, nil
}

// EventType provides the type for a end-point event
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

// Event is a store event for a end-point
type Event struct {
	Type   EventType
	Object *regapi.TerminationEndPoint
}
