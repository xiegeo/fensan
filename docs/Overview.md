# Overview


## 1. Global File Space

The global file space consists of data identifiers, verification algrithms, and locators.

Unlike https, which verifies the server or location providing the data. FenSan verifies the data directly based on identifiers. This allows locators to choose arbitery locations and allow anyone to provide data caching services. 

### 1.a Content addressable StaticID 

StaticID -> (hash, length) -> Merkle tree -> binary data

StaticIDs provide a globally unique name space that maps cryptographic hashes to immutable data.

https://godoc.org/github.com/xiegeo/fensan/hashtree

### 1.b Public key signed DynamicID (Single Publisher)

DynamicID -> (public key, topic) -> (public key, topic, version number, StaticID, meta data, signature) -> updated StaticID -> updated contends

see also [The Self Updating Document](The Self Updating Document.md), [Extended Dynamic Content Schemes](Extended Dynamic Content Schemes.md)


### 1.c Locators 

Locators are what builds Fensan into a network from individual file servers.

Given a StaticID or DynamicID, locators looks for good sources quickly. 

[More on locators](Locators.md)

