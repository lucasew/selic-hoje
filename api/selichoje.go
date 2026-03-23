package handler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Handler fetches the latest Selic rate from the BCB API and serves it as plain text.
// It is designed as a serverless entrypoint for Vercel. Errors from the BCB API
// result in a 500 response containing the error string. Responses are aggressively
// cached (1 hour) to avoid rate limits and minimize latency.
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	w.Header().Set("Cache-Control", "max-age=3600")
	data, err := requestData()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}
	w.WriteHeader(200)
	fmt.Fprint(w, getOnlySelic(data))
}

// getOnlySelic parses the CSV response from the BCB API to extract the daily Selic rate.
// The CSV is expected to be structurally rigid: 4 lines total, with the data on the 3rd line.
// It isolates the rate column (index 1) and standardizes the decimal separator to a period.
// Edge case: If the input format drifts (e.g., more/fewer lines or columns), it silently returns
// an empty string.
func getOnlySelic(in string) string {
	lines := strings.Split(in, "\n")
	if len(lines) != 4 {
		return ""
	}
	values := strings.Split(lines[2], ";")
	if len(values) != 11 {
		return ""
	}
	return strings.ReplaceAll(values[1], ",", ".")
}

// requestData constructs and executes a POST request to the BCB API to export the Selic rate CSV.
// Nuance: The API uses a POST endpoint that typically expects form data, requiring specific
// headers like 'Content-Type' and a browser-like 'User-Agent' to avoid being blocked.
// Limitation: Currently sends a hardcoded date filter (07/04/2021) within the request body.
// Performance Implication: Uses http.DefaultClient without an explicit timeout, which could
// lead to resource exhaustion if the BCB API hangs.
func requestData() (string, error) {
	var err error
	req := new(http.Request)
	req.URL, err = url.Parse("https://www3.bcb.gov.br/novoselic/rest/taxaSelicApurada/pub/exportarCsv")
	if err != nil {
		return "", err
	}
	req.Header = http.Header{}
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:86.0) Gecko/20100101 Firefox/86.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "pt-BR,en-US;q=0.7,en;q=0.3")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Referer", "https://www3.bcb.gov.br/novoselic/pesquisa-taxa-apurada.jsp")
	buf := bytes.NewBufferString("filtro=%7B%22dataInicial%22%3A%2207%2F04%2F2021%22%2C%22dataFinal%22%3A%2207%2F04%2F2021%22%7D&parametrosOrdenacao=%5B%5D")
	req.Body = rcwrap{r: buf}
	req.Method = "POST"
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	err = res.Body.Close()
	if err != nil {
		return "", err
	}
	return string(data), err
}

// rcwrap wraps an io.Reader to implement the io.ReadCloser interface.
// This is necessary because http.Request.Body requires an io.ReadCloser,
// but bytes.NewBufferString only provides an io.Reader.
type rcwrap struct {
	r interface{}
}

// Read delegates the Read call to the underlying io.Reader.
func (r rcwrap) Read(b []byte) (int, error) {
	return r.r.(io.Reader).Read(b)
}

// Close implements the Closer interface with a no-op, satisfying io.ReadCloser
// without requiring the underlying reader to actually be closeable.
func (rcwrap) Close() error {
	return nil
}
