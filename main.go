package main

import (
	"flag"
	"log"
	"os"
	"runtime"
)

//Declaration des variables
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

		WhatToDo()

	} else if Brique != "" && IP != "" && CMD != "" {
		//Si l'un des flags est manquant erreur sinon appelle de la function ConnexionSSH présent dans le package ssh
		//log.Printf("Read file with identifiants")

		//Extension des fichiers
		extFile := [2]string{"Id.txt", "Pass.txt"}

		if runtime.GOOS == "windows" {
			//environnement de dev (e-buro)
			for i := range extFile {
				chiffer, _ := ReadFromFile(Brique + extFile[i])
				id := Decrypt(string(chiffer), "KEY_TO_REPLACE")
				Identifiants = append(Identifiants, id)
			}
		} else {
			for i := range extFile {
				//environnement d'exploitation Linux
				pathFile := "/usr/etc/script/" + Brique + extFile[i]
				chiffer, _ := ReadFromFile(pathFile)
				id := Decrypt(string(chiffer), "20Ders3CGEvita20")
				Identifiants = append(Identifiants, id)
			}
		}

		//Lecture des fichiers contenant les IDs
		ConnexionSSH(Identifiants, IP, CMD)

	} else {
		os.Exit(1)
		log.Printf("Argument manquant pour initialiser une connexion ssh\n")
	}

}
