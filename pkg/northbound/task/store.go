// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package task

import (
	"context"
	"errors"
	"github.com/atomix/go-client/pkg/client/util/net"
	"io"
	"time"

	"github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/gogo/protobuf/proto"
	taskapi "github.com/onosproject/onos-e2sub/api/e2/task/v1beta1"
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

	tasks, err := database.GetMap(context.Background(), "subscription-tasks")
	if err != nil {
		return nil, err
	}

	return &atomixStore{
		tasks: tasks,
	}, nil
}

// NewLocalStore returns a new local subscription task store
func NewLocalStore() (Store, error) {
	_, address := atomix.StartLocalNode()
	return newLocalStore(address)
}

// newLocalStore creates a new local subscription task store
func newLocalStore(address net.Address) (Store, error) {
	name := primitive.Name{
		Namespace: "local",
		Name:      "subscription-tasks",
	}

	session, err := primitive.NewSession(context.TODO(), primitive.Partition{ID: 1, Address: address})
	if err != nil {
		return nil, err
	}

	tasks, err := _map.New(context.Background(), name, []*primitive.Session{session})
	if err != nil {
		return nil, err
	}

	return &atomixStore{
		tasks: tasks,
	}, nil
}

// Store stores end-point registry information
type Store interface {
	io.Closer

	// Create creates a task in the store
	Create(ctx context.Context, sub *taskapi.SubscriptionTask) error

	// Update updates a task in the store
	Update(ctx context.Context, sub *taskapi.SubscriptionTask) error

	// Delete deletes an task from the store
	Get(ctx context.Context, id taskapi.ID) (*taskapi.SubscriptionTask, error)

	// Delete deletes an task from the store
	Delete(ctx context.Context, id taskapi.ID) error

	// List streams tasks to the given channel
	List(ctx context.Context) ([]taskapi.SubscriptionTask, error)

	// Watch streams task events to the given channel
	Watch(ctx context.Context, ch chan<- taskapi.Event, opts ...WatchOption) error
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

// atomixStore is the implementation of the task Store
type atomixStore struct {
	tasks  _map.Map
	closer func() error
}

func (s *atomixStore) Create(ctx context.Context, task *taskapi.SubscriptionTask) error {
	if task.ID == "" {
		return errors.New("ID cannot be empty")
	}

	bytes, err := proto.Marshal(task)
	if err != nil {
		return err
	}

	// Create the task in the map only if it does not already exist
	entry, err := s.tasks.Put(ctx, string(task.ID), bytes, _map.IfNotSet())
	if err != nil {
		return err
	}
	task.Revision = taskapi.Revision(entry.Version)
	return err
}

func (s *atomixStore) Update(ctx context.Context, task *taskapi.SubscriptionTask) error {
	if task.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if task.Revision == 0 {
		return errors.New("object must contain a revision on update")
	}

	bytes, err := proto.Marshal(task)
	if err != nil {
		return err
	}

	// Update the task in the map
	entry, err := s.tasks.Put(ctx, string(task.ID), bytes, _map.IfVersion(_map.Version(task.Revision)))
	if err != nil {
		return err
	}
	task.Revision = taskapi.Revision(entry.Version)
	return err
}

func (s *atomixStore) Get(ctx context.Context, id taskapi.ID) (*taskapi.SubscriptionTask, error) {
	if id == "" {
		return nil, errors.New("ID cannot be empty")
	}

	entry, err := s.tasks.Get(ctx, string(id))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	return decodeObject(entry)
}

func (s *atomixStore) Delete(ctx context.Context, id taskapi.ID) error {
	if id == "" {
		return errors.New("ID cannot be empty")
	}

	_, err := s.tasks.Remove(ctx, string(id))
	return err
}

func (s *atomixStore) List(ctx context.Context) ([]taskapi.SubscriptionTask, error) {
	mapCh := make(chan *_map.Entry)
	if err := s.tasks.Entries(context.Background(), mapCh); err != nil {
		return nil, err
	}

	tasks := make([]taskapi.SubscriptionTask, 0)
	for entry := range mapCh {
		if task, err := decodeObject(entry); err == nil {
			tasks = append(tasks, *task)
		}
	}
	return tasks, nil
}

func (s *atomixStore) Watch(ctx context.Context, ch chan<- taskapi.Event, opts ...WatchOption) error {
	watchOpts := make([]_map.WatchOption, 0)
	for _, opt := range opts {
		watchOpts = opt.apply(watchOpts)
	}

	mapCh := make(chan *_map.Event)
	if err := s.tasks.Watch(context.Background(), mapCh, watchOpts...); err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for event := range mapCh {
			if task, err := decodeObject(event.Entry); err == nil {
				var eventType taskapi.EventType
				switch event.Type {
				case _map.EventNone:
					eventType = taskapi.EventType_NONE
				case _map.EventInserted:
					eventType = taskapi.EventType_CREATED
				case _map.EventUpdated:
					eventType = taskapi.EventType_UPDATED
				case _map.EventRemoved:
					eventType = taskapi.EventType_REMOVED
				}
				ch <- taskapi.Event{
					Type: eventType,
					Task: *task,
				}
			}
		}
	}()
	return nil
}

func (s *atomixStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = s.tasks.Close(ctx)
	cancel()
	if s.closer != nil {
		return s.closer()
	}
	return nil
}

func decodeObject(entry *_map.Entry) (*taskapi.SubscriptionTask, error) {
	task := &taskapi.SubscriptionTask{}
	if err := proto.Unmarshal(entry.Value, task); err != nil {
		return nil, err
	}
	task.ID = taskapi.ID(entry.Key)
	task.Revision = taskapi.Revision(entry.Version)
	return task, nil
}
