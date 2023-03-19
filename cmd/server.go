package cmd

import (
	"fmt"
	"github.com/cdgeass/ssh-client/config"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.AddCommand(serverListCmd)

	serverCmd.AddCommand(serverAddCmd)
	serverAddCmd.Flags().StringVar(&name, "name", "", "Server name")
	serverAddCmd.Flags().StringVar(&host, "host", "", "Server host")
	serverAddCmd.Flags().IntVar(&port, "port", 22, "Server port")
	serverAddCmd.Flags().StringVar(&user, "user", "", "Server user")
	serverAddCmd.Flags().StringVar(&password, "password", "", "Server password")

	serverCmd.AddCommand(serverEditCmd)
	serverEditCmd.Flags().IntVar(&index, "index", -1, "Server name")
	serverEditCmd.Flags().StringVar(&name, "name", "", "Server name")
	serverEditCmd.Flags().StringVar(&host, "host", "", "Server host")
	serverEditCmd.Flags().IntVar(&port, "port", 22, "Server port")
	serverEditCmd.Flags().StringVar(&user, "user", "", "Server user")
	serverEditCmd.Flags().StringVar(&password, "password", "", "Server password")

	serverCmd.AddCommand(serverDeleteCmd)
	serverDeleteCmd.Flags().IntVar(&index, "index", -1, "Server name")
	serverDeleteCmd.Flags().StringVar(&name, "name", "", "Server name")
	serverDeleteCmd.Flags().StringVar(&host, "host", "", "Server host")
	serverDeleteCmd.Flags().IntVar(&port, "port", 22, "Server port")
	serverDeleteCmd.Flags().StringVar(&user, "user", "", "Server user")
	serverDeleteCmd.Flags().StringVar(&password, "password", "", "Server password")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server manager",
}

var serverListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Servers: ")
		if conf.Servers != nil {
			for i, server := range conf.Servers {
				fmt.Println(fmt.Sprintf("> %d: %s", i, server.Name))
			}
		}
	},
}

var serverAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		for _, server := range conf.Servers {
			if server.Name == name {
				log.Fatalln(fmt.Sprintf("Server %s existed", server.Name))
				return
			}
		}

		conf.Servers = append(conf.Servers, config.Server{
			Name:     name,
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
		})
		if err := conf.Save(); err != nil {
			log.Fatalln("Failed to save server: ", err)
		}
	},
}

var serverEditCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"e"},
	Run: func(cmd *cobra.Command, args []string) {
		if index == -1 && name == "" {
			log.Fatalln("Please select a server by enter index or name")
		}

		var server *config.Server
		for i, s := range conf.Servers {
			if i == index || s.Name == name {
				index = i
				server = &s
				break
			}
		}
		if server == nil {
			log.Fatalln("Please enter servername")
		}

		if name != "" {
			server.Name = name
		}
		if host != "" {
			server.Host = host
		}
		if port != server.Port {
			server.Port = port
		}
		if user != "" {
			server.User = user
		}
		if password != "" {
			server.Password = password
		}

		conf.Servers[index] = *server
		if err := conf.Save(); err != nil {
			log.Fatalln("Failed to edit server info: ", err)
		}
	},
}

var serverDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d"},
	Run: func(cmd *cobra.Command, args []string) {
		if index == -1 && name == "" {
			log.Fatalln("Please select a server by enter index or name")
		}

		for i, s := range conf.Servers {
			if i == index || s.Name == name {
				conf.Servers = append(conf.Servers[:i], conf.Servers[i+1:]...)
				if err := conf.Save(); err != nil {
					log.Fatalln("Failed to delete server: ", err)
				}
				break
			}
		}
	},
}
