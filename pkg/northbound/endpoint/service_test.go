// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package endpoint

import (
	"context"
	"net"
	"sync"
	"testing"

	regapi "github.com/onosproject/onos-api/go/onos/e2sub/endpoint"
	store "github.com/onosproject/onos-e2sub/pkg/store/endpoint"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func newTestService() (northbound.Service, error) {
	endPointStore, err := store.NewLocalStore()
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
	client := regapi.NewE2RegistryServiceClient(conn)

	_, err := client.AddTermination(context.Background(), &regapi.AddTerminationRequest{
		Endpoint: &regapi.TerminationEndpoint{
			ID: "1", IP: "10.10.10.1", Port: 111,
		},
	})
	assert.NoError(t, err)

	_, err = client.AddTermination(context.Background(), &regapi.AddTerminationRequest{
		Endpoint: &regapi.TerminationEndpoint{
			ID: "2", IP: "10.10.10.2", Port: 222,
		},
	})
	assert.NoError(t, err)

	res, err := client.ListTerminations(context.Background(), &regapi.ListTerminationsRequest{})
	assert.NoError(t, err)
	assert.Condition(t, func() bool {
		return len(res.Endpoints) == 2 &&
			(res.Endpoints[0].ID == regapi.ID("1") || res.Endpoints[1].ID == regapi.ID("1"))
	})

	_, err = client.RemoveTermination(context.Background(), &regapi.RemoveTerminationRequest{
		ID: "1",
	})
	assert.NoError(t, err)

	res, err = client.ListTerminations(context.Background(), &regapi.ListTerminationsRequest{})
	assert.NoError(t, err)
	assert.Condition(t, func() bool {
		return len(res.Endpoints) == 1 && res.Endpoints[0].ID == regapi.ID("2")
	})
}

func TestWatchBasics(t *testing.T) {
	conn := createServerConnection(t)
	client := regapi.NewE2RegistryServiceClient(conn)

	_, err := client.AddTermination(context.Background(), &regapi.AddTerminationRequest{
		Endpoint: &regapi.TerminationEndpoint{
			ID: "1", IP: "10.10.10.1", Port: 111,
		},
	})
	assert.NoError(t, err)

	res, err := client.WatchTerminations(context.Background(), &regapi.WatchTerminationsRequest{})
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)

	var pause sync.WaitGroup
	pause.Add(1)
	go func() {
		wr, err := res.Recv()
		assert.NoError(t, err)
		assert.Equal(t, regapi.EventType_NONE, wr.Event.Type)
		assert.Equal(t, regapi.ID("1"), wr.Event.Endpoint.ID)
		pause.Done()

		wr, err = res.Recv()
		assert.NoError(t, err)
		assert.Equal(t, regapi.EventType_ADDED, wr.Event.Type)
		assert.Equal(t, regapi.ID("2"), wr.Event.Endpoint.ID)

		wr, err = res.Recv()
		assert.NoError(t, err)
		assert.Equal(t, regapi.EventType_REMOVED, wr.Event.Type)
		assert.Equal(t, regapi.ID("1"), wr.Event.Endpoint.ID)

		wg.Done()
	}()

	// Pause before adding a new item to validate that existing items are processed first
	pause.Wait()
	_, err = client.AddTermination(context.Background(), &regapi.AddTerminationRequest{
		Endpoint: &regapi.TerminationEndpoint{
			ID: "2", IP: "10.10.10.2", Port: 222,
		},
	})
	assert.NoError(t, err)

	_, err = client.RemoveTermination(context.Background(), &regapi.RemoveTerminationRequest{
		ID: "1",
	})
	assert.NoError(t, err)

	wg.Wait()
}

func TestBadAdd(t *testing.T) {
	conn := createServerConnection(t)
	client := regapi.NewE2RegistryServiceClient(conn)

	_, err := client.AddTermination(context.Background(), &regapi.AddTerminationRequest{
		Endpoint: &regapi.TerminationEndpoint{},
	})
	assert.Error(t, err)
}

func TestBadRemove(t *testing.T) {
	conn := createServerConnection(t)
	client := regapi.NewE2RegistryServiceClient(conn)

	_, err := client.RemoveTermination(context.Background(), &regapi.RemoveTerminationRequest{})
	assert.Error(t, err)
}
