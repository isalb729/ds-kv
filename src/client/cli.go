package client

import "github.com/isalb729/ds-kv/src/rpc"

func Concurrent(cli *KvCli) {

}

func Sequential(cli *KvCli) {
	get()
	put()
	del()
}
