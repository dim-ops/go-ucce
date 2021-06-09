package ssh

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	timeout = 10 * time.Minute
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

	e, _, err := expect.SpawnSSH(config, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()

	e.ExpectBatch([]expect.Batcher{
		&expect.BCas{[]expect.Caser{
			&expect.Case{R: regexp.MustCompile(`router#`), T: expect.OK()},
			&expect.Case{R: regexp.MustCompile(`Login: `), S: *user,
				T: expect.Continue(expect.NewStatus(codes.PermissionDenied, "wrong username")), Rt: 3},
			&expect.Case{R: regexp.MustCompile(`Password: `), S: *pass1, T: expect.Next(), Rt: 1},
				T: expect.Continue(expect.NewStatus(codes.PermissionDenied, "wrong password")), Rt: 1},
		}},
	}, timeout)

}
