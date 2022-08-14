package gith

import (
    "encoding/json"
    "os/exec"
    "strings"

    "github.com/hrbust-liu/gith/src/common"
)

type GitHFileHandle struct {
}

type CommitIdInfo struct {
    CommitId       string
    date           string
    Message        string
    FileMap        map[string]string
    ParentCommitId string
}

func CommitIdConvertDir(commitId string) string {
    pathDir := common.GithDataDir + "/" + commitId[:common.CommitDirLen]
    return pathDir
}

func CommitIdConvertPath(commitId string) string {
    path := CommitIdConvertDir(commitId) + "/" + commitId[common.CommitDirLen:]
    return path
}

type MetaBranch struct {
    BranchMap map[string]string
}

// Current
func (g GitHFileHandle) UpdateCurrentBranch(branch string) error {
    return common.WriteData(common.CurrentBranch, branch)
}

func (g GitHFileHandle) GetCurrentBranch() (string, error) {
    data, err := common.ReadData(common.CurrentBranch)
    return data.(string), err
}

func (g GitHFileHandle) DumpTmpCommitIdInfo(commitIdInfo CommitIdInfo) error {
    return g.DumpInfo(common.CurrentTmp, commitIdInfo)
}

func (g GitHFileHandle) GetTmpCommitIdInfo() (CommitIdInfo, error) {
    return g.GetCommitIdInfoByPath(common.CurrentTmp)
}

func (g GitHFileHandle) GetCurrentCommitIdInfo() (commitIdInfo CommitIdInfo, err error) {
    // 查询分支名
    branchName, err := g.GetCurrentBranch()
    if err != nil {
        return commitIdInfo, err
    }
    // 查询分支映射数据
    metaBranch, err := g.GetMetaBranchMap()
    if err != nil {
        return commitIdInfo, err
    }

    // 查询分支 commit 信息
    branchId, found := metaBranch.BranchMap[branchName]
    if !found {
        panic("[Data Corrupt!] Can't found branch! branchName: " + branchName)
    }
    commitIdInfo, err = g.GetCommitIdInfoByCommitId(branchId)
    return commitIdInfo, err
}

func (g GitHFileHandle) GetCommitIdInfoByCommitId(commitId string) (CommitIdInfo, error) {
    commitIdPath := CommitIdConvertPath(commitId)
    return g.GetCommitIdInfoByPath(commitIdPath)
}

// Meta
func (g GitHFileHandle) UpdateCurrentBrachCommitId(commitId string) (err error) {
    // 查询分支名
    branchName, err := g.GetCurrentBranch()
    if err != nil {
        return err
    }
    // 查询分支映射数据
    metaBranch, err := g.GetMetaBranchMap()
    if err != nil {
        return err
    }
    metaBranch.BranchMap[branchName] = commitId
    return g.UpdateMetaBranchMap(metaBranch)
}

func (g GitHFileHandle) UpdateMetaBranchMap(metaBranch MetaBranch) (err error) {
    dataBytes, err := json.Marshal(&metaBranch)
    return common.WriteDataByte(common.BranchData, dataBytes)
}

func (g GitHFileHandle) GetMetaBranchMap() (metaBranch MetaBranch, err error) {
    dataBytes, err := common.ReadDataByte(common.BranchData)

    metaBranch = MetaBranch{}
    json.Unmarshal(dataBytes, &metaBranch)
    return metaBranch, err
}

// Data
func (g GitHFileHandle) DumpDataFile(fileName string, commitId string) error {
    dstFilePath := CommitIdConvertPath(commitId)
    return common.Copy(fileName, dstFilePath)
}

// Common
func (g GitHFileHandle) DumpCommitIdInfo(commitIdInfo CommitIdInfo) error {
    path := CommitIdConvertPath(commitIdInfo.CommitId)
    return g.DumpInfo(path, commitIdInfo)
}

func (g GitHFileHandle) DumpInfo(path string, commitIdInfo CommitIdInfo) error {
    return g.Serialize(path, commitIdInfo)
}

func (g GitHFileHandle) Serialize(filePath string, commitIdInfo CommitIdInfo) (err error) {
    dataBytes, err := json.Marshal(&commitIdInfo)
    return common.WriteDataByte(filePath, dataBytes)
}

func (g GitHFileHandle) GetCommitIdInfoByPath(filePath string) (CommitIdInfo, error) {
    return g.UnSerialize(filePath)
}

func (g GitHFileHandle) UnSerialize(filePath string) (CommitIdInfo, error) {
    dataBytes, err := common.ReadDataByte(filePath)

    commitIdInfo := CommitIdInfo{FileMap: map[string]string{}}
    json.Unmarshal(dataBytes, &commitIdInfo)
    return commitIdInfo, err
}

func (g GitHFileHandle) FileIsSame(srcFilePath string, dstFilePath string) (bool, error) {
    // 调用校验和方法生成校验字符串
    fileSum1, err := common.FileCheckSum(srcFilePath)
    if err != nil {
        common.ERROR("fileCheckSum error! filePath: " + srcFilePath + "err: " + err.Error())
        return false, err
    }
    fileSum2, err := common.FileCheckSum(dstFilePath)
    if err != nil {
        common.ERROR("fileCheckSum error! filePath: " + dstFilePath + "err: " + err.Error())
        return false, err
    }

    // 比较校验和字符串
    if strings.Compare(fileSum1, fileSum2) != 0 {
        return false, nil
    }
    return true, nil
}

// 过滤没有变化的文件
func (g GitHFileHandle) FilterUnChangedFileList(fileList []string, fileMap map[string]string) ([]string, error) {

    var changedFileList []string
    for fileId := range fileList {
        fileName := fileList[fileId]
        fileCommitId, ok := fileMap[fileName]
        if !ok {
            changedFileList = append(changedFileList, fileName)
            continue
        }
        isSame, err := g.FileIsSame(fileName, CommitIdConvertPath(fileCommitId))
        if err != nil {
            return nil, err
        }
        if isSame {
            continue
        }
        changedFileList = append(changedFileList, fileName)
    }
    return changedFileList, nil
}

// newfile -> oldfile 的增量
func (g GitHFileHandle) FileDiff(newfile string, oldfile string) (string, error) {
    cmd := exec.Command("bash", "-c", "diff "+oldfile+" "+newfile)
    out, _ := cmd.Output()
    return string(out), nil

    // dmp := diffmatchpatch.New()
    // newText := ""
    // oldText := ""
    // if newfile != "" {
    //     newText, _ = common.ReadFile(newfile)
    // }
    // if oldfile != "" {
    //     oldText, _ = common.ReadFile(oldfile)
    // }
    // diffs := dmp.DiffMain(oldText, newText, false)

    // fmt.Println("---------")
    // fmt.Println(oldText)
    // fmt.Println("---------")
    // fmt.Println(newText)
    // fmt.Println("---------")
    // fmt.Println(diffs)
    // fmt.Println("---------")
    // text := dmp.DiffPrettyText(diffs)
    // fmt.Println(text)
    // fmt.Println("---------")
    // diffmatchpatch, _ := dmp.PatchFromText(text)
    // fmt.Println(diffmatchpatch)
    // fmt.Println("---------")

    // // newtext, _ := dmp.PatchApply(diffs, oldText)
    // // fmt.Println(newtext)
    // fmt.Println("---------")

    // return dmp.DiffPrettyText(diffs), nil
}
