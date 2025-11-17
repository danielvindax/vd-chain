/**
 * metrics.ts
 * 
 * Purpose: Export metrics for Prometheus scraping.
 * 
 * Includes:
 * 1. Counter: count requests, events
 * 2. Gauge: current value (e.g., active connections, queue size)
 * 3. Histogram: measure request latency
 * 4. Summary: quick latency statistics (optional)
 * 
 * Usage:
 * 1. Import and mount /metrics endpoint in Express app or HTTP service.
 * 2. When service processes request or event, increment Counter/Gauge/observe Histogram.
 */

import express from 'express';
import client from 'prom-client';

const register = client.register;

// ==== 1. Counter ====
// Count total requests by service
export const requestCounter = new client.Counter({
  name: 'service_requests_total',
  help: 'Total number of requests processed',
  labelNames: ['service', 'endpoint', 'method', 'status_code']
});

// Count important events, e.g., order processed
export const eventCounter = new client.Counter({
  name: 'service_events_total',
  help: 'Total number of events processed',
  labelNames: ['service', 'event_type']
});

// ==== 2. Gauge ====
// Current value, e.g., active connections, cache size
export const activeConnections = new client.Gauge({
  name: 'service_active_connections',
  help: 'Current number of active connections',
  labelNames: ['service']
});

// ==== 3. Histogram ====
// Measure request latency
export const requestLatency = new client.Histogram({
  name: 'service_request_latency_seconds',
  help: 'Request latency in seconds',
  labelNames: ['service', 'endpoint', 'method'],
  buckets: [0.005, 0.01, 0.05, 0.1, 0.5, 1, 5] // buckets from 5ms to 5s
});

// ==== 4. Summary ==== (optional)
// Similar to Histogram but focused on percentiles
export const requestLatencySummary = new client.Summary({
  name: 'service_request_latency_summary_seconds',
  help: 'Latency summary',
  labelNames: ['service', 'endpoint', 'method'],
  percentiles: [0.5, 0.9, 0.99]
});

// ==== 5. Metrics endpoint ====
// Prometheus scrapes metrics from here
export function createMetricsRouter(serviceName: string): express.Router {
  const router = express.Router();
  
  router.get('/metrics', async (_req: express.Request, res: express.Response) => {
    res.set('Content-Type', register.contentType);
    res.end(await register.metrics());
  });
  
  return router;
}

// ==== 6. Example usage ====
// When processing request or event, call as follows:

export function trackRequest(
  serviceName: string,
  endpoint: string,
  method: string,
  statusCode: number
) {
  requestCounter.labels(serviceName, endpoint, method, statusCode.toString()).inc();
}

export function trackRequestLatency(
  serviceName: string,
  endpoint: string,
  method: string
): () => void {
  const end = requestLatency.startTimer({ service: serviceName, endpoint, method });
  const endSummary = requestLatencySummary.startTimer({ service: serviceName, endpoint, method });
  
  return () => {
    end();
    endSummary();
  };
}

export function trackEvent(serviceName: string, eventType: string) {
  eventCounter.labels(serviceName, eventType).inc();
}

export function setActiveConnections(serviceName: string, value: number) {
  activeConnections.labels(serviceName).set(value);
}

