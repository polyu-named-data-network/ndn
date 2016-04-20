# Named Data Network
## agent / proxy
This program is designed on top of TCP/IP, it serve as the relay

## Features
 - Data Naming
 - Security
 - Routing and Forwarding
 - Caching
 - Pending Interest Table
 - Transport

## Use case
### Bind as service provider
1. connect to proxy (a.k.a. local agent)
2. register
   1. provide data name
   2. loop and wait for request (/callback?)
### Exact content lookup
### Fuzzy content lookup

## package structure
### config
 - seed

### network
 - agent
 - proxy

### protocol
 - packet

### data
 - Content Store (CS)
 - Pending Interest Table (PIT)
 - Forwarding Information Base (FIB)

## Stages
### Version 1
 - workable PIT
 - workable FIB
 - Content Name (Exact match only)
 - not support asking peer yet (if not found in FIB)
### Version 2
 - use one connection for generic json exchange
### Version further
 - add support to CS
 - add hook count
