package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"os"
	"path/filepath"

	"golang.org/x/term"
)

func PasswordInput(mess string) []byte {
	fmt.Print(mess + " ")
	bytepw, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Print("\n")
	check(err)
	return bytepw
}

func StoreNewPassword() {
	firstInput := PasswordInput("Please enter the new password:")
	secondInput := PasswordInput("Please confirm the password again:")
	if len(firstInput) != 0 && firstInput != nil && bytes.Compare(firstInput, secondInput) == 0 {
		path, err := os.Getwd()
		check(err)
		passDir := filepath.Join(path, "resources", "pass")
		var f *os.File
		f, err = os.OpenFile(passDir, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		check(err)
		defer func(f *os.File) {
			err := f.Close()
			check(err)
		}(f)
		hash := sha256.Sum256(secondInput)
		fmt.Printf("%x", hash)
		fmt.Println()
		_, err = f.WriteString(hex.EncodeToString(hash[:]) + "\n")
		check(err)
		return
	}
	fmt.Println("The passwords did not match or was nil.")
}

func Login() *route53.Client {
	password := PasswordInput("Login password:")
	hash := sha256.Sum256(password)
	passwordString := hex.EncodeToString(hash[:])
	path, err := os.Getwd()
	check(err)
	passDir := filepath.Join(path, "resources", "pass")
	var f *os.File
	f, err = os.Open(passDir)
	check(err)
	defer func(f *os.File) {
		err := f.Close()
		check(err)
	}(f)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if passwordString == scanner.Text() {
			fmt.Println("Login success.")
			cfg, err := config.LoadDefaultConfig(context.Background(),
				config.WithRegion("aws-global"),
			)
			check(err)
			return route53.NewFromConfig(cfg)
		}
	}
	err = scanner.Err()
	check(err)
	fmt.Println("Login fail.")
	return nil
}
