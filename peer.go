/*
* @Author: andrea
* @Date:   2018-02-08 11:18:43
* @Last Modified by:   Andrea Golin
* @Last Modified time: 2018-02-14 16:59:09
 */

package lockan

import (
	"bufio"
	//"encoding/binary"
	"bytes"
	"fmt"
	"log"
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
	commands   map[int]*Command
}

type peer interface {
	Init() int64
	inputLoop(reader *bufio.Reader)
	print()
	listCommands()
}

func (p Peer) listCommands() {
	for _, value := range p.commands {
		log.Println(value)
	}
}

func (p Peer) Init() int64 {

	log.Printf("%s", "Starting...")

	/**
	 * Variable init
	 */
	p.quit = make(chan bool)
	p.printer = make(chan string)
	reader := bufio.NewReader(os.Stdin)
	InitCommandsList()

	/**
	 * Start listening
	 * @type net.Listener
	 */
	sPort := strconv.Itoa(p.port)
	l, err := net.Listen("tcp", "localhost:"+sPort)
	if err != nil {
		panic(err)
	}
	log.Println("Listening to localhost " + sPort)

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
		textSlice := strings.Fields(text)

		/**
		 * TODO: to be rewritten with switch/case
		 * TODO: change commands input from stdin to TCP socket
		 */
		if strings.Compare("ping", text) == 0 {
			p.printer <- "pong"
		} else if strings.Compare("stop", text) == 0 {
			p.printer <- "resquested Shutdown"
			p.quit <- true
		} else if strings.Compare("commands-list", text) == 0 {
			for _, value := range p.commands {
				log.Printf("%d: %s \n", value.id, value.name)
			}
		} else if strings.Compare("lok", textSlice[0]) == 0 {
			ParseCommands(textSlice[1:])
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
		log.Println(msg)
	}
}

/**
 * [handleRequest description]
 * @param  {[type]} conn    net.Conn      [description]
 * @param  {[type]} printer chan          string        [description]
 * @return {[type]}         [description]
 *
 * TODO: implement REAL netcode
 */
func (p Peer) handleRequest(conn net.Conn) {

	buf := make([]byte, 35)

	_, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}

	/**
	 * Logging to console for debug
	 */
	log.Printf("%s", "Received connection")

	ParseNetInput(buf)

	conn.Write([]byte("Message corectly received: \n"))
	conn.Close()

}

func DummySend(dstPort string) {

	log.Printf("%s: %s", "Dest port:", dstPort)

	conn, err := net.Dial("tcp", "localhost:"+dstPort)
	if err != nil {
		log.Printf("%s", err)
		panic("Ping Error")
	}

	defer conn.Close()

	/*pSize    [2]byte
	pType    [1]byte
	pPayload [32]byte*/

	var pSize [2]byte
	var pType [1]byte

	pSize[0] = 8
	pSize[1] = 16
	pType[0] = 2
	dummyPayload := []byte("ASD")

	log.Printf("%s: %T", "Dummy size type", pSize)
	log.Printf("%s: %T", "Dummy type type", pType)

	s := [][]byte{pSize[:], pType[:], dummyPayload}
	fmt.Printf("%b", s)
	cs := bytes.Join(s, []byte(""))
	conn.Write(cs)

}

/**
 * TODO: rewrite bootstrap logic
 * TODO: rewrite command logic
 */
func Start(port *int) {
	peers := make(map[string]Peer)
	scan := &Command{1, "Scan", "Scan"}
	commands := make(map[int]*Command)
	commands[1] = scan
	peer := &Peer{debug: true, maxpeers: 10, serverport: 678, serverhost: "localhost", peers: peers, port: *port, commands: commands}
	peer.Init()
}
