// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package registry

import (
	"context"
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

	// Load loads an end-point from the store
	Load(endPointID regapi.ID) (*regapi.TerminationEndPoint, error)

	// Store stores an end-point in the store
	Store(point *regapi.TerminationEndPoint) error

	// Delete deletes an end-point from the store
	Delete(regapi.ID) error

	// List streams end-points to the given channel
	List(chan<- *regapi.TerminationEndPoint) error

	// Watch streams end-point events to the given channel
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

// atomixStore is the implementation of the end-point Store
type atomixStore struct {
	endPoints _map.Map
	closer    func() error
}

func (s *atomixStore) Load(endPointID regapi.ID) (*regapi.TerminationEndPoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	entry, err := s.endPoints.Get(ctx, string(endPointID))
	if err != nil {
		return nil, err
	} else if entry == nil {
		return nil, nil
	}
	return decodeObject(entry)
}

func (s *atomixStore) Store(endPoint *regapi.TerminationEndPoint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	bytes, err := proto.Marshal(endPoint)
	if err != nil {
		return err
	}

	// Put the end-pPoint in the map using an optimistic lock if this is an update
	_, err = s.endPoints.Put(ctx, string(endPoint.ID), bytes)

	if err != nil {
		return err
	}

	return err
}

func (s *atomixStore) Delete(id regapi.ID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := s.endPoints.Remove(ctx, string(id))
	return err
}

func (s *atomixStore) List(ch chan<- *regapi.TerminationEndPoint) error {
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

func (s *atomixStore) Watch(ch chan<- *Event, opts ...WatchOption) error {
	watchOpts := make([]_map.WatchOption, 0)
	for _, opt := range opts {
		watchOpts = opt.apply(watchOpts)
	}

	mapCh := make(chan *_map.Event)
	if err := s.endPoints.Watch(context.Background(), mapCh, watchOpts...); err != nil {
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
