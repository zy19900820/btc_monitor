package ding

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/xerrors"
)

const DING_ERROR_URL = "https://oapi.dingtalk.com/robot/send?access_token=d5e507b4ef9122b60af42e7418050a67eb43d797315a2504469de01570684ba9"
const ERROR_WORD = "提醒 "

func DingInfo(msg string, bAll bool) error {
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
