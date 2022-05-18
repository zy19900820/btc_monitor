package main

import (
	"encoding/json"
	"golang.org/x/xerrors"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

const NORMAL = "震荡"
const BUY = "买入"
const SELL = "卖出"

const RICH_200_DATA_DIR = "/root/rich200"

type RICH_200_DATE struct {
	Position []float64  //大户近期持仓数组
}

func getLatestFil() (string, error) {
	fileInfoList, err := ioutil.ReadDir(RICH_200_DATA_DIR)
	if err != nil {
		log.Println(err)
		return "", err
	}
	//fmt.Println(len(fileInfoList))
	sort.Slice(fileInfoList, func(i, j int) bool{
		if fileInfoList[i].ModTime().Unix() > fileInfoList[j].ModTime().Unix() {
			return true
		} else {
			return false
		}
	})

	if len(fileInfoList) >= 1 {
		return fileInfoList[0].Name(), nil
	}

	return "", xerrors.New("fil not exist")
}

func getSecodeLatestFil() (string, error) {
	fileInfoList, err := ioutil.ReadDir(RICH_200_DATA_DIR)
	if err != nil {
		log.Println(err)
		return "", err
	}
	//fmt.Println(len(fileInfoList))
	sort.Slice(fileInfoList, func(i, j int) bool{
		if fileInfoList[i].ModTime().Unix() > fileInfoList[j].ModTime().Unix() {
			return true
		} else {
			return false
		}
	})

	if len(fileInfoList) >= 2 {
		return fileInfoList[1].Name(), nil
	}

	return "", xerrors.New("fil not exist")
}

func loadRichFile(fileName string) (RICH_200_DATE, error) {
	var richData RICH_200_DATE
	filePtr, err := os.Open(fileName)
	if err != nil {
		return richData, err
	}
	defer filePtr.Close()

	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&richData)
	if err != nil {
		return richData, err
	}
	return richData, nil
}

//获取大于等于12组的最新数据
func getFullRichData() (RICH_200_DATE, error) {
	//获取最新的文件名
	latestFileName, err := getLatestFil()
	if err != nil {
		return RICH_200_DATE{}, err
	}
	richDataFirst, err := loadRichFile(RICH_200_DATA_DIR + "/" + latestFileName)
	if err != nil {
		return richDataFirst, err
	}
	if len(richDataFirst.Position) >= 12 {
		return richDataFirst, nil
	}

	//数据小于12组 则重新加载上一个文件
	var richData RICH_200_DATE
	latestSecondFileName, err := getSecodeLatestFil()
	if err != nil {
		return richData, err
	}
	richDataSecode, err := loadRichFile(RICH_200_DATA_DIR + "/" + latestSecondFileName)
	if err != nil {
		return richData, err
	}
	//倒数第二新的数据+最新的数据组合
	for i := 287; i < 300; i++ {
		richData.Position = append(richData.Position, richDataSecode.Position[i])
	}
	for i := 0; i < len(richDataFirst.Position); i++ {
		richData.Position = append(richData.Position, richDataFirst.Position[i])
	}

	return richData, nil
}

func getStatus(richData RICH_200_DATE) string {
	if len(richData.Position) != 12 {
		return "持仓数据不为12组"
	}


	buyNum := 0
	sellNum := 0

	for i := 1; i < 12; i++ {
		if richData.Position[i] - richData.Position[i - 1] > 10 {
			buyNum++
		} else if richData.Position[i] - richData.Position[i - 1] < -10 {
			sellNum++
		}
	}
	if buyNum > 7 {
		return BUY
	} else if sellNum > 7 {
		return SELL
	} else {
		return NORMAL
	}
}

func writeRichStruct(fileName string, data RICH_200_DATE) error {
	//写入最新的数据
	filePtr, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer filePtr.Close()

	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func writeStatusInfo() error {
	latestFileName, err := getLatestFil()
	if err != nil {
		if strings.Contains(err.Error(), "fil not exist") {
			var data RICH_200_DATE
			data.Position = append(data.Position, g_latestAddrInfo.Total200)
			err = writeRichStruct(RICH_200_DATA_DIR + "/" + time.Now().Format(time.RFC3339), data)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	richDataFirst, err := loadRichFile(RICH_200_DATA_DIR + "/" + latestFileName)
	if err != nil {
		return err
	}

	//写入新文件
	if len(richDataFirst.Position) == 300 {
		var data RICH_200_DATE
		data.Position = append(data.Position, g_latestAddrInfo.Total200)
		err = writeRichStruct(RICH_200_DATA_DIR + "/" + time.Now().Format(time.RFC3339), data)
		if err != nil {
			return err
		}
		return nil
	} else {
		//添加老文件写入
		richDataFirst.Position = append(richDataFirst.Position, g_latestAddrInfo.Total200)
		err = writeRichStruct(RICH_200_DATA_DIR + "/" + latestFileName, richDataFirst)
		if err != nil {
			return err
		}
		return nil
	}
}