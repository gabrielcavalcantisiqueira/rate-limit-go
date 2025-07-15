# Rate Limiter em Go / Rate Limiter in Go

# Português

Projeto de um **rate limiter** implementado em Go para controlar o número máximo de requisições por segundo, baseado em endereço IP ou token de acesso.

### Descrição

Este rate limiter pode ser configurado para limitar requisições por IP ou por token, com diferentes limites e tempos de bloqueio configuráveis. A persistência das informações é feita via Redis, garantindo eficiência e escalabilidade.

O sistema funciona como middleware para um servidor web, retornando HTTP 429 quando o limite for ultrapassado.

### Configuração

Crie um arquivo `.env` na raiz do projeto com as variáveis abaixo (exemplo):

```
REDIS_ADDR=redis:6379
DEFAULT_LIMIT=5
DEFAULT_WINDOW=1s
DEFAULT_LOCK=5m
TOKENS_LIMITS=abc123:10:1s:2m,xyz789:100:1s:1m
``` 

 - REDIS_ADDR: endereço do Redis (ex: redis:6379 para Docker Compose)
 - DEFAULT_LIMIT: número padrão máximo de requisições por segundo (por IP)
 - DEFAULT_WINDOW: janela de tempo para contagem (ex: 1s)
 - DEFAULT_LOCK: tempo que o IP ou token fica bloqueado após ultrapassar o limite
 - TOKENS_LIMITS: configurações específicas para tokens no formato token:limite:janela:bloqueio, separados por vírgula

### Rodando com Docker Compose
O projeto inclui um docker-compose.yml que sobe o Redis e a aplicação:

``` docker-compose up -d -build ```

Isso iniciará o Redis e o serviço Go que escuta na porta 8080.

### Testando o rate limiter
Você pode testar o rate limiter fazendo requisições para o serviço, usando o IP ou o header API_KEY para limitar conforme configurado.

Exemplos com curl:

``` curl -H http://localhost:8080/ping ```

``` curl -H "API_KEY: abc123" http://localhost:8080/ping ``` 

Se o limite for ultrapassado, o serviço retornará HTTP 429 com a mensagem:

``` 429 - Too Many Requests: you have reached the maximum number of requests or actions allowed within a certain time frame ```

# English

Project of a **rate limiter** implemented in Go to control the maximum number of requests per second, based on IP address or access token.

### Description

This rate limiter can be configured to limit requests by IP or by token, with different configurable limits and blocking times. The data persistence is done via Redis, ensuring efficiency and scalability.

The system works as middleware for a web server, returning HTTP 429 when the limit is exceeded.

### Configuration

Create a `.env` file in the root of the project with the following variables (example):

```
REDIS_ADDR=redis:6379
DEFAULT_LIMIT=5
DEFAULT_WINDOW=1s
DEFAULT_LOCK=5m
TOKENS_LIMITS=abc123:10:1s:2m,xyz789:100:1s:1m
``` 


- `REDIS_ADDR`: Redis address (e.g., `redis:6379` for Docker Compose)
- `DEFAULT_LIMIT`: default maximum number of requests per second (per IP)
- `DEFAULT_WINDOW`: time window for counting (e.g., `1s`)
- `DEFAULT_LOCK`: time the IP or token remains blocked after exceeding the limit
- `TOKENS_LIMITS`: specific token configurations in the format `token:limit:window:lock`, separated by commas

### Running with Docker Compose

The project includes a `docker-compose.yml` file that runs Redis and the application:

``` docker-compose up -d --build ```

This will start Redis and the Go service listening on port 8080.

### Testing the rate limiter
You can test the rate limiter by making requests to the service, using either the IP or the API_KEY header to limit as configured.

Example with curl:

``` curl -H http://localhost:8080/ping ```

``` curl -H "API_KEY: abc123" http://localhost:8080/ping ``` 

If the limit is exceeded, the service will return HTTP 429 with the message:

``` 429 - Too Many Requests: you have reached the maximum number of requests or actions allowed within a certain time frame ```