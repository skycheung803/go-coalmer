package coalmer

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/imroc/req/v3"
)

type APIFetcher struct {
	Client *req.Client
}

func NewAPIFetcher(debug bool) *APIFetcher {
	Client := req.C()
	Client.ImpersonateChrome() // 伪装 HTTP 指纹
	if debug {
		Client.DevMode() // 启用开发模式
	}

	return &APIFetcher{
		Client: Client,
	}
}

func apiSearchParse(p SearchData) (string, error) {
	sp := V2Search{}
	sp.DefaultDatabases = []string{"DATASET_TYPE_MERCARI", "DATASET_TYPE_BEYOND"}
	sp.IndexRouting = "INDEX_ROUTING_UNSPECIFIED"
	if p.Limit > 0 {
		sp.PageSize = p.Limit
	} else {
		sp.PageSize = DefaultPageSize
	}

	if p.Page > 0 {
		sp.PageToken = fmt.Sprintf("v1:%d", p.Page)
	}

	// searchCondition
	sp.SearchCondition.HasCoupon = false
	//sp.SearchCondition.Status = []string{"STATUS_ON_SALE"} //"STATUS_ON_SALE", "STATUS_TRADING", "STATUS_SOLD_OUT"
	sp.SearchCondition.Keyword = p.Keyword

	if p.PriceMin > p.PriceMax {
		p.PriceMin, p.PriceMax = p.PriceMax, p.PriceMin
	}

	if p.PriceMin > 0 {
		sp.SearchCondition.PriceMin = p.PriceMin
	}
	if p.PriceMax > 0 {
		sp.SearchCondition.PriceMax = p.PriceMax
	}

	if len(p.CategoryId) > 0 {
		sp.SearchCondition.CategoryId = p.CategoryId
	}

	if len(p.ConditionId) > 0 {
		sp.SearchCondition.ItemConditionId = p.ConditionId
	}

	if len(p.ColorId) > 0 {
		sp.SearchCondition.ColorId = p.ColorId
	}

	if p.Sort != "" && p.Order != "" {
		sp.SearchCondition.Sort = strings.ToUpper("SORT_" + p.Sort)
		sp.SearchCondition.Order = strings.ToUpper("ORDER_" + p.Order)
	} else {
		sp.SearchCondition.Sort = SearchOptionSortCreatedTime
		sp.SearchCondition.Order = SearchOptionOrderDESC
	}

	if len(p.Status) > 0 {
		if ok := Contains(p.Status, "on_sale"); ok {
			sp.SearchCondition.Status = append(sp.SearchCondition.Status, "STATUS_ON_SALE")
		}

		if ok := Contains(p.Status, "sold_out"); ok {
			sp.SearchCondition.Status = append(sp.SearchCondition.Status, "STATUS_SOLD_OUT", "STATUS_TRADING")
		}
	}

	if len(p.ItemTypes) > 0 {
		if ok := Contains(p.ItemTypes, "mercari"); ok {
			sp.SearchCondition.ItemTypes = append(sp.SearchCondition.ItemTypes, "ITEM_TYPE_MERCARI")
		}

		if ok := Contains(p.ItemTypes, "beyond"); ok {
			sp.SearchCondition.ItemTypes = append(sp.SearchCondition.ItemTypes, "ITEM_TYPE_BEYOND")
		}
	}

	sp.SearchSessionId = generateSearchSessionId(DefaultLengthSearchSessionId)
	sp.ServiceFrom = "suruga"
	sp.WithItemBrand = true
	sp.WithItemPromotions = true
	sp.WithItemSize = false
	sp.WithItemSizes = true
	sp.WithShopName = false

	res, err := json.Marshal(sp)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// Search product by any search param
func (a *APIFetcher) Search(params SearchData) (response SearchResponse, err error) {
	queryData, err := apiSearchParse(params)
	if err != nil {
		return
	}

	headers, err := generateHeader(searchParams.URL, searchParams.Method)
	if err != nil {
		return
	}

	_, err = a.Client.R().SetHeaders(headers).SetBody(queryData).SetSuccessResult(&response).Post(searchParams.URL)
	if err != nil {
		return
	}
	if len(response.Items) > 0 {
		response.Result = "OK"
	}
	return
}

// ProductDetail product detail
func (a *APIFetcher) Detail(id string) (response ItemResultResponse, err error) {
	related := RelatedResponse{}
	var err2 error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		if len(id) > 15 {
			response, err = a.ShopItem(id) //B2C
		} else {
			response, err = a.Item(id) // C2C
		}
		wg.Done()
	}()

	go func() {
		related, err2 = a.Related(id, "15")
		wg.Done()
	}()
	wg.Wait()

	if err != nil {
		return
	}
	if err2 == nil {
		response.Data.Related = related.Data.Items
	}

	return response, nil
}

// Item is C2C product detail
func (a *APIFetcher) Item(id string) (response ItemResultResponse, err error) {
	reqVal := url.Values{}
	reqVal.Add("id", id)
	link := fmt.Sprintf("%s?%s", itemParams.URL, reqVal.Encode())

	headers, err := generateHeader(link, itemParams.Method)
	if err != nil {
		return
	}

	var similarLooksResponse SimilarLooksResponse
	var err2 error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		_, err = a.Client.R().SetHeaders(headers).SetSuccessResult(&response).Get(link)
		wg.Done()
	}()

	go func() {
		s := SimilarData{
			ItemID: id,
		}
		similarLooksResponse, err2 = a.SimilarLooks(s)
		wg.Done()
	}()
	wg.Wait()

	if err != nil {
		return
	}
	response.Data.Url = fmt.Sprintf("%s/%s", webItemURL, id)
	response.Data.Categories = response.Data.ParentCategories
	response.Data.Categories = append(response.Data.Categories, response.Data.ItemCategoryNtiers)

	if err2 == nil {
		response.Data.SimilarLooks = similarLooksResponse.Items
	}

	return
}

// ShopItem is B2C product detail
func (a *APIFetcher) ShopItem(id string) (response ItemResultResponse, err error) {
	reqVal := url.Values{}
	reqVal.Add("view", "FULL")
	link := fmt.Sprintf("%s/%s?%s", shopItemParams.URL, id, reqVal.Encode())

	headers, err := generateHeader(link, shopItemParams.Method)
	if err != nil {
		return
	}

	resp, err := a.Client.R().SetHeaders(headers).Get(link)
	if err != nil {
		return
	}
	return ProductDetailResult(resp.Bytes())
}

// MerchantProducts merchant products
func (a *APIFetcher) Seller(seller_id, pager_id string) (SellerProductsResponse, error) {
	if len(seller_id) > 15 {
		return a.ShopProducts(seller_id, pager_id) // B2C
	} else {
		return a.SellerProducts(seller_id, pager_id) // C2C
	}
}

// SellerProducts is C2C merchant products
func (a *APIFetcher) SellerProducts(seller_id, pager_id string) (response SellerProductsResponse, err error) {
	reqVal := url.Values{}
	reqVal.Add("seller_id", seller_id)
	reqVal.Add("limit", "30")                        // hardcode
	reqVal.Add("status", "on_sale,trading,sold_out") // hardcode
	if pager_id != "" {
		reqVal.Add("max_pager_id", pager_id)
	}
	link := fmt.Sprintf("%s?%s", sellerProductParams.URL, reqVal.Encode())
	headers, err := generateHeader(link, sellerProductParams.Method)
	if err != nil {
		return
	}
	var sellerProfile SellerProfileResponse
	var err2 error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		_, err = a.Client.R().SetHeaders(headers).SetSuccessResult(&response).Get(link)
		wg.Done()
	}()

	go func() {
		sellerProfile, err2 = a.Profile(seller_id)
		wg.Done()
	}()
	wg.Wait()

	if err != nil {
		return
	}

	if response.Meta.HasNext {
		l := len(response.Data)
		response.Meta.PagerId = strconv.FormatInt(response.Data[l-1].PagerId, 10)
	}

	if err2 == nil {
		response.Profile = sellerProfile.Data
	}

	return
}

// ShopProducts is B2C  merchant products
func (a *APIFetcher) ShopProducts(shop_id, pager_id string) (response SellerProductsResponse, err error) {
	queryData, err := shopProductsGQL(shop_id, pager_id)
	if err != nil {
		return
	}

	headers, err := generateHeader(shopProductParams.URL, sellerProductParams.Method)
	if err != nil {
		return
	}

	resp, err := a.Client.R().SetHeaders(headers).SetBody(queryData).Post(shopProductParams.URL)
	if err != nil {
		return
	}

	return ShopProductsResult(resp.Bytes())
}

// Related for itemID products
func (a *APIFetcher) Related(item_id, limit string) (response RelatedResponse, err error) {
	reqVal := url.Values{}
	reqVal.Add("item_id", item_id)
	reqVal.Add("limit", limit) // 15
	link := fmt.Sprintf("%s?%s", relatedParams.URL, reqVal.Encode())

	headers, err := generateHeader(link, relatedParams.Method)
	if err != nil {
		return
	}

	_, err = a.Client.R().SetHeaders(headers).SetSuccessResult(&response).Get(link)
	if err != nil {
		return
	}
	return
}

// 目前只有 similar Looks
func similarParse(s SimilarData) (string, error) {
	so := SimilarItemOptions{}
	so.ItemID = s.ItemID
	so.ItemType = "ITEM_TYPE_MERCARI"
	so.ComponentType = "COMPONENT_TYPE_SIMILAR_LOOKS_ON_ITEM_THUMBNAIL"

	//@todo load more page
	slo := SimilarLooksOptions{
		PageSize:     120,
		ItemStatuses: []string{"ITEM_STATUS_ON_SALE"},
	}
	so.ComponentOptions = append(so.ComponentOptions, ComponentOption{SimilarLooksOptions: slo})
	res, err := json.Marshal(so)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// SimilarLooks
func (a *APIFetcher) SimilarLooks(s SimilarData) (response SimilarLooksResponse, err error) {
	queryData, err := similarParse(s)
	if err != nil {
		return
	}

	headers, err := generateHeader(similarLooksParams.URL, similarLooksParams.Method)
	if err != nil {
		return
	}

	data := SimilarLooksData{}
	_, err = a.Client.R().SetHeaders(headers).SetBody(queryData).SetSuccessResult(&data).Post(similarLooksParams.URL)
	if err != nil {
		return
	}

	response.Result = "OK"
	response.Meta.PageToke = data.LoadMoreToken
	items := []SimilarItem{}
	for _, v := range data.Contents {
		items = append(items, v.ItemContent.Item)
	}
	response.Items = items

	return response, nil
}

/**
* Profile get user profile
 */
func (a *APIFetcher) Profile(user_id string) (response SellerProfileResponse, err error) {
	reqVal := url.Values{}
	reqVal.Add("user_id", user_id)
	reqVal.Add("_user_format", "profile")
	link := fmt.Sprintf("%s?%s", profileParams.URL, reqVal.Encode())

	headers, err := generateHeader(link, profileParams.Method)
	if err != nil {
		return
	}
	_, err = a.Client.R().SetHeaders(headers).SetSuccessResult(&response).Get(link)
	if err != nil {
		return
	}

	if response.Data.Code == "" {
		response.Data.Code = strconv.Itoa(response.Data.ID)
	}

	return
}
