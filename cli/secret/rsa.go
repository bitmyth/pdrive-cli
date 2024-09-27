package secret

// https://dev.to/elioenaiferrari/asymmetric-cryptography-with-golang-2ffd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"log"
	"os"
)

type RSA struct {
	publicKeyPEM  string
	privateKeyPEM string
}

func NewRSA(privateKeyPEMFile string) *RSA {
	rsa := &RSA{}
	if privateKeyPEMFile == "" {
		return rsa
	}

	privateKey, err := os.ReadFile(privateKeyPEMFile)
	if err != nil {
		println(privateKeyPEMFile, " keyfile not exists")
		return rsa
	}

	rsa.privateKeyPEM = string(privateKey)
	println("load keyfile ", privateKeyPEMFile)

	publicKeyPEM, err := rsa.ExtractPublicKey(privateKey)
	if err == nil {
		rsa.publicKeyPEM = string(publicKeyPEM)
	}

	return rsa
}

func (r RSA) Encrypt(data []byte) ([]byte, error) {
	publicKeyPEM := r.publicKeyPEM
	publicKeyBlock, _ := pem.Decode([]byte(publicKeyPEM))
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	plaintext := data
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), plaintext)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func (r RSA) Decrypt(hexCipher string) ([]byte, error) {
	privateKeyPEM := r.privateKeyPEM
	privateKeyBlock, _ := pem.Decode([]byte(privateKeyPEM))
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// Parse the hex string to bytes
	bytes, err := hex.DecodeString(hexCipher)
	if err != nil {
		log.Fatalf("Error decoding hex string: %v", err)
		return nil, err
	}
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, bytes)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (r RSA) ExtractPublicKey(privateKeyPEM []byte) ([]byte, error) {
	// Decode the PEM file
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		log.Fatal("Failed to decode PEM block containing the private key")
	}

	// Parse the private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Error parsing private key: %v", err)
	}

	// Extract the public key
	publicKey := privateKey.Public()
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return publicKeyPEM, nil
}
