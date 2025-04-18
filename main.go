package main

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rornic/starlingexporter/internal/client"
	"github.com/rornic/starlingexporter/internal/metrics"
)

func initialiseStarlingClient() client.StarlingClient {
	accessToken := getTokenFromEnvironment()

	endpoint := "https://api.starlingbank.com/api/v2"
	sandbox := strings.ToLower(os.Getenv("STARLING_SANDBOX")) == "true"
	if sandbox {
		slog.Info("using sandbox environment")
		endpoint = strings.Replace(endpoint, "api", "api-sandbox", 1)
	}

	client := client.NewStarlingHttpClient(accessToken, endpoint)
	return &client

}

func getTokenFromEnvironment() string {
	var accessToken string

	accessTokenPath := os.Getenv("STARLING_ACCESS_TOKEN_PATH")
	if accessTokenPath != "" {
		slog.Info("found access token path in environment")
		accessTokenBytes, err := os.ReadFile(accessTokenPath)
		if err != nil {
			slog.Error("error reading access token file", "error", err)
			os.Exit(1)
		}
		accessToken = string(accessTokenBytes)
	}
	if accessToken == "" {
		accessToken = os.Getenv("STARLING_ACCESS_TOKEN")
	}
	if accessToken == "" {
		slog.Error("no access token available from STARLING_ACCESS_TOKEN_PATH or STARLING_ACCESS_TOKEN. Exiting.")
		os.Exit(1)
	}
	accessToken = strings.TrimSpace(accessToken)
	slog.Info("found access token")

	return accessToken
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	client := initialiseStarlingClient()
	metrics.Record(client)
	http.Handle("/metrics", promhttp.Handler())

	slog.Info("listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("error running server", "error", err)
		os.Exit(1)
	}
}
