# Go Shorty - Um Encurtador de URLs com Go

Olá! Este é o **Shorty URL**, um projeto que construí para explorar e aprofundar meus conhecimentos em Go, criando uma aplicação web robusta, eficiente e moderna do zero. Mais do que apenas um encurtador de URLs, este é um estudo de caso sobre boas práticas de desenvolvimento, arquitetura de software e automação com Docker.

## Visão Geral

O Go Shorty é um serviço que transforma URLs longas e complexas em links curtos e fáceis de compartilhar. Ele foi construído com uma API RESTful em Go no backend, um banco de dados PostgreSQL para persistência e uma interface de usuário simples em HTML e Tailwind CSS para demonstrar a funcionalidade.

## Tecnologias Utilizadas

- **Backend**: **Go 1.25+** com o roteador **Chi** pela sua leveza e compatibilidade com a `net/http` padrão.
- **Frontend**: **HTML5**, **Tailwind CSS** e **JavaScript** para uma interface limpa e interativa.
- **Banco de Dados**: **PostgreSQL 15**, escolhido por sua robustez e confiabilidade.
- **Containerização**: **Docker** e **Docker Compose** para criar ambientes de desenvolvimento e produção consistentes e isolados.
- **Hot-Reload**: **Air** para um fluxo de desenvolvimento ágil com recarregamento automático.

---

## Arquitetura e Decisões de Design

Esta seção é um mergulho técnico nas escolhas que moldaram o projeto.

### 1. O Algoritmo de Geração de Código

A funcionalidade central é a geração de um código único e curto. A abordagem precisava ser segura, eficiente e minimizar colisões.

**generateRandomCode.go (`internal/service/generateRandomCode.go`):**

```go
package service

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomCode(length int) (string, error) {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)

	for i := range result {
		// Gera um número aleatório criptograficamente seguro
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		// Mapeia o número para um caractere do charset
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
```

**Por que essa abordagem?**

- **Segurança com `crypto/rand`**: Em vez de usar `math/rand` (que é pseudoaleatório e inadequado para fins de segurança), utilizei `crypto/rand`. Ele gera números aleatórios a partir de fontes de entropia do sistema operacional, tornando os códigos gerados imprevisíveis.
- **Base de Caracteres (`charset`)**: Com 62 caracteres alfanuméricos, um código de 6 caracteres nos dá **56.800.235.584** (62^6) combinações possíveis, o que é mais do que suficiente para uma aplicação de pequeno a médio porte.
- **Tratamento de Colisões**: Embora raras, colisões (gerar um código que já existe) podem acontecer. A lógica no `handler` (`internal/handler/handler.go`) aborda isso de forma pragmática:

  ```go
  // ... (dentro do handler Shorten)
  var code string
  var err error
  maxAttempts := 5 // Tenta até 5 vezes para evitar loops infinitos
  for i := 0; i < maxAttempts; i++ {
      code, err = service.GenerateRandomCode(6)
      // ... (tratamento de erro)

      _, err = h.store.GetURL(code) // Verifica se o código já existe
      if err != nil { // Se err != nil (esperando sql.ErrNoRows), o código não foi encontrado e está livre
          break
      }
  }
  ```

### 2. Backend: Estrutura e Organização

O projeto segue a estrutura padrão de projetos Go, separando claramente as responsabilidades.

- `cmd/main.go`: O ponto de entrada da aplicação. Responsável por inicializar o banco de dados, o servidor e injetar as dependências.
- `internal/`: O coração da lógica de negócios, inacessível para outros projetos.
  - `handler/`: Camada de apresentação. Lida com as requisições HTTP, decodifica JSON e envia respostas.
  - `store/`: Camada de acesso a dados. Abstrai toda a comunicação com o PostgreSQL.
  - `service/`: Contém a lógica de negócio pura, como a geração de códigos.
  - `server/`: Configuração e inicialização do servidor HTTP e das rotas.

### 3. Frontend: Uma UI para a API

O frontend em `public/index.html` é intencionalmente simples. Seu objetivo é ser uma interface de demonstração para a API.

- **Tecnologia**: HTML, Tailwind CSS e JavaScript puro.
- **Funcionalidade**: Um formulário que envia a URL para a API via `fetch` e exibe o resultado.

  ```javascript
  // Snippet de public/index.html
  document
    .getElementById("shortenForm")
    .addEventListener("submit", async (e) => {
      e.preventDefault();
      // ... (lógica de UI para loading)

      try {
        const response = await fetch("/api/shorten", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ url: urlInput.value }),
        });

        const result = await response.json();

        if (!response.ok) {
          throw new Error(result.error || "Ocorreu um erro");
        }

        // ... (exibe a URL encurtada na UI)
      } catch (error) {
        // ... (exibe a mensagem de erro na UI)
      }
    });
  ```

### 4. Ambiente de Desenvolvimento com Hot-Reload

Para agilizar o desenvolvimento, configurei um ambiente Docker com hot-reload.

- **`Dockerfile.dev`**: Uma imagem de desenvolvimento que instala o **Air**, um live-reloader para Go.
- **`docker-compose.dev.yml`**: Orquestra os serviços. A chave é o `volume` que mapeia o código-fonte local para dentro do contêiner:
  ```yaml
  services:
    app:
      # ...
      volumes:
        - ./:/app # Mapeia o diretório atual para /app no contêiner
  ```
- **`.air.toml`**: Arquivo de configuração que instrui o Air a monitorar arquivos `.go` e reiniciar o servidor a cada alteração.

Essa configuração permite que eu edite o código na minha máquina e veja as alterações refletidas instantaneamente na aplicação rodando no Docker, sem intervenção manual.

---

## Como Executar

### Ambiente de Desenvolvimento (com Hot-Reload)

Este é o modo recomendado para desenvolver.

1.  **Clone o repositório**
2.  **Inicie os containers**:
    ```bash
    docker-compose -f docker-compose.dev.yml up --build
    ```
3.  **Acesse a aplicação**:
    - Frontend: `http://localhost:4200`
    - API: `http://localhost:4200/api`

Qualquer alteração nos arquivos `.go` reiniciará o servidor automaticamente.

### Ambiente de Produção

Este modo cria uma imagem Go otimizada e menor.

1.  **Inicie os containers**:
    ```bash
    docker-compose up --build
    ```
2.  **Acesse a aplicação**:
    - Frontend e API: `http://localhost:8080`

## API Endpoints

### Encurtar URL

- **`POST /api/shorten`**
- **Body**: `{ "url": "https://sua-url-longa.com" }`
- **Resposta de Sucesso (201 Created)**: `{ "data": { "code": "aBcDeF" } }`

### Redirecionar para URL Original

- **`GET /{code}`**
- Redireciona com um status `308 Permanent Redirect` para a URL original.

## Testando a API com `curl`

```bash
# 1. Encurtar uma URL
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com/thenopholo"}'

# Resposta esperada: {"data":{"code":"someCode"}}

# 2. Acessar a URL encurtada (use o código da resposta anterior)
curl -L http://localhost:8080/someCode
```

## Parar a Aplicação

```bash
# Para o ambiente de desenvolvimento
docker-compose -f docker-compose.dev.yml down

# Para o ambiente de produção
docker-compose down

# Para remover os dados do banco de dados (cuidado!)
docker-compose down -v
```

## Autor

Desenvolvido com entusiasmo como um projeto de estudo em Go. Sinta-se à vontade para explorar, fazer perguntas e se inspirar!
