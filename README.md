# ds-kv
TODO:
get familiar with zk and zk client
trry cluster
finish basic
unit test
finish advanced
unit test
write shell client
loadtest


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

