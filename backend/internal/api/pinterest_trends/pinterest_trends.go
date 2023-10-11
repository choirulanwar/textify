package pinterest_trends_api

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/choirulanwar/textify/backend/pkg/helper"
	"github.com/choirulanwar/textify/backend/pkg/http"
	"gorm.io/gorm"
)

type PinterestTrendsAPIService struct {
	Period  string `json:"period"`
	Country string `json:"country"`
	EndDate string `json:"end_date"`
	DB      *gorm.DB
}

func NewPinterestTrend(period string, country string, date string, db *gorm.DB) *PinterestTrendsAPIService {
	if country == "" {
		country = "US"
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	return &PinterestTrendsAPIService{
		Period:  period,
		Country: country,
		EndDate: date,
		DB:      db,
	}
}

func (s *PinterestTrendsAPIService) GetTrends(keywordsToInclude string) (*PTExplorerResp, error) {
	settingInfo, err := helper.GetSettingInfo(s.DB)
	if err != nil {
		return nil, err
	}

	url := s.buildTrendsURL(keywordsToInclude)

	response, err := http.MakeGetRequest(url, "", settingInfo.SessionPinterest)
	if err != nil {
		return nil, err
	}

	var pinterestTrends *PTExplorerResp
	err = json.Unmarshal(response, &pinterestTrends)
	if err != nil {
		return nil, err
	}

	return pinterestTrends, nil
}

func (s *PinterestTrendsAPIService) GetRelated(term string) ([]*PTRelatedResp, error) {
	settingInfo, err := helper.GetSettingInfo(s.DB)
	if err != nil {
		return nil, err
	}

	lookback := s.getLookback()

	url := s.buildRelatedURL(term, lookback)

	response, err := http.MakeGetRequest(url, "", settingInfo.SessionPinterest)
	if err != nil {
		return nil, err
	}

	var relatedsByTerm []*PTRelatedResp
	err = json.Unmarshal(response, &relatedsByTerm)
	if err != nil {
		return nil, err
	}

	return relatedsByTerm, nil
}

func (s *PinterestTrendsAPIService) GetMetric(term string) ([]*PTMetricResp, error) {
	settingInfo, err := helper.GetSettingInfo(s.DB)
	if err != nil {
		return nil, err
	}

	url := s.buildMetricURL(term)

	response, err := http.MakeGetRequest(url, "", settingInfo.SessionPinterest)
	if err != nil {
		return nil, err
	}

	var metricsByTerm []*PTMetricResp
	err = json.Unmarshal(response, &metricsByTerm)
	if err != nil {
		return nil, err
	}

	return metricsByTerm, nil
}

func (s *PinterestTrendsAPIService) buildTrendsURL(keywordsToInclude string) string {
	url := fmt.Sprintf("https://trends.pinterest.com/top_trends_filtered/?endDate=%s&country=%s&trendsPreset=%s&numTermsToReturn=100",
		s.EndDate, s.Country, s.Period)
	if keywordsToInclude != "" {
		url = fmt.Sprintf("%s&keywordsToInclude=%s", url, strings.ToLower(keywordsToInclude))
	}
	return url
}

func (s *PinterestTrendsAPIService) buildRelatedURL(term string, lookback int) string {
	url := fmt.Sprintf("https://trends.pinterest.com/related_terms/?requestTerm=%s&country=%s&endDate=%s&aggregation=2&lookback=%d",
		strings.ToLower(term), s.Country, s.EndDate, lookback)
	return url
}

func (s *PinterestTrendsAPIService) buildMetricURL(term string) string {
	url := fmt.Sprintf("https://trends.pinterest.com/metrics/?terms=%s&country=%s&end_date=%s&days=365&aggregation=2&shouldMock=false&normalize_against_group=true&predicted_days=0",
		strings.ToLower(term), s.Country, s.EndDate)
	return url
}

func (s *PinterestTrendsAPIService) getLookback() int {
	switch s.Period {
	case "1", "4":
		return 30
	case "2":
		return 365
	case "3":
		return 90
	default:
		return 365
	}
}
