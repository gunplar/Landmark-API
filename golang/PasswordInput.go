package main
import (
    "fmt"
    "golang.org/x/term"
    "os"
    "path/filepath"
    "bytes"
    "crypto/sha256"
    "encoding/hex"
)

func PasswordInput(mess string) []byte {
    fmt.Print(mess+" ")
    bytepw, err := term.ReadPassword(int(os.Stdin.Fd()))
    fmt.Print("\n")
    check(err)
    return bytepw
}

func StoreNewPassword() {
    firstInput := PasswordInput("Please enter the new password")
    secondInput := PasswordInput("Please confirm the password again")
    if bytes.Compare(firstInput,secondInput) == 0 {
        path, err := os.Getwd()
        check(err)
        passDir := filepath.Join(path,"resources","pass")
        var f *os.File 
        f, err = os.OpenFile(passDir, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
        check(err)
        hash := sha256.Sum256(secondInput)
        fmt.Printf("%x", hash)
        _, err = f.WriteString(hex.EncodeToString(hash[:])+"\n")
        check(err)
        defer f.Close()
    }
}