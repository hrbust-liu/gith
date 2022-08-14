package gith

import "github.com/hrbust-liu/gith/src/common"

func (githCommand GitHCommand) Branch(args []string) error {
    // 读取 branch 文件，获取当前指向 branch，打印信息
    currentBranchName, err := githCommand.fileHandle.GetCurrentBranch()
    if err != nil {
        return err
    }

    metaBranch, err := githCommand.fileHandle.GetMetaBranchMap()
    if err != nil {
        return err
    }
    for branchName := range metaBranch.BranchMap {
        branchType := "  "
        if branchName == currentBranchName {
            branchType = " *"
        }
        common.STDOUT(branchType + branchName)
    }
    return nil
}
