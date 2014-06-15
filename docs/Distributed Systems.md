This is my view on distributed system, which also explains some guiding principles for Fensan

### Biolgy and society as distributed systems


 System | Member | Fault-tolerants
--------|--------|----------------
Tissue | Cell | High, cells die all the time, cell divisions creat redundant copies.
Organ | Tissue | Low, in most cases missing tissue can not be regrown. Most of the time organs can still operate at a reduced capacity, adapt, and make do with less.
Human Body | Organ | Very low, many organs are singulair and critical, few organs are duplicated as keeping them is expansive, and many critical functions are separated across different organs. The body is still ok over all, becasue it is made of redundant cells, so that organs only fail under extremely traumatic events.
Society | Person | Extremely high, individuals are highly adoptive and self sufficient. Even when a majority dissappears or are cut off, small isolated minorities can still survive. 

Most existing computer system are like the body made of organs, a web service can be made up of databases, application servers, client side application, and other middle ware. Each is a different organ.

Unlike the biological organ, servers are not made up of individual cell, but hardware and code that are difficult to build redundance into. Power outs and disconnects are also more frequent. To build redundancy on top of organs is expansive and non-trival.

The current best practice 

