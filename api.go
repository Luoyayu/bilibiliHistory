package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func getHistory(max, accessKey string) (data *retStruct) {
	params := map[string]string{
		"access_key": accessKey,
		"appkey":     APPKEY,
		"build":      "5500300",
		"business":   "all",
		"channel":    "bilibili140",
		"max":        max,
		"max_tp":     "3",
		"mobi_app":   "android",
		"platform":   "android",
		"ps":         "20",
		"ts":         fmt.Sprint(time.Now().Unix()),
	}
	params["sign"] = getSign(params)

	l := url.Values{}
	c := &http.Client{}
	for k, v := range params {
		l.Add(k, v)
	}
	r, _ := http.NewRequest("GET",
		"https://app.biliapi.net/x/v2/history/cursor"+"?"+l.Encode(),
		nil)

	if r != nil {
		r.Header.Add("user-agent", "Mozilla/5.0 BiliDroid/5.37.0")
	}

	if resp, err := c.Do(r); err == nil {
		b, _ := ioutil.ReadAll(resp.Body)
		var ret retStruct
		_ = json.Unmarshal(b, &ret)

		return &ret
	}
	return nil
}
