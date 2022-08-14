package gith

import (
    "fmt"

    "github.com/hrbust-liu/gith/src/common"
)

type fileDiffGroup struct {
    FileName    string
    DiffData    string
    OldFilePath string
    NewFilePath string
}

// childMap is new commit(child modified base on parent)
func (gt GitHCommand) fileMapDiff(childMap map[string]string, parentMap map[string]string) (map[string]fileDiffGroup, error) {
    onlyChild, _, curDiffDist, onlyParent := common.MapSetOpt(
        childMap, parentMap)

    fileMap := map[string]fileDiffGroup{}
    for _, fileName := range onlyParent {
        parentFileCommitId, _ := parentMap[fileName]
        parentFilePath := CommitIdConvertPath(parentFileCommitId)

        curFilePath := ""
        diffData, err := gt.fileHandle.FileDiff(
            curFilePath, parentFilePath)
        if err != nil {
            return nil, err
        }
        fileMap[fileName] = fileDiffGroup{
            FileName:    fileName,
            DiffData:    diffData,
            OldFilePath: parentFilePath,
            NewFilePath: "",
        }
    }
    for _, fileName := range curDiffDist {

        parentFileCommitId, _ := parentMap[fileName]
        parentFilePath := CommitIdConvertPath(parentFileCommitId)

        childFileCommitId, _ := childMap[fileName]
        childFilePath := CommitIdConvertPath(childFileCommitId)
        diffData, err := gt.fileHandle.FileDiff(childFilePath, parentFilePath)
        if err != nil {
            return nil, err
        }
        fileMap[fileName] = fileDiffGroup{
            FileName:    fileName,
            DiffData:    diffData,
            OldFilePath: parentFilePath,
            NewFilePath: childFilePath,
        }
    }
    for _, fileName := range onlyChild {
        parentFilePath := ""

        childFileCommitId, _ := childMap[fileName]
        childFilePath := CommitIdConvertPath(childFileCommitId)
        diffData, err := gt.fileHandle.FileDiff(childFilePath, parentFilePath)
        if err != nil {
            return nil, err
        }
        fileMap[fileName] = fileDiffGroup{
            FileName:    fileName,
            DiffData:    diffData,
            OldFilePath: parentFilePath,
            NewFilePath: childFilePath,
        }
    }
    return fileMap, nil
}

func (githCommand GitHCommand) Diff(args []string) error {
    commitId := ""
    if len(args) > 0 {
        commitId = args[0]
    } else {
        branchCommitIdInfo, err := githCommand.fileHandle.GetCurrentCommitIdInfo()
        if err != nil {
            return err
        }
        commitId = branchCommitIdInfo.CommitId
    }
    // 获取和上一个 commitId 不同的文件
    curCommitIdInfo, err := githCommand.fileHandle.GetTmpCommitIdInfo()
    if err != nil {
        return err
    }
    dstCommitIdInfo, err := githCommand.fileHandle.GetCommitIdInfoByCommitId(commitId)
    if err != nil {
        return err
    }

    fileMap, err := githCommand.fileMapDiff(curCommitIdInfo.FileMap, dstCommitIdInfo.FileMap)
    if err != nil {
        return err
    }

    // 用 diff 工具获取每个文件的 diff
    for _, diffGroup := range fileMap {
        if diffGroup.DiffData == "" {
            continue
        }
        common.STDOUT(
            // fileName commitId..commitId
            fmt.Sprintf("[%s%s%s %s..%s]",
                common.OutYello, diffGroup.FileName, common.OutEnd,
                curCommitIdInfo.CommitId, dstCommitIdInfo.CommitId))
        fmt.Println(diffGroup.DiffData)
    }
    return nil
}
