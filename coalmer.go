package coalmer

import "github.com/go-rod/rod"

// DataFetcher interface defines methods for retrieving data from a data source.
type DataFetcher interface {
	// Search searches for items based on the provided search data.
	Search(s SearchData) (SearchResponse, error)

	// Detail retrieves detailed information about a specific item.
	Detail(itemId string) (ItemResultResponse, error)

	// Seller retrieves a list of products from a specific seller.
	Seller(sellerId string, pager_id string) (SellerProductsResponse, error)
}

type mode string

const (
	ModeWeb mode = "WEB"
	ModeAPI mode = "API"
)

type Coalmer struct {
	Mode    mode
	Debug   bool
	Fetcher DataFetcher
	Browser *rod.Browser
}

type CoalmerOption func(*Coalmer)

// WithBrowser
func WithBrowser(browser *rod.Browser) CoalmerOption {
	return func(c *Coalmer) {
		c.Browser = browser
	}
}

// WithBrowserMode
func WithBrowserMode() CoalmerOption {
	return func(c *Coalmer) {
		c.Mode = ModeWeb
	}
}

// WithDebug
func WithDebug(debug bool) CoalmerOption {
	return func(c *Coalmer) {
		c.Debug = debug
	}
}

/*
NewCoalmer creates a new Coalmer instance with the given mode and debug flag.
The mode can be either "web" or "api", and the debug flag determines whether to show the browser or not .
default is "api" and debug is false.
*/

func NewCoalmer(options ...CoalmerOption) *Coalmer {
	c := &Coalmer{
		Mode:  ModeAPI,
		Debug: false,
	}

	for _, f := range options {
		f(c)
	}

	var fetcher DataFetcher
	if c.Mode == ModeWeb {
		//fetcher = NewWebFetcher(!c.Debug) //debug no headless show browser
		fetcher = NewWebFetcher(WithClient(c.Browser), WithHeadless(!c.Debug))
	} else {
		fetcher = NewAPIFetcher(c.Debug)
	}

	c.Fetcher = fetcher
	return c
}
