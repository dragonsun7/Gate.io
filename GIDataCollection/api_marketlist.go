/*
	交易市场详细行情
	返回所有系统支持的交易市场的详细行情和币种信息，包括币种名，市值，供应量，最新价格，涨跌趋势，价格曲线等。
	http://data.gate.io/api2/1/marketlist
*/

package main

import (
	"github.com/buger/jsonparser"
	"strconv"
	"errors"
	"fmt"
	"strings"
)

type ApiMarketList struct {
	Api
	Result bool					`json:"result"`
	Data []ApiMarketListItem	`json:"data"`
}

type ApiMarketListItem struct {
	Pair 		string			`json:"pair"`					// 交易对
	Symbol 		string			`json:"symbol"`					// 币种标识
	Name 		string			`json:"name"`					// 币种名称
	NameEn 		string			`json:"name_en"`				// 英文名称
	NameCn 		string			`json:"name_cn"`				// 中文名称
	CurrA 		string			`json:"curr_a"`					// 被兑换货币
	CurrB 		string			`json:"curr_b"`					// 兑换货币
	CurrSuffix 	string			`json:"curr_suffix"`			// 货币类型后缀
	Supply 		int64			`json:"supply"`					// 币种供应量
	MarketCap 	int64			`json:"marketcap,string"`		// 总市值
	Trend		string			`json:"trend"`					// 24小时趋势(up—涨, down—跌, flat—平)
	Price 		float64			`json:"rate,string"`			// 当前价格
	VolA		float64			`json:"vol_a"`					// 被兑换货币交易量
	VolB		float64			`json:"vol_b,string"`			// 兑换货币交易量
	Percent		float64			`json:"rate_percent,string"`	// 涨跌百分比
}

func (api *ApiMarketList) Init(pg *Postgres) (*ApiMarketList) {
	api.desc = "交易市场详细行情"
	api.uri = "marketlist"
	api.pg = pg
	return api
}

func (api *ApiMarketList) Request() ([]byte, error) {
	return api.httpGet(api.uri)
}

func (api *ApiMarketList) Parser(body []byte) (error) {
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

	var data []ApiMarketListItem
	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var item ApiMarketListItem

		sMarketCap, _ := jsonparser.GetString(value, "marketcap")
		sPrice, _ := jsonparser.GetString(value, "rate")
		sVolB, _ := jsonparser.GetString(value, "vol_b")
		sPercent, _ := jsonparser.GetString(value, "rate_percent")

		item.Pair, _ = jsonparser.GetString(value, "pair")
		item.Symbol, _ = jsonparser.GetString(value, "symbol")
		item.Name, _ = jsonparser.GetString(value, "name")
		item.NameEn, _ = jsonparser.GetString(value, "name_en")
		item.NameCn, _ = jsonparser.GetString(value, "name_cn")
		item.CurrA, _ = jsonparser.GetString(value, "curr_a")
		item.CurrB, _ = jsonparser.GetString(value, "curr_b")
		item.CurrSuffix, _ = jsonparser.GetString(value, "curr_suffix")
		item.Supply, _ = jsonparser.GetInt(value, "supply")
		item.MarketCap = StrToInt(sMarketCap)
		item.Price = StrToFloat(sPrice)
		item.VolA, _ = jsonparser.GetFloat(value, "vol_a")
		item.VolB = StrToFloat(sVolB)
		item.Percent = StrToFloat(sPercent)
		item.Trend, _ = jsonparser.GetString(value, "trend")

		data = append(data, item)
	}, "data")
	if err != nil {
		return err
	}

	api.Data = data
	return nil
}

func (api *ApiMarketList) Save() (error) {
	for _, item := range api.Data  {

		// bs_pairs
		sql := "SELECT COUNT(*) AS row_count FROM bs_pairs WHERE pair = $1"
		dataSet, err := api.pg.Query(sql, item.Pair)
		if err != nil {
			return err
		}

		rowCount := dataSet[0]["row_count"].(int64)
		if rowCount == 0 {
			return errors.New("没有找到交易对：" + item.Pair)
		}

		sql = `
UPDATE
	bs_pairs
SET
	symbol = $1, 
	name = $2, 
	name_en = $3, 
	name_cn = $4,
	curr_a = $5,
	curr_b = $6,
	supply = $7,
	market_cap = $8
WHERE 
	pair = $9
`
		_, err = api.pg.Exec(sql, item.Symbol, item.Name, item.NameCn, item.NameEn, item.CurrA, item.CurrB, item.Supply,
			item.MarketCap, item.Pair)
		if err != nil {
			fmt.Println(item)
			return err
		}

		parts := strings.Split(item.Pair, "_")
		currB := strings.ToUpper(parts[1])
		match := parts[0] + "%"
		sql = `
UPDATE
	bs_pairs
SET
	symbol = $1, 
	name = $2, 
	name_en = $3, 
	name_cn = $4,
	curr_a = $5,
	curr_b = $6
WHERE 
	pair like $7 
`
		_, err = api.pg.Exec(sql, item.Symbol, item.Name, item.NameCn, item.NameEn, item.CurrA, currB, match)
		if err != nil {
			return err
		}

		// tr_market
		sql = "SELECT COUNT(*) AS row_count FROM tr_market WHERE pair = $1"
		dataSet, err = api.pg.Query(sql, item.Pair)
		if err != nil {
			return err
		}

		rowCount = dataSet[0]["row_count"].(int64)
		if rowCount == 0 {
			return errors.New("没有找到交易对：" + item.Pair)
		}

		trend := 0
		if item.Trend == "up" {
			trend = 1
		}
		if item.Trend == "down" {
			trend = -1
		}
		sql = "UPDATE tr_market SET price = $1, vol_a = $2, vol_b = $3, percent = $4, trend = $5 WHERE pair = $6"
		_, err = api.pg.Exec(sql, item.Price, item.VolA, item.VolB, item.Percent, trend, item.Pair)
		if err != nil {
			fmt.Println("x")
			return err
		}
	}

	return nil
}