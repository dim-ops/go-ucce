package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func ConnexionSSH(Identifiants []string, IP, CMD string) {

	var commands []string

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
			"yes",
		}
	} else {
		//les commandes à passer
		commands = []string{
			CMD,
			"exit",
		}
	}

	//log.Printf("Connexion ssh is pending...")

	//Port SSH à utiliser
	port := "22"

	//Configuration du client SSH
	config := &ssh.ClientConfig{
		User: Identifiants[0],
		Auth: []ssh.AuthMethod{
			ssh.Password(Identifiants[1]),
		},
		//Ignore les clés SSH
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		//Si pas de réponse au bout de 10sec => Stop
		Timeout: 10 * time.Second,
	}

	//Connection à la VM
	client, err := ssh.Dial("tcp", IP+":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	//Lancement de la connexion SSH
	sess, err := client.NewSession()
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

	//Output ssh
	sess.Stdout = os.Stdout

	//Error ssh
	sess.Stderr = os.Stderr

	//Permet l'envoie des inputs
	stdin, err := sess.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	//Lanchement du shell à distance
	err = sess.Shell()
	if err != nil {
		log.Fatal(err)
	}

	//Envoie des inputs => CMD
	for _, cmd := range commands {

		_, err = fmt.Fprintf(stdin, "%s\n", cmd)

		//fmt.Println(cmd)
		//Besoin de ralentir l'execution du script, sinon il lance le exit avant le VOS n'ait eu le temps d'executer la premiere CMD
		time.Sleep(10 * time.Second)

		if err != nil {
			log.Fatal(err)
		}
	}

	//Attend que la ou les commandes ssh n'execute
	err = sess.Wait()
	if err != nil {
		log.Fatal(err)
	}

}
