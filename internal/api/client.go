package api

import (
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	URL url.URL
	Key string
}

func NewClient(url url.URL, key string) *Client {
	return &Client{URL: url, Key: key}
}

func (c *Client) GetPerformanceInformation(ISIN string) (PerformanceInformation, error) {
	c.URL.RawQuery = "organization=PCPK21&api-version=beta-1.0.0"
	uri := c.URL.JoinPath("beta/web/calc/union-tools-pcpk/public/product/performanceInformation").JoinPath(ISIN).String()

	request, _ := http.NewRequest("GET", uri, nil)
	request.Header.Set("x-api-key", c.Key)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return PerformanceInformation{}, err
	}

	if response.StatusCode != http.StatusOK {
		return PerformanceInformation{}, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var performanceInformation PerformanceInformation
	if err = performanceInformation.Decode(response.Body); err != nil {
		return PerformanceInformation{}, err
	}

	return performanceInformation, nil
}
