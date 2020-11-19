package endpointsdynamic

/*
This defines a basic set of data for our Star Wars Schema.

This data is hard coded for the sake of the demo, but you could imagine
fetching this data from a backend service rather than from hardcoded
JSON objects in a more complex demo.
*/

//Endpoint ...
type Endpoint struct {
	ClientID   string `json:"clientId"`
	SiteID     string `json:"siteId"`
	EndpointID string `json:"endpointId"`
}

//Site ...
type Site struct {
	ClientID  string     `json:"clientId"`
	SiteID    string     `json:"siteId"`
	Code      string     `json:"siteCode"`
	Name      string     `json:"siteName"`
	Endpoints []Endpoint `json:"endpoints"`
}

var E1 = &Endpoint{"6500023", "6500023", "37866524-cc91-4d64-b5db-b912eaf4339e"}
var E2 = &Endpoint{"6500023", "6500023", "37866524-cc91-4d64-b5db-b912eaf4339u"}
var E3 = &Endpoint{"6500023", "6500023", "37866524-cc91-4d64-b5db-b912eaf4338u"}

var S1 = &Site{
	"6500023",
	"6500023",
	"test",
	"test-name",
	[]Endpoint{*E1, *E2, *E3},
}

var endpoints = map[string]*Endpoint{
	"37866524-cc91-4d64-b5db-b912eaf4339e": E1,
	"37866524-cc91-4d64-b5db-b912eaf4339u": E2,
	"37866524-cc91-4d64-b5db-b912eaf4338u": E3,
}
var sites = map[string]*Site{
	"6500023": S1,
}

func GetSite(id string) interface{} {
	if site, ok := sites[id]; ok {
		return site
	}
	return nil
}
func GetEndpoint(id string) interface{} {
	if endpoint, ok := endpoints[id]; ok {
		return endpoint
	}
	return nil
}

// dummy fetcher
func GetSites() *Site {
	return S1
}

// dummy fetcher
func GetEndpoints() *Endpoint {
	return E1
}
