package common

import (
    "fmt"
)

const debug = false

func INFO(str string) {
    if debug {
        fmt.Printf(str)
    }
}

func ERROR(str string) {
    if debug {
        fmt.Printf(str)
    }
}

func STDOUT(str string) {
    fmt.Println(str)

}
