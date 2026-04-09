import setupAPIClient from '../../components/api/api';

export type CertificadoClienteResumo = {
  tipo_certificado?: string;
  nome_certificado?: string;
  emitido_para?: string;
  emitido_por?: string;
  cnpj?: string;
  validade_de?: string;
  validade_ate?: string;
};

export default function CertificadoClienteService() {
  const getByEmpresa = async (empresaId: string) => {
    const apiClient = setupAPIClient(undefined);
    const response = await apiClient.get('/api/certificado-cliente', {
      params: { empresa_id: empresaId },
    });
    return { data: response.data as { certificado?: CertificadoClienteResumo } };
  };

  const upload = async (params: { empresaId: string; arquivo: File; senha_certificado: string; titular_nome?: string; cnpj?: string }) => {
    const apiClient = setupAPIClient(undefined);
    const body = new FormData();
    body.append('empresa_id', params.empresaId);
    body.append('arquivo', params.arquivo);
    body.append('senha_certificado', params.senha_certificado);
    if (params.titular_nome?.trim()) {
      body.append('titular_nome', params.titular_nome.trim());
    }
    if (params.cnpj?.trim()) {
      body.append('cnpj', params.cnpj.trim());
    }
    const response = await apiClient.post('/api/certificado-cliente/upload', body, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return { data: response.data };
  };

  return { getByEmpresa, upload };
}
