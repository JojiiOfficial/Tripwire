package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Accept     bool   `cli:"a,accept" usage:"Specify wether to drop or accept the incoming connections"`
	Port       int    `cli:"*p,port" usage:"Specify the port to apply the wire to"`
	Output     string `cli:"o,output" usage:"Specify log file path"`
	DeleteRule bool   `cli:"d,delete" usage:"wether to delete the rule"`
	LogLevel   int    `cli:"l,log-level" usage:"Specify the log level" deft:"6"`
}

func main() {
	cli.RunWithArgs(new(argT), os.Args, func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		_ = argv
		if os.Getuid() != 0 {
			fmt.Println("You need to be root!")
			os.Exit(1)
			return nil
		}

		ChainName := "Tripwire\\[" + strconv.Itoa(argv.Port) + "\\]"
		LogIdentifier := ChainName

		errorHandler := func(err error, cmd string) {
			fmt.Println("Error running " + cmd + ": " + err.Error())
			os.Exit(2)
		}

		if argv.DeleteRule {
			if chainExisits(ChainName) != nil {
				fmt.Println("Chain doesn't exist")
				return nil
			}
			runCommand(errorHandler, "iptables -F "+ChainName)
			runCommand(errorHandler, "iptables -X "+ChainName)
			runCommand(errorHandler, "rm /etc/rsyslog.d/"+ChainName)
			runCommand(errorHandler, "systemctl restart rsyslog.service")
			fmt.Println("Deleted chain " + ChainName + " successfully")
		} else {
			ruleAction := "ACCEPT"
			if !argv.Accept {
				ruleAction = "DROP"
			}
			if len(argv.Output) == 0 {
				fmt.Println("You need to set the outputfile. Use the -o or --output argument")
				return nil
			}
			if argv.Port <= 0 || argv.Port > 65535 {
				fmt.Println("You port must be between 0 and 65535!")
				return nil
			}
			if chainExisits(ChainName) == nil {
				fmt.Println("This port already has a rule! Try deleting it with -d")
				return nil
			}
			confFile := argv.Output
			if !strings.HasPrefix(confFile, ".conf") {
				confFile = confFile + ".conf"
			}
			runCommand(errorHandler, "iptables -N "+ChainName)
			runCommand(errorHandler, "iptables -A "+ChainName+" -j LOG --log-prefix "+LogIdentifier+" --log-level "+strconv.Itoa(argv.LogLevel))
			runCommand(errorHandler, "iptables -A "+ChainName+" -j "+ruleAction)
			runCommand(errorHandler, "echo \":msg,contains,"+LogIdentifier+" /var/log/"+confFile+"\" > /etc/rsyslog.d/"+ChainName)
			runCommand(errorHandler, "systemctl restart rsyslog")
			fmt.Println("Created chain " + ChainName + " successfully")
		}
		fmt.Println("Done")

		return nil
	})
}

func chainExisits(chainName string) error {
	res, err := runCommand(nil, "iptables -L "+chainName)
	if err != nil {
		return err
	}
	if strings.HasSuffix(res, "iptables: No chain") {
		return nil
	}
	return nil
}

func runCommand(errorHandler func(error, string), sCmd string) (outb string, err error) {
	out, err := exec.Command("su", "-c", sCmd).Output()
	output := string(out)
	if len(strings.ReplaceAll(strings.Trim(output, " "), "\n", "")) > 0 {
		fmt.Println(output)
	}
	if err != nil {
		if errorHandler != nil {
			errorHandler(err, sCmd)
		}
		return "", err
	}
	return output, nil
}
