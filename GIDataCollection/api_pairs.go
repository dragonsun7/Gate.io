/*
	所有交易对
	返回所有系统支持的交易对
	http://data.gate.io/api2/1/pairs
*/

package main

import (
	"encoding/json"
)

type ApiPairs struct {
	Api
	pairs []string
}

func (api *ApiPairs) Init(pg *Postgres) (*ApiPairs) {
	api.desc = "所有交易对"
	api.uri = "pairs"
	api.pg = pg
	return api
}

func (api *ApiPairs) Request() ([]byte, error) {
	return api.httpGet(api.uri)
}

func (api *ApiPairs) Parser(body []byte) (error) {
	return json.Unmarshal(body, &api.pairs)
}

func (api *ApiPairs) Save() (error) {
	for _, pair := range api.pairs {

		// bs_pair
		sql := "SELECT COUNT(*) AS row_count FROM bs_pairs WHERE pair = $1"
		dataSet, err := api.pg.Query(sql, pair)
		if err != nil {
			return err
		}

		rowCount := dataSet[0]["row_count"].(int64)
		if rowCount == 0 {
			sql = "INSERT INTO bs_pairs (pair, create_at) VALUES ($1, current_timestamp)"
			_, err := api.pg.Exec(sql, pair)
			if err != nil {
				return err
			}
		}

		// tr_market
		sql = "SELECT COUNT(*) AS row_count FROM tr_market WHERE pair = $1"
		dataSet, err = api.pg.Query(sql, pair)
		if err != nil {
			return err
		}

		rowCount = dataSet[0]["row_count"].(int64)
		if rowCount == 0 {
			sql = "INSERT INTO tr_market (pair) VALUES ($1)"
			_, err := api.pg.Exec(sql, pair)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
