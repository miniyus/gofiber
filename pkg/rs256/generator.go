package rs256

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path"
)

func Generate(savePath string, bitSize int) {
	reader := rand.Reader

	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)

	publicKey := key.PublicKey

	saveGobKey(path.Join(savePath, "private.key"), key)
	savePEMKey(path.Join(savePath, "private.pem"), key)

	saveGobKey(path.Join(savePath, "public.key"), publicKey)
	savePublicPEMKey(path.Join(savePath, "public.pem"), publicKey)
}

func saveGobKey(fileName string, key interface{}) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			checkError(err)
		}
	}(outFile)

	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	checkError(err)
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			checkError(err)
		}
	}(outFile)

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicPEMKey(fileName string, pubKey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubKey)
	checkError(err)

	var pemKey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemFile, err := os.Create(fileName)
	checkError(err)
	defer func(pemFile *os.File) {
		err := pemFile.Close()
		if err != nil {
			checkError(err)
		}
	}(pemFile)

	err = pem.Encode(pemFile, pemKey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func GobDecode(filename string) map[string]string {
	var pubKey map[string]string

	open, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	decoder := gob.NewDecoder(open)

	err = decoder.Decode(&pubKey)
	if err != nil {
		log.Fatal(err)
	}

	return pubKey
}

func PrivatePemDecode(filename string) *rsa.PrivateKey {
	privatePem, err := os.ReadFile(filename)
	if err != nil {
		log.Println("not exists private.pem")
	}

	b, _ := pem.Decode(privatePem)
	privateKey, _ := x509.ParsePKCS1PrivateKey(b.Bytes)

	return privateKey
}

func PublicPemDecode(filename string) *rsa.PublicKey {
	publicPem, err := os.ReadFile(filename)
	publicKey := new(rsa.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	b, _ := pem.Decode(publicPem)
	_, err = asn1.Unmarshal(b.Bytes, &publicKey)

	if err != nil {
		log.Fatal(err)
	}

	return publicKey
}
