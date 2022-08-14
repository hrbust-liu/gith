package common

const (
    GithDir       = ".gith"                     // 全部 gith 信息
    CurrentConfig = GithDir + "/config"         // gith 配置文件
    CurrentDir    = GithDir + "/current"        // current 信息目录
    CurrentBranch = CurrentDir + "/branch"      // current 指向分支
    CurrentTmp    = CurrentDir + "/tmp"         // current 指向分支
    MetaDataDir   = GithDir + "/meta"           // 元数据目录
    BranchData    = MetaDataDir + "/branchFile" // 分支信息
    GithDataDir   = GithDir + "/data"           // commitId，data 信息
    CommitDirLen  = 2

    DirMode = 0755

    // 输出
    OutRed    = "\033[1;31m"
    OutGreen  = "\033[1;32m"
    OutYello  = "\033[1;33m"
    OutBlue   = "\033[1;34m"
    OutPurple = "\033[1;35m"
    OutCyan   = "\033[1;36m"
    OutEnd    = "\033[m"
)

type GlobalConfig struct {
    // 规则
    FilterFile []string
}

var GlobalConfigOnce *GlobalConfig
