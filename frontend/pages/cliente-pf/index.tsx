import { EmpresasPage } from '../empresas';
import { useTenantIdQuery } from '../../components/hooks/useClientGuards';

export default function ClientePF() {
  const { data: tenantid = '' } = useTenantIdQuery();
  return <EmpresasPage dados={tenantid} tipoPessoa="PF" />;
}
