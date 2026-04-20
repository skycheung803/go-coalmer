package coalmer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"slices"

	"github.com/imroc/req/v3"
)

type APIFetcher struct {
	Client *req.Client
}

func NewAPIFetcher(debug bool) *APIFetcher {
	Client := req.C()
	Client.ImpersonateChrome()
	if debug {
		Client.DevMode()
	}

	return &APIFetcher{
		Client: Client,
	}
}
func SearchConditionParse(p SearchData) (V2SearchRequestDetail, error) {
	sp := V2SearchRequestDetail{}
	sp.HasCoupon = false
	sp.Keyword = p.Keyword

	if p.PriceMin > p.PriceMax {
		p.PriceMin, p.PriceMax = p.PriceMax, p.PriceMin
	}

	if p.PriceMin > 0 {
		sp.PriceMin = p.PriceMin
	}
	if p.PriceMax > 0 {
		sp.PriceMax = p.PriceMax
	}

	if len(p.CategoryId) > 0 {
		sp.CategoryId = p.CategoryId
	} else {
		sp.CategoryId = []int{}
	}

	if len(p.ConditionId) > 0 {
		sp.ItemConditionId = p.ConditionId
	} else {
		sp.ItemConditionId = []int{}
	}

	if len(p.ColorId) > 0 {
		sp.ColorId = p.ColorId
	} else {
		sp.ColorId = []int{}
	}

	if p.Sort != "" && p.Order != "" {
		sp.Sort = strings.ToUpper("SORT_" + p.Sort)
		sp.Order = strings.ToUpper("ORDER_" + p.Order)
	} else {
		sp.Sort = SearchOptionSortScore
		sp.Order = SearchOptionOrderDESC
	}

	if len(p.Status) > 0 {
		if ok := Contains(p.Status, "on_sale"); ok {
			sp.Status = append(sp.Status, "STATUS_ON_SALE")
		}

		if ok := Contains(p.Status, "sold_out"); ok {
			sp.Status = append(sp.Status, "STATUS_SOLD_OUT", "STATUS_TRADING")
		}
	}

	if len(p.ItemTypes) > 0 {
		if ok := Contains(p.ItemTypes, "mercari"); ok {
			sp.ItemTypes = append(sp.ItemTypes, "ITEM_TYPE_MERCARI")
		}

		if ok := Contains(p.ItemTypes, "beyond"); ok {
			sp.ItemTypes = append(sp.ItemTypes, "ITEM_TYPE_BEYOND")
		}
	}

	if sp.ItemTypes == nil {
		sp.ItemTypes = []string{}
	}
	if sp.SellerId == nil {
		sp.SellerId = []string{}
	}

	if sp.BrandId == nil {
		sp.BrandId = []int{}
	}

	if sp.SizeId == nil {
		sp.SizeId = []any{}
	}

	if sp.SizeId == nil {
		sp.SizeId = []any{}
	}

	if sp.SKUIds == nil {
		sp.SKUIds = []any{}
	}

	if sp.ShippingFromArea == nil {
		sp.ShippingFromArea = []any{}
	}

	if sp.ShippingMethod == nil {
		sp.ShippingMethod = []any{}
	}

	if sp.ShippingPayerId == nil {
		sp.ShippingPayerId = []any{}
	}

	return sp, nil
}

func apiSearchParse(p SearchData) (string, error) {
	sp := V2Search{}
	//sp.DefaultDatabases = []string{"DATASET_TYPE_MERCARI", "DATASET_TYPE_BEYOND"}
	sp.ServiceFrom = "suruga"
	sp.WithItemBrand = true
	sp.WithItemPromotions = true
	sp.WithItemSize = false
	sp.WithItemSizes = true
	sp.WithShopName = false
	sp.IndexRouting = "INDEX_ROUTING_UNSPECIFIED"

	sp.SearchCondition, _ = SearchConditionParse(p)
	sp.SearchSessionId = generateSearchSessionId(DefaultLengthSearchSessionId)

	if p.Limit > 0 {
		sp.PageSize = p.Limit
	} else {
		sp.PageSize = DefaultPageSize
	}

	if p.Page > 0 {
		sp.PageToken = fmt.Sprintf("v1:%d", p.Page)
	}

	res, err := json.Marshal(sp)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (a *APIFetcher) Index(limit int) (response IndexProductsResponse, err error) {
	reqVal := url.Values{}
	reqVal.Add("limit", strconv.Itoa(limit))
	reqVal.Add("type", "category")
	link := fmt.Sprintf("%s?%s", indexParams.URL, reqVal.Encode())

	headers, err := generateHeader(link, indexParams.Method)
	if err != nil {
		return
	}
	_, err = a.Client.R().SetHeaders(headers).SetSuccessResult(&response).Get(link)
	if err != nil {
		return
	}

	if response.Meta.HasNext {
		l := len(response.Data)
		response.Meta.PagerId = strconv.FormatInt(response.Data[l-1].PagerId, 10)
	}
	return response, nil
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

	var errMsg merror
	resp, err := a.Client.R().SetHeaders(headers).SetBody(queryData).SetSuccessResult(&response).SetErrorResult(&errMsg).Post(searchParams.URL)
	if err != nil {
		return
	}

	if resp.IsErrorState() {
		response.Result = "error"
		err = fmt.Errorf("mercari search api error: %s", errMsg.Message)
		return
	}

	for k, v := range response.Items {
		response.Items[k].IsNoPrice2 = v.IsNoPrice
	}

	response.Result = "OK"
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
	reqVal.Add("include_auction", "true")
	link := fmt.Sprintf("%s?%s", itemParams.URL, reqVal.Encode())

	headers, err := generateHeader(link, itemParams.Method)
	if err != nil {
		return
	}

	var similarLooksResponse SimilarLooksResponse
	var err2 error
	var errMsg merror
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		_, err = a.Client.R().SetHeaders(headers).SetSuccessResult(&response).SetErrorResult(&errMsg).Get(link)
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

	if errMsg.Code != 0 {
		response.Result = "error"
		err = fmt.Errorf("mercari Item api error: %s", errMsg.Message)
		return
	}

	if response.Data.ProductId == "" {
		response.Result = "error"
		err = errors.New("item not found")
		return
	}

	response.Data.Url = fmt.Sprintf("%s/%s", webItemURL, id)
	response.Data.Categories = response.Data.ParentCategories
	response.Data.Categories = append(response.Data.Categories, response.Data.ItemCategoryNtiers)

	if err2 == nil {
		response.Data.SimilarLooks = similarLooksResponse.Items
	}

	if slices.Contains(response.Data.Hashtags, "価格がつけられないもの") {
		response.Data.IsNoPrice = true
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

	var errMsg merror
	resp, err := a.Client.R().SetHeaders(headers).SetErrorResult(&errMsg).Get(link)
	if err != nil {
		return
	}

	if errMsg.Code != 0 {
		response.Result = "error"
		err = fmt.Errorf("mercari ShopItem api error: %s", errMsg.Message)
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
	var errMsg merror
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		_, err = a.Client.R().SetHeaders(headers).SetSuccessResult(&response).SetErrorResult(&errMsg).Get(link)
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

	if errMsg.Code != 0 {
		response.Result = "error"
		err = fmt.Errorf("mercari SellerProducts api error: %s", errMsg.Message)
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
		response.Result = "error"
		return
	}

	headers, err := generateHeader(shopProductParams.URL, sellerProductParams.Method)
	if err != nil {
		response.Result = "error"
		return
	}

	var errMsg merror
	resp, err := a.Client.R().SetHeaders(headers).SetBody(queryData).SetErrorResult(&errMsg).Post(shopProductParams.URL)
	if err != nil {
		response.Result = "error"
		return
	}

	if errMsg.Code != 0 {
		response.Result = "error"
		err = fmt.Errorf("mercari ShopProducts api error: %s", errMsg.Message)
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

func ImageSearchParse(params SearchData) (string, error) {

	if params.ImageUri == "" {
		return "", fmt.Errorf("image_uri is required")
	}

	searchCondition, _ := SearchConditionParse(params)
	searchCondition.Sort = SearchOptionSortSimilarity
	searchCondition.Order = SearchOptionOrderDESC

	ImageSearchCondition := ImageSearchCondition{
		ImageUri:        params.ImageUri,
		SearchCondition: searchCondition,
	}
	pageToken := ""
	if params.Page > 0 {
		pageToken = fmt.Sprintf("v1:%d", params.Page)
	}

	searchData := ImageSearchData{
		Config: ImageSearchConfig{
			ResponseToggles: []string{"WITH_FILTERING", "WITH_CATEGORY_FACETS_SUGGEST"},
		},
		ImageSearchCondition: ImageSearchCondition,
		PageSize:             30,
		PageToken:            pageToken,
		SearchSessionId:      generateSearchSessionId(DefaultLengthSearchSessionId),
	}

	res, err := json.Marshal(searchData)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (a *APIFetcher) SearchByImage(params SearchData) (response ImageSearchResponse, err error) {
	queryData, err := ImageSearchParse(params)
	if err != nil {
		return
	}

	headers, err := generateHeader(imageSearchParams.URL, imageSearchParams.Method)
	if err != nil {
		return
	}

	var errMsg merror
	resp, err := a.Client.R().SetHeaders(headers).SetBody(queryData).SetSuccessResult(&response).SetErrorResult(&errMsg).Post(imageSearchParams.URL)
	if err != nil {
		return
	}

	if resp.IsErrorState() {
		response.Result = "error"
		err = fmt.Errorf("mercari image search api error: %s", errMsg.Message)
		return
	}

	response.Result = "OK"
	return
}
