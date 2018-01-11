/*
	交易市场订单参数
	返回所有系统支持的交易市场的参数信息，包括交易费，最小下单量，价格精度等。
	http://data.gate.io/api2/1/marketinfo
*/

package main

import (
	"github.com/buger/jsonparser"
	"errors"
	"strconv"
)

type ApiMarketInfo struct {
	Api
	Result bool					`json:"result,string"`
	Pairs []ApiMarketInfoPair	`json:"pairs"`
}

type ApiMarketInfoPair struct {
	Pair string
	Info ApiMarketInfoPairInfo
}

type ApiMarketInfoPairInfo struct {
	Decimal int64				`json:"decimal_places"`
	MinAmount float64			`json:"min_amount"`
	Fee float64					`json:"fee"`
}

func (api *ApiMarketInfo) Init(pg *Postgres) (*ApiMarketInfo) {
	api.desc = "交易市场订单参数"
	api.uri = "marketinfo"
	api.pg = pg
	return api
}

func (api *ApiMarketInfo) Request() ([]byte, error) {
	return api.httpGet(api.uri)
}

func (api *ApiMarketInfo) Parser(body []byte) (error) {
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

	var pairs []ApiMarketInfoPair
	_, err =jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

		jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			var err1 error
			var pair ApiMarketInfoPair

			pair.Pair = string(key)

			pair.Info.Decimal, err1 = jsonparser.GetInt(value, "decimal_places")
			if err1 != nil {
				return err1
			}

			pair.Info.MinAmount, err1 = jsonparser.GetFloat(value, "min_amount")
			if err1 != nil {
				return err1
			}

			pair.Info.Fee, err1 = jsonparser.GetFloat(value, "fee")
			if err1 != nil {
				return err1
			}

			pairs = append(pairs, pair)

			return nil
		})

	}, "pairs")
	if err != nil {
		return err
	}

	api.Pairs = pairs
	return nil
}

func (api *ApiMarketInfo) Save() (error) {
	for _, pair := range api.Pairs {
		sql := "UPDATE bs_pairs SET precision = $1, min_amount = $2, fee = $3 WHERE pair = $4"
		_, err := api.pg.Exec(sql, pair.Info.Decimal, pair.Info.MinAmount, pair.Info.Fee, pair.Pair)
		if err != nil {
			return err
		}
	}

	return nil
}