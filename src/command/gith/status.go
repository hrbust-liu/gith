package gith

import "github.com/hrbust-liu/gith/src/common"

func (githCommand GitHCommand) Status(args []string) error {
    if !HadInit() {
        common.STDOUT("not init!")
        return nil
    }
    // 拿到 . 下所有文件列表
    fileList, err := findFileList([]string{"."})
    if err != nil {
        return err
    }
    if len(fileList) == 0 {
        return err
    }

    // 跟 Map 对比，是否有新增或修改的文件
    // 获取暂存区(tmpFile)信息
    tmpCommitIdInfo, err := githCommand.fileHandle.GetTmpCommitIdInfo()
    if err != nil {
        return nil
    }

    // 展示暂存区与上一个日志节点的 diff
    currentCommitIdInfo, err := githCommand.fileHandle.GetCurrentCommitIdInfo()
    currentMap, _, modefiedMap, tmpMap := common.MapSetOpt(
        currentCommitIdInfo.FileMap, tmpCommitIdInfo.FileMap)
    if len(currentMap) == 0 && len(tmpMap) == 0 && len(modefiedMap) == 0 {
        common.STDOUT(common.OutBlue + "[Temp] Not Changed!" + common.OutEnd)
    } else {
        common.STDOUT(common.OutBlue + "[Temp] ChangeFile:" + common.OutEnd)
        if len(tmpMap) != 0 {
            changeType := common.OutGreen + "add" + common.OutEnd
            for _, fileName := range tmpMap {
                common.STDOUT(
                    "    " + changeType + ": " + fileName)
            }
        }
        if len(modefiedMap) != 0 {
            changeType := common.OutBlue + "modefiled" + common.OutEnd
            for _, fileName := range modefiedMap {
                common.STDOUT(
                    "    " + changeType + ": " + fileName)
            }
        }
        if len(currentMap) != 0 {
            changeType := common.OutRed + "del" + common.OutEnd
            for _, fileName := range currentMap {
                common.STDOUT(
                    "    " + changeType + ": " + fileName)
            }
        }

    }

    AddFileList, commonFileList, removeFileList :=
        common.SliceSetOpt(fileList, common.GetMapKeys(tmpCommitIdInfo.FileMap))
    changedFileList, err := githCommand.fileHandle.FilterUnChangedFileList(
        commonFileList, tmpCommitIdInfo.FileMap)
    if err != nil {
        return err
    }

    if len(AddFileList) == 0 && len(changedFileList) == 0 && len(removeFileList) == 0 {
        common.STDOUT(common.OutBlue + "[UnStoreArea] Not Changed!" + common.OutEnd)
    } else {
        common.STDOUT(common.OutBlue + "[UnStoreArea] Changed!" + common.OutEnd)
        if len(AddFileList) != 0 {
            changeType := common.OutGreen + "add" + common.OutEnd
            for _, fileName := range AddFileList {
                common.STDOUT(
                    "    " + changeType + ": " + fileName)
            }
        }
        if len(changedFileList) != 0 {
            changeType := common.OutBlue + "modefiled" + common.OutEnd
            for _, fileName := range changedFileList {
                common.STDOUT(
                    "    " + changeType + ": " + fileName)
            }
        }
        if len(removeFileList) != 0 {
            changeType := common.OutRed + "del" + common.OutEnd
            for _, fileName := range removeFileList {
                common.STDOUT(
                    "    " + changeType + ": " + fileName)
            }

        }
    }

    return nil
}
