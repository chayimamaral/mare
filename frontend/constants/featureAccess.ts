/** Slugs alinhados a `public.modulo_plataforma.slug` e `backend/internal/auth/feature_slugs.go`. */
export const FEATURE = {
  core: 'core',
  nfe: 'nfe',
  integra_contador: 'integra_contador',
  caixa_postal: 'caixa_postal',
  compromissos: 'compromissos',
  monitor: 'monitor',
} as const;

export type FeatureSlug = (typeof FEATURE)[keyof typeof FEATURE];

/**
 * SUPER ignora matriz. Token legado (claim ausente em usuariorole) não restringe.
 * Caso contrário, exige o slug no array (comparação case-insensitive).
 */
export function sessionAllowsFeature(
  role: string | null | undefined,
  featureSlugs: string[] | undefined,
  slug: string,
): boolean {
  const r = String(role ?? '').toUpperCase();
  if (r === 'SUPER') {
    return true;
  }
  if (featureSlugs === undefined) {
    return true;
  }
  const want = slug.toLowerCase();
  return featureSlugs.some((s) => String(s).toLowerCase() === want);
}

/** Rotas que exigem slug específico (exceto legado / SUPER). */
export const PATH_REQUIRES_FEATURE_SLUG: Record<string, string> = {
  '/agenda': FEATURE.compromissos,
  '/agenda-arvore': FEATURE.compromissos,
  '/compromissos': FEATURE.compromissos,
  '/compromissos-empresas': FEATURE.compromissos,
  '/compromissos-por-natureza': FEATURE.compromissos,
  '/compromissos-visao': FEATURE.compromissos,
  '/nfe/consulta': FEATURE.nfe,
  '/nfe/manutencao': FEATURE.nfe,
  '/nfe/sincronizacao': FEATURE.nfe,
  '/monitor': FEATURE.monitor,
  '/configuracoes/geracao-guias': FEATURE.integra_contador,
  '/configuracoes/certificado-digital': FEATURE.integra_contador,
};

/** Rotas de cadastro/operacao base que exigem `core` quando a matriz está ativa no token. */
export const PATH_REQUIRES_CORE = new Set<string>([
  '/cliente-pf',
  '/clientes',
  '/empresas',
  '/estados',
  '/feriados',
  '/grupopassos',
  '/municipios',
  '/obrigacoes',
  '/passos',
  '/regimes-tributarios',
  '/registro',
  '/rotinas',
  '/rotinas-pf',
  '/salario-minimo',
  '/tipoempresa',
  '/cnae',
  '/enquadramento-juridico',
]);
