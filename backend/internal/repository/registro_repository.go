package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RegistroUpdateInput struct {
	CNPJ        string
	CEP         string
	Endereco    string
	Bairro      string
	Cidade      string
	Estado      string
	Telefone    string
	Email       string
	IE          string
	IM          string
	RazaoSocial string
	Fantasia    string
	Observacoes string
}

type RegistroCreateInput struct {
	Nome     string
	Email    string
	Password string
}

type DadosComplementaresRecord struct {
	Tenantid    string         `json:"tenantid"`
	CNPJ        sql.NullString `json:"cnpj"`
	CEP         sql.NullString `json:"cep"`
	Endereco    sql.NullString `json:"endereco"`
	Bairro      sql.NullString `json:"bairro"`
	Cidade      sql.NullString `json:"cidade"`
	Estado      sql.NullString `json:"estado"`
	Telefone    sql.NullString `json:"telefone"`
	Email       sql.NullString `json:"email"`
	IE          sql.NullString `json:"ie"`
	IM          sql.NullString `json:"im"`
	RazaoSocial sql.NullString `json:"razaosocial"`
	Fantasia    sql.NullString `json:"fantasia"`
	Observacoes sql.NullString `json:"observacoes"`
}

type RegistroUserRecord struct {
	ID       string `json:"id"`
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	TenantID string `json:"tenantid"`
	Active   bool   `json:"active"`
}

type RegistroRepository struct {
	pool *pgxpool.Pool
}

func NewRegistroRepository(pool *pgxpool.Pool) *RegistroRepository {
	return &RegistroRepository{pool: pool}
}

func (r *RegistroRepository) DetailByTenant(ctx context.Context, tenantID string) (DadosComplementaresRecord, error) {
	if tenantID == "" {
		return DadosComplementaresRecord{}, nil
	}

	const query = `SELECT tenantid, cnpj, cep, endereco, bairro, cidade, estado, telefone, email, ie, im, razaosocial, fantasia, observacoes FROM public.dadoscomplementares WHERE tenantid::text = $1 LIMIT 1`

	var record DadosComplementaresRecord
	err := r.pool.QueryRow(ctx, query, tenantID).Scan(
		&record.Tenantid,
		&record.CNPJ,
		&record.CEP,
		&record.Endereco,
		&record.Bairro,
		&record.Cidade,
		&record.Estado,
		&record.Telefone,
		&record.Email,
		&record.IE,
		&record.IM,
		&record.RazaoSocial,
		&record.Fantasia,
		&record.Observacoes,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Some tenants (e.g. manually provisioned SUPER tenant) may not have
			// dadoscomplementares yet; return an empty payload instead of 400.
			return DadosComplementaresRecord{Tenantid: tenantID}, nil
		}
		return DadosComplementaresRecord{}, fmt.Errorf("detail registro: %w", err)
	}
	return record, nil
}

func (r *RegistroRepository) UpdateByUser(ctx context.Context, userID string, input RegistroUpdateInput) (DadosComplementaresRecord, error) {
	const tenantQuery = `SELECT tenantid FROM public.usuario WHERE id = $1`
	var tenantID string
	if err := r.pool.QueryRow(ctx, tenantQuery, userID).Scan(&tenantID); err != nil {
		return DadosComplementaresRecord{}, fmt.Errorf("find tenant by user: %w", err)
	}

	const query = `
		UPDATE public.dadoscomplementares
		SET cnpj = $1,
			cep = $2,
			endereco = $3,
			bairro = $4,
			cidade = $5,
			estado = $6,
			telefone = $7,
			email = $8,
			ie = $9,
			im = $10,
			razaosocial = $11,
			fantasia = $12,
			observacoes = $13
		WHERE tenantid = $14
		RETURNING tenantid, cnpj, cep, endereco, bairro, cidade, estado, telefone, email, ie, im, razaosocial, fantasia, observacoes`

	var record DadosComplementaresRecord
	err := r.pool.QueryRow(
		ctx,
		query,
		input.CNPJ,
		input.CEP,
		input.Endereco,
		input.Bairro,
		input.Cidade,
		input.Estado,
		input.Telefone,
		input.Email,
		input.IE,
		input.IM,
		input.RazaoSocial,
		input.Fantasia,
		input.Observacoes,
		tenantID,
	).Scan(
		&record.Tenantid,
		&record.CNPJ,
		&record.CEP,
		&record.Endereco,
		&record.Bairro,
		&record.Cidade,
		&record.Estado,
		&record.Telefone,
		&record.Email,
		&record.IE,
		&record.IM,
		&record.RazaoSocial,
		&record.Fantasia,
		&record.Observacoes,
	)
	if err != nil {
		return DadosComplementaresRecord{}, fmt.Errorf("update registro: %w", err)
	}
	return record, nil
}

func (r *RegistroRepository) Create(ctx context.Context, input RegistroCreateInput) (RegistroUserRecord, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return RegistroUserRecord{}, fmt.Errorf("begin tx registro create: %w", err)
	}
	defer tx.Rollback(ctx)

	const existsQuery = `SELECT EXISTS(SELECT 1 FROM public.usuario WHERE email = $1)`
	var exists bool
	if err := tx.QueryRow(ctx, existsQuery, input.Email).Scan(&exists); err != nil {
		return RegistroUserRecord{}, fmt.Errorf("check user exists: %w", err)
	}
	if exists {
		return RegistroUserRecord{}, fmt.Errorf("Usuario ja cadastrado")
	}

	const tenantQuery = `
		INSERT INTO public.tenant (nome, active, plano, contato)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	var tenantID string
	if err := tx.QueryRow(ctx, tenantQuery, input.Nome, true, "DEMO", input.Nome).Scan(&tenantID); err != nil {
		return RegistroUserRecord{}, fmt.Errorf("create tenant: %w", err)
	}

	const dadosQuery = `
		INSERT INTO public.dadoscomplementares (tenantid)
		VALUES ($1)`
	if _, err := tx.Exec(ctx, dadosQuery, tenantID); err != nil {
		return RegistroUserRecord{}, fmt.Errorf("create dadoscomplementares: %w", err)
	}

	const userQuery = `
		INSERT INTO public.usuario (nome, email, password, role, tenantid, active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, nome, email, role, tenantid, active`

	var record RegistroUserRecord
	if err := tx.QueryRow(ctx, userQuery, input.Nome, input.Email, input.Password, "ADMIN", tenantID, true).Scan(
		&record.ID,
		&record.Nome,
		&record.Email,
		&record.Role,
		&record.TenantID,
		&record.Active,
	); err != nil {
		return RegistroUserRecord{}, fmt.Errorf("create user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return RegistroUserRecord{}, fmt.Errorf("commit registro create: %w", err)
	}

	return record, nil
}
