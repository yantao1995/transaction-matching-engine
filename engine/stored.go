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
	bidsName       = "/bids.json"
	asksName       = "/asks.json"
	pairName       = "pairs.json"
)

//序列化
func serialize(data interface{}) []byte {
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
	pairs := []string{} //记录本次启动的pairs
	for pair, pool := range eg.pools {
		pairs = append(pairs, pair)
		bids, asks := pool.GetOrders()
		os.MkdirAll(filePathPrefix+pair, fs.ModeDir)
		save(filePathPrefix+pair+bidsName, bids)
		save(filePathPrefix+pair+asksName, asks)
	}
	save(filePathPrefix+pairName, pairs)
}

//存储
func save(filePath string, data interface{}) {
	os.Remove(filePath)
	file, err := os.Create(filePath)
	if err != nil {
		panic("创建文件失败" + err.Error())
	}
	file.Write(serialize(data))
}

//加载		只加载pairs内的文件
func Load(pairs []string) {
	files, err := os.ReadDir(filePathPrefix)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(filePathPrefix + "内的dump文件读取异常:" + err.Error())
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
		if info.IsDir() {
			needPair[info.Name()]++
		}
	}
	//包含 启动pair && 文件存在 的才载入
	for pair, count := range needPair {
		if count == 2 {
			eg.pools[pair].SetOrders(readOrders(filePathPrefix+pair+bidsName), readOrders(filePathPrefix+pair+asksName))
			fmt.Println("成功加载文件：", pair)
		}
	}
}

//读取存储的订单数据
func readOrders(filePath string) []*models.Order {
	bts, err := os.ReadFile(filePath)
	if err != nil {
		panic(filePath + "读取文件失败" + err.Error())
	}
	data := []*models.Order{}
	err = json.Unmarshal(bts, &data)
	if err != nil {
		panic(filePath + "解析文件内容失败" + filePath)
	}
	return data
}

//记录上次启动的pair
func ReadPairs() []string {
	filePath := filePathPrefix + pairName
	data := []string{}
	bts, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(filePath + "读取文件失败" + err.Error())
		return data
	}
	err = json.Unmarshal(bts, &data)
	if err != nil {
		panic(filePath + "解析文件内容失败" + filePath)
	}
	for k := range data {
		data[k] = strings.ToUpper(data[k])
	}
	return data
}
