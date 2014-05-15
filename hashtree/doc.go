/*
package hashtree implements merkle hash trees based on sha256 and sha244 from
sha2's family.

At a glance:

	Leaf Nodes: 1 to 1024 bytes hashed to 32 bytes using sha256. (lh)

	Inner Nodes: merges 2 32 byte hashes to 1 using sha244's initial hash values
		and sha256 without any padding. (ih)

	A file or blob is then identified by the root inner hash and data length.

	Root Hash =        ih(d,c)
	                    /    \
	d, c      =     ih(a,b),   c
	                 /    \       \
	a, b, c   =   lh(b1), lh(b2), lh(b3)

	Where b1 and b2 are blocks of 1024 bytes, and b3 must be 1 to 1024 Bytes



*/
package hashtree
