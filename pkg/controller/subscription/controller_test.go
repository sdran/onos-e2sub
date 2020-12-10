// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package subscription

import (
	"context"
	"testing"
	"time"

	epapi "github.com/onosproject/onos-api/go/onos/e2sub/endpoint"
	subapi "github.com/onosproject/onos-api/go/onos/e2sub/subscription"
	taskapi "github.com/onosproject/onos-api/go/onos/e2sub/task"
	epstore "github.com/onosproject/onos-e2sub/pkg/store/endpoint"
	substore "github.com/onosproject/onos-e2sub/pkg/store/subscription"
	taskstore "github.com/onosproject/onos-e2sub/pkg/store/task"

	"github.com/onosproject/onos-lib-go/pkg/controller"
	"github.com/stretchr/testify/assert"
)

type testController struct {
	cntrl     *controller.Controller
	subStore  substore.Store
	epStore   epstore.Store
	taskStore taskstore.Store
}

func createController(t *testing.T) testController {
	subStore, err := substore.NewLocalStore()
	assert.NoError(t, err)

	epStore, err := epstore.NewLocalStore()
	assert.NoError(t, err)

	taskStore, err := taskstore.NewLocalStore()
	assert.NoError(t, err)

	cntrl := NewController(subStore, epStore, taskStore)
	assert.NotNil(t, cntrl)

	return testController{
		cntrl:     cntrl,
		subStore:  subStore,
		epStore:   epStore,
		taskStore: taskStore,
	}
}

func destroyController(t *testing.T, c testController) {
	c.cntrl.Stop()
	assert.NoError(t, c.subStore.Close())
	assert.NoError(t, c.epStore.Close())
	assert.NoError(t, c.taskStore.Close())
}

func checkTask(t *testing.T, task taskapi.SubscriptionTask, taskID taskapi.ID, subID subapi.ID, epID epapi.ID) {
	assert.Equal(t, taskID, task.ID)
	assert.Equal(t, subID, task.SubscriptionID)
	assert.Equal(t, epID, task.EndpointID)
}

func checkEvent(t *testing.T, event taskapi.Event, eventType taskapi.EventType, task taskapi.SubscriptionTask) {
	checkTask(t, task, event.Task.ID, event.Task.SubscriptionID, event.Task.EndpointID)
	assert.Equal(t, eventType, event.Type)
	assert.Equal(t, task.ID, event.Task.ID)
}

func createSubscription(subID string, e2ID string) subapi.Subscription {
	return subapi.Subscription{
		ID:       subapi.ID(subID),
		Revision: 1,
		AppID:    "app1",
		Details: &subapi.SubscriptionDetails{
			E2NodeID: subapi.E2NodeID(e2ID),
		},
		Lifecycle: subapi.Lifecycle{Status: subapi.Status_ACTIVE},
	}
}

func createEP(epID string) epapi.TerminationEndpoint {
	return epapi.TerminationEndpoint{
		ID:       epapi.ID(epID),
		Revision: 0,
		IP:       "127.0.0.1",
		Port:     555,
	}
}

func nextTaskEvent(t *testing.T, ch chan taskapi.Event) (taskapi.Event, taskapi.SubscriptionTask) {
	t.Helper()
	var event taskapi.Event
	var task taskapi.SubscriptionTask
	select {
	case event = <-ch:
		task = event.Task
		break
	case <-time.After(15 * time.Second):
		t.Error("Task Event channel timed out")
		break
	}
	return event, task
}

func nextSubEvent(t *testing.T, ch chan subapi.Event) subapi.Event {
	t.Helper()
	var event subapi.Event
	select {
	case event = <-ch:
		break
	case <-time.After(15 * time.Second):
		t.Error("Sub Event channel timed out")
		break
	}
	return event
}

// TestAddSubscription tests adding a new subscription and the resulting events
func TestAddSubscription(t *testing.T) {
	// Set up a controller to test with
	const (
		subID  = "sub1"
		epID   = "ep1"
		taskID = taskapi.ID(subID + ":" + epID)
	)
	c := createController(t)
	assert.NoError(t, c.cntrl.Start())

	// Make an end point and put it in the store
	ep := createEP(epID)
	assert.NoError(t, c.epStore.Create(context.Background(), &ep))

	// Make a subscription and put it in the store
	sub := createSubscription(subID, epID)
	assert.NoError(t, c.subStore.Create(context.TODO(), &sub))

	// Make a channel for task events
	ch := make(chan taskapi.Event)
	assert.NoError(t, c.taskStore.Watch(context.TODO(), ch))

	// Make sure the subscription creation made a task
	event, task := nextTaskEvent(t, ch)
	checkTask(t, task, taskID, subID, epID)
	checkEvent(t, event, taskapi.EventType_CREATED, task)

	// clean up
	close(ch)
	destroyController(t, c)
}

func TestDeleteSubscription(t *testing.T) {
	// Set up a controller to test
	const (
		subAddID = "sub2"
		epID     = "ep2"
		taskID   = taskapi.ID(subAddID + ":" + epID)
	)
	c := createController(t)
	assert.NoError(t, c.cntrl.Start())

	// Make an end point
	ep := createEP(epID)
	assert.NoError(t, c.epStore.Create(context.Background(), &ep))

	// Make a subscription and put it in the store
	subAdd := createSubscription(subAddID, epID)
	assert.NoError(t, c.subStore.Create(context.TODO(), &subAdd))

	// Watch for task events
	taskCh := make(chan taskapi.Event)
	assert.NoError(t, c.taskStore.Watch(context.TODO(), taskCh))

	// Watch for subscription events
	subCh := make(chan subapi.Event)
	assert.NoError(t, c.subStore.Watch(context.TODO(), subCh))

	// Get and check the subscription created event
	event, task := nextTaskEvent(t, taskCh)
	checkTask(t, task, taskID, subAddID, epID)
	checkEvent(t, event, taskapi.EventType_CREATED, task)

	// Update the subscription to mark it for deletion
	subAdd.Lifecycle = subapi.Lifecycle{Status: subapi.Status_PENDING_DELETE}
	assert.NoError(t, c.subStore.Update(context.TODO(), &subAdd))

	// Get and check the task updated event
	event, task = nextTaskEvent(t, taskCh)
	checkTask(t, task, taskID, subAddID, epID)
	checkEvent(t, event, taskapi.EventType_UPDATED, task)

	// Mark the task as completed
	task.Lifecycle = taskapi.Lifecycle{
		Phase:  taskapi.Phase_CLOSE,
		Status: taskapi.Status_COMPLETE,
	}
	assert.NoError(t, c.taskStore.Update(context.TODO(), &task))

	// Get and check the subscription updated event
	subCreatedEvent := nextSubEvent(t, subCh)
	assert.Equal(t, subapi.EventType_UPDATED, subCreatedEvent.Type)

	// Get and check the subscription removed event
	subCreatedEvent = nextSubEvent(t, subCh)
	assert.Equal(t, subapi.EventType_REMOVED, subCreatedEvent.Type)

	// Get and check the task update event
	event, task = nextTaskEvent(t, taskCh)
	checkTask(t, task, taskID, subAddID, epID)
	checkEvent(t, event, taskapi.EventType_UPDATED, task)

	// Get and check the task removed event
	event, task = nextTaskEvent(t, taskCh)
	checkTask(t, task, taskID, subAddID, epID)
	checkEvent(t, event, taskapi.EventType_REMOVED, task)

	// Clean up
	close(taskCh)
	close(subCh)
	destroyController(t, c)
}
