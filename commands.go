/*
* @Author: Andrea Golin
* @Date:   2018-02-09 17:43:51
* @Last Modified by:   Andrea Golin
* @Last Modified time: 2018-02-13 10:17:27
 */
package lockan

import (
	"fmt"
)

type Command struct {
	id       int
	name     string
	function string
}

type iCommand interface {
	Run()
	Status()
	Stop()
	Parse()
}

func ParseCommands(values []string) {
	for _, value := range values {
		fmt.Printf("%q \n", value)
	}
}
