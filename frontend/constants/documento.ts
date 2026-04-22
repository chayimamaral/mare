export const onlyDigits = (value: string): string => String(value ?? '').replace(/\D/g, '');

const allDigitsEqual = (digits: string): boolean => /^(\d)\1+$/.test(digits);

export function isValidCPF(value: string): boolean {
  const cpf = onlyDigits(value);
  if (cpf.length !== 11 || allDigitsEqual(cpf)) return false;

  const calc = (base: string, factor: number): number => {
    let total = 0;
    for (const ch of base) {
      total += Number(ch) * factor--;
    }
    const mod = total % 11;
    return mod < 2 ? 0 : 11 - mod;
  };

  const d1 = calc(cpf.slice(0, 9), 10);
  const d2 = calc(cpf.slice(0, 10), 11);
  return cpf === `${cpf.slice(0, 9)}${d1}${d2}`;
}

export function isValidCNPJ(value: string): boolean {
  const cnpj = onlyDigits(value);
  if (cnpj.length !== 14 || allDigitsEqual(cnpj)) return false;

  const calc = (base: string, weights: number[]): number => {
    const total = base.split('').reduce((sum, ch, idx) => sum + Number(ch) * weights[idx], 0);
    const mod = total % 11;
    return mod < 2 ? 0 : 11 - mod;
  };

  const w1 = [5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2];
  const w2 = [6, ...w1];
  const d1 = calc(cnpj.slice(0, 12), w1);
  const d2 = calc(cnpj.slice(0, 12) + d1, w2);
  return cnpj === `${cnpj.slice(0, 12)}${d1}${d2}`;
}

