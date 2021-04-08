package handler

import (
	"testing"
)

func TestParse(t *testing.T) {
    const in = `Taxa Selic - Dados diários;Filtros aplicados: Data inicial: 07/04/2021 / Data final: 07/04/2021.;;;;;;;;;
Data;Taxa (% a.a.);Fator diário;Financeiro (R$);Operações;Média;Mediana;Moda;Desvio padrão;Índice de curtose;
07/04/2021;2,65;1,00010379;1.445.859.744.518,52;765;2,65;2,64;2,65;0,014;526,421;
`
    const expectedOut = "2.65"
    out := getOnlySelic(in)
    if (out != expectedOut) {
        t.Errorf("expected '%s' got '%s'", expectedOut, out)
    }
    // Personal debugging
    // data, err := requestData()
    // if err != nil {
    //     t.Error(err)
    // }
    // println("DATA:")
    // println(data)
    // println("PARSED:")
    // println(getOnlySelic(data))
}
