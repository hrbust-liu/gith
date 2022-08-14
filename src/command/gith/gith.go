package gith

import (
    "os"

    "github.com/hrbust-liu/gith/src/common"
)

type GitHCommand struct {
    fileHandle GitHFileHandle
}

func HadInit() bool {
    _, err := os.Stat(common.GithDir)
    if err == nil || !os.IsNotExist(err) {
        return true
    }
    return false
}
