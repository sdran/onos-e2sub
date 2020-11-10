// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package subscription

import (
	"context"
	subapi "github.com/onosproject/onos-e2sub/api/e2/subscription/v1beta1"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSubscriptionStore(t *testing.T) {
	_, address := atomix.StartLocalNode()

	store1, err := newLocalStore(address)
	assert.NoError(t, err)
	defer store1.Close()

	store2, err := newLocalStore(address)
	assert.NoError(t, err)
	defer store2.Close()

	ch := make(chan subapi.Event)
	err = store2.Watch(context.Background(), ch)
	assert.NoError(t, err)

	sub1 := &subapi.Subscription{
		ID:    "subscription-1",
		AppID: subapi.AppID(1),
	}
	sub2 := &subapi.Subscription{
		ID:    "subscription-2",
		AppID: subapi.AppID(2),
	}

	// Create a new subscription
	err = store1.Create(context.TODO(), sub1)
	assert.NoError(t, err)
	assert.Equal(t, subapi.ID("subscription-1"), sub1.ID)
	assert.NotEqual(t, subapi.Revision(0), sub1.Revision)

	// Get the subscription
	sub1, err = store2.Get(context.TODO(), "subscription-1")
	assert.NoError(t, err)
	assert.NotNil(t, sub1)
	assert.Equal(t, subapi.ID("subscription-1"), sub1.ID)
	assert.NotEqual(t, subapi.Revision(0), sub1.Revision)

	// Create another subscription
	err = store2.Create(context.TODO(), sub2)
	assert.NoError(t, err)
	assert.Equal(t, subapi.ID("subscription-2"), sub2.ID)
	assert.NotEqual(t, subapi.Revision(0), sub2.Revision)

	// Verify events were received for the subscriptions
	subscriptionEvent := nextEvent(t, ch)
	assert.Equal(t, subapi.ID("subscription-1"), subscriptionEvent.ID)
	subscriptionEvent = nextEvent(t, ch)
	assert.Equal(t, subapi.ID("subscription-2"), subscriptionEvent.ID)

	// Update one of the subscriptions
	sub2.ServiceModel = &subapi.ServiceModel{
		ID: subapi.ServiceModelID("service-model-2"),
	}
	revision := sub2.Revision
	err = store1.Update(context.TODO(), sub2)
	assert.NoError(t, err)
	assert.NotEqual(t, revision, sub2.Revision)

	// Read and then update the subscription
	sub2, err = store2.Get(context.TODO(), "subscription-2")
	assert.NoError(t, err)
	assert.NotNil(t, sub2)
	sub2.State.Status = subapi.Status_PENDING_DELETE
	revision = sub2.Revision
	err = store1.Update(context.TODO(), sub2)
	assert.NoError(t, err)
	assert.NotEqual(t, revision, sub2.Revision)

	// Verify that concurrent updates fail
	sub11, err := store1.Get(context.TODO(), "subscription-1")
	assert.NoError(t, err)
	sub12, err := store2.Get(context.TODO(), "subscription-1")
	assert.NoError(t, err)

	sub11.State.Status = subapi.Status_PENDING_DELETE
	err = store1.Update(context.TODO(), sub11)
	assert.NoError(t, err)

	sub12.State.Status = subapi.Status_PENDING_DELETE
	err = store2.Update(context.TODO(), sub12)
	assert.Error(t, err)

	// Verify events were received again
	subscriptionEvent = nextEvent(t, ch)
	assert.Equal(t, subapi.ID("subscription-2"), subscriptionEvent.ID)
	subscriptionEvent = nextEvent(t, ch)
	assert.Equal(t, subapi.ID("subscription-2"), subscriptionEvent.ID)
	subscriptionEvent = nextEvent(t, ch)
	assert.Equal(t, subapi.ID("subscription-1"), subscriptionEvent.ID)

	// List the subscriptions
	subs, err := store1.List(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, subs, 2)

	// Delete a subscription
	err = store1.Delete(context.TODO(), sub2.ID)
	assert.NoError(t, err)
	sub2, err = store2.Get(context.TODO(), "subscription-2")
	assert.NoError(t, err)
	assert.Nil(t, sub2)

	sub := &subapi.Subscription{
		ID:    "subscription-1",
		AppID: subapi.AppID(1),
	}

	err = store1.Create(context.TODO(), sub)
	assert.Error(t, err)

	sub = &subapi.Subscription{
		ID:    "subscription-2",
		AppID: subapi.AppID(2),
	}

	err = store1.Create(context.TODO(), sub)
	assert.NoError(t, err)

	ch = make(chan subapi.Event)
	err = store1.Watch(context.TODO(), ch, WithReplay())
	assert.NoError(t, err)

	sub = nextEvent(t, ch)
	assert.NotNil(t, sub)
	sub = nextEvent(t, ch)
	assert.NotNil(t, sub)
}

func nextEvent(t *testing.T, ch chan subapi.Event) *subapi.Subscription {
	select {
	case c := <-ch:
		return &c.Subscription
	case <-time.After(5 * time.Second):
		t.FailNow()
	}
	return nil
}
