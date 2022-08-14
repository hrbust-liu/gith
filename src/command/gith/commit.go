package gith

import (
    "flag"
    "fmt"
    "os"
)

func (githCommand GitHCommand) Commit(args []string) error {
    // 解析参数
    var cmd = flag.NewFlagSet("", flag.ExitOnError)
    var message string
    cmd.StringVar(&message, "m", "default message", "message")
    cmd.Usage = func() {
        fmt.Fprintf(os.Stderr, "Gith Usage of params:\n")
        cmd.PrintDefaults()
    }
    cmd.Parse(os.Args[2:])

    // TODO 检查暂存区与 current 之间是否有增量数据

    // 获取暂存区(tmpFile)信息, 记录 message，dump 到 data 目录。
    // 切换 branch 指向
    tmpCommitIdInfo, err := githCommand.fileHandle.GetTmpCommitIdInfo()
    if err != nil {
        return nil
    }
    tmpCommitIdInfo.Message = message
    githCommand.fileHandle.DumpCommitIdInfo(tmpCommitIdInfo)
    githCommand.fileHandle.UpdateCurrentBrachCommitId(tmpCommitIdInfo.CommitId)

    // 更新 tmp 目录
    // 必须先把 Parent 更新到当前。再生产新 CommitId
    tmpCommitIdInfo.ParentCommitId = tmpCommitIdInfo.CommitId
    tmpCommitIdInfo.Message = ""
    tmpCommitIdInfo.CommitId = GenCommitId()

    err = githCommand.fileHandle.DumpTmpCommitIdInfo(tmpCommitIdInfo)
    if err != nil {
        return err
    }
    return nil
}
