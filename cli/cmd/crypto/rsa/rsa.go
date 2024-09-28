package rsa

import (
	"encoding/hex"
	"fmt"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"github.com/bitmyth/pdrive-cli/cli/secret"
)

func Encrypt(f *factory.Factory, plainText string) error {
	//cfg, _ := f.Config()
	rsa := secret.NewRSA(f.KeyFile)
	cipher, err := rsa.Encrypt([]byte(plainText))
	if err != nil {
		return err
	}

	content := hex.EncodeToString(cipher)

	cs := f.IOStreams.ColorScheme()
	infoColor := cs.Cyan

	fmt.Fprintln(f.IOStreams.Out, infoColor("RSA encrypted"))
	fmt.Fprintln(f.IOStreams.Out, cs.Green(content))
	return nil
}

func Decrypt(f *factory.Factory, cipher string) error {
	r := secret.NewRSA(f.KeyFile)

	plaintext, err := r.Decrypt(cipher)
	if err != nil {
		return err
	}

	cs := f.IOStreams.ColorScheme()
	infoColor := cs.Cyan

	fmt.Fprintln(f.IOStreams.Out, infoColor("RSA decrypted"))
	fmt.Fprintln(f.IOStreams.Out, cs.Green(string(plaintext)))
	return nil
}
