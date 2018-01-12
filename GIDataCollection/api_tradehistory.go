/*
	历史成交记录
	返回最新80条历史成交记录
	http://data.gate.io/api2/1/tradeHistory/[CURR_A]_[CURR_B]

	返回从[TID]往后的最多1000历史成交记录：
	http://data.gate.io/api2/1/tradeHistory/[CURR_A]_[CURR_B]/[TID]
	请替换 [CURR_A] and [CURR_B] 为您需要查看的币种.
*/

package main

import (
	"strconv"
	"github.com/buger/jsonparser"
	"errors"
	"encoding/json"
)

type ApiTradeHistory struct {
	Api
	Pair   string
	Tid    int64
	Result bool
	Histories []ApiTradeHistoryItem
}

type ApiTradeHistoryItem struct {
	TradeID		int64		`json:"tradeID,string"`		// tradeID
	TimeStamp	int64		`json:"timestamp,string"`	// timestamp
	Type		string		`json:"type"`				// 交易类型, buy买 sell卖
	Price		float64		`json:"rate"`				// 币种单价
	Amount		float64		`json:"amount"`				// 成交币种数量
	Total		float64		`json:"total"`				// 订单总额
	Date		string		`json:"date"`				// 订单时间
}

// 如果 tid = 0，则返回最新的80条
func (api *ApiTradeHistory) Init(pg *Postgres, pair string, tid int64) (*ApiTradeHistory) {
	stid := strconv.Itoa(int(tid))

	api.desc = "历史成交记录"
	if tid != 0 {
		api.desc += "tid = " + stid
	}

	api.uri = "tradeHistory/" + pair
	if tid != 0 {
		api.uri += "/" + stid
	}

	api.pg = pg
	api.Pair = pair
	api.Tid = tid

	return api
}

func (api *ApiTradeHistory) Request() ([]byte, error) {
	return api.httpGet(api.uri)
}

func (api *ApiTradeHistory) Parser(body []byte) (error) {
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

	var histories []ApiTradeHistoryItem
	_, err =jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var history ApiTradeHistoryItem
		json.Unmarshal(value, &history)
		histories = append(histories, history)
	}, "data")
	if err != nil {
		return err
	}

	api.Histories = histories
	return nil
}

func (api *ApiTradeHistory) Save() (error) {
	for _, history := range api.Histories {
		sql := "SELECT COUNT(*) AS row_count FROM tr_history WHERE trade_id = $1 AND pair = $2"
		dataSet, err := api.pg.Query(sql, history.TradeID, api.Pair)
		if err != nil {
			return err
		}

		rowCount := dataSet[0]["row_count"].(int64)
		if rowCount == 0 {
			sql = "INSERT INTO tr_history (trade_id, pair, type, price, amount, total, timestamp) VALUES ($1, $2, $3, $4, $5, $6, to_timestamp($7))"
			_, err := api.pg.Exec(sql, history.TradeID, api.Pair, history.Type, history.Price, history.Amount,
				history.Total, history.TimeStamp)
			if err != nil {
				return err
			}
		} else {
			sql = `
UPDATE 
	tr_history 
SET
	type = $1,
	price = $2,
	amount = $3,
	total = $4,
	timestamp = to_timestamp($5)
WHERE
	trade_id = $6
	AND pair = $7
`
			_, err = api.pg.Exec(sql, history.Type, history.Price, history.Amount, history.Total, history.TimeStamp,
				history.TradeID, api.Pair)
			if err != nil {
				return err
			}
		}
	}

	return nil
}