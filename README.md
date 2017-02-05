  <img src="https://raw.githubusercontent.com/openbaton/openbaton.github.io/master/images/openBaton.png" width="250"/>

  Copyright © 2015-2016 [Open Baton](http://openbaton.org).
  Licensed under [Apache v2 License](http://www.apache.org/licenses/LICENSE-2.0).
  
  
# Go libraries for OpenBaton

[![GoDoc](https://godoc.org/github.com/openbaton/go-openbaton?status.svg)](https://godoc.org/github.com/openbaton/go-openbaton)

`go-openbaton` contains several packages that can be used to write services that interface with the [Open Baton][openbaton] [NFVO][nfvo] using the Go language.

## Packages

- [catalogue](https://github.com/openbaton/go-openbaton/tree/master/catalogue): provides a partial implementation of the Open Baton catalogue.
- [catalogue/messages](https://github.com/openbaton/go-openbaton/tree/master/catalogue/messages): defines the default message types for NFVO-VNFM communication, plus facilities to handle their serialisation.
- [plugin](https://github.com/openbaton/go-openbaton/tree/master/plugin): provides a runtime to develop and execute plugins for the NFVO.
- [vnfm](https://github.com/openbaton/go-openbaton/tree/master/vnfm): provides a runtime to develop and execute VNFManagers in Go.
- [vnfm/channel](https://github.com/openbaton/go-openbaton/tree/master/vnfm/channel): a set of interfaces that provide an abstraction above which API the VNFM uses to connect to the NFVO.
- [vnfm/amqp](https://github.com/openbaton/go-openbaton/tree/master/vnfm/): implements a `channel` that uses AMQP to connect with the NFVO.
- [vnfm/config](https://github.com/openbaton/go-openbaton/tree/master/vnfm/config): provides facilities for parsing VNFM configuration files.

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
Follow us on Twitter @[openbaton][openbaton-twitter].

## Licensing and distribution
Copyright © [2015-2017] Open Baton project

Licensed under the Apache License, Version 2.0 (the "License");

you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## Support
The Open Baton project provides community support through the Open Baton Public Mailing List and through StackOverflow using the tags openbaton.

## Supported by
  <img src="https://raw.githubusercontent.com/openbaton/openbaton.github.io/master/images/fokus.png" width="250"/><img src="https://raw.githubusercontent.com/openbaton/openbaton.github.io/master/images/tu.png" width="150"/>

[fokus-logo]: https://raw.githubusercontent.com/openbaton/openbaton.github.io/master/images/fokus.png
[openbaton]: http://openbaton.org
[openbaton-doc]: http://openbaton.org/documentation
[openbaton-github]: http://github.org/openbaton
[openbaton-logo]: https://raw.githubusercontent.com/openbaton/openbaton.github.io/master/images/openBaton.png
[openbaton-mail]: mailto:users@openbaton.org
[openbaton-twitter]: https://twitter.com/openbaton
[nfvo]: https://github.com/openbaton/NFVO
[NFV MANO]:http://docbox.etsi.org/ISG/NFV/Open/Published/gs_NFV-MAN001v010101p%20-%20Management%20and%20Orchestration.pdf
[tub-logo]: https://raw.githubusercontent.com/openbaton/openbaton.github.io/master/images/tu.png
