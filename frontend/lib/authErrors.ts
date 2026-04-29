/** Erro de autenticação com código opcional da API (ex.: TENANT_INACTIVE_VECX). */
export class AuthError extends Error {
  readonly code?: string;

  constructor(message: string, code?: string) {
    super(message);
    this.name = 'AuthError';
    this.code = code;
  }
}
