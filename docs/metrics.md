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

#### Клиентские метрики:

| Метрика | Тип | Описание | Единицы | Метки |
|---------|-----|----------|---------|-------|
| `http_client_request_duration` | Histogram | Длительность HTTP запросов клиента | ms | method, route, status_code |
| `http_client_request_count` | Counter | Количество HTTP запросов клиента | 1 | method, route, status_code |
| `http_client_errors_count` | Counter | Количество ошибок HTTP запросов клиента | 1 | method, route, status_code |

### Метки HTTP метрик:

- `http.request.method` - HTTP метод (GET, POST, PUT, DELETE)
- `http.route` - маршрут API (например, `/api/v1/orders`)
- `http.response.status_code` - код ответа (200, 400, 500)
- `otelogen.operation_id` - ID операции (CreateOrder, PayOrder)

## Дашборды Grafana

Проект использует два специализированных дашборда для разных аудиторий и задач:

### 1. Business Metrics Dashboard (`business_metrics.json`)

Дашборд для бизнес-команды и руководства, отображает ключевые бизнес-показатели и KPI.

#### Скорости (Rates) - Мониторинг активности:

1. **Orders Created** - скорость создания заказов (заказов/минуту)
   - Запрос: `sum(rate(orders_created_total[5m]))`

2. **Orders Paid** - скорость оплаты заказов (заказов/минуту)
   - Запрос: `sum(rate(orders_paid_total[5m]))`

3. **Orders Cancelled** - скорость отмены заказов (заказов/минуту)
   - Запрос: `sum(rate(orders_cancelled_total[5m]))`

4. **Total Revenue** - скорость получения выручки (валюта/минуту)
   - Запрос: `sum(rate(orders_revenue_total[5m]))`

5. **Order Conversion Rate** - конверсия заказов (%)
   - Запрос: `sum(rate(orders_paid_total[5m])) / sum(rate(orders_created_total[5m])) * 100`

6. **Assembly In Progress** - текущее количество ракет в сборке
   - Запрос: `sum(assembly_in_progress)`

7. **Rockets Assembled** - скорость сборки ракет (ракет/минуту)
   - Запрос: `sum(rate(rockets_assembled_total[5m]))`

8. **Assembly Errors** - скорость ошибок сборки (ошибок/минуту)
   - Запрос: `sum(rate(assembly_errors_total[5m]))`

#### Распределения и тренды:

9. **Order Value Distribution** - распределение стоимости заказов (p50, p95, p99)
   - Запросы: `histogram_quantile(0.50/0.95/0.99, rate(order_value_bucket[5m]))`

10. **Revenue by Currency** - выручка по валютам
    - Запрос: `sum(rate(orders_revenue_total[5m])) by (currency)`

11. **Assembly by Rocket Type** - сборка по типам ракет
    - Запрос: `sum(rate(rockets_assembled_total[5m])) by (rocket_type)`

12. **Assembly Errors by Type** - ошибки сборки по типам
    - Запрос: `sum(rate(assembly_errors_total[5m])) by (error_type, rocket_type)`

#### Общие количества (Totals) - Итоговая статистика:

13. **Total Orders Created** - общее количество созданных заказов
    - Запрос: `sum(orders_created_total)`

14. **Total Orders Paid** - общее количество оплаченных заказов
    - Запрос: `sum(orders_paid_total)`

15. **Total Orders Cancelled** - общее количество отмененных заказов
    - Запрос: `sum(orders_cancelled_total)`

16. **Total Revenue** - общая сумма выручки
    - Запрос: `sum(orders_revenue_total)`

17. **Total Rockets Assembled** - общее количество собранных ракет
    - Запрос: `sum(rockets_assembled_total)`

18. **Total Assembly Errors** - общее количество ошибок сборки
    - Запрос: `sum(assembly_errors_total)`

#### Круговые диаграммы:

19. **Orders Status Distribution** - распределение заказов по статусам
    - Запросы: `sum(orders_created_total/orders_paid_total/orders_cancelled_total) by (status)`

20. **Rockets by Type Distribution** - распределение ракет по типам
    - Запрос: `sum(rockets_assembled_total) by (rocket_type)`

### 2. Technical Metrics Dashboard (`technical_metrics.json`)

Дашборд для технических команд, отображает производительность и технические показатели.

#### HTTP Server метрики:

1. **HTTP Server Requests** - запросы к серверу по методам и маршрутам
   - Запрос: `sum(rate(http_server_request_count[5m])) by (http_request_method, http_route)`

2. **HTTP Server Errors** - ошибки сервера с кодами ответов
   - Запрос: `sum(rate(http_server_errors_count[5m])) by (http_request_method, http_route, http_response_status_code)`

3. **HTTP Server Response Time** - время ответа сервера (p50, p95, p99)
   - Запросы: `histogram_quantile(0.50/0.95/0.99, rate(http_server_request_duration_bucket[5m])) by (http_request_method, http_route)`

#### HTTP Client метрики:

4. **HTTP Client Requests** - запросы клиентов
   - Запрос: `sum(rate(http_client_request_count[5m])) by (http_request_method, http_route)`

5. **HTTP Client Errors** - ошибки клиентов
   - Запрос: `sum(rate(http_client_errors_count[5m])) by (http_request_method, http_route, http_response_status_code)`

6. **HTTP Client Response Time** - время ответа клиентов (p50, p95, p99)
   - Запросы: `histogram_quantile(0.50/0.95/0.99, rate(http_client_request_duration_bucket[5m])) by (http_request_method, http_route)`

#### Метрики производительности:

7. **Order Processing Duration** - время обработки заказов (p50, p95, p99)
   - Запросы: `histogram_quantile(0.50/0.95/0.99, rate(order_duration_seconds_bucket[5m]))`

8. **Assembly Duration** - длительность сборки (p50, p95, p99)
   - Запросы: `histogram_quantile(0.50/0.95/0.99, rate(assembly_duration_seconds_bucket[5m]))`

#### Аналитические метрики:

9. **Error Rate by Endpoint** - процент ошибок по endpoint'ам
   - Запрос: `sum(rate(http_server_errors_count[5m])) by (http_route) / sum(rate(http_server_request_count[5m])) by (http_route) * 100`

10. **Request Rate by Status Code** - запросы по статус кодам
    - Запрос: `sum(rate(http_server_request_count[5m])) by (http_response_status_code)`

11. **Average Response Time by Operation** - среднее время ответа по операциям
    - Запрос: `histogram_quantile(0.50, rate(http_server_request_duration_bucket[5m])) by (otelogen_operation_id)`

#### Настройки дашбордов:

- **Обновление**: каждые 5 секунд
- **Временной диапазон**: последний час
- **Цветовая схема**: зелено-желто-красная (зеленый = норма, желтый = предупреждение, красный = аномалия)
- **Пороги предупреждений**: настроены индивидуально для каждого типа метрик

