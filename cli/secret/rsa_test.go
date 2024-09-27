package secret

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	err = os.WriteFile("private.pem", privateKeyPEM, 0644)
	if err != nil {
		panic(err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	err = os.WriteFile("public.pem", publicKeyPEM, 0644)
	if err != nil {
		panic(err)
	}
}

func TestRSAEncrypt(t *testing.T) {
	ciphertext, err := RSA{}.Encrypt([]byte("hello, world!"))
	if err != nil {
		panic(err)
		return
	}

	fmt.Printf("Encrypted: %x\n", ciphertext) // "7bc039a9618692d08b2650abb31c59d40efaba71a7ee37e3169668585a6181fbe9bae803abda594eea52f5aa30a3293416901965c8368a20de09042dd75d32081c97f087579250092868b355b0ed3417807c1fc2786651794b2c70102a48a57f2f64806009295fd8fbe12d446d862ce99c8470e2de65c230273e512c2437303773fd231d0faddf2014d84295bbbe17db6f94d94723f95b76c7db1fe3f9278cb203b584813e5f5e32e79d654a7745649f3317d3e35961111d35ed6c1d396a803198bf4121c0ea7ce7b91ed3e02113184547ed9400e3a02a95c393fb5b0264275dc49ce7a5a547a544e1a88a23ff77cf04d5d13736c0506dd27c89e77bd2c3acb3"
}

func TestRSADecrypt(t *testing.T) {
	hexCipher := "7bc039a9618692d08b2650abb31c59d40efaba71a7ee37e3169668585a6181fbe9bae803abda594eea52f5aa30a3293416901965c8368a20de09042dd75d32081c97f087579250092868b355b0ed3417807c1fc2786651794b2c70102a48a57f2f64806009295fd8fbe12d446d862ce99c8470e2de65c230273e512c2437303773fd231d0faddf2014d84295bbbe17db6f94d94723f95b76c7db1fe3f9278cb203b584813e5f5e32e79d654a7745649f3317d3e35961111d35ed6c1d396a803198bf4121c0ea7ce7b91ed3e02113184547ed9400e3a02a95c393fb5b0264275dc49ce7a5a547a544e1a88a23ff77cf04d5d13736c0506dd27c89e77bd2c3acb3"
	plaintext, err := RSA{}.Decrypt(hexCipher)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decrypted: %s\n", plaintext)
}
