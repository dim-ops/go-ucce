package main

import "github.com/DimProject/ucce-cisco/cmd"

//Declaration des variables
// var Identifiants []string
// var Brique string
// var IP string
// var CMD string

// func main() {
// 	flag.StringVar(&Brique, "brique", "", "quel type de VOS")
// 	flag.StringVar(&IP, "ip", "", "quelle IP")
// 	flag.StringVar(&CMD, "cmd", "", "quelle commande executer")
// 	flag.Parse()
// 	flag.Args()

// 	if Brique != "" && IP != "" && CMD != "" {

// 		IP = IP + ":22"
// 		conn, err := Connect(Identifiants, IP)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		err = conn.SendCommands(CMD)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 	} else {
// 		os.Exit(1)
// 		log.Printf("Argument manquant pour initialiser une connexion ssh\n")
// 	}

// }

func main() {
	cmd.Execute()
}
