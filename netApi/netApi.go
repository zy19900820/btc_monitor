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

type NET_RSP_ADDRESS_INFO struct {
	Count int64 `json:"final_balance"`
}

func parseMulAddressInfo(str string) ([]LOCAL_BTC_ADDR_INFO, error) {
	//str := `{"1P5ZEDWTKTFGxQjZphgWPQUpe554WKDfHQ":{"final_balance":12937154385778,"n_tx":753,"total_received":22903342834381},"34xp4vRoCGJym3xR7yCVPFHoCNxv4Twseo":{"final_balance":25259723636577,"n_tx":757,"total_received":119037494563185},"3LYJfcfHPXYJreMsASk2jkn69LWEYKzexb":{"final_balance":12535113665929,"n_tx":45,"total_received":12535113665929},"3Kzh9qAqVWQhEsfQz7zEQL1EuSx5tyNLNS":{"final_balance":6548821843046,"n_tx":1473,"total_received":113757489644930},"bc1qgdjqv0av3q56jvd82tkdjpy7gdp9ut8tlqmgrpmv24sq90ecnvqqjwvw97":{"final_balance":16800998579715,"n_tx":96,"total_received":356165894996868},"bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h":{"final_balance":13980561651354,"n_tx":348890,"total_received":2913702325831397},"1FeexV6bAHb8ybZjqQMjJrcCrHGW9sb6uF":{"final_balance":7995721806726,"n_tx":420,"total_received":7995721806726},"bc1qazcm763858nkj2dj986etajv6wquslv8uxwczt":{"final_balance":9464330459658,"n_tx":43,"total_received":9464330459658},"bc1qa5wkgaew2dkv56kfvj49j0av5nml45x9ek9hz6":{"final_balance":6937017661874,"n_tx":47,"total_received":6937017661874},"37XuVSEpWW4trkfmvWzegTHQt7BdktSKUs":{"final_balance":9450534198112,"n_tx":135,"total_received":9450606038446}}`
	map1 := make(map[string]NET_RSP_ADDRESS_INFO)
	err := json.Unmarshal([]byte(str), &map1)
	if err != nil {
		return []LOCAL_BTC_ADDR_INFO{}, err
	}

	var rspInfos []LOCAL_BTC_ADDR_INFO
	for k, v := range map1 {
		tmpInfo := LOCAL_BTC_ADDR_INFO{Addr: k, Count: float64(v.Count) / 100000000}
		rspInfos = append(rspInfos, tmpInfo)
	}
	return rspInfos, nil
}

func GetMulAddressInfo(adds []string) ([]LOCAL_BTC_ADDR_INFO, error) {
	if len(adds) > 40 {
		return []LOCAL_BTC_ADDR_INFO{}, xerrors.New("adds too long")
	}
	var addUrls string
	for _, address := range adds {
		addUrls += address + "|"
	}
	addUrls = addUrls[:len(addUrls)-1]
	url := "https://blockchain.info/balance?active=" + addUrls
	res, err := http.Get(url)
	if err != nil {
		return []LOCAL_BTC_ADDR_INFO{}, err
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return []LOCAL_BTC_ADDR_INFO{}, err
	}
	if res.StatusCode != 200 {
		return []LOCAL_BTC_ADDR_INFO{}, xerrors.New("http error")
	}
	return parseMulAddressInfo(string(robots))
}

type LATEST_HEIGHT struct {
	Height int64 `json:"height"`
}

func GetLatestBlock() (int64, error) {
	url := "https://blockchain.info/latestblock"
	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return 0, err
	}
	if res.StatusCode != 200 {
		return 0, xerrors.New("http error")
	}

	var latestHeight LATEST_HEIGHT
	err = json.Unmarshal(robots, &latestHeight)
	if err != nil {
		return 0, err
	}
	return latestHeight.Height, nil
}

func GetHeightBlockInfo(height int64) (string, error) {
	url := "https://blockchain.info/block-height/" + strconv.FormatInt(height, 10)
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
