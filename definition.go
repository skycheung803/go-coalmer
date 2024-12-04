package coalmer

import "net/http"

const (
	ApiURL  = "https://api.mercari.jp/"
	RootURL = "https://jp.mercari.com/"
	ShopURL = "https://mercari-shops.com/"
)

const DefaultLengthSearchSessionId = 32

const (
	SearchOptionSortScore       = "SORT_SCORE"
	SearchOptionSortPrice       = "SORT_PRICE"
	SearchOptionSortNumLikes    = "SORT_NUM_LIKES"
	SearchOptionSortCreatedTime = "SORT_CREATED_TIME"

	SearchOptionOrderDESC = "ORDER_DESC"
	SearchOptionOrderASC  = "ORDER_ASC"
)

var (
	webSearchURL      = RootURL + "search"
	webItemURL        = RootURL + "item"
	webShopsItemURL   = RootURL + "shops/product/"
	webSellerURL      = RootURL + "user/profile/"
	webShopsSellerURL = ShopURL + "shops/"
)

var searchParams = struct {
	URL    string
	Method string
}{
	URL:    ApiURL + "v2/entities:search",
	Method: http.MethodPost,
}

var itemParams = struct {
	URL    string
	Method string
}{
	URL:    ApiURL + "items/get",
	Method: http.MethodGet,
}

var shopItemParams = struct {
	URL    string
	Method string
}{
	URL:    ApiURL + "v1/marketplaces/shops/products",
	Method: http.MethodGet,
}

var relatedParams = struct {
	URL    string
	Method string
}{
	URL:    ApiURL + "items/related_items",
	Method: http.MethodGet,
}

var profileParams = struct {
	URL    string
	Method string
}{
	URL:    ApiURL + "users/get_profile",
	Method: http.MethodGet,
}

// similarLooks items
var similarLooksParams = struct {
	URL    string
	Method string
}{
	URL:    ApiURL + "v2/relateditems/component",
	Method: http.MethodPost,
}

var sellerProductParams = struct {
	URL    string
	Method string
}{
	URL:    ApiURL + "items/get_items",
	Method: http.MethodGet,
}

var shopProductParams = struct {
	URL    string
	Method string
}{
	URL:    ShopURL + "graphql",
	Method: http.MethodPost,
}

type xerror struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SearchData struct {
	Keyword           string
	BrandId           []int
	CategoryId        []int
	ColorId           []int
	ConditionId       []int
	PriceMax          int
	PriceMin          int
	Order             string
	Sort              string
	Page              int
	Limit             int
	SearchConditionId string
}

// entities:search request body
type V2Search struct {
	DefaultDatabases   []string              `json:"defaultDatabases"` // have default value
	IndexRouting       string                `json:"indexRouting"`     // have default value
	PageSize           int                   `json:"pageSize"`
	PageToken          string                `json:"pageToken,omitempty"` // pagination, stucture like "v1:1" "v1:2", not appear in first page
	SearchCondition    V2SearchRequestDetail `json:"searchCondition"`
	SearchSessionId    string                `json:"searchSessionId"`
	ServiceFrom        string                `json:"serviceFrom"`    // have default value
	ThumbnailTypes     []int                 `json:"thumbnailTypes"` // default empty
	UserId             string                `json:"userId"`
	WithItemBrand      bool                  `json:"withItemBrand"`
	WithItemPromotions bool                  `json:"withItemPromotions"`
	WithItemSize       bool                  `json:"withItemSize"`
	WithItemSizes      bool                  `json:"withItemSizes"`
	WithShopName       bool                  `json:"withShopName"`
}

// searchCondition part of request body
type V2SearchRequestDetail struct {
	Attributes       []any    `json:"attributes"`      // default empty
	BrandId          []int    `json:"brandId"`         // default empty
	CategoryId       []int    `json:"categoryId"`      // default empty
	ColorId          []int    `json:"colorId"`         // default empty
	HasCoupon        bool     `json:"hasCoupon"`       // default false
	ItemConditionId  []int    `json:"itemConditionId"` // default empty
	ItemTypes        []any    `json:"itemTypes"`       // default empty
	Keyword          string   `json:"keyword"`
	ExcludeKeyword   string   `json:"excludeKeyword"`   // TODO: check if this can achieve what it means
	PriceMax         int      `json:"priceMax"`         // default empty
	PriceMin         int      `json:"priceMin"`         // default empty
	SellerId         []string `json:"sellerId"`         // default empty
	ShippingFromArea []any    `json:"shippingFromArea"` // default empty
	ShippingMethod   []any    `json:"shippingMethod"`   // default empty
	ShippingPayerId  []any    `json:"shippingPayerId"`  // default empty
	SizeId           []any    `json:"sizeId"`           // default empty
	SKUIds           []any    `json:"skuIds"`           // default empty
	Order            string   `json:"order"`
	Sort             string   `json:"sort"`
	Status           []string `json:"status"` // default empty
}

// entities:search response body
type SearchResponse struct {
	Result            string                 `json:"result"` // OK or error
	Errors            []xerror               `json:"errors,omitempty"`
	Items             []Item                 `json:"items"`
	Components        []any                  `json:"components,omitempty"`
	SearchConditionId string                 `json:"searchConditionId,omitempty"`
	Meta              map[string]interface{} `json:"meta"`
}

// Search() result item that function return
type Item struct {
	ProductId     string   `json:"id"`
	ProductName   string   `json:"name"`
	Price         string   `json:"price"`
	Thumbnails    []string `json:"thumbnails"`
	ItemType      string   `json:"itemType"` // ITEM_TYPE_MERCARI ITEM_TYPE_BEYOND
	Condition     string   `json:"itemConditionId"`
	ShippingPayer string   `json:"shippingPayerId,omitempty"` // 0(or 2): by seller
	Status        string   `json:"status"`
	Seller        string   `json:"sellerId,omitempty"`
	Buyer         string   `json:"buyerId,omitempty"`
	Created       string   `json:"created"`
	Updated       string   `json:"updated"`
}

type Name_Id_Unit struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// Item() response body
type ItemResultResponse struct {
	Result string                 `json:"result"`
	Errors []xerror               `json:"errors,omitempty"`
	Data   MercariDetail          `json:"data"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}

// Item() response body item part
type MercariDetail struct {
	Url                string         `json:"url"`
	ProductId          string         `json:"id"`
	ProductName        string         `json:"name"`
	Price              int            `json:"price"`
	Seller             ItemSellerInfo `json:"seller"`
	Status             string         `json:"status"`
	Description        string         `json:"description"`
	Condition          Name_Id_Unit   `json:"item_condition"`
	Brand              Name_Id_Unit   `json:"item_brand"`
	Like               int            `json:"num_likes"`
	Comment            int            `json:"num_comments"`
	Photos             []string       `json:"photos"`
	AnonymousShipping  bool           `json:"is_anonymous_shipping"`
	ShippingDuration   Name_Id_Unit   `json:"shipping_duration"`
	ShippingFrom       Name_Id_Unit   `json:"shipping_from_area"`
	ShippingMethod     Name_Id_Unit   `json:"shipping_method"`
	ShippingPayer      Name_Id_Unit   `json:"shipping_payer"`
	ItemSize           Name_Id_Unit   `json:"item_size"`
	Colors             []Name_Id_Unit `json:"colors"`
	ItemCategory       Name_Id_Unit   `json:"item_category"`
	ItemCategoryNtiers Name_Id_Unit   `json:"item_category_ntiers,omitempty"`
	ParentCategories   []Name_Id_Unit `json:"parent_categories_ntiers,omitempty"`
	Categories         []Name_Id_Unit `json:"categories"`
	Created            int64          `json:"created"`
	Updated            int64          `json:"updated"`
	UpdatedStr         string         `json:"UpdatedStr"`
	Related            []RelatedItem  `json:"related,omitempty"`
	SimilarLooks       []SimilarItem  `json:"similar_looks,omitempty"`
}

type ItemSellerInfo struct {
	ID           int64            `json:"id"`
	Code         string           `json:"code"`
	Name         string           `json:"name"`
	QuickShipper bool             `json:"quick_shipper"`
	NumSell      int32            `json:"num_sell_items"`
	Avatar       string           `json:"photo_thumbnail_url"`
	Created      int64            `json:"created"`
	SmsAuth      string           `json:"register_sms_confirmation"`
	SmsAuthAt    string           `json:"register_sms_confirmation_at"`
	Score        int              `json:"score"` // =good-bad
	NumRatings   int              `json:"num_ratings"`
	Rating       int              `json:"star_rating_score"`
	Ratings      ItemSellerRating `json:"ratings"`
}

type ItemSellerRating struct {
	Good   int32
	Bad    int32
	Normal int32
}

type RelatedItem struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Price      int      `json:"price"`
	Status     string   `json:"status"`
	Thumbnail  string   `json:"thumbnail"`
	Thumbnails []string `json:"thumbnails"`
	ItemType   string   `json:"item_type"`
}

// related item response body
type RelatedResponse struct {
	Result string                 `json:"result"`
	Meta   map[string]interface{} `json:"meta"`
	Data   struct {
		Items []RelatedItem `json:"items"`
	} `json:"data"`
}

type SellerItem struct {
	RelatedItem
	PagerId int64 `json:"pager_id"`
	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

type Ratings struct {
	Good   int `json:"good"`
	Normal int `json:"normal"`
	Bad    int `json:"bad"`
}

type ProfileData struct {
	ID                       int     `json:"id"`
	Code                     string  `json:"code"`
	Name                     string  `json:"name"`
	PhotoURL                 string  `json:"photo_url"`
	PhotoThumbnailURL        string  `json:"photo_thumbnail_url"`
	RegisterSMSConfirmation  string  `json:"register_sms_confirmation"`
	Ratings                  Ratings `json:"ratings"`
	PolarizedRatings         Ratings `json:"polarized_ratings"`
	NumRatings               int     `json:"num_ratings"`
	StarRatingScore          int     `json:"star_rating_score"`
	IsFollowable             bool    `json:"is_followable"`
	IsFollowingRequester     bool    `json:"is_following_requester"`
	IsBlocked                bool    `json:"is_blocked"`
	FollowingCount           int     `json:"following_count"`
	FollowerCount            int     `json:"follower_count"`
	KYCType                  string  `json:"kyc_type"`
	Score                    int     `json:"score"`
	Created                  int64   `json:"created"`
	Proper                   bool    `json:"proper"`
	Email                    string  `json:"email"`
	PhoneNumber              string  `json:"phone_number"`
	Introduction             string  `json:"introduction"`
	IVCode                   string  `json:"iv_code"`
	IsOfficial               bool    `json:"is_official"`
	NumSellItems             int     `json:"num_sell_items"`
	NumTicket                int     `json:"num_ticket"`
	BounceMailFlag           string  `json:"bounce_mail_flag"`
	IsFollowing              bool    `json:"is_following"`
	CurrentPoint             int     `json:"current_point"`
	CurrentSales             int     `json:"current_sales"`
	IsOrganizationalUser     bool    `json:"is_organizational_user"`
	OrganizationalUserStatus string  `json:"organizational_user_status"`
}

type SellerProfileResponse struct {
	Result string      `json:"result"`
	Data   ProfileData `json:"data"`
	Meta   struct{}    `json:"meta"`
}

// seller product item response body
type SellerProductsResponse struct {
	Result string `json:"result"`
	Meta   struct {
		HasNext bool   `json:"has_next"`
		PagerId string `json:"pager_id"`
	} `json:"meta"`
	Data    []SellerItem `json:"data"`
	Profile ProfileData  `json:"profile"`
}

type SimilarData struct {
	ItemID    string
	PageToken string
	Limit     int
}

type ComponentOption struct {
	SimilarLooksOptions SimilarLooksOptions `json:"similarLooksOptions"`
}

type SimilarLooksOptions struct {
	PageSize     int      `json:"pageSize"`
	ItemStatuses []string `json:"itemStatuses"`
}

type SimilarItemOptions struct {
	ItemID           string            `json:"itemId"`
	ItemType         string            `json:"itemType"`
	ComponentType    string            `json:"componentType"`
	ComponentOptions []ComponentOption `json:"componentOptions"`
}

type SimilarItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	Status    string `json:"status"`
	Thumbnail string `json:"thumbnail"`
	Type      string `json:"type"`
}

type ItemContent struct {
	Item SimilarItem `json:"item"`
}

type Content struct {
	Index       int         `json:"index"`
	ItemContent ItemContent `json:"itemContent"`
}

type SimilarLooksData struct {
	Index         int    `json:"index"`
	ComponentType string `json:"componentType"`
	DataType      string `json:"dataType"`
	Header        struct {
		Title string `json:"title"`
	} `json:"header"`
	Contents      []Content `json:"contents"`
	LoadMoreToken string    `json:"loadMoreToken"`
}

type SimilarLooksResponse struct {
	Result string `json:"result"` // OK or error
	Meta   struct {
		PageToke string `json:"page_token"`
	} `json:"meta"`
	//Components any             `json:"components"`
	Items []SimilarItem `json:"items"`
}

/*

{
    "result": "error",
    "errors": [
        {
            "code": "InvisibleItemException",
            "message": "該当する商品は削除されています。"
        }
    ],
    "meta": {}
}
*/
