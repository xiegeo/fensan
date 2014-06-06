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

Locators are what builds Fensan into a network from individual file servers, and each file/ID/Collection forms it's own connected network.

Given a StaticID or DynamicID, locators looks for good sources quickly. 

[More on locators](Locators.md)

## 2. Resource Management

Resource management primarily focuses on bandwith and storage useage. Both are shared by many users running many tasks; with variable quality, supply, and demand; and costly maintenances and expansion.

### 2.1 Publishing / Backing Up / Redandency Factors

### 2.2 Subscriping / Retraving / Caching

### 2.3 Hierarchical / Local Networks

### 2.4 Market Economy



## 3. Secrets and Identity Management

