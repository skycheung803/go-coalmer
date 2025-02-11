package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/skycheung803/go-coalmer"
)

var (
	Coalmer *coalmer.Coalmer
)

func init() {
	log.Println("start init ~~~ ")

	Coalmer = coalmer.NewCoalmer(coalmer.WithDebug(false))
	//browser := coalmer.LaunchBrowser(true)
	//Coalmer = coalmer.NewCoalmer(coalmer.WithBrowserMode(), coalmer.WithBrowser(browser))
	log.Println("init  finish~~~ ")
}

func main() {
	log.Println("--------------------- start ---------------------")
	//index(20)
	detail("m43989853551")
	//detail("m43250958296")
	//detail("z7DYa2QbrbC2LwyXV9edZ7")
	//detail("8NhEyjQNedXNuJPv7YQEiP")
	//detail("8NhEyjQNedXNuJPv7YQEiPppppp")
	log.Println("--------------------- index  Finish -----------------------------")
	//search()
	log.Println("--------------------- search Finish -----------------------------")
	//detail("m12958444254")
	//detail("m58457143925")
	log.Println("--------------------- detail 1 Finish -----------------------------")
	//detail("m97792326581")
	log.Println("--------------------- detail 2 Finish -----------------------------")
	//seller()
	log.Println("--------------------- seller Finish -----------------------------")
	//runtime.Goexit()
}

func index(limit int) {

	res, err := Coalmer.Fetcher.Index(limit)
	if err != nil {
		panic(err)
	}
	coalmer.Dump(res)
	//fmt.Println(res.Data.ProductName)
}

func detail(id string) {
	//res, err := Coalmer.Fetcher.Detail("m97792326581") // m97792326581 7HnYy2wC4begbaif4BXTf5
	res, err := Coalmer.Fetcher.Detail(id) // m97792326581 7HnYy2wC4begbaif4BXTf5
	coalmer.Dump(res)
	if err != nil {
		panic(err)
	}
	//coalmer.Dump(res)
	//fmt.Println(res.Data.ProductName)
	//fmt.Println("-----------------------detail 1---------------------------")

	//res, err = Coalmer.Fetcher.Detail("7HnYy2wC4begbaif4BXTf5")
	/* res, err = Coalmer.Fetcher.Detail("m12958444254")
	if err != nil {
		panic(err)
	} */
	//coalmer.Dump(res)
	//fmt.Println(res.Data.ProductName)
	//log.Println("---------------------Finish-----------------------------")
}

func seller() {
	//res, err := Coalmer.Fetcher.Seller("755873977", "") // 755873977   GqU4Yahsuz6LW3NZZR53T8
	res, err := Coalmer.Fetcher.Seller("GqU4Yahsuz6LW3NZZR53T89999999", "") // 755873977   GqU4Yahsuz6LW3NZZR53T8
	//coalmer.Dump(res)
	if err != nil {
		panic(err)
	}
	//coalmer.Dump(res)
	log.Println(res.Profile.Name)
}

func search() {
	/*
		p := coalmer.SearchData{
			Keyword:     "iPhone 15",
			ConditionId: []int{1},
			Sort:        "price",
			Order:       "asc",
			//Page:              5,
			//SearchConditionId: "1cx0zHGljZB0xHGsdaVBob25lIDE1",
		}
	*/
	p := coalmer.SearchData{}
	//params := `{"keyword":"agd","category_id":[14],"price_min":5000,"price_max":8000,"sort":"","item_condition_id":[],"page":0}`
	//params := `{"keyword":"999UNION HORLOGERE 腕時計 19015N007SS 鑑定済み ブランド","category_id":[414],"sort":"","item_condition_id":[],"item_types":["beyond"],"price_min":0,"price_max":0,"page":0,"brand_id":[],"color_id":[],"status":[]}`
	params := `{"keyword":"UNION HORLOGERE 腕時計 19015N007SS 鑑定済み ブランド","category_id":[414],"sort":"","item_condition_id":[],"item_types":["beyond"],"price_min":0,"price_max":0,"page":0,"brand_id":[],"color_id":[],"status":[]}`
	json.Unmarshal([]byte(params), &p)
	fmt.Printf("%+v\n", p)
	//coalmer.Dump(p)
	res, err := Coalmer.Fetcher.Search(p)
	coalmer.Dump(res)
	if err != nil {
		panic(err)
	}

	log.Println(len(res.Items))
	//coalmer.Dump(res)
}
