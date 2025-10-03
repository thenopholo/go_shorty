.PHONY: help dev-start dev-stop dev-restart dev-logs prod-start prod-stop prod-restart clean rebuild test

# Variáveis
DOCKER_COMPOSE_DEV = docker-compose -f docker-compose.dev.yml
DOCKER_COMPOSE_PROD = docker-compose

# Comando padrão - mostra ajuda
help:
	@echo "Shorty URL - Comandos Disponíveis"
	@echo ""
	@echo "Desenvolvimento (Hot Reload):"
	@echo "  make dev-start       - Inicia ambiente de desenvolvimento com hot reload"
	@echo "  make dev-stop        - Para ambiente de desenvolvimento"
	@echo "  make dev-restart     - Reinicia ambiente de desenvolvimento"
	@echo "  make dev-logs        - Exibe logs do ambiente de desenvolvimento"
	@echo ""
	@echo "Produção:"
	@echo "  make prod-start      - Inicia ambiente de produção"
	@echo "  make prod-stop       - Para ambiente de produção"
	@echo "  make prod-restart    - Reinicia ambiente de produção"
	@echo ""
	@echo "Limpeza:"
	@echo "  make clean           - Remove containers, volumes e imagens"
	@echo "  make rebuild         - Rebuild completo (dev)"
	@echo ""
	@echo "Testes:"
	@echo "  make test            - Executa testes"
	@echo ""

# Desenvolvimento
dev-start:
	@echo "Iniciando ambiente de desenvolvimento com hot reload..."
	$(DOCKER_COMPOSE_DEV) up --build

dev-start-d:
	@echo "Iniciando ambiente de desenvolvimento em background..."
	$(DOCKER_COMPOSE_DEV) up -d --build

dev-stop:
	@echo "Parando ambiente de desenvolvimento..."
	$(DOCKER_COMPOSE_DEV) down

dev-restart:
	@echo "Reiniciando ambiente de desenvolvimento..."
	$(DOCKER_COMPOSE_DEV) restart

dev-logs:
	@echo "Exibindo logs do ambiente de desenvolvimento..."
	$(DOCKER_COMPOSE_DEV) logs -f

dev-logs-app:
	@echo "Exibindo logs da aplicação..."
	$(DOCKER_COMPOSE_DEV) logs -f app

# Produção
prod-start:
	@echo "Iniciando ambiente de produção..."
	$(DOCKER_COMPOSE_PROD) up --build

prod-start-d:
	@echo "Iniciando ambiente de produção em background..."
	$(DOCKER_COMPOSE_PROD) up -d --build

prod-stop:
	@echo "Parando ambiente de produção..."
	$(DOCKER_COMPOSE_PROD) down

prod-restart:
	@echo "Reiniciando ambiente de produção..."
	$(DOCKER_COMPOSE_PROD) restart

prod-logs:
	@echo "Exibindo logs do ambiente de produção..."
	$(DOCKER_COMPOSE_PROD) logs -f

# Limpeza
clean:
	@echo "Limpando containers, volumes e imagens..."
	$(DOCKER_COMPOSE_DEV) down -v --rmi local
	$(DOCKER_COMPOSE_PROD) down -v --rmi local
	@echo "Limpeza concluída!"

clean-all:
	@echo "Limpeza completa (incluindo volumes do banco)..."
	$(DOCKER_COMPOSE_DEV) down -v
	$(DOCKER_COMPOSE_PROD) down -v
	docker volume prune -f
	@echo "Limpeza completa concluída!"

rebuild:
	@echo "Rebuild completo do ambiente de desenvolvimento..."
	$(DOCKER_COMPOSE_DEV) down -v
	$(DOCKER_COMPOSE_DEV) up --build --force-recreate

# Testes
test:
	@echo "Executando testes..."
	go test -v ./...

test-coverage:
	@echo "Executando testes com cobertura..."
	go test -v -cover ./...

# Database
db-shell:
	@echo "Acessando shell do PostgreSQL..."
	$(DOCKER_COMPOSE_DEV) exec postgres_db psql -U admin -d shorty_url

# Shell
shell:
	@echo "Acessando shell do container da aplicação..."
	$(DOCKER_COMPOSE_DEV) exec app sh

# Formatação e Lint
fmt:
	@echo "Formatando código..."
	go fmt ./...

lint:
	@echo "Executando linter..."
	golangci-lint run ./...

# Instalação de dependências
deps:
	@echo "Instalando dependências..."
	go mod download
	go mod tidy

# Status
status:
	@echo "Status dos containers:"
	@docker ps -a | grep shorty_url || echo "Nenhum container encontrado"
