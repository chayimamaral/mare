/* eslint-disable @next/next/no-img-element */

import React, { useContext, useEffect, useRef, useState } from 'react';
import AppConfig from '../../../layout/AppConfig';
import { Button } from 'primereact/button';
import { Password } from 'primereact/password';
import { LayoutContext } from '../../../layout/context/layoutcontext';
import { InputText } from 'primereact/inputtext';
import { classNames } from 'primereact/utils';
import { Page } from '../../../types/types';
import { Toast } from 'primereact/toast';
import { Dialog } from 'primereact/dialog';

import AuthContext from '../../../components/context/AuthContext';
import Link from 'next/link';
import { AuthError } from '../../../lib/authErrors';

export const LoginPage: Page = () => {
  const [password, setPassword] = useState('');
  const [email, setEmail] = useState('');
  const [mounted, setMounted] = useState(false);
  const { layoutConfig } = useContext(LayoutContext);
  //const { signIn } = useContext(AuthContext)
  const { signIn } = useContext(AuthContext);
  const [isInvalid, setIsInvalid] = useState(false);
  const [vecxBlockVisible, setVecxBlockVisible] = useState(false);
  const toast = useRef<Toast>(null);

  const containerClassName = classNames('surface-ground flex align-items-center justify-content-center min-h-screen min-w-screen overflow-hidden', { 'p-input-filled': layoutConfig.inputStyle === 'filled' });

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return <div className={containerClassName}></div>;
  }

  async function handleLogin() {
    if (!email || !password) {
      setIsInvalid(true);
      toast?.current?.show({ severity: 'error', summary: 'Erro', detail: 'Preencha todos os campos!', life: 3000 });
      return;
    }
    try {
      await signIn({
        email,
        password,
      });
    } catch (error) {
      if (error instanceof AuthError && error.code === 'TENANT_INACTIVE_VECX') {
        setVecxBlockVisible(true);
        return;
      }
      const message = error instanceof Error ? error.message : 'Falha ao realizar login';
      toast?.current?.show({ severity: 'error', summary: 'Erro no login', detail: message, life: 4000 });
    }
  }

  return (
    <div className={containerClassName}>
      <Dialog
        header="Utilização não autorizada"
        visible={vecxBlockVisible}
        onHide={() => setVecxBlockVisible(false)}
        modal
        dismissableMask
        style={{ width: 'min(32rem, 95vw)' }}
        headerStyle={{
          background: '#b71c1c',
          color: '#fff',
          borderRadius: '6px 6px 0 0',
        }}
        contentStyle={{ padding: '1.25rem' }}
      >
        <p className="m-0 line-height-3 text-900">
          Esta utilização não está autorizada. Todos os dados estão salvos. Entre em contato diretamente com a VECX para
          regularizar a situação.
        </p>
      </Dialog>
      <div className="flex flex-column align-items-center justify-content-center">
        <div style={{ borderRadius: '56px', padding: '0.3rem', background: 'linear-gradient(180deg, var(--primary-color) 10%, rgba(33, 150, 243, 0) 30%)' }}>
          <div className="w-full surface-card py-8 px-5 sm:px-8" style={{ borderRadius: '53px' }}>
            <div className="text-center mb-5">
              <img src="/vecx_logo.svg" alt="Vecontab logo" className="mb-5 w-16rem flex-shrink-0" />

              <div className="text-900 text-3xl font-medium mb-3"></div>
              <span className="text-600 font-medium">Faça login para continuar</span>
            </div>

            <div>


              <label htmlFor="login-email" className="block text-900 text-xl font-medium mb-2">
                Email
              </label>
              <InputText
                id="login-email"
                name="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                type="email"
                autoComplete="username"
                placeholder="Email"
                className={`w-full md:w-30rem mb-5 ${isInvalid ? 'p-invalid' : ''}`}
                style={{ padding: '1rem' }}
              />

              <label htmlFor="login-password" className="block text-900 font-medium text-xl mb-2">
                Senha
              </label>
              <Password
                inputId="login-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                feedback={false}
                placeholder="Senha"
                toggleMask
                className={`w-full md:w-30rem mb-5 ${isInvalid ? 'p-invalid' : ''}`}
                inputClassName="w-full p-3 md:w-30rem"
              ></Password>

              <div className="flex align-items-center justify-content-between mb-5 gap-5">

                <span className="font-medium ml-2 text-center" style={{ color: 'var(--primary-color)' }}>
                  Ainda não tem conta ? Cadastre-se{' '}
                  <Link href='/auth/register' className="font-bold no-underline cursor-pointer">
                    aqui
                  </Link>
                </span>

              </div>

              <div className="card flex justify-content-center">
                <Toast ref={toast} />
                <Button label="Acessar" className="w-full p-3 text-xl" onClick={handleLogin}></Button>
              </div>

            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

LoginPage.getLayout = function getLayout(page) {
  return (
    <React.Fragment>
      {page}
      <AppConfig simple />
    </React.Fragment>
  );
};
export default LoginPage;