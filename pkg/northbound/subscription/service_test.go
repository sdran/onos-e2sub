// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package subscription

import (
	"context"
	subapi "github.com/onosproject/onos-e2sub/api/e2/subscription/v1beta1"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"sync"
	"testing"
)

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func newTestService() (northbound.Service, error) {
	endPointStore, err := NewLocalStore()
	if err != nil {
		return nil, err
	}
	return &Service{
		store: endPointStore,
	}, nil
}

func createServerConnection(t *testing.T) *grpc.ClientConn {
	lis = bufconn.Listen(1024 * 1024)
	s, err := newTestService()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	server := grpc.NewServer()
	s.Register(server)

	go func() {
		if err := server.Serve(lis); err != nil {
			assert.NoError(t, err, "Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	return conn
}

func TestServiceBasics(t *testing.T) {
	conn := createServerConnection(t)
	client := subapi.NewE2SubscriptionServiceClient(conn)

	_, err := client.AddSubscription(context.Background(), &subapi.AddSubscriptionRequest{
		Subscription: &subapi.Subscription{
			ID: "1", AppID: "foo", ServiceModel: &subapi.ServiceModel{ID: "sm1"},
		},
	})
	assert.NoError(t, err)

	_, err = client.AddSubscription(context.Background(), &subapi.AddSubscriptionRequest{
		Subscription: &subapi.Subscription{
			ID: "2", AppID: "foo", ServiceModel: &subapi.ServiceModel{ID: "sm2"},
		},
	})
	assert.NoError(t, err)

	res, err := client.ListSubscriptions(context.Background(), &subapi.ListSubscriptionsRequest{})
	assert.NoError(t, err)
	assert.Condition(t, func() bool {
		return len(res.Subscriptions) == 2 &&
			(res.Subscriptions[0].ID == subapi.ID("1") || res.Subscriptions[1].ID == subapi.ID("1"))
	})

	_, err = client.RemoveSubscription(context.Background(), &subapi.RemoveSubscriptionRequest{
		ID: "1",
	})
	assert.NoError(t, err)

	res, err = client.ListSubscriptions(context.Background(), &subapi.ListSubscriptionsRequest{})
	assert.NoError(t, err)
	assert.Condition(t, func() bool {
		return len(res.Subscriptions) == 1 && res.Subscriptions[0].ID == subapi.ID("2")
	})
}

func TestWatchBasics(t *testing.T) {
	conn := createServerConnection(t)
	client := subapi.NewE2SubscriptionServiceClient(conn)

	_, err := client.AddSubscription(context.Background(), &subapi.AddSubscriptionRequest{
		Subscription: &subapi.Subscription{
			ID: "1", AppID: "foo", ServiceModel: &subapi.ServiceModel{ID: "sm1"},
		},
	})
	assert.NoError(t, err)

	res, err := client.WatchSubscriptions(context.Background(), &subapi.WatchSubscriptionsRequest{})
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)

	var pause sync.WaitGroup
	pause.Add(1)
	go func() {
		e, err := res.Recv()
		assert.NoError(t, err)
		assert.Equal(t, subapi.EventType_NONE, e.Event.Type)
		assert.Equal(t, subapi.ID("1"), e.Event.Subscription.ID)
		pause.Done()

		e, err = res.Recv()
		assert.NoError(t, err)
		assert.Equal(t, subapi.EventType_ADDED, e.Event.Type)
		assert.Equal(t, subapi.ID("2"), e.Event.Subscription.ID)

		wg.Done()
	}()

	// Pause before adding a new item to validate that existing items are processed first
	pause.Wait()
	_, err = client.AddSubscription(context.Background(), &subapi.AddSubscriptionRequest{
		Subscription: &subapi.Subscription{
			ID: "2", AppID: "foo", ServiceModel: &subapi.ServiceModel{ID: "sm2"},
		},
	})
	assert.NoError(t, err)

	_, err = client.RemoveSubscription(context.Background(), &subapi.RemoveSubscriptionRequest{
		ID: "1",
	})
	assert.NoError(t, err)

	wg.Wait()
}

func TestBadAdd(t *testing.T) {
	conn := createServerConnection(t)
	client := subapi.NewE2SubscriptionServiceClient(conn)

	_, err := client.AddSubscription(context.Background(), &subapi.AddSubscriptionRequest{
		Subscription: &subapi.Subscription{},
	})
	assert.Error(t, err)
}

func TestBadRemove(t *testing.T) {
	conn := createServerConnection(t)
	client := subapi.NewE2SubscriptionServiceClient(conn)

	_, err := client.RemoveSubscription(context.Background(), &subapi.RemoveSubscriptionRequest{})
	assert.Error(t, err)
}
