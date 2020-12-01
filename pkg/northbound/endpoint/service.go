// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package endpoint

import (
	"context"

	epapi "github.com/onosproject/onos-api/go/onos/e2sub/endpoint"

	store "github.com/onosproject/onos-e2sub/pkg/store/endpoint"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("northbound", "endpoint")

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
	epapi.RegisterE2RegistryServiceServer(r, server)
}

var _ northbound.Service = &Service{}

// Server implements the gRPC service for managing of subscriptions
type Server struct {
	endPointStore store.Store
}

// E2RegistryClientFactory : Default E2RegistryClientFactory creation.
var E2RegistryClientFactory = func(cc *grpc.ClientConn) epapi.E2RegistryServiceClient {
	return epapi.NewE2RegistryServiceClient(cc)
}

// CreateE2RegistryClient creates and returns a new topo entity client
func CreateE2RegistryClient(cc *grpc.ClientConn) epapi.E2RegistryServiceClient {
	return E2RegistryClientFactory(cc)
}

// AddTermination adds an E2 end-point
func (s *Server) AddTermination(ctx context.Context, req *epapi.AddTerminationRequest) (*epapi.AddTerminationResponse, error) {
	log.Infof("Received AddTerminationRequest %+v", req)
	ep := req.Endpoint
	err := s.endPointStore.Create(ctx, ep)
	if err != nil {
		log.Warnf("AddTerminationRequest %+v failed: %v", req, err)
		return nil, errors.Status(err).Err()
	}
	res := &epapi.AddTerminationResponse{}
	log.Infof("Sending AddTerminationResponse %+v", res)
	return res, nil
}

// GetTermination retrieves information about a specific termination end-point
func (s *Server) GetTermination(ctx context.Context, req *epapi.GetTerminationRequest) (*epapi.GetTerminationResponse, error) {
	log.Infof("Received GetSubscriptionRequest %+v", req)
	ep, err := s.endPointStore.Get(ctx, req.ID)
	if err != nil {
		log.Warnf("GetTerminatonRequest %+v failed: %v", req, err)
		return nil, errors.Status(err).Err()
	}
	res := &epapi.GetTerminationResponse{
		Endpoint: ep,
	}
	log.Infof("Sending GetTerminationResponse %+v", res)
	return res, nil
}

// RemoveTermination removes a subscription
func (s *Server) RemoveTermination(ctx context.Context, req *epapi.RemoveTerminationRequest) (*epapi.RemoveTerminationResponse, error) {
	log.Infof("Received RemoveTerminationRequest %+v", req)
	err := s.endPointStore.Delete(ctx, req.ID)
	if err != nil {
		log.Warnf("RemoveTerminationRequest %+v failed: %v", req, err)
		return nil, errors.Status(err).Err()
	}
	res := &epapi.RemoveTerminationResponse{}
	log.Infof("Sending RemoveTerminationResponse %+v", res)
	return res, nil
}

// ListTerminations returns the list of current existing termination end-points
func (s *Server) ListTerminations(ctx context.Context, req *epapi.ListTerminationsRequest) (*epapi.ListTerminationsResponse, error) {
	log.Infof("Received ListTerminationsRequest %+v", req)
	eps, err := s.endPointStore.List(ctx)
	if err != nil {
		log.Warnf("ListTerminationsRequest %+v failed: %v", req, err)
		return nil, errors.Status(err).Err()
	}

	res := &epapi.ListTerminationsResponse{
		Endpoints: eps,
	}
	log.Infof("Sending ListTerminationsResponse %+v", res)
	return res, nil
}

// WatchTerminations streams termination end-point changes
func (s *Server) WatchTerminations(req *epapi.WatchTerminationsRequest, server epapi.E2RegistryService_WatchTerminationsServer) error {
	log.Infof("Received WatchTerminationsRequest %+v", req)
	var watchOpts []store.WatchOption
	if !req.Noreplay {
		watchOpts = append(watchOpts, store.WithReplay())
	}

	ch := make(chan epapi.Event)
	if err := s.endPointStore.Watch(server.Context(), ch, watchOpts...); err != nil {
		log.Warnf("WatchTerminationsRequest %+v failed: %v", req, err)
		return errors.Status(err).Err()
	}

	return s.Stream(server, ch)
}

// Stream is the ongoing stream for WatchTerminations request
func (s *Server) Stream(server epapi.E2RegistryService_WatchTerminationsServer, ch chan epapi.Event) error {
	for event := range ch {
		res := &epapi.WatchTerminationsResponse{
			Event: event,
		}

		log.Infof("Sending WatchTerminationsResponse %+v", res)
		if err := server.Send(res); err != nil {
			log.Warnf("WatchTerminationsResponse %+v failed: %v", res, err)
			return err
		}
	}
	return nil
}
