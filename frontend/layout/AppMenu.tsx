/* eslint-disable @next/next/no-img-element */

import React, { useEffect, useMemo, useState } from 'react';
import AppMenuitem from './AppMenuitem';
import { MenuProvider } from './context/menucontext';
import { AppMenuItem } from '../types/types';
import setupAPIClient from '../components/api/api';
// Onde está o nome: logoIntegra
// Onde está o local: '../assets/logo_integracontador_limpo.avif'
import logoIntegra from '../public/logo_integracontador.avif';
import { sessionAllowsFeature } from '../constants/featureAccess';

type MenuSession = {
  role: string | null;
  /** Ausente no JSON = token legado (não restringe por módulo). Array (vazio ou não) = claim presente no JWT. */
  featureSlugs: string[] | undefined;
};

const AppMenu = () => {
  const [menuSession, setMenuSession] = useState<MenuSession | null>(null);
  /** Evita divergência SSR/cliente: até montar no browser, o menu não usa o papel (igual ao HTML do servidor). */
  const [menuMounted, setMenuMounted] = useState(false);

  useEffect(() => {
    setMenuMounted(true);
  }, []);

  useEffect(() => {
    const api = setupAPIClient(undefined);
    api
      .get('/api/usuariorole')
      .then((r) => {
        const logado = r.data?.logado as { role?: string; feature_slugs?: string[] } | undefined;
        const hasFeatureKey = Boolean(logado && Object.prototype.hasOwnProperty.call(logado, 'feature_slugs'));
        setMenuSession({
          role: logado?.role ?? null,
          featureSlugs: hasFeatureKey ? logado?.feature_slugs ?? [] : undefined,
        });
      })
      .catch(() => setMenuSession({ role: null, featureSlugs: undefined }));
  }, []);

  const roleForMenu = menuMounted ? menuSession?.role ?? null : null;

  const model: AppMenuItem[] = useMemo(() => {
    const podeGerenciarUsuarios = roleForMenu === 'ADMIN' || roleForMenu === 'SUPER';
    const podeVerMonitor = roleForMenu === 'ADMIN' || roleForMenu === 'SUPER';
    const isSuper = roleForMenu === 'SUPER';
    const podeConsultarNFe =
      roleForMenu === 'SUPER' ||
      roleForMenu === 'ADMIN' ||
      roleForMenu === 'USER' ||
      roleForMenu === 'REPRESENTANTE';

    const hasFeature = (slug: string): boolean => {
      if (!menuMounted || menuSession === null) {
        return true;
      }
      return sessionAllowsFeature(menuSession.role, menuSession.featureSlugs, slug);
    };

    return [
      {
        label: 'Home',
        items: [
          { label: 'Dashboard', icon: 'pi pi-fw pi-home', to: '/' },
          {
            label: 'Compromissos Fiscais/Tributários',
            icon: 'pi pi-fw pi-list',
            visible: hasFeature('compromissos'),
            items: [
              {
                label: 'Compromissos por empresas',
                icon: 'pi pi-fw pi-list',
                to: '/compromissos-empresas',
              },
              {
                label: 'Compromissos por natureza',
                icon: 'pi pi-fw pi-sitemap',
                to: '/compromissos-por-natureza',
              },
              {
                label: 'Compromissos (visão corrida)',
                icon: 'pi pi-fw pi-table',
                to: '/compromissos-visao',
              },
            ],
          },
          {
            label: 'Fluxos de Processos',
            icon: 'pi pi-fw pi-calendar',
            visible: hasFeature('compromissos'),
            items: [
              {
                label: 'Fluxo em Árvore',
                icon: 'pi pi-fw pi-sitemap',
                to: '/agenda-arvore',
              },
              {
                label: 'Fluxo em Agenda (calendário)',
                icon: 'pi pi-fw pi-calendar',
                to: '/agenda',
              },
            ],
          },
          {
            label: 'Manutenção de Empresas',
            icon: 'pi pi-fw pi-table',
            visible: hasFeature('core') || hasFeature('nfe'),
            items: [
              {
                label: 'Manutenção de NFe',
                icon: 'pi pi-fw pi-file',
                to: '/nfe/manutencao',
                visible: hasFeature('nfe'),
              },
              {
                label: 'Sincronização de NFe',
                icon: 'pi pi-fw pi-sync',
                to: '/nfe/sincronizacao',
                visible: hasFeature('nfe'),
              },
              {
                label: 'Manutenção de Empresas',
                icon: 'pi pi-fw pi-building',
                to: '/empresas',
              },
            ],
          },
          {
            label: 'Manutenção de Cliente PF (IRPF)',
            icon: 'pi pi-fw pi-user',
            to: '/cliente-pf',
            visible: hasFeature('core'),
          },
          {
            label: 'Caixa Postal',
            icon: 'pi pi-fw pi-envelope',
            to: '/caixa-postal',
            visible: hasFeature('caixa_postal'),
          },
        ],
      },
      {
        label: 'Operações',
        visible: hasFeature('core'),
        items: [
          {
            label: 'Cadastros',
            icon: 'pi pi-fw pi-database',
            items: [
              {
                label: 'Cadastros Básicos',
                icon: 'pi pi-fw pi-bookmark',
                items: [
                  {
                    label: 'Clientes',
                    icon: 'pi pi-fw pi-id-card',
                    to: '/clientes',
                  },
                  {
                    label: 'Feriados',
                    icon: 'pi pi-fw pi-table',
                    to: '/feriados',
                  },
                  {
                    label: 'Municípios',
                    icon: 'pi pi-fw pi-building',
                    to: '/municipios',
                  },
                  {
                    label: 'Estados',
                    icon: 'pi pi-fw pi-flag',
                    to: '/estados',
                  },
                ],
              },
              {
                label: 'Cadastros Operacionais',
                icon: 'pi pi-fw pi-sitemap',
                items: [
                  {
                    label: 'Configurações Fiscais',
                    icon: 'pi pi-fw pi-table',
                    items: [

                      {
                        label: 'Regras de Obrigações',
                        icon: 'pi pi-fw pi-money-bill',
                        to: '/obrigacoes',
                      },

                    ],
                  },
                  {

                    label: 'Processos para Empresas', // O item pai agora agrupa os dois
                    icon: 'pi pi-fw pi-briefcase',
                    items: [
                      {
                        label: 'Processos',
                        icon: 'pi pi-fw pi-list',
                        to: '/rotinas', // Mantive o path atual para não quebrar seus links
                      },
                      {
                        label: 'Etapas',
                        icon: 'pi pi-fw pi-check-square',
                        to: '/passos',
                      },
                    ]
                  },
                  {
                    label: 'Processos PF (IRPF / Carnê-Leão)',
                    icon: 'pi pi-fw pi-user',
                    to: '/rotinas-pf',
                  },
                ],
              },
              {
                label: 'Cadastros Contábeis',
                icon: 'pi pi-fw pi-sitemap',
                items: [
                  {
                    label: 'Enquadramento Jurídico',
                    icon: 'pi pi-fw pi-table',
                    to: '/tipoempresa',
                  },
                  {
                    label: 'Regime tributário',
                    icon: 'pi pi-fw pi-percentage',
                    to: '/regimes-tributarios',
                  },
                  {
                    label: 'CNAE',
                    icon: 'pi pi-fw pi-table',
                    to: '/cnae',
                  },
                  {
                    label: 'Salário mínimo nacional',
                    icon: 'pi pi-fw pi-money-bill',
                    to: '/salario-minimo',
                  },
                ],
              },
            ],
          },
        ],
      },
      {
        label: 'Diversos',
        icon: 'pi pi-fw pi-briefcase',
        to: '/pages',
        items: [
          {
            label: 'Usuários',
            icon: 'pi pi-fw pi-user',
            visible: podeGerenciarUsuarios,
            items: [
              {
                label: 'Usuários',
                icon: 'pi pi-fw pi-users',
                to: '/usuarios',
              },
            ],
          },
          {
            label: 'Monitor',
            icon: 'pi pi-fw pi-chart-line',
            visible: podeVerMonitor && hasFeature('monitor'),
            items: [
              {
                label: 'Operações',
                icon: 'pi pi-fw pi-list',
                to: '/monitor',
              },
            ],
          },
          {
            label: 'Configurações APIs',
            icon: 'pi pi-fw pi-cog',
            visible: isSuper,
            items: [
              {
                label: 'Integra Contador - Serpro',
                iconSrc: '/microservice-icon.svg',
                items: [
                  {
                    label: 'Chave de Autenticação',
                    icon: 'pi pi-fw pi-key',
                    to: '/configuracoes/api-integra-contador',
                    visible: isSuper,
                  },
                  {
                    label: 'Execução de Serviços',
                    icon: 'pi pi-fw pi-play',
                    to: '/configuracoes/integra-contador-servicos',
                    visible: isSuper,
                  },
                  {
                    label: 'Catálogo de Serviços',
                    icon: 'pi pi-fw pi-sitemap',
                    to: '/catalogo-servicos',
                    visible: isSuper,
                  },
                  {
                    label: 'Matriz de Conformidade Fiscal',
                    icon: 'pi pi-fw pi-check-square',
                    to: '/matriz-conformidade-fiscal',
                    visible: isSuper,
                  },
                  {
                    label: 'Tabela de Preços de Consumo - Integra Contador',
                    icon: 'pi pi-fw pi-wallet',
                    to: '/configuracoes/integra-contador-tabela-consumo',
                    visible: isSuper,
                  },
                  {
                    label: 'NFe - Serpro',
                    icon: 'pi pi-fw pi-cog',
                    to: '/configuracoes/api-integra-contador',
                    visible: isSuper,
                  },

                ],
              },
            ],
          },
          {
            label: 'Configurações do Tenant',
            icon: 'pi pi-fw pi-building',
            visible: roleForMenu === 'SUPER' || roleForMenu === 'ADMIN',
            items: [
              {
                label: 'Integra Contador - Tenant',
                iconSrc: '/microservice-icon.svg',
                items: [
                  {
                    label: 'Geração de Guias',
                    icon: 'pi pi-fw pi-file',
                    to: '/configuracoes/geracao-guias',
                    visible: (roleForMenu === 'SUPER' || roleForMenu === 'ADMIN') && hasFeature('integra_contador'),
                  },
                  {
                    label: 'Certificado Digital',
                    icon: 'pi pi-fw pi-shield',
                    to: '/configuracoes/certificado-digital',
                    visible: (roleForMenu === 'SUPER' || roleForMenu === 'ADMIN') && hasFeature('integra_contador'),
                  },
                  {
                    label: 'Consulta NFe',
                    icon: 'pi pi-fw pi-file-check',
                    to: '/nfe/consulta',
                    visible: podeConsultarNFe && hasFeature('nfe'),
                  },
                ],
              },
            ],
          },
          {
            label: 'Gestão de Tenants',
            icon: 'pi pi-fw pi-server',
            visible: roleForMenu === 'SUPER' || roleForMenu === 'REPRESENTANTE',
            items: [
              {
                label: 'Manutenção',
                icon: 'pi pi-fw pi-table',
                to: '/tenants',
              },
              {
                label: 'Representantes comerciais',
                icon: 'pi pi-fw pi-users',
                to: '/representantes',
                visible: roleForMenu === 'SUPER',
              },
              {
                label: 'Broadcast & Sessões',
                icon: 'pi pi-fw pi-megaphone',
                to: '/admin/broadcast',
                visible: roleForMenu === 'SUPER',
              },
            ],
          },
          {
            label: 'Sobre',
            icon: 'pi pi-fw pi-pencil',
            to: '/pages/landing',
          },
          {
            label: 'Utilitários',
            icon: 'pi pi-fw pi-wrench',
            visible: isSuper,
            items: [
              {
                label: 'Hardware Manager',
                icon: 'pi pi-fw pi-desktop',
                to: '/utilitarios/hardware-manager',
                visible: isSuper,
              },
            ],
          },
        ],
      },
    ];
  }, [roleForMenu, menuMounted, menuSession]);

  return (
    <MenuProvider>
      <ul className="layout-menu">
        {model.map((item, i) => {
          return !item?.seperator ? (
            <AppMenuitem item={item} root={true} index={i} key={`root-${i}-${item.label}`} />
          ) : (
            <li className="menu-separator"></li>
          );
        })}
      </ul>
    </MenuProvider>
  );
};

export default AppMenu;
