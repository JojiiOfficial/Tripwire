package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/mkideal/cli"
)

type childT struct {
	cli.Helper
}
type argList struct {
	cli.Helper
	Port int `cli:"*p,port" usage:"Shows the chain to a given port" dft:"-1"`
}

var list = &cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Desc:    "Lists all tripwires",
	Argv:    func() interface{} { return new(argList) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argList)
		list, err := runCommand(errorHandler, "iptables -L")
		if err != nil {
			fmt.Println("Error viewing all chains")
			return nil
		}
		lines := strings.Split(list, "\n")
		fmt.Println("\033[1;32m\tPort\033[0m  \t\t\t\033[1;32m Logfile\033[0m")
		len := 0
		for _, line := range lines {
			if strings.HasPrefix(line, "Chain Tripwire") {
				viewed := viewChain(strings.Trim(strings.Split(line, " ")[1], " "), argv.Port)
				if viewed {
					len++
				}
			}
		}
		if len == 0 {
			fmt.Println("\t---- No tripwire chain found ----\n")
			return nil
		}
		return nil
	},
}

func viewChain(chainName string, portFilter int) bool {
	port, err := strconv.Atoi(strings.ReplaceAll(chainName, "Tripwire", ""))
	if err != nil {
		return false
	}
	if portFilter != -1 && port != portFilter {
		return false
	}
	ChainName := generateChainname(port)
	LogFile, err := readFile("/etc/rsyslog.d/" + ChainName + ".conf")
	LogFile = strings.Trim(strings.Split(LogFile, "then")[1], " ")
	fmt.Println("Tripwire on port " + strconv.Itoa(port) + "\033[1;37m\t--->\033[0m  " + LogFile)
	return true
}

func readFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(data), "\n", ""), nil
}
