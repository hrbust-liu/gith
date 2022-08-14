package command

import (
    "fmt"
    "os"
    "strings"

    "github.com/hrbust-liu/gith/src/command/gith"
    "github.com/hrbust-liu/gith/src/common"
)

type GitParse interface {
    Init() error
    Run() error
}

type FuncGroup struct {
    NameList []string
    FuncPtr  func(args []string) error
    Message  string
}

type GitHParse struct {
    solveMap   map[string]FuncGroup
    nameIdList []string // 每个注册的方法，记录到 list 中
    nameMap    map[string]string
    gitCommand GitCommand
}

func CreateParse() GitParse {
    return GitHParse{
        solveMap:   make(map[string]FuncGroup),
        nameIdList: []string{},
        nameMap:    make(map[string]string),
        gitCommand: gith.GitHCommand{},
    }
}

/*
目录结构
./gith
 -- config
 -- current
   -- tmpId
   -- branchId
   -- commitId
 -- data-dir
   -- HashDir
     -- HashFile

current ->
    tmpId
    branchId(包含了 commit-id)
    commitId(current) ->
        date[string]
        FileMap[fileName -> filePath]
        commidId(father)[commitFilePath]

状态转换
INit 创建目录，生成 BranchId，Current -> BranchId
Add/Remove 更新当前暂存(tmpId)内容
Commit  将 tmpId 内容生成 CommitId
        更新 branchId->commitId 和 current->commitId
Checkout 切换 Current->branchId 和 current->commitId
Reset 切换 current->commitId
Log 沿着 CommitId 查看信息
*/

func (gitParse GitHParse) Init() error {
    args := os.Args

    // debug 信息
    common.INFO(fmt.Sprintf("%T -- %v\n", args, args))
    for id := 1; id < len(args); id++ {
        common.INFO(fmt.Sprintf("%#v -- %#v\n", args, args[id]))
    }

    // 初始化全局信息
    common.GlobalConfigOnce = &common.GlobalConfig{
        FilterFile: []string{".gith", ".vscode"},
    }

    // 初始化函数

    gitParse.InitFunc([]string{"init"}, gitParse.InitInit(), "init code repo")

    gitParse.InitFunc([]string{""}, gitParse.InitInit(), "")

    gitParse.InitFunc([]string{"status", "st"}, gitParse.InitStatus(), "show repo status")
    gitParse.InitFunc([]string{"add"}, gitParse.InitAdd(), "add code file")
    gitParse.InitFunc([]string{"commit"}, gitParse.InitCommit(), "commit add file")

    gitParse.InitFunc([]string{""}, gitParse.InitInit(), "")

    gitParse.InitFunc([]string{"reset"}, gitParse.InitReset(), "reset code commit version")
    gitParse.InitFunc([]string{"branch", "br"}, gitParse.InitBranch(), "show branch info")
    gitParse.InitFunc([]string{"checkout", "co"}, gitParse.InitCheckout(), "checkout code version")
    gitParse.InitFunc([]string{"cherry-pick", "cp"}, gitParse.InitCherryPick(), "cherry-pick code")

    gitParse.InitFunc([]string{""}, gitParse.InitInit(), "")

    gitParse.InitFunc([]string{"diff"}, gitParse.InitDiff(), "diff comitId")
    gitParse.InitFunc([]string{"log"}, gitParse.InitLog(), "show log")

    gitParse.InitFunc([]string{"help"}, gitParse.InitHelp(), "")
    return nil
}

func (gitParse *GitHParse) InitFunc(nameList []string, funcName func(args []string) error, message string) {
    for id, name := range nameList {
        gitParse.nameMap[name] = nameList[0]
        if id != 0 {
            continue
        }

        gitParse.nameIdList = append(gitParse.nameIdList, name)
        gitParse.solveMap[name] = FuncGroup{
            NameList: nameList,
            FuncPtr:  funcName,
            Message:  message,
        }
    }
}

func (gitParse GitHParse) InitInit() func(args []string) error {
    return gitParse.gitCommand.Init
}

func (gitParse GitHParse) InitAdd() func(args []string) error {
    return gitParse.gitCommand.Add
}

func (gitParse GitHParse) InitCommit() func(args []string) error {
    return gitParse.gitCommand.Commit
}

func (gitParse GitHParse) InitStatus() func(args []string) error {
    return gitParse.gitCommand.Status
}

func (gitParse GitHParse) InitLog() func(args []string) error {
    return gitParse.gitCommand.Log
}

func (gitParse GitHParse) InitReset() func(args []string) error {
    return gitParse.gitCommand.Reset
}

func (gitParse GitHParse) InitBranch() func(args []string) error {
    return gitParse.gitCommand.Branch
}

func (gitParse GitHParse) InitCheckout() func(args []string) error {
    return gitParse.gitCommand.Checkout
}

func (gitParse GitHParse) InitDiff() func(args []string) error {
    return gitParse.gitCommand.Diff
}

func (gitParse GitHParse) InitCherryPick() func(args []string) error {
    return gitParse.gitCommand.CherryPick
}

func (gitParse GitHParse) InitHelp() func(args []string) error {
    return func(args []string) error {
        helpMsg := common.OutYello + "help:\n" + common.OutEnd
        for _, name := range gitParse.nameIdList {
            funcGroup := gitParse.solveMap[name]
            // }
            // for _, funcGroup := range gitParse.solveMap {

            helpMsg += "    " + common.OutBlue + strings.Join(funcGroup.NameList, "/") + common.OutEnd
            helpMsg += " " + funcGroup.Message
            helpMsg += "\n"
        }
        fmt.Println(helpMsg)
        return nil
    }
}

func (gitParse GitHParse) Run() error {
    args := os.Args

    if len(args) == 1 {
        gitParse.solveMap["help"].FuncPtr(nil)
    } else if _, ok := gitParse.solveMap[gitParse.nameMap[args[1]]]; ok {
        fun, _ := gitParse.solveMap[gitParse.nameMap[args[1]]]
        err := fun.FuncPtr(args[2:])
        if err != nil {
            common.STDOUT(err.Error())
        }
    } else {
        gitParse.solveMap["help"].FuncPtr(nil)
    }
    return nil
}
