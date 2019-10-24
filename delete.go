package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mkideal/cli"
)

type deleteT struct {
	cli.Helper
	Port       int    `cli:"*p,port" usage:"Specify the port to apply the wire to"`
	OutputFile string `cli:"o,output" usage:"Specify the logfile" dft:"/var/log/<ChainName>"`
}

var deleteCMD = &cli.Command{
	Name: "delete",
	Desc: "Deletes a tripwire chain",
	Argv: func() interface{} { return new(deleteT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*deleteT)
		ChainName := generateChainname(argv.Port)
		if chainExisits(ChainName) != nil {
			fmt.Println("Chain doesn't exist")
			return nil
		}
		runCommand(errorHandler, "iptables -F "+ChainName)
		deleteRuleForChain(errorHandler, ChainName)
		runCommand(errorHandler, "iptables -X "+ChainName)
		runCommand(errorHandler, "rm /etc/rsyslog.d/"+ChainName+".conf")
		runCommand(errorHandler, "systemctl restart rsyslog.service")
		if argv.OutputFile != "/var/log/<ChainName>" {
			if _, err := os.Stat("/var/log/" + argv.OutputFile); err == nil {
				runCommand(errorHandler, "rm /var/log/"+argv.OutputFile)
				fmt.Println("Deleted logfile /var/log/" + argv.OutputFile)
			}
		}
		fmt.Println("Deleted chain " + ChainName + " successfully")
		return nil
	},
}

func deleteRuleForChain(errorHandler func(error, string), chainName string) {
	data, _ := runCommand(errorHandler, "iptables -L INPUT --line-numbers")
	lines := strings.Split(data, "\n")
	for _, i := range lines {
		if strings.HasPrefix(i, "num  target") || strings.HasPrefix(i, "Chain") || len(strings.Trim(i, " ")) == 0 {
			continue
		}
		e := strings.Trim(strings.Split(i, " ")[0], " ")
		in, err := strconv.Atoi(e)
		if err != nil {
			fmt.Println("Couldn't delete rule")
			continue
		}
		if strings.Contains(i, chainName) {
			fmt.Println("Deleting rule " + i)
			runCommand(errorHandler, "iptables -D INPUT "+strconv.Itoa(in))
			return
		}
	}
}
