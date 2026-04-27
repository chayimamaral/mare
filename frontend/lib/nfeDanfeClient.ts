import api from '../components/api/apiClient';

export type NFEDanfeIdentificacao = {
    chave: string;
    modelo: string;
    serie: string;
    numero: string;
    emissao_em: string;
    saida_entrada_em: string;
    protocolo: string;
    codigo_status: string;
    data_autorizacao: string;
    evento_descricao: string;
    ambiente: string;
    situacao: string;
    natureza_operacao: string;
    tipo_operacao: string;
    destino_operacao: string;
    consumidor_final: string;
    presenca_comprador: string;
    processo_emissao: string;
    versao_processo: string;
    tipo_emissao: string;
    finalidade: string;
    forma_pagamento: string;
    digest_value: string;
    data_inclusao_bd: string;
};

export type NFEDanfePessoa = {
    nome: string;
    nome_fantasia?: string;
    cnpj_cpf: string;
    ie: string;
    ie_substituto?: string;
    im?: string;
    cnae?: string;
    crt?: string;
    crt_descricao?: string;
    indicador_ie_dest?: string;
    indicador_ie_descricao?: string;
    logradouro: string;
    numero: string;
    bairro: string;
    municipio: string;
    municipio_codigo?: string;
    municipio_cod_nome?: string;
    uf: string;
    cep: string;
    pais_codigo?: string;
    pais_nome?: string;
    pais_cod_nome?: string;
    telefone?: string;
    email?: string;
    isuf?: string;
    cod_mun_fato_gerador_icms?: string;
    endereco_completo?: string;
};

export type NFEDanfeItemICMS = {
    origem: string;
    tributacao: string;
};

export type NFEDanfeItemIPI = {
    cl_enq: string;
    c_enq: string;
    c_selo: string;
    cnpj_prod: string;
    q_selo: string;
    cst: string;
    q_unid: string;
    v_unid: string;
    v_ipi: string;
    v_bc: string;
    p_ipi: string;
};

export type NFEDanfeItem = {
    codigo: string;
    descricao: string;
    ncm: string;
    extipi: string;
    cfop: string;
    unidade: string;
    quantidade: string;
    valor_unitario: string;
    valor_total: string;
    valor_desconto: string;
    valor_frete: string;
    valor_seguro: string;
    valor_outros: string;
    indicador_total_nf: string;
    indicador_total_desc?: string;
    cean: string;
    cean_trib: string;
    u_trib: string;
    q_trib: string;
    v_un_trib: string;
    valor_total_tributos: string;
    x_ped?: string;
    n_item_ped?: string;
    n_fci?: string;
    base_icms: string;
    valor_icms: string;
    valor_ipi: string;
    aliquota_icms: string;
    aliquota_ipi: string;
    icms: NFEDanfeItemICMS;
    ipi: NFEDanfeItemIPI;
    pis_cst?: string;
    cofins_cst?: string;
};

export type NFEDanfeView = {
    identificacao: NFEDanfeIdentificacao;
    emitente: NFEDanfePessoa;
    destinatario: NFEDanfePessoa;
    itens: NFEDanfeItem[];
    totais: {
        base_icms: string;
        valor_icms: string;
        valor_icms_desonerado: string;
        base_icms_st: string;
        valor_st: string;
        valor_ii: string;
        valor_ipi: string;
        valor_pis: string;
        valor_cofins: string;
        valor_produtos: string;
        valor_frete: string;
        valor_seguro: string;
        valor_desconto: string;
        valor_outros: string;
        valor_total_tributos: string;
        valor_nota: string;
    };
    transporte: {
        modalidade: string;
        transportador: string;
        cnpj_cpf: string;
        ie: string;
        endereco: string;
        municipio: string;
        placa: string;
        uf: string;
        rntc: string;
        quantidade_volumes: string;
        volumes: Array<{
            quantidade: string;
            especie: string;
            marca: string;
            numero: string;
            peso_liquido: string;
            peso_bruto: string;
        }>;
    };
    cobranca: {
        numero_fatura: string;
        valor_original: string;
        valor_desconto: string;
        valor_liquido: string;
        duplicatas: Array<{
            numero: string;
            vencimento: string;
            valor: string;
        }>;
        pagamentos: Array<{
            forma: string;
            valor: string;
            cnpj_credenciadora: string;
            bandeira: string;
            autorizacao: string;
        }>;
    };
    adicionais: {
        tp_imp?: string;
        informacoes_complementares: string;
        informacoes_fisco: string;
    };
};

export function parseDanfeErrorMessage(e: unknown): string {
    const err = e as { response?: { data?: unknown }; message?: string };
    let msg =
        (err?.response?.data as { error?: string; message?: string })?.error
        ?? (err?.response?.data as { message?: string })?.message
        ?? err?.message
        ?? 'Falha ao gerar DANFE';
    const raw = err?.response?.data;
    if (typeof raw === 'string' && raw.trim().startsWith('{')) {
        try {
            const o = JSON.parse(raw) as { error?: string };
            if (o?.error) {
                msg = o.error;
            }
        } catch {
            /* ignore */
        }
    }
    return String(msg);
}

export async function fetchDanfeHtmlFromRetorno(retorno: string): Promise<string> {
    const { data } = await api.post<string>(
        '/api/serpro/nfe/documento/danfe-html',
        { retorno },
        { responseType: 'text' as any },
    );
    const html = String(data ?? '').trim();
    if (!html) {
        throw new Error('DANFE_VAZIO');
    }
    return html;
}

export async function fetchDanfeJsonByChave(chave: string): Promise<NFEDanfeView> {
    const { data } = await api.get<NFEDanfeView>('/api/serpro/nfe/documento/danfe-json', {
        params: { chave },
        timeout: 20000,
    });
    return data;
}
