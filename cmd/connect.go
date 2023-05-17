package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/cdgeass/ssh-client/config"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.Flags().IntVar(&index, "index", -1, "Server index")
	connectCmd.Flags().StringVar(&name, "name", "", "Server name")
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a server",
	Run: func(cmd *cobra.Command, args []string) {
		if index == -1 && name == "" {
			log.Fatalln("Please select a server by enter index or name")
		}

		var server config.Server
		for i, s := range conf.Servers {
			if i == index || s.Name == name {
				server = s
				break
			}
		}

		connect(server)
	},
}

func connect(server config.Server) {
	log.Println("Connecting to " + server.Name)

	//var hostKey ssh.PublicKey
	conf := &ssh.ClientConfig{
		User: server.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port), conf)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
		return
	}
	defer session.Close()

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
		ssh.VSTATUS:       1,
	}

	fd := int(os.Stdin.Fd())

	sysType := runtime.GOOS
	if sysType == "windows" {
		// Set windows.ENABLE_VIRTUAL_TERMINAL_INPUT
		setConsoleMode(fd)
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer term.Restore(fd, oldState)

	// Request pseudo terminal
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	if err := session.RequestPty("xterm", height, width, modes); err != nil {
		log.Fatal("Request for pseudo terminal failed: ", err)
		return
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatal("Failed to start shell: ", err)
	}

	if err := session.Wait(); err != nil {
		log.Fatal("Failed to wait: ", err)
	}
}

func setConsoleMode(fd int) {
	var st uint32
	if err := windows.GetConsoleMode(windows.Handle(fd), &st); err != nil {
		log.Fatal("Failed to set console mode: ", err)
		return
	}
	raw := st | windows.ENABLE_VIRTUAL_TERMINAL_INPUT
	if err := windows.SetConsoleMode(windows.Handle(fd), raw); err != nil {
		log.Fatal("Failed to set console mode: ", err)
		return
	}
}
