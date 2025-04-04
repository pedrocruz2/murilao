FROM node:16-alpine

WORKDIR /app

# Copiar package.json e instalar dependências
COPY package*.json ./
RUN npm install --unsafe-perm

# Copiar o restante dos arquivos
COPY . .

EXPOSE 3000

CMD ["npm", "start"]