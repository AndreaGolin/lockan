/*
* @Author: andrea
* @Date:   2018-02-08 11:18:43
* @Last Modified by:   Andrea Golin
* @Last Modified time: 2018-02-08 13:30:46
 */

package lockan

import (
	"fmt"
	"net"
)

type Peer struct {
	debug      bool
	maxpeers   int64
	serverport int64
	serverhost string
	peers      map[string]Peer
}

type peer interface {
	init() int64
}

func (p Peer) init() int64 {
	l, err := net.Listen("tcp", "localhost:678")
	if err != nil {
		panic(err)
	}

	defer l.Close()

	fmt.Println("Listening on localhost, 678 port")

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go handleRequest(conn)
	}

	return 10
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	conn.Write([]byte("Message received: \n"))
	conn.Write(buf)
	conn.Write([]byte("\n"))
	conn.Close()
}

func Test() {
	peers := make(map[string]Peer)
	peer := &Peer{debug: true, maxpeers: 10, serverport: 666, serverhost: "localhost", peers: peers}
	fmt.Println(peer.init())
}
