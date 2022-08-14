package gith

import (
    "fmt"

    "github.com/hrbust-liu/gith/src/common"
)

/*
   CommitId: XXXXX
   Message: ???
   date: XX-XX-XX

   CommitId: XXXXX
   ....
*/
func (githCommand GitHCommand) Log(args []string) error {
    commitIdInfo, err := githCommand.fileHandle.GetCurrentCommitIdInfo()
    fmt.Printf("%#v\n", commitIdInfo)
    if err != nil {
        return err
    }
    for {
        common.STDOUT(
            fmt.Sprintf("    %sCommitId%s: %s\n", common.OutYello, common.OutEnd, commitIdInfo.CommitId) +
                fmt.Sprintf("    Date: %s\n", commitIdInfo.date) +
                fmt.Sprintf("    Msg: %s\n", commitIdInfo.Message))
        if commitIdInfo.ParentCommitId == "" {
            break
        }

        commitIdInfo, err = githCommand.fileHandle.GetCommitIdInfoByCommitId(commitIdInfo.ParentCommitId)
        if err != nil {
            return err
        }
    }
    return nil
}
