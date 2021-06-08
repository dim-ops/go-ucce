package main

import (
	"DERS_3CG_CISCO/expect/ssh"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

var Identifiants []string

var Brique string
var IP string
var CMD string

func main() {

	flag.StringVar(&Brique, "brique", "", "quel type de VOS")
	flag.StringVar(&IP, "ip", "", "quelle IP")
	flag.StringVar(&CMD, "cmd", "", "quelle commande executer")
	flag.Parse()
	flag.Args()

	file, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.Printf("Read file with identifiants")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		Identifiants = append(Identifiants, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if Brique != "" && IP != "" && CMD != "" {
		ssh.ConnexionSSH(Identifiants, IP, CMD)
	} else {
		fmt.Printf("Argument manquant pour initialiser une connexion ssh\n")
	}

	//encrypt.FirstStep()

}
