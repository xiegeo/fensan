The Self Updating Document
==

Features
--
- Protected by public key signature.
- Updates by version number.
- Updates future downloading locations if changed.
- Download content by demand or cache as early as possible.
- Does not track history in the sync level, gave app the freedom to implement history as fit. But factual history can still be technically observed/logged. 


Format 
--
<small>{repeat}  [optional]  |or  ,concat</small>

	id, [sources], [...], sign, [content]
	
	where
	
	id = pk, [topic], version, content_hash
	
	sources = {server locations | alternate protocols | ...}
	


"pk" is the public key of the signer. The key infostructure is outside the scope of this document... (to do)

"topic" is used to allow one public key to sign many documents.
If "topic" is left out, it is the same as empty string "". (so we don't need to seperate NULL and "")

"topic" cannot be changed and does not have to be descriptive. (An app may use content to store redirect infomation)

??? "topic" could follow a URI like scheme to avoid collision. Or/And use different public keys for different...

"topic" should not be used for storage of other information, such as file name and type. which belongs in content, the context it is linked from, or [...].

"pk, [topic]" identifies a document and is the key used for subscribing to updates.

"version" is an integer 0 to 2^63 − 1 (chosen to be within a signed 64 bit, and large enough for time stamps).
Only the document with the larges known version number is accepted as current.  
If the version number is equal, then "sign" is used to break even (larger wins), allowing the network to stay consistent under attacks or misconfigurations.

Sources should include multiple locations where updates can be subscribed to and content can be downloaded.

... is whatever that maybe added in the future or for application specific features, they can be safely ignored but must be included in signature calculations

Validation
--

- "pk, [topic]" must be as requested.
- "version", or "version" and "sign" must be greater than existing.
- "sign" must validate "id, [sources], [...]"
- If hash(content) != content_hash, then content is corrupted or missing, the rest of "id, [sources], [...], sign" are still valid.
- An update can come from any source, even those not included in the "sources" list. Allowing  external source finders.

Anything other than content are public and not encrypted. Such that any server can work with them.  
They should not reveal private information and function similare to public certificates and DNS (as in publicly redistributable by anyone).

Use cases
---

### Backup / Subscription ###

Alice produce an document, Bob makes a local copy of the latest version.

setup:

- Alice gave Bob "pk, [topic], [sources]", where Alice has the private key for "pk" and intends or already publish in "pk, [topic]" on sources
- Bob subscribe to "pk, [topic]" on sources.
- Alice update "pk, [topic]" on sources with increasing version number, Bob recive.
- Bob resubscribe to "pk, [topic]" on sources when sources change.

For backup, Alice request Bob to do the subscription.

### Bidirectional Sync ###

Alice and Bob may both want to contribute to a document, or Alice may have many devices she uses to edit the same document, even when disconnected.

To avoid lost edits by having another device produce a higher version number before it can be merged. No two devices can ever publish to the same "pk, [topic]". This is best enforced by not sharing private keys of devices. Only the public keys of devices are exported, signed as a trusted, and subscript to.

Just as a higher version number implies the coverage of content for all previous versions, bidirectional sync must also include references to the last versions of all sources it merged from. This will tell content editors when merges are needed, in the case that both versions are newer than referenced. 

References must include "pk", "[topic]", "version", and "sign", functioning as secure [Vector clocks](http://en.wikipedia.org/wiki/Vector_clock). Such refercencs can be inclueded in content only, encrypted from servers.

see also [Extended Dynamic Content Schemes](Extended Dynamic Content Schemes.md)