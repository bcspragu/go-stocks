package stock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	VanguardSuggestUrl = "https://api.vanguard.com/rs/sae/01/autosuggest.json?types=funds&limit=10&query=%s"
	VanguardFundUrl    = "https://personal.vanguard.com/us/JSP/Funds/VGITab/VGIFundOverviewTabContent.jsf?FundIntExt=INT&FundId=%s"
)

type vanguardRetreiver struct{}

func NewVanguard() Retreiver {
	return &vanguardRetreiver{}
}

func (v *vanguardRetreiver) Price(stock string) (float64, error) {
	f, err := findFund(stock)
	if err != nil {
		return 0.0, err
	}

	return getPrice(f.FundID)
}

// Example query: https://api.vanguard.com/rs/sae/01/autosuggest.json?types=funds,term&limit=10&limits=6,4&query=VXUS
// Example response {"type":"autosuggest","results":[{"tickerSymbol":"VXUS","fundID":"3369","term":"Vanguard Total International Stock ETF (VXUS) (3369)","type":"fund"}]}
type suggestResponse struct {
	Results []fund
}

type fund struct {
	TickerSymbol string
	FundID       string
}

func findFund(stock string) (*fund, error) {
	suggestURL := fmt.Sprintf(VanguardSuggestUrl, stock)

	resp, err := http.Get(suggestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sResp suggestResponse
	err = json.NewDecoder(resp.Body).Decode(&sResp)
	if err != nil {
		return nil, err
	}
	if len(sResp.Results) == 0 {
		return nil, fmt.Errorf("No results found for %s", stock)
	}

	// Search the results for a fund with the given name
	var f *fund
	for _, fund := range sResp.Results {
		if fund.TickerSymbol == stock {
			f = &fund
			break
		}
	}

	if f == nil {
		return nil, fmt.Errorf("No results found for %s", stock)
	}

	return f, nil
}

func getPrice(fundID string) (float64, error) {
	fundUrl := fmt.Sprintf(VanguardFundUrl, fundID)
	doc, err := goquery.NewDocument(fundUrl)
	if err != nil {
		return 0.0, err
	}

	p := doc.Find("#performanceTabletbody0 tr").Eq(0).Find("td").Eq(1).Text()
	// Confirm p has a leading dollar sign
	p = strings.TrimSpace(p)
	if p[0] != '$' {
		return 0.0, fmt.Errorf("Expected field to start with dollar sign, found %s", p)
	}
	return strconv.ParseFloat(p[1:], 64)
}
