package coalmer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

// ProductDetailResult mapper the B2C merchant product struct to C2C merchant product struct
func ProductDetailResult(payload []byte) (detail ItemResultResponse, err error) {
	result, err := simplejson.NewJson(payload)
	if err != nil {
		return
	}

	data := MercariDetail{}
	data.ProductId = result.Get("name").MustString()
	data.Url = fmt.Sprintf("%s%s", webShopsItemURL, data.ProductId)
	data.ProductName = result.Get("displayName").MustString()
	data.Price, _ = strconv.Atoi(result.Get("price").MustString())

	createTime, _ := time.Parse(time.RFC3339, result.Get("createTime").MustString())
	updateTime, _ := time.Parse(time.RFC3339, result.Get("updateTime").MustString())
	data.Created = int64(createTime.Unix())
	data.Updated = int64(updateTime.Unix())

	productTags := result.Get("productTags").MustStringArray()
	if len(productTags) > 0 {
		data.Status = productTags[0]
	} else {
		data.Status = "on_sale"
	}
	productDetail := result.Get("productDetail")
	data.Description = productDetail.Get("description").MustString()
	data.Condition.Name = productDetail.Get("condition").Get("displayName").MustString()
	data.Brand.Id, _ = strconv.ParseInt(productDetail.Get("brand").Get("brandId").MustString(), 10, 64)
	data.Brand.Name = productDetail.Get("brand").Get("displayName").MustString()

	data.ShippingFrom.Name = productDetail.Get("shippingFromArea").Get("displayName").MustString()
	data.ShippingMethod.Id, _ = strconv.ParseInt(productDetail.Get("shippingMethod").Get("shippingMethodId").MustString(), 10, 64)
	data.ShippingMethod.Name = productDetail.Get("shippingMethod").Get("displayName").MustString()

	data.ShippingDuration.Id, _ = strconv.ParseInt(productDetail.Get("shippingDuration").Get("shippingDurationId").MustString(), 10, 64)
	data.ShippingDuration.Name = productDetail.Get("shippingDuration").Get("displayName").MustString()

	shop := productDetail.Get("shop")
	data.Seller.Code = shop.Get("name").MustString()
	data.Seller.Name = shop.Get("displayName").MustString()
	data.Seller.Avatar = shop.Get("thumbnail").MustString()
	data.Seller.Rating = shop.Get("shopStats").Get("score").MustFloat64()
	data.Seller.Score, _ = strconv.Atoi(shop.Get("shopStats").Get("reviewCount").MustString())
	data.Seller.NumSell = int32(len(shop.Get("shopItems").MustArray()))

	categories := []Name_Id_Unit{}
	for _, v := range productDetail.Get("categories").MustArray() {
		vMap, ok := v.(map[string]interface{})
		if ok {
			id, _ := strconv.Atoi(vMap["categoryId"].(string))
			categories = append(categories, Name_Id_Unit{
				Id:   int64(id),
				Name: vMap["displayName"].(string),
			})
		}

	}
	data.Categories = categories
	data.Photos = productDetail.Get("photos").MustStringArray()

	data.ShippingPayer.Id, _ = strconv.ParseInt(productDetail.Get("shippingPayer").Get("shippingPayerId").MustString(), 10, 64)
	data.ShippingPayer.Name = productDetail.Get("shippingPayer").Get("displayName").MustString()
	data.ShippingPayer.Code = strings.ToLower(productDetail.Get("shippingPayer").Get("code").MustString())

	shippingFeeConfig := productDetail.Get("shippingFeeConfig")
	// 检查 shippingFeeConfig 是否为 nil
	if shippingFeeConfig.Interface() != nil {
		data.ShippingPayer.MinFee = shippingFeeConfig.Get("minFeePrice").MustInt()
		data.ShippingPayer.MaxFee = shippingFeeConfig.Get("maxFeePrice").MustInt()
		fees := shippingFeeConfig.Get("fees").MustArray()
		for _, fee := range fees {
			feeMap := fee.(map[string]interface{})
			// 检查 displayName 是否为 "大阪"
			if feeMap["displayName"] == "大阪" {
				data.ShippingPayer.Fee = feeMap["price"].(int)
				break
			}
		}
	}

	detail.Result = "OK"
	detail.Data = data
	return detail, nil
}

// ShopProductsResult mapper the B2C merchant products struct to C2C merchant products struct
func ShopProductsResult(payload []byte) (result SellerProductsResponse, err error) {
	info, err := simplejson.NewJson(payload)
	if err != nil {
		return
	}

	products := info.Get("data").Get("products")
	pageInfo := products.Get("pageInfo")

	result.Meta.HasNext = pageInfo.Get("hasNextPage").MustBool()
	result.Meta.PagerId = pageInfo.Get("endCursor").MustString()

	items := []SellerItem{}
	for _, item := range products.Get("edges").MustArray() {
		node, ok := item.(map[string]interface{})["node"].(map[string]interface{})
		if ok {
			priceValue, _ := node["price"].(json.Number)
			priceInt, _ := priceValue.Int64()
			item := SellerItem{
				RelatedItem: RelatedItem{
					ID:    node["id"].(string),
					Name:  node["name"].(string),
					Price: int(priceInt),
				},
			}

			inStock, ok := node["inStock"].(bool)
			if !ok {
				item.Status = "sold_out"
			} else {
				if inStock {
					item.Status = "on_sale"
				} else {
					item.Status = "sold_out"
				}
			}

			thumbnails := make([]string, 0)
			for _, image := range node["assets"].([]interface{}) {
				imageMap, ok := image.(map[string]interface{})
				if ok {
					thumbnails = append(thumbnails, imageMap["imageUrl"].(string))
				}
			}
			item.Thumbnails = thumbnails
			items = append(items, item)
		}
	}

	result.Data = items
	//result.Meta = meta
	result.Result = "OK"
	return result, nil
}
