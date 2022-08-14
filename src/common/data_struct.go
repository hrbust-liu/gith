package common

func RemoveRepeatElement(list []string) []string {
    // 创建一个临时map用来存储数组元素
    temp := make(map[string]bool)
    index := 0
    for _, v := range list {
        // 遍历数组元素，判断此元素是否已经存在map中
        _, ok := temp[v]
        if ok {
            list = append(list[:index], list[index+1:]...)
        } else {
            temp[v] = true
        }
        index++
    }
    return list
}

func InList(data string, dataList []string) bool {
    for dataId := range dataList {
        if dataList[dataId] == data {
            return true
        }
    }
    return false
}

// 同时计算交集 & 差集, 时间复杂度 O(nlgn)
func SliceSetOpt(slice1, slice2 []string) ([]string, []string, []string) {
    sliceMap1 := make(map[string]int)
    sliceMap2 := make(map[string]int)
    intersectSlice := make([]string, 0)
    onlySlice1 := make([]string, 0)
    onlySlice2 := make([]string, 0)
    for _, v := range slice1 {
        sliceMap1[v]++
    }

    for _, v := range slice2 {
        sliceMap2[v]++
        _, isFound := sliceMap1[v]
        if sliceMap2[v] != 1 {
            continue
        }
        if isFound {
            intersectSlice = append(intersectSlice, v)
        } else {
            onlySlice2 = append(onlySlice2, v)
        }
    }
    sliceMap1 = make(map[string]int)

    for _, v := range slice1 {
        sliceMap1[v]++
        _, isFound := sliceMap2[v]
        if sliceMap1[v] != 1 {
            continue
        }
        if !isFound {
            onlySlice1 = append(onlySlice1, v)
        }
    }
    return onlySlice1, intersectSlice, onlySlice2
}

func GetMapKeys(m map[string]string) []string {
    j := 0
    keys := make([]string, len(m))
    for k := range m {
        keys[j] = k
        j++
    }
    return keys

}

func MapSetOpt(m1 map[string]string, m2 map[string]string) (onlyMap1 []string, commonMapValueSame []string, commonMapValueDiff []string, onlyMap2 []string) {
    // fmt.Printf("%#v \n", m1)
    // fmt.Printf("%#v \n", m2)
    for key, value := range m1 {
        res, ok := m2[key]
        if ok {
            if value == res {
                commonMapValueSame = append(commonMapValueSame, key)
            } else {
                commonMapValueDiff = append(commonMapValueDiff, key)
            }
        } else {
            onlyMap1 = append(onlyMap1, key)
        }
    }
    for key := range m2 {
        _, ok := m1[key]
        if !ok {
            onlyMap2 = append(onlyMap2, key)
        }
    }
    return onlyMap1, commonMapValueSame, commonMapValueDiff, onlyMap2
}
