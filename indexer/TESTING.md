# Hướng dẫn Test Prometheus + Grafana

## Bước 1: Cài đặt Dependencies

```bash
# Cài đặt tất cả dependencies (bao gồm prom-client)
pnpm install
```

## Bước 2: Build tất cả Services

```bash
# Build tất cả services và packages
pnpm run build:all
```

## Bước 3: Chạy Docker Compose

```bash
# Chạy tất cả services bao gồm Prometheus và Grafana
docker-compose -f docker-compose-local-deployment.yml up
```

Hoặc chạy ở background:
```bash
docker-compose -f docker-compose-local-deployment.yml up -d
```

## Bước 4: Kiểm tra Services đang chạy

```bash
# Xem tất cả containers
docker-compose -f docker-compose-local-deployment.yml ps

# Kiểm tra logs của một service cụ thể
docker-compose -f docker-compose-local-deployment.yml logs comlink
docker-compose -f docker-compose-local-deployment.yml logs prometheus
docker-compose -f docker-compose-local-deployment.yml logs grafana
```

## Bước 5: Test Metrics Endpoints

### Test metrics từ các services:

```bash
# Comlink (cùng port với API)
curl http://localhost:3002/metrics

# Socks (cùng port với API)
curl http://localhost:3003/metrics

# Ender (port riêng)
curl http://localhost:9101/metrics

# Vulcan (port riêng)
curl http://localhost:9102/metrics

# Roundtable (port riêng)
curl http://localhost:9103/metrics
```

Bạn sẽ thấy output dạng Prometheus metrics như:
```
# HELP service_requests_total Tổng số request đã xử lý
# TYPE service_requests_total counter
service_requests_total{service="comlink",endpoint="/health",method="GET",status_code="200"} 1

# HELP service_request_latency_seconds Latency của request (giây)
# TYPE service_request_latency_seconds histogram
...
```

## Bước 6: Kiểm tra Prometheus

1. Mở browser: http://localhost:9090

2. Vào **Status > Targets** để xem các targets Prometheus đang scrape:
   - Tất cả services phải có state = **UP** (màu xanh)

3. Test query trong Prometheus:
   ```
   # Số request mỗi giây
   rate(service_requests_total[1m])
   
   # Latency p95
   histogram_quantile(0.95, rate(service_request_latency_seconds_bucket[5m]))
   
   # Active connections
   service_active_connections
   ```

## Bước 7: Kiểm tra Grafana

1. Mở browser: http://localhost:3000

2. Đăng nhập:
   - Username: `admin`
   - Password: `admin`

3. Kiểm tra Data Source:
   - Vào **Configuration > Data Sources**
   - Prometheus phải được auto-provisioned và có status **Green**

4. Xem Dashboard:
   - Vào **Dashboards > Browse**
   - Tìm dashboard "Indexer Services Overview"
   - Hoặc tạo dashboard mới với các queries từ Bước 6

## Bước 8: Generate Traffic để test Metrics

### Test Comlink API:
```bash
# Gọi API để generate metrics
curl http://localhost:3002/health
curl http://localhost:3002/v4/time
```

### Test Socks WebSocket:
```bash
# Sử dụng websocat hoặc wscat
# Install: npm install -g wscat
wscat -c ws://localhost:3003
```

Sau khi generate traffic, refresh Prometheus và Grafana để xem metrics được update.

## Bước 9: Kiểm tra Exporters

### Redis Exporter:
```bash
curl http://localhost:9121/metrics | grep redis
```

### Postgres Exporter:
```bash
curl http://localhost:9187/metrics | grep postgres
```

## Troubleshooting

### Nếu metrics không hiển thị:

1. **Kiểm tra service có chạy không:**
   ```bash
   docker-compose -f docker-compose-local-deployment.yml ps
   ```

2. **Kiểm tra logs:**
   ```bash
   docker-compose -f docker-compose-local-deployment.yml logs <service-name>
   ```

3. **Kiểm tra metrics endpoint trực tiếp:**
   ```bash
   curl http://localhost:<port>/metrics
   ```

4. **Kiểm tra Prometheus targets:**
   - Vào http://localhost:9090/targets
   - Xem service nào có lỗi (màu đỏ)

### Nếu Prometheus không scrape được:

1. **Kiểm tra network:**
   ```bash
   docker network ls
   docker network inspect indexer_default
   ```

2. **Kiểm tra prometheus.yml:**
   ```bash
   cat prometheus.yml
   ```

3. **Reload Prometheus config:**
   ```bash
   curl -X POST http://localhost:9090/-/reload
   ```

### Nếu Grafana không kết nối được Prometheus:

1. **Kiểm tra Prometheus có chạy không:**
   ```bash
   curl http://localhost:9090/-/healthy
   ```

2. **Kiểm tra Grafana logs:**
   ```bash
   docker-compose -f docker-compose-local-deployment.yml logs grafana
   ```

3. **Test connection từ Grafana:**
   - Vào Configuration > Data Sources > Prometheus
   - Click "Save & Test"

## Cleanup

Để dừng và xóa tất cả containers:

```bash
docker-compose -f docker-compose-local-deployment.yml down

# Xóa cả volumes (xóa data Prometheus và Grafana)
docker-compose -f docker-compose-local-deployment.yml down -v
```

