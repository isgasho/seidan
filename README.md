# Seidan

Seidan is still in development. Please come back later.

# Why?

Most cloud software out there is either way too complex for what it does,
or does not provide the right tools to work with large clusters. It should be
possible to grow clusters easily, add/remove services in a few clicks,
and know that any lost node is going to be replaced instantly.

# Concept

Seidan is a cluster management system that runs on each single node ("star"),
forming groups ("constellation"? maybe a bit too long) which are all part of
a cluster.

Nodes are grouped by latency.

Each node has its own sub-CA which is signed by the global CA. Typically, the
global CA is kept offline and only used when new nodes are started, and private
keys are never shared outside of a given node.

## Daemons (Planets?)

Seidan can then run daemons on the specified node and re-run them in case of
failure, monitor various information about the machine (load, memory, disk,
etc), apply OS updates, etc.

Each started daemon can have a key and certificate issued by the local node CA
which can then be used to identify that specific daemon on that specific node
toward other daemons running on other nodes of the same cluster.

This approach allows nodes of a given cluster to be able to authenticate any
other node using strong encryption, ensuring not only security of
communications between nodes, but also secure authentication, and discovery.

## Outside configuration

Seidan makes available data on the running cluster to the local processes so
configuration can happen on various third party or local services, such as
high availability, etc.

## Services

Each launched daemon can expose itself to the rest of the cluster as a service
and accept various connections from other services. Based on resources usage
new instances can be launched automatically. Small clusters with less needs can
have a single node running multiple low traffic daemons, and cluster can grow
automatically depending on load.

## Non-cloud

Seidan does not depend on any kind of cloud technology such as kubernetes,
aws, google cloud, docker, etc. Its setup is simple, adding nodes is just a
matter of launching seidan on a new machine and signing its certificate
request. Clusters can be made of any kind of machine running an UNIX operating
system, mixed in any kind of configuration, in any country.

## Services

Seidan does the following:
* Node discovery, cluster joining
* Handling of one sub-CA per node, signature of certificates for subnodes
* Configuration storage
* Detection of latency groups (typically, datacenters)
* Nodes monitoring
* Download, launching and configuration of software

The following will be provided by other software:
* Inter-node virtual networking (vpn-like, encrypted, peer to peer and secure)
* Logging facilities
* High throughput decentralized database
* pkg access daemon

