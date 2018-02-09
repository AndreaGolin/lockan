/*
* @Author: andrea
* @Date:   2018-02-08 11:18:43
* @Last Modified by:   Andrea Golin
* @Last Modified time: 2018-02-09 17:38:23
 */

package lockan

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	//"time"
)

type Peer struct {
	debug      bool
	maxpeers   int64
	serverport int64
	port       int
	serverhost string
	peers      map[string]Peer
	quit       chan bool
	printer    chan string
}

type peer interface {
	Init() int64
	inputLoop(reader *bufio.Reader)
	print()
}

func (p Peer) Init() int64 {

	/**
	 * Variable init
	 */
	p.quit = make(chan bool)
	p.printer = make(chan string)
	reader := bufio.NewReader(os.Stdin)

	/**
	 * Start listening
	 * @type net.Listener
	 */
	sPort := strconv.Itoa(p.port)
	l, err := net.Listen("tcp", "localhost:"+sPort)
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening to localhost " + sPort)

	var tcpConn net.TCPConn
	fmt.Println(tcpConn)

	/**
	 * Start side thread for:
	 * 	printing
	 * 	pooling stdin
	 */
	go p.print()
	go p.inputLoop(reader)

	defer l.Close()

	/**
	 * Wait for quit channel to fire in background
	 * @return error
	 */
	go func() {
		for {
			select {
			case <-p.quit:
				l.Close()
			}
		}
	}()

	/**
	 * Start looping for tcp connection
	 */
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go p.handleRequest(conn)
	}

	return 1
}

/**
 * [inputLoop description]
 * @param  {[type]} printer chan          string        [description]
 * @param  {[type]} reader  *bufio.Reader [description]
 * @param  {[type]} quit    chan          bool          [description]
 * @return {[type]}         [description]
 */
func (p Peer) inputLoop(reader *bufio.Reader) {
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare("ping", text) == 0 {
			p.printer <- "pong"
		} else if strings.Compare("stop", text) == 0 {
			p.printer <- "resquested Shutdown"
			p.quit <- true
		} else {
			p.printer <- "received input \n"
		}
	}
}

/**
 * [print description]
 * @param  {[type]} print chan          string [description]
 * @return {[type]}       [description]
 */
func (p Peer) print() {
	for {
		msg := <-p.printer
		fmt.Println(msg)
	}
}

/**
 * [handleRequest description]
 * @param  {[type]} conn    net.Conn      [description]
 * @param  {[type]} printer chan          string        [description]
 * @return {[type]}         [description]
 */
func (p Peer) handleRequest(conn net.Conn) {

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}

	p.printer <- "Received connection"

	conn.Write([]byte("Message corectly received: \n"))
	conn.Close()

}

/**
 * [Start description]
 */
func Start(port *int) {
	peers := make(map[string]Peer)
	peer := &Peer{debug: true, maxpeers: 10, serverport: 678, serverhost: "localhost", peers: peers, port: *port}
	peer.Init()
}
