package coalmer

import (
	"encoding/json"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/skycheung803/colly"
)

func GetQueryParam(link, key string) string {
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}
	query, _ := url.ParseQuery(u.RawQuery)
	return query.Get(key)
}

func GetQueryParamInt(url, key string) int {
	i, err := strconv.Atoi(GetQueryParam(url, key))
	if err != nil {
		return 0
	}
	return i
}

func GetQueryParamFloat(url, key string) float64 {
	f, err := strconv.ParseFloat(GetQueryParam(url, key), 64)
	if err != nil {
		return 0
	}
	return f
}

func GetQueryParamBool(url, key string) bool {
	return GetQueryParam(url, key) == "true"
}

func GetQueryParamSlice(url, key string) []string {
	return strings.Split(GetQueryParam(url, key), ",")
}

func GetQueryParamIntSlice(url, key string) []int {
	s := GetQueryParamSlice(url, key)
	ret := make([]int, len(s))
	for i, v := range s {
		ret[i], _ = strconv.Atoi(v)
	}
	return ret
}

func Dump(d interface{}) {
	data, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		log.Panicf("dump err:%v", err)
	}
	log.Println(string(data))
}

func ParseHtml(html string) (*colly.HTMLElement, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	htmlElement := colly.HTMLElement{
		Name: "html",
		//Request:  r.Request,
		//Response: r,
		DOM:   dom.Find("html"),
		Index: 0,
	}

	return &htmlElement, nil
}

func parsePrice(priceStr string) (price string) {
	price = strings.ReplaceAll(priceStr, ",", "")
	re := regexp.MustCompile("[0-9]+")
	prices := re.FindAllString(price, -1)
	if len(prices) > 0 {
		price = prices[0]
	}
	return price
}

func parsePriceInt(priceStr string) (price int) {
	price = 0
	priceStr = parsePrice(priceStr)
	price, _ = strconv.Atoi(priceStr)
	return price
}

// int slice to string
func IntSliceToString(intSlice []int) string {
	strSlice := make([]string, len(intSlice))
	for i, v := range intSlice {
		strSlice[i] = strconv.Itoa(v)
	}
	return strings.Join(strSlice, ",")
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
