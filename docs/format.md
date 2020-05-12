Route map data exchange format
==============================

> :information_source: This document is for **Version 1**.
> 
> :heavy_check_mark: This is the **Current Version**.


This document describes the data exchange format that customers use to send 
route maps to NS1.


### Basic format

Routemaps are JSON documents with the following basic format. Note that comments
are not supported.

```json
{
  "meta": {
    "version": 1
  },
  "map": [
    {
      "networks": [],
      "labels": []
    }
  ]
}
```

### Fields, data types, and usage semantics

| Field | Data type | Description |
| ----- | --------- | ------------|
| `meta` | Object | Metadata related to the exchange format itself. This object is always required. |
| `[meta] version` | Integer | (Required)  Version of the exchange data format.  Determines the semantics of the remainder of the exchange. |
| `map` | Array of Objects | Each member object declares a segment of the map as defined by the list of network addresses contained within. |
| `[map] networks` | Array of Strings | (Required)  IPv4 and/or IPv6 Network addresses in CIDR presentation form.  Can include both v4 and v6 addresses in a single network definition. |
| `[map] labels` | Array of Strings | (Required)  Array of labels associated with the networks. These are arbitrary ASCII-only strings used to associate DNS answers with map segments. **Note**: The order in which you list the labels is important as it determines the order in which the answers are emitted. For example, a list of `["b", "c", "a"]` means emit DNS answers for "b" first, then "c", and then "a". |


### Limits

* IPv4 prefix length up to /26
* IPv6 prefix length up to /64


### Simple example

The following example map defines three networks. Note that the first map 
segment associates networks 192.168.10.0/24 - 192.168.60.0/24 with two labels.

```json
{
  "meta": {
    "version": 1
  },
  "map": [
    {
      "networks": [
        "192.168.10.0/24",
        "192.168.20.0/24",
        "192.168.40.0/24",
        "192.168.50.0/24",
        "192.168.60.0/24"
      ],
      "labels": ["hkg", "sin"]
    },
    {
      "networks": [
        "192.168.132.0/24",
        "192.168.133.0/24",
        "192.168.135.0/24",
        "192.168.144.0/20"
      ],
      "labels": ["nrt"]
    },
    {
      "networks": [
        "172.16.5.0.0/24",
        "172.16.1.0/24",
        "192.168.111.0/24",
        "fd0d:82a7:b e8b:ae00::/56",
        "fde2:d85e:f372:7408::/64"
      ],
      "labels": ["syd"]
    }
  ]
}
```
