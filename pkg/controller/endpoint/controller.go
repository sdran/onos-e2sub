// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package endpoint

import (
	"context"

	"time"

	epapi "github.com/onosproject/onos-api/go/onos/e2sub/endpoint"

	"github.com/onosproject/onos-e2sub/pkg/store/endpoint"
	"github.com/onosproject/onos-lib-go/pkg/controller"
	"github.com/onosproject/onos-lib-go/pkg/env"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("controller", "endpoint")

const defaultTimeout = 30 * time.Second

// NewController returns a new endpoint controller
func NewController(endpoints endpoint.Store, client *kubernetes.Clientset) *controller.Controller {
	c := controller.NewController("Endpoint")
	c.Watch(&Watcher{
		endpoints: endpoints,
	})
	c.Watch(&PodWatcher{
		client:    client,
		namespace: env.GetPodNamespace(),
		endpoints: endpoints,
	})
	c.Reconcile(&Reconciler{
		namespace: env.GetPodNamespace(),
		endpoints: endpoints,
		client:    client,
	})
	return c
}

// Reconciler is a endpoint reconciler
type Reconciler struct {
	namespace string
	endpoints endpoint.Store
	client    *kubernetes.Clientset
}

// Reconcile reconciles the state of a endpoint
func (r *Reconciler) Reconcile(id controller.ID) (controller.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Get the endpoint from the store
	endpoint, err := r.endpoints.Get(ctx, id.Value.(epapi.ID))
	if err != nil {
		if errors.IsNotFound(err) {
			return controller.Result{}, nil
		}
		return controller.Result{}, err
	}

	log.Infof("Reconciling Endpoint %+v", endpoint)

	// Get the pod associated with the endpoint
	_, err = r.client.CoreV1().Pods(r.namespace).Get(string(endpoint.ID), metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Warnf("Failed to reconcile Endpoint %+v: %s", endpoint, err)
			return controller.Result{}, err
		}

		// If the pod not exist, delete the endpoint
		log.Infof("Deleting orphaned Endpoint %+v", endpoint)
		err := r.endpoints.Delete(ctx, endpoint.ID)
		if err != nil && !errors.IsNotFound(err) {
			log.Warnf("Failed to delete orphaned Endpoint %+v: %s", endpoint, err)
			return controller.Result{}, err
		}
	}
	return controller.Result{}, nil
}
