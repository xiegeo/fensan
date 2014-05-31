# Overview


## Global File Space


### Contend addressable static contend 

StaticID -> (hash, length) -> Merkle tree -> binary data


### Public key signed dynamic contend

DynamicID -> (public key, topic) -> (public key, topic, version number, StaticID, meta data, signature) -> updated StaticID -> updated contends

