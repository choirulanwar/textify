package google_trends_api

import (
	"encoding/json"
	"fmt"
	net_http "net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/choirulanwar/textify/backend/pkg/http"
)

type GoogleTrendsAPIService struct {
	Period    string `json:"period"`
	Timezone  int    `json:"timezone"`
	HL        string `json:"hl"`
	CookieVal string `json:"cookieVal"`
	Property  string `json:"property"`
	Category  int    `json:"category"`
	Country   string `json:"country"`
}

type Query struct {
	Q    string `json:"q"`
	Time string `json:"time"`
}

const (
	baseUrl     = "trends.google.com"
	exploreUrl  = "/trends/api/explore"
	relatedUrl  = "/trends/api/widgetdata/relatedsearches"
	interestUrl = "/trends/api/widgetdata/multiline"
)

func NewGoogleTrend(property string, category string, country string) *GoogleTrendsAPIService {
	if category == "" {
		category = "0"
	}
	c, err := strconv.Atoi(category)
	if err != nil {
		return nil
	}

	if country == "" {
		country = "US"
	}

	return &GoogleTrendsAPIService{
		Timezone:  -420,
		HL:        "en-US",
		CookieVal: "",
		Property:  property,
		Category:  c,
		Country:   country,
	}
}

func (s *GoogleTrendsAPIService) GetWidgets(param *Query) (*GTExplorerResp, error) {
	c := &comparison{
		Keyword:                param.Q,
		Geo:                    s.Country,
		Backend:                "CM",
		Property:               s.Property,
		Timezone:               s.Timezone,
		GranularTimeResolution: false,
		HL:                     s.HL,
		Category:               strconv.Itoa(s.Category),
	}

	switch param.Time {
	case "now 1-H", "now 4-H", "now 1-d", "now 7-d", "today 1-m", "today 3-m", "today 12-m", "today 5-y":
		c.Time = param.Time

	default:
		formatTime := "2006-01-02"
		timeSplit := strings.Split(param.Time, " ")
		startTime, err := time.Parse(formatTime, timeSplit[0])
		if err != nil {
			return &GTExplorerResp{}, err
		}
		endTime, err := time.Parse(formatTime, timeSplit[1])
		if err != nil {
			return &GTExplorerResp{}, err
		}
		sT := startTime.Format(formatTime)
		eT := endTime.Format(formatTime)
		c.Time = fmt.Sprintf("%s %s", sT, eT)
	}

	reqBody := map[string]interface{}{
		"comparisonItem": []comparison{*c},
		"category":       s.Category,
		"property":       s.Property,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return &GTExplorerResp{}, err
	}

	qs := url.Values{
		"hl":  []string{s.HL},
		"tz":  []string{strconv.Itoa(s.Timezone)},
		"req": []string{string(reqBodyBytes)},
	}

	client := http.NewGTHttpClient()

	req, err := client.Request("POST", baseUrl, exploreUrl, qs, &net_http.Transport{
		DisableKeepAlives: false,
	})

	if err != nil {
		return &GTExplorerResp{}, err
	}

	var explorerResponse *GTExplorerResp
	err = json.Unmarshal(req[4:], &explorerResponse)
	if err != nil {
		return &GTExplorerResp{}, err
	}

	return explorerResponse, nil
}

func (s *GoogleTrendsAPIService) GetRelated(explorer *GTExplorerResp) (*GTRelatedResp, error) {
	relatedReq := explorer.Widgets[3].Request
	relatedToken := explorer.Widgets[3].Token

	reqBodyBytes, err := json.Marshal(relatedReq)
	if err != nil {
		return nil, err
	}

	qs := url.Values{
		"hl":    []string{s.HL},
		"tz":    []string{strconv.Itoa(s.Timezone)},
		"req":   []string{string(reqBodyBytes)},
		"token": []string{relatedToken},
	}

	client := http.NewGTHttpClient()

	req, err := client.Request("GET", baseUrl, relatedUrl, qs, nil)
	if err != nil {
		return &GTRelatedResp{}, err
	}

	var relatedResponse *GTRelatedResp
	err = json.Unmarshal(req[5:], &relatedResponse)
	if err != nil {
		return &GTRelatedResp{}, err
	}

	return relatedResponse, nil
}

func (s *GoogleTrendsAPIService) GetInterest(explorer *GTExplorerResp) (*GTInterestResp, error) {
	interestRequest := explorer.Widgets[0].Request
	interestToken := explorer.Widgets[0].Token

	reqBodyBytes, err := json.Marshal(interestRequest)
	if err != nil {
		return nil, err
	}

	qs := url.Values{
		"hl":    []string{s.HL},
		"tz":    []string{strconv.Itoa(s.Timezone)},
		"req":   []string{string(reqBodyBytes)},
		"token": []string{interestToken},
	}

	client := http.NewGTHttpClient()

	req, err := client.Request("GET", baseUrl, interestUrl, qs, &net_http.Transport{
		DisableKeepAlives: false,
	})
	if err != nil {
		return &GTInterestResp{}, err
	}

	var interestResponse *GTInterestResp
	err = json.Unmarshal(req[5:], &interestResponse)
	if err != nil {
		return &GTInterestResp{}, err
	}

	return interestResponse, nil
}
