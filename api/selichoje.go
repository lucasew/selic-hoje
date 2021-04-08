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
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
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

type rcwrap struct {
	r interface{}
}

func (r rcwrap) Read(b []byte) (int, error) {
	return r.r.(io.Reader).Read(b)
}

func (rcwrap) Close() error {
	return nil
}
