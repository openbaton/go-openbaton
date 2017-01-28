Go libraries for OpenBaton
==========================

[![GoDoc](https://godoc.org/github.com/mcilloni/go-openbaton?status.svg)](https://godoc.org/github.com/mcilloni/go-openbaton)

`go-openbaton` contains several packages that can be used to write services that interface with the [OpenBaton][openbaton] [NFVO][nfvo] using the Go language.

## Packages

- [catalogue](https://github.com/mcilloni/go-openbaton/tree/master/catalogue): provides a partial implementation of the OpenBaton catalogue.
- [catalogue/messages](https://github.com/mcilloni/go-openbaton/tree/master/catalogue/messages): defines the default message types for NFVO-VNFM communication, plus facilities to handle their serialisation.
- [plugin](https://github.com/mcilloni/go-openbaton/tree/master/plugin): provides a runtime to develop and execute plugins for the NFVO.
- [vnfm](https://github.com/mcilloni/go-openbaton/tree/master/vnfm): provides a runtime to develop and execute VNFManagers in Go.
- [vnfm/channel](https://github.com/mcilloni/go-openbaton/tree/master/vnfm/channel): a set of interfaces that provide an abstraction above which API the VNFM uses to connect to the NFVO.
- [vnfm/amqp](https://github.com/mcilloni/go-openbaton/tree/master/vnfm/): implements a `channel` that uses AMQP to connect with the NFVO.
- [vnfm/config](https://github.com/mcilloni/go-openbaton/tree/master/vnfm/config): provides facilities for parsing VNFM configuration files.

## Issue tracker

Issues and bug reports should be posted to the GitHub Issue Tracker of this project

# What is Open Baton?

Open Baton is an open source project providing a comprehensive implementation of the ETSI Management and Orchestration (MANO) specification and the TOSCA Standard.

Open Baton provides multiple mechanisms for interoperating with different VNFM vendor solutions. It has a modular architecture which can be easily extended for supporting additional use cases. 

It integrates with OpenStack as standard de-facto VIM implementation, and provides a driver mechanism for supporting additional VIM types. It supports Network Service management either using the provided Generic VNFM and Juju VNFM, or integrating additional specific VNFMs. It provides several mechanisms (REST or PUB/SUB) for interoperating with external VNFMs. 

It can be combined with additional components (Monitoring, Fault Management, Autoscaling, and Network Slicing Engine) for building a unique MANO comprehensive solution.

## Source Code and documentation

The Source Code of the other Open Baton projects can be found [here][openbaton-github] and the documentation can be found [here][openbaton-doc] .

## News and Website

Check the [Open Baton Website][openbaton]

## Licensing and distribution
Licensed under the Apache License, Version 2.0. See LICENSE for further details.

[openbaton]: http://openbaton.org
[openbaton-doc]: http://openbaton.org/documentation
[openbaton-github]: http://github.org/openbaton
[nfvo]: https://github.com/openbaton/NFVO
[NFV MANO]:http://docbox.etsi.org/ISG/NFV/Open/Published/gs_NFV-MAN001v010101p%20-%20Management%20and%20Orchestration.pdf
