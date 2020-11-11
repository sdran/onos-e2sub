// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package subscription

import (
	"context"
	"fmt"
	endpointapi "github.com/onosproject/onos-e2sub/api/e2/registry/v1beta1"
	subapi "github.com/onosproject/onos-e2sub/api/e2/subscription/v1beta1"
	taskapi "github.com/onosproject/onos-e2sub/api/e2/task/v1beta1"
	"github.com/onosproject/onos-e2sub/pkg/northbound/registry"
	"github.com/onosproject/onos-e2sub/pkg/northbound/subscription"
	"github.com/onosproject/onos-e2sub/pkg/northbound/task"
	"github.com/onosproject/onos-lib-go/pkg/controller"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("subscription", "controller")

const defaultTimeout = 30 * time.Second

// NewController returns a new network controller
func NewController(subs subscription.Store, endpoints registry.Store, tasks task.Store) *controller.Controller {
	c := controller.NewController("Subscription")
	c.Watch(&Watcher{
		subs: subs,
	})
	c.Watch(&TerminationEndpointWatcher{
		subs:      subs,
		endpoints: endpoints,
	})
	c.Watch(&TaskWatcher{
		subs:  subs,
		tasks: tasks,
	})
	c.Reconcile(&Reconciler{
		subs:      subs,
		endpoints: endpoints,
		tasks:     tasks,
	})
	return c
}

// Reconciler is a device change reconciler
type Reconciler struct {
	subs      subscription.Store
	endpoints registry.Store
	tasks     task.Store
}

// Reconcile reconciles the state of a device change
func (r *Reconciler) Reconcile(id controller.ID) (controller.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	sub, err := r.subs.Get(ctx, id.Value.(subapi.ID))
	if err != nil {
		return controller.Result{}, err
	}

	log.Infof("Reconciling Subscription %+v", sub)

	switch sub.Lifecycle.Status {
	case subapi.Status_ACTIVE:
		return r.reconcileActiveSubscription(sub)
	case subapi.Status_PENDING_DELETE:
		return r.reconcileDeletedSubscription(sub)
	}
	return controller.Result{}, nil
}

func (r *Reconciler) reconcileActiveSubscription(sub *subapi.Subscription) (controller.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// List the termination endpoints
	endpointsCh := make(chan *endpointapi.TerminationEndPoint)
	err := r.endpoints.List(ctx, endpointsCh)
	if err != nil {
		return controller.Result{}, err
	}

	// Create a list of termination endpoints
	endpoints := make([]endpointapi.TerminationEndPoint, 0)
	for endpoint := range endpointsCh {
		endpoints = append(endpoints, *endpoint)
	}

	// Get the first termination endpoint
	if len(endpoints) == 0 {
		return controller.Result{}, nil
	}

	// TODO: Use mastership to support multiple endpoints
	endpoint := &endpoints[0]

	// If a subscription task was not found, create one
	taskID := taskapi.ID(fmt.Sprintf("%s:%s", sub.ID, endpoint.ID))
	task, err := r.tasks.Get(ctx, taskID)
	if err != nil {
		return controller.Result{}, err
	}

	if task == nil {
		task := &taskapi.SubscriptionTask{
			ID:                    taskapi.ID(fmt.Sprintf("%s:%s", sub.ID, endpoint.ID)),
			SubscriptionID:        sub.ID,
			TerminationEndpointID: endpoint.ID,
		}
		err := r.tasks.Create(ctx, task)
		if err != nil {
			return controller.Result{}, err
		}
	}
	return controller.Result{}, nil
}

func (r *Reconciler) reconcileDeletedSubscription(sub *subapi.Subscription) (controller.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// List the subscription tasks
	tasks, err := r.tasks.List(ctx)
	if err != nil {
		return controller.Result{}, err
	}

	// Filter the subscription tasks by this subscription
	subTasks := make([]taskapi.SubscriptionTask, 0, len(tasks))
	for _, task := range tasks {
		if task.SubscriptionID == sub.ID {
			subTasks = append(subTasks, task)
		}
	}

	// If the subscription tasks are empty, delete the subscription
	if len(subTasks) == 0 {
		err := r.subs.Delete(ctx, sub.ID)
		if err != nil {
			return controller.Result{}, err
		}
		return controller.Result{}, nil
	}

	// Ensure all subscription tasks are marked closed and delete tasks already closed
	for _, task := range subTasks {
		if task.Lifecycle.Phase != taskapi.Phase_CLOSE {
			task.Lifecycle.Phase = taskapi.Phase_CLOSE
			task.Lifecycle.Status = taskapi.Status_PENDING
			err := r.tasks.Update(ctx, &task)
			if err != nil {
				return controller.Result{}, err
			}
		}
		if task.Lifecycle.Phase == taskapi.Phase_CLOSE && task.Lifecycle.Status == taskapi.Status_COMPLETE {
			err = r.tasks.Delete(ctx, task.ID)
			if err != nil {
				return controller.Result{}, err
			}
		}
	}
	return controller.Result{}, nil
}