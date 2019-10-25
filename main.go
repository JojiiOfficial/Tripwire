package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mkideal/cli"
)

func main() {
	if err := cli.Root(root,
		cli.Tree(help),
		cli.Tree(addCMD),
		cli.Tree(deleteCMD),
		cli.Tree(list),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type argT struct {
	cli.Helper
}

var help = cli.HelpCommand("display help information")
var errorHandler = func(err error, cmd string) {
	fmt.Println("Error running " + cmd + ": " + err.Error())
	os.Exit(2)
}

var root = &cli.Command{
	Desc: "this is root command",
	Argv: func() interface{} { return new(argT) },
	Fn: func(ctx *cli.Context) error {
		fmt.Println("Usage: tripwire <help/add/delete/list> [-h,-p,-o,-l,-a]")
		return nil
	},
}

func checkRoot() bool {
	if os.Getuid() != 0 {
		fmt.Println("You need to be root to run this command!")
		return false
	}
	return true
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
	if err != nil {
		if errorHandler != nil {
			errorHandler(err, sCmd)
		}
		return "", err
	}
	return output, nil
}

func generateChainname(port int) string {
	return "Tripwire" + strconv.Itoa(port)
}
