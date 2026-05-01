LOCAL_DIR="/home/camaral/backups/vecx"

echo "Limpando a pasta local dos backups"

rm -rf "$LOCAL_DIR"/*

# No início do deploy.sh

./bkp_db.sh || { echo "Backup do banco de dados falhou, deploy cancelado"; exit 1; }

# Se o backup falhar e você quiser parar o deploy:
./backup.sh || { echo "Backup falhou, deploy cancelado"; exit 1; }
./bkp_drive.sh