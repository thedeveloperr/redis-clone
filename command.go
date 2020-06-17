package main

import (
	"strconv"
	"strings"
)

type Command struct {
	fullText string
}

func (c Command) parse() (commandType string, key string, parsedArguments [][2]string) {
	commandComponents := strings.Fields(c.fullText)
	if commandComponents[0] == "GET" && len(commandComponents) == 2 {
		commandType = "GET"
		key = commandComponents[1]
		return
	}
	if commandComponents[0] == "SET" && len(commandComponents) == 3 {
		commandType = "SET"
		key = commandComponents[1]
		parsedArguments = [][2]string{
			{commandComponents[2], ""},
		}
		return
	}

	if commandComponents[0] == "EXPIRE" && len(commandComponents) == 3 {
		if _, err := strconv.ParseInt(commandComponents[2], 10, 64); err != nil {
			return
		}
		commandType = "EXPIRE"
		key = commandComponents[1]
		parsedArguments = [][2]string{
			{commandComponents[2], ""},
		}
		return
	}
	if commandComponents[0] == "ZRANK" && len(commandComponents) == 3 {
		commandType = "ZRANK"
		key = commandComponents[1]
		parsedArguments = [][2]string{
			{commandComponents[2], ""},
		}
		return
	}
	if commandComponents[0] == "ZADD" && len(commandComponents) >= 4 &&
		len(commandComponents)%2 == 0 {
		commandType = "ZADD"
		key = commandComponents[1]
		for i := 2; i < len(commandComponents); i = i + 2 {
			parsedArguments = append(parsedArguments, [2]string{
				commandComponents[i],
				commandComponents[i+1],
			})
		}
		return
	}
	if commandComponents[0] == "ZRANGE" && len(commandComponents) >= 4 {
		commandType = "ZRANGE"
		key = commandComponents[1]
		if len(commandComponents) == 5 && commandComponents[4] == "WITHSCORES" {
			parsedArguments = [][2]string{
				{commandComponents[2], commandComponents[3]},
				{commandComponents[4], ""},
			}

		} else if len(commandComponents) == 4 {
			parsedArguments = [][2]string{
				{commandComponents[2], commandComponents[3]},
			}
		}
		return
	}
	return

}
