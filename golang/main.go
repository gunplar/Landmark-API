package main
import (

)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    StoreNewPassword()
}
