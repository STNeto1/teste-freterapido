# 📄 API — Fluxo de Cotações e Métricas

Este documento descreve o fluxo de requisições da API, suas validações, tentativas de integração com serviços externos e respostas possíveis.

---

## 🚚 POST `/quote`

### ✅ **Fluxo de sucesso**

```
Cliente → [POST /quote] → API → Validação → OK
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
Cliente → [POST /quote] → API → Validação → Inválido → Responde 400
```

* Se os dados forem inválidos, a API retorna **HTTP 400 Bad Request** imediatamente.

---

### ❌ **Fluxo de falha na cotação**

```
Cliente → [POST /quote] → API → Validação → OK
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
    price Float64,
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
       avg(price) AS avg_price
FROM last_quotes
GROUP BY carrier
ORDER BY carrier;

-- LIMIT <n> pode ser usado para limitar
-- é adicionado dinamicamente pelo servidor.
-- Go não tem uma forma melhor de fazer isso
-- sem mandar um `N ou MAX_INT`.
WITH last_quotes AS (SELECT *
                     FROM quotes
                     ORDER BY timestamp DESC
                     LIMIT -1)
SELECT min(price) AS lowest_price,
       max(price) AS highest_price
FROM last_quotes;
```
