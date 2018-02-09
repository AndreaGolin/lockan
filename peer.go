/*
* @Author: andrea
* @Date:   2018-02-08 11:18:43
* @Last Modified by:   Andrea Golin
* @Last Modified time: 2018-02-09 16:56:02
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

	/**
	 * Start side thread for:
	 * 	printing
	 * 	pooling stdin
	 */
	go print(p.printer)
	go inputLoop(p.printer, reader, p.quit)

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
		go handleRequest(conn, p.printer)
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
func inputLoop(printer chan string, reader *bufio.Reader, quit chan bool) {
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare("ping", text) == 0 {
			printer <- "pong"
		} else if strings.Compare("stop", text) == 0 {
			printer <- "resquested Shutdown"
			quit <- true
		}
	}
}

/**
 * [print description]
 * @param  {[type]} print chan          string [description]
 * @return {[type]}       [description]
 */
func print(print chan string) {
	for {
		msg := <-print
		fmt.Println(msg)
	}
}

/**
 * [handleRequest description]
 * @param  {[type]} conn    net.Conn      [description]
 * @param  {[type]} printer chan          string        [description]
 * @return {[type]}         [description]
 */
func handleRequest(conn net.Conn, printer chan string) {

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}

	printer <- "Received connection"

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
