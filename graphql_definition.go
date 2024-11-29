package coalmer

import (
	"encoding/json"
)

func shopProductsGQL(shop_id, pager_id string) (gql string, err error) {

	type variables struct {
		ShopId  string `json:"shopId"`
		Cursor  string `json:"cursor,omitempty"` // pagination, not appear in first page
		InStock bool   `json:"inStock,omitempty"`
		OrderBy string `json:"orderBy,omitempty"`
	}

	var shopProducts = struct {
		OperationName string    `json:"operationName"`
		Variables     variables `json:"variables"`
		Query         string    `json:"query"`
	}{
		OperationName: "ShopProducts",
		Variables: variables{
			ShopId: shop_id,
			Cursor: pager_id,
		},
		Query: `query ShopProducts($shopId: String!, $inStock: Boolean, $inDualPrice: Boolean, $orderBy: OrderBy, $cursor: String) {
  products(
    shopId: $shopId
    inStock: $inStock
    inDualPrice: $inDualPrice
    orderBy: $orderBy
    after: $cursor
    first: 90
    status: [STATUS_OPENED]
  ) {
    pageInfo {
      ...PageInfo
      __typename
    }
    edges {
      node {
        ...ProductNode
        __typename
      }
      __typename
    }
    __typename
  }
  }
fragment PageInfo on PageInfo {
  hasNextPage
  endCursor
  __typename
  }
fragment ProductNode on Product {
  id
  price
  name
  status {
    id
    name
    __typename
  }
  inStock
  assets(count: 1) {
    id
    imageUrl(options: {presets: [Small]})
    __typename
  }
  dualPrice {
    dualPriceType
    basePrice
    __typename
  }
  __typename
  }`,
	}

	jsonData, err := json.Marshal(shopProducts)
	if err != nil {
		return
	}
	return string(jsonData), nil
}
