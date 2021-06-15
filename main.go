package main

import (
	"DERS_3CG_CISCO/expect/chiffrement"
	"DERS_3CG_CISCO/expect/ssh"
	"flag"
	"log"
)

//Declaration des variables
var ExtFile []string
var Identifiants []string
var Brique string
var IP string
var CMD string
var ENCRYPT bool

func main() {
	//les flags correspondent aux variables à passer au strict
	//go run .\main.go -brique cuic -ip 10.218.115.26 -cmd status
	////La valeur des StringVar est "" par défaut
	flag.StringVar(&Brique, "brique", "", "quel type de VOS")
	flag.StringVar(&IP, "ip", "", "quelle IP")
	flag.StringVar(&CMD, "cmd", "", "quelle commande executer")
	//La valeur de ENCRYPT est false par défaut
	flag.BoolVar(&ENCRYPT, "encrypt", false, "chiffrement des identifiants")
	flag.Parse()
	flag.Args()

	//Si ENCRYPT == true, l'utilisateur va chiffrer les identifiants
	if ENCRYPT == true {

		chiffrement.WhatToDo()

	} else if Brique != "" && IP != "" && CMD != "" {
		//Si l'un des flags est manquant erreur sinon appelle de la function ConnexionSSH présent dans le package ssh
		log.Printf("Read file with identifiants")

		//Extension des fichiers
		extFile := [2]string{"Id.txt", "Pass.txt"}

		//Lecture des fichiers contenant les IDs
		for i := range extFile {
			chiffer, _ := chiffrement.ReadFromFile(Brique + extFile[i])
			id := chiffrement.Decrypt(string(chiffer), "testtesttesttest")
			Identifiants = append(Identifiants, id)
		}

		ssh.ConnexionSSH(Identifiants, IP, CMD)

	} else {
		log.Printf("Argument manquant pour initialiser une connexion ssh\n")
	}

}
