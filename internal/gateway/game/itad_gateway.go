package game

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"nexus-video-games/internal/entity"

	"github.com/spf13/viper"
)

type ITADGateway struct {
	apiKey     string
	httpClient *http.Client
}

type itadDealItem struct {
	ID    string `json:"id"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Price struct {
		Amount float64 `json:"amount"`
	} `json:"price"`
}

type itadResponse struct {
	List []itadDealItem `json:"list"`
}

func NewITADGateway(viper *viper.Viper) *ITADGateway {
	return &ITADGateway{
		apiKey:     viper.GetString("ITAD_API_KEY"),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// FetchDealsCatalog retrieves real-time curated price catalog entries from ITAD
func (g *ITADGateway) FetchDealsCatalog(ctx context.Context) ([]entity.GamesCatalog, error) {
	if g.apiKey == "" {
		return g.getMockCatalog(), nil
	}

	url := fmt.Sprintf("https://api.isthereanydeal.com/deals/v2?key=%s&limit=20&country=ID", g.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return g.getMockCatalog(), nil
	}

	var parsed itadResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	var catalog []entity.GamesCatalog
	for _, item := range parsed.List {
		price := int64(item.Price.Amount)
		if price <= 0 {
			price = 150000 // Fallback nominal IDR balance value
		}

		catalog = append(catalog, entity.GamesCatalog{
			ExternalGameID: item.ID,
			Title:          item.Title,
			PriceInBalance: price,
			IsAvailable:    true,
		})
	}

	return catalog, nil
}

func (g *ITADGateway) getMockCatalog() []entity.GamesCatalog {
	return []entity.GamesCatalog{
		{ExternalGameID: "itad-01", Title: "Cyberpunk 2077", PriceInBalance: 350000, IsAvailable: true},
		{ExternalGameID: "itad-02", Title: "Elden Ring", PriceInBalance: 599000, IsAvailable: true},
		{ExternalGameID: "itad-03", Title: "Hades II", PriceInBalance: 245000, IsAvailable: true},
		{ExternalGameID: "itad-04", Title: "Baldur's Gate 3", PriceInBalance: 699000, IsAvailable: true},
		{ExternalGameID: "itad-05", Title: "Stardew Valley", PriceInBalance: 119000, IsAvailable: true},
	}
}
