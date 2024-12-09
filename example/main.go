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

	Coalmer = coalmer.NewCoalmer()
	//browser := coalmer.LaunchBrowser(false)
	//Coalmer = coalmer.NewCoalmer(coalmer.WithBrowserMode(), coalmer.WithBrowser(browser))
	log.Println("init  finish~~~ ")
}

func main() {
	log.Println("--------------------- start ---------------------")
	search()
	log.Println("--------------------- search Finish -----------------------------")
	detail("m12958444254")
	//detail("m58457143925")
	log.Println("--------------------- detail 1 Finish -----------------------------")
	detail("m97792326581")
	log.Println("--------------------- detail 2 Finish -----------------------------")

	seller()
	log.Println("--------------------- seller Finish -----------------------------")
}

func detail(id string) {
	//res, err := Coalmer.Fetcher.Detail("m97792326581") // m97792326581 7HnYy2wC4begbaif4BXTf5
	res, err := Coalmer.Fetcher.Detail(id) // m97792326581 7HnYy2wC4begbaif4BXTf5
	if err != nil {
		panic(err)
	}
	//coalmer.Dump(res)
	fmt.Println(res.Data.ProductName)
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
	res, err := Coalmer.Fetcher.Seller("755873977", "") // 755873977   GqU4Yahsuz6LW3NZZR53T8
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
	params := `{"keyword":"","category_id":[4],"sort":"","item_condition_id":[],"item_types":["beyond"],"price_min":0,"price_max":0,"page":0,"brand_id":[],"color_id":[],"status":[]}`
	json.Unmarshal([]byte(params), &p)
	fmt.Printf("%+v\n", p)
	//coalmer.Dump(p)
	res, err := Coalmer.Fetcher.Search(p)
	if err != nil {
		panic(err)
	}

	log.Println(len(res.Items))
	//coalmer.Dump(res)

}
