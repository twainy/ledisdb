package main

import (
	"flag"
	"fmt"
	"ledis/client"
	"regexp"
	"strconv"
	"strings"
)

var ip = flag.String("h", "127.0.0.1", "ledisdb server ip (default 127.0.0.1)")
var port = flag.Int("p", 6380, "ledisdb server port (default 6380)")
var socket = flag.String("s", "", "ledisdb server socket, overwrite ip and port")
var dbn = flag.Int("n", 0, "ledisdb database number(default 0)")

func main() {
	flag.Parse()

	cfg := new(client.Config)
	if len(*socket) > 0 {
		cfg.Addr = *socket
	} else {
		cfg.Addr = fmt.Sprintf("%s:%d", *ip, *port)
	}

	cfg.MaxIdleConns = 1

	c := client.NewClient(cfg)
	sendSelect(c, *dbn)

	SetCompletionHandler(completionHandler)
	setHistoryCapacity(100)

	reg, _ := regexp.Compile(`'.*?'|".*?"|\S+`)

	prompt := ""

	for {
		if *dbn > 0 && *dbn < 16 {
			prompt = fmt.Sprintf("%s[%d]>", cfg.Addr, *dbn)
		} else {
			prompt = fmt.Sprintf("%s>", cfg.Addr)
		}

		cmd, err := line(prompt)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}

		cmds := reg.FindAllString(cmd, -1)
		if len(cmds) == 0 {
			continue
		} else {
			addHistory(cmd)

			args := make([]interface{}, len(cmds[1:]))

			for i := range args {
				args[i] = strings.Trim(string(cmds[1+i]), "\"'")
			}

			cmd := cmds[0]
			if strings.ToLower(cmd) == "help" || cmd == "?" {
				printHelp(cmds)
			} else {
				if len(cmds) == 2 && strings.ToLower(cmds[0]) == "select" {
					if db, _ := strconv.Atoi(cmds[1]); db < 16 && db >= 0 {
						*dbn = db
					}

				}

				r, err := c.Do(cmds[0], args...)

				if err != nil {
					fmt.Printf("%s", err.Error())
				} else {
					printReply(cmd, r)
				}

				fmt.Printf("\n")
			}

		}
	}
}

func printReply(cmd string, reply interface{}) {
	switch reply := reply.(type) {
	case int64:
		fmt.Printf("(integer) %d", reply)
	case string:
		fmt.Printf("%s", reply)
	case []byte:
		fmt.Printf("%q", reply)
	case nil:
		fmt.Printf("(nil)")
	case client.Error:
		fmt.Printf("%s", string(reply))
	case []interface{}:
		for i, v := range reply {
			fmt.Printf("%d) ", i+1)
			if v == nil {
				fmt.Printf("(nil)")
			} else {
				fmt.Printf("%q", v)
			}
			if i != len(reply)-1 {
				fmt.Printf("\n")
			}
		}
	default:
		fmt.Printf("invalid ledis reply")
	}
}

func printGenericHelp() {
	msg :=
		`ledis-cli
Type:	"help <command>" for help on <command>
	`
	fmt.Println(msg)
}

func printCommandHelp(arr []string) {
	fmt.Println()
	fmt.Printf("\t%s %s \n", arr[0], arr[1])
	fmt.Printf("\tGroup: %s \n", arr[2])
	fmt.Println()
}

func printHelp(cmds []string) {
	args := cmds[1:]
	if len(args) == 0 {
		printGenericHelp()
	} else if len(args) > 1 {
		fmt.Println()
	} else {
		cmd := strings.ToUpper(args[0])
		for i := 0; i < len(helpCommands); i++ {
			if helpCommands[i][0] == cmd {
				printCommandHelp(helpCommands[i])
			}
		}
	}
}

func sendSelect(ledis_client *client.Client, index int) {
	if index > 16 || index < 0 {
		index = 0
		fmt.Println("index out of range, should less than 16")
	}
	_, err := ledis_client.Do("select", index)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}

func completionHandler(in string) []string {
	var keyWords []string
	for _, i := range helpCommands {
		if strings.HasPrefix(i[0], strings.ToUpper(in)) {
			keyWords = append(keyWords, i[0])
		}
	}
	return keyWords
}
