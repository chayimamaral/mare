#!/bin/bash
# para essa execução, antes foi configurado o seguinte:
#git config --global alias.save '!git add . && git commit -m "$1" && git push origin main && echo "🚀 Enviado com sucesso!"'

# Verifica se a mensagem de commit foi enviada como argumento
if [ -z "$1" ]; then
    echo "❌ Erro: Forneça uma mensagem para o commit."
    echo "Exemplo: ./git_save.sh 'Refatoração do domínio clientes'"
    exit 1
fi

# Executa os comandos do Git
echo "---------------------------------------"
echo "📦 Adicionando arquivos..."
git add .

echo "💾 Criando commit: \"$1\""
if git commit -m "$1"; then
    echo "🚀 Enviando para o GitHub (main)..."
    git push origin main
    echo "---------------------------------------"
    echo "✅ Sucesso! Código atualizado."
else
    echo "⚠️ Nada para commitar ou erro no processo."
fi