# Pourquoi avoir développé un applicatif Golang pour administrer nos appliances Cisco ?

## Introduction

Sur Evita v10 l'ensemble des appliances Cisco étaient déjà supervisées via des plugin command Nagios mais tous les identifiants de connexions étaient en clair dans un fichier "caché" (le fichier était appelé dans les scripts shells donc très bien caché).

Les scripts interrogent les appliances Cisco via SSH en passant des commandes dans le CLI (CLI spécifique aux appliances).

Afin de repartir sur des bases plus saines il fallait améliorer la sécurité de ces scripts sans trop les modifier car ils fonctionnent toujours.

Ma première idée était de demander à l'équipe applicative de créer des utilisateurs ayant uniquement des droits de lecture, mais l'équipe ayant la tête dans le guidon ce n'était pas leur priorité. Puis avoir un "outil" permettant d'aller plus loin, en les rendant administrable peut être utile (Shutdow, configurer le SNMP...)

La seconde comme vous l'avez compris était de créer un applicatif intermédiaire qui se chargerait de la connexion SSH vers nos appliances de manière sécurisée.

## Choix du langage

Dans un premier temps, n'importe quel SysOps penserait à utiliser Python car c'est un des langages les plus courrant mais le problème serait juste déplacé. Vous pourriez gérer des identifiants chiffrés dans un fichier par exemple mais la clé de chiffrement serait en clair dans le code. Pas fou. De plus, Evita est un environnement fermé si j'avais besoin de libraries Python externes, j'aurai du les transférer à la main (à cette époque je ne connaissais pas l'existence de l'artifactory)

C'est pourquoi j'ai décidé de me tourner vers Golang qui est un langage compilé, la clé de chiffrement sera dans l'executable qui n'est pas lisible ; je ne serai pas embêté par la gestion des libraries externes vu que tout ce qui est nécessaire sera dans l'executable.

## Difficultés rencontrées

Comme je l'ai évoqué précédemment ces appliances Cisco ont un CLI et un terminal compliqués ce qui ne simplifie pas les les choses.

Si on s'écarte du sujet durant quelques lignes, je dois vous dire que nous n'avons jamais réussi à faire fonctionner ansible avec ces appliances que ce soit en les traitant comme une machine Linux, comme un routeur ou en modifiant les paramètres ssh d'ansible (en utilisant paramiko par exemple).

### Terminal et inputs

J'ai rencontré également beaucoup de contraintes à développer un outil pouvant les administrer. Il m'a fallu à peine 10min en me basant sur une librarie Golang pour développer quelque chose qui se connecte en SSH, envoie une ou des commande(s) et qui récupère l'output vers une machine Linux.

Un fois testé sur nos appliances, une erreur Java (oui oui) apparait car mon programme n'arrive pas à executer la commande. A ce moment là je commence à comprendre pourquoi Ansible me renvoyait soit une erreur (même si différente), soit un timeout.

Je comprends qu'il faut que je me détourne de l'idée de développer le plus proprement possible et faire simplement quelque chose de fonctionnel et fiable.

Après de multiples recherches et tests, j'arrive à trouver quelque chose qui fonctionne : passer mes commandes/inputs via **fmt.Fprintf** (_Fprintf formats according to a format specifier and writes to w. It returns the number of bytes written and any write error encountered.)_ une fois avoir établie ma connexion ssh.

### Enchainer les commandes

Second problème, enchainer les inputs et couper la session SSH une fois la ou les commandes executées. Je ne peux pas lancer toutes mes commandes d'un coup comme je le ferai vers un VM linux. Si linux m'aurait enchainé nativement toutes mes commandes avec ces appliances ce n'est pas possible, une commande envoyée avant avoir recu l'intégralité de l'output d'une commande précédente ou l'affichage de la bannière ssh est perdue.

Dans un premier temps je suis parti sur une solution un peu sale qui consitait à utiliser un timer entre l'envoie de deux commandes, ce qui fonctionnait mais le temps d'execution était très long (puis nous ne pouvons jamais être sûr du temps nécessaire à l'execution de la commande).

Pour répondre à ce problème je décide d'ajouter une goroutine qui me détectera en tâche de fond le moment où la commande a fini de s'executer ; en d'autre terme quand le terminal sera à nouveau disponible ou bien quand une confirmation sera nécessaire "Enter (yes/no)". La gestion du terminal intéractif était un autre de mes soucis :sweat_smile: 

## Fonctionnalités

Le code est assez simple, il existe 3 fonctionnalités :
1. Chiffrement, permet à l'utilisateur de chiffrer les identifiants.
1. Déchiffrement, permet à l'utilisateur de déchiffrer les identifiants.
1. Utilisation permet d'envoyer des commandes à une appliance via SSH

## Identifiants

Chaque identifiant a son propre fichier. C'est un choix volontaire pour simplifier le code et cela peut permettre à quelqu'un de l'équipe de monter en compétence sur Go en stockant tous les identifiants dans un fichier unique.

## Execution

Il y a 3 flags :
- brique => Type d'appliance à interroger
- ip => ip de la machine distante
- cmd => abbréviation de la commande à executer

## Resultat

Une seule ligne à modifier dans les scripts shell déjà existant et 1 sec de gagné sur le temps d'execution des scripts, ce qui est pas mal :smile: Même si le code n'est pas des plus élégants je vous l'accorde.

## Déploiement

Le déploiement de cet outil ce fait via une CICD très simple que je n'ai pas vraiment besoin de détailler.

La seule chose à savoir est que la clé de chiffrement n'est pas stocké dans le code quand il est poussé dans le Gitlab, un string "KEY_TO_REPLACE" est modifié avant la compilation du code dans la CICD via une simple commande shell :

```  - sed -i "s/KEY_TO_REPLACE/$ENCRYPT_KEY/" main.go ```

La variable $ENCRYPT_KEY fait référence à une variable Gitlab créé pour la CICD.

![1](uploads/ccd9a09ab41e34e504ded2c88fbcdaea/1.png)

![2](uploads/cb3397ed4bfd92714bd48cd08af55f14/2.png)