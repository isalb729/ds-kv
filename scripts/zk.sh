#!/bin/bash
ZK=../zk
ZKLIST=$(find $ZK/* -maxdepth 0 -type d | sed  "s/^.*\///")
if [ -z "$1" ]
then
  CMD=restart
else
  CMD=$1
fi
for zk in $ZKLIST
do
    $ZK/$zk/bin/zkServer.sh $CMD
done
if [[ ($CMD = restart) || ($CMD = start) ]]
then
  ZKLIST=($ZKLIST)
  PS=$(ps -ef | grep zk- | grep -v grep | awk '{print $2}')
  if [ ${#PS[@]} -ne ${#ZKLIST[@]} ]
  then
    echo "FAIL; only ${#PS[@]} zks are started"
  else
    echo SUCCESS
  fi
fi


