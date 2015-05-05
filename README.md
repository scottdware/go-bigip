## go-bigip
[![GoDoc](https://godoc.org/github.com/scottdware/go-bigip?status.svg)](https://godoc.org/github.com/scottdware/go-bigip) [![Travis-CI](https://travis-ci.org/scottdware/go-bigip.svg?branch=master)](https://travis-ci.org/scottdware/go-bigip)

A Go package that interacts with F5 BIG-IP systems using the REST API.

Some of the tasks you can do are as follows:

* Get a detailed list of all nodes, pools, vlans, routes, trunks, route domains, self IP's, virtual servers on the BIG-IP system.
* Create/delete nodes, pools, vlans, routes, trunks, route domains, self IP's, virtual servers, etc.
* Modify individual settings for all of the above.
* Change the status of nodes and individual pool members (enable/disable).

> **Note**: You must be on version 11.4+

### Examples & Documentation
Visit the [GoDoc][godoc-go-bigip] page for package documentation and examples.

### License
go-bigip is licensed under the [MIT License][license].

[godoc-go-bigip]: http://godoc.org/github.com/scottdware/go-bigip
[license]: https://github.com/scottdware/go-bigip/blob/master/LICENSE