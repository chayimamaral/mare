package domain

import "time"

// ActiveSessionRow representa uma linha em vecx_audit.active_sessions (EF-929).
type ActiveSessionRow struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	UserEmail     string    `json:"user_email"`
	TenantID      string    `json:"tenant_id"`
	TenantName    string    `json:"tenant_name"`
	TenantCNPJ    string    `json:"tenant_cnpj"`
	TenantContato string    `json:"tenant_contato"`
	LoggedAt      time.Time `json:"logged_at"`
	IPAddress     string    `json:"ip_address"`
	Active        bool      `json:"active"`
}
