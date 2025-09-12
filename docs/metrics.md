# Метрики MSA Rocket Factory

## Архитектура сбора метрик

Система сбора метрик построена на основе OpenTelemetry и использует следующую архитектуру:

```
Application → OpenTelemetry Meter → OTLP Exporter → OTLP Collector → Prometheus → Grafana
```

### Компоненты:

1. **OpenTelemetry MeterProvider** - центральный объект для управления метриками в приложении
2. **OTLP Exporter** - отправляет метрики в коллектор по протоколу OTLP (gRPC)
3. **OTLP Collector** - принимает метрики и пересылает их в Prometheus
4. **Prometheus** - хранит временные ряды метрик
5. **Grafana** - визуализирует метрики в дашбордах

### Модель отправки:

- **Push Model** - приложение активно отправляет метрики в коллектор
- Периодическая отправка каждые 10 секунд
- Таймаут отправки: 5 секунд
- Адрес коллектора: `localhost:4317` (gRPC)

## Бизнес-метрики

### Сервис Order (Заказы)

#### Метрики заказов:

| Метрика | Тип | Описание | Единицы | Метки |
|---------|-----|----------|---------|-------|
| `orders_total` | Counter | Общее количество заказов | 1 | currency, status |
| `orders_created_total` | Counter | Количество созданных заказов | 1 | currency, status |
| `orders_paid_total` | Counter | Количество оплаченных заказов | 1 | currency, status |
| `orders_cancelled_total` | Counter | Количество отмененных заказов | 1 | currency, status |
| `orders_revenue_total` | Counter | Суммарная выручка от заказов | currency | currency, status |
| `order_value` | Histogram | Стоимость заказов | currency | currency, status |
| `order_duration_seconds` | Histogram | Время обработки заказов | s | operation |

#### События записи метрик:

- **Создание заказа**: увеличивается `orders_total`, `orders_created_total`, `orders_revenue_total`, записывается `order_value`
- **Оплата заказа**: увеличивается `orders_paid_total`
- **Обработка заказа**: записывается `order_duration_seconds`

### Сервис Assembly (Сборка ракет)

#### Метрики сборки:

| Метрика | Тип | Описание | Единицы | Метки |
|---------|-----|----------|---------|-------|
| `rockets_assembled_total` | Counter | Количество собранных ракет | 1 | rocket_type, status |
| `assembly_in_progress` | UpDownCounter | Количество ракет в процессе сборки | 1 | rocket_type, status |
| `assembly_errors_total` | Counter | Количество ошибок при сборке | 1 | rocket_type, error_type, status |
| `assembly_duration_seconds` | Histogram | Длительность сборки ракет | s | rocket_type, status |

#### События записи метрик:

- **Начало сборки**: увеличивается `assembly_in_progress`
- **Завершение сборки**: увеличивается `rockets_assembled_total`, записывается `assembly_duration_seconds`, уменьшается `assembly_in_progress`
- **Ошибка сборки**: увеличивается `assembly_errors_total`, уменьшается `assembly_in_progress`

## Технические метрики

### HTTP API метрики (OpenAPI)

Автоматически генерируемые метрики для всех HTTP endpoints:

#### Серверные метрики:

| Метрика | Тип | Описание | Единицы | Метки |
|---------|-----|----------|---------|-------|
| `http_server_request_duration` | Histogram | Длительность HTTP запросов | ms | method, route, status_code |
| `http_server_request_count` | Counter | Количество HTTP запросов | 1 | method, route, status_code |
| `http_server_errors_count` | Counter | Количество ошибок HTTP запросов | 1 | method, route, status_code |
| `otelogen_server_request_count` | Counter | Количество HTTP запросов (OpenAPI) | 1 | http_route |
| `otelogen_server_errors_count` | Counter | Количество ошибок HTTP запросов (OpenAPI) | 1 | http_route, http_response_status_code |

#### Клиентские метрики:

| Метрика | Тип | Описание | Единицы | Метки |
|---------|-----|----------|---------|-------|
| `http_client_request_duration` | Histogram | Длительность HTTP запросов клиента | ms | method, route, status_code |
| `http_client_request_count` | Counter | Количество HTTP запросов клиента | 1 | method, route, status_code |
| `http_client_errors_count` | Counter | Количество ошибок HTTP запросов клиента | 1 | method, route, status_code |

### gRPC метрики

Метрики для межсервисного взаимодействия через gRPC:

| Метрика | Тип | Описание | Единицы | Метки |
|---------|-----|----------|---------|-------|
| `grpc_client_requests_total` | Counter | Количество gRPC запросов клиента | 1 | service |
| `grpc_client_request_errors_total` | Counter | Количество ошибок gRPC запросов | 1 | service, grpc_code |

### Kafka метрики

Метрики для асинхронного взаимодействия через Kafka:

#### Producer метрики:

| Метрика | Тип | Описание | Единицы | Метки |
|---------|-----|----------|---------|-------|
| `kafka_producer_messages_total` | Counter | Количество отправленных сообщений | 1 | topic |
| `kafka_producer_messages_failed_total` | Counter | Количество неудачных отправок | 1 | topic |

#### Consumer метрики:

| Метрика | Тип | Описание | Единицы | Метки |
|---------|-----|----------|---------|-------|
| `kafka_consumer_messages_total` | Counter | Количество обработанных сообщений | 1 | topic, group_id |
| `kafka_consumer_messages_failed_total` | Counter | Количество неудачных обработок | 1 | topic, group_id |
| `kafka_consumer_offset_lag` | Gauge | Отставание consumer от producer | 1 | topic, group_id, partition |
| `kafka_consumer_rebalancing_total` | Counter | Количество ребалансировок | 1 | group_id |

### Метки метрик:

#### HTTP метки:
- `http.request.method` - HTTP метод (GET, POST, PUT, DELETE)
- `http.route` - маршрут API (например, `/api/v1/orders`)
- `http.response.status_code` - код ответа (200, 400, 500)
- `otelogen.operation_id` - ID операции (CreateOrder, PayOrder)

#### gRPC метки:
- `service` - имя gRPC сервиса
- `grpc_code` - код ответа gRPC (OK, UNAVAILABLE, DEADLINE_EXCEEDED)

#### Kafka метки:
- `topic` - название топика Kafka
- `group_id` - ID группы consumer'ов
- `partition` - номер партиции

## Дашборды Grafana

Проект использует два специализированных дашборда для разных аудиторий и задач:

### 1. Business Metrics Dashboard (`business_metrics.json`)

Дашборд для бизнес-команды и руководства, отображает ключевые бизнес-показатели и KPI.

#### Основные показатели (Stat Panels):

1. **Total Revenue** - общая выручка от всех заказов
   - Запрос: `sum(orders_revenue_total)`
   - Единицы: USD
   - Пороги: зеленый (норма), красный при >80

2. **Total Orders** - общее количество созданных заказов
   - Запрос: `sum(orders_total)`
   - Единицы: штуки

3. **Total Rockets Assembled** - общее количество собранных ракет
   - Запрос: `sum(rockets_assembled_total)`
   - Единицы: штуки

4. **Total Orders Paid** - общее количество оплаченных заказов
   - Запрос: `sum(orders_paid_total)`
   - Единицы: штуки

5. **Total Orders Cancelled** - общее количество отмененных заказов
   - Запрос: `sum(orders_cancelled_total)`
   - Единицы: штуки

6. **Order Conversion Rate** - процент конверсии заказов (оплаченные/созданные)
   - Запрос: `sum(rate(orders_paid_total[5m])) / sum(rate(orders_created_total[5m])) * 100`
   - Единицы: проценты
   - Пороги: красный <50%, желтый 50-80%, зеленый >80%

7. **Assembly Duration Distribution** - среднее время сборки ракеты
   - Запрос: `assembly_duration_seconds_sum / assembly_duration_seconds_count`
   - Единицы: секунды

8. **Assembly In Progress** - количество ракет в процессе сборки
   - Запрос: `sum(assembly_in_progress)`
   - Единицы: штуки
   - Пороги: зеленый (норма), желтый >5, красный >10

#### Скорости (Rates) - Мониторинг активности:

9. **Orders Created Rate** - скорость создания заказов (заказов/минуту)
   - Запрос: `sum(rate(orders_created_total[5m]))`
   - Пороги: красный (низкая активность), желтый 0.1-1, зеленый >1

10. **Orders Paid Rate** - скорость оплаты заказов (заказов/минуту)
    - Запрос: `sum(rate(orders_paid_total[5m]))`
    - Пороги: красный (низкая активность), желтый 0.05-0.5, зеленый >0.5

11. **Orders Cancelled Rate** - скорость отмены заказов (заказов/минуту)
    - Запрос: `sum(rate(orders_cancelled_total[5m]))`
    - Пороги: зеленый (норма), желтый 0.1-0.5, красный >0.5

12. **Rockets Assembled Rate** - скорость сборки ракет (ракет/минуту)
    - Запрос: `sum(rate(rockets_assembled_total[5m]))`
    - Пороги: красный (низкая активность), желтый 0.1-0.5, зеленый >0.5

#### Распределения и аналитика:

13. **Orders Status Distribution** - распределение заказов по статусам (круговая диаграмма)
    - Запросы: 
      - `sum(orders_created_total) by (status)`
      - `sum(orders_paid_total) by (status)`
      - `sum(orders_cancelled_total) by (status)`

14. **Order Value Distribution** - распределение стоимости заказов по перцентилям
    - Запросы:
      - `histogram_quantile(0.50, sum(order_value_bucket) by (le))` - 50-й перцентиль
      - `histogram_quantile(0.95, sum(order_value_bucket) by (le))` - 95-й перцентиль
      - `histogram_quantile(0.99, sum(order_value_bucket) by (le))` - 99-й перцентиль
    - Единицы: USD

15. **Revenue by Currency** - выручка по валютам во времени
    - Запрос: `sum(orders_revenue_total) by (currency)`
    - Единицы: USD

16. **Assembly Duration Distribution** - распределение длительности сборки ракет по временным интервалам
    - Запрос: `count(assembly_duration_seconds_bucket) by (le)`
    - Тип: гистограмма
    - Диапазон: 0-30 секунд

### 2. Technical Metrics Dashboard (`technical_metrics.json`)

Дашборд для технических команд, отображает производительность и технические показатели.

#### HTTP API метрики:

1. **HTTP Requests Rate by Endpoint** - общее количество HTTP запросов в секунду с детализацией по endpoint
   - Запрос: `sum(rate(http_server_request_count[5m])) by (http_route) or sum(rate(otelogen_server_request_count[5m])) by (http_route)`
   - Единицы: запросов/сек
   - Тип: временной ряд

2. **HTTP Errors Rate by Endpoint** - общее количество ошибок HTTP запросов в секунду с детализацией по endpoint и HTTP кодам ошибок
   - Запрос: `sum(rate(http_server_errors_count[5m])) by (http_route, http_response_status_code) or sum(rate(otelogen_server_errors_count[5m])) by (http_route, http_response_status_code)`
   - Единицы: ошибок/сек
   - Пороги: зеленый (норма), красный >0.1

#### gRPC метрики:

3. **gRPC Requests Rate by Service** - общее количество gRPC запросов в секунду с детализацией по клиентам
   - Запрос: `sum(rate(grpc_client_requests_total[5m])) by (service)`
   - Единицы: запросов/сек
   - Тип: временной ряд

4. **gRPC Errors Rate by Service** - общее количество ошибок gRPC запросов в секунду с детализацией по сервисам и gRPC кодам ошибок
   - Запрос: `sum(rate(grpc_client_request_errors_total[5m])) by (service, grpc_code)`
   - Единицы: ошибок/сек
   - Пороги: зеленый (норма), красный >0.1

#### Kafka Producer метрики:

5. **Kafka Producer Messages Rate by Topic** - общее количество запросов в секунду в Kafka producer
   - Запрос: `sum(rate(kafka_producer_messages_total[5m])) by (topic)`
   - Единицы: сообщений/сек
   - Тип: временной ряд

6. **Kafka Producer Errors Rate by Topic** - общее количество ошибок в Kafka producer
   - Запрос: `sum(rate(kafka_producer_messages_failed_total[5m])) by (topic)`
   - Единицы: ошибок/сек
   - Пороги: зеленый (норма), красный >0.1

#### Kafka Consumer метрики:

7. **Kafka Consumer Messages Rate by Topic** - общее количество запросов в секунду в Kafka consumer
   - Запрос: `sum(rate(kafka_consumer_messages_total[5m])) by (topic, group_id)`
   - Единицы: сообщений/сек
   - Тип: временной ряд

8. **Kafka Consumer Errors Rate by Topic** - общее количество ошибок в Kafka consumer
   - Запрос: `sum(rate(kafka_consumer_messages_failed_total[5m])) by (topic, group_id)`
   - Единицы: ошибок/сек
   - Пороги: зеленый (норма), красный >0.1

#### Kafka Consumer Lag и ребалансировки:

9. **Kafka Consumer Lag** - график роста consumer lag
   - Запрос: `kafka_consumer_offset_lag`
   - Единицы: сообщений
   - Пороги: зеленый (норма), желтый >100, красный >1000
   - Детализация: по топикам, группам и партициям

10. **Kafka Consumer Rebalancing Count** - количество ребалансировок в consumer
    - Запрос: `sum(kafka_consumer_rebalancing_total) by (group_id)`
    - Единицы: количество
    - Пороги: зеленый (норма), желтый >1, красный >5

#### Настройки дашбордов:

**Business Metrics Dashboard:**
- **Обновление**: каждые 5 секунд
- **Временной диапазон**: последний час
- **Теги**: business, orders, assembly, revenue
- **UID**: business-metrics
- **Версия**: 40

**Technical Metrics Dashboard:**
- **Обновление**: каждые 5 секунд
- **Временной диапазон**: последний час
- **Теги**: technical, http, grpc, kafka, metrics
- **UID**: technical-metrics
- **Версия**: 1

**Общие настройки:**
- **Цветовая схема**: зелено-желто-красная (зеленый = норма, желтый = предупреждение, красный = аномалия)
- **Пороги предупреждений**: настроены индивидуально для каждого типа метрик
- **Datasource**: Prometheus (UID: prometheus-uid)
- **Timezone**: browser (локальное время пользователя)

