package main

import (
	conf "btc_monitor/config"
	"btc_monitor/netApi"
	"btc_monitor/sql"
	"log"
)

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
// func getCoin(addr string) int {
// 	res, err := http.Get("https://btc.com/" + addr)
// 	if err != nil {
// 		return 0
// 	}
// 	robots, err := ioutil.ReadAll(res.Body)
// 	res.Body.Close()
// 	if err != nil {
// 		return 0
// 	}

// 	page := strings.Replace(string(robots), "\n", " ", -1)
// 	index := strings.Index(page, ">Balance<")

// 	indexNextFirst := strings.Index(page[index:], "<dd>")
// 	indexNext := strings.Index(page[index+indexNextFirst:], "<span")

// 	coinNumStr := page[index+indexNextFirst+4 : index+indexNextFirst+indexNext]
// 	coinNumStr = strings.Replace(coinNumStr, " ", "", -1)
// 	coinNumStr = strings.Replace(coinNumStr, ",", "", -1)

// 	indexDot := strings.Index(coinNumStr, ".")
// 	coinNumStr = coinNumStr[:indexDot]
// 	coinNum, _ := strconv.Atoi(coinNumStr)
// 	return coinNum
// }

// func getAddr(result string) string {
// 	//fmt.Println(result)
// 	index1 := strings.Index(result, "btc.com/")
// 	index2 := strings.Index(result[index1:], "\"")
// 	addrStr := result[index1+8 : index1+index2]
// 	addrStr = strings.Replace(addrStr, " ", "", -1)
// 	return addrStr
// }

// func getRate(num, total float64) string {
// 	rate := float64(num) / float64(total) * 100
// 	rateStr := strconv.FormatFloat(rate, 'f', -1, 64)
// 	return rateStr + "%"
// }

// func sendChange(latestInfo LOCAL_BTC_ADDR_INFO) {
// 	for _, oldInfo := range g_OldAddrInfo.LocalAddrInfos {
// 		//币安热钱包 跳过
// 		if latestInfo.Addr == "1NDyJtNTjmwk5xPNhjgAMu4HDHigtobu1s" {
// 			return
// 		}

// 		if oldInfo.Addr == latestInfo.Addr {
// 			if (latestInfo.Count-oldInfo.Count) > 10 || (latestInfo.Count-oldInfo.Count) < -10 {
// 				alias := latestInfo.Alias
// 				if alias == "" {
// 					alias = "大佬钱包"
// 				}

// 				msg := "addr:" + latestInfo.Addr + " 数量变化:" + strconv.FormatFloat(oldInfo.Count, 'f', -1, 64) + "->" + strconv.FormatFloat(latestInfo.Count, 'f', -1, 64)
// 				if (latestInfo.Count - oldInfo.Count) > 10 {
// 					msg += " 增持:" + strconv.FormatFloat(latestInfo.Count-oldInfo.Count, 'f', -1, 64)
// 				} else {
// 					msg += " 减持:" + strconv.FormatFloat(oldInfo.Count-latestInfo.Count, 'f', -1, 64)
// 				}
// 				msg = msg + " 排名变化:" + strconv.Itoa(oldInfo.Ranking) + "->" + strconv.Itoa(latestInfo.Ranking) +
// 					" 地址类型:" + alias + " 时间:" + time.Now().Format(time.RFC3339)
// 				if latestInfo.Addr == "1P5ZEDWTKTFGxQjZphgWPQUpe554WKDfHQ" {
// 					msg = msg + " 四哥出动了！！！！！！！！！！！！！"
// 					dingInfo(msg, true)
// 				} else if latestInfo.Addr == "19iqYbeATe4RxghQZJnYVFU4mjUUu76EA6" {
// 					msg = msg + " 短线A6哥出手！！！！！！！！！！！！！"
// 					dingInfo(msg, false)
// 				} else {
// 					dingInfo(msg, false)
// 				}
// 				if latestInfo.Alias == "" {
// 					sendHanbi = true
// 				}
// 			}
// 			return
// 		}
// 	}

// 	alias := latestInfo.Alias
// 	if alias == "" {
// 		alias = "大佬钱包"
// 	}

// 	msg := "addr:" + latestInfo.Addr + " 数量变化:" + "原未上榜未知数量" + "->" + strconv.FormatFloat(latestInfo.Count, 'f', -1, 64) +
// 		" 排名变化:" + "原未上榜未知排名" + "->" + strconv.Itoa(latestInfo.Ranking) +
// 		" 地址类型:" + alias + " 时间:" + time.Now().Format(time.RFC3339)
// 	dingInfo(msg, false)
// 	if latestInfo.Alias == "" {
// 		sendHanbi = true
// 	}
// }

var gCfg conf.ServerConf

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//get config
	cfg, err := conf.GetConfig()
	if err != nil {
		log.Println(err)
		return
	}

	gCfg = cfg
	log.Println("config load over:", cfg)

	err = sql.InitMysqlDB(cfg.Mysql.Name, cfg.Mysql.Password, cfg.Mysql.Db_ip, int(cfg.Mysql.Db_port), cfg.Mysql.Db)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("init db over")

	bInitNet, err := sql.CheckNetInit()
	if err != nil {
		log.Println(err)
		return
	}

	if bInitNet {
		//从网络获取
		netInfos, err := netApi.GetLatestAddrInfo()
		if err != nil {
			log.Println(err)
			return
		}

		err = sql.InsertNetInfos(netInfos)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
