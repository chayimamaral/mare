# backendgo

Primeira base do porte do backend Node/TypeScript para Go.

## O que ja esta portado

- bootstrap HTTP em Go
- conexao PostgreSQL com pool
- autenticacao JWT
- login em `/session` e `/api/session`
- usuario logado em `/me`, `/usuariorole` e `/usuariotenant`
- modulo `estado` em `/estados`, `/estado`, `/deleteestado` e `/ufscidade`
- modulo `cidade` em `/cidades`, `/cidade` (POST/PUT/DELETE) e `/cidadeslite`
- modulo `tenant` em `/tenant` (POST/GET/PUT) e `/tenants`
- usuario em `/usuarios` (list) e `/usuario` (create)
- modulo `tipoempresa` em `/tiposempresa`, `/tipoempresa`, `/deletetipoempresa` e `/tiposempresalite`
- modulo `passo` em `/passos`, `/passo`, `/deletepasso`, `/getPassoById` e `/passosporcidade`
- modulo `grupopassos` em `/grupopassos`, `/grupopasso`, `/deletegrupopasso` e `/getgrupopassobyid`
- modulo `feriado` em `/feriados`, `/feriado` e `/deleteferiado`
- modulo `empresa` em `/empresas`, `/empresa`, `/updateempresa`, `/deleteempresa` e `/iniciarprocesso`
- modulo `cnae` em `/cnaes`, `/cnae`, `/deletecnae`, `/cnaelite` e `/validacnae`
- modulo `agenda` em `/agendalist` e `/agendadetalhes`
- modulo `rotina` em `/rotinas`, `/rotina`, `/deleterotina`, `/rotinaitens`, `/rotinaitemcreate`, `/rotinaitemupdate`, `/rotinaitemdelete`, `/listrotinas`, `/listrotinaslite`, `/listrotinaitensselected`, `/salvarselecao` e `/removepassoselecionado`
- modulo `registro` em `/registro` (POST/GET/PUT)
- modulo `node` em `/node`, `/family` e `/recurso`
- hardening inicial de tenant em `empresa` e `agenda`, sempre usando o tenant do JWT em vez de confiar em query/body
- hardening de role em rotas administrativas de escrita (ADMIN/SUPER)
- hardening de criacao de usuario: apenas ADMIN/SUPER, com validacao de role alvo (`USER|ADMIN|SUPER`) e bloqueio de elevacao para `SUPER` por ADMIN
- reducao progressiva de `map[string]any` com respostas internas tipadas nos modulos `agenda`, `empresa`, `node`, `tenant`, `auth`, `grupopassos`, `feriado`, `estado`, `rotina`, `registro` (service + repository), `user`, `cidade`, `passo`, `cnae` (service + repository), e `tipoempresa` (service + repository lite)

## Variaveis de ambiente

- `SERVER_PORT`: porta HTTP, padrao `3333`
- `PG_URL`: string de conexao PostgreSQL
- `JWT_SECRET`: segredo do token JWT
- `PG_SSL_ROOT_CERT`: caminho do certificado CA
- `PG_SSL_INSECURE`: `true` ou `false`

## Rodando localmente

```bash
go mod tidy
go run ./cmd/api
```

## Estrategia de migracao

1. manter o contrato HTTP existente
2. portar modulo por modulo, com SQL explicito
3. trocar stubs `501` por handlers reais
4. validar o frontend contra o backend Go antes de desligar o Node

## Arquitetura

O backend Go foi organizado em camadas. A divisao principal existe para separar responsabilidades e reduzir acoplamento entre HTTP, regra de negocio e acesso ao banco.

Fluxo principal:

`HTTP -> handler -> service -> repository -> PostgreSQL`

### Handler

Fica em `internal/httpapi/handlers`.

Responsabilidades:

- receber request HTTP
- ler `query`, `params` e `body`
- extrair informacoes do JWT pelo `context`
- validar entrada basica
- chamar o service
- transformar erro em resposta HTTP

O handler nao deve conhecer SQL. Ele tambem nao deve concentrar regra de negocio que pertence ao dominio.

### Service

Fica em `internal/service`.

Responsabilidades:

- concentrar regra de negocio
- aplicar autorizacao funcional quando isso nao for apenas transporte HTTP
- orquestrar mais de um repository quando necessario
- definir o contrato de uso da funcionalidade

O service existe para impedir que a regra de negocio fique espalhada entre controller/handler e SQL.

### Repository

Fica em `internal/repository`.

Responsabilidades:

- encapsular acesso ao PostgreSQL
- montar SQL explicitamente
- executar queries e transacoes
- mapear resultados do banco

O repository nao deve conhecer detalhes de HTTP. Ele recebe dados ja preparados pela camada superior.

## Por que separar em 3 pastas

A separacao parece redundante no inicio, mas resolve problemas diferentes:

- manutencao: fica claro onde mexer quando o problema esta no HTTP, na regra ou no banco
- teste: permite validar regra de negocio sem depender diretamente de request HTTP
- evolucao: trocar detalhes do banco ou do transporte afeta menos arquivos
- seguranca: fica mais facil centralizar regras como tenant do JWT e checks de role

Na pratica, a redundancia ruim nao esta nas pastas. Ela aparece quando tipos muito parecidos comecam a ser copiados entre as camadas sem necessidade.

## Sobre `model`

Nao foi criada uma pasta `model` generica de proposito.

Em Go, uma pasta `model` costuma virar um deposito misturando:

- entidade de negocio
- payload HTTP
- filtro de listagem
- resultado de query
- resposta de API

Quando isso acontece, uma struct que deveria pertencer a uma camada passa a ser usada por todas, e o acoplamento aumenta.

Por isso, neste porte, os tipos ficaram proximos da camada onde sao usados:

- structs de request/response perto dos handlers
- inputs e orchestration perto dos services
- filtros e structs de persistencia perto dos repositories

## O que esta bom e o que ainda pode melhorar

O desenho atual favorece a migracao incremental com baixo risco, porque preserva o contrato do backend original e deixa o SQL visivel.

Ao mesmo tempo, ainda ha pontos para evoluir:

- reduzir uso de `map[string]any`
- reduzir duplicacao de tipos muito parecidos entre handler/service/repository
- introduzir entidades compartilhadas apenas onde houver ganho real

## Direcao futura para modelos de dominio

Se o projeto evoluir alem da migracao, a recomendacao e criar uma camada de dominio explicita, por exemplo:

- `internal/domain/empresa.go`
- `internal/domain/user.go`
- `internal/domain/tenant.go`

Essa camada deve concentrar apenas entidades centrais do negocio.

Ela nao deve substituir DTOs HTTP nem structs especificas de query. A ideia e manter:

- handler com payloads HTTP
- service com contratos de caso de uso
- repository com tipos de persistencia
- domain com entidades realmente compartilhadas

## Matriz de permissoes por rota

Legenda:

- `PUBLICO`: sem token
- `AUTH`: exige token JWT valido
- `ADMIN/SUPER`: exige token + role `ADMIN` ou `SUPER`

### Autenticacao e cadastro inicial

- `POST /session`: `PUBLICO`
- `POST /registro`: `PUBLICO`
- `GET /registro`: `AUTH`
- `PUT /registro`: `AUTH`

### Tenant e usuario

- `POST /tenant`: `PUBLICO`
- `GET /tenant`: `AUTH`
- `PUT /tenant`: `ADMIN/SUPER`
- `GET /tenants`: `AUTH` (retorno depende do role)
- `GET /me`: `AUTH`
- `GET /usuarios`: `AUTH` (service restringe para `ADMIN`/`SUPER`)
- `GET /usuariorole`: `AUTH`
- `GET /usuariotenant`: `AUTH`
- `POST /usuario`: `ADMIN/SUPER`

### Cadastros administrativos

- `POST/PUT/DELETE /cidade`: `ADMIN/SUPER`
- `GET /cidades` e `GET /cidadeslite`: `AUTH`
- `POST /estado`, `PUT /estado`, `PUT /deleteestado`: `ADMIN/SUPER`
- `GET /estados` e `GET /ufscidade`: `AUTH`
- `POST /tipoempresa`, `PUT /tipoempresa`, `PUT /deletetipoempresa`: `ADMIN/SUPER`
- `GET /tiposempresa` e `GET /tiposempresalite`: `AUTH`
- `POST /passo`, `PUT /passo`, `PUT /deletepasso`: `ADMIN/SUPER`
- `GET /passos`, `GET /getPassoById`, `GET /passosporcidade`: `AUTH`
- `POST /grupopassos`, `PUT /grupopasso`, `PUT /deletegrupopasso`: `ADMIN/SUPER`
- `GET /grupopassos`, `GET /getgrupopassobyid`: `AUTH`

### Rotinas, feriados, empresas, cnae

- `POST /rotina`, `PUT /rotina`, `PUT /deleterotina`: `ADMIN/SUPER`
- `GET /rotinas`, `GET /rotinaitens`, `GET /listrotinas`, `GET /listrotinaslite`, `GET /listrotinaitensselected`: `AUTH`
- `GET /rotinaitemcreate`, `GET /rotinaitemupdate`, `GET /rotinaitemdelete`, `PUT /salvarselecao`, `PUT /removepassoselecionado`: `ADMIN/SUPER`
- `POST /feriado`, `PUT /feriado`, `PUT /deleteferiado`: `ADMIN/SUPER`
- `GET /feriados`: `AUTH`
- `POST /empresa`, `PUT /updateempresa`, `PUT /deleteempresa`, `PUT /iniciarprocesso`: `ADMIN/SUPER`
- `GET /empresas`: `AUTH`
- `POST /cnae`, `PUT /cnae`, `PUT /deletecnae`: `ADMIN/SUPER`
- `GET /cnaes`, `GET /cnaelite`, `POST /validacnae`: `AUTH`

### Agenda e arvore de passos

- `GET /agendalist`: `AUTH`
- `GET /agendadetalhes`: `AUTH`
- `GET /node`, `GET /family`, `GET /recurso`: `AUTH`

### Endpoints auxiliares

- `GET /healthz`: `PUBLICO`

Observacao:

- As rotas sao expostas em dois prefixos: raiz (`/`) e espelho em `/api`.
- O hardening de tenant em recursos sensiveis (`agenda` e `empresa`) usa o tenant do JWT para evitar acesso cruzado entre tenants.

## Checklist de homologacao por perfil

Use este roteiro para validar rapidamente se as permissoes estao corretas no ambiente.

### 1. Perfil USER

Esperado:

- consegue autenticar em `/session`
- consegue ler rotas `AUTH` (ex.: `/me`, `/agendalist`, `/empresas`, `/cnaes`)
- recebe `403` nas rotas `ADMIN/SUPER` (ex.: `POST /empresa`, `PUT /deletecnae`, `POST /usuario`)

Teste minimo sugerido:

1. fazer login como USER
2. chamar `GET /api/me` e confirmar `200`
3. chamar `POST /api/empresa` e confirmar `403`
4. chamar `POST /api/usuario` e confirmar `403`

### 2. Perfil ADMIN

Esperado:

- consegue ler rotas `AUTH`
- consegue operar rotas `ADMIN/SUPER` do proprio tenant
- nao consegue criar usuario com role `SUPER`
- ao criar usuario, tenant deve seguir o tenant do token (nao o tenant informado no body)

Teste minimo sugerido:

1. fazer login como ADMIN
2. chamar `POST /api/usuario` com role `USER` e confirmar `200`
3. chamar `POST /api/usuario` com role `SUPER` e confirmar `403`
4. chamar `PUT /api/updateempresa` com `id` de empresa de outro tenant e confirmar resposta sem alteracao

### 3. Perfil SUPER

Esperado:

- consegue ler rotas `AUTH`
- consegue operar rotas `ADMIN/SUPER`
- consegue criar usuario em tenant alvo (quando informado)
- consegue listar tenants em `/tenants`

Teste minimo sugerido:

1. fazer login como SUPER
2. chamar `GET /api/tenants` e confirmar retorno com lista
3. chamar `POST /api/usuario` com tenant alvo valido e confirmar `200`

### 4. Tenant isolation (smoke test)

Esperado:

- `agenda` e `empresa` nao permitem acesso cruzado entre tenants via query/body

Teste minimo sugerido:

1. com token do Tenant A, chamar `GET /api/agendadetalhes?agenda_id=<id_do_tenant_B>`
2. confirmar retorno vazio ou sem dados do Tenant B
3. com token do Tenant A, chamar `PUT /api/deleteempresa` em empresa do Tenant B
4. confirmar que nao houve alteracao no registro do Tenant B

### 5. Criterio de aceite rapido

- rotas `PUBLICO`: respondem sem token
- rotas `AUTH`: `401` sem token e `200` com token valido
- rotas `ADMIN/SUPER`: `403` para USER e `200` para ADMIN/SUPER
- nenhuma alteracao cruzada entre tenants em `agenda` e `empresa`

## Proximos passos sugeridos

1. concluir `tenant hardening` nos recursos restantes sensiveis (ownership em update/delete/detail)
2. revisar as regras ADMIN/SUPER junto ao frontend para confirmar se algum fluxo de escrita precisa de ajuste fino
3. reduzir `map[string]any` nos modulos mais estaveis
4. avaliar criacao de `internal/domain` depois da migracao estar funcional de ponta a ponta
