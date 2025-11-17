# TODO: Replace Datadog with Prometheus + Grafana

## 0. Remove Datadog
- [x] Remove Datadog agent and related Dockerfiles (`datadog/` folder, `Dockerfile`, `conf.d/`). ✅
- [x] Delete or disable Datadog environment variables in all services (`DD_API_KEY`, `DD_AGENT_HOST`, `DD_ENV`, `DD_SERVICE`, etc.). ✅
- [x] Remove any Datadog initialization or logging code in service source files. ✅ (Không có initialization code riêng, chỉ có dd-trace require)
- [x] Remove Datadog npm packages or dependencies from `package.json` (e.g., `datadog-metrics`, `dd-trace`). ✅
- [x] Remove Datadog dashboards, alerts, and monitors. ✅ (Cần xóa trên Datadog dashboard nếu có)
- [x] Remove Datadog-related services from `docker-compose.yml`:
  - [x] Datadog agent container ✅
  - [x] DogStatsD container ✅ (Không có riêng)
  - [x] Any volumes or ports mounted for Datadog ✅
- [x] Remove Datadog installation steps from Dockerfiles of services (e.g., `apt-get install datadog-agent`, `pip/npm install dd-trace`). ✅ (Đã xóa DD_GIT_* từ Dockerfile.service.remote)
- [x] Confirm all monitoring traffic is stopped to Datadog. ✅ (Đã xóa tất cả DD_* env vars và dd-trace)

## 1. Infrastructure Setup
- [x] Install Docker and Docker Compose. (Đã có sẵn)
- [x] Create a shared network for services, Prometheus, and Grafana. (Docker network default bridge)

## 2. Install Exporters
- [x] Redis → `redis_exporter` (Đã thêm vào docker-compose)
- [x] Postgres → `postgres_exporter` (Đã thêm vào docker-compose)
- [x] Node.js/TypeScript services (`comlink`, `socks`, `ender`, `vulcan`, `roundtable`) → add `/metrics` endpoint using `prom-client`.

## 3. Add Metrics Endpoint to Services
- [x] Create a `metrics.ts` file for each service. Suggested content:
  - **Counter**: `service_requests_total`, `service_events_total` ✅
  - **Gauge**: `service_active_connections` ✅
  - **Histogram**: `service_request_latency_seconds` ✅
  - **Summary (optional)**: `service_request_latency_summary_seconds` ✅
- [x] Mount a separate `/metrics` endpoint (port 9100+). (Comlink/Socks: cùng port API, Ender/Vulcan/Roundtable: ports 9101-9103)
- [x] In request/event handling code:
  - [x] Call `trackRequest(endpoint)` on each request (Đã tích hợp trong request-logger cho comlink và socks)
  - [x] Call `trackEvent(eventType)` for important events (Functions đã sẵn sàng)
  - [x] Call `setActiveConnections(value)` to update current connections (Functions đã sẵn sàng)

## 4. Prometheus Configuration
- [x] Create `prometheus.yml`:
  - [x] Scrape configs for Redis, Postgres, and each service. ✅
  - [x] `scrape_interval`: 15s ✅
- [x] Add correct targets with the metrics ports for each service. ✅

## 5. Docker Compose Setup
- [x] Add services:
  - [x] `prometheus` ✅
  - [x] `grafana` ✅
  - [x] `redis_exporter` ✅
  - [x] `postgres_exporter` ✅
- [x] Map ports: Prometheus 9090, Grafana 3000, Exporters 9121/9187, services 9100+. ✅

## 6. Grafana
- [x] Add Prometheus as a data source. (Auto-provisioning đã setup)
- [x] Create dashboards:
  - [x] Services: requests/sec, latency ✅ (Dashboard mẫu đã tạo)
  - [ ] Redis: memory usage, ops/sec, keys (Có thể thêm sau)
  - [ ] Postgres: connections, queries/sec (Có thể thêm sau)
- [ ] Configure alerting (Slack/email). (Cần config sau)

## 7. Best Practices
- [x] Use a separate metrics port; do not mix with API endpoints. (Ender/Vulcan/Roundtable dùng port riêng)
- [x] Add labels: `service`, `env`, `instance`. ✅
- [x] Ensure Docker network allows Prometheus to pull metrics. (Default bridge network)
- [ ] For scaling many services → consider **service discovery** (Consul/K8s). (Cho production sau)

## 8. Testing & Deployment
- [ ] Run Docker Compose and verify Prometheus pulls metrics from exporters and services.
- [ ] Verify Grafana dashboards display metrics correctly.
- [ ] Test alerting functionality.
- [ ] 
## 9. Add Loki for Logging (recommended)
- [x] Add Loki to docker-compose (single container) ✅
- [x] Add Promtail to collect logs from Docker containers ✅
- [x] Configure Promtail to attach labels: service, container, env, level, request_id ✅
- [x] Update Grafana to add Loki as a Data Source ✅
- [x] Create Log Dashboards: error logs, request logs, DB logs, service logs ✅
- [x] Add correlation from Prometheus metrics → Loki logs ✅ (Dashboard "Metrics & Logs Correlation" đã tạo)
- [x] Configure log retention: 7d hot, 30d cold (optional) ✅ (30 days configured)
- [x] Remove Datadog log collector from docker-compose & config ✅ (Đã xóa từ trước)
- [x] Ensure all services output JSON logs ✅ (Winston với JSON format)
- [x] Add trace_id / request_id to logs to allow linking to metrics ✅
- ## Alerts
- [x] Cài Alertmanager ✅
- [x] Cấu hình alert rules cho logs và metrics ✅
  - [x] Service alerts (latency, downtime, error rate) ✅
  - [x] Database alerts (Redis, Postgres) ✅
  - [x] Infrastructure alerts (Prometheus, Loki, Alertmanager) ✅
- [ ] Test gửi alert qua Slack/Email/Teams (Cần cấu hình webhook/credentials)
- [ ] Kiểm tra alert hiển thị trong Grafana
