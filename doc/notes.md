### requirement:
zookeeper:
1. configuration management(meta data of server)
2. cluster management(register and deregister)
3. lock service(global and server read write lock)
4. naming service(todo:?)
5. at least 3 nodes

master/slave architecture
1.  master works as message queue
1. master node manage the meta data(including locations of the data nodes)
2. data node, access the data directly
3. The node information is registered in Zookeeper and will be used as metadata by master.
4. Communication between client - master node- data nodes should use RPC.
5. Data as key-value the type of key and value are BOTH String.
6. Data is distributed by Key (determine the location of data).
7. Whats more:
    data replication
    How to ensure Load Balance? (qos, everytime get a key, master return the flow and server pair)
    What if data node(s) is added dynamically (scalability)? (jump consistent hashing)
8. 3 operations for data accessing: 
    PUT
    READ
    DELETE
9. How to deal with concurrent data accessing ?
    i.e. concurrency control 
10. At least 2 standby(idle, ready) nodes for Data node for fault tolerance
How to keep the consistency between primary node and standby node?
11. What about standby node for Master? (Optional)

# ds-kv
TODO:
1. zk cluster connect, config manage, lock, naming service, cluster manage
2. zk client, create delete get, ephemeral, sequential
3. server establish
4. server basic
unit test
5. server cluster
6. finish advanced
unit test
7. write shell client
8. loadtest


advanced:
shared lock
data realloc when deregistering
master replica
multi master run on the same port