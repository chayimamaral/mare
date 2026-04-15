package repository

import (
	"testing"
	"time"
)

func TestResolveCompetenciaMensal(t *testing.T) {
	loc := time.Local
	tests := []struct {
		name       string
		dataInicio time.Time
		agora      time.Time
		want       time.Time
	}{
		{
			name:       "vigencia futura mantem mes de referencia",
			dataInicio: time.Date(2026, 5, 25, 0, 0, 0, 0, loc),
			agora:      time.Date(2026, 4, 15, 9, 0, 0, 0, loc),
			want:       time.Date(2026, 5, 25, 0, 0, 0, 0, loc),
		},
		{
			name:       "vigencia hoje mantem mes atual",
			dataInicio: time.Date(2026, 4, 15, 0, 0, 0, 0, loc),
			agora:      time.Date(2026, 4, 15, 12, 0, 0, 0, loc),
			want:       time.Date(2026, 4, 15, 0, 0, 0, 0, loc),
		},
		{
			name:       "vigencia passada avanca para proximo mes",
			dataInicio: time.Date(2026, 4, 10, 0, 0, 0, 0, loc),
			agora:      time.Date(2026, 4, 15, 12, 0, 0, 0, loc),
			want:       time.Date(2026, 5, 10, 0, 0, 0, 0, loc),
		},
		{
			name:       "dia 31 respeita fim do mes ao avancar",
			dataInicio: time.Date(2026, 1, 31, 0, 0, 0, 0, loc),
			agora:      time.Date(2026, 2, 5, 12, 0, 0, 0, loc),
			want:       time.Date(2026, 2, 28, 0, 0, 0, 0, loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveCompetenciaMensal(tt.dataInicio, tt.agora)
			if !got.Equal(tt.want) {
				t.Fatalf("resolveCompetenciaMensal(%s, %s) = %s; want %s",
					tt.dataInicio.Format("2006-01-02"),
					tt.agora.Format("2006-01-02"),
					got.Format("2006-01-02"),
					tt.want.Format("2006-01-02"),
				)
			}
		})
	}
}
