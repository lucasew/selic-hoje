package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// metaSelicURL is BCB SGS series 432 (Meta Selic, % a.a.) — the Copom target
// rate commonly quoted as "a Selic". Open Data JSON API, no scrape.
const metaSelicURL = "https://api.bcb.gov.br/dados/serie/bcdata.sgs.432/dados/ultimos/1?formato=json"

// httpClient bounds upstream latency so a hung BCB does not pin the function.
var httpClient = &http.Client{Timeout: 15 * time.Second}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	w.Header().Set("Cache-Control", "max-age=3600")
	rate, err := fetchMetaSelic()
	if err != nil {
		ReportError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, rate)
}

type sgsPoint struct {
	Data  string `json:"data"`
	Valor string `json:"valor"`
}

func fetchMetaSelic() (string, error) {
	req, err := http.NewRequest(http.MethodGet, metaSelicURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("bcb sgs 432: unexpected status %s", res.Status)
	}
	// Cap body size: series response is tiny; bound memory anyway.
	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return "", err
	}
	return parseMetaSelic(body)
}

// parseMetaSelic extracts the latest Meta Selic value from an SGS JSON array.
func parseMetaSelic(body []byte) (string, error) {
	var points []sgsPoint
	if err := json.Unmarshal(body, &points); err != nil {
		return "", fmt.Errorf("bcb sgs 432: decode: %w", err)
	}
	if len(points) == 0 || points[0].Valor == "" {
		return "", fmt.Errorf("bcb sgs 432: empty series")
	}
	return points[0].Valor, nil
}

func ReportError(err error) {
	log.Println(err)
}
