package validation

import "unicode"

func onlyDigits(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsDigit(r) {
			out = append(out, r)
		}
	}
	return string(out)
}

func allDigitsEqual(s string) bool {
	if len(s) == 0 {
		return false
	}
	first := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != first {
			return false
		}
	}
	return true
}

func IsValidCPF(raw string) bool {
	cpf := onlyDigits(raw)
	if len(cpf) != 11 || allDigitsEqual(cpf) {
		return false
	}

	calc := func(base string, factor int) int {
		total := 0
		for i := 0; i < len(base); i++ {
			total += int(base[i]-'0') * factor
			factor--
		}
		mod := total % 11
		if mod < 2 {
			return 0
		}
		return 11 - mod
	}

	d1 := calc(cpf[:9], 10)
	d2 := calc(cpf[:10], 11)
	return cpf == cpf[:9]+string(rune('0'+d1))+string(rune('0'+d2))
}

func IsValidCNPJ(raw string) bool {
	cnpj := onlyDigits(raw)
	if len(cnpj) != 14 || allDigitsEqual(cnpj) {
		return false
	}

	calc := func(base string, weights []int) int {
		total := 0
		for i := 0; i < len(base); i++ {
			total += int(base[i]-'0') * weights[i]
		}
		mod := total % 11
		if mod < 2 {
			return 0
		}
		return 11 - mod
	}

	w1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	w2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	d1 := calc(cnpj[:12], w1)
	base2 := cnpj[:12] + string(rune('0'+d1))
	d2 := calc(base2, w2)
	return cnpj == cnpj[:12]+string(rune('0'+d1))+string(rune('0'+d2))
}

