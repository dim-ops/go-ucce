package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Connection struct {
	*ssh.Client
	user     string
	password string
}

func SelectCommand(CMD, brique string) string {
	var command string
	//Choix de la commande
	switch CMD {
	case "status":
		command = "show status"
	case "service":
		command = "utils service list"
	case "replication":
		command = "utils dbreplication runtimestate"
	case "jtapi_users":
		command = "run sql select name from ApplicationUser where name like '%jtapi%'"
	case "licence_cuic":
		command = "show cuic license-info"
	case "shutdown":
		if brique != "cusp" {
			CMD = "utils system shutdown"
			//commands = []string{CMD, "\n"}
		} else {
			CMD = "shutdown"
			//commands = []string{CMD}
		}

	}

	log.Println("Commande sélectionnée")

	return command
}

func Connect(Identifiants []string, IP string) (*Connection, error) {

	sshConfig := &ssh.ClientConfig{
		User: Identifiants[0],
		Auth: []ssh.AuthMethod{
			ssh.Password(Identifiants[1]),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", IP, sshConfig)
	if err != nil {
		return nil, err
	}

	log.Println("Fonctione Connect appelé")

	return &Connection{conn, Identifiants[0], Identifiants[1]}, nil
}

func (conn *Connection) SendCommands(CMD string) error {

	command := SelectCommand(CMD, Brique)
	//Lancement de la connexion SSH
	sess, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	//Configuration par défaut/basique d'un Terminal
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	//Lancement d'un PTY (Pseudo terminal), chose obligatoire pour fonctionner avec les VOS
	if err := sess.RequestPty("xterm", 0, 40, modes); err != nil {
		sess.Close()
	}

	//Error ssh
	sess.Stderr = os.Stderr

	//Permet l'envoie des inputs
	stdin, err := sess.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	//Output ssh
	//Nécessaire pour les pluginCommands
	stdout, err := sess.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	var output []byte
	//Tourne en tâche de fond
	//Analyse chaque ligne renvoyée dans le terminal
	//Si contient (yes/no), la fonction le détécte et envoie yes pour arrêter le VOS
	go func(stdin io.WriteCloser, stdout io.Reader, output *[]byte) {
		var (
			line string
			r    = bufio.NewReader(stdout)
		)
		for {
			b, err := r.ReadByte()
			if err != nil {
				break
			}

			*output = append(*output, b)

			if b == byte('\n') {
				line = ""
				continue
			}

			line += string(b)

			if strings.Contains(line, "admin") {
				_, err = fmt.Fprintf(stdin, "%s\n", "exit")
				if err != nil {
					log.Fatal(err)
				}
			}

			if strings.Contains(line, "Enter (yes/no)") {
				_, err = fmt.Fprintf(stdin, "%s\n", "yes")
				if err != nil {
					log.Fatal(err)
				}
			}

		}

		fmt.Print(string(*output))
	}(stdin, stdout, &output)

	//Lanchement du shell à distance
	err = sess.Shell()
	if err != nil {
		log.Fatal(err)
	}

	//Envoie des inputs => CMD
	_, err = fmt.Fprintf(stdin, "%s\n", command)
	if err != nil {
		log.Fatal(err)
	}

	//Attend que la ou les commandes ssh n'execute
	err = sess.Wait()
	if err != nil {
		log.Fatal(err)
	}

	return err
}
