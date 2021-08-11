// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package endpoint

import (
	"context"
	"sync"

	epapi "github.com/onosproject/onos-api/go/onos/e2sub/endpoint"
	"github.com/onosproject/onos-e2sub/pkg/store/endpoint"
	"github.com/onosproject/onos-lib-go/pkg/controller"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const queueSize = 100

// Watcher is a endpoint watcher
type Watcher struct {
	endpoints endpoint.Store
	cancel    context.CancelFunc
	mu        sync.Mutex
}

// Start starts the endpoint watcher
func (w *Watcher) Start(ch chan<- controller.ID) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil {
		return nil
	}

	endpointCh := make(chan epapi.Event, queueSize)
	ctx, cancel := context.WithCancel(context.Background())
	err := w.endpoints.Watch(ctx, endpointCh)
	if err != nil {
		cancel()
		return err
	}
	w.cancel = cancel

	go func() {
		for request := range endpointCh {
			ch <- controller.NewID(request.Endpoint.ID)
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

// PodWatcher is a pod watcher
type PodWatcher struct {
	client    *kubernetes.Clientset
	namespace string
	endpoints endpoint.Store
	cancel    context.CancelFunc
	mu        sync.Mutex
}

// Start starts the pod watcher
func (w *PodWatcher) Start(ch chan<- controller.ID) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	watch, err := w.client.CoreV1().Pods(w.namespace).Watch(metav1.ListOptions{})
	if err != nil {
		cancel()
		return err
	}
	w.cancel = cancel

	go func() {
		<-ctx.Done()
		watch.Stop()
	}()

	go func() {
		for event := range watch.ResultChan() {
			pod := event.Object.(*corev1.Pod)
			endpoint, err := w.endpoints.Get(ctx, epapi.ID(pod.Name))
			if err != nil {
				log.Error(err)
			} else {
				ch <- controller.NewID(endpoint.ID)
			}
		}
		close(ch)
	}()
	return nil
}

// Stop stops the pod watcher
func (w *PodWatcher) Stop() {
	w.mu.Lock()
	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
	}
	w.mu.Unlock()
}

var _ controller.Watcher = &PodWatcher{}
