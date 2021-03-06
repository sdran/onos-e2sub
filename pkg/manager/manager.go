// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package manager

import (
	endpointctrl "github.com/onosproject/onos-e2sub/pkg/controller/endpoint"
	subctrl "github.com/onosproject/onos-e2sub/pkg/controller/subscription"
	"github.com/onosproject/onos-e2sub/pkg/northbound/endpoint"
	"github.com/onosproject/onos-e2sub/pkg/northbound/subscription"
	"github.com/onosproject/onos-e2sub/pkg/northbound/task"
	regstore "github.com/onosproject/onos-e2sub/pkg/store/endpoint"
	substore "github.com/onosproject/onos-e2sub/pkg/store/subscription"
	taskstore "github.com/onosproject/onos-e2sub/pkg/store/task"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var log = logging.GetLogger("manager")

// Config is a manager configuration
type Config struct {
	CAPath   string
	KeyPath  string
	CertPath string
	GRPCPort int
	E2Port   int
}

// NewManager creates a new manager
func NewManager(config Config) *Manager {
	log.Info("Creating Manager")
	return &Manager{
		Config: config,
	}
}

// Manager is a manager for the E2T service
type Manager struct {
	Config Config
}

// Run starts the manager and the associated services
func (m *Manager) Run() {
	log.Info("Running Manager")
	if err := m.Start(); err != nil {
		log.Fatal("Unable to run Manager", err)
	}
}

// Start starts the manager
func (m *Manager) Start() error {
	err := m.startNorthboundServer()
	if err != nil {
		return err
	}
	return nil
}

// startNorthboundServer starts the northbound gRPC server
func (m *Manager) startNorthboundServer() error {
	s := northbound.NewServer(northbound.NewServerCfg(
		m.Config.CAPath,
		m.Config.KeyPath,
		m.Config.CertPath,
		int16(m.Config.GRPCPort),
		true,
		northbound.SecurityConfig{}))

	endpointStore, err := regstore.NewAtomixStore()
	if err != nil {
		return err
	}

	subStore, err := substore.NewAtomixStore()
	if err != nil {
		return err
	}

	taskStore, err := taskstore.NewAtomixStore()
	if err != nil {
		return err
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	endpointController := endpointctrl.NewController(endpointStore, kubeClient)
	err = endpointController.Start()
	if err != nil {
		return err
	}

	subController := subctrl.NewController(subStore, endpointStore, taskStore)
	err = subController.Start()
	if err != nil {
		return err
	}

	s.AddService(logging.Service{})
	s.AddService(endpoint.NewService(endpointStore))
	s.AddService(subscription.NewService(subStore))
	s.AddService(task.NewService(taskStore))

	doneCh := make(chan error)
	go func() {
		err := s.Serve(func(started string) {
			log.Info("Started NBI on ", started)
			close(doneCh)
		})
		if err != nil {
			doneCh <- err
		}
	}()
	return <-doneCh
}

// Close kills the channels and manager related objects
func (m *Manager) Close() {
	log.Info("Closing Manager")
}
