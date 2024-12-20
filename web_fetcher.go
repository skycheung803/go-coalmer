package coalmer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

type WebFetcher struct {
	headless bool
	Client   *rod.Browser
}

type WebFetcherOption func(*WebFetcher)

func WithClient(browser *rod.Browser) WebFetcherOption {
	return func(wf *WebFetcher) {
		wf.Client = browser
	}
}

func WithHeadless(headless bool) WebFetcherOption {
	return func(wf *WebFetcher) {
		wf.headless = headless
	}
}

// LaunchBrowser quick launch browser with headless mode
func LaunchBrowser(headless bool) *rod.Browser {
	path, _ := launcher.LookPath()
	serviceURL, err := launcher.New().Bin(path).Headless(headless).Launch()
	if err != nil {
		log.Fatal(err)
	}

	browser := rod.New().ControlURL(serviceURL).NoDefaultDevice()
	if err := browser.Connect(); err != nil {
		log.Fatal(err)
	}

	return browser
}

// NewWebFetcher create new web fetcher
func NewWebFetcher(options ...WebFetcherOption) *WebFetcher {
	wf := &WebFetcher{
		headless: true, // 设置默认值为 true
	}

	for _, option := range options {
		option(wf)
	}

	if wf.Client == nil {
		browser := LaunchBrowser(wf.headless)
		wf.Client = browser
	}

	return wf
}

// WebSearchParse parse search condition to url
// @link  https://jp.mercari.com/search?brand_id=7572&price_min=10000&price_max=50000&category_id=76&item_condition_id=1,2&status=on_sale,sold_out&item_types=beyond,mercari&color_id=10,12&page_token=v1:3
func webSearchParse(p SearchData) string {
	reqVal := url.Values{}
	if p.Page > 0 {
		page_token := fmt.Sprintf("v1:%d", p.Page)
		reqVal.Add("page_token", page_token)
	}

	if p.SearchConditionId != "" {
		reqVal.Add("search_condition_id", p.SearchConditionId)
		link := fmt.Sprintf("%s?%s", webSearchURL, reqVal.Encode())
		return link
	}

	if p.Keyword != "" {
		reqVal.Add("keyword", p.Keyword)
	}

	if len(p.BrandId) > 0 {
		brand_id := strconv.Itoa(p.BrandId[0])
		reqVal.Add("brand_id", brand_id)
	}

	if len(p.CategoryId) > 0 {
		category_id := strconv.Itoa(p.CategoryId[0])
		reqVal.Add("category_id", category_id)
	}

	if len(p.ConditionId) > 0 {
		item_condition_id := IntSliceToString(p.ConditionId)
		reqVal.Add("item_condition_id", item_condition_id)
	}

	if len(p.ColorId) > 0 {
		color_id := IntSliceToString(p.ColorId)
		reqVal.Add("color_id", color_id)
	}

	if p.PriceMin > p.PriceMax {
		p.PriceMin, p.PriceMax = p.PriceMax, p.PriceMin
	}

	if p.PriceMin > 0 {
		reqVal.Add("price_min", strconv.Itoa(p.PriceMin))
	}
	if p.PriceMax > 0 {
		reqVal.Add("price_max", strconv.Itoa(p.PriceMax))
	}

	//sort=price&order=asc
	if p.Sort != "" && p.Order != "" {
		reqVal.Add("sort", p.Sort)
		reqVal.Add("order", p.Order)
	} else {
		reqVal.Add("sort", "created_time") //default
		reqVal.Add("order", "desc")
	}

	if len(p.Status) > 0 {
		status := strings.Join(p.Status, ",") // on_sale  trading  sold_out
		reqVal.Add("status", status)
	}

	if len(p.ItemTypes) > 0 {
		types := strings.Join(p.ItemTypes, ",") // mercari  beyond
		reqVal.Add("item_types", types)
	}

	link := fmt.Sprintf("%s/search?%s", RootURL, reqVal.Encode())
	return link
}

func (w *WebFetcher) getHtml(link string, wait_selector string) (html string) {
	//log.Println(link)
	//@link https://go-rod.github.io/#/context-and-timeout?id=timeout
	page := stealth.MustPage(w.Client) //隐身模式
	//defer page.MustClose()

	page.Timeout(time.Second * 3).MustNavigate(link).MustWaitStable()
	page.Timeout(time.Second * 1).MustEval(`() => {window.scrollTo({top: document.body.scrollHeight,behavior: 'smooth'});}`)
	if wait_selector != "" {
		page.Timeout(time.Second * 3).MustElement(wait_selector).MustWaitStable()
	}
	html = page.MustHTML()
	return
}

func (w *WebFetcher) Index(limit int) (response IndexProductsResponse, err error) {
	html := w.getHtml(webIndexURL, "#item-grid")
	htmlElement, err := ParseHtml(html)
	if err != nil {
		return
	}

	var res searchSelector
	if err = htmlElement.Unmarshal(&res); err != nil {
		return
	}

	items := make([]SellerItem, 0)
	for _, item := range res.Items {
		if item.ID == "" && item.Name == "" {
			continue
		}
		status := "on_sale"
		if strings.Contains(item.Status, "売り切れ") {
			status = "sold_out"
		}

		if strings.Contains(item.Label, "HK") {
			labels := strings.Split(item.Label, " ")
			l := len(labels)
			if l >= 3 {
				item.Price = parsePrice(labels[l-2])
			}
		}

		items = append(items, SellerItem{
			RelatedItem: RelatedItem{
				ID:         item.ID,
				ItemType:   item.ItemType,
				Name:       item.Name,
				Thumbnails: []string{item.Image},
				Price:      parsePriceInt(item.Price),
				Status:     status,
			},
		})

		limit--
		if limit == 0 {
			break
		}
	}

	response.Result = "OK"
	response.Data = items

	return
}

func (w *WebFetcher) Search(params SearchData) (response SearchResponse, err error) {
	link := webSearchParse(params)
	html := w.getHtml(link, "#item-grid")
	htmlElement, err := ParseHtml(html)
	if err != nil {
		return
	}

	var res searchSelector
	if err = htmlElement.Unmarshal(&res); err != nil {
		return
	}

	items := make([]Item, 0)
	for _, item := range res.Items {
		if item.ID == "" && item.Name == "" {
			continue
		}
		status := "on_sale"
		if strings.Contains(item.Status, "売り切れ") {
			status = "sold_out"
		}

		if strings.Contains(item.Label, "HK") {
			labels := strings.Split(item.Label, " ")
			l := len(labels)
			if l >= 3 {
				item.Price = parsePrice(labels[l-2])
			}
		}

		items = append(items, Item{
			ProductId:   item.ID,
			ItemType:    item.ItemType,
			ProductName: item.Name,
			Thumbnails:  []string{item.Image},
			Price:       item.Price,
			Status:      status,
		})
	}

	response.Result = "OK"
	response.Items = items

	next_page_token := ""
	prev_page_token := ""
	search_condition_id := ""
	if res.Meta.Pager.NextPageUrl != "" {
		next_page_token = GetQueryParam(res.Meta.Pager.NextPageUrl, "page_token")
		search_condition_id = GetQueryParam(res.Meta.Pager.NextPageUrl, "search_condition_id")
	}

	if res.Meta.Pager.PrevPageUrl != "" {
		prev_page_token = GetQueryParam(res.Meta.Pager.PrevPageUrl, "page_token")
	}

	meta := map[string]interface{}{
		"nextPageToken":     next_page_token,
		"previousPageToken": prev_page_token,
	}
	response.Meta = meta
	response.SearchConditionId = search_condition_id
	return
}

func (w *WebFetcher) Detail(id string) (response ItemResultResponse, err error) {
	if len(id) > 15 {
		return w.ShopItem(id) //B2C
	} else {
		return w.Item(id) // C2C
	}
}

func (w *WebFetcher) Item(id string) (response ItemResultResponse, err error) {
	link := fmt.Sprintf("%s/%s", webItemURL, id)
	html := w.getHtml(link, "#main")
	htmlElement, err := ParseHtml(html)
	if err != nil {
		return
	}

	var product detailProduct
	var res detailSelector
	if err = htmlElement.Unmarshal(&res); err != nil {
		err = fmt.Errorf("detail page can not unmarshal Result struct: %s", err.Error())
		return
	}

	if res.Product == "" {
		err = fmt.Errorf("detail page can not get Product json data")
		return
	}

	err = json.Unmarshal([]byte(res.Product), &product)
	if err != nil {
		err = fmt.Errorf("detail page can not unmarshal Product struct: %s", err.Error())
		return
	}

	response.Result = "OK"
	response.Data.Url = product.Offers.URL
	response.Data.ProductId = product.ProductID
	response.Data.ProductName = product.Name
	response.Data.Price = product.Offers.Price
	response.Data.Photos = product.Image
	response.Data.Description = product.Description
	////on_sale trading sold_out
	if strings.Contains(product.Offers.Availability, "SoldOut") {
		response.Data.Status = "sold_out"
	} else {
		response.Data.Status = "on_sale"
	}

	for _, v := range res.Data.Info {
		if strings.Contains(v.Title, "ブランド") {
			response.Data.Brand.Id = int64(GetQueryParamInt(v.Url, "brand_id"))
			response.Data.Brand.Name = v.Body
		}

		if strings.Contains(v.Title, "商品の状態") {
			response.Data.Condition.Name = v.Body
		}

		if strings.Contains(v.Title, "発送元の地域") {
			response.Data.ShippingFrom.Name = v.Body
		}

		if strings.Contains(v.Title, "商品のサイズ") {
			response.Data.ItemSize.Name = v.Body
		}
	}

	for _, v := range res.Data.Categories {
		response.Data.Categories = append(response.Data.Categories, Name_Id_Unit{
			Id:   int64(GetQueryParamInt(v.Url, "category_id")),
			Name: v.Name,
		})
	}

	l := len(res.Data.Categories)
	if l > 0 {
		response.Data.ItemCategory = Name_Id_Unit{
			Id:   int64(res.Data.Categories[l-1].Id),
			Name: res.Data.Categories[l-1].Name,
		}
	}

	response.Data.ShippingPayer.Name = strings.TrimSpace(strings.Replace(res.Data.ShippingPayerStr, "(税込)", "", 1))

	if res.Data.SellerInfo.Url != "" {
		IDStr := strings.Replace(res.Data.SellerInfo.Url, "/user/profile/", "", 1)
		id, _ := strconv.Atoi(IDStr)
		response.Data.Seller.ID = int64(id)
		response.Data.Seller.Code = IDStr
	}

	if res.Data.SellerInfo.Avatar != "" {
		response.Data.Seller.Avatar = res.Data.SellerInfo.Avatar
	}

	if res.Data.SellerInfo.Label != "" {
		label := strings.Split(res.Data.SellerInfo.Label, ",")
		l := len(label)
		if l > 0 {
			response.Data.Seller.Name = label[0]
		}

		if l > 1 {
			NumRatings, _ := strconv.Atoi(strings.TrimSpace(strings.Replace(label[1], "件のレビュー", "", 1)))
			response.Data.Seller.NumRatings = NumRatings
		}

		if l > 2 {
			Rating, _ := strconv.ParseFloat(strings.TrimSpace(strings.Replace(label[2], "段階評価中5", "", 1)), 64)
			response.Data.Seller.Rating = Rating
		}

		if l > 3 && strings.Contains(label[3], "本人確認") {
			response.Data.Seller.SmsAuth = "yes"
		}
	}

	for _, v := range res.Data.SellerBadges {
		if strings.Contains(v, "24時間以内発送") {
			response.Data.Seller.QuickShipper = true
		}
	}
	response.Data.UpdatedStr = res.Data.UpdatedStr

	for _, v := range res.Data.Related {
		if strings.Contains(v.Label, "HK") {
			labels := strings.Split(v.Label, " ")
			l := len(labels)
			if l >= 3 {
				v.Price = labels[l-2]
			}
		}
		originalPrice := parsePriceInt(v.Price)

		status := "on_sale"
		if strings.Contains(v.Status, "売り切れ") {
			status = "sold_out"
		}

		thumbnails := []string{v.Thumbnail}
		response.Data.Related = append(response.Data.Related, RelatedItem{
			ID:         v.ID,
			Name:       v.Name,
			Price:      originalPrice,
			Status:     status,
			Thumbnail:  v.Thumbnail,
			Thumbnails: thumbnails,
			ItemType:   v.ItemType,
		})
	}

	response.Data.Colors = []Name_Id_Unit{} //@todo
	response.Meta = map[string]interface{}{}

	return response, nil
}

func (w *WebFetcher) ShopItem(id string) (response ItemResultResponse, err error) {
	link := fmt.Sprintf("%s%s", webShopsItemURL, id)
	html := w.getHtml(link, "#main")
	htmlElement, err := ParseHtml(html)
	if err != nil {
		return
	}

	var product detailShopProduct
	var res detailShopSelector
	if err = htmlElement.Unmarshal(&res); err != nil {
		err = fmt.Errorf("shops detail page can not unmarshal Result struct: %s", err.Error())
		return
	}

	if res.Product == "" {
		err = fmt.Errorf("shops detail page can not get Product json data")
		return
	}

	err = json.Unmarshal([]byte(res.Product), &product)
	if err != nil {
		err = fmt.Errorf("shops detail page can not unmarshal Product struct: %s", err.Error())
		return
	}

	response.Result = "OK"
	response.Data.ProductId = id
	response.Data.Url = product.Offers.URL
	response.Data.ProductName = product.Name
	response.Data.Price = product.Offers.Price
	response.Data.Description = product.Description
	response.Data.Photos = res.Data.Photos
	////on_sale trading sold_out
	if strings.Contains(product.Offers.Availability, "SoldOut") {
		response.Data.Status = "sold_out"
	} else {
		response.Data.Status = "on_sale"
	}

	for _, v := range res.Data.Info {
		if strings.Contains(v.Title, "ブランド") {
			response.Data.Brand.Id = int64(GetQueryParamInt(v.Url, "brand_id"))
			response.Data.Brand.Name = v.Body
		}

		if strings.Contains(v.Title, "商品の状態") {
			response.Data.Condition.Name = v.Body
		}

		if strings.Contains(v.Title, "発送元の地域") {
			response.Data.ShippingFrom.Name = v.Body
		}

		if strings.Contains(v.Title, "商品のサイズ") {
			response.Data.ItemSize.Name = v.Body
		}
	}

	for k, v := range res.Data.Categories {
		//@todo 5 level
		//https://jp.mercari.com/search?category_id=345
		//https://jp.mercari.com/category/345
		if k < 4 {
			response.Data.Categories = append(response.Data.Categories, Name_Id_Unit{
				Id:   int64(GetQueryParamInt(v.Url, "category_id")),
				Name: v.Name,
			})
		}
	}

	l := len(response.Data.Categories)
	if l > 0 {
		response.Data.ItemCategory = Name_Id_Unit{
			Id:   int64(response.Data.Categories[l-1].Id),
			Name: response.Data.Categories[l-1].Name,
		}
	}

	response.Data.ShippingPayer.Name = strings.TrimSpace(strings.Replace(res.Data.ShippingPayerStr, "(税込)", "", 1))

	response.Data.Seller.Code = id
	if res.Data.SellerInfo.Avatar != "" {
		response.Data.Seller.Avatar = res.Data.SellerInfo.Avatar
	}

	if res.Data.SellerInfo.Label != "" {
		response.Data.Seller.Name = res.Data.SellerInfo.Label
	}

	response.Data.Seller.NumRatings, _ = strconv.Atoi(res.Data.SellerInfo.NumRatingsStr)
	response.Data.Seller.Rating, _ = strconv.ParseFloat(res.Data.SellerInfo.StarRatingScoreStr, 64)
	response.Data.UpdatedStr = res.Data.UpdatedStr

	for _, v := range res.Data.Related {
		if strings.Contains(v.Label, "HK") {
			labels := strings.Split(v.Label, " ")
			l := len(labels)
			if l >= 3 {
				v.Price = labels[l-2]
			}
		}
		originalPrice := parsePriceInt(v.Price)

		status := "on_sale"
		if strings.Contains(v.Status, "売り切れ") {
			status = "sold_out"
		}

		thumbnails := []string{v.Thumbnail}
		response.Data.Related = append(response.Data.Related, RelatedItem{
			ID:         v.ID,
			Name:       v.Name,
			Price:      originalPrice,
			Status:     status,
			Thumbnail:  v.Thumbnail,
			Thumbnails: thumbnails,
			ItemType:   v.ItemType,
		})
	}

	response.Data.Colors = []Name_Id_Unit{} //@todo
	response.Meta = map[string]interface{}{}

	return
}

func (w *WebFetcher) Seller(seller_id, pager_id string) (response SellerProductsResponse, err error) {
	if len(seller_id) > 15 {
		return w.ShopProducts(seller_id, pager_id) // B2C
	} else {
		return w.SellerProducts(seller_id, pager_id) // C2C
	}
}

// SellerProducts is C2C merchant products
func (w *WebFetcher) SellerProducts(seller_id, pager_id string) (response SellerProductsResponse, err error) {
	link := fmt.Sprintf("%s%s?status=on_sale", webSellerURL, seller_id)
	html := w.getHtml(link, "#main")
	htmlElement, err := ParseHtml(html)
	if err != nil {
		return
	}

	var res sellerSelector
	if err = htmlElement.Unmarshal(&res); err != nil {
		return
	}

	response.Result = "OK"
	for _, v := range res.Goods {
		response.Data = append(response.Data, SellerItem{
			RelatedItem: RelatedItem{
				ID:         v.ID,
				Name:       v.Name,
				Price:      parsePriceInt(v.Price),
				Thumbnail:  v.Image,
				Thumbnails: []string{v.Image},
				ItemType:   v.ItemType,
			},
		})
	}

	id, _ := strconv.Atoi(seller_id)
	profile := ProfileData{
		ID:           id,
		Name:         res.Meta.Seller.Name,
		PhotoURL:     res.Meta.Seller.PhotoURL,
		Introduction: res.Meta.Seller.Introduction,
	}

	response.Profile = profile
	response.Profile.NumRatings, _ = strconv.Atoi(res.Meta.Seller.NumRatingsStr)
	response.Profile.StarRatingScore, _ = strconv.Atoi(res.Meta.Seller.StarRatingScoreStr)
	response.Profile.NumSellItems, _ = strconv.Atoi(res.Meta.Seller.NumSellItemsStr)
	response.Profile.FollowerCount, _ = strconv.Atoi(res.Meta.Seller.FollowerCountStr)
	response.Profile.FollowingCount, _ = strconv.Atoi(res.Meta.Seller.FollowingCountStr)

	return response, nil
}

// // ShopProducts is B2C  merchant products
func (w *WebFetcher) ShopProducts(seller_id, pager_id string) (response SellerProductsResponse, err error) {
	//link := fmt.Sprintf("%s%s?in_stock=true", webShopsSellerURL, seller_id) //order_by=PRICE_ASC
	//link := fmt.Sprintf("%s%s?order_by=PRICE_ASC", webShopsSellerURL, seller_id) //order_by=PRICE_ASC
	link := fmt.Sprintf("%s%s", webShopsSellerURL, seller_id) //order_by=PRICE_ASC
	html := w.getHtml(link, "#__next")
	htmlElement, err := ParseHtml(html)
	if err != nil {
		return
	}
	var res sellerShopSelector
	if err = htmlElement.Unmarshal(&res); err != nil {
		return
	}

	response.Result = "OK"
	for _, v := range res.Goods {
		status := "on_sale"
		if v.Status != "" {
			status = "sold_out"
		}
		response.Data = append(response.Data, SellerItem{
			RelatedItem: RelatedItem{
				ID:         getIdFromLink(v.Url),
				Name:       v.Name,
				Price:      parsePriceInt(v.Price),
				Thumbnail:  v.Image,
				Thumbnails: []string{v.Image},
				ItemType:   "ITEM_TYPE_BEYOND",
				Status:     status,
			},
		})
	}

	profile := ProfileData{
		Code:     seller_id,
		Name:     res.Meta.Seller.Name,
		PhotoURL: res.Meta.Seller.PhotoURL,
	}

	shop_info, err := simplejson.NewJson([]byte(res.Props))
	if err == nil {
		info := shop_info.Get("props").Get("pageProps").Get("__APOLLO_STATE__").MustMap()
		for k, v := range info {
			if strings.Contains(k, "Shop:") {
				jsonData, err := json.Marshal(v)
				if err != nil {
					log.Fatal("Error marshalling JSON:", err)
					break
				}

				var data map[string]interface{}
				if err := json.Unmarshal(jsonData, &data); err != nil {
					log.Fatal("Error unmarshalling JSON:", err)
					break
				}

				//profile.Name = data["name"].(string)
				profile.Introduction = data["description"].(string)
				//profile.FollowerCount = data["followerCount"].(int)
				//profile.StarRatingScore = int(math.Round(data["reviewStats"]["score"].(float64)))
				break
			}
		}
	}

	response.Profile = profile
	response.Profile.NumRatings, _ = strconv.Atoi(res.Meta.Seller.NumRatingsStr)
	response.Profile.StarRatingScore = len(res.Meta.Seller.StarRatingScoreStr)
	response.Profile.NumSellItems, _ = strconv.Atoi(res.Meta.Seller.NumSellItemsStr)
	response.Profile.FollowerCount, _ = strconv.Atoi(res.Meta.Seller.FollowerCountStr)
	//response.Profile.FollowingCount, _ = strconv.Atoi(res.Meta.Seller.FollowingCountStr)

	return response, nil
}

func getIdFromLink(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}
	if u.Path == "" {
		return ""
	}
	parts := strings.Split(u.Path, "/")
	return parts[len(parts)-1]
}
