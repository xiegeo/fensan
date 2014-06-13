Warning, IMHO, the CAP theorem, and any other attempts  to define how a distributed system should behaive as if it is a single system, is harmful. The following have changed definitions to better fit a distributed system that is not trying to emulate a single system.

## Consistency

- Content addressable immutable static data by hashing.

- Single producer for dynamic data by signing. Producer also versions the dynamic data for globle ordering of changes. Client side merging for multipal producers.

- References/links to static and dynamic data have different formats. This allows linking to the latest version of a database, by a single dynamic link to a snapshot, which contains collections of static links, that represents the immutable internal structure of the snapshot.

## Availability

- Any server can serve cached static data, which never invalidates. Clients can also request servers to reserve long term storage for some data.

- Any server can serve the latest dynamic data it knows about. Clients can check against it's own version to make sure the version is not regressing, and request multipal servers in case some are behind.

- Any server means anyone can add availability. No one can dictate to others which servers must be or not connected to. You can receive updates as long as the producer is somewhere up stream.

## Partition Tolerance

- The is no globle consensus that need to be reached at run time. 


### Trade-offs in Availability vs Delta Consistency
Or what to do when producers and consumers are partitioned. 

Producers for dynamic data can add an expiration date for each version. When the time passes and no updates arrive, a warning can be displayed. Alternatively, the topic, which is part of the identifier of dynamic data, can be appended with a time stamp, for example the date for a daily newsletter. This allows users to retrive the news for a date or fails, instead of retriving a news which maybe days late.
