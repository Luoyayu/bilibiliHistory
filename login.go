package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type List struct {
	Title            string   `json:"title"`
	Covers           []string `json:"covers"`
	Uri              string   `json:"uri"`
	Name             string   `json:"name"`
	Mid              int      `json:"mid"`
	Goto             string   `json:"goto"`
	Badge            string   `json:"badge"`
	ViewAt           int64    `json:"view_at"`
	DisplayAttention int      `json:"display_attention"`
	Duration         int      `json:"duration"`
}

type Cursor struct {
	Max   int `json:"max"`
	MaxTp int `json:"max_tp"`
	Ps    int `json:"ps"`
}
type dataStruct struct {
	Hash         string `json:"hash"`
	Key          string `json:"key"`
	AccessToken  string `json:"access_token"`
	Mid          int    `json:"mid"`
	RefreshToken string `json:"refresh_token"`
	List         []List `json:"list"`
	Cursor       Cursor `json:"cursor"`
}

type retStruct struct {
	Message    string `json:"message"`
	Ts         int64    `json:"ts"`
	Code       int    `json:"code"`
	dataStruct `json:"data"`
}

func getKey() (flag bool, hash, key []byte) {
	flag = false
	params := map[string]string{
		"appkey": APPKEY,
	}
	params["sign"] = getSign(params)
	l := url.Values{}
	c := &http.Client{}
	for k, v := range params {
		l.Add(k, v)
	}

	r, _ := http.NewRequest("POST",
		strings.Join(
			[]string{LOGINURL, "api/oauth2/getKey"}, "/")+"?"+l.Encode(),
		nil)

	if r != nil {
		r.Header.Add("user-agent", "Mozilla/5.0 BiliDroid/5.37.0")
	}

	resp, err := c.Do(r)
	if err != nil {
		log.Fatal(err)
	} else {
		var ret retStruct
		b, _ := ioutil.ReadAll(resp.Body)
		if DEBUG {
			log.Println(string(b))
		}
		defer resp.Body.Close()

		_ = json.Unmarshal(b, &ret)
		if ret.Code != 0 {
			log.Fatal("get Key failed! Error message: ", ret.Message)
		} else {
			flag = true
			if DEBUG {
				log.Println(ret.Hash)
				log.Println(ret.Key)
			}

			hash = []byte(ret.Hash)
			key = []byte(ret.Key)
		}
	}
	return
}

// return AccessToken
func doLogin(username, password string) string {
	ok, r1, r2 := getKey()
	if ok {
		resp := login(username, rsaEncrypt(password, r1, r2))
		if DEBUG {
			log.Println(resp)
		}
		if resp.Code == 0 {
			return resp.AccessToken
		}
	}
	return ""
}

func login(username, password string) (ret retStruct) {
	params := map[string]string{
		"appkey":   APPKEY,
		"password": password,
		"username": username,
	}
	params["sign"] = getSign(params)

	l := url.Values{}
	c := &http.Client{}
	for k, v := range params {
		l.Add(k, v)
	}

	r, _ := http.NewRequest("POST",
		strings.Join(
			[]string{LOGINURL, "api/oauth2/login"}, "/")+"?"+l.Encode(),
		nil)

	if r != nil {
		r.Header.Add("user-agent", "Mozilla/5.0 BiliDroid/5.37.0")
	}
	resp, err := c.Do(r)
	if err != nil {
		log.Fatal(err)
	} else {
		b, _ := ioutil.ReadAll(resp.Body)
		if DEBUG {
			log.Println("登录返回值")
			log.Println(string(b))
		}
		_ = json.Unmarshal(b, &ret)
	}
	return
}

func rsaEncrypt(password string, hash, key []byte) string {
	block, _ := pem.Decode(key)
	if block == nil {
		log.Fatal("private key error!")
	}
	pubKey, _ := x509.ParsePKIXPublicKey(block.Bytes)
	em, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), []byte(string(hash)+password))
	if err != nil {
		log.Fatal(err)
	} else {
		if DEBUG {
			log.Println("RSA BASE64 加密结果: ")
			log.Println(base64.StdEncoding.EncodeToString(em))
		}
		return base64.StdEncoding.EncodeToString(em)
	}
	return ""
}
