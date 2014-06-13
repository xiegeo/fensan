## The Problem of removing stall data from deduplicated data stores

The old way to manage storage space uses the file system or database with explicit delete command. This makes it difficult for shared resources.

Specifically, I want to address content addressable and immutable data stores. Where storage is shared by deduplication. 

This problem is very similar to garbage collection (GC) of managed heap, freeing allocated memory that are no longer needed. The common solutions are reference counting and tracing. This looks like the natural solutions when extended to managing a shared data store. Infact, existing data stores, such as the Git blob store, uses tracing. This works well for such systems when deletion is sparse, and a full content transversal is also used for compaction.

When extending such a system to storage as a service. There are additional features that we would like to have.

1. End user privacy: allow the user to encrypt it's document tree, but leave leafs deduplicable, _without disclosing that the user is interested in those leafs_. This also means the server can not block users from data it’s not supposed to access (user can access any data that it has a lookup key for). This system assumes read rights is controlled by encryption, allowing anyone to cache a blob, using allowed bandwidth, but only those with the matching decryption key can truly read it.
2. Support a very large store space that can be easily subdivided without garbage collection having to transverse the whole.

## Proposal: Time To Live

Simply attach a Time To Live (TTL) to storage blobs after deduplication. A user can increase the time as desired. The server periodically scans TTL to remove outdated blobs. Unlike tracing, this scan is trivial to do incrementally.

The TTL should be implemented in large increments, such a month, to disassociate the time stamp from user's access time. The server can technically log that a user extended a blob's TTL, but such a log is not needed to provide service. Privacy conscious users can choose services that promise to not log. The server may still be forced to log some users or some blobs when required by law.

Server implementation is vastly simplyfied over tracing. This strategy also mirrors the payment structures of online storage services (paying storage time).

## Maintain The TTL, Maintain The Data

What if we want to keep something forever? This is what we think we are doing every time we press save. This is what we want when we create anything digital. 

All these things are merely illusions in existing fallible storage that we need from a storage as a service. When using such a services, it has traditionally been the provider that regularly checks for faults, such as hard drive breakdowns and recover invisibly from the user. With a TTL, it requires an agent of the user to also regularly check on storeage points.

We can look at several use cases for data storage

1. Private storage: A user with a collection of hard drives or other long term storage.

	A process that regularly extend TTL perfectly compliments regular checking of those drives, online or offline. For online drives that can be used for media servers and constant backup, the TTL extension can be done on top of a full content read verification. In effect, the TTL records the next data verification deadline. This is complimentary to network liveness checking as they discover faults in different parts of the system. For offline drives, that can be used for sticker nets and cheaper long term storage, the TTL can be configured to match the expected lifetime of the storage device and indexed on live devices capable of notifying the user to bring off line storage online for verification and TTL extension before expiration. 

2. Cloud Storage

	A cloud storage provider is likely to serve many users. With the internet being much slower than an internal network, the end users should not make full data verification themselves. This leaves byte level verification as an exercise for the provider to provide a reliable service. What TTL does is to free the provider from maintaining data structures to link users and data blobs. The user is also free to index data blobs without worrying about server compatibility and privacy leakages.


## Problem with time

Any time based system functions depends on our definition of time. Without careful consideration, TTL based garbage collection would not work properly if time is errors, for example a time server that’s configured a few years ahead could lead a cluster to consider all it’s contents to be outdated.

The system also have to use civil time, and not some internal event counter that’s more regular and dependable, as that is the time the end users communicate in.

As stated above, we only need time in large increments. I suggest the TTL to only record the month and year. A blob is outdated after the end of the stated month, and users should update the TTL way before, such as 6 month prior for an extension of a year. This also makes the usage of different time zones inconsequential.

The store should also keep track of large changes in time, to prevent unintended massive deletions from a skip in time. When the store restarts with or receives a time update far in the future, there are several possibilities:

1. The time was misconfigured or became misconfigured.

2. It was powered down for a long time.

3. It was segmented from the network for a long time. This situation can not be detected from a time skip, but the effect and recovery is identical to 2, therefore it should be tracked and treated as such.

A simple solution to such problems is for the store to be put into read only mode, with the GC stopped. This mode is turned on automatically when time skip is higher that a preconfigured setting, and can be turned back to normal mode with an administrative command, after time settings are confirmed and clients had a chance to renew TTL. 

