## Distributed Key-Value Storage System

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
A three-node zookeeper cluster is built in the very beginning of this lab. The detailed establishment steps are shown as follows.
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
        * [test.go](./src/cmd/client/test.go)
      * [client-shell](./src/cmd/client-shell)
        * [shell.go](./src/cmd/client-shell/shell.go)
      * [server](./src/cmd/server)
      * [server_main.go](./src/cmd/server/server_main.go)
    * [utils](./src/utils)
      * [config.go](./src/utils/config.go)
      * [collection.go](./src/utils/collection.go)
      * [conv.go](./src/utils/conv.go)
      * [hash.go](./src/utils/hash.go)
      * [prime.go](./src/utils/prime.go)
      * [label.go](./src/utils/label.go)
      * [data.go](./src/utils/data.go)
    * [client](./src/client)
      * [op.go](./src/client/op.go)
    * [rpc](./src/rpc)
      * [pb](./src/rpc/pb)
        * [master.pb.go](./src/rpc/pb/master.pb.go)
        * [slave.pb.go](./src/rpc/pb/slave.pb.go)
      * [proto](./src/rpc/proto)
        * [master.proto](./src/rpc/proto/master.proto)
        * [slave.proto](./src/rpc/proto/slave.proto)
      * [master.go](./src/rpc/master.go)
      * [slave.go](./src/rpc/slave.go)
    * [server](./src/server)
      * [master.go](./src/server/master.go)
      * [slave.go](./src/server/slave.go)
    * [zookeeper](./src/zookeeper)
    * [lock.go](./src/zookeeper/lock.go)
 * [Makefile](./Makefile)
 * [scripts](./scripts)
     * [source_proxy.sh](./scripts/source_proxy.sh)
     * [zk.sh](./scripts/zk.sh)
     * [tree.sh](./scripts/tree.sh)
     * [run.sh](./scripts/run.sh)
 * [report.md](./report.md)


Source code is placed under `src` directory and `scripts` are mainly used for development.
### Design of the system
Next i'm going to talk about some high level design of the storage system which includes scalability, consistency and some other topics.

The code or implementation itself is more complicated so for more details please go for the source code and comments.
#### Nodes
There are two types of nodes: master and data. Both have standby nodes to improve availability.
 
Master nodes help manage the metadata of data nodes and transform standby nodes. Client can ask the master nodes for address of data nodes and operate on data node directly. 

#### Registration and deregistration
When a node registers, it creates an ephemeral znode in zookeeper with the name of its address. In this way, there are no identical server names. Data node registration will also create the data directory and move some data from other servers to itself(discussed in the next part).
Then a goroutine keeps fetching information of standby data nodes, which is used for backup data.

Master registration and deregistration do nothing except create and delete znode.

Data node deregistration will move its data to other servers and notify the master(in case of server failures).

Master node watches the data znode in zookeeper and when a data node register, master helps label the server, that is, write an integer value to this server's znode to implement consistent hash algorithm(discussed in the next part).

The label value is between 0 and maximum server number(abbreviated MSN). Imagine a ring with MSN seats on it. Every time a node register, it takes a seat and tries to be as far as possible away from the adjacent nodes.

For example, if MSN equals 137 and there are two nodes labelled 0 and 68. Then next server registered can be labeled 102 or 103.

Both registration and deregistration needs a shared "register" lock with the help of zookeeper. It's implemented in the same way as "master" lock mentioned above.
#### Data
Data are stored as key value pairs in files. Both key and value have the type of string. Data file operations can be found in `src/server/utils/data.go`.
Normally a key value pair is only stored in one active data server and some of the standby data servers. 

Data are distributed to different data node by key. A basic hash method is applied to key and it determines which server the key goes to based on consistent hash. 

To make it more clear I'll give a simple example here. Suppose MSN is 137(default value), and two nodes are labelled 0 and 68. 
A key has a hash value of 1395665. Then divide the hash value by MSN and the remainder is 46. It's closer to 68 then to 0, so the key is stored in the server labelled 68.

The data are also distributed in the server itself to make the whole system more scalable. 

For example, one level storage has only one file to store the data while two level storage can have 3*5 files.
The numbers of directories/files in each level are increasing prime numbers. In a two level storage server with 3*5 files, the key with hash value 1395665 will be stored in file `2/0` because 1395665 mod 3 = 2 and 1395665 mod 5 = 0.

We can prove that each key has the same probability to be stored in each file by chinese remainder theorem.
#### Standby master node
There is only one available master node at any time, the one keeping the "master" lock implemented with zookeeper. 
The lock details can be found in `src/zookeeper/lock.go`. 

Standby master will try to grab the lock when starting up and get the lock when the old master exits.

The client end can specify multiple master addresses. There is a goroutine running to ping the master every 300ms. 
If it can't ping the master, the client will try to connect other master addresses. 
#### Standby data node
Standby data node create a `sb` znode when registering. And it will store all the latest data after the startup. It won't provide backup for earlier data because it's too costly to block all servers to transfer data.

In this way, if data is created before any standby node, it could be lost if the server faces something terrible like power failure. The standby node can register itself as data node after a data server passes out, which means its znode disappeared but the deregistration function is not executed.

This is detected by master node, which get the notification of data server deregistration. 

When a data znode disappear and the master didn't get the notification, the master will delete one of the standby data znode. And then that standby node knows he is appointed as new data server.
He will register himself as data server, and remove all the data that shouldn't be stored in this server anymore.
#### Operations
Four client operations are provided: get, put, delete, dump. Dump is only used for debugging purpose which shows the data storage information of all servers.
See `src/client/op` for more information. 

Get operation will ask master for address of data server and then send a get request to it. Put and delete operations work in similar way. Dump can be seen as advanced get operation. Each put or delete operation will be sent to standby node as well.

Concurrent data accessing and strong consistency is supported. There is a read write lock for each data file. Data locks are file-level, which means operations on two keys stored in the same file will be sequential.
This is more scalable than key-level lock and performs better than server-level lock. And more storage levels can help improve the performance.

### Practice
Here I'll list some practice and also tests on the key value system.

First establish the zookeeper cluster. (ps: zk directory is not submitted)

```make zk```

This will start all the zookeeper nodes in `zk` directory.

Then generate the executable files.

```make build```

In `scripts/run.sh` it has shown some ways to run the servers or the clients.

To run a master:

```
cd src/cmd/server
go run server_main.go -cfg=../cfg.yaml -type=master -addr=:9666
# or
# cd build
# ./server -cfg=cfg.yaml -type=master -addr=:9666
```

Run another one as standby:

```
cd src/cmd/server
go run server_main.go -cfg=../cfg.yaml -type=master -addr=:9667
# or
# cd build
# ./server -cfg=cfg.yaml -type=master -addr=:9667
```

Run two data servers on random port:
```
cd src/cmd/server
go run server_main.go -cfg=../cfg.yaml -type=slave -addr=127.0.0.1: -data=data/server1
go run server_main.go -cfg=../cfg.yaml -type=slave -addr=127.0.0.1: -data=data/server2
# or
# cd build
# ./server -cfg=cfg.yaml -type=slave -addr=127.0.0.1: -data=data/server1
# ./server -cfg=cfg.yaml -type=slave -addr=127.0.0.1: -data=data/server2
```

Run two standby data node:
```
cd src/cmd/server
go run server_main.go -cfg=../cfg.yaml -type=slave-sb -addr=127.0.0.1: -data=data/sb1
go run server_main.go -cfg=../cfg.yaml -type=slave-sb -addr=127.0.0.1: -data=data/sb2
# or
# cd build
# ./server -cfg=cfg.yaml -type=slave-sb -addr=127.0.0.1: -data=data/sb1
# ./server -cfg=cfg.yaml -type=slave-sb -addr=127.0.0.1: -data=data/sb2
```

I have written a test client end and also an interactive shell client.
To run the shell and do some operations:
```
cd src/cmd/client-shell
go run shell.go -addr=:9666,127.0.0.1:9667
# or
# cd build
# ./shell -addr=:9666,127.0.0.1:9667
```
```
This is a distributed key value system.
Supported operations include get, put, del, dump(only for debugging).
kv@ds$ get a
not found
kv@ds$ put a 1
created
kv@ds$ del a
deleted
kv@ds$ get a
not found
kv@ds$ del a
not found
kv@ds$ put b 2
created
kv@ds$ put c 3
created
kv@ds$ put e 4
created
kv@ds$ dump
-------Data server 127.0.0.1:43309 with label 0-------
    key: c value: 3 label: 34
    key: e value: 4 label: 2

-------Data server 127.0.0.1:39755 with label 68-------
    key: b value: 2 label: 85

kv@ds$ 
```

Test the sequential operations:
```
cd src/cmd/client
go run client_main.go test.go -op=sequential -addr=:9666,127.0.0.1:9667
# or 
# cd build
# ./client -op=sequential -addr=:9666,127.0.0.1:9667
# result:
get key: os, err: <nil>, val: 100
del key: ds err: <nil>
put key: os val: 100 err: <nil>
put key: ds val: 98 err: <nil>
put key: ca val: 97 err: <nil>
put key: st val: 96 err: <nil>
del key: st err: <nil>
get key: os, err: <nil>, val: 100
get key: ds, err: <nil>, val: 98
get key: ca, err: <nil>, val: 97
get key: st, err: not found, val: 
DUMPING ALL!!!
-------Data server 127.0.0.1:43309 with label 0-------
    key: ds value: 98 label: 105
    key: ca value: 97 label: 130
    key: os value: 100 label: 134
    key: c value: 3 label: 34
    key: e value: 4 label: 2

-------Data server 127.0.0.1:39755 with label 68-------
    key: b value: 2 label: 85
```


Test the concurrent operations:
```
cd src/cmd/client
go run client_main.go test.go -op=concurrent -addr=:9666,127.0.0.1:9667
# or 
# cd build
# ./client -op=sequential -addr=:9666,127.0.0.1:9667
# result:
put key: os val: 100 err: <nil>
2 get key: os, err: <nil>, val: 100
1 get key: os, err: <nil>, val: 100
put key: os val: 99 err: <nil>
put key: os val: 98 err: <nil>
3 get key: os, err: <nil>, val: 99
put key: os val: 97 err: <nil>
1 get key: os, err: <nil>, val: 97
2 get key: os, err: <nil>, val: 97
3 get key: os, err: <nil>, val: 97
1 get key: os, err: <nil>, val: 97
2 get key: os, err: <nil>, val: 97
3 get key: os, err: <nil>, val: 97
1 get key: os, err: <nil>, val: 97
2 get key: os, err: <nil>, val: 97
3 get key: os, err: <nil>, val: 97
1 get key: os, err: <nil>, val: 97
2 get key: os, err: <nil>, val: 97
DUMPING ALL!!!
-------Data server 127.0.0.1:43309 with label 0-------
    key: ds value: 98 label: 105
    key: ca value: 97 label: 130
    key: os value: 97 label: 134
    key: c value: 3 label: 34
    key: e value: 4 label: 2

-------Data server 127.0.0.1:39755 with label 68-------
    key: b value: 2 label: 85
```

Test the standby:
```
cd src/cmd/client
go run client_main.go test.go -op=crazyloop -addr=:9666,127.0.0.1:9667
# or 
# cd build
# ./client -op=crazyloop -addr=:9666,127.0.0.1:9667
# kill one data server
# and then a standby server will turn to data server

# kill the master run on 9666
# after a few rpc failure in clients(300ms interval)
# the client request turns back to normal 
# and the standby master run on 9777 is the new master now
```

