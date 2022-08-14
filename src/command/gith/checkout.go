package gith

import (
    "github.com/hrbust-liu/gith/src/common"
)

/*
co branchName [切换(已有)分支]
co -b branchName [新建分支 & 切换分支]
*/
func (gt GitHCommand) Checkout(args []string) error {
    // 解析参数
    createBranch := false
    branchName := ""
    if len(args) == 0 {
        return nil
    } else if len(args) == 1 {
        branchName = args[0]
    } else if len(args) == 2 {
        if args[0] == "-b" {
            createBranch = true
            branchName = args[1]
        } else {
            return nil
        }
    }

    // 当前分支的信息
    currentCommitIdInfo, err := gt.fileHandle.GetCurrentCommitIdInfo()
    if err != nil {
        return err
    }
    tmpCommitIdInfo, err := gt.fileHandle.GetTmpCommitIdInfo()
    if err != nil {
        return err
    }

    if createBranch {
        // 新建分支 branchMap 生成新的 commitId 节点
        metaBrach, err := gt.fileHandle.GetMetaBranchMap()
        if err != nil {
            return err
        }
        _, found := metaBrach.BranchMap[branchName]
        if found {
            common.STDOUT("Found Branch " + branchName + "!")
            return nil
        }
        metaBrach.BranchMap[branchName] = currentCommitIdInfo.CommitId
        gt.fileHandle.UpdateMetaBranchMap(metaBrach)
    } else {
        // 非新建分支
        // 检查 tmpCommitIdInfo 必须和 branchCommitIdInfo 一致
        currentMap, _, modefiedMap, tmpMap := common.MapSetOpt(
            currentCommitIdInfo.FileMap, tmpCommitIdInfo.FileMap)

        if len(currentMap) != 0 || len(modefiedMap) != 0 || len(tmpMap) != 0 {
            common.STDOUT("Has File Modefied!\nCan't checkout branch")
        }

        // 将 current 切换到目标分支的 commitIdInfo

        // 查询分支映射数据
        metaBranch, err := gt.fileHandle.GetMetaBranchMap()
        if err != nil {
            return err
        }

        // 查询分支 commit 信息
        branchId, found := metaBranch.BranchMap[branchName]
        if !found {
            common.STDOUT("Must First Create NewBranch: " + branchName)
            return nil
        }
        commitIdInfo, err := gt.fileHandle.GetCommitIdInfoByCommitId(branchId)
        if err != nil {
            return err
        }
        currentCommitIdInfo = commitIdInfo
    }

    // tmp 切换到新的 FileMap
    gt.CheckoutFileMap(tmpCommitIdInfo.FileMap, currentCommitIdInfo.FileMap)
    tmpCommitIdInfo.FileMap = currentCommitIdInfo.FileMap
    gt.fileHandle.DumpTmpCommitIdInfo(tmpCommitIdInfo)

    // 分支切换
    err = gt.fileHandle.UpdateCurrentBranch(branchName)
    if err != nil {
        return err
    }

    return nil
}
