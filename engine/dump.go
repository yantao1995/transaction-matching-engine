package engine

import (
	"bytes"
	"encoding/gob"
	"transaction-matching-engine/pool"
)

var (
	filePathPrefix = "./dump/"
)

type dump struct {
}

//序列化
func (dp *dump) Serialize(mp *pool.MatchPool) []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(mp); err != nil {
		panic("序列化失败:" + err.Error())
	}
	return buffer.Bytes()
}

//反序列化
func (dp *dump) Deserialize(data []byte) *pool.MatchPool {
	mp := &pool.MatchPool{}
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(mp); err != nil {
		panic("反序列化失败:" + err.Error())
	}
	return mp
}

// //存储   文件存在就删除了重新创建
// func Dump() {
// 	dp := &dump{}
// 	eg := GetMatchEngine(nil)
// 	os.MkdirAll(filePathPrefix, fs.ModeDir)
// 	for pair, pool := range eg.pools {
// 		filePath := filePathPrefix + pair
// 		os.Remove(filePath)
// 		file, err := os.Create(filePath)
// 		if err != nil {
// 			panic("创建文件失败" + err.Error())
// 		}
// 		file.Write(dp.Serialize(pool))
// 	}
// }

//加载		只加载pairs内的文件
// func Load(pairs []string) {
// 	files, err := os.ReadDir(filePathPrefix)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return
// 		}
// 		panic("dump文件读取异常:" + err.Error())
// 	}
// 	if len(files) == 0 {
// 		return
// 	}
// 	dp := &dump{}
// 	eg := GetMatchEngine(pairs)
// 	needPair := map[string]bool{}
// 	for k := range pairs {
// 		needPair[pairs[k]] = true
// 	}
// 	for _, info := range files {
// 		if needPair[info.Name()] {
// 			filePath := filePathPrefix + info.Name()
// 			bts, err := os.ReadFile(filePath)
// 			if err != nil {
// 				panic("读取文件失败" + err.Error())
// 			}
// 			eg.pools[info.Name()] = dp.Deserialize(bts)
// 		}
// 	}
// }
