## Simple Distributed Trust

Trust management attempts to generate a white list or group of trusted entities (public keys that can sign certain documents or actions).

It needs to be distributed so that each end user fully controls who to trust (no universal Central Authority (CA)), well also receiving recommendations from others to collectivly manage trust.

Trust is directional, A trusts B does not mean B trusts A.

Trust is not fussy, an entitie can ether do something or not. It can't maybe do something. You can work around this by creating multiple groups of trust with different capabilities. This document will concentrat on creating a single group.

Trust is not transitive, A trust B and B trust C does not mean A trust C. Otherwise users can not select to only trust a subset of a community, and adding or removing any entity require consensus.

Trust can be distributed explicties, A trust who B trust and B trust C, then A trust C. In extention, A can trust who B trust trust, and so on. This describes maximum degrees of separation (degs) that trust is allowed to propagate in a chain, encodable as two intergers T and P. 

T is take, or how many degs from that user to trust. 

P is publish, or how many degs from that user to suggest others to trust.

T ≥ P

T allows users to trust something without suggesting others to do the same, otherwise T = P.

Sudo code to get all trusted entities:
	
	input: Ts list[(Entity, int)] //all my takes
	input: Pm map[Entity, list[(Entity, int)]] //every Entity's publish


	//trusted is a map of ents to trustablity T,
	//unknown ents have a T of 0, or untrusted.
	var trusted map[Entity, int] default 0 
	
	for each (ent, T) in Ts{
		add(ent, T)
	} 
	
	for each (ent, T) in trusted{
		iff T > 0 then ent is trusted
	}
	
	where:
	
	func add(ent Entity, T int){
		oldT = trusted[ent]
		if oldT < T{
			//increase trust for this ent
			trusted[ent] = T
			if T > 1{
				Ps = Pm[ent]
				for each (toAdd, P) in Ps{
					adjusted = min(T-1, P)
					add(toAdd, adjusted)
				} 
			}
		}
	}
	
In the above code Pm is made of many [dynamic files](The Self Updating Document.md). The code should rerun eveytime any update occurs.

