// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package task

import (
	"context"
	"testing"
	"time"

	taskapi "github.com/onosproject/onos-api/go/onos/e2sub/task"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptionStore(t *testing.T) {
	_, address := atomix.StartLocalNode()

	store1, err := newLocalStore(address)
	assert.NoError(t, err)
	defer store1.Close()

	store2, err := newLocalStore(address)
	assert.NoError(t, err)
	defer store2.Close()

	ch := make(chan taskapi.Event)
	err = store2.Watch(context.Background(), ch)
	assert.NoError(t, err)

	task1 := &taskapi.SubscriptionTask{
		ID: "task-1",
	}
	task2 := &taskapi.SubscriptionTask{
		ID: "task-2",
	}

	// Create a new task
	err = store1.Create(context.TODO(), task1)
	assert.NoError(t, err)
	assert.Equal(t, taskapi.ID("task-1"), task1.ID)
	assert.NotEqual(t, taskapi.Revision(0), task1.Revision)

	// Get the task
	task1, err = store2.Get(context.TODO(), "task-1")
	assert.NoError(t, err)
	assert.NotNil(t, task1)
	assert.Equal(t, taskapi.ID("task-1"), task1.ID)
	assert.NotEqual(t, taskapi.Revision(0), task1.Revision)

	// Create another task
	err = store2.Create(context.TODO(), task2)
	assert.NoError(t, err)
	assert.Equal(t, taskapi.ID("task-2"), task2.ID)
	assert.NotEqual(t, taskapi.Revision(0), task2.Revision)

	// Verify events were received for the tasks
	taskEvent := nextEvent(t, ch)
	assert.Equal(t, taskapi.ID("task-1"), taskEvent.ID)
	taskEvent = nextEvent(t, ch)
	assert.Equal(t, taskapi.ID("task-2"), taskEvent.ID)

	// Update one of the tasks
	task2.Lifecycle.Status = taskapi.Status_COMPLETE
	revision := task2.Revision
	err = store1.Update(context.TODO(), task2)
	assert.NoError(t, err)
	assert.NotEqual(t, revision, task2.Revision)

	// Read and then update the task
	task2, err = store2.Get(context.TODO(), "task-2")
	assert.NoError(t, err)
	assert.NotNil(t, task2)
	task2.Lifecycle.Phase = taskapi.Phase_CLOSE
	task2.Lifecycle.Status = taskapi.Status_PENDING
	revision = task2.Revision
	err = store1.Update(context.TODO(), task2)
	assert.NoError(t, err)
	assert.NotEqual(t, revision, task2.Revision)

	// Verify that concurrent updates fail
	task11, err := store1.Get(context.TODO(), "task-1")
	assert.NoError(t, err)
	task12, err := store2.Get(context.TODO(), "task-1")
	assert.NoError(t, err)

	task11.Lifecycle.Phase = taskapi.Phase_CLOSE
	task11.Lifecycle.Status = taskapi.Status_PENDING
	err = store1.Update(context.TODO(), task11)
	assert.NoError(t, err)

	task12.Lifecycle.Phase = taskapi.Phase_CLOSE
	task12.Lifecycle.Status = taskapi.Status_PENDING
	err = store2.Update(context.TODO(), task12)
	assert.Error(t, err)

	// Verify events were received again
	taskEvent = nextEvent(t, ch)
	assert.Equal(t, taskapi.ID("task-2"), taskEvent.ID)
	taskEvent = nextEvent(t, ch)
	assert.Equal(t, taskapi.ID("task-2"), taskEvent.ID)
	taskEvent = nextEvent(t, ch)
	assert.Equal(t, taskapi.ID("task-1"), taskEvent.ID)

	// List the tasks
	tasks, err := store1.List(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, tasks, 2)

	// Delete a task
	err = store1.Delete(context.TODO(), task2.ID)
	assert.NoError(t, err)
	task2, err = store2.Get(context.TODO(), "task-2")
	assert.Error(t, err)
	assert.True(t, errors.IsNotFound(err))
	assert.Nil(t, task2)

	task := &taskapi.SubscriptionTask{
		ID: "task-1",
	}

	err = store1.Create(context.TODO(), task)
	assert.Error(t, err)

	task = &taskapi.SubscriptionTask{
		ID: "task-2",
	}

	err = store1.Create(context.TODO(), task)
	assert.NoError(t, err)

	ch = make(chan taskapi.Event)
	err = store1.Watch(context.TODO(), ch, WithReplay())
	assert.NoError(t, err)

	task = nextEvent(t, ch)
	assert.NotNil(t, task)
	task = nextEvent(t, ch)
	assert.NotNil(t, task)
}

func nextEvent(t *testing.T, ch chan taskapi.Event) *taskapi.SubscriptionTask {
	select {
	case c := <-ch:
		return &c.Task
	case <-time.After(5 * time.Second):
		t.FailNow()
	}
	return nil
}
