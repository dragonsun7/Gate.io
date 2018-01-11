/*
	单项交易行情
	返回最新，最高，最低 交易行情和交易量，每10秒钟更新:
	替换 [CURR_A] and [CURR_B] 为您需要查看的币种.
	http://data.gate.io/api2/1/ticker/[CURR_A]_[CURR_B]
*/

package main

import (
	"github.com/buger/jsonparser"
	"strconv"
	"errors"
)

type ApiTicker struct {
	Api
	Pair		string
	Result 		bool		`json:"result,string"`
	BaseVol		float64		`json:"baseVolume"`		// 交易量
	QuoteVol	float64		`json:"quoteVolume"`	// 兑换货币交易量
	High24		float64		`json:"high24hr"`		// 24小时最高价
	Low24		float64		`json:"low24hr"`		// 24小时最低价
	HighestBid	float64		`json:"highestBid"`		// 买方最高价
	LowestAsk	float64		`json:"lowestAsk"`		// 卖方最低价
	Last		float64		`json:"last"`			// 最新成交价
	Percent		float64		`json:"percentChange"`	// 涨跌百分比
}

func (api *ApiTicker) Init(pg *Postgres, pair string) (*ApiTicker) {
	api.desc = "单项交易行情" + pair
	api.uri = "ticker/" + pair
	api.pg = pg
	api.Pair = pair
	return api
}

func (api *ApiTicker) Request() ([]byte, error) {
	return api.httpGet(api.uri)
}

func (api *ApiTicker) Parser(body []byte) (error) {
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

	api.BaseVol, _ = jsonparser.GetFloat(body, "baseVolume")
	api.QuoteVol, _ = jsonparser.GetFloat(body, "quoteVolume")
	api.High24, _ = jsonparser.GetFloat(body, "high24hr")
	api.Low24, _ = jsonparser.GetFloat(body, "low24hr")
	api.HighestBid, _ = jsonparser.GetFloat(body, "highestBid")
	api.LowestAsk, _ = jsonparser.GetFloat(body, "lowestAsk")
	api.Last, _ = jsonparser.GetFloat(body, "last")
	api.Percent, _ = jsonparser.GetFloat(body, "percentChange")

	return nil
}

func (api *ApiTicker) Save() (error) {
	sql := `
UPDATE
	tr_market
SET
	base_vol = $1,
	quote_vol = $2,
	high24 = $3,
	low24 = $4,
	highest_bid = $5,
	lowest_ask = $6,
	last = $7,
	percent_change = $8
WHERE
	pair = $9
`
	_, err := api.pg.Exec(sql, api.BaseVol, api.QuoteVol, api.High24, api.Low24, api.HighestBid, api.LowestAsk,
		api.Last, api.Percent, api.Pair)
	if err != nil {
		return err
	}

	return nil
}
