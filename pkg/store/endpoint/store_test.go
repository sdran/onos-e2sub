// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package endpoint

import (
	"context"
	regapi "github.com/onosproject/onos-e2sub/api/e2/endpoint/v1beta1"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEndPointStore(t *testing.T) {
	_, address := atomix.StartLocalNode()

	store1, err := newLocalStore(address)
	assert.NoError(t, err)
	defer store1.Close()

	store2, err := newLocalStore(address)
	assert.NoError(t, err)
	defer store2.Close()

	ch := make(chan regapi.Event)
	err = store2.Watch(context.Background(), ch)
	assert.NoError(t, err)

	ep1 := &regapi.TerminationEndpoint{ID: "ep1", IP: "10.10.10.1", Port: 111}
	ep2 := &regapi.TerminationEndpoint{ID: "ep2", IP: "10.10.10.2", Port: 222}

	// Create a new end-point in one store
	err = store1.Create(context.TODO(), ep1)
	assert.NoError(t, err)
	assert.NotEqual(t, regapi.Revision(0), ep1.Revision)

	// Get it from the other store
	ep, err := store2.Get(context.TODO(), "ep1")
	assert.NoError(t, err)
	assert.NotNil(t, ep)
	assert.Equal(t, ep1.ID, ep.ID)
	assert.NotEqual(t, regapi.Revision(0), ep1.Revision)

	// Create another end-point
	err = store2.Create(context.TODO(), ep2)
	assert.NoError(t, err)
	assert.NotEqual(t, regapi.Revision(0), ep2.Revision)

	// Verify events were received for the two end-points
	event := nextEvent(t, ch)
	assert.Equal(t, regapi.ID("ep1"), event.ID)
	event = nextEvent(t, ch)
	assert.Equal(t, regapi.ID("ep2"), event.ID)

	// List the end-points
	eps, err := store1.List(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, eps, 2)

	// Delete an end-point
	err = store1.Delete(context.TODO(), ep2.ID)
	assert.NoError(t, err)
	ep, err = store2.Get(context.TODO(), "ep2")
	assert.NoError(t, err)
	assert.Nil(t, ep)

	_ = &regapi.TerminationEndpoint{ID: "ep2"}
	err = store1.Create(context.TODO(), ep2)
	assert.NoError(t, err)

	ch = make(chan regapi.Event)
	err = store1.Watch(context.TODO(), ch, WithReplay())
	assert.NoError(t, err)

	ep = nextEvent(t, ch)
	assert.NotNil(t, ep)
	ep = nextEvent(t, ch)
	assert.NotNil(t, ep)
}

func nextEvent(t *testing.T, ch chan regapi.Event) *regapi.TerminationEndpoint {
	select {
	case c := <-ch:
		return &c.Endpoint
	case <-time.After(5 * time.Second):
		t.FailNow()
	}
	return nil
}
