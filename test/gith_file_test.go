package test

import (
    "fmt"
    "testing"

    "github.com/hrbust-liu/gith/src/command"
    "github.com/hrbust-liu/gith/src/common"
)

func TestSerialize(t *testing.T) {
    gitHFileHandle := command.GitHFileHandle{}
    filePath := t.TempDir() + "/tmpfile"

    commitIdInfo := command.CommitIdInfo{
        CommitId: filePath,
        FileMap:  map[string]string{"A": "B"},
    }
    gitHFileHandle.Serialize(filePath, commitIdInfo)
    res, err := gitHFileHandle.UnSerialize(filePath)
    common.AssertEqual(t, err, nil)
    fmt.Printf("%T %#v\n", res, res)
    fmt.Printf("%T %#v\n", commitIdInfo, commitIdInfo)
    commitIdInfo = command.CommitIdInfo{}
    fmt.Printf("----\n")
    common.AssertEqual(t, res, commitIdInfo)
}
