# 📄 API — Fluxo de Cotações e Métricas

Este documento descreve o fluxo de requisições da API, suas validações, tentativas de integração com serviços externos e respostas possíveis.

---

## Como executar

### Local
- É preciso ter um clickhouse local configurado de forma específica:
  - Database: "freterapido"
  -	Username: "default"
  - Password: "admin"

```shell
$ go run cmd/webserver/main.go start --registered-number="25438296000158" --token="1d52a9b6b78cf07b08586152459a5c90" --platform-code="5AKVkHqCn" --dispatcher-zip-code=29161376 --clickhouse-addr="192.168.1.7:9000"
```

### Container
- É preciso ter docker compose local. Dockerfile contém comando equivalente ao rodar local, apenas configurado para network do docker.
```shell
$ docker compose up
```

### Testes
- Testes podem ser rodados localmente
```shell
$ go test -v ./...
```
- Também podem ser encontrados no repositório na aba de actions ou no último push.

---

## Como testar

### POST `/quotes`
```shell
$ curl --request POST \
  --url http://localhost:8080/quotes \
  --header 'content-type: application/json' \
  --data '{
  "recipient": {
    "address": {
      "zipcode": "01311000"
    }
  },
  "volumes": [
    {
      "category": 7,
      "amount": 1,
      "unitary_weight": 5,
      "price": 349,
      "sku": "abc-teste-123",
      "height": 0.2,
      "width": 0.2,
      "length": 0.2
    },
    {
      "category": 7,
      "amount": 2,
      "unitary_weight": 4,
      "price": 556,
      "sku": "abc-teste-527",
      "height": 0.4,
      "width": 0.6,
      "length": 0.15
    }
  ]
}'
```

### GET `/metrics`
```shell
$ curl --request GET \
  --url http://localhost:8080/metrics

$ curl --request GET \
  --url http://localhost:8080/metrics?last_quotes=1000

```

---

## 🚚 POST `/quotes`

### ✅ **Fluxo de sucesso**

```
Cliente → [POST /quotes] → API → Validação → OK
                                      ↓
                               Tenta cotação → 2xx
                                      ↓
                  Responde 2xx ao cliente
                  → Em paralelo (goroutine) → Salva dados no Clickhouse
```

* A requisição é validada.
* A cotação é solicitada na API Frete Rápido.
* Se a cotação for bem-sucedida (2xx), os dados são salvos no Clickhouse de forma **assíncrona** (não bloqueia a resposta).
* O cliente recebe resposta **2xx**.

---

### ❌ **Fluxo de validação inválida**

```
Cliente → [POST /quotes] → API → Validação → Inválido → Responde 400
```

* Se os dados forem inválidos, a API retorna **HTTP 400 Bad Request** imediatamente.

---

### ❌ **Fluxo de falha na cotação**

```
Cliente → [POST /quotes] → API → Validação → OK
                                      ↓
                      Tenta cotação (1ª tentativa) → Falha
                      Tenta cotação (2ª tentativa) → Falha
                      ...
                      Tenta cotação (Nª tentativa) → Falha
                                      ↓
                               Responde 424 / 4xx
```

* Caso a validação seja OK, mas todas as tentativas de cotação falhem, a API responde com **HTTP 424 Failed Dependency** ou outro **4xx** adequado.

---

## 📊 GET `/metrics`

### ✅ **Fluxo de sucesso**

```
Cliente → [GET /metrics] → API → Validação → OK
                                      ↓
                              Consulta métricas → 2xx
                                      ↓
                               Responde 2xx
```

* A requisição é validada.
* As métricas são consultadas no Clickhouse.
* Se bem-sucedido, responde **HTTP 2xx** com os dados.

---

### ❌ **Fluxo de falha na consulta**

```
Cliente → [GET /metrics] → API → Validação → OK
                                      ↓
                    Tenta métricas (1ª tentativa) → Falha
                    Tenta métricas (2ª tentativa) → Falha
                    ...
                    Tenta métricas (Nª tentativa) → Falha
                                      ↓
                               Responde 424 / 4xx
```

* Se todas as tentativas de consulta falharem, responde com **HTTP 424** ou outro **4xx**.

---

### ❌ **Parâmetro inválido**

```
Cliente → [GET /metrics?last_quotes=a] → API → Validação → Inválido → Responde 4xx
```

* Se o parâmetro `last_quotes` for inválido (ex: não numérico), responde **HTTP 4xx**.

---

### ✅ **Parâmetro válido**

```
Cliente → [GET /metrics?last_quotes=100] → API → Validação → OK
                                              ↓
                                      Consulta métricas → 2xx
                                              ↓
                                       Responde 2xx
```

* Consulta limitada pelas últimas cotações solicitadas.
* Retorna **2xx** com o resultado.

---

## 📌 **Observações**

* 📌 **Validação:** Sempre ocorre antes de interações externas.
* 🔄 **Tentativas:** Operações críticas podem ter múltiplas tentativas.
* ⚙️ **Assíncrono:** O salvamento no Clickhouse é feito em segundo plano após resposta de sucesso na cotação.

-----

### Regras adicionais
- recipient.address.zipcode -> regex de CEP
- volumes[*].category -> https://dev.freterapido.com.br/common/tipos_de_volumes/
- volumes[*].amount -> >0
- volumes[*].height -> >0.0
- volumes[*].width -> >0.0
- volumes[*].length -> >0.0
- volumes[*].unitary_price -> >0.0
- volumes[*].unitary_weight -> >0.0


-----

## Schema do Clickhouse

```sql
-- particionar por mês do timestamp PODE ser uma boa
-- mas para o momento não é preciso ^^
CREATE TABLE quotes
(
    id UUID DEFAULT generateUUIDv7(),
    name String,
    service String,
    deadline UInt8,
    price Decimal(10, 2),
    timestamp DateTime DEFAULT now()
)
ENGINE = MergeTree()
ORDER BY (service, timestamp);
```

### Inserção de registros
```sql
-- timestamp pode ser omitido já que
-- o registro será inserido no momento
-- da requisição.
insert into quotes (name, service, deadline, price)
values (?, ?, ?, ?);
```

### Pesquisa de registros
```sql
-- LIMIT <n> pode ser usado para limitar
-- é adicionado dinamicamente pelo servidor.
-- Go não tem uma forma melhor de fazer isso
-- sem mandar um `N ou MAX_INT`.
WITH last_quotes AS (SELECT *
                     FROM quotes
                     ORDER BY timestamp DESC
                     LIMIT -1)
SELECT name       AS carrier,
       count()    AS total_quotes,
       sum(price) AS total_price,
       avg(price) AS avg_price,
       min(price) as min_price,
       max(price) as max_price
FROM last_quotes
GROUP BY carrier
ORDER BY carrier;
```
