package ssh

import (
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func ConnexionSSH(Identifiants []string, IP, CMD string) {

	log.Printf("Connexion ssh is pending...")
	port := "22"

	// ssh client config
	config := &ssh.ClientConfig{
		User: Identifiants[0],
		Auth: []ssh.AuthMethod{
			ssh.Password(Identifiants[1]),
		},
		// allow any host key to be used (non-prod)
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		// optional tcp connect timeout
		Timeout: 10 * time.Second,
	}

	// connect
	client, err := ssh.Dial("tcp", IP+":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// start session
	sess, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	// setup standard out and error
	// uses writer interface
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	// run single command
	err = sess.Run(CMD)
	if err != nil {
		log.Fatal(err)
	}
}
