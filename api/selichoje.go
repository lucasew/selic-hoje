package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	error_reporter "github.com/lucasew/bcb-selic-hoje/src/reporter"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	w.Header().Set("Cache-Control", "max-age=3600")
	data, err := requestData()
	if err != nil {
		error_reporter.ReportError(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	fmt.Fprint(w, getOnlySelic(data))
}

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

func requestData() (string, error) {
	buf := bytes.NewBufferString("filtro=%7B%22dataInicial%22%3A%2207%2F04%2F2021%22%2C%22dataFinal%22%3A%2207%2F04%2F2021%22%7D&parametrosOrdenacao=%5B%5D")
	req, err := http.NewRequest("POST", "https://www3.bcb.gov.br/novoselic/rest/taxaSelicApurada/pub/exportarCsv", buf)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:86.0) Gecko/20100101 Firefox/86.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "pt-BR,en-US;q=0.7,en;q=0.3")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Referer", "https://www3.bcb.gov.br/novoselic/pesquisa-taxa-apurada.jsp")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
