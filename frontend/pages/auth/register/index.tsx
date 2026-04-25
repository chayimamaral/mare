/* eslint-disable @next/next/no-img-element */

import React, { useContext, useState, useRef, useEffect } from 'react';
import AppConfig from '../../../layout/AppConfig';
import { Button } from 'primereact/button';
import { Password } from 'primereact/password';
import { LayoutContext } from '../../../layout/context/layoutcontext';
import { InputText } from 'primereact/inputtext';
import { classNames } from 'primereact/utils';
import { Page } from '../../../types/types';

import AuthContext from "../../../components/context/AuthContext";
import { Toast } from 'primereact/toast';
import Link from 'next/link';

export const RegisterPage: Page = () => {
  const [password, setPassword] = useState('');
  const [email, setEmail] = useState('');
  const [nome, setNome] = useState('');
  const [empresaNome, setEmpresaNome] = useState('');
  const { layoutConfig } = useContext(LayoutContext);
  const { signUp } = useContext(AuthContext);
  const toast = useRef<Toast>(null);
  const [isInvalid, setIsInvalid] = useState(false);

  const containerClassName = classNames('surface-ground flex align-items-center justify-content-center min-h-screen min-w-screen overflow-hidden', { 'p-input-filled': layoutConfig.inputStyle === 'filled' });

  useEffect(() => {
    const clearForm = () => {
      setNome('');
      setEmail('');
      setEmpresaNome('');
      setPassword('');
      setIsInvalid(false);
    };

    clearForm();

    const timeoutId = window.setTimeout(clearForm, 50);
    return () => window.clearTimeout(timeoutId);
  }, []);

  async function handleRegister() {

    if (!nome || !email || !password || !empresaNome) {
      setIsInvalid(true);
      toast?.current?.show({ severity: 'error', summary: 'Erro', detail: 'Preencha todos os campos!', life: 3000 });
      return;
    }

    try {
      const created = await signUp({
        nome,
        email,
        password,
        empresa_nome: empresaNome.trim(),
      })

      toast?.current?.show({
        severity: 'success',
        summary: 'Sucesso',
        detail: `Conta criada. Tenant: ${created.tenantid || 'n/d'} | Schema: ${created.tenant_schema || 'n/d'}`,
        life: 5000
      });
    } catch (err) {
      setIsInvalid(true);
      const message = err instanceof Error ? err.message : 'Erro ao criar conta';
      toast?.current?.show({ severity: 'error', summary: 'Erro', detail: message, life: 3500 });
      return;
    }
  }

  return (
    <div className={containerClassName}>
      <div className="flex flex-column align-items-center justify-content-center">
        <div style={{ borderRadius: '56px', padding: '0.3rem', background: 'linear-gradient(180deg, var(--primary-color) 10%, rgba(33, 150, 243, 0) 30%)' }}>
          <div className="w-full surface-card py-8 px-5 sm:px-8" style={{ borderRadius: '53px' }}>
            <div className="text-center mb-5">
              <img src="/vecx_logo.svg" alt="Vecontab logo" className="mb-5 w-16rem flex-shrink-0" />
              <div className="text-900 text-3xl font-medium mb-3">

              </div>
              <span className="text-600 font-medium">Crie sua conta no VECX!</span>
            </div>
            <form autoComplete="off" onSubmit={(e) => e.preventDefault()}>
              {/* Campos isca para o navegador não injetar credenciais do ultimo login no form de cadastro */}
              <input type="text" name="username" autoComplete="username" style={{ display: 'none' }} tabIndex={-1} />
              <input type="password" name="current-password" autoComplete="current-password" style={{ display: 'none' }} tabIndex={-1} />
              <label htmlFor="register-user-name" className="block text-900 text-xl font-medium mb-2">
                Nome do Usuário
              </label>
              <InputText
                id="register-user-name"
                name="nome"
                value={nome}
                onChange={(e) => setNome(e.target.value)}
                type="text"
                autoComplete="new-password"
                placeholder="Nome"
                className={`w-full md:w-30rem mb-5 ${isInvalid ? 'p-invalid' : ''}`}
                style={{ padding: '1rem' }}
              />

              <label htmlFor="register-email" className="block text-900 text-xl font-medium mb-2">
                Email (login)
              </label>
              <InputText
                id="register-email"
                name="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                type="email"
                autoComplete="new-password"
                placeholder="Email"
                className={`w-full md:w-30rem mb-5 ${isInvalid ? 'p-invalid' : ''}`}
                style={{ padding: '1rem' }}
              />

              <label htmlFor="register-company-name" className="block text-900 text-xl font-medium mb-2">
                Empresa/Escritório
              </label>
              <InputText
                id="register-company-name"
                name="empresa_nome"
                value={empresaNome}
                onChange={(e) => setEmpresaNome(e.target.value)}
                type="text"
                autoComplete="new-password"
                placeholder="ex: Vec Contabilidade"
                className={`w-full md:w-30rem mb-5 ${isInvalid ? 'p-invalid' : ''}`}
                style={{ padding: '1rem' }}
              />

              <label htmlFor="register-password" className="block text-900 font-medium text-xl mb-2">
                Senha
              </label>
              <Password
                inputId="register-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Senha"
                toggleMask
                pt={{
                  input: {
                    name: 'register_new_password',
                    autoComplete: 'new-password',
                  },
                }}
                className={`w-full md:w-30rem mb-5 ${isInvalid ? 'p-invalid' : ''}`}
                inputClassName="w-full p-3 md:w-30rem"
              ></Password>

              <div className="flex align-items-center justify-content-between mb-5 gap-5">
                <span className="font-medium no-underline ml-2 text-center" style={{ color: 'var(--primary-color)' }}>
                  Já possui conta ? Faça login <Link href='/auth/login'><strong>aqui</strong></Link>
                </span>
              </div>

              <div className="card flex justify-content-center">
                <Toast ref={toast} />
                <Button type="submit" label="Acessar" className="w-full p-3 text-xl" onClick={handleRegister}></Button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
};

RegisterPage.getLayout = function getLayout(page) {
  return (
    <React.Fragment>
      {page}
      <AppConfig simple />
    </React.Fragment>
  );
};

export default RegisterPage;