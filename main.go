package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/getsentry/sentry-go"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/fsrv-xyz/stonkschamp/internal/api"
)

type InfluxConfig struct {
	Bucket string
	Org    string
}

func main() {
	influxConfig := InfluxConfig{}
	app := kingpin.New(os.Args[0], "health-export")
	unionInvenstmentApiToken := app.Flag("ui-api.token", "union investment api token").Envar("UI_API_TOKEN").Required().String()
	unionInvestmentISIN := app.Flag("isin", "union investment fond isin").Envar("ISIN").Required().String()

	// influxdb configuration
	influxdbAuthToken := app.Flag("influx.token", "influxdb auth token").Envar("INFLUXDB_TOKEN").Required().String()
	influxdbURL := app.Flag("influx.url", "influxdb url").Envar("INFLUXDB_URL").Required().URL()
	app.Flag("influx.bucket", "influxdb bucket").Envar("INFLUXDB_BUCKET").Required().StringVar(&influxConfig.Bucket)
	app.Flag("influx.org", "influxdb org").Envar("INFLUXDB_ORG").Required().StringVar(&influxConfig.Org)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	// initialize sentry instrumentation
	sentry.Init(sentry.ClientOptions{TracesSampleRate: 1.0, Transport: sentry.NewHTTPSyncTransport(), Debug: true, EnableTracing: true})
	sentryRootSpan := sentry.StartSpan(context.Background(), "stonkschamp", sentry.WithTransactionName("stonkschamp"))
	defer sentryRootSpan.Finish()
	defer sentry.Flush(2 * time.Second)

	influxClient := influxdb2.NewClient((*influxdbURL).String(), *influxdbAuthToken)

	client := api.NewClient(url.URL{
		Scheme: "https",
		Host:   "internal.api.union-investment.de",
	}, *unionInvenstmentApiToken)

	sentryDataFetchSpan := sentryRootSpan.StartChild("get performance information")
	sentryDataFetchSpan.SetTag("isin", *unionInvestmentISIN)

	performanceInformation, err := client.GetPerformanceInformation(*unionInvestmentISIN)

	sentryDataFetchSpan.Finish()
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	writeAPI := influxClient.WriteAPIBlocking(influxConfig.Org, influxConfig.Bucket)

	sentryDataIngresSpan := sentryRootSpan.StartChild("write performance information")
	defer sentryDataIngresSpan.Finish()

	for _, performanceMetric := range performanceInformation.ProductsPerformance[0].TimeSeries.TimeSeries {

		point := influxdb2.NewPoint(
			"metrics",
			map[string]string{
				"isin":     performanceInformation.ProductsPerformance[0].Isin,
				"currency": performanceInformation.ProductsPerformance[0].Currency,
			},
			map[string]interface{}{
				"value": performanceMetric.Value,
			},
			performanceMetric.Date.Time,
		)

		writePointError := writeAPI.WritePoint(context.Background(), point)
		if writePointError != nil {
			sentry.CaptureException(writePointError)
			log.Println(writePointError)
		}
	}
}
