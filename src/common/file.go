package common

import (
    "crypto/sha512"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "os"
)

func ReadJson(filePath string) (info interface{}, err error) {
    filePtr, err := os.Open(filePath)
    if err != nil {
        fmt.Println("文件打开失败 [Err:%s]", err.Error())
        return nil, err
    }
    defer filePtr.Close()

    // 创建json解码器
    decoder := json.NewDecoder(filePtr)
    err = decoder.Decode(&info)
    if err != nil {
        return nil, err
    }
    return info, err
}

func Mkdir(dirPath string) error {
    return os.Mkdir(dirPath, DirMode)
}

func WriteJson(filePath string, info interface{}) error {

    // 创建文件
    filePtr, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer filePtr.Close()

    // 创建Json编码器
    encoder := json.NewEncoder(filePtr)
    err = encoder.Encode(info)
    return err
}

func WriteData(filePath string, data interface{}) error {
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    p, err := json.Marshal(data)
    if err != nil {
        return err
    }

    file.Write(p)

    return nil
}

func WriteDataByte(filePath string, data []byte) error {
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()
    file.Write(data)
    return nil
}

func ReadDataByte(filePath string) ([]byte, error) {
    f, err := os.Open(filePath)
    if err != nil {
        return nil, nil
    }
    defer f.Close()
    return ioutil.ReadAll(f)
}

func ReadData(filePath string) (data interface{}, err error) {
    f, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    dataBytes, err := ioutil.ReadAll(f)
    json.Unmarshal(dataBytes, &data)

    return data, err
}

func Copy(srcFilePath string, dstFilePath string) error {
    sourceFileStat, err := os.Stat(srcFilePath)
    if err != nil {
        return err
    }

    if !sourceFileStat.Mode().IsRegular() {
        return fmt.Errorf("%s is not a regular file", srcFilePath)
    }

    source, err := os.Open(srcFilePath)
    if err != nil {
        return err
    }
    defer source.Close()

    destination, err := os.Create(dstFilePath)
    if err != nil {
        return err
    }
    defer destination.Close()
    _, err = io.Copy(destination, source)

    return err
}

// 生成文件校验和函数
func FileCheckSum(fileName string) (string, error) {
    f, err := os.Open(fileName)
    if err != nil {
        ERROR("FileOpen error! filePath: " + fileName + "err: " + err.Error())
        return "", err
    }

    defer f.Close()

    // h := md5.New()
    // h := sha256.New()
    // h := sha1.New()
    h := sha512.New()

    if _, err := io.Copy(h, f); err != nil {
        ERROR("ReadFile error! filePath: " + fileName + "err: " + err.Error())
        return "", err
    }

    // 格式化为16进制字符串
    return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func RemoveFile(filePath string) error {
    return os.Remove(filePath)
}

func ReadFile(filePath string) (string, error) {
    contents, err := ioutil.ReadFile(filePath)
    if err != nil {
        return "", err
    }
    return string(contents), nil
}

func WriteFile(filePath string, data string) error {
    file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
    if err != nil {
        return err
    }
    file.Write([]byte(data))
    return nil
}
