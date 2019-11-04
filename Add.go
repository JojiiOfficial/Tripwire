package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mkideal/cli"
)

type addT struct {
	cli.Helper
	Accept     bool   `cli:"a,accept" usage:"Specify wether to drop or accept the incoming connections"`
	Port       int    `cli:"*p,port" usage:"Specify the port to apply the wire to"`
	OutputFile string `cli:"o,output" usage:"Specify the logfile" dft:"/var/log/<ChainName>"`
	LogLevel   int    `cli:"l,log-level" usage:"Specify the log level" dft:"6"`
}

var addCMD = &cli.Command{
	Name:    "add",
	Aliases: []string{"a", "ad"},
	Desc:    "Adds a tripwire chain",
	Argv:    func() interface{} { return new(addT) },
	Fn: func(ctx *cli.Context) error {
		if !checkRoot() {
			return nil
		}
		argv := ctx.Argv().(*addT)
		ChainName := generateChainname(argv.Port)
		LogIdentifier := ChainName + " "
		ruleAction := "ACCEPT"
		if !argv.Accept {
			ruleAction = "DROP"
		}
		if argv.Port <= 0 || argv.Port > 65535 {
			fmt.Println("You port must be between 0 and 65535!")
			return nil
		}
		if chainExisits(ChainName) == nil {
			fmt.Println("This port already has a rule! Try deleting it with -d")
			return nil
		}
		if argv.Port == 22 && !argv.Accept {
			reader := bufio.NewReader(os.Stdin)
			y, i := confirmInput("If you continue iptables will block port 22. If you are currently logged in via SSH on port 22 you won't be able to use SSH!\nContinue? [y/n] > ", reader)
			if i != 1 || !y {
				fmt.Println("Exiting")
				return nil
			}
		}
		outFile := argv.OutputFile
		if argv.OutputFile == "/var/log/<ChainName>" {
			outFile = "/var/log/" + ChainName
		}
		if !strings.HasPrefix(outFile, "/") {
			outFile = "/var/log/" + outFile
		}
		runCommand(errorHandler, "iptables -N "+ChainName)
		runCommand(errorHandler, "iptables -A "+ChainName+" -p tcp -m tcp -m state --state NEW --dport "+strconv.Itoa(argv.Port)+" -j LOG --log-prefix \""+LogIdentifier+"\" --log-level "+strconv.Itoa(argv.LogLevel))
		runCommand(errorHandler, "iptables -A "+ChainName+" -p tcp -m tcp --dport "+strconv.Itoa(argv.Port)+" -j "+ruleAction)
		runCommand(errorHandler, "iptables -I INPUT -j "+ChainName)
		runCommand(errorHandler, "echo \"if \\$msg contains '"+LogIdentifier+"' then "+outFile+"\" > /etc/rsyslog.d/"+ChainName+".conf")
		runCommand(errorHandler, "systemctl restart rsyslog")
		runCommand(nil, "touch "+outFile)
		fmt.Println("Created chain " + ChainName + " successfully")
		fmt.Println("All logs for port (" + strconv.Itoa(argv.Port) + ") will be in " + outFile)
		return nil
	},
}
