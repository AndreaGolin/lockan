/*
* @Author: Andrea Golin
* @Date:   2018-02-09 17:43:51
* @Last Modified by:   Andrea Golin
* @Last Modified time: 2018-02-14 16:44:33
 */
package lockan

import (
	"log"
	"strings"
)

type Command struct {
	id       int
	name     string
	function string
}

/**
 * @brief      List of commands.
 */
type CommandList struct {
	id          int
	name        string
	commands    map[int]*Command
	status      int32
	commandsDef map[int]string
}

/**
 * @brief      { function_description }
 *
 * @return     { description_of_the_return_value }
 */
func InitCommandsList() {
	commandsDef := make(map[int]string)

	commandsDef[0] = "scan"
	commandsDef[1] = "ping"
	commandsDef[2] = "connect"
	commandsDef[3] = "cut"

	commands := make(map[int]*Command)
	CommandList := &CommandList{id: 1, name: "Test", commands: commands, status: 1, commandsDef: commandsDef}

	log.Printf("%s | %v", "Finished initializing commands list", CommandList)
}

/**
 * @brief      { function_description }
 *
 * @param      values  The values
 *
 * @return     { description_of_the_return_value }
 */
func ParseCommands(values []string) {
	for _, value := range values {
		log.Printf("%s: %q \n", "Received command", value)
		if strings.Compare("ping", value) == 0 {
			DummySend()
		}
	}
}
