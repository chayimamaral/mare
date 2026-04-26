package auth

import "github.com/chayimamaral/vecontab/backend/internal/domain"

func TenantClaimsFromDomain(t domain.Tenant) TenantClaims {
	return TenantClaims{
		ID:          t.ID,
		Active:      t.Active,
		Nome:        t.Nome,
		Contato:     t.Contato,
		Plano:       t.Plano,
		SchemaName:  t.SchemaName,
		IsVecMaster: t.IsVecMaster,
	}
}
