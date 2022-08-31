package engine

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"transaction-matching-engine/models"
)

var (
	filePathPrefix = "./dump/"
)

//序列化
func serialize(data []*models.Order) []byte {
	bts, _ := json.Marshal(data)
	return bts
}

//反序列化
func deserialize(bts []byte) []*models.Order {
	data := []*models.Order{}
	json.Unmarshal(bts, &data)
	return data
}

//存储   文件存在就删除了重新创建
func Dump() {
	eg := GetMatchEngine(nil)
	os.MkdirAll(filePathPrefix, fs.ModeDir)
	for pair, pool := range eg.pools {
		bids, asks := pool.GetOrders()
		save(filePathPrefix+pair+"_bids", bids)
		save(filePathPrefix+pair+"_asks", asks)
	}
}

//存储
func save(filePath string, orders []*models.Order) {
	os.Remove(filePath)
	file, err := os.Create(filePath)
	if err != nil {
		panic("创建文件失败" + err.Error())
	}
	file.Write(serialize(orders))
}

//加载		只加载pairs内的文件
func Load(pairs []string) {
	files, err := os.ReadDir(filePathPrefix)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic("dump文件读取异常:" + err.Error())
	}
	if len(files) == 0 {
		return
	}
	eg := GetMatchEngine(pairs)
	needPair := map[string]int{}
	for k := range pairs {
		needPair[pairs[k]]++
	}
	for _, info := range files {
		nameSlice := strings.Split(info.Name(), "_")
		if len(nameSlice) != 2 {
			panic("文件名异常")
		}
		needPair[nameSlice[0]]++
	}
	//包含 bids 和 asks 的才载入
	for pair, count := range needPair {
		if count == 3 {
			eg.pools[pair].SetOrders(read(filePathPrefix+pair+"_bids"), read(filePathPrefix+pair+"_asks"))
			fmt.Println("成功加载文件：", pair)
		}
	}
}

func read(filePath string) []*models.Order {
	bts, err := os.ReadFile(filePath)
	if err != nil {
		panic("读取文件失败" + err.Error())
	}
	data := []*models.Order{}
	err = json.Unmarshal(bts, &data)
	if err != nil {
		panic("解析文件内容失败" + filePath)
	}
	return data
}
