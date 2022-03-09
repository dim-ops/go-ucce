package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Connection struct {
	*ssh.Client
	user     string
	password string
}

func SelectCommand(CMD string) ([]string, int) {
	var commands []string
	var timeToWait int

	//Choix de la commande
	switch CMD {
	case "status":
		CMD = "show status"
	case "service":
		CMD = "utils service list"
	case "replication":
		CMD = "utils dbreplication runtimestate"
	case "jtapi_users":
		CMD = "run sql select name from ApplicationUser where name like '%jtapi%'"
	case "licence_cuic":
		CMD = "show cuic license-info"
	case "shutdown":
		CMD = "utils system shutdown"
	}

	if CMD == "utils system shutdown" {
		//les commandes à passer
		commands = []string{
			CMD,
			//"yes",
		}

		//Besoin d'attendre plus pour éteindre le VOS
		timeToWait = 25

	} else {
		//les commandes à passer
		commands = []string{
			CMD,
			"exit",
		}

		// L'output est un peu longue à charger, il faut attendre 20sec pour être sûr de la capturer
		timeToWait = 20
	}

	return commands, timeToWait
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

	return &Connection{conn, Identifiants[0], Identifiants[1]}, nil
}

func (conn *Connection) SendCommands(CMD string) error {

	cmds, timeToWait := SelectCommand(CMD)
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

	if cmds[0] != "utils system shutdown" {
		//Output ssh
		sess.Stdout = os.Stdout
	} else {
		stdout, err := sess.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		var output []byte

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

				if strings.Contains(line, "Enter (yes/no)") {
					_, err = fmt.Fprintf(stdin, "%s\n", "yes")
					if err != nil {
						break
					}
				}
			}
		}(stdin, stdout, &output)
	}

	//Lanchement du shell à distance
	err = sess.Shell()
	if err != nil {
		log.Fatal(err)
	}

	//Envoie des inputs => CMD
	for _, cmd := range cmds {

		_, err = fmt.Fprintf(stdin, "%s\n", cmd)

		//fmt.Println(cmd)
		//Besoin de ralentir l'execution du script, sinon il lance le exit avant le VOS n'ait eu le temps d'executer la premiere CMD
		time.Sleep(time.Duration(timeToWait) * time.Second)

		if err != nil {
			log.Fatal(err)
		}
	}

	//Attend que la ou les commandes ssh n'execute
	err = sess.Wait()
	if err != nil {
		log.Fatal(err)
	}

	return err
}
