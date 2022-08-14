package gith

import (
    "os/exec"

    "github.com/hrbust-liu/gith/src/common"
)

func (gt GitHCommand) ApplyDiff(fileName string, diffFileGroup fileDiffGroup) error {
    tmpFilePath := CommitIdConvertPath(GenCommitId())
    common.WriteFile(tmpFilePath, diffFileGroup.DiffData)
    cmd := exec.Command("bash", "-c", "patch "+fileName+" "+tmpFilePath)
    cmd.Output()
    return nil

    // dmp := diffmatchpatch.New()
    // fileData := ""
    // if fileName != "" {
    //     fileData, _ = common.ReadFile(fileName)
    // }

    // patchData := diffFileGroup.DiffData
    // diffmatchpatch, _ := dmp.PatchFromText(patchData)
    // fmt.Println("+++++++++++++++++")
    // fmt.Println(patchData)
    // fmt.Println("+++++++++++++++++")
    // fmt.Println(diffmatchpatch)
    // fmt.Println("+++++++++++++++++")
    // fmt.Println(fileData)
    // fmt.Println("+++++++++++++++++")
    // newtext, _ := dmp.PatchApply(diffmatchpatch, fileData)
    // common.WriteFile(fileName, newtext)

    return nil
}
func (gt GitHCommand) CherryPick(args []string) error {
    commitId := ""
    if len(args) == 0 {
        return nil
    } else if len(args) == 1 {
        commitId = args[0]
    } else {
        return nil
    }
    // 用 diff 获取(commitId和其父Id差距)并保存
    cpCommitInfo, err := gt.fileHandle.GetCommitIdInfoByCommitId(commitId)
    if err != nil {
        return err
    }
    cppCommitInfo, err := gt.fileHandle.GetCommitIdInfoByCommitId(cpCommitInfo.ParentCommitId)
    if err != nil {
        return err
    }
    fileMap, _ := gt.fileMapDiff(cpCommitInfo.FileMap, cppCommitInfo.FileMap)

    // 将 diff 应用到当前 commitId 上
    for fileName, diffFileGroup := range fileMap {
        gt.ApplyDiff(fileName, diffFileGroup)
    }

    tmpCommitIdInfo, err := gt.fileHandle.GetTmpCommitIdInfo()
    if err != nil {
        return nil
    }

    // 生成 FileMap
    newFileMap := tmpCommitIdInfo.FileMap
    for fileName, fileGroup := range fileMap {
        if fileGroup.NewFilePath == "" {
            // 无目标文件直接删除
            delete(newFileMap, fileName)
            common.RemoveFile(fileName)
            continue
        }
        // 生成文件唯一 commitId，保存到 data，并且更新 fileMap
        newFileMap[fileName] = GenCommitId()
        gt.fileHandle.DumpDataFile(fileName, tmpCommitIdInfo.FileMap[fileName])
    }

    // 生成 commitId
    cpCommitInfo.CommitId = GenCommitId()
    currentCommitIdInfo, _ := gt.fileHandle.GetCurrentCommitIdInfo()
    cpCommitInfo.ParentCommitId = currentCommitIdInfo.CommitId
    cpCommitInfo.FileMap = newFileMap
    gt.fileHandle.UpdateCurrentBrachCommitId(cpCommitInfo.CommitId)

    err = gt.fileHandle.DumpCommitIdInfo(cpCommitInfo)
    if err != nil {
        return err
    }

    tmpCommitIdInfo.FileMap = newFileMap
    gt.fileHandle.DumpTmpCommitIdInfo(tmpCommitIdInfo)
    return nil
}
