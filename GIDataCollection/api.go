package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"crypto/hmac"
	"crypto/sha512"
	"strings"
	"time"
)

/* ---------------------------------------- */

type IApi interface {
	GetDesc() (string)
	Request() ([]byte, error)
	Parser(body []byte) (error)
	Save() (error)
}

func ApiDo(iapi IApi) (err error) {
	start := time.Now()
	fmt.Print(fmt.Sprintf("    开始处理接口【%s】...", iapi.GetDesc()))
	defer func() {
		status := "成功"
		if err != nil {
			status = "失败"
		}
		fmt.Println(status, "! 耗时:", time.Since(start))
	}()

	var body []byte
	body, err = iapi.Request()
	if err != nil {
		return
	}

	err = iapi.Parser(body)
	if err != nil {
		return
	}

	return iapi.Save()
}

/* ---------------------------------------- */

type Api struct {
	desc string
	uri string
	pg *Postgres
}

func (api *Api) GetDesc() (string) {
	return api.desc
}

func (api *Api) jointUrl(uri string) (string) {
	return APIPrefix + uri
}

func (api *Api) getSign(params string) string {
	key := []byte(APISecret)
	mac := hmac.New(sha512.New, key)
	mac.Write([]byte(params))

	return fmt.Sprintf("%x", mac.Sum(nil))
}

func (api *Api) httpGet(uri string) ([]byte, error) {
	resp, err := http.Get(api.jointUrl(uri))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (api *Api) httpPost(uri, params string) ([]byte, error) {
	resp, err := http.Post(api.jointUrl(uri), "application/x-www-form-urlencoded", strings.NewReader(params))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (api *Api) clientDo(method, uri, params string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, api.jointUrl(uri), strings.NewReader(params))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("key", APIKey)
	req.Header.Set("sign", api.getSign(params))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
