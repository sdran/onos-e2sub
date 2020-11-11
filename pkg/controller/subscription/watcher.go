// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package subscription

import (
	"context"
	regapi "github.com/onosproject/onos-e2sub/api/e2/registry/v1beta1"
	subapi "github.com/onosproject/onos-e2sub/api/e2/subscription/v1beta1"
	taskapi "github.com/onosproject/onos-e2sub/api/e2/task/v1beta1"
	"github.com/onosproject/onos-e2sub/pkg/store/registry"
	"github.com/onosproject/onos-e2sub/pkg/store/subscription"
	"github.com/onosproject/onos-e2sub/pkg/store/task"
	"github.com/onosproject/onos-lib-go/pkg/controller"
	"sync"
)

const queueSize = 100

// Watcher is a subscription watcher
type Watcher struct {
	subs   subscription.Store
	cancel context.CancelFunc
	mu     sync.Mutex
}

// Start starts the subscription watcher
func (w *Watcher) Start(ch chan<- controller.ID) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil {
		return nil
	}

	subCh := make(chan subapi.Event, queueSize)
	ctx, cancel := context.WithCancel(context.Background())
	err := w.subs.Watch(ctx, subCh)
	if err != nil {
		cancel()
		return err
	}
	w.cancel = cancel

	go func() {
		for request := range subCh {
			ch <- controller.NewID(request.Subscription.ID)
		}
		close(ch)
	}()
	return nil
}

// Stop stops the subscription watcher
func (w *Watcher) Stop() {
	w.mu.Lock()
	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
	}
	w.mu.Unlock()
}

var _ controller.Watcher = &Watcher{}

// TerminationEndpointWatcher is a termination endpoint watcher
type TerminationEndpointWatcher struct {
	subs      subscription.Store
	endpoints registry.Store
	cancel    context.CancelFunc
	mu        sync.Mutex
}

// Start starts the channel watcher
func (w *TerminationEndpointWatcher) Start(ch chan<- controller.ID) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil {
		return nil
	}

	endpointCh := make(chan regapi.Event, queueSize)
	ctx, cancel := context.WithCancel(context.Background())
	err := w.endpoints.Watch(ctx, endpointCh)
	if err != nil {
		cancel()
		return err
	}
	w.cancel = cancel

	go func() {
		for range endpointCh {
			subs, err := w.subs.List(ctx)
			if err == nil {
				for _, sub := range subs {
					ch <- controller.NewID(sub.ID)
				}
			}
		}
		close(ch)
	}()
	return nil
}

// Stop stops the channel watcher
func (w *TerminationEndpointWatcher) Stop() {
	w.mu.Lock()
	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
	}
	w.mu.Unlock()
}

var _ controller.Watcher = &TerminationEndpointWatcher{}

// TaskWatcher is a termination endpoint watcher
type TaskWatcher struct {
	subs   subscription.Store
	tasks  task.Store
	cancel context.CancelFunc
	mu     sync.Mutex
}

// Start starts the channel watcher
func (w *TaskWatcher) Start(ch chan<- controller.ID) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil {
		return nil
	}

	taskCh := make(chan taskapi.Event, queueSize)
	ctx, cancel := context.WithCancel(context.Background())
	err := w.tasks.Watch(ctx, taskCh)
	if err != nil {
		cancel()
		return err
	}
	w.cancel = cancel

	go func() {
		for event := range taskCh {
			sub, err := w.subs.Get(ctx, event.Task.SubscriptionID)
			if err == nil {
				ch <- controller.NewID(sub.ID)
			}
		}
		close(ch)
	}()
	return nil
}

// Stop stops the channel watcher
func (w *TaskWatcher) Stop() {
	w.mu.Lock()
	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
	}
	w.mu.Unlock()
}

var _ controller.Watcher = &TaskWatcher{}
