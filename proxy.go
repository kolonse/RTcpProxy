// proxy project proxy.go
package main

import (
	"flag"
	"fmt"
	"net"
)

var local = flag.String("l", "127.0.0.1:8001", "-l=<127.0.0.1:8001> local listen addr")
var remote = flag.String("r", "127.0.0.1:8002", "-r=<127.0.0.1:8002> remote proxy addr")

func main() {
	flag.Parse()
	l, e := net.Listen("tcp", *local)
	if e != nil {
		panic(e)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go func(c net.Conn) {
			cc, err1 := net.Dial("tcp", *remote)
			defer c.Close()
			if err1 != nil {
				fmt.Printf("%v proxy to %v except:%v\n", c.RemoteAddr(), cc.LocalAddr(), err1.Error())
				return
			}
			e := make(chan bool, 1)
			f := func(c1, c2 net.Conn) {
				buff := make([]byte, 1024)
				for {
					n, err2 := c1.Read(buff)
					if err2 != nil {
						fmt.Printf("%v data to %v except:%v\n", c1.RemoteAddr(), c2.LocalAddr(), err2.Error())
						break
					}
					c2.Write(buff[:n])
				}
				c1.Close()
				c2.Close()
				e <- true
			}
			go f(c, cc)
			go f(cc, c)
			<-e
			<-e
		}(conn)
	}
}
