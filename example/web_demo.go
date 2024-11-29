/*
package main

import (
	"github.com/skycheung803/go-coalmer"
)

var (
	webFetcher *coalmer.WebFetcher
)

func init() {
	webFetcher = coalmer.NewWebFetcher(true)
}

func main() {
	seller()
	//detail()
	//related()
	//search()
	//similarLooks()
	//runtime.Goexit()
}

func search() {
	p := coalmer.SearchData{
		Keyword:           "iPhone 15",
		ConditionId:       []int{1},
		Page:              5,
		SearchConditionId: "1cx0zHGljZB0xHGsdaVBob25lIDE1",
	}
	res, err := webFetcher.Search(p)
	if err != nil {
		panic(err)
	}

	//fmt.Println(len(res.Items))
	coalmer.Dump(res)

}

func detail() {
	res, err := webFetcher.Detail("m13616196726")
	if err != nil {
		panic(err)
	}
	coalmer.Dump(res)
}

func seller() {
	res, err := webFetcher.Seller("436032940", "")
	if err != nil {
		panic(err)
	}
	coalmer.Dump(res)
}
*/