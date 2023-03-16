package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	config := LoadConfig()

	app := &cli.App{
		Name:        "sshg",
		Version:     "1.0.0",
		Description: "This is a cli ssh client written by go",
		Commands: []*cli.Command{
			{
				Name:    "connect",
				Aliases: []string{"c"},
				Usage:   "connect to a server",
				Action: func(context *cli.Context) error {
					var arg string

					args := context.Args()
					if !args.Present() {
						fmt.Println("Servers:")
						for i, server := range config.Servers {
							fmt.Printf("  %d: %s\n", i, server.Name)
						}

						scanner := bufio.NewScanner(os.Stdin)
						for scanner.Scan() {
							arg = scanner.Text()
							break
						}
					} else {
						arg = args.First()
					}

					for i, server := range config.Servers {
						if fmt.Sprintf("%d", i) == arg {
							Terminal(server)
						} else if server.Name == arg {
							Terminal(server)
						}
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
