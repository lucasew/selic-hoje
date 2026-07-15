package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// func main() {
//     data, err := requestData()
//     if err != nil {
//         panic(err)
//     }
//     println(data)
// }

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	w.Header().Set("Cache-Control", "max-age=3600")
	data, err := requestData()
	if err != nil {
		ReportError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	fmt.Fprint(w, getOnlySelic(data))
}

// getOnlySelic returns the most recent non-zero daily Selic rate (% a.a.)
// from a BCB novoselic CSV export. Rows are newest-first; stub/zero days are skipped.
func getOnlySelic(in string) string {
	in = strings.TrimPrefix(in, "\ufeff")
	for _, line := range strings.Split(in, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Taxa Selic") || strings.HasPrefix(line, "Data;") {
			continue
		}
		values := strings.Split(line, ";")
		if len(values) < 2 {
			continue
		}
		rate := strings.TrimSpace(values[1])
		if rate == "" || rate == "0" {
			continue
		}
		return strings.ReplaceAll(rate, ",", ".")
	}
	return ""
}

// exportForm builds the BCB export body for a window ending at now (last lookbackDays days).
func exportForm(now time.Time, lookbackDays int) string {
	end := now
	start := now.AddDate(0, 0, -lookbackDays)
	filtro := fmt.Sprintf(`{"dataInicial":"%s","dataFinal":"%s"}`,
		start.Format("02/01/2006"),
		end.Format("02/01/2006"))
	v := url.Values{}
	v.Set("filtro", filtro)
	v.Set("parametrosOrdenacao", "[]")
	return v.Encode()
}

func requestData() (string, error) {
	const endpoint = "https://www3.bcb.gov.br/novoselic/rest/taxaSelicApurada/pub/exportarCsv"
	// 14-day window covers weekends/holidays so the latest published day is included.
	form := exportForm(time.Now(), 14)
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form))
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:86.0) Gecko/20100101 Firefox/86.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,en-US;q=0.7,en;q=0.3")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Referer", "https://www3.bcb.gov.br/novoselic/pesquisa-taxa-apurada.jsp")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("bcb selic export: unexpected status %s", res.Status)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ReportError(err error) {
	log.Println(err)
}
