### requirement:
zookeeper:
1. configuration management
2. cluster management 
3. lock service
4. naming service 
5. at least 3 nodes
show the build-up details in report

master/slave architecture
1. master node manage the meta data(including locations of the data nodes)
2. data node, access the data directly
3. The node information is registered in Zookeeper and will used as metadata by master.
4. Communication between client - master node- data nodes should use RPC
5. Data as key-value the type of key and value are BOTH String
6. Data is distributed by Key (determine the location of data)
7. Whats more:
    How to ensure Load Balance?
    What if data node(s) is added dynamically (scalability)?
8. 3 operations for data accessing: 
    PUT
    READ
    DELETE
9. How to deal with concurrent data accessing ?
    i.e. concurrency control 
10. At least 2 standby(idle, ready) nodes for Data node for fault tolerance
How to keep the consistency between primary node and standby node?
11. What about standby node for Master? (Optional)

