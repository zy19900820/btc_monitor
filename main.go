package main

import (
	conf "btc_monitor/config"
	"btc_monitor/ding"
	"btc_monitor/netApi"
	"btc_monitor/sql"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
var gBlockHeight int64
var gAddressAlias map[string]string

func main() {
	gAddressAlias = make(map[string]string)
	gAddressAlias["34xp4vRoCGJym3xR7yCVPFHoCNxv4Twseo"] = "Binance-coldwallet"
	gAddressAlias["3M219KR5vEneNb47ewrPfWyb5jQ2DjxRP6"] = "Binance-coldwallet"
	gAddressAlias["3Kzh9qAqVWQhEsfQz7zEQL1EuSx5tyNLNS"] = "Gemini"
	gAddressAlias["bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h"] = "Binance-wallet"
	gAddressAlias["385cR5DM96n1HvBDMzLHPYcw89fZAXULJP"] = "Bittrex-coldwallet"

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

		for i := 0; i < len(netInfos)/40; i++ {
			var adds []string

			for j := 0; j < 40; j++ {
				adds = append(adds, netInfos[i*40+j].Addr)
			}

			result, err := netApi.GetMulAddressInfo(adds)
			if err != nil {
				log.Println(err)
				return
			}
			err = sql.UpdateNetInfos(result)
			if err != nil {
				log.Println(err)
				return
			}
		}
	} else {
		//从数据库获取前200大户地址
		addressInfos, err := sql.GetSqlTop200Addr()
		if err != nil {
			log.Println(err)
			return
		}

		//更新地址余额
		for i := 0; i < len(addressInfos)/40; i++ {
			var adds []string

			for j := 0; j < 40; j++ {
				adds = append(adds, addressInfos[i*40+j])
			}

			result, err := netApi.GetMulAddressInfo(adds)
			if err != nil {
				log.Println(err)
				return
			}
			err = sql.UpdateNetInfos(result)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(Cors())

	router.GET("/hello", HandleHello)

	go listenLatestBlock()
	go updateDaily()
	log.Println("ready start server")
	router.Run(":" + strconv.Itoa(cfg.Port))
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		//startTime := time.Now().Format("2006-01-02 15:04:05")
		c.Next()

		//responseBody := bodyLogWriter.body.String()

		//endTime := time.Now().Format("2006-01-02 15:04:05")

		// if c.Request.Method == "POST" {
		// 	c.Request.ParseForm()
		// }
	}
}

type RESPONSE struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HandleHello(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}
	context.JSON(http.StatusOK, response)
}

func listenLatestBlock() {
	for {
		height, err := netApi.GetLatestBlock()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 300)
			continue
		}
		//log.Println(height)
		if gBlockHeight == height {
			time.Sleep(time.Second * 300)
			continue
		}
		gBlockHeight = height
		heightBlockInfo, err := netApi.GetHeightBlockInfo(height)
		//heightBlockInfo, err := netApi.GetHeightBlockInfo(741261)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 300)
			continue
		}
		//解析获得转账数量>10地址
		myTxInfos := webInfoToAddress(heightBlockInfo)

		var needRecordAddress []string
		for _, myTxInfo := range myTxInfos {
			for _, input := range myTxInfo.Inputs {
				needRecordAddress = append(needRecordAddress, input.Address)
			}
			for _, out := range myTxInfo.Outputs {
				needRecordAddress = append(needRecordAddress, out.Address)
			}
		}
		//log.Println("needrecordAddress:", needRecordAddress)
		//更新数据库
		var walletInfoResults []netApi.LOCAL_BTC_ADDR_INFO
		for i := 0; i < len(needRecordAddress)/40; i++ {
			var adds []string

			for j := 0; j < 40; j++ {
				adds = append(adds, needRecordAddress[i*40+j])
			}

			result, err := netApi.GetMulAddressInfo(adds)
			if err != nil {
				log.Println(err)
				return
			}
			walletInfoResults = append(walletInfoResults, result...)
			//地址余额大于8000的更新
			// var more8000AddressInfo []netApi.LOCAL_BTC_ADDR_INFO
			// for _, res := range result {
			// 	if res.Count > 8000 {
			// 		more8000AddressInfo = append(more8000AddressInfo, res)
			// 	}
			// }
			// log.Println(more8000AddressInfo)
			err = sql.UpdateNetInfos(result)
			if err != nil {
				log.Println(err)
				return
			}
		}
		for _, myTxInfo := range myTxInfos {
			//1.没有大额转账不报
			haveBigTransfer := false
			for _, input := range myTxInfo.Inputs {
				if input.Value > 500 {
					haveBigTransfer = true
				}
			}
			if !haveBigTransfer {
				continue
			}
			//2.大额转账转入自己不报
			transferOwn := false
			for _, input := range myTxInfo.Inputs {
				if input.Value > 500 {
					for _, output := range myTxInfo.Outputs {
						if output.Address == input.Address {
							if (input.Value - output.Value) < 100 {
								transferOwn = true
								break
							}
						}
					}
				}
			}
			if transferOwn {
				continue
			}
			//3.更换钱包不报 转出后余额为0 转入后为第一笔交易
			changeWallet := false
			for _, input := range myTxInfo.Inputs {
				if input.Value > 500 {
					haveBigTransfer = true
				}

				inBalance := getNowBalance(input.Address, walletInfoResults)
				if inBalance < 1 {
					//转出后自己钱包币很少 判断转入钱包原先数量是否很少
					for _, output := range myTxInfo.Outputs {
						if output.Value > 500 {
							outBalance := getNowBalance(output.Address, walletInfoResults)
							if outBalance-output.Value < 1 {
								changeWallet = true
								break
							}
						}
					}
					if changeWallet {
						break
					}
				}
			}
			if changeWallet {
				continue
			}
			//4.不是大户变动不报 钱包余额大于8000
			bBigUser := false
			for _, input := range myTxInfo.Inputs {
				if input.Value > 500 {
					nowBalance := getNowBalance(input.Address, walletInfoResults)
					if nowBalance > 8000 {
						bBigUser = true
						break
					}
				}
			}
			for _, output := range myTxInfo.Outputs {
				if output.Value > 500 {
					nowBalance := getNowBalance(output.Address, walletInfoResults)
					if nowBalance > 8000 {
						bBigUser = true
						break
					}
				}
			}
			if !bBigUser {
				continue
			}
			//5.剩下把这笔交易汇报 input output在10以内
			if len(myTxInfo.Inputs) > 10 || len(myTxInfo.Outputs) > 10 {
				continue
			}
			//这里要ding了
			var dingString string
			for _, input := range myTxInfo.Inputs {
				dingString += "地址:" + input.Address
				alias, ok := gAddressAlias[input.Address]
				if ok {
					dingString += " alias:" + alias
				}
				dingString += " 转出数量:" + strconv.FormatFloat(input.Value, 'f', -1, 64)
				dingString += " 转出后余额:" + strconv.FormatFloat(getNowBalance(input.Address, walletInfoResults), 'f', -1, 64)
			}
			ding.DingInfo(dingString, true)

			dingString = ""
			for _, output := range myTxInfo.Outputs {
				dingString += "地址:" + output.Address
				alias, ok := gAddressAlias[output.Address]
				if ok {
					dingString += " alias:" + alias
				}
				dingString += " 转入数量:" + strconv.FormatFloat(output.Value, 'f', -1, 64)
				dingString += " 转入后余额:" + strconv.FormatFloat(getNowBalance(output.Address, walletInfoResults), 'f', -1, 64)
			}
			ding.DingInfo(dingString, true)
		}
		time.Sleep(time.Second * 300)
	}
}

func getNowBalance(address string, walletInfoResults []netApi.LOCAL_BTC_ADDR_INFO) float64 {
	var nowBalance float64
	for _, walletInfoResult := range walletInfoResults {
		if walletInfoResult.Addr == address {
			nowBalance = walletInfoResult.Count
			break
		}
	}
	return nowBalance
}

type PREV_OUT struct {
	Tx_index           int64               `json:"tx_index"`
	Value              int64               `json:"value"`
	N                  int64               `json:"n"`
	Type               int64               `json:"type"`
	Spent              bool                `json:"spent"`
	Script             string              `json:"script"`
	Spending_outpoints []SPENDING_OUTPOINT `json:"spending_outpoints"`
	Addr               string              `json:"addr"`
}

type SPENDING_OUTPOINT struct {
	Tx_index int64 `json:"tx_index"`
	N        int64 `json:"n"`
}

type INPUT struct {
	Sequence int64    `json:"sequence"`
	Witness  string   `json:"witness"`
	Script   string   `json:"script"`
	Index    int64    `json:"index"`
	Prev_out PREV_OUT `json:"prev_out"`
}

type OUTPUT struct {
	Type               int64               `json:"type"`
	Spent              bool                `json:"spent"`
	Value              int64               `json:"value"`
	Spending_outpoints []SPENDING_OUTPOINT `json:"spending_outpoints"` //好像不重要
	N                  int64               `json:"n"`
	Tx_index           int64               `json:"tx_index"`
	Script             string              `json:"script"`
	Addr               string              `json:"addr"`
}

type TX_INFO struct {
	Hash         string   `json:"hash"`
	Ver          int64    `json:"ver"`
	Vin_sz       int64    `json:"vin_sz"`
	Vout_sz      int64    `json:"vout_sz"`
	Size         int64    `json:"size"`
	Weight       int64    `json:"weight"`
	Fee          int64    `json:"fee"`
	Relayed_by   string   `json:"relayed_by"`
	Lock_time    int64    `json:"lock_time"`
	Tx_index     int64    `json:"tx_index"`
	Double_spend bool     `json:"double_spend"` //双花是否要判断
	Time         int64    `json:"time"`
	Block_index  int64    `json:"block_index"`
	Block_height int64    `json:"block_height"`
	Inputs       []INPUT  `json:"inputs"`
	Out          []OUTPUT `json:"out"`
}

type BLOCK_INFO struct {
	Hash        string    `json:"hash"`
	Ver         int64     `json:"ver"`
	Prev_block  string    `json:"prev_block"`
	Mrkl_root   string    `json:"mrkl_root"`
	Time        int64     `json:"time"`
	Bits        int64     `json:"bits"`
	Next_block  []string  `json:"next_block"`
	Fee         int64     `json:"fee"`
	Nonce       int64     `json:"nonce"`
	N_tx        int64     `json:"n_tx"`        //交易数量？
	Size        int64     `json:"size"`        //快大小？
	Block_index int64     `json:"block_index"` //快高度?
	Main_chain  bool      `json:"main_chain"`
	Height      int64     `json:"height"` //快高度？
	Weight      int64     `json:"weight"`
	Tx          []TX_INFO `json:"tx"`
}

type IN_OUT_PUT struct {
	Address string  `json:"string"`
	Value   float64 `json:"value"`
}

type MY_TX_INFO struct {
	Inputs  []IN_OUT_PUT `json:"inputs"`
	Outputs []IN_OUT_PUT `json:"outputs"`
}

func webInfoToAddress(webInfo string) []MY_TX_INFO {
	var parseResults []MY_TX_INFO
	//log.Println(webInfo)
	map1 := make(map[string][]BLOCK_INFO)
	err := json.Unmarshal([]byte(webInfo), &map1)
	if err != nil {
		log.Println(err)
		return parseResults
	}
	//log.Println(map1)
	for _, v := range map1 {
		//第一层是k：block  v：[map]
		for _, header := range v {
			//741261 高度解析出来也只有一个大的map 不确定其他区块是不是不一样
			//应该是的 看内容是快头信息
			for _, tx := range header.Tx {
				//这里是交易信息
				var parseResult MY_TX_INFO
				for _, input := range tx.Inputs {
					if input.Prev_out.Value > 10*100000000 {
						parseResult.Inputs = append(parseResult.Inputs, IN_OUT_PUT{Address: input.Prev_out.Addr, Value: float64(input.Prev_out.Value) / 100000000})
					}
				}

				for _, out := range tx.Out {
					if out.Value > 10*100000000 {
						parseResult.Outputs = append(parseResult.Outputs, IN_OUT_PUT{Address: out.Addr, Value: float64(out.Value) / 100000000})
					}
				}
				if len(parseResult.Inputs) != 0 {
					parseResults = append(parseResults, parseResult)
				}
			}
		}
	}
	return parseResults
}

func updateDaily() {
	for {
		err := sql.UpdateDaily(20)
		if err != nil {
			log.Println(err)
		}
		err = sql.UpdateDaily(50)
		if err != nil {
			log.Println(err)
		}
		err = sql.UpdateDaily(100)
		if err != nil {
			log.Println(err)
		}
		err = sql.UpdateDaily(180)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Hour * 24)
	}
}
