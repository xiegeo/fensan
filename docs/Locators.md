## Locators 

Locators are what builds Fensan into a network from individual file servers.

Given a StaticID or DynamicID, locators looks for good sources quickly. 

The fastest location methods will be build first, including Current (the currently connected server that served the ID), and Linked (locations that are included in the URL with the ID). Corrent and Linked are simular to using relative and absolute URLs in HTTP, by themselves it should allow Fensan to be at least equvilent to HTTP in serving static content in terms of speed and reliability.

Other methods can be consitered in the future, such as peer exchange, centralized indexing, and distributed indexing. They currently work well in peer to peer (P2P) distribution of large files. But Fensan is not a traditional P2P file sharing network, the details will be considered after Fensan has had a bit of use to see what's needed.

When locators list many sources, a selection should be made to try good sources first. Good sources should be fast in reaction time, have high bandwith for large downloads, and cheap in case of metered bandwith or costs to user's account on server. Many of these properties are difficult to predict and can change at anytime. Heristics for selection algrithms will be expermented on. 

Sources are not required to be trusted, so locators have the freedom to innovate. however, some use cases may have privacy implications: A user may not want to brodcast his interest in some files to the world. The file maybe encrypted, but the ID can still correlate a group of users. Or the file is public, and the user want to avoid tracking. Locators should inclued ways to limit this disclosure or be possible to disable.
