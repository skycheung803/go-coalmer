package coalmer

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
}

/*
NewCoalmer creates a new Coalmer instance with the given mode and debug flag.
The mode can be either "web" or "api", and the debug flag determines whether to show the browser or not .
default is "api" and debug is false.
*/

func NewCoalmer(options ...func(*Coalmer)) *Coalmer {
	c := &Coalmer{
		Mode:    ModeAPI,
		Debug:   false,
		Fetcher: nil,
	}

	for _, f := range options {
		f(c)
	}

	var fetcher DataFetcher
	if c.Mode == ModeWeb {
		fetcher = NewWebFetcher(!c.Debug) //debug no headless show browser
	} else {
		fetcher = NewAPIFetcher(c.Debug)
	}

	c.Fetcher = fetcher
	return c
}

// WithBrowserMode
func WithBrowserMode() func(*Coalmer) {
	return func(c *Coalmer) {
		c.Mode = ModeWeb
	}
}

// WithDebug
func WithDebug(debug bool) func(*Coalmer) {
	return func(c *Coalmer) {
		c.Debug = debug
	}
}
