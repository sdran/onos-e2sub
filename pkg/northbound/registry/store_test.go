// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package registry

import (
	regapi "github.com/onosproject/onos-e2sub/api/e2/registry/v1beta1"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func validate(t *testing.T, ep *regapi.TerminationEndPoint, id string, ip string, port uint32) {
	assert.Equal(t, regapi.ID(id), ep.ID)
	assert.Equal(t, regapi.IP(ip), ep.IP)
	assert.Equal(t, regapi.Port(port), ep.Port)
}

func TestStoreBasics(t *testing.T) {
	store, _ := NewLocalStore()

	assert.NoError(t, store.Store(&regapi.TerminationEndPoint{ID: "1", IP: "10.10.10.1", Port: 111}))

	ch := make(chan *regapi.TerminationEndPoint)
	assert.NoError(t, store.List(ch))
	for ep := range ch {
		validate(t, ep, "1", "10.10.10.1", 111)
	}

	assert.NoError(t, store.Store(&regapi.TerminationEndPoint{ID: "2", IP: "10.10.10.2", Port: 222}))
	assert.NoError(t, store.Store(&regapi.TerminationEndPoint{ID: "3", IP: "10.10.10.3", Port: 333}))

	ch = make(chan *regapi.TerminationEndPoint)
	assert.NoError(t, store.List(ch))

	count := 0
	for range ch {
		count = count + 1
	}

	assert.NoError(t, store.Delete("3"))
	assert.NoError(t, store.Delete("1"))

	ch = make(chan *regapi.TerminationEndPoint)
	assert.NoError(t, store.List(ch))
	count = 0
	for ep := range ch {
		validate(t, ep, "2", "10.10.10.2", 222)
		count = count + 1
	}
	assert.Equal(t, 1, count)
	assert.NoError(t, store.Close())
}

func TestStoreWatch(t *testing.T) {
	store, _ := NewLocalStore()

	assert.NoError(t, store.Store(&regapi.TerminationEndPoint{ID: "1", IP: "10.10.10.1", Port: 111}))

	ch := make(chan *Event)
	assert.NoError(t, store.Watch(ch, WithReplay()))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		e := <-ch
		assert.Equal(t, EventNone, e.Type)
		validate(t, e.Object, "1", "10.10.10.1", 111)

		e = <-ch
		assert.Equal(t, EventInserted, e.Type)
		validate(t, e.Object, "2", "10.10.10.2", 222)

		e = <-ch
		assert.Equal(t, EventRemoved, e.Type)
		validate(t, e.Object, "1", "10.10.10.1", 111)

		e = <-ch
		assert.Equal(t, EventInserted, e.Type)
		validate(t, e.Object, "3", "10.10.10.3", 333)

		e = <-ch
		assert.Equal(t, EventRemoved, e.Type)
		validate(t, e.Object, "2", "10.10.10.2", 222)

		wg.Done()
		close(ch)
	}()

	assert.NoError(t, store.Store(&regapi.TerminationEndPoint{ID: "2", IP: "10.10.10.2", Port: 222}))
	assert.NoError(t, store.Delete("1"))
	assert.NoError(t, store.Store(&regapi.TerminationEndPoint{ID: "3", IP: "10.10.10.3", Port: 333}))
	assert.NoError(t, store.Delete("2"))
	wg.Wait()
}
