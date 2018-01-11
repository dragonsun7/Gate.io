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
	fmt.Println("")

	//times := 1
	//for true {
	//	fmt.Print("等待...")
	//
	//	time := time.NewTimer(time.Second * 10)
	//	<-time.C
	//
	//	fmt.Println("采集数据...", times)

		pairsApi := new(ApiPairs).Init(pg)
		err = ApiDo(pairsApi)
		if err != nil {
			fmt.Println(err)
			//continue
		}

		marketInfoApi := new(ApiMarketInfo).Init(pg)
		err = ApiDo(marketInfoApi)
		if err != nil {
			fmt.Println(err)
			//continue
		}
		fmt.Println(marketInfoApi)

	//	times++
	//	fmt.Println("采集结束！")
	//	fmt.Println("")
	//}
}
