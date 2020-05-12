Pulsar Route Maps
=================

> This project is in [active development](https://github.com/ns1/community/blob/master/project_status/ACTIVE_DEVELOPMENT.md) and is currently available as **early access**.

Utilities for working with Pulsar Route Maps.

Route maps are a way to express routing policy based on the network locations of
real users. Locations are specified as IPv4 and/or IPv6 CIDR blocks (as
opposed to Geo/ASN groups) for ultimate precision. Route maps scale to **millions
of network addresses** so you can create groupings on a granular level
for very precise and accurate control over routing.


### What can I do with this utility?

This utility streamlines CRUD operations (Create/Read/Update/Delete) for
route maps. It also does validation and linting of maps before they are
uploaded.

See [our documentation](docs/README.md) for more information.


### Who is using Route Maps?

Many of our customers!  [Dropbox wrote a blog post](https://dropbox.tech/infrastructure/intelligent-dns-based-load-balancing-at-dropbox)
about their experience using route maps to implement a global load balancing
policy.


Installation
------------

We provide binary releases for most platforms so installation is easy. Just:

1. Download the **latest release** here: https://github.com/ns1/pulsar-routemap/releases/latest
(older releases are [available here](https://github.com/ns1/pulsar-routemap/releases)).
1. Expand the release archive which includes the `routemap(.exe)` executable.
1. Run the `routemap(.exe)` command from here or copy to some location on your
Path.

Verify your new install:

```sh
$ ./routemap --version
routemap version x.y.z from <timestamp> (Git commit)
```


Quick start
-----------

You'll need an NS1 API key to get started. You can get or create one via NS1's
portal here:  https://my.nsone.net/#/account/settings

Use your API key as an argument to all commands. For example,

```sh
$ routemap --api-key xxxxxxxxxxxxxxxxx list
```

You can also add your API key as an environment variable to save providing it
on the command line:

```sh
$ export NS1_APIKEY=xxxxxxxxxxxxxxxxx

# No longer need to add --api-key param
$ routemap list
```

For information on available commands and options try:

```sh
$ routemap help
```

More detail can be found in [our documentation](docs/README.md).

Contributing
------------

Pull Requests and issues are welcome. See the [NS1 Contribution Guidelines](https://github.com/ns1/community) 
for more information.


License
-------

Copyright (C) 2020, NSONE, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
