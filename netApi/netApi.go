package netApi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/xerrors"
)

type ADDR_INFO struct {
	Addr      string  `json:"addr"`
	Balance   float64 `json:"balance"`
	TxCnt     int64   `json:"txCnt"`
	AddrAlias string  `json:"addrAlias"`
}

type RSP_DATA struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data []ADDR_INFO `json:"data"`
}

type LOCAL_BTC_ADDR_INFO struct {
	Addr  string  //钱包地址
	Count float64 //数量
	Alias string  //别名
}

func GetLatestAddrInfo() ([]LOCAL_BTC_ADDR_INFO, error) {
	var rspBtcAddrInfo []LOCAL_BTC_ADDR_INFO
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
			return rspBtcAddrInfo, err
		}

		var rsp RSP_DATA
		err = json.Unmarshal([]byte(page), &rsp)
		if err != nil {
			return rspBtcAddrInfo, err
		}
		if rsp.Code != 1 {
			return rspBtcAddrInfo, xerrors.New(rsp.Msg)
		}

		for _, addrInfo := range rsp.Data {
			var localAddrInfo LOCAL_BTC_ADDR_INFO

			localAddrInfo.Addr = addrInfo.Addr
			localAddrInfo.Count = addrInfo.Balance
			localAddrInfo.Alias = addrInfo.AddrAlias

			rspBtcAddrInfo = append(rspBtcAddrInfo, localAddrInfo)
		}
	}
	return rspBtcAddrInfo, nil
}

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
