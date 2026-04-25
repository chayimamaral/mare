package service

import (
	"testing"
	"time"
)

func TestProximaConsultaDistribuicaoSC(t *testing.T) {
	now := time.Now().UTC()
	cases := []struct {
		cstat    int
		qt       int
		minHours float64
		maxNil   bool
	}{
		{110, 0, 1, false},
		{117, 0, 12, false},
		{118, 49, 12, false},
		{118, 50, 0, true},
		{138, 1, 12, false},
		{108, 0, 1, false},
		{657, 0, 1, false},
		{143, 0, 0, true},
		{146, 0, 0, true},
		{999, 0, 0, true},
	}
	for _, tc := range cases {
		got := proximaConsultaDistribuicaoSC(tc.cstat, tc.qt)
		if tc.maxNil {
			if got != nil {
				t.Fatalf("cstat=%d qt=%d: esperado nil, obteve %v", tc.cstat, tc.qt, got)
			}
			continue
		}
		if got == nil {
			t.Fatalf("cstat=%d qt=%d: esperado horário, obteve nil", tc.cstat, tc.qt)
		}
		d := got.Sub(now).Hours()
		if d+0.01 < tc.minHours || d > tc.minHours+0.05 {
			t.Fatalf("cstat=%d qt=%d: delta horas=%.2f queria ~%.1f", tc.cstat, tc.qt, d, tc.minHours)
		}
	}
}
