# Sistema Web Modularizado

## Estrutura
- **Frontend**: React
- **Backend**: Go/Gin
- **Database**: PostgreSQL
- **Web Server**: Nginx

## Como Executar
1. Crie a pasta uploads: `mkdir -p uploads && chmod 777 uploads`
2. Ajuste permissões: `chmod 644 postgres/init.sql`
3. Inicie os serviços: `docker-compose up -d`
4. Acesse: http://localhost

## Correção de Problemas
- **Permissões Frontend**: `chmod -R 777 frontend`
- **Permissões Backend**: `chmod -R 777 backend`
- **Reiniciar Contêineres**: `docker-compose restart`

## Funcionalidades
- Cadastro/login de usuários
- Gerenciamento de produtos
- Upload de imagens