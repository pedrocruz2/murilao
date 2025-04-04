#!/bin/bash

# Cores para saída
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Reiniciando aplicação do zero ===${NC}"

# Parar todos os containers e remover volumes
echo -e "${BLUE}Parando containers existentes...${NC}"
docker-compose down -v

# Remover imagens para garantir rebuild completo
echo -e "${BLUE}Removendo imagens existentes...${NC}"
docker rmi $(docker images -q murilao_frontend murilao_backend murilao_nginx) 2>/dev/null || true

# Corrigir permissões
echo -e "${BLUE}Corrigindo permissões dos diretórios...${NC}"
mkdir -p uploads
chmod -R 777 uploads
chmod -R 777 frontend
chmod -R 777 backend

# Atualizar o Dockerfile do frontend para usar --unsafe-perm
echo -e "${BLUE}Atualizando Dockerfile do frontend...${NC}"
cat > frontend/Dockerfile << EOL
FROM node:16-alpine

WORKDIR /app

# Copiar tudo e definir permissões
COPY . /app/
RUN chmod -R 777 /app

# Instalar dependências com --unsafe-perm para ignorar verificações de permissão
RUN npm install --unsafe-perm

EXPOSE 3000

CMD ["npm", "start"]
EOL

# Build e iniciar a aplicação
echo -e "${BLUE}Construindo e iniciando a aplicação...${NC}"
docker-compose build
docker-compose up -d

# Monitorar os logs e verificar se a aplicação está funcionando
echo -e "${BLUE}Verificando status dos containers...${NC}"
sleep 5
docker ps

echo -e "${GREEN}=== Aplicação reiniciada! ===${NC}"
echo -e "${GREEN}Acesse http://localhost no navegador${NC}"
echo -e "${BLUE}Para acompanhar os logs, execute:${NC} docker-compose logs -f"
echo -e "${BLUE}Em caso de problemas com o frontend, verifique:${NC} docker logs app_frontend"
echo -e "${BLUE}Em caso de problemas com o backend, verifique:${NC} docker logs app_backend"