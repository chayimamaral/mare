
labels: enhancement, implementação

---
### Resumo

Manter o cadastro de municípios (cidades) com listagem paginada sob demanda (lazy loading), filtros na grade, vínculo obrigatório com estado (UF), inclusão/edição/exclusão e exportação CSV, consumindo as APIs `/api/cidades`, `/api/cidade` e catálogo de UFs para o diálogo.

### Atores

-  Usuário autenticado com acesso à rota (protegida por sessão/cookie).
- Sistema: carrega estados via `EstadoService`; persiste via `MunicipioService`.

### Contexto e Objetivo

O município é referência geográfica usada em empresas, feriados, rotinas, obrigações etc. A tela centraliza o CRUD. 
 `frontend/pages/municipios/index.tsx`, alinhado ao serviço `frontend/services/cruds/MunicipioService.ts` (endpoints `cidades` / `cidade` / `cidadeslite` onde aplicável).

### Fluxo principal

1. O sistema apresenta a página **Cadastro de Municípios** após **SSR** com `canSSRAuth` e chamada a **`GET /api/registro`**; falha redireciona para `/`. **[RN — autenticação/registro]**

2. ## Fluxos alternativos / extensões- 

**[FA1] Pick-list UF:** Dropdown PrimeReact, `options={estados}`, `optionLabel="nome"`, `dataKey="id"`, valor em estado separado `estado` (não só `municipio.uf`).- 

**[FA2] Lista “lite”:** Outras telas usam `getMunicipiosLite` / `GET /api/cidadeslite`; esta tela usa lista completa lazy em `/api/cidades`.

- **[FA3] Filtros no estado inicial:** `lazyState.filters` inclui chaves `nome`, `codigo`, `municipio` (contains), mas a busca rápida do cabeçalho e `handleClear` mexem majoritariamente em **`nome`** — **revisar consistência** backend/frontend se código/“municipio” forem requisito.

### Regras de negócio (RN)

- **[RN1] Lazy list:** Parâmetros espelham o objeto `lazyState` (first, rows, page, sortField, sortOrder, filters).
- **[RN2] Ordenação:** `sortOrder` 1 ou -1 mapeado para o DataTable conforme `lazyState.sortOrder`.
- **[RN3] Busca cabeçalho:** Somente Enter; filtro por nome (contains). Tooltip orienta o usuário.
- **[RN4] Criação:** Código sempre persistido em **maiúsculas** no payload (`toUpperCase()` no save).
- **[RN5] UF obrigatória:** Se não houver `estado.id`, Toast “Estado obrigatório” e não envia.
- **[RN6] Exclusão:** Só prossegue se `municipio.nome` preenchido e `municipio.id` presente.
- **[RN7] Modelo de dados na UI:** Entidade local como `Vec.Cidade` (`id`, `nome`, `codigo`, `ufid`, `uf?: { nome, ... }`); coluna “Estado” exibe `rowData.uf?.nome`.

### Tratamento de exceções (E)

- **[E1] Erro de API (create/update/delete):** Toast de erro com `apiErr` — prioriza `AxiosError.response.data.error`, senão `message`, senão texto fixo.

- **[E2] SSR:** Falha em `/api/registro` no `getServerSideProps` → redirect `/` (exceto se o template de projeto evoluir para outro destino).

- **[E3] Validação de formulário:** `submitted` exibe `p-invalid` e mensagens para nome/código vazios no diálogo; UF usa Toast específico se faltar.

### Modelagem e impacto técnico

**Banco de dados (PostgreSQL)**  
Documentar no backend (tabelas `cidade`/`municipio` e FK para `estado`/UF conforme migrações reais do repositório). O frontend envia `ufid` como vínculo.

**Backend (Go)**  

- **`GET /api/cidades`**: lista paginada/filtrada conforme query derivada de `lazyEvent`.  

- **`POST/PUT/DELETE /api/cidade`**: corpo `{ params: { id?, nome, codigo, ufid, ... } }` alinhado ao handler existente.  

- **`GET /api/registro`**: pré-condição SSR desta rota.  

- Garantir regras de unicidade (ex.: código+UF) e integridade referencial na exclusão — **validar no serviço/repositório** (não só no frontend).

**Frontend (React / PrimeReact)**  

- **Página:** `frontend/pages/municipios/index.tsx` — componente `Municipios`, estado `lazyState`, `DataTable` lazy, `Dialog` CRUD, `Toast`, `Toolbar`.  

- **Serviço:** `frontend/services/cruds/MunicipioService.ts` — `getMunicipios`, `createMunicipio`, `updateMunicipio`, `deleteMunicipio`, `getMunicipiosLite`.  

- **Tipos:** `Vec.Cidade` / `Vec.Estado` em `frontend/types/vec.d.ts`.  

- **Estados:** `EstadoService().getUFCidade` para options do Dropdown.

### Critérios de aceite

- [ ] Listagem lazy retorna dados e `totalRecords` coerentes com o backend após paginar/ordenar/filtrar.
- [ ] Criar município com nome, código e UF válidos persiste e aparece na lista; código gravado conforme regra de maiúsculas acordada com API.
- [ ] Editar mantém UF correta após carregar estados (inclusive race entre `editMunicipio` e `loadLazyEstados`).
- [ ] Excluir remove o registro ou exibe erro do backend (ex.: FK) de forma clara.
- [ ] Usuário não autenticado ou falha de registro não permanece na página (redirect).
- [ ] Exportar CSV gera arquivo a partir do DataTable.
- [ ] (Opcional / dívida técnica) Unificar filtros `codigo` / `municipio` no `lazyState` com o que o backend realmente interpreta; alinhar `handleClear` com o estado completo dos filtros.

### Observações e documentos de apoio

- **Autor (template):** Carlos Amaral  
- **Documento:** `docs/requisitos/EF-50-municipios.md` (sugerido)  
- **https://github.com/chayimamaral/vecontab/issues/50
- **Dependências:** Cadastro de Estados (UF); consumidores: Empresas, Feriados, Rotinas, Obrigações.  
- **Referência de template:** [.github/ISSUE_TEMPLATE/requisito-complexo.md](https://github.com/chayimamaral/vecontab/blob/main/.github/ISSUE_TEMPLATE/requisito-complexo.md)
