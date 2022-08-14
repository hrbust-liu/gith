package gith

import (
    "encoding/hex"
    "io/ioutil"
    "math/rand"
    "os"
    "path"
    "time"

    "github.com/hrbust-liu/gith/src/common"
)

func GenCommitId() string {
    rand.Seed(time.Now().UnixNano())
    uLen := 20
    b := make([]byte, uLen)
    rand.Read(b)
    rand_str := hex.EncodeToString(b)[0:uLen]
    pathDir := CommitIdConvertDir(rand_str)

    _, err := os.Stat(pathDir)
    if os.IsNotExist(err) {
        common.Mkdir(pathDir)
    }

    path := CommitIdConvertPath(rand_str)

    _, err = os.Stat(path)
    if os.IsExist(err) {
        return GenCommitId()
    }
    return rand_str
}

func findFileList(args []string) (result []string, resErr error) {
    for pathId := range args {
        // 如果是文件，直接 add 进去
        // 如果是目录，list 目录，挨个递归
        filePath := args[pathId]
        _, fileName := path.Split(filePath)
        if common.InList(fileName, common.GlobalConfigOnce.FilterFile) {
            continue
        }
        stat, err := os.Stat(filePath)
        if os.IsNotExist(err) {
            continue
        }
        if stat.IsDir() {
            fileInfoList, err := ioutil.ReadDir(filePath)
            if err != nil {
                return nil, err
            }
            for i := range fileInfoList {
                fileList, err := findFileList([]string{filePath + "/" + fileInfoList[i].Name()})
                if err != nil {
                    return nil, err
                }
                result = append(result, fileList...)
            }
        } else {
            result = append(result, filePath)
        }
    }
    // 需要进行去重
    return common.RemoveRepeatElement(result), resErr
}
