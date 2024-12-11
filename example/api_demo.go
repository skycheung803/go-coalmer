/*
package main

import (
	"fmt"

	"github.com/skycheung803/go-coalmer"
)

var (
	apier *coalmer.APIFetcher
)

func init() {
	apier = coalmer.NewAPIFetcher(false)
}

func main() {
	detail()
	//related()
	//search()
	//seller()
	//similarLooks()
	//profile()
}

func detail() {
	res, err := apier.Detail("m13616196726")
	if err != nil {
		panic(err)
	}
	coalmer.Dump(res)

	fmt.Println("--------------------------------------------------")
	res, err = apier.Detail("mscq8am7HoHsqAisb8Qxga")
	if err != nil {
		panic(err)
	}
	coalmer.Dump(res)
}

func search() {
	p := coalmer.SearchData{
		Keyword: "iPhone 15",
	}
	res, err := apier.Search(p)
	if err != nil {
		panic(err)
	}
	coalmer.Dump(res)
}

func seller() {
	res, err := apier.Seller("436032940", "")
	if err != nil {
		panic(err)
	}
	coalmer.Dump(res)

	fmt.Println("--------------------------------------------------")
	res, err = apier.Seller("TcPxqCaTaFgNtgHpcLzhCG", "")
	if err != nil {
		panic(err)
	}
	coalmer.Dump(res)
}

func related() {
	res, err := apier.Related("m13616196726", "15")
	if err != nil {
		panic(err)
	}
	coalmer.Dump(res)
}

func similarLooks() {
	s := coalmer.SimilarData{
		ItemID: "m13616196726",
		//"mscq8am7HoHsqAisb8Qxga",
	}
	res, err := apier.SimilarLooks(s)
	if err != nil {
		panic(err)
	}
	coalmer.Dump(res)
}

func profile() {
	//res, err := apier.Profile("182093486")
	res, err := apier.Profile("TcPxqCaTaFgNtgHpcLzhCG")
	if err != nil {
		panic(err)
	}
	coalmer.Dump(res)
}
*/