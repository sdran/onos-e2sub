// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package main

import (
	"flag"
	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("main")

const probeFile = "/tmp/healthy"

func main() {
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")

	flag.Parse()

	opts, err := certs.HandleCertPaths(*caPath, *keyPath, *certPath, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("not using gRPC server just yet %p", opts)

	log.Info("Starting onos-e2sub")
	//cfg := manager.Config{
	//	CAPath:   *caPath,
	//	KeyPath:  *keyPath,
	//	CertPath: *certPath,
	//	GRPCPort: 5150,
	//}
	//mgr := manager.NewManager(cfg)
	//mgr.Run()
	//
	//if err := ioutil.WriteFile(probeFile, []byte("onos-e2sub"), 0644); err != nil {
	//	log.Fatalf("Unable to write probe file %s", probeFile)
	//}
	//defer os.Remove(probeFile)
	//
	//sigCh := make(chan os.Signal, 1)
	//signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	//<-sigCh
}
