/*
* @Author: andrea
* @Date:   2018-02-08 11:18:43
* @Last Modified by:   Andrea Golin
* @Last Modified time: 2018-02-09 14:35:25
 */

package lockan

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Peer struct {
	debug      bool
	maxpeers   int64
	serverport int64
	serverhost string
	peers      map[string]Peer
	printer    chan string
	lifelink   chan int
}

type peer interface {
	init() int64
}

func (p Peer) init() int64 {
	p.printer = make(chan string)
	p.lifelink = make(chan int)

	go mainPeerLoop(p.printer, p.lifelink)
	go print(p.printer)

	reader := bufio.NewReader(os.Stdin)
	go inputLoop(p.printer, reader, p.lifelink)

	l, err := net.Listen("tcp", "localhost:678")
	if err != nil {
		panic(err)
	}

	defer l.Close()

	for {
		_, err := l.Accept()
		if err != nil {
			panic(err)
		}

		d := <-p.lifelink
		if d == 1 {
			fmt.Println("received kill imput")
			l.Close()
			break
		}
		/*go handleRequest(conn, p.printer, p.lifelink, l)*/

	}

	return 1
}

func inputLoop(printer chan string, reader *bufio.Reader, lifelink chan int) {
	for {
		/*text, _ := reader.ReadString('\n')*/
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare("hi", text) == 0 {
			printer <- "hi to you!"
		} else {
			printer <- "input given"
			lifelink <- 1
		}
	}
}

func print(print chan string) {
	for {
		msg := <-print
		fmt.Println(msg)
	}
}

func handleRequest(conn net.Conn, printer chan string, lifelink chan int, l net.Listener) {

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}

	conn.Write([]byte("Message corectly received: \n"))

	message := "Messaggio"
	printer <- message

	for {
		select {
		case <-lifelink:
			fmt.Println("received sigkill")
			conn.Close()
			l.Close()
		}
	}

	conn.Close()

}

func listenForStdin(reader *bufio.Reader) {
	fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	if strings.Compare("hi", text) == 0 {
		fmt.Println("hello, Yourself")
	} else {
		fmt.Println(text)
	}
}

func mainPeerLoop(printer chan string, lifelink chan int) {
	for {
		select {
		/*case <-lifelink:
		message := "Kill!"
		printer <- message*/

		}
	}
}

func Test() {
	peers := make(map[string]Peer)
	peer := &Peer{debug: true, maxpeers: 10, serverport: 666, serverhost: "localhost", peers: peers}
	fmt.Println(peer.init())
}
