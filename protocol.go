/*
* @Author: Andrea Golin
* @Date:   2018-02-14 10:50:39
* @Last Modified by:   Andrea Golin
* @Last Modified time: 2018-02-14 17:03:32
 */
package lockan

import (
	"bytes"
	"encoding/hex"
	"log"
)

type LockPacket struct {
	pSize    [2]byte
	pType    [1]byte
	pPayload [32]byte
}

type iLockPacket interface {
	Dump()
	Compose(pSize [2]byte, pType [1]byte, pPayload [32]byte) *LockPacket
}

/**
 * Dump packet
 *
 * @param  {[type]} p *LockPacket)  Dump( [description]
 * @return {[type]}   [description]
 */
func (p *LockPacket) Dump() {
	s := [][]byte{p.pSize[:], p.pType[:], p.pPayload[:]}
	d := bytes.Join(s, []byte(""))
	log.Printf("%s: %T", "Packet type", p)
	log.Printf("%s: %T", "Packet size", p.pSize)
	log.Printf("%s: %T", "Packet type", p.pType)
	log.Printf("%s: %b", "Packet payload", p.pPayload)
	pPSlice := p.pPayload[:]
	log.Printf("%s%s", "Packet payload hex dump \n", hex.Dump(pPSlice))
	log.Printf("%s: %b", "Package complete bytes", d)
	log.Printf("%s%s", "Packet payload hex dump \n", hex.Dump(d))
}

/**
 * Parse net input and fetch instructions
 */
func ParseNetInput(buf []byte) {
	log.Printf("%s %b", "Parsing net data...", buf)
	pSize := buf[:1]
	log.Printf("%s: %b", "Size", pSize)
	pType := buf[2]
	log.Printf("%s: %d", "Type", pType)
	pPayload := buf[3:]
	log.Printf("%s: %b", "Payload", pPayload)
	log.Printf("%s: %s", "Payload string parse", string(pPayload))
}

/**
 *
 */
func (p *LockPacket) Compose(pSize [2]byte, pType [1]byte, pPayload [32]byte) *LockPacket {
	p.pSize = pSize
	p.pType = pType
	return p
}
