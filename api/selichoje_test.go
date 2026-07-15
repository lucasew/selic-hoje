package handler

import (
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	const in = `Taxa Selic - Dados diários;Filtros aplicados: Data inicial: 07/04/2021 / Data final: 07/04/2021.;;;;;;;;;
Data;Taxa (% a.a.);Fator diário;Financeiro (R$);Operações;Média;Mediana;Moda;Desvio padrão;Índice de curtose;
07/04/2021;2,65;1,00010379;1.445.859.744.518,52;765;2,65;2,64;2,65;0,014;526,421;
`
	const expectedOut = "2.65"
	out := getOnlySelic(in)
	if out != expectedOut {
		t.Errorf("expected '%s' got '%s'", expectedOut, out)
	}
}

func TestParseSkipsZeroAndPicksLatest(t *testing.T) {
	// Shape matches current BCB export (BOM, multi-day, newest first, zero stub day).
	const in = "\ufeff" + `Taxa Selic - Dados diários;Filtros aplicados: Data inicial: 01/07/2026 / Data final: 15/07/2026.;;;;;;;;;
Data;Taxa (% a.a.);Fator diário;Financeiro (R$);Operações;Taxa média;Taxa mediana;Taxa modal;Desvio padrão;Índice de curtose;
15/07/2026;0;0,00000000;0;0;0;0;0;0,0000;0,0000;
14/07/2026;14,15;1,00052531;1.352.202.436.721,45;858;14,15;14,14;14,15;0,0241;410,0909;
13/07/2026;14,15;1,00052531;1.333.848.576.222,51;836;14,15;14,14;14,15;0,0234;473,1097;
`
	if out := getOnlySelic(in); out != "14.15" {
		t.Errorf("expected 14.15 got %q", out)
	}
}

func TestExportForm(t *testing.T) {
	now := time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)
	got := exportForm(now, 14)
	if !strings.Contains(got, "01%2F07%2F2026") {
		t.Fatalf("form missing start date: %s", got)
	}
	if !strings.Contains(got, "15%2F07%2F2026") {
		t.Fatalf("form missing end date: %s", got)
	}
	if !strings.Contains(got, "parametrosOrdenacao") {
		t.Fatalf("form missing parametrosOrdenacao: %s", got)
	}
}
