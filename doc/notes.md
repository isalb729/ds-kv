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

##功能
再做client可链接多个master，client-shell，重试机制(go func 监听ping)
再做master和data的standby，每次操作把数据放到standby，这个操作的锁和原操作一样

data standby:
1. 注册 /data-sb　数据节点，必须为空　否则清空
2. 每次put和del操作传到所有standby利用go func,处理时加本地锁
3. 节点注销的最后一个操作是通知sb
4. 发现节点爆炸后（或是所有节点都注销）注册自己此时抢占锁　(），将数据全部转移

##非功能
测试并发　写并发脚本
一致性检查
写注释
写makefile和script和report
提交
思考答辩模式