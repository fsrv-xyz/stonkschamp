package api

import (
	"encoding/json"
	"io"
	"time"
)

func (p *PerformanceInformation) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(p)
}

type PerformanceInformation struct {
	ProductsPerformance []ProductsPerformance `json:"productsPerformance"`
}

type TimeSeries struct {
	Isin       string                     `json:"isin"`
	TimeSeries map[string]TimeSeriesEntry `json:"timeSeries"`
}

type ProductsPerformance struct {
	Isin           string     `json:"isin"`
	Name           string     `json:"name"`
	Currency       string     `json:"currency"`
	CurrencySymbol string     `json:"currencySymbol"`
	TimeSeries     TimeSeries `json:"timeSeries"`
}

type TimeSeriesEntry struct {
	Date  CustomTime `json:"date"`
	Value float64    `json:"value"`
}

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(b []byte) error {
	var err error
	t.Time, err = time.Parse("2006-01-02", string(b[1:len(b)-1]))
	return err
}
