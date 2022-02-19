package main
import (
    "os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    d1 := []byte("hello\ngo\n")
    err1 := os.WriteFile("/d/Landmark-API/golang/resources/pass", d1, 0644)
    check(err1)
}
