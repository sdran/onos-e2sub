// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package subscription

import (
	"context"
	"fmt"
	"time"

	subapi "github.com/onosproject/onos-api/go/onos/e2sub/subscription"
	taskapi "github.com/onosproject/onos-api/go/onos/e2sub/task"
	"github.com/onosproject/onos-e2sub/pkg/store/endpoint"
	"github.com/onosproject/onos-e2sub/pkg/store/subscription"
	"github.com/onosproject/onos-e2sub/pkg/store/task"
	"github.com/onosproject/onos-lib-go/pkg/controller"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("controller", "subscription")

const defaultTimeout = 30 * time.Second

// NewController returns a new network controller
func NewController(subs subscription.Store, endpoints endpoint.Store, tasks task.Store) *controller.Controller {
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
	endpoints endpoint.Store
	tasks     task.Store
}

// Reconcile reconciles the state of a device change
func (r *Reconciler) Reconcile(id controller.ID) (controller.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	sub, err := r.subs.Get(ctx, id.Value.(subapi.ID))
	if err != nil {
		if errors.IsNotFound(err) {
			return controller.Result{}, nil
		}
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
	endpoints, err := r.endpoints.List(ctx)
	if err != nil {
		log.Warnf("Failed to reconcile Subscription %+v: %s", sub, err)
		return controller.Result{}, err
	}

	// Get the first termination endpoint
	if len(endpoints) == 0 {
		log.Warnf("No endpoints found for Subscription %+v", sub)
		return controller.Result{}, nil
	}

	// TODO: Use mastership to support multiple endpoints
	endpoint := &endpoints[0]

	// If a subscription task was not found, create one
	taskID := taskapi.ID(fmt.Sprintf("%s:%s", sub.ID, endpoint.ID))
	_, err = r.tasks.Get(ctx, taskID)
	if errors.IsNotFound(err) {
		log.Infof("Assigning Subscription %+v to TerminationEndpoint %+v", sub, endpoint)
		task := &taskapi.SubscriptionTask{
			ID:             taskapi.ID(fmt.Sprintf("%s:%s", sub.ID, endpoint.ID)),
			SubscriptionID: sub.ID,
			EndpointID:     endpoint.ID,
		}
		err := r.tasks.Create(ctx, task)
		if err != nil && !errors.IsAlreadyExists(err) {
			log.Warnf("Failed to assign Subscription %+v to TerminationEndpoint %+v: %s", sub, endpoint, err)
			return controller.Result{}, err
		}
	} else if err != nil {
		log.Warnf("Failed to reconcile Subscription %+v: %s", sub, err)
		return controller.Result{}, err
	}
	return controller.Result{}, nil
}

func (r *Reconciler) reconcileDeletedSubscription(sub *subapi.Subscription) (controller.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// List the subscription tasks
	tasks, err := r.tasks.List(ctx)
	if err != nil {
		log.Warnf("Failed to reconcile Subscription %+v: %s", sub, err)
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
		log.Infof("Deleting Subscription %+v", sub)
		err := r.subs.Delete(ctx, sub.ID)
		if err != nil && !errors.IsNotFound(err) {
			log.Warnf("Failed to reconcile Subscription %+v: %s", sub, err)
			return controller.Result{}, err
		}
		return controller.Result{}, nil
	}

	// Ensure all subscription tasks are marked closed and delete tasks already closed
	for _, task := range subTasks {
		if task.Lifecycle.Phase != taskapi.Phase_CLOSE {
			log.Infof("Closing SubscriptionTask %+v", task)
			task.Lifecycle.Phase = taskapi.Phase_CLOSE
			task.Lifecycle.Status = taskapi.Status_PENDING
			updateTask := task
			err := r.tasks.Update(ctx, &updateTask)
			if err != nil {
				log.Warnf("Failed to reconcile Subscription %+v: %s", sub, err)
				return controller.Result{}, err
			}
		}
		if task.Lifecycle.Phase == taskapi.Phase_CLOSE && task.Lifecycle.Status == taskapi.Status_COMPLETE {
			log.Infof("Deleting SubscriptionTask %+v", task)
			err = r.tasks.Delete(ctx, task.ID)
			if err != nil && !errors.IsNotFound(err) {
				log.Warnf("Failed to reconcile Subscription %+v: %s", sub, err)
				return controller.Result{}, err
			}
		}
	}
	return controller.Result{}, nil
}
