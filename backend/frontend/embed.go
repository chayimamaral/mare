package frontend

import "embed"

// FS contém o export estático do Next (pasta out). Em produção o Dockerfile
// preenche frontend/out a partir de frontend_static.
//
//go:embed all:out
var FS embed.FS
