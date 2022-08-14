package gith

import (
    "os"

    "github.com/hrbust-liu/gith/src/common"
)

func (githCommand GitHCommand) Init(args []string) error {
    _, err := os.Stat(common.GithDir)
    if err == nil || !os.IsNotExist(err) {
        common.STDOUT("Gith had Init!")
        return nil
    }
    common.Mkdir(common.GithDir)
    common.Mkdir(common.CurrentDir)
    common.Mkdir(common.MetaDataDir)
    common.Mkdir(common.GithDataDir)

    commitIdInfo := CommitIdInfo{
        CommitId:       GenCommitId(),
        date:           "",
        FileMap:        map[string]string{},
        ParentCommitId: "",
    }
    err = githCommand.fileHandle.DumpCommitIdInfo(commitIdInfo)
    if err != nil {
        return err
    }

    tmpCommitIdInfo := CommitIdInfo{
        CommitId:       GenCommitId(),
        date:           "",
        FileMap:        map[string]string{},
        ParentCommitId: "",
    }
    err = githCommand.fileHandle.DumpTmpCommitIdInfo(tmpCommitIdInfo)
    if err != nil {
        return err
    }

    err = githCommand.fileHandle.UpdateMetaBranchMap(
        MetaBranch{
            BranchMap: map[string]string{"master": commitIdInfo.CommitId},
        },
    )
    if err != nil {
        return err
    }

    err = githCommand.fileHandle.UpdateCurrentBranch("master")
    if err != nil {
        return err
    }

    common.STDOUT("Gith Init succ!")
    return nil
}
