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
	OutputFile string `cli:"o,output" usage:"Specify the logfile" dft:"/var/log/<ChainName>"`
	DeleteRule bool   `cli:"d,delete" usage:"wether to delete the rule"`
	LogLevel   int    `cli:"l,log-level" usage:"Specify the log level" dft:"6"`
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

		ChainName := "Tripwire" + strconv.Itoa(argv.Port)
		LogIdentifier := "Tripwire" + strconv.Itoa(argv.Port)

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
			deleteRuleForChain(errorHandler, ChainName)
			runCommand(errorHandler, "iptables -X "+ChainName)
			runCommand(errorHandler, "rm /etc/rsyslog.d/"+ChainName+".conf")
			runCommand(errorHandler, "systemctl restart rsyslog.service")
			if argv.OutputFile != "/var/log/<ChainName>" {
				if _, err := os.Stat("/var/log/" + ChainName); err == nil {
					runCommand(errorHandler, "/var/log/"+ChainName)
					fmt.Println("Deleted logfile /var/log/" + ChainName)
				}
			}
			fmt.Println("Deleted chain " + ChainName + " successfully")
		} else {
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
			outFile := argv.OutputFile
			if argv.OutputFile == "/var/log/<ChainName>" {
				outFile = "/var/log/" + ChainName
			}
			runCommand(errorHandler, "iptables -N "+ChainName)
			runCommand(errorHandler, "iptables -A "+ChainName+" -p tcp -m tcp -m state --state NEW --dport "+strconv.Itoa(argv.Port)+" -j LOG --log-prefix \""+LogIdentifier+"\" --log-level "+strconv.Itoa(argv.LogLevel))
			runCommand(errorHandler, "iptables -A "+ChainName+" -p tcp -m tcp --dport "+strconv.Itoa(argv.Port)+" -j "+ruleAction)
			runCommand(errorHandler, "iptables -I INPUT -j "+ChainName)
			runCommand(errorHandler, "echo \"if \\$msg contains '"+LogIdentifier+"' then /var/log/"+outFile+"\" > /etc/rsyslog.d/"+ChainName+".conf")
			runCommand(errorHandler, "systemctl restart rsyslog")
			runCommand(errorHandler, "touch /var/log/"+outFile)
			fmt.Println("Created chain " + ChainName + " successfully")
			fmt.Println("All logs for port (" + strconv.Itoa(argv.Port) + ") will be in /var/log/" + outFile)
		}
		fmt.Println("Done")

		return nil
	})
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
		//fmt.Println(output)
	}
	if err != nil {
		if errorHandler != nil {
			errorHandler(err, sCmd)
		}
		return "", err
	}
	return output, nil
}
