package main

import (
    "github.com/hrbust-liu/gith/src/command"
)

func main() {
    gt := command.CreateParse()
    gt.Init()
    gt.Run()
}
