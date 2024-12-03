package main

import (
	"log"

	"github.com/skycheung803/go-coalmer"
)

var (
	Coalmer *coalmer.Coalmer
)

func init() {
	log.Println("start init ~~~ ")
	//Coalmer = coalmer.NewCoalmer()
	Coalmer = coalmer.NewCoalmer(coalmer.WithBrowserMode())
	log.Println("init  finish~~~ ")
}

func main() {
	log.Println("--------------------- start ---------------------")
	search()
	log.Println("--------------------- search Finish -----------------------------")
	//detail()
	log.Println("--------------------- detail Finish -----------------------------")
	//seller()
	log.Println("--------------------- seller Finish -----------------------------")
}

func detail() {
	res, err := Coalmer.Fetcher.Detail("f5s7W74buc5RAViYnqPggD")
	if err != nil {
		panic(err)
	}
	coalmer.Dump(res)
	//fmt.Println("--------------------------------------------------")

	/* res, err = Coalmer.Fetcher.Detail("m86235728638")
	if err != nil {
		panic(err)
	}
	coalmer.Dump(res) */
	//log.Println("---------------------Finish-----------------------------")
}

func seller() {
	res, err := Coalmer.Fetcher.Seller("GqU4Yahsuz6LW3NZZR53T8", "")
	if err != nil {
		panic(err)
	}
	coalmer.Dump(res)
}

func search() {
	p := coalmer.SearchData{
		Keyword:     "iPhone 15",
		ConditionId: []int{1},
		Sort:        "price",
		Order:       "asc",
		//Page:              5,
		//SearchConditionId: "1cx0zHGljZB0xHGsdaVBob25lIDE1",
	}
	res, err := Coalmer.Fetcher.Search(p)
	if err != nil {
		panic(err)
	}

	//fmt.Println(len(res.Items))
	coalmer.Dump(res)

}
