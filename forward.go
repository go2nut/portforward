package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var source = flag.String("source", "0.0.0.0:80", "source host:port")
var dest = flag.String("dest", "127.0.0.1:8000", "dest host:port")

func main() {
	flag.Parse()
	forwardPort(*source, *dest)
}

func forwardPort(source, dest string) {
	l, err := net.Listen("tcp", source)
	if err != nil {
		fmt.Println(err, err.Error())
		os.Exit(0)
	}

	fmt.Println(fmt.Sprintf("start forward traffic from %s to %s.", source, dest))

	for {
		sConn, err := l.Accept()
		if err != nil {
			continue
		}

		dTcpAddr, _ := net.ResolveTCPAddr("tcp4", dest)
		dConn, err := net.DialTCP("tcp", nil, dTcpAddr)
		if err != nil {
			fmt.Println(err)
			sConn.Write([]byte(fmt.Sprintf("can't connect %s", dest)))
			sConn.Close()
			continue
		}
		go io.Copy(sConn, dConn)
		go io.Copy(dConn, sConn)
	}

	fmt.Println("quit port forward.")
}
