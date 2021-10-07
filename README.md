# Fonctionnement

## 3 Fonctionnalités 

1 - Chiffrement <br />
2 - Déchiffrement <br />
3 - Execution (envoie de commande ssh) <br />

## Chiffrement/Déchiffrement

Ces deux opérations sont gérés par le script il suffit de lancer le script avec la commande suivante : go run main.go et vous serez guidé. <br />

## Execution

Execution -brique -ip -cmd <br />
Ex : expect -brique cuic -ip A.B.C.D -cmd status // go run main.go -brique cuic -ip A.B.C.D -cmd status <br />

## Identifiants
Chaque identifiant a son propre fichier. C'est un choix volontaire pour simplifier le code et cela peut permettre à quelqu'un de l'équipe de monter en compétence sur Go en stockant tous les identifiants dans un fichier unique. <br />
