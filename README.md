# onos-e2sub
E2 Subscription management module for ONOS SD-RAN (µONOS Architecture)

## Overview

The subscription service is responsible for managing the lifecycle of E2 AP subscriptions requested by applications, coordinating the assignment of subscriptions to E2 Termination nodes, and providing fault-tolerance for subscriptions and indications requested by xApps and executed by E2T. The subscription service does not interact with E2 nodes itself, but acts as a broker  between xApps and E2T.

The northbound of the subscription service is a set of gRPC services specified by the [onos-api]. xApps use the northbound API to create and manage subscriptions. Once a subscription has been created, the E2 Subscription service is responsible for assigning the subscription to an E2T node and ensuring the E2T progresses to propagate the subscription to the appropriate E2 node(s). Once a subscription has been assigned to an E2T node, the xApp is notified so it can open an indications stream to the appropraite E2T node.

The E2 termination is shipped as a [Docker] image and deployed with [Helm]. To build the Docker image, run `make images`.

[onos-api]: https://github.com/onosproject/onos-api
[Docker]: https://www.docker.com/
[Helm]: https://helm.sh
