package main

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
)

func getSign(params map[string]string) string {
	q := url.Values{}
	for k, v := range params {
		q.Add(k, v) // sorted
	}
	m := md5.New()
	m.Write([]byte(q.Encode() + APPSECRET))
	return hex.EncodeToString(m.Sum(nil))
}
