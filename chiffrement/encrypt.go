package chiffrement

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

//Doc : https://gist.github.com/josephspurrier/12cc5ed76d2228a41ceb
//https://eli.thegreenplace.net/2019/aes-encryption-of-files-in-go/

func Decrypt(cipherstring string, keystring string) string {
	// Byte array of the string
	ciphertext := []byte(cipherstring)

	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		panic("Text is too short")
	}

	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]

	// Remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext)
}

func encrypt(plainstring, keystring string) string {
	// Byte array of the string
	plaintext := []byte(plainstring)

	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Empty array of 16 + plaintext length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return string(ciphertext)
}

func readline() string {
	bio := bufio.NewReader(os.Stdin)
	line, _, err := bio.ReadLine()
	if err != nil {
		fmt.Println(err)
	}
	return string(line)
}

func writeToFile(data, file string) {
	ioutil.WriteFile(file, []byte(data), 777)
}

func ReadFromFile(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	return data, err
}

func WhatToDo() {
	key := "20Ders3CGEvita20"

	for {
		fmt.Print("Qu'est ce que tu veux faire ? (encrypt/descrypt/exit)\n")
		line := readline()

		switch line {
		case "help":
			fmt.Println("You can:\nencrypt\ndecrypt\nexit")
		case "exit":
			os.Exit(0)
		case "encrypt":
			fmt.Print("Quel identifiant veux-tu chiffrer ?\n")
			line2 := readline()
			cipherId := encrypt(line2, key)
			fmt.Print("Quel mot de passe veux-tu ciffrer ?\n")
			line3 := readline()
			cipherPass := encrypt(line3, key)
			fmt.Print("Pour quelle brique sont ces identifiants ?\n")
			line4 := readline()
			writeToFile(cipherId, line4+"Id.txt")
			writeToFile(cipherPass, line4+"Pass.txt")
			fmt.Println("Wrote to files")
		case "decrypt":
			fmt.Print("What is the name of the file to decrypt: ")
			line2 := readline()
			if ciphertext, err := ReadFromFile(line2); err != nil {
				fmt.Println("File is not found")
			} else {
				plaintext := Decrypt(string(ciphertext), key)
				fmt.Println(plaintext)
			}
		}
	}
}
