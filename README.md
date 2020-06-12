# ds-kv
TODO:
1. zk cluster
2. zk client, create delete get, ephemeral, sequential
3. server establish
4. server basic
unit test
5. server cluster
6. finish advanced
unit test
7. write shell client
8. loadtest


## zookeeper deployment
download
https://mirrors.sonic.net/apache/zookeeper/zookeeper-3.6.1/apache-zookeeper-3.6.1-bin.tar.gz
mkdir data
zk-1/conf/zoo.cfg

tickTime = 2000
dataDir = /path/to/zookeeper/data
clientPort = 12181
initLimit = 5
syncLimit = 2

bin/zkServer.sh start

