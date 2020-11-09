// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package subscription

import (
	"context"
	subapi "github.com/onosproject/onos-e2sub/api/e2/subscription/v1beta1"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func validate(t *testing.T, sub *subapi.Subscription, id string, aid string, sid string) {
	assert.Equal(t, subapi.ID(id), sub.ID)
	assert.Equal(t, subapi.AppID(aid), sub.AppID)
	assert.Equal(t, subapi.ServiceModelID(sid), sub.ServiceModel.ID)
}

func TestStoreBasics(t *testing.T) {
	store, _ := NewLocalStore()
	ctx := context.Background()

	assert.NoError(t, store.Store(ctx, &subapi.Subscription{ID: "1", AppID: "foo", ServiceModel: &subapi.ServiceModel{ID: "sm1"}}))

	sub, err := store.Get(ctx, "1")
	assert.NoError(t, err)
	validate(t, sub, "1", "foo", "sm1")

	assert.NoError(t, store.Store(ctx, &subapi.Subscription{ID: "2", AppID: "foo", ServiceModel: &subapi.ServiceModel{ID: "sm2"}}))
	assert.NoError(t, store.Store(ctx, &subapi.Subscription{ID: "3", AppID: "bar", ServiceModel: &subapi.ServiceModel{ID: "sm1"}}))

	ch := make(chan *subapi.Subscription)
	assert.NoError(t, store.List(ctx, ch))

	count := 0
	for range ch {
		count = count + 1
	}

	assert.NoError(t, store.Delete(ctx, "3"))
	assert.NoError(t, store.Delete(ctx, "1"))

	sub, err = store.Get(ctx, "2")
	assert.NoError(t, err)
	validate(t, sub, "2", "foo", "sm2")
	assert.NoError(t, store.Close())
}

func TestStoreWatch(t *testing.T) {
	store, _ := NewLocalStore()
	ctx := context.Background()

	assert.NoError(t, store.Store(ctx, &subapi.Subscription{ID: "1", AppID: "foo", ServiceModel: &subapi.ServiceModel{ID: "sm1"}}))

	ch := make(chan *Event)
	assert.NoError(t, store.Watch(ctx, ch, WithReplay()))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		e := <-ch
		assert.Equal(t, EventNone, e.Type)
		validate(t, e.Object, "1", "foo", "sm1")

		e = <-ch
		assert.Equal(t, EventInserted, e.Type)
		validate(t, e.Object, "2", "foo", "sm2")

		e = <-ch
		assert.Equal(t, EventRemoved, e.Type)
		validate(t, e.Object, "1", "foo", "sm1")

		e = <-ch
		assert.Equal(t, EventInserted, e.Type)
		validate(t, e.Object, "3", "bar", "sm1")

		e = <-ch
		assert.Equal(t, EventRemoved, e.Type)
		validate(t, e.Object, "2", "foo", "sm2")

		wg.Done()
		close(ch)
	}()

	assert.NoError(t, store.Store(ctx, &subapi.Subscription{ID: "2", AppID: "foo", ServiceModel: &subapi.ServiceModel{ID: "sm2"}}))
	assert.NoError(t, store.Delete(ctx, "1"))
	assert.NoError(t, store.Store(ctx, &subapi.Subscription{ID: "3", AppID: "bar", ServiceModel: &subapi.ServiceModel{ID: "sm1"}}))
	assert.NoError(t, store.Delete(ctx, "2"))
	wg.Wait()
}
