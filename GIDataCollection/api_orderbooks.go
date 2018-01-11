/*
	市场深度
	返回系统支持的所有交易对的市场深度（委托挂单），其中 asks 是委卖单, bids 是委买单。
	http://data.gate.io/api2/1/orderBooks
*/

package main

import "github.com/buger/jsonparser"

type ApiOrderBooks struct {
	Api
	OrderBooks []ApiOrderBook
}

func (api *ApiOrderBooks) Init(pg *Postgres) (*ApiOrderBooks) {
	api.desc = "市场深度"
	api.uri = "orderBooks"
	api.pg = pg

	return api
}

func (api *ApiOrderBooks) Request() ([]byte, error) {
	return api.httpGet(api.uri)
}

func (api *ApiOrderBooks) Parser(body []byte) (error) {
	var orderBooks []ApiOrderBook

	err := jsonparser.ObjectEach(body, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		orderBook := new(ApiOrderBook).Init(api.pg, string(key))

		err1 := orderBook.Parser(value)
		if err1 != nil {
			return err1
		}

		orderBooks = append(orderBooks, *orderBook)
		return nil
	})
	if err != nil {
		return err
	}

	api.OrderBooks = orderBooks
	return nil
}

func (api *ApiOrderBooks) Save() (error) {
	return nil
}