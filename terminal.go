package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func Terminal(server Server) {
	log.Println("Connecting to " + server.Name)

	//var hostKey ssh.PublicKey
	config := &ssh.ClientConfig{
		User: server.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port), config)
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
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
		return
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}

	if err := session.Wait(); err != nil {
		log.Fatal("failed to wait: ", err)
	}
}
