/* eslint-disable @next/next/no-img-element */
import React from 'react';
import Link from 'next/link';
import { Button } from 'primereact/button';
import AppConfig from '../../../layout/AppConfig';
import { Page } from '../../../types/types';

const SobreVecxPage: Page = () => {
  return (
    <div className="surface-ground min-h-screen">
      <div className="py-6 px-4 mx-auto" style={{ maxWidth: '52rem' }}>
        <div className="flex flex-wrap align-items-center justify-content-between gap-3 mb-4">
          <Link href="/pages/landing">
            <Button type="button" label="Voltar ao início" icon="pi pi-arrow-left" className="p-button-text" />
          </Link>
          <Link href="/" className="text-600 hover:text-primary">
            Ir ao sistema
          </Link>
        </div>

        <div className="surface-card border-round-xl border-1 border-solid border-200 shadow-1 p-4 md:p-5">
          <h1 className="text-900 text-3xl md:text-4xl font-medium mt-0 mb-4">Sobre o VecX</h1>

          <p className="text-700 line-height-3 text-lg mb-4">
            VEC é a união de Valéria, Eduardo e Carlos, três irmãos que são sócios da VecX.
          </p>

          <h2 className="text-900 text-2xl font-medium mt-5 mb-3">Mas e o &quot;X&quot;?</h2>

          <section className="mb-5">
            <h3 className="text-900 text-xl font-medium mt-0 mb-2">1. Multiplicação de resultados</h3>
            <p className="text-700 line-height-3 mb-0">
              Na matemática, o &quot;X&quot; é o símbolo da multiplicação. Isso simboliza que o VecX não é apenas uma ferramenta de registro, mas um
              multiplicador da eficiência, do lucro e do tempo para os seus clientes. Ele pega o &quot;VEC&quot; (a base de vocês) e multiplica os
              resultados contábeis.
            </p>
          </section>

          <section className="mb-5">
            <h3 className="text-900 text-xl font-medium mt-0 mb-2">2. O ponto de encontro (cross-functional)</h3>
            <p className="text-700 line-height-3 mb-2">
              O &quot;X&quot; é o ponto onde duas linhas se cruzam. Isso representa o VecX como o ponto de encontro central onde tudo se integra:
            </p>
            <ul className="text-700 line-height-3 pl-3 my-0">
              <li className="mb-2">
                O encontro entre a tecnologia avançada (Go, PostgreSQL) e o conhecimento humano (VEC).
              </li>
              <li className="mb-0">
                O ponto onde os dados do cliente se cruzam com as obrigações fiscais (NF-e, LCDPR), gerando clareza contábil.
              </li>
            </ul>
          </section>

          <section className="mb-5">
            <h3 className="text-900 text-xl font-medium mt-0 mb-2">3. A nova geração (exponencial / experiência)</h3>
            <p className="text-700 line-height-3 mb-2">
              Na linguagem corporativa moderna, o &quot;X&quot; é frequentemente usado para denotar o &quot;próximo nível&quot; ou
              &quot;exponencial&quot;:
            </p>
            <ul className="text-700 line-height-3 pl-3 my-0">
              <li className="mb-2">
                <strong className="text-900">Performance exponencial:</strong> a promessa de levar a contabilidade a um nível nunca antes visto.
              </li>
              <li className="mb-0">
                <strong className="text-900">Experiência do usuário (UX):</strong> o &quot;X&quot; também simboliza o foco em uma experiência simples
                e intuitiva (o &quot;X&quot; de &quot;Experience&quot;).
              </li>
            </ul>
          </section>

          <section className="mb-5">
            <h3 className="text-900 text-xl font-medium mt-0 mb-2">4. A letra da engenharia (eixo XY)</h3>
            <p className="text-700 line-height-3 mb-0">
              O design do VecX esconde um segredo matemático: enquanto o X estabelece o eixo horizontal da nossa base técnica e fundamento lógico (o
              hub), o corte sutil à direita revela um Y à esquerda. Esse Y representa o eixo da verticalidade e do crescimento dos nossos clientes. É
              a união da engenharia (X) com a ascensão (Y), provando que para cada dado processado, existe um resultado que se eleva.
            </p>
          </section>
        </div>
      </div>
    </div>
  );
};

SobreVecxPage.getLayout = function getLayout(page) {
  return (
    <React.Fragment>
      {page}
      <AppConfig simple />
    </React.Fragment>
  );
};

export default SobreVecxPage;
