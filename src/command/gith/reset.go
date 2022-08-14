package gith

import (
    "flag"
    "fmt"
    "os"

    "github.com/hrbust-liu/gith/src/common"
)

// 代码切换到目标 map
func (githCommand GitHCommand) CheckoutFileMap(
    curFileMap map[string]string, dstFileMap map[string]string) error {
    onlyDst, _, curDiffDist, onlyCur := common.MapSetOpt(
        dstFileMap, curFileMap)

    for _, fileName := range onlyDst {
        fileCommitId := dstFileMap[fileName]
        common.Copy(CommitIdConvertPath(fileCommitId), fileName)
    }
    for _, fileName := range curDiffDist {
        fileCommitId := dstFileMap[fileName]
        common.Copy(CommitIdConvertPath(fileCommitId), fileName)
    }
    for _, fileName := range onlyCur {
        common.RemoveFile(fileName)
    }
    return nil
}

func (gt GitHCommand) Reset(args []string) error {
    // 解析参数
    var cmd = flag.NewFlagSet("", flag.ExitOnError)
    var hardCommit, softCommit string
    cmd.StringVar(&hardCommit, "hard", "", "message")
    cmd.StringVar(&softCommit, "soft", "", "message")
    cmd.Usage = func() {
        fmt.Fprintf(os.Stderr, "Gith Usage of params:\n")
        cmd.PrintDefaults()
    }
    cmd.Parse(os.Args[2:])

    isHard := false
    dstCommitId := ""
    if hardCommit != "" {
        isHard = true
        dstCommitId = hardCommit
    } else if softCommit != "" {
        isHard = false
        dstCommitId = softCommit
    } else {
        return nil
    }

    common.STDOUT("Skip to " + dstCommitId)
    // TODO 判断是否存在

    // 通过 Map 区别，更新本地文件
    dstCommitIdInfo, err := gt.fileHandle.GetCommitIdInfoByCommitId(dstCommitId)
    if err != nil {
        return err
    }
    tmpCommitIdInfo, err := gt.fileHandle.GetTmpCommitIdInfo()
    if err != nil {
        return err
    }
    if isHard {
        // 文件需要切换好
        gt.CheckoutFileMap(tmpCommitIdInfo.FileMap, dstCommitIdInfo.FileMap)

        // 切换 FileMap
        tmpCommitIdInfo.FileMap = dstCommitIdInfo.FileMap
    }
    tmpCommitIdInfo.date = common.GetCurrentDate()
    gt.fileHandle.DumpTmpCommitIdInfo(tmpCommitIdInfo)
    gt.fileHandle.UpdateCurrentBrachCommitId(dstCommitIdInfo.CommitId)

    return nil
}
