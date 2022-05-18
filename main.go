package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/xerrors"
)

type LOCAL_BTC_ADDR_INFO struct {
	Addr    string  //钱包地址
	Count   float64 //数量
	Ranking int     //排名
	Alias   string  //别名
}

type BTC_INFO struct {
	LocalAddrInfos []LOCAL_BTC_ADDR_INFO
	Total10        float64
	Total20        float64
	Total50        float64
	Total100       float64
	Total200       float64
}

//func first() {
//	btcInfo := BTC_INFO{}
//	btcInfo.Count = 0
//
//	btcInfo.ExchangeCodeAddr = append(btcInfo.ExchangeCodeAddr, "34xp4vRoCGJym3xR7yCVPFHoCNxv4Twseo")
//	btcInfo.ExchangeCodeAddr = append(btcInfo.ExchangeCodeAddr, "35hK24tcLEWcgNA4JxpvbkNkoAcDGqQPsP")
//	btcInfo.ExchangeCodeAddr = append(btcInfo.ExchangeCodeAddr, "385cR5DM96n1HvBDMzLHPYcw89fZAXULJP")
//	btcInfo.ExchangeCodeAddr = append(btcInfo.ExchangeCodeAddr, "1AnwDVbwsLBVwRfqN2x9Eo4YEJSPXo2cwG")
//	btcInfo.ExchangeCodeAddr = append(btcInfo.ExchangeCodeAddr, "14eQD1QQb8QFVG8YFwGz7skyzsvBLWLwJS")
//
//	btcInfo.ExchangeHotAddr = append(btcInfo.ExchangeHotAddr, "1NDyJtNTjmwk5xPNhjgAMu4HDHigtobu1s")
//	btcInfo.ExchangeHotAddr = append(btcInfo.ExchangeHotAddr, "3Kzh9qAqVWQhEsfQz7zEQL1EuSx5tyNLNS")
//
//	buf := getPage()
//	//fmt.Println(buf)
//
//	//解析正则表达式，如果成功返回解释器
//	buf = strings.Replace(buf, "\n", " ", -1)
//	//fmt.Println(buf)
//	reg1 := regexp.MustCompile(`<tr>.*?</tr>`)
//	if reg1 == nil {
//		fmt.Println("regexp err")
//		return
//	}
//	//根据规则提取关键信息
//	result1 := reg1.FindAllStringSubmatch(buf, -1)
//	result1 = result1[1 : 101]
//
//	for i, result := range result1 {
//		var btcAddrInfo BTC_ADDR_INFO
//		btcAddrInfo.Ranking = i + 1
//		btcAddrInfo.Addr = getAddr(result[0])
//		btcAddrInfo.Count = getCoin(btcAddrInfo.Addr)
//		btcInfo.AddrInfos = append(btcInfo.AddrInfos, btcAddrInfo)
//		//fmt.Println(result)
//	}
//
//	filePtr, err := os.Create("info.json")
//	if err != nil {
//		fmt.Println("Create file failed", err.Error())
//		return
//	}
//	defer filePtr.Close()
//
//	encoder := json.NewEncoder(filePtr)
//
//	err = encoder.Encode(btcInfo)
//	if err != nil {
//		fmt.Println("Encoder failed", err.Error())
//
//	} else {
//		fmt.Println("Encoder success")
//	}
//}

func getPage(i int) (string, error) {
	url := "https://btc.tokenview.com/api/address/richrange/btc/" + strconv.Itoa(i) + "/10"
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", xerrors.New("http error")
	}
	return string(robots), nil
}

func getTotal() int {
	res, err := http.Get("https://blockchain.info/q/totalbc")
	if err != nil {
		return 0
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return 0
	}
	totalStr := string(robots)[0 : len(robots)-8]
	total, _ := strconv.Atoi(totalStr)
	return total
}

const DING_ERROR_URL = "https://oapi.dingtalk.com/robot/send?access_token=d5e507b4ef9122b60af42e7418050a67eb43d797315a2504469de01570684ba9"
const ERROR_WORD = "提醒 "

const DING_BOBO_URL = "https://oapi.dingtalk.com/robot/send?access_token=ba1edbe70946dace18b0e65f2a4245b4a33481952f49596792183ce281cb0125"
const BOBO_WORD = "来看看呗 "

func dingInfo(msg string, bAll bool) {
	dingInfo1(msg, bAll)
	//dingInfo2(msg, bAll)
}

func dingInfo1(msg string, bAll bool) error {
	client := &http.Client{}
	type sendDing struct {
		Msgtype string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
		At struct {
			IsAtAll bool `json:"isAtAll"`
		} `json:"at"`
	}
	var cmd sendDing
	cmd.Msgtype = "text"
	cmd.Text.Content = ERROR_WORD + msg
	cmd.At.IsAtAll = bAll
	a, _ := json.Marshal(cmd)
	req, err := http.NewRequest("POST", DING_ERROR_URL,
		strings.NewReader(string(a)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	res := struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}
	if res.Errmsg == "ok" {
		return nil
	} else {
		return xerrors.New(res.Errmsg)
	}
}

func dingInfo2(msg string, bAll bool) error {
	client := &http.Client{}
	type sendDing struct {
		Msgtype string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
		At struct {
			IsAtAll bool `json:"isAtAll"`
		} `json:"at"`
	}
	var cmd sendDing
	cmd.Msgtype = "text"
	cmd.Text.Content = BOBO_WORD + msg
	cmd.At.IsAtAll = bAll
	a, _ := json.Marshal(cmd)
	req, err := http.NewRequest("POST", DING_BOBO_URL,
		strings.NewReader(string(a)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	res := struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}
	if res.Errmsg == "ok" {
		return nil
	} else {
		return xerrors.New(res.Errmsg)
	}
}

func getCoin(addr string) int {
	res, err := http.Get("https://btc.com/" + addr)
	if err != nil {
		return 0
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return 0
	}

	page := strings.Replace(string(robots), "\n", " ", -1)
	index := strings.Index(page, ">Balance<")

	indexNextFirst := strings.Index(page[index:], "<dd>")
	indexNext := strings.Index(page[index+indexNextFirst:], "<span")

	coinNumStr := page[index+indexNextFirst+4 : index+indexNextFirst+indexNext]
	coinNumStr = strings.Replace(coinNumStr, " ", "", -1)
	coinNumStr = strings.Replace(coinNumStr, ",", "", -1)

	indexDot := strings.Index(coinNumStr, ".")
	coinNumStr = coinNumStr[:indexDot]
	coinNum, _ := strconv.Atoi(coinNumStr)
	return coinNum
}

func getAddr(result string) string {
	//fmt.Println(result)
	index1 := strings.Index(result, "btc.com/")
	index2 := strings.Index(result[index1:], "\"")
	addrStr := result[index1+8 : index1+index2]
	addrStr = strings.Replace(addrStr, " ", "", -1)
	return addrStr
}

func getRate(num, total float64) string {
	rate := float64(num) / float64(total) * 100
	rateStr := strconv.FormatFloat(rate, 'f', -1, 64)
	return rateStr + "%"
}

func loadInfo() error {
	filePtr, err := os.Open(jsonPath)
	if err != nil {
		return err
	}
	defer filePtr.Close()

	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&g_OldAddrInfo)
	if err != nil {
		dingInfo(err.Error(), false)
	}
	return nil
}

func sendChange(latestInfo LOCAL_BTC_ADDR_INFO) {
	for _, oldInfo := range g_OldAddrInfo.LocalAddrInfos {
		//币安热钱包 跳过
		if latestInfo.Addr == "1NDyJtNTjmwk5xPNhjgAMu4HDHigtobu1s" {
			return
		}

		if oldInfo.Addr == latestInfo.Addr {
			if (latestInfo.Count-oldInfo.Count) > 10 || (latestInfo.Count-oldInfo.Count) < -10 {
				alias := latestInfo.Alias
				if alias == "" {
					alias = "大佬钱包"
				}

				msg := "addr:" + latestInfo.Addr + " 数量变化:" + strconv.FormatFloat(oldInfo.Count, 'f', -1, 64) + "->" + strconv.FormatFloat(latestInfo.Count, 'f', -1, 64)
				if (latestInfo.Count - oldInfo.Count) > 10 {
					msg += " 增持:" + strconv.FormatFloat(latestInfo.Count-oldInfo.Count, 'f', -1, 64)
				} else {
					msg += " 减持:" + strconv.FormatFloat(oldInfo.Count-latestInfo.Count, 'f', -1, 64)
				}
				msg = msg + " 排名变化:" + strconv.Itoa(oldInfo.Ranking) + "->" + strconv.Itoa(latestInfo.Ranking) +
					" 地址类型:" + alias + " 时间:" + time.Now().Format(time.RFC3339)
				if latestInfo.Addr == "1P5ZEDWTKTFGxQjZphgWPQUpe554WKDfHQ" {
					msg = msg + " 四哥出动了！！！！！！！！！！！！！"
					dingInfo(msg, true)
				} else if latestInfo.Addr == "19iqYbeATe4RxghQZJnYVFU4mjUUu76EA6" {
					msg = msg + " 短线A6哥出手！！！！！！！！！！！！！"
					dingInfo(msg, false)
				} else {
					dingInfo(msg, false)
				}
				if latestInfo.Alias == "" {
					sendHanbi = true
				}
			}
			return
		}
	}

	alias := latestInfo.Alias
	if alias == "" {
		alias = "大佬钱包"
	}

	msg := "addr:" + latestInfo.Addr + " 数量变化:" + "原未上榜未知数量" + "->" + strconv.FormatFloat(latestInfo.Count, 'f', -1, 64) +
		" 排名变化:" + "原未上榜未知排名" + "->" + strconv.Itoa(latestInfo.Ranking) +
		" 地址类型:" + alias + " 时间:" + time.Now().Format(time.RFC3339)
	dingInfo(msg, false)
	if latestInfo.Alias == "" {
		sendHanbi = true
	}
}

type ADDR_INFO struct {
	Addr      string  `json:"addr"`
	Balance   float64 `json:"balance"`
	TxCnt     int     `json:"txCnt"`
	AddrAlias string  `json:"addrAlias"`
}

type RSP_DATA struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data []ADDR_INFO `json:"data"`
}

var g_latestAddrInfo BTC_INFO
var g_OldAddrInfo BTC_INFO
var sendHanbi bool

const g_diff_num = 10

const jsonPath = "/root/info.json"

func getLatestAddrInfo() error {
	for i := 1; i <= 20; i++ {
		page, err := getPage(i)
		if err != nil {
			for i := 0; i < 10; i++ {
				time.Sleep(time.Second * 10)
				page, err = getPage(i)
				if err == nil {
					break
				}
			}
		}
		if err != nil {
			return err
		}

		var rsp RSP_DATA
		err = json.Unmarshal([]byte(page), &rsp)
		if err != nil {
			dingInfo(page, false)
			return err
		}
		if rsp.Code != 1 {
			return xerrors.New(rsp.Msg)
		}

		for j, addrInfo := range rsp.Data {
			var localAddrInfo LOCAL_BTC_ADDR_INFO

			localAddrInfo.Addr = addrInfo.Addr
			localAddrInfo.Count = addrInfo.Balance
			localAddrInfo.Alias = addrInfo.AddrAlias
			localAddrInfo.Ranking = (i-1)*10 + j + 1

			g_latestAddrInfo.LocalAddrInfos = append(g_latestAddrInfo.LocalAddrInfos, localAddrInfo)
		}
	}
	return nil
}

func writeInfo() error {
	//写入最新的数据
	filePtr, err := os.Create(jsonPath)
	if err != nil {
		return err
	}
	defer filePtr.Close()

	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(g_latestAddrInfo)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	sendHanbi = false

	//获取最新的链上排行数据
	err := getLatestAddrInfo()
	if err != nil {
		//dingInfo(err.Error() + "1", false)
		return
	}

	//获取本地数据信息 没有则写入结束
	err = loadInfo()
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			err = writeInfo()
			if err != nil {
				dingInfo(err.Error(), false)
			}
			return
		}
		dingInfo(err.Error(), false)
	}
	//log.Println(g_latestAddrInfo)

	addrTotalNum := 0.0
	total := getTotal()
	dingString := "比特币总流通量:" + strconv.FormatFloat(float64(total), 'f', -1, 64)
	for i, latestAddrInfo := range g_latestAddrInfo.LocalAddrInfos {
		//发送地址变动
		sendChange(latestAddrInfo)

		if latestAddrInfo.Alias != "" {
			continue
		}

		addrTotalNum += latestAddrInfo.Count
		if i == 9 {
			g_latestAddrInfo.Total10 = addrTotalNum
			dingString = dingString + "\n" + "钱包地址前10总量:" + strconv.FormatFloat(float64(addrTotalNum), 'f', -1, 64)
			if addrTotalNum > g_OldAddrInfo.Total10+10 {
				dingString += "\n" + "增持:" + strconv.FormatFloat(float64(addrTotalNum)-g_OldAddrInfo.Total10, 'f', -1, 64)
			} else if addrTotalNum < g_OldAddrInfo.Total10-10 {
				dingString += "\n" + "减持:" + strconv.FormatFloat(g_OldAddrInfo.Total10-float64(addrTotalNum), 'f', -1, 64)
			}
			dingString += "\n" + "占比:" + getRate(addrTotalNum, float64(total))
		} else if i == 19 {
			g_latestAddrInfo.Total20 = addrTotalNum
			dingString = dingString + "\n" + "钱包地址前20总量:" + strconv.FormatFloat(float64(addrTotalNum), 'f', -1, 64)
			if addrTotalNum > g_OldAddrInfo.Total20+10 {
				dingString += "\n" + "增持:" + strconv.FormatFloat(float64(addrTotalNum)-g_OldAddrInfo.Total20, 'f', -1, 64)
			} else if addrTotalNum < g_OldAddrInfo.Total20-10 {
				dingString += "\n" + "减持:" + strconv.FormatFloat(g_OldAddrInfo.Total20-float64(addrTotalNum), 'f', -1, 64)
			}
			dingString += "\n" + "占比:" + getRate(addrTotalNum, float64(total))
		} else if i == 49 {
			g_latestAddrInfo.Total50 = addrTotalNum
			dingString = dingString + "\n" + "钱包地址前50总量:" + strconv.FormatFloat(float64(addrTotalNum), 'f', -1, 64)
			if addrTotalNum > g_OldAddrInfo.Total50+10 {
				dingString += "\n" + "增持:" + strconv.FormatFloat(float64(addrTotalNum)-g_OldAddrInfo.Total50, 'f', -1, 64)
			} else if addrTotalNum < g_OldAddrInfo.Total50-10 {
				dingString += "\n" + "减持:" + strconv.FormatFloat(g_OldAddrInfo.Total50-float64(addrTotalNum), 'f', -1, 64)
			}
			dingString += "\n" + "占比:" + getRate(addrTotalNum, float64(total))
		} else if i == 99 {
			g_latestAddrInfo.Total100 = addrTotalNum
			dingString = dingString + "\n" + "钱包地址前100总量:" + strconv.FormatFloat(float64(addrTotalNum), 'f', -1, 64)
			if addrTotalNum > g_OldAddrInfo.Total100+10 {
				dingString += "\n" + "增持:" + strconv.FormatFloat(float64(addrTotalNum)-g_OldAddrInfo.Total100, 'f', -1, 64)
			} else if addrTotalNum < g_OldAddrInfo.Total100-10 {
				dingString += "\n" + "减持:" + strconv.FormatFloat(g_OldAddrInfo.Total100-float64(addrTotalNum), 'f', -1, 64)
			}
			dingString += "\n" + "占比:" + getRate(addrTotalNum, float64(total))
		} else if i == 199 {
			g_latestAddrInfo.Total200 = addrTotalNum
			dingString = dingString + "\n" + "钱包地址前200总量:" + strconv.FormatFloat(float64(addrTotalNum), 'f', -1, 64)
			if addrTotalNum > g_OldAddrInfo.Total200+10 {
				dingString += "\n" + "增持:" + strconv.FormatFloat(float64(addrTotalNum)-g_OldAddrInfo.Total200, 'f', -1, 64)
				if addrTotalNum > g_OldAddrInfo.Total200+1000 {
					infoString := "前200大额加仓:" + strconv.FormatFloat(float64(addrTotalNum)-g_OldAddrInfo.Total200, 'f', -1, 64)
					dingInfo(infoString, true)
				}
			} else if addrTotalNum < g_OldAddrInfo.Total200-10 {
				dingString += "\n" + "减持:" + strconv.FormatFloat(g_OldAddrInfo.Total200-float64(addrTotalNum), 'f', -1, 64)
				if addrTotalNum < g_OldAddrInfo.Total200-1000 {
					infoString := "前200大额减仓:" + strconv.FormatFloat(g_OldAddrInfo.Total200-float64(addrTotalNum), 'f', -1, 64)
					dingInfo(infoString, true)
				}
			}
			dingString += "\n" + "占比:" + getRate(addrTotalNum, float64(total))
		}
		//fmt.Println(result)
	}

	if sendHanbi {
		bAll := false
		richData, err := getFullRichData()
		if err == nil {
			richData.Position = richData.Position[len(richData.Position)-12:]
			beforeStatus := getStatus(richData)

			richData.Position = append(richData.Position, g_latestAddrInfo.Total200)
			richData.Position = richData.Position[len(richData.Position)-12:]
			nowStatus := getStatus(richData)
			dingString += "\n当前状态:" + nowStatus
			if nowStatus != beforeStatus {
				dingString += "\n 状态从 " + beforeStatus + " 变成 " + nowStatus
				bAll = true
			}
		}
		if bAll {
			dingInfo(dingString, true)
		} else {
			dingInfo(dingString, false)
		}
	}

	err = writeInfo()
	if err != nil {
		dingInfo(err.Error(), false)
	}

	err = writeStatusInfo()
	if err != nil {
		dingInfo(err.Error(), false)
	}
	return
}
