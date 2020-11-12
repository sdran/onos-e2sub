// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package endpoint

import (
	"context"
	"errors"
	"github.com/atomix/go-client/pkg/client/util/net"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"io"
	"time"

	"github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/gogo/protobuf/proto"
	epapi "github.com/onosproject/onos-e2sub/api/e2/endpoint/v1beta1"
	"github.com/onosproject/onos-e2sub/pkg/config"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
)

var log = logging.GetLogger("store", "endpoint")

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

	endpoints, err := database.GetMap(context.Background(), "endpoints")
	if err != nil {
		return nil, err
	}

	return &atomixStore{
		endpoints: endpoints,
	}, nil
}

// NewLocalStore returns a new local end-point store
func NewLocalStore() (Store, error) {
	_, address := atomix.StartLocalNode()
	return newLocalStore(address)
}

// newLocalStore creates a new local subscription task store
func newLocalStore(address net.Address) (Store, error) {
	name := primitive.Name{
		Namespace: "local",
		Name:      "endpoints",
	}

	session, err := primitive.NewSession(context.TODO(), primitive.Partition{ID: 1, Address: address})
	if err != nil {
		return nil, err
	}

	endpoints, err := _map.New(context.Background(), name, []*primitive.Session{session})
	if err != nil {
		return nil, err
	}

	return &atomixStore{
		endpoints: endpoints,
	}, nil
}

// Store stores end-point registry information
type Store interface {
	io.Closer

	// Create stores an end-point in the store
	Create(ctx context.Context, point *epapi.TerminationEndpoint) error

	// Gets an end-point from the store
	Get(ctx context.Context, id epapi.ID) (*epapi.TerminationEndpoint, error)

	// Delete deletes an end-point from the store
	Delete(ctx context.Context, id epapi.ID) error

	// List streams end-points to the given channel
	List(ctx context.Context) ([]epapi.TerminationEndpoint, error)

	// Watch streams end-point events to the given channel
	Watch(ctx context.Context, ch chan<- epapi.Event, opts ...WatchOption) error
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
	endpoints _map.Map
}

func (s *atomixStore) Create(ctx context.Context, ep *epapi.TerminationEndpoint) error {
	if ep.ID == "" {
		return errors.New("ID cannot be empty")
	}

	log.Infof("Creating TerminationEndpoint %+v", ep)
	bytes, err := proto.Marshal(ep)
	if err != nil {
		log.Errorf("Failed to create TerminationEndpoint %+v: %s", ep, err)
		return err
	}

	// Put the end-point in the map using an optimistic lock if this is an update
	entry, err := s.endpoints.Put(ctx, string(ep.ID), bytes, _map.IfNotSet())
	if err != nil {
		log.Errorf("Failed to create TerminationEndpoint %+v: %s", ep, err)
		return err
	}
	ep.Revision = epapi.Revision(entry.Version)
	return err
}

func (s *atomixStore) Get(ctx context.Context, id epapi.ID) (*epapi.TerminationEndpoint, error) {
	if id == "" {
		return nil, errors.New("ID cannot be empty")
	}
	entry, err := s.endpoints.Get(ctx, string(id))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	return decodeObject(entry)
}

func (s *atomixStore) Delete(ctx context.Context, id epapi.ID) error {
	if id == "" {
		return errors.New("ID cannot be empty")
	}

	log.Infof("Deleting TerminationEndpoint %s", id)
	_, err := s.endpoints.Remove(ctx, string(id))
	if err != nil {
		log.Errorf("Failed to delete TerminationEndpoint %s: %s", id, err)
		return err
	}
	return nil
}

func (s *atomixStore) List(ctx context.Context) ([]epapi.TerminationEndpoint, error) {
	mapCh := make(chan *_map.Entry)
	if err := s.endpoints.Entries(context.Background(), mapCh); err != nil {
		return nil, err
	}

	eps := make([]epapi.TerminationEndpoint, 0)
	for entry := range mapCh {
		if ep, err := decodeObject(entry); err == nil {
			eps = append(eps, *ep)
		}
	}
	return eps, nil
}

func (s *atomixStore) Watch(ctx context.Context, ch chan<- epapi.Event, opts ...WatchOption) error {
	watchOpts := make([]_map.WatchOption, 0)
	for _, opt := range opts {
		watchOpts = opt.apply(watchOpts)
	}

	mapCh := make(chan *_map.Event)
	if err := s.endpoints.Watch(ctx, mapCh, watchOpts...); err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for event := range mapCh {
			if ep, err := decodeObject(event.Entry); err == nil {
				var eventType epapi.EventType
				switch event.Type {
				case _map.EventNone:
					eventType = epapi.EventType_NONE
				case _map.EventInserted:
					eventType = epapi.EventType_ADDED
				case _map.EventRemoved:
					eventType = epapi.EventType_REMOVED
				}
				ch <- epapi.Event{
					Type:     eventType,
					Endpoint: *ep,
				}
			}
		}
	}()
	return nil
}

func (s *atomixStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = s.endpoints.Close(ctx)
	defer cancel()
	return s.endpoints.Close(ctx)
}

func decodeObject(entry *_map.Entry) (*epapi.TerminationEndpoint, error) {
	ep := &epapi.TerminationEndpoint{}
	if err := proto.Unmarshal(entry.Value, ep); err != nil {
		return nil, err
	}
	ep.ID = epapi.ID(entry.Key)
	return ep, nil
}
