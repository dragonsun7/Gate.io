package main

import (
	//"time"
	"fmt"
	"log"
)

func main() {
	fmt.Print("连接数据库...")
	pg := new(Postgres)
	err := pg.Open()
	defer pg.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("成功！")

	//times := 1
	//for true {
	//	fmt.Print("等待...")
	//
	//	time := time.NewTimer(time.Second * 10)
	//	<-time.C
	//
	//	fmt.Println("采集数据...", times)

		//pairsApi := new(ApiPairs).Init(pg)
		//err = ApiDo(pairsApi)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//marketInfoApi := new(ApiMarketInfo).Init(pg)
		//err = ApiDo(marketInfoApi)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//markerListApi := new(ApiMarketList).Init(pg)
		//err = ApiDo(markerListApi)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//tickerApi := new(ApiTicker).Init(pg, "btc_usdt")
		//err = ApiDo(tickerApi)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//tickersApi := new(ApiTickers).Init(pg)
		//err = ApiDo(tickersApi)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//orderBookApi := new(ApiOrderBook).Init(pg, "btc_usdt")
		//err = ApiDo(orderBookApi)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//orderBooksApi := new(ApiOrderBooks).Init(pg)
		//err = ApiDo(orderBooksApi)
		//if err != nil {
		//	fmt.Println(err)
		//}

		tradeHistoryApi := new(ApiTradeHistory).Init(pg, "btc_usdt", 0)
		err = ApiDo(tradeHistoryApi)
		if err != nil {
			fmt.Println(err)
		}


	//	times++
	//	fmt.Println("采集结束！")
	//	fmt.Println("")
	//}
}
