// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package task

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
)

import (
	"context"
	taskapi "github.com/onosproject/onos-e2sub/api/e2/task/v1beta1"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("northbound", "task")

// NewService creates a new subscription service
func NewService(store Store) northbound.Service {
	return &Service{
		store: store,
	}
}

// Service is a Service implementation for subscription service.
type Service struct {
	store Store
}

// Register registers the Service with the gRPC server.
func (s *Service) Register(r *grpc.Server) {
	server := &Server{
		store: s.store,
	}
	taskapi.RegisterE2SubscriptionTaskServiceServer(r, server)
}

var _ northbound.Service = &Service{}

// Server implements the gRPC service for managing of subscriptions
type Server struct {
	store Store
}

func (s *Server) GetSubscriptionTask(ctx context.Context, req *taskapi.GetSubscriptionTaskRequest) (*taskapi.GetSubscriptionTaskResponse, error) {
	log.Debugf("Received GetSubscriptionTaskRequest %+v", req)
	task, err := s.store.Get(ctx, req.ID)
	if err != nil {
		log.Warnf("GetSubscriptionTaskRequest %+v failed: %v", req, err)
		return nil, err
	}
	res := &taskapi.GetSubscriptionTaskResponse{
		Task: task,
	}
	log.Debugf("Sending GetSubscriptionTaskResponse %+v", res)
	return res, nil
}

func (s *Server) ListSubscriptionTasks(ctx context.Context, req *taskapi.ListSubscriptionTasksRequest) (*taskapi.ListSubscriptionTasksResponse, error) {
	log.Debugf("Received ListSubscriptionTasksRequest %+v", req)
	tasks, err := s.store.List(ctx)
	if err != nil {
		log.Warnf("ListSubscriptionTasksRequest %+v failed: %v", req, err)
		return nil, err
	}

	res := &taskapi.ListSubscriptionTasksResponse{
		Task: tasks,
	}
	log.Debugf("Sending ListSubscriptionTasksResponse %+v", res)
	return res, nil
}

func (s *Server) WatchSubscriptionTasks(req *taskapi.WatchSubscriptionTasksRequest, server taskapi.E2SubscriptionTaskService_WatchSubscriptionTasksServer) error {
	log.Debugf("Received WatchSubscriptionTasksRequest %+v", req)
	var watchOpts []WatchOption
	if !req.Noreplay {
		watchOpts = append(watchOpts, WithReplay())
	}

	ch := make(chan taskapi.Event)
	if err := s.store.Watch(server.Context(), ch, watchOpts...); err != nil {
		log.Warnf("WatchSubscriptionTasksRequest %+v failed: %v", req, err)
		return err
	}

	for event := range ch {
		res := &taskapi.WatchSubscriptionTasksResponse{
			Event: event,
		}

		log.Debugf("Sending WatchSubscriptionTasksResponse %+v", res)
		if err := server.Send(res); err != nil {
			log.Warnf("WatchSubscriptionTasksResponse %+v failed: %v", res, err)
			return err
		}
	}
	return nil
}

func (s *Server) UpdateSubscriptionTask(ctx context.Context, req *taskapi.UpdateSubscriptionTaskRequest) (*taskapi.UpdateSubscriptionTaskResponse, error) {
	log.Debugf("Received UpdateSubscriptionTaskRequest %+v", req)
	err := s.store.Update(ctx, req.Task)
	if err != nil {
		log.Warnf("UpdateSubscriptionTaskRequest %+v failed: %v", req, err)
		return nil, err
	}
	res := &taskapi.UpdateSubscriptionTaskResponse{}
	log.Debugf("Sending UpdateSubscriptionTaskResponse %+v", res)
	return res, nil
}
