package secret

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
}

func (RSA) Encrypt(data []byte) ([]byte, error) {
	publicKeyPEM, err := os.ReadFile("public.pem")
	if err != nil {
		return nil, err
	}
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
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

func (RSA) Decrypt(hexCipher string) ([]byte, error) {
	privateKeyPEM, err := os.ReadFile("private.pem")
	if err != nil {
		return nil, err
	}
	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
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
