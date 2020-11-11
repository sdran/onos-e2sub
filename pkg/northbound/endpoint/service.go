// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package endpoint

import (
	"context"
	regapi "github.com/onosproject/onos-e2sub/api/e2/endpoint/v1beta1"
	store "github.com/onosproject/onos-e2sub/pkg/store/endpoint"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("northbound", "ricapi", "registry")

// NewService creates a new registry service
func NewService(store store.Store) northbound.Service {
	return &Service{
		store: store,
	}
}

// Service is a Service implementation for subscription service.
type Service struct {
	store store.Store
}

// Register registers the Service with the gRPC server.
func (s *Service) Register(r *grpc.Server) {
	server := &Server{
		endPointStore: s.store,
	}
	regapi.RegisterE2RegistryServiceServer(r, server)
}

var _ northbound.Service = &Service{}

// Server implements the gRPC service for managing of subscriptions
type Server struct {
	endPointStore store.Store
}

// E2RegistryClientFactory : Default E2RegistryClientFactory creation.
var E2RegistryClientFactory = func(cc *grpc.ClientConn) regapi.E2RegistryServiceClient {
	return regapi.NewE2RegistryServiceClient(cc)
}

// CreateE2RegistryClient creates and returns a new topo entity client
func CreateE2RegistryClient(cc *grpc.ClientConn) regapi.E2RegistryServiceClient {
	return E2RegistryClientFactory(cc)
}

// AddTermination adds an E2 end-point
func (s *Server) AddTermination(ctx context.Context, req *regapi.AddTerminationRequest) (*regapi.AddTerminationResponse, error) {
	log.Debugf("Received AddTerminationRequest %+v", req)
	ep := req.Endpoint
	err := s.endPointStore.Create(ctx, ep)
	if err != nil {
		log.Warnf("AddTerminationRequest %+v failed: %v", req, err)
		return nil, err
	}
	res := &regapi.AddTerminationResponse{}
	log.Debugf("Sending AddTerminationResponse %+v", res)
	return res, nil
}

// GetTermination retrieves information about a specific termination end-point
func (s *Server) GetTermination(ctx context.Context, req *regapi.GetTerminationRequest) (*regapi.GetTerminationResponse, error) {
	log.Debugf("Received GetSubscriptionRequest %+v", req)
	ep, err := s.endPointStore.Get(ctx, req.ID)
	if err != nil {
		log.Warnf("GetTerminatonRequest %+v failed: %v", req, err)
		return nil, err
	}
	res := &regapi.GetTerminationResponse{
		Endpoint: ep,
	}
	log.Debugf("Sending GetTerminationResponse %+v", res)
	return res, nil
}

// RemoveTermination removes a subscription
func (s *Server) RemoveTermination(ctx context.Context, req *regapi.RemoveTerminationRequest) (*regapi.RemoveTerminationResponse, error) {
	log.Debugf("Received RemoveTerminationRequest %+v", req)
	err := s.endPointStore.Delete(ctx, req.ID)
	if err != nil {
		log.Warnf("RemoveTerminationRequest %+v failed: %v", req, err)
		return nil, err
	}
	res := &regapi.RemoveTerminationResponse{}
	log.Debugf("Sending RemoveTerminationResponse %+v", res)
	return res, nil
}

// ListTerminations returns the list of current existing termination end-points
func (s *Server) ListTerminations(ctx context.Context, req *regapi.ListTerminationsRequest) (*regapi.ListTerminationsResponse, error) {
	log.Debugf("Received ListTerminationsRequest %+v", req)
	eps, err := s.endPointStore.List(ctx)
	if err != nil {
		log.Warnf("ListTerminationsRequest %+v failed: %v", req, err)
		return nil, err
	}

	res := &regapi.ListTerminationsResponse{
		Endpoints: eps,
	}
	log.Debugf("Sending ListTerminationsResponse %+v", res)
	return res, nil
}

// WatchTerminations streams termination end-point changes
func (s *Server) WatchTerminations(req *regapi.WatchTerminationsRequest, server regapi.E2RegistryService_WatchTerminationsServer) error {
	log.Debugf("Received WatchTerminationsRequest %+v", req)
	var watchOpts []store.WatchOption
	if !req.Noreplay {
		watchOpts = append(watchOpts, store.WithReplay())
	}

	ch := make(chan regapi.Event)
	if err := s.endPointStore.Watch(server.Context(), ch, watchOpts...); err != nil {
		log.Warnf("WatchTerminationsRequest %+v failed: %v", req, err)
		return err
	}

	return s.Stream(server, ch)
}

// Stream is the ongoing stream for WatchTerminations request
func (s *Server) Stream(server regapi.E2RegistryService_WatchTerminationsServer, ch chan regapi.Event) error {
	for event := range ch {
		res := &regapi.WatchTerminationsResponse{
			Event: event,
		}

		log.Debugf("Sending WatchTerminationsResponse %+v", res)
		if err := server.Send(res); err != nil {
			log.Warnf("WatchTerminationsResponse %+v failed: %v", res, err)
			return err
		}
	}
	return nil
}
