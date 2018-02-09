/*
* @Author: Andrea Golin
* @Date:   2018-02-09 17:43:51
* @Last Modified by:   Andrea Golin
* @Last Modified time: 2018-02-09 17:52:24
 */
package lockan

type Command struct {
	id       int
	name     string
	function string
}

type iCommand interface {
	Run()
	Status()
	Stop()
}
