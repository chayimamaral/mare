# MARE Frontend - Interface em Next.js + React

Interface moderna e responsiva para o sistema MARE de gestão de contabilidade. Construída com Next.js 16, React 18 e PrimeReact 10.

## 📋 Índice

- [Stack Tecnológico](#stack-tecnológico)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Requisitos](#requisitos)
- [Instalação](#instalação)
- [Execução](#execução)
- [Desenvolvimento](#desenvolvimento)

## 💻 Stack Tecnológico

| Tecnologia | Versão | Propósito |
|------------|--------|----------|
| **Next.js** | 16.2.1 | Framework React Full Stack |
| **React** | 18.3.1 | Biblioteca UI |
| **TypeScript** | 5.1.3 | Type-safe development |
| **PrimeReact** | 10.9.7 | Componentes UI premium |
| **FullCalendar** | 6.1.20 | Calendário avançado |
| **Axios** | 1.4.0 | Cliente HTTP |
| **Node.js** | 22.12.0+ | Runtime JavaScript |
| **npm** | 11.11.0+ | Gerenciador de pacotes |

## 📁 Estrutura do Projeto

```
frontend/
├── pages/                    # Rotas do Next.js
│   ├── _app.tsx             # Wrapper da aplicação
│   ├── _document.tsx        # HTML base
│   ├── index.tsx            # Home page
│   ├── 404.tsx              # Página not found
│   ├── auth/
│   │   ├── login/           # Página de login
│   │   ├── register/        # Página de cadastro
│   │   └── ...
│   ├── agenda/              # Calendário
│   ├── empresas/            # Gestão de empresas
│   ├── usuarios/            # Gestão de usuários
│   ├── rotinas/             # Workflows
│   ├── estados/             # Localidades
│   ├── municipios/          # Cidades
│   ├── registro/            # Registros contábeis
│   └── ...
├── components/              # Componentes reutilizáveis
│   ├── api/
│   │   ├── api.ts          # Configuração Axios
│   │   └── apiClient.ts    # Cliente HTTP
│   ├── context/
│   │   └── AuthContext.tsx  # Context de autenticação
│   ├── errors/
│   │   └── AuthTokenError.ts
│   ├── utils/
│   │   ├── canSSRAuth.ts   # SSR com autenticação
│   │   ├── canSSRGuest.ts  # SSR para visitantes
│   │   ├── crudUtils.ts    # Utilitários CRUD
│   │   └── withServerSideProps.ts
│   ├── toolbar/            # Componentes toolbar
│   └── ...
├── layout/                  # Layout e estrutura
│   ├── AppTopbar.tsx       # Barra superior
│   ├── AppSidebar.tsx      # Menu lateral
│   ├── AppMenu.tsx         # Menu principal
│   ├── AppConfig.tsx       # Configurações
│   ├── AppFooter.tsx       # Rodapé
│   ├── layout.tsx          # Layout principal
│   └── context/
│       ├── layoutcontext.tsx
│       └── menucontext.tsx
├── services/               # Serviços da aplicação
│   ├── cruds/             # Serviços de entidade
│   │   ├── UsuarioService.ts
│   │   ├── EmpresaService.ts
│   │   ├── AgendaService.ts
│   │   └── ...
│   └── utils/             # Utilitários
├── styles/                # Estilos globais
│   └── layout/
├── types/                 # Type definitions
│   ├── layout.d.ts
│   ├── types.d.ts
│   └── vec.d.ts
├── public/                # Assets estáticos
│   ├── demo/
│   ├── layout/
│   ├── scripts/
│   └── themes/
├── package.json          # Dependências e scripts
├── tsconfig.json         # Configuração TypeScript
├── next.config.js        # Configuração Next.js
│ └── Dockerfile            # Container Docker
└── README.md             # Este arquivo
```

## 🔧 Requisitos

- **Node.js:** 22.0+ (com npm 11.0+)
- **Backend:** Rodando em `http://localhost:3333`

Verificar versões:

```bash
node --version    # v22.12.0+
npm --version     # 11.11.0+
```

## ⚙️ Instalação

### 1. Instalar dependências

```bash
npm install
```

### 2. Configurar Backend

Garantir que o backend está rodando em `http://localhost:3333`

```bash
# Terminal 1 - Backend
cd ../backend
go run ./cmd/api

# Terminal 2 - Frontend
npm run dev
```

## 🚀 Execução

### Desenvolvimento

```bash
npm run dev
```

Frontend estará em http://localhost:3000

**Características:**
- ✅ Hot reload automático
- ✅ Fast Refresh ativado
- ✅ Webpack bundler (otimizado)
- ✅ TypeScript type checking

### Build para Produção

```bash
npm run build
npm start
```

### Ver bundle size

```bash
npm run analyze
```

## 🎯 Features Principais

### 📅 Agenda Inteligente
- Calendário integrado com FullCalendar
- Sincronização em tempo real
- Eventos por tenant

### 👥 Gestão de Usuários
- Multi-tenant com isolamento
- Controle de roles (ADMIN/USER)
- Validação de permissões

### 🏢 Cadastro de Empresas
- Validação de CNAE
- Integração com localidades
- Suporte multi-tenant

### 📊 Relatórios
- Registros contábeis
- Rotinas e workflows
- Dashboard de atividades

## 🛠️ Desenvolvimento

### Autenticação

```typescript
const { signIn } = useContext(AuthContext);

await signIn({
  email: "user@example.com",
  password: "senha123"
});
```

### Lazy Loading (Paginação)

```typescript
const [lazyState, setLazyState] = useState({
  page: 1,
  rows: 20,
  sortField: 'nome',
  sortOrder: 1,
  filters: { ... }
});

<DataTable
  lazy
  onPage={setLazyState}
  totalRecords={total}
  // ...
/>
```

### Serviço CRUD

```typescript
export default function MyService() {
  const apiClient = setupAPIClient(undefined);
  
  return {
    list: (params) => apiClient.get('/api/myendpoint', { params }),
    get: (id) => apiClient.get(`/api/myendpoint/${id}`),
    create: (data) => apiClient.post('/api/myendpoint', data),
    update: (id, data) => apiClient.put(`/api/myendpoint/${id}`, data),
    delete: (id) => apiClient.delete(`/api/myendpoint/${id}`)
  };
}
```

### Proteger Página

```typescript
// Apenas autenticados
export const getServerSideProps = withAuthServerSideProps(...);

// Apenas visitantes
export const getServerSideProps = canSSRGuest(...);
```

## 🐛 Troubleshooting

### "Cannot find module 'X'"

```bash
npm install
# ou limpar cache
rm -rf node_modules package-lock.json
npm install
```

### API retorna 401

```bash
# Limpar token e fazer login novamente
localStorage.clear()
# Verificar se backend está rodando
curl http://localhost:3333/healthcheck
```

### TypeScript errors

```bash
npm run build
# ou
tsc --noEmit
```

### DOM nesting error

```typescript
// ❌ ERRADO
<Link href="/page">
  <a>Link</a>
</Link>

// ✅ CORRETO
<Link href="/page">
  <span>Link</span>
</Link>
```

## 📦 Deployment

### Vercel (Recomendado)

```bash
npm install -g vercel
vercel
```

### Self-hosted

```bash
npm run build
npm start
```

### Docker

```bash
docker build -t mare-frontend .
docker run -p 3000:3000 mare-frontend
```

## 🔐 Acessibilidade

Melhorias implementadas:
- ✅ `aria-label` em buttons sem texto
- ✅ Navegação por teclado
- ✅ Bom contraste de cores
- ✅ Estrutura semântica HTML

## 📝 Convenções

- **Componentes:** PascalCase (`AppTopbar.tsx`)
- **Funções:** camelCase (`handleSubmit`)
- **Constantes:** UPPER_CASE (`MAX_ITEMS`)
- **Arquivos:** kebab-case (`my-component.tsx`)

## 🔄 Scripts Disponíveis

```bash
npm run dev       # Iniciar dev server
npm run build     # Build para produção
npm start         # Rodar build em produção
npm run lint      # Verificar linting
npm run type-check # Type checking
npm run analyze   # Analisar bundle size
```

## 📞 Contato

Consulte [README principal](../README.md) para visão geral do projeto.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js/) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/deployment) for more details.
