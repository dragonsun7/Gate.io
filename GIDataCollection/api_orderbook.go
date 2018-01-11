/*
	当前市场深度
	返回当前市场深度（委托挂单），其中 asks 是委卖单, bids 是委买单。
	请替换 [CURR_A] and [CURR_B] 为您需要查看的币种.
	http://data.gate.io/api2/1/orderBook/[CURR_A]_[CURR_B]
*/

package main

import (
	"github.com/buger/jsonparser"
	"strconv"
	"errors"
	"encoding/json"
)

type ApiOrderBook struct {
	Api
	Pair string
	Result bool
	asks []ApiDepth
	bids []ApiDepth
}

type ApiDepth struct {
	Amount 	float64
	Price 	float64
}

func (api *ApiOrderBook) Init(pg *Postgres, pair string) (*ApiOrderBook) {
	api.desc = "当前市场深度" + pair
	api.uri = "orderBook/" + pair
	api.Pair = pair
	api.pg = pg

	return api
}

func (api *ApiOrderBook) Request() ([]byte, error) {
	return api.httpGet(api.uri)
}

func (api *ApiOrderBook) Parser(body []byte) (error) {
	result, err := jsonparser.GetString(body, "result")
	if err != nil {
		return err
	}

	api.Result, err = strconv.ParseBool(result)
	if err != nil {
		return err
	}

	if !api.Result {
		return errors.New("接口返回失败")
	}

	var asks []ApiDepth
	_, err =jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var values []float64
		json.Unmarshal(value, &values)

		var ask ApiDepth
		ask.Amount = values[0]
		ask.Price = values[1]

		asks = append(asks, ask)
	}, "asks")
	if err != nil {
		return err
	}

	var bids []ApiDepth
	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var values []float64
		json.Unmarshal(value, &values)

		var bid ApiDepth
		bid.Amount = values[0]
		bid.Price = values[1]

		bids = append(bids, bid)
	}, "bids")
	if err != nil {
		return err
	}

	api.asks = asks
	api.bids = bids
	return nil
}

func (api *ApiOrderBook) Save() (error) {
	return nil
}
