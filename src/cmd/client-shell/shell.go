package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/isalb729/ds-kv/src/client"
	"os"
	"strings"
)

func main() {
	addr := flag.String("addr", "", "master or slave listening address")
	flag.Parse()
	addrList := strings.Split(*addr, ",")
	for _, addr := range addrList {
		if (addr)[0] == ':' {
			addr = "127.0.0.1" + addr
		}
	}
	cli := client.Connect(addrList)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("This is a distributed key value system.")
	fmt.Println("Supported operations include get, put, del, dump(only for debugging).")
	for {
		fmt.Print("kv@ds$ ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		err = runCommand(cmdString, cli)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
func runCommand(commandStr string, cli *client.KvCli) error {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)
	if len(arrCommandStr) == 0 {
		return nil
	}
	// Operations.
	switch arrCommandStr[0] {
	case "exit":
		os.Exit(0)
	case "put":
		if len(arrCommandStr) != 3 {
			fmt.Fprintf(os.Stdout, "WRONG FORMAT\n")
		}
		err, created := cli.Put(arrCommandStr[1], arrCommandStr[2])
		if err != nil {
			fmt.Println(err)
		} else if created {
			fmt.Println("created")
		} else {
			fmt.Println("updated")
		}
	case "get":
		if len(arrCommandStr) != 2 {
			fmt.Fprintf(os.Stdout, "WRONG FORMAT\n")
		}
		err, val := cli.Get(arrCommandStr[1])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(val)
		}
	case "del":
		if len(arrCommandStr) != 2 {
			fmt.Fprintf(os.Stdout, "WRONG FORMAT\n")
		}
		err, deleted := cli.Del(arrCommandStr[1])
		if err != nil {
			fmt.Println(err)
		} else if deleted {
			fmt.Println("deleted")
		} else {
			fmt.Println("not found")
		}
	case "dump":
		if len(arrCommandStr) != 1 {
			fmt.Fprintf(os.Stdout, "WRONG FORMAT\n")
		}
		err, rsp := cli.DumpAll()
		if err != nil {
			fmt.Println(err)
		} else {
			for _, data := range rsp {
				fmt.Printf("-------Data server %s with label %d-------\n", data.Host, data.Label)
				for _, kvl := range data.Kvls {
					fmt.Printf("    key: %s value: %s label: %d\n", kvl.Key, kvl.Value, kvl.Label)
				}
				fmt.Println()
			}
		}
	default:
		fmt.Fprintf(os.Stdout, "UNKNOWN COMMAND %s\n", arrCommandStr[0])
	}
	return nil
}
