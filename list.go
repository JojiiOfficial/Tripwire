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

var list = &cli.Command{
	Name: "list",
	Desc: "Lists all tripwires",
	Fn: func(ctx *cli.Context) error {
		list, err := runCommand(errorHandler, "iptables -L")
		if err != nil {
			fmt.Println("Error viewing all chains")
			return nil
		}
		lines := strings.Split(list, "\n")
		if len(lines) == 0 {
			fmt.Println("No tripwire chain found")
			return nil
		}
		for _, line := range lines {
			if strings.HasPrefix(line, "Chain Tripwire") {
				viewChain(strings.Trim(strings.Split(line, " ")[1], " "))
			}
		}
		return nil
	},
}

func viewChain(chainName string) {
	port, err := strconv.Atoi(strings.ReplaceAll(chainName, "Tripwire", ""))
	if err != nil {
		return
	}
	ChainName := generateChainname(port)
	LogFile, err := readFile("/etc/rsyslog.d/" + ChainName + ".conf")
	LogFile = strings.Trim(strings.Split(LogFile, "then")[1], " ")
	fmt.Println("Tripwire on port " + strconv.Itoa(port) + "\t--->  " + LogFile)
}

func readFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(data), "\n", ""), nil
}
