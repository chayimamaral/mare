package service

import (
	"context"

	"github.com/chayimamaral/vecontab/backendgo/internal/repository"
)

type ListUsersInput struct {
	First     int
	Rows      int
	SortField string
	SortOrder int
	Nome      string
	TenantID  string
	Role      string
}

type CreateUserInput struct {
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	TenantID string `json:"tenantId"`
}

type UserService struct {
	users *repository.UserRepository
}

type TotalRecordsResponse struct {
	ResTotal int64 `json:"resTotal"`
}

type ListUsersResponse struct {
	Usuarios     []repository.UserListItem `json:"usuarios,omitempty"`
	TotalRecords *TotalRecordsResponse     `json:"totalRecords,omitempty"`
	Error        string                    `json:"error,omitempty"`
}

type CreateUserResponse struct {
	Usuarios []repository.UserListItem `json:"usuarios"`
}

func NewUserService(users *repository.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) Detail(ctx context.Context, userID string) (repository.UserDetailResponse, error) {
	return s.users.Detail(ctx, userID)
}

func (s *UserService) UserRole(ctx context.Context, userID string) (repository.UserRoleResponse, error) {
	return s.users.UserRole(ctx, userID)
}

func (s *UserService) TenantID(ctx context.Context, userID string) (repository.UserTenantIDResponse, error) {
	return s.users.TenantID(ctx, userID)
}

func (s *UserService) List(ctx context.Context, input ListUsersInput) (ListUsersResponse, error) {
	if input.Role != "ADMIN" && input.Role != "SUPER" {
		return ListUsersResponse{Error: "Usuario nao autorizado"}, nil
	}

	usuarios, total, err := s.users.ListByTenant(
		ctx,
		input.TenantID,
		input.First,
		input.Rows,
		input.SortField,
		input.SortOrder,
		input.Nome,
	)
	if err != nil {
		return ListUsersResponse{}, err
	}

	return ListUsersResponse{
		Usuarios:     usuarios,
		TotalRecords: &TotalRecordsResponse{ResTotal: total},
	}, nil
}

func (s *UserService) Create(ctx context.Context, input CreateUserInput) (CreateUserResponse, error) {
	usuarios, err := s.users.Create(ctx, input.Nome, input.Email, input.Password, input.Role, input.TenantID)
	if err != nil {
		return CreateUserResponse{}, err
	}

	return CreateUserResponse{Usuarios: usuarios}, nil
}
