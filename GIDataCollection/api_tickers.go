/*
	所有交易行情
	返回系统支持的所有交易对的 最新，最高，最低 交易行情和交易量，每10秒钟更新:
	http://data.gate.io/api2/1/tickers
*/

package main

import "github.com/buger/jsonparser"

type ApiTickers struct {
	Api
	Tickers []ApiTicker
}

func (api *ApiTickers) Init(pg *Postgres) (* ApiTickers) {
	api.desc = "所有交易行情"
	api.uri = "tickers"
	api.pg = pg
	return api
}

func (api *ApiTickers) Request() ([]byte, error) {
	return api.httpGet(api.uri)
}

func (api *ApiTickers) Parser(body []byte) (error) {
	var tickers []ApiTicker

	err := jsonparser.ObjectEach(body, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		ticker := new(ApiTicker).Init(api.pg, string(key))

		err1 := ticker.Parser(value)
		if err1 != nil {
			return err1
		}

		tickers = append(tickers, *ticker)
		return nil
	})
	if err != nil {
		return err
	}

	api.Tickers = tickers
	return nil
}

func (api *ApiTickers) Save() (error) {
	for _, ticker := range api.Tickers  {
		err := ticker.Save()
		if err != nil {
			return err
		}
	}

	return nil
}
