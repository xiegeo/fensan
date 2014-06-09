# Overview


## 1. Global File Space

The global file space consists of data identifiers, verification algrithms, and locators.

Unlike https, which verifies the server or location providing the data. FenSan verifies the data directly based on identifiers. This allows locators to choose arbitery locations and allow anyone to provide data caching services. 

### 1.a Content addressable StaticID 

StaticID -> (hash, length) -> Merkle tree -> binary data

StaticIDs provide a globally unique name space that maps cryptographic hashes to immutable data.

https://godoc.org/github.com/xiegeo/fensan/hashtree

The binary data from a StaticID can represent anything, including plaintext files, encrypted blobs, and collections, which are lists of StaticIDs or other links, which can be used to represent folders.

### 1.b Public key signed DynamicID (Single Publisher)

DynamicID -> (public key, topic) -> (public key, topic, version number, StaticID, meta data, signature) -> updated StaticID -> updated contends

see also [The Self Updating Document](The Self Updating Document.md), [Extended Dynamic Content Schemes](Extended Dynamic Content Schemes.md)

### 1.d Locators 

Locators are what builds Fensan into a network from individual file servers, and each file/ID/Collection/Task forms it's own connected network.

Given a StaticID or DynamicID, locators looks for good sources quickly. 

[More on locators](Locators.md)

## 2. Basic Server Features

Server and clients are not differentiated by functionally. The server represents an instance that is long running, have some data a client want, or prove storage.

### 2.0 bootstraping
The server starts off reading a local configeration file, that helps it connect to the network. This should be the only file, other than the program itself, that it needs outside of the network or it's own database, both of which are nil on first start up.

### 2.1 Retain
What the server should keep locally. Retain increase data redandency

[TTL based](Time To Live for Garbage Collection of Shared Storage.md)

### 2.2 Subscripe 
Retain for dynamic contents, push based, to always keep up to date.

Server settings should be done using subscripe. This allows a cluster of servers to be configered together. Using collections, servers can share some but not all configerations.

### 2.3 Proxy
Retrave and cache remote contents on clients request. This allows well positioned servers to cache content; and clients to use less connections and imporve privacy.

Cached content that are not retained are not asked to increase redandency, they can be purged whenever.

### 2.4 Users
The server need to enact user requests for the above features.


### 2.5 Tasks
Tasks are how the server management it self. 

Transient tasks are short lived tasks that fill user requests. The task states only live in memory. The tasks can 'crash' when the server stops or when relevant clients disconnect. 

Resident tasks are long lived tasks that should continue even after the server restarts. Some Resident tasks are run periodically.

(Future: allow users to run arbitery programs in a sandbox. This allow new features without hard coding them in the server. )

## 3. Resource Management

Resource management primarily focuses on bandwith and storage useage. Both are shared by many users running many tasks; with variable quality, supply, and demand; and costly maintenances and expansion.

Resource management has two sides: a client selecting from many servers, and a server selecting which clients to serve. 

 

### 3.1 Publishing / Backing Up / Redandency Factors
Clients decide where data is stored, and how many copies exist, based on intened usage. 

### 3.2 Subscriping / Retraving / Caching
Anyone else who know about a piece data, va a link, can copy over the data. They 

### 3.3 Hierarchical / Local Networks

### 3.4 Market Economy



## 4. Secrets and Identity Management

