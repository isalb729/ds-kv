## Distributed Key-Value Storage System
517021910851-于亚杰

In this lab I've built a distributed key-value storage system based on zookeeper and client-server architecture. The system is aimed to achieve high availability and data consistency, without much consideration of performance.

### Dependencies
This repo mainly relies on the following third-party projects:
+ gRPC & protobuf - rpc code generation
+ zookeeper - for distributed services
+ go-zookeeper - zookeeper client written in golang

### Lab environment
+ golang 1.13
    + with go module as package management tools
+ protocol buffer
    + only used in development to generate rpc code
+ zookeeper 3.6.1
    + three nodes as a cluster
+ ubuntu 18.04

### Zookeeper cluster
A three-node zookeeper cluster is built in the very beginning of this lab. The detailed steps are shown as follows.
1. Download zookeeper from https://mirrors.sonic.net/apache/zookeeper/zookeeper-3.6.1/apache-zookeeper-3.6.1-bin.tar.gz
2. Extract the file and put it under a `zk` directory. Make 3 copy of the directory respectively named `zk-1`, `zk-2` and `zk-3`.
3. Modify the configuration files, e.g. add the following content to `zk-1/conf/zoo.cfg`. 
    ```
    tickTime = 2000
    dataDir = /home/blasi/Desktop/distributedsystem/ds-kv/zk/zk-1/data
    clientPort = 2181
    initLimit = 5
    syncLimit = 2
    maxClientCnxns=100
    server.1=localhost:2888:3888
    server.2=localhost:2889:3889
    server.3=localhost:2890:3890
    ```
    This configuration mainly specifies the port it runs on, data directory and information about the other nodes in the cluster.
4. Then to establish the whole cluster locally, run `make zk` under the project directory. The script is actually put in `scripts/zk.sh` which can start or stop all the zookeeper node under `zk` directory.

### Directory structure
 * [go.sum](./go.sum) 
 * [go.mod](./go.mod)
 * [src](./src)
   * [cmd](./src/cmd)
     * [cfg.yaml](./src/cmd/cfg.yaml)
     * [client](./src/cmd/client)
       * [client_main.go](./src/cmd/client/client_main.go)
     * [client-shell](./src/cmd/client-shell)
       * [shell.go](./src/cmd/client-shell/shell.go)
     * [server](./src/cmd/server)
     * [server_main.go](./src/cmd/server/server_main.go)
   * [zookeeper](./src/zookeeper)
     * [lock.go](./src/zookeeper/lock.go)
   * [utils](./src/utils)
     * [config.go](./src/utils/config.go)
     * [collection.go](./src/utils/collection.go)
     * [conv.go](./src/utils/conv.go)
     * [hash.go](./src/utils/hash.go)
     * [prime.go](./src/utils/prime.go)
     * [label.go](./src/utils/label.go)
     * [data.go](./src/utils/data.go)
   * [rpc](./src/rpc)
     * [pb](./src/rpc/pb)
       * [master.pb.go](./src/rpc/pb/master.pb.go)
       * [slave.pb.go](./src/rpc/pb/slave.pb.go)
     * [master.go](./src/rpc/master.go)
     * [proto](./src/rpc/proto)
       * [master.proto](./src/rpc/proto/master.proto)
       * [slave.proto](./src/rpc/proto/slave.proto)
     * [slave.go](./src/rpc/slave.go)
   * [test](./src/test)
     * [test.go](./src/test/test.go)
   * [client](./src/client)
     * [cli.go](./src/client/cli.go)
     * [op.go](./src/client/op.go)
   * [server](./src/server)
   * [master.go](./src/server/master.go)
   * [slave.go](./src/server/slave.go)
 * [Makefile](./Makefile)
 * [doc](./doc)
   * [lab 5 description.docx](./doc/lab 5 description.docx)
   * [Lab5 介绍 2020.6.2.pptx](./doc/Lab5 介绍 2020.6.2.pptx)
   * [notes.md](./doc/notes.md)
   * [report.md](./doc/report.md)
 * [scripts](./scripts)
     * [source_proxy.sh](./scripts/source_proxy.sh)
     * [zk.sh](./scripts/zk.sh)
     * [run.sh](./scripts/run.sh)

Source code is placed under `src` directory and `scripts` are mainly used for development.
### Design of the system
Next i'm going to talk about some high level design of the storage system which includes scalability, consistency and some other topics.
The code or implementation itself is more complicated so for more details please go for the source code and comments.
#### Nodes
There are two types of nodes: master and data. Both have standby nodes to improve availability.
 
Master nodes help manage the metadata of data nodes and transform standby nodes. Client can ask the master nodes for address of data nodes and operate on data node directly. 

#### Registration and deregistration
When a node registers, it creates an ephemeral znode in zookeeper with the name of its address. 
In this way, there are no identical server names. Data node registration will also create the data directory and move some data from other servers to itself(discussed in the next part).
Then a goroutine keeps fetching information of standby data nodes, which is used for backup data. Master registration and deregistration do nothing except create and delete znode.
Data node deregistration will move its data to other servers and notify the master(in case of server failures).
Master node watches the data znode in zookeeper and when a data node register, master helps label the server, that is, write an integer value to this server's znode to implement jump consistent hash algorithm(discussed in the next part).
The label value is between 0 and maximum server number(abbreviated MSN). Imagine a ring with MSN seats on it. Every time a node register, it takes a seat and tries to be as far as possible away from the adjacent nodes, just like social distancing.
For example, if MSN equals 137 and there are two nodes labelled 0 and 68. Then next server registered can be labeled 102 or 103.

Both registration and deregistration needs a shared "register" lock with the help of zookeeper. It's implemented in the same way as "master" lock mentioned above.
#### Data
Data are stored as key value pairs in files. Both key and value have the type of string. Data file operations can be found in `src/server/utils/data.go`.
Normally a key value pair is only stored in one active data server and some of the standby data servers. 

Data are distributed to different data node by key. A basic hash method is applied to key and it determines which server the key goes to based on jump consistent hash. 
To make it more clear I'll give a simple example here. Suppose MSN is 137(default value), and two nodes are labelled 0 and 68. 
A key has a hash value of 1395665. Then divide the hash value by MSN and the remainder is 46. It's closer to 68 then to 0, so the key is stored in the server labelled 68.

The data are also distributed in the server itself to make the whole system more scalable. For example, one level storage has only one file to store the data while two level storage can have 3*5 files.
The numbers of directories/files in each level are increasing prime numbers. In a two level storage server with 3*5 files, the key with hash value 1395665 will be stored in file `2/0` because 1395665 mod 3 = 2 and 1395665 mod 5 = 0.
We can prove that each key has the same probability to be stored in each file by chinese remainder theorem.
#### Standby master node
There is only one available master node at any time, the one keeping the "master" lock implemented with zookeeper. 
The lock details can be found in `src/zookeeper/lock.go`. 
Standby master will try to grab the lock when starting up and get the lock when the old master exits.
The client end can specify multiple master addresses. There is a goroutine running to ping the master every 300ms. 
If it can't ping the master, the client will try to connect other master addresses. 
#### Standby data node
Standby data node create a `sb` znode when registering. And it will stores all the latesing data after the startup. It won't provide backup for earlier data because it's too costly to block all servers to transfer data.
In this way, if data is created before any standby node, it could be lost if the server faces something terrible like power failure. The standby node can register itself as data node after a data server passes out, which means its znode disappeared but the deregistration function is not executed.
This is detected by master node, which get the notification of data server deregistration. When a data znode disappear and the master didn't get the notification, the master will delete one of the standby data znode. And then that standby node knows he is appointed as new data server.
He will register himself as data server, and remove all the data that shouldn't be stored in this server anymore.
#### Operations
Four client operations are provided: get, put, delete, dump. Dump is only used for debugging purpose which shows the data storage information of all servers.
See `src/client/op` for more information. Get operation will ask master for address of data server and then send a get request to it.
Put and delete operations work in similar way. Dump can be seen as advanced get operation.
Concurrent data accessing and strong consistency is supported. Data locks are file-level, which means operations on two keys stored in the same file will be sequential.
This is more scalable than key-level lock and performs better than server-level lock. And more storage levels can help improve the performance.

### Practice
