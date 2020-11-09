// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package subscription

import (
	"context"
	subapi "github.com/onosproject/onos-e2sub/api/e2/subscription/v1beta1"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("northbound", "ricapi", "subscription")

// NewService creates a new subscription service
func NewService() (northbound.Service, error) {
	subscriptionStore, err := NewAtomixStore()
	if err != nil {
		return nil, err
	}
	return &Service{
		store: subscriptionStore,
	}, nil
}

// Service is a Service implementation for subscription service.
type Service struct {
	northbound.Service
	store Store
}

// Register registers the Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	server := &Server{
		subscriptionStore: s.store,
	}
	subapi.RegisterE2SubscriptionServiceServer(r, server)
}

// Server implements the gRPC service for managing of subscriptions
type Server struct {
	subscriptionStore Store
}

// AddSubscription adds a subscription
func (s *Server) AddSubscription(ctx context.Context, req *subapi.AddSubscriptionRequest) (*subapi.AddSubscriptionResponse, error) {
	log.Debugf("Received AddSubscriptionRequest %+v", req)
	sub := req.Subscription
	err := s.subscriptionStore.Store(sub)
	if err != nil {
		log.Warnf("AddSubscriptionRequest %+v failed: %v", req, err)
		return nil, err
	}
	res := &subapi.AddSubscriptionResponse{
		Subscription: sub,
	}
	log.Debugf("Sending AddSubscriptionResponse %+v", res)
	return res, nil
}

// GetSubscription retrieves information about a specific subscription in the list of existing subscriptions
func (s *Server) GetSubscription(ctx context.Context, req *subapi.GetSubscriptionRequest) (*subapi.GetSubscriptionResponse, error) {
	log.Debugf("Received GetSubscriptionRequest %+v", req)
	sub, err := s.subscriptionStore.Get(req.ID)
	if err != nil {
		log.Warnf("GetSubscriptionRequest %+v failed: %v", req, err)
		return nil, err
	}
	res := &subapi.GetSubscriptionResponse{
		Subscription: sub,
	}
	log.Debugf("Sending GetSubscriptionResponse %+v", res)
	return res, nil
}

// RemoveSubscription removes a subscription
func (s *Server) RemoveSubscription(ctx context.Context, req *subapi.RemoveSubscriptionRequest) (*subapi.RemoveSubscriptionResponse, error) {
	log.Debugf("Received RemoveSubscriptionRequest %+v", req)
	err := s.subscriptionStore.Delete(req.ID)
	if err != nil {
		log.Warnf("RemoveSubscriptionRequest %+v failed: %v", req, err)
		return nil, err
	}
	res := &subapi.RemoveSubscriptionResponse{}
	log.Debugf("Sending RemoveSubscriptionResponse %+v", res)
	return res, nil
}

// ListSubscriptions returns the list of current existing subscriptions
func (s *Server) ListSubscriptions(ctx context.Context, req *subapi.ListSubscriptionsRequest) (*subapi.ListSubscriptionsResponse, error) {
	log.Debugf("Received ListSubscriptionsRequest %+v", req)
	ch := make(chan *subapi.Subscription)
	err := s.subscriptionStore.List(ch)
	if err != nil {
		log.Warnf("ListSubscriptionsRequest %+v failed: %v", req, err)
		return nil, err
	}

	subs := make([]subapi.Subscription, 0)
	for entry := range ch {
		subs = append(subs, *entry)
	}

	res := &subapi.ListSubscriptionsResponse{
		Subscriptions: subs,
	}
	log.Debugf("Sending ListSubscriptionsResponse %+v", res)
	return res, nil
}

// WatchSubscriptions streams subscription changes
// WatchTerminations streams termination end-point changes
func (s *Server) WatchSubscriptions(req *subapi.WatchSubscriptionsRequest, server subapi.E2SubscriptionService_WatchSubscriptionsServer) error {
	log.Debugf("Received WatchTerminationsRequest %+v", req)
	var watchOpts []WatchOption
	if !req.Noreplay {
		watchOpts = append(watchOpts, WithReplay())
	}

	ch := make(chan *Event)
	if err := s.subscriptionStore.Watch(ch, watchOpts...); err != nil {
		log.Warnf("WatchTerminationsRequest %+v failed: %v", req, err)
		return err
	}

	return s.Stream(server, ch)
}

// Stream is the ongoing stream for WatchSubscriptions request
func (s *Server) Stream(server subapi.E2SubscriptionService_WatchSubscriptionsServer, ch chan *Event) error {
	for event := range ch {
		var t subapi.EventType
		switch event.Type {
		case EventNone:
			t = subapi.EventType_NONE
		case EventInserted:
			t = subapi.EventType_ADDED
		case EventRemoved:
			t = subapi.EventType_REMOVED
		}

		res := &subapi.WatchSubscriptionsResponse{
			Type:         t,
			Subscription: *event.Object,
		}

		log.Debugf("Sending WatchSubscriptionsResponse %+v", res)
		if err := server.Send(res); err != nil {
			log.Warnf("WatchSubscriptionsResponse %+v failed: %v", res, err)
			return err
		}
	}
	return nil
}
