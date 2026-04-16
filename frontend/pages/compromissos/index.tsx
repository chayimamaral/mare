import { useEffect } from 'react';
import { useRouter } from 'next/router';

/** Rota antiga; o cadastro unificado está em `/obrigacoes`. */
const CompromissosRedirect = () => {
  const router = useRouter();

  useEffect(() => {
    void router.replace('/obrigacoes');
  }, [router]);

  return null;
};

export default CompromissosRedirect;

  // sem processamento adicional
