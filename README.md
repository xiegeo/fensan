FenSan
======

Dispersed Data Cache Store

Dispersed VS Distributed

- Dispersed implies less organization, and more organic and natrual goodness.

- FenSan minimize build in structures, everything from read/write rights (controlled by client side encryption), file/folder metadata (just more encrypted blobs), to replication factors (how many servers a client asks for storeage) are all implemented by clients. This does not exclued servers from doing any operations on the clients behalf, but servers can not unilaterally operate on the data without clients decryption keys other than copy and remove.

- How minimal? The servers does not even need to track which users need which blobs saved. The server can't since folders can be encrypted client side. To remove stale data, the server keeps a coarse Time To Live on each blob, which any client may request or increase using credit. 

- Unlike most distributed systems, FenSan can not slipt-brain, by not sharing brains in anyway. Each server exits on it's own, with ablilty to support all native operations without outside contact. Only need to contact other servers when it need new information. As such, there is no global state, only local state on the server level. It pushes all data merging to the application level. 

- FenSan is disperse in Chinese PinYin, FenBu (distribute) sounds stupid and overly generic, comming up with good names is hard.


