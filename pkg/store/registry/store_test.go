// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package registry

import (
	"context"
	regapi "github.com/onosproject/onos-e2sub/api/e2/registry/v1beta1"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func validate(t *testing.T, ep regapi.TerminationEndPoint, id string, ip string, port uint32) {
	assert.Equal(t, regapi.ID(id), ep.ID)
	assert.Equal(t, regapi.IP(ip), ep.IP)
	assert.Equal(t, regapi.Port(port), ep.Port)
}

func TestStoreBasics(t *testing.T) {
	store, _ := NewLocalStore()
	ctx := context.Background()

	assert.NoError(t, store.Store(ctx, &regapi.TerminationEndPoint{ID: "1", IP: "10.10.10.1", Port: 111}))

	ep, err := store.Get(ctx, "1")
	assert.NoError(t, err)
	validate(t, *ep, "1", "10.10.10.1", 111)

	assert.NoError(t, store.Store(ctx, &regapi.TerminationEndPoint{ID: "2", IP: "10.10.10.2", Port: 222}))
	assert.NoError(t, store.Store(ctx, &regapi.TerminationEndPoint{ID: "3", IP: "10.10.10.3", Port: 333}))

	ch := make(chan *regapi.TerminationEndPoint)
	assert.NoError(t, store.List(ctx, ch))

	count := 0
	for range ch {
		count = count + 1
	}

	assert.NoError(t, store.Delete(ctx, "3"))
	assert.NoError(t, store.Delete(ctx, "1"))

	ch = make(chan *regapi.TerminationEndPoint)
	assert.NoError(t, store.List(ctx, ch))
	count = 0
	for ep := range ch {
		validate(t, *ep, "2", "10.10.10.2", 222)
		count = count + 1
	}
	assert.Equal(t, 1, count)
	assert.NoError(t, store.Close())
}

func TestStoreWatch(t *testing.T) {
	store, _ := NewLocalStore()
	ctx := context.Background()

	assert.NoError(t, store.Store(ctx, &regapi.TerminationEndPoint{ID: "1", IP: "10.10.10.1", Port: 111}))

	ch := make(chan regapi.Event)
	assert.NoError(t, store.Watch(ctx, ch, WithReplay()))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		e := <-ch
		assert.Equal(t, regapi.EventType_NONE, e.Type)
		validate(t, e.EndPoint, "1", "10.10.10.1", 111)

		e = <-ch
		assert.Equal(t, regapi.EventType_ADDED, e.Type)
		validate(t, e.EndPoint, "2", "10.10.10.2", 222)

		e = <-ch
		assert.Equal(t, regapi.EventType_REMOVED, e.Type)
		validate(t, e.EndPoint, "1", "10.10.10.1", 111)

		e = <-ch
		assert.Equal(t, regapi.EventType_ADDED, e.Type)
		validate(t, e.EndPoint, "3", "10.10.10.3", 333)

		e = <-ch
		assert.Equal(t, regapi.EventType_REMOVED, e.Type)
		validate(t, e.EndPoint, "2", "10.10.10.2", 222)

		wg.Done()
		close(ch)
	}()

	assert.NoError(t, store.Store(ctx, &regapi.TerminationEndPoint{ID: "2", IP: "10.10.10.2", Port: 222}))
	assert.NoError(t, store.Delete(ctx, "1"))
	assert.NoError(t, store.Store(ctx, &regapi.TerminationEndPoint{ID: "3", IP: "10.10.10.3", Port: 333}))
	assert.NoError(t, store.Delete(ctx, "2"))
	wg.Wait()
}
