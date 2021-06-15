package main

import (
	"DERS_3CG_CISCO/expect/ssh"
	"flag"
	"log"
	"os"
)

//Declaration des variables
var Identifiants []string
var Brique string
var IP string
var CMD string

func main() {
	//les flags correspondent aux variables à passer au strict
	//go run .\main.go -brique cuic -ip 10.218.115.26 -cmd status
	flag.StringVar(&Brique, "brique", "", "quel type de VOS")
	flag.StringVar(&IP, "ip", "", "quelle IP")
	flag.StringVar(&CMD, "cmd", "", "quelle commande executer")
	flag.Parse()
	flag.Args()

	//Ouverture du fichier avec les identifiants
	file, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	//Go à la particularité de détecter la fin de l'utilisation d'un fichier, pour cela GO utilise defer pour fermer le fichier une fois son utilisation terminées.
	defer file.Close()

	log.Printf("Read file with identifiants")

	// //Lecture du fichier contenant les informations
	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	//Ajouts des identifiants dans le tableau []Identifiants
	// 	Identifiants = append(Identifiants, scanner.Text())
	// }

	// //Si n'arrive pas à lire le fichier erreur
	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	//Si l'un des flags est manquant erreur sinon appelle de la function ConnexionSSH présent dans le package ssh
	if Brique != "" && IP != "" && CMD != "" {
		ssh.ConnexionSSH(Identifiants, IP, CMD)
	} else {
		log.Printf("Argument manquant pour initialiser une connexion ssh\n")
	}

}

// func main() {
// 	chiffrement.WhatToDo()
// }
