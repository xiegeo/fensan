### Extended dynamic content schemes 

More Generally, dynamic content look ups can be seen as searching on a dataset, as the dataset changes (by adding data), the search returns new results.

To limit abuse, the search results are filtered. To promote new or more important data, the search results are ordered. 

(Alternatively, dynamic content schemes that remove the pre-trusted key requirement just need ways to include new keys)


[In the case of single publisher](The Self Updating Document.md), a search is performed on (public key, topic), filtered to only accept matching signatures, and ordered by version number, then the top result is served.

A possible scheme for an email like system. To not limit who can send, the public key is ignored and a search is performed on topic only, which include the recipents public key. Instead of just verifying the signature, some form of proof-of-work (POW) can be used to increase the cost to spammers. Ordering is less important in this case as the recipent should read every message or else increase the difficulty of POW. Of course, once a communitcation is established, the parties can use any alternative means to carry on. 

TODO: other schemes: tag/keyword search, forum/comments

To error on the sized of user privacy, servers should not provide unlimited search capblities to normal users. Only data which are specifically submited to be included in a search can be returned as a result. Even if some data appear to be formated in a dynamic content scheme, it should not be automatically include. As an example, normal users should not be able to do range requests on StaticIDs, this limit the knowledge of a StaticID to end users which have the link and server that the users use. Giving an attacker a hard time to even get the ciphertext by not implementing a feature is a win.

