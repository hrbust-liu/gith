package gith

import "github.com/hrbust-liu/gith/src/common"

func (githCommand GitHCommand) Add(args []string) error {
    if !HadInit() {
        common.STDOUT("not init!")
        return nil
    }
    // 获取要 Add fileList
    fileList, err := findFileList(args)
    if err != nil {
        return err
    }
    if len(fileList) == 0 {
        return err
    }

    // 获取暂存区(tmpFile)信息
    tmpCommitIdInfo, err := githCommand.fileHandle.GetTmpCommitIdInfo()
    if err != nil {
        return nil
    }

    // 过滤没有变化的文件
    var changedFileList []string
    for fileId := range fileList {
        fileName := fileList[fileId]
        fileCommitId, ok := tmpCommitIdInfo.FileMap[fileName]
        if !ok {
            changedFileList = append(changedFileList, fileName)
            continue
        }
        isSame, err := githCommand.fileHandle.FileIsSame(
            fileName, CommitIdConvertPath(fileCommitId))
        if err != nil {
            return err
        }
        if isSame {
            continue
        }
        changedFileList = append(changedFileList, fileName)
    }

    // 将文件记录到 data 目录
    for fileId := range changedFileList {
        fileName := changedFileList[fileId]

        tmpCommitIdInfo.FileMap[fileName] = GenCommitId()
        githCommand.fileHandle.DumpDataFile(fileName, tmpCommitIdInfo.FileMap[fileName])
    }

    // 更新暂存区(tmpFile)信息
    err = githCommand.fileHandle.DumpTmpCommitIdInfo(tmpCommitIdInfo)
    if err != nil {
        return err
    }
    return nil
}
