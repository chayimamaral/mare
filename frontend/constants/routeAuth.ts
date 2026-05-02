/**
 * Rotas que usam o shell do dashboard (menu lateral) e exigem sessão válida para exibir conteúdo.
 * Mantido alinhado a useClientGuards (AUTH_REQUIRED_ROUTES + ROLE_RESTRICTED + home).
 */
import type { UserRole } from './routeAuthTypes';

export const AUTH_REQUIRED_ROUTES = new Set<string>([
    '/',
    '/agenda',
    '/agenda-arvore',
    '/catalogo-servicos',
    '/cliente-pf',
    '/clientes',
    '/cnae',
    '/caixa-postal',
    '/compromissos',
    '/compromissos-empresas',
    '/compromissos-por-natureza',
    '/compromissos-visao',
    '/configuracoes/api-integra-contador',
    '/configuracoes/certificado-digital',
    '/configuracoes/geracao-guias',
    '/configuracoes/integra-contador-servicos',
    '/configuracoes/integra-contador-tabela-consumo',
    '/empresas',
    '/estados',
    '/feriados',
    '/grupopassos',
    '/matriz-conformidade-fiscal',
    '/monitor',
    '/utilitarios/hardware-manager',
    '/utilitarios/monitoramento-global',
    '/nfe/consulta',
    '/nfe/manutencao',
    '/nfe/sincronizacao',
    '/municipios',
    '/obrigacoes',
    '/pages/empty',
    '/pages/timeline',
    '/passos',
    '/regimes-tributarios',
    '/registro',
    '/representantes',
    '/rotinas',
    '/rotinas-pf',
    '/salario-minimo',
    '/tenants',
    '/tipoempresa',
    '/usuarios',
]);

export const ROLE_RESTRICTED_ROUTES: Partial<Record<string, UserRole[]>> = {
    '/admin/broadcast': ['SUPER'],
    '/catalogo-servicos': ['SUPER'],
    '/configuracoes/api-integra-contador': ['SUPER'],
    '/configuracoes/integra-contador-servicos': ['SUPER'],
    '/configuracoes/integra-contador-tabela-consumo': ['SUPER'],
    '/matriz-conformidade-fiscal': ['SUPER'],
    '/representantes': ['SUPER'],
    '/monitor': ['SUPER', 'ADMIN'],
    '/tenants': ['SUPER', 'REPRESENTANTE'],
    '/usuarios': ['SUPER', 'ADMIN'],
    '/utilitarios/hardware-manager': ['SUPER'],
    '/utilitarios/monitoramento-global': ['SUPER'],
};

export function routeUsesAuthenticatedDashboard(pathname: string): boolean {
    if (AUTH_REQUIRED_ROUTES.has(pathname)) {
        return true;
    }
    return Boolean(ROLE_RESTRICTED_ROUTES[pathname]);
}
