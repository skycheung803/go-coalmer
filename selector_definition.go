package coalmer

type searchSelector struct {
	Items []struct {
		ID       string `selector:"a > div" attr:"id" json:"id"`
		ItemType string `selector:"a > div" attr:"itemtype" json:"itemType"`
		Label    string `selector:"a > div" attr:"aria-label"`
		Name     string `selector:"a > div > span" json:"name"`
		Url      string `selector:"a" attr:"href" json:"url"`
		Image    string `selector:"picture > img" attr:"src" json:"image"`
		Price    string `selector:"a span.merPrice > span[class^='number']" json:"price"`
		Status   string `selector:"figure div[data-testid='thumbnail-sticker']" attr:"aria-label" json:"status"`
	} `selector:"#item-grid > ul > li" json:"items"`

	Meta struct {
		Pager struct {
			PrevPageUrl string `selector:"#search-result div[data-testid='pagination-prev-button'] > a" attr:"href" json:"prev_page_url"`
			NextPageUrl string `selector:"#search-result div[data-testid='pagination-next-button'] > a" attr:"href" json:"next_page_url"`
		}
	} `json:"meta"`
}

type detailProduct struct {
	ProductID string   `json:"productID"`
	Image     []string `json:"image"`
	//Image       string `json:"image"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Brand       string `json:"brand"`
	Offers      struct {
		URL           string `json:"url"`
		Availability  string `json:"availability"`
		Price         int    `json:"price"`
		PriceCurrency string `json:"priceCurrency"`
	} `json:"offers"`
}

type detailShopProduct struct {
	ProductID string `json:"productID"`
	//Image     []string `json:"image"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Brand       struct {
		Name string `json:"name"`
	} `json:"brand"`
	Offers struct {
		URL           string `json:"url"`
		Availability  string `json:"availability"`
		Price         int    `json:"price"`
		PriceCurrency string `json:"priceCurrency"`
	} `json:"offers"`
}

type relatedSelector struct {
	ID         string   `selector:"div.merItemThumbnail" attr:"id"`
	Name       string   `selector:"picture > img" attr:"alt"`
	Currency   string   `selector:"figure .merPrice span[class^='currency']"`
	Price      string   `selector:"figure .merPrice span[class^='number']"`
	Status     string   `selector:"figure div[data-testid='thumbnail-sticker']" attr:"aria-label"`
	Thumbnail  string   `selector:"picture > img" attr:"src"`
	Thumbnails []string `selector:"-" attr:"src"`
	ItemType   string   `selector:"div.merItemThumbnail" attr:"itemtype"`
	Label      string   `selector:"div.merItemThumbnail" attr:"aria-label"`
}

type detailSelector struct {
	Data struct {
		CurrentPrice struct {
			Currency string `selector:"span.currency"`
			PriceStr string `selector:"span.price"`
		} `selector:"#item-info div[data-testid='converted-currency-section']"`

		Info []struct {
			Title string `selector:"div[class^='title_']"`
			Body  string `selector:"div[class^='body_']"`
			Url   string `selector:"div[class^='body_'] a" attr:"href"`
		} `selector:"#item-info div.merDisplayRow"`

		Categories []struct {
			Id   int    `selector:"-" json:"id"`
			Url  string `selector:"a" attr:"href"`
			Name string `selector:"a" json:"name"`
		} `selector:"#item-info div[data-testid='item-detail-category'] .merTextLink"`

		SellerInfo struct {
			Label  string `selector:"~" attr:"aria-label"`
			Url    string `selector:"~" attr:"href"`
			Avatar string `selector:"picture > img" attr:"src"`
		} `selector:"#item-info a[data-location='item_details:seller_info']"`

		SellerBadges     []string `selector:"#item-info div[data-testid='seller-badge']"`
		ShippingPayerStr string   `selector:"#item-info > section:nth-child(1) > section:nth-child(2) p[class*='caption']"`
		UpdatedStr       string   `selector:"#item-info > section:nth-child(2) > p.merText"`

		Related []relatedSelector `selector:"#item-grid > ul > li"`
	} `selector:"#main" json:"data"`
	Product string `selector:"script[type='application/ld+json']"`
}

type detailShopSelector struct {
	Data struct {
		Photos []string `selector:"div.slick-slider div.slick-list img" attr:"src"`

		Info []struct {
			Title string `selector:"div[class^='title_']"`
			Body  string `selector:"div[class^='body_']"`
			Url   string `selector:"div[class^='body_'] a" attr:"href"`
		} `selector:"#product-info div.merDisplayRow"`

		Categories []struct {
			Id   int    `selector:"-" json:"id"`
			Url  string `selector:"a" attr:"href"`
			Name string `selector:"a" json:"name"`
		} `selector:"#product-info div[data-testid='product-detail-category'] .merTextLink"`

		SellerInfo struct {
			Label              string `selector:"picture > img" attr:"alt"`
			Url                string `selector:"~" attr:"href"`
			Avatar             string `selector:"picture > img" attr:"src"`
			NumRatingsStr      string `selector:"div.merRating span[class^='count']"`
			StarRatingScoreStr string `selector:"div.merRating" attr:"aria-label"`
		} `selector:"#product-info a[data-location='item_details:shop_info']"`

		SellerBadges     []string `selector:"#product-info div[data-testid='seller-badge']"`
		ShippingPayerStr string   `selector:"#product-info > section:nth-child(1) > section:nth-child(2) p[class*='caption']"`
		UpdatedStr       string   `selector:"#product-info > section:nth-child(2) > p.merText"`

		Related []relatedSelector `selector:"#item-grid > ul > li"`
	} `selector:"#main" json:"data"`
	Product string `selector:"script[type='application/ld+json']"`
}

type sellerSelector struct {
	Goods []struct {
		ID       string `selector:"a > div" attr:"id" json:"id"`
		ItemType string `selector:"a > div" attr:"itemType" json:"itemType"`
		Name     string `selector:"a > div > span" json:"name"`
		Url      string `selector:"a" attr:"href" json:"url"`
		Image    string `selector:"picture > img" attr:"src" json:"image"`
		Price    string `selector:"a span.merPrice > span[class^='number']" json:"price"`
	} `selector:"#item-grid > ul > li" json:"goods" `

	Meta struct {
		Seller struct {
			//ID                      int    `json:"id"`
			Name                    string `selector:"#avatarImage img" attr:"alt" json:"name"`
			PhotoURL                string `selector:"#avatarImage img" attr:"src" json:"photo_url"`
			RegisterSMSConfirmation string `json:"register_sms_confirmation"`
			Ratings                 struct {
				Good   int `selector:"-" json:"good"`
				Normal int `selector:"-" json:"normal"`
				Bad    int `selector:"-" json:"bad"`
			} `json:"ratings"`
			//NumRatings         int    `selector:"-" json:"num_ratings"`
			NumRatingsStr string `selector:"div.avatar-info div.merRating span[class^='count']"`
			//StarRatingScore    int    `selector:"-"  json:"star_rating_score"`
			StarRatingScoreStr string `selector:"div.avatar-info div.merRating" attr:"aria-label"`
			//NumSellItems       int    `selector:"-" json:"num_sell_items"`
			NumSellItemsStr string `selector:"div.user-info-supplement > a  span.merText"`
			//FollowerCount      int    `selector:"-" json:"follower_count"`
			FollowerCountStr string `selector:"div.user-info-supplement a:nth-child(2)  span.merText"`
			//FollowingCount     int    `selector:"-" json:"following_count"`
			FollowingCountStr string `selector:"div.user-info-supplement a:nth-child(3)  span.merText"`
			//Score             int    `selector:"-" json:"score"`
			Introduction string `selector:"div.merShowMore pre" json:"introduction"`
			//Created         int    `json:"created"`
		} `selector:"#main div[data-testid='profile-info']" json:"seller"`

		Pager struct {
			PrevPageUrl string `selector:"#search-result div[data-testid='pagination-prev-button'] > a" attr:"href" json:"prev_page_url"`
			NextPageUrl string `selector:"#search-result div[data-testid='pagination-next-button'] > a" attr:"href" json:"next_page_url"`
		}
	} `json:"meta"`
}

type sellerShopSelector struct {
	Goods []struct {
		ID       string `selector:"-"`
		ItemType string `selector:"a > div" attr:"itemType"`
		Name     string `selector:"img" attr:"alt"`
		Url      string `selector:"a" attr:"href"`
		Image    string `selector:"img" attr:"src"`
		Price    string `selector:"a p.chakra-text"`
		Status   string `selector:"a div[data-testid='soldout-label']" attr:"class"`
	} `selector:"div.css-gaste1 > div > div > div > div"`

	Meta struct {
		Seller struct {
			//ID                      int    `json:"id"`
			Name     string `selector:"p[data-testid='shop-name']"`
			PhotoURL string `selector:"img" attr:"src"`
			//RegisterSMSConfirmation string `json:"register_sms_confirmation"`
			Ratings struct {
				Good   int `selector:"-"`
				Normal int `selector:"-"`
				Bad    int `selector:"-"`
			} `json:"ratings"`
			//NumRatings         int      `selector:"-" json:"num_ratings"`
			NumRatingsStr string `selector:"p[data-testid='rating-count']"`
			//StarRatingScore    int      `selector:"-"  json:"star_rating_score"`
			StarRatingScoreStr []string `selector:"div[data-testid='star-rating'] > div.chakra-stack svg" `
			//NumSellItems       int      `selector:"-" json:"num_sell_items"`
			NumSellItemsStr string `selector:"div.user-info-supplement > a  span.merText"`
			//FollowerCount      int      `selector:"-" json:"follower_count"`
			FollowerCountStr string `selector:"p[data-testid='follows-count'] > span:nth-child(2)"`
			//FollowingCount     int      `selector:"-" json:"following_count"`
			FollowingCountStr string `selector:"div.user-info-supplement a:nth-child(3)  span.merText"`
			//Score              int      `selector:"-" json:"score"`
			//Introduction       string   `selector:"div.merShowMore pre" json:"introduction"`
		} `selector:"#__next li[data-testid='list-item']" `

		Pager struct {
			PrevPageUrl string `selector:"#search-result div[data-testid='pagination-prev-button'] > a" attr:"href"`
			NextPageUrl string `selector:"#search-result div[data-testid='pagination-next-button'] > a" attr:"href"`
		}
	} `json:"meta"`

	Props string `selector:"#__NEXT_DATA__"`
}
