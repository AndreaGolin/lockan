/*
* @Author: Andrea Golin
* @Date:   2018-02-14 09:10:17
* @Last Modified by:   Andrea Golin
* @Last Modified time: 2018-02-14 09:12:40
 */
package main

import (
	"flag"
	"github.com/AndreaGolin/lockan"
	"log"
)

func main() {
	log.Printf("%s", "Starting bootstrap file...")
	prtNmb := flag.Int("port", 678, "an int")
	flag.Parse()
	log.Printf("%s %d", "Given port: ", prtNmb)
	lockan.Start(prtNmb)
}
