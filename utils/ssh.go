package ssh

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

type Connection struct {
	*ssh.Client
	user     string
	password string
}

// func SelectCommand(CMD, brique string) string {
// 	var command string
// 	//Choix de la commande
// 	switch CMD {
// 	case "status":
// 		command = "show status"
// 	case "service":
// 		command = "utils service list"
// 	case "replication":
// 		command = "utils dbreplication runtimestate"
// 	case "jtapi_users":
// 		command = "run sql select name from ApplicationUser where name like '%jtapi%'"
// 	case "licence_cuic":
// 		command = "show cuic license-info"
// 	case "shutdown":
// 		if brique != "cusp" {
// 			CMD = "utils system shutdown"
// 			//commands = []string{CMD, "\n"}
// 		} else {
// 			CMD = "shutdown"
// 			//commands = []string{CMD}
// 		}

// 	}

// 	return command
// }

func Connect(user, password, host string) (*Connection, error) {

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, err
		//return nil, fmt.Errorf("Failed to dial: %s", err)
	}

	return &Connection{conn, user, password}, nil
}

func (conn *Connection) newSession() (*ssh.Session, error) {
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 0, 40, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	return session, nil
}

func (conn *Connection) SendCommands(cmd string) error {

	//Lancement de la connexion SSH
	sess, err := conn.newSession()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

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

	//Lanchement du shell à distance
	err = sess.Shell()
	if err != nil {
		log.Fatal(err)
	}

	//Envoie des inputs => CMD
	_, err = fmt.Fprintf(stdin, "%s\n", cmd)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	var output []byte
	//Tourne en tâche de fond
	//Analyse chaque ligne renvoyée dans le terminal
	//Si contient (yes/no), la fonction le détécte et envoie yes pour arrêter le VOS
	go func(stdin io.WriteCloser, stdout io.Reader, output *[]byte) {
		defer wg.Done()
		var (
			line string
			r    = bufio.NewReader(stdout)
		)
		for {
			wg.Add(1)
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
			outputString := string(*output)
			if strings.HasPrefix(line, "admin") && strings.HasSuffix(line, ":") && strings.Contains(outputString, cmd) {

				_, err = fmt.Fprintf(stdin, "%s\n", "exit")
				if err != nil {
					log.Fatal(err)
				}
				break
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

	wg.Wait()

	//Attend que la ou les commandes ssh n'execute
	err = sess.Wait()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
