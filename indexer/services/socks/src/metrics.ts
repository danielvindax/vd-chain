/**
 * metrics.ts
 * 
 * Purpose: Export metrics for Prometheus scraping.
 */

import { logger } from '@dydxprotocol-indexer/base';
import express from 'express';
import client from 'prom-client';

const register = client.register;

export const requestCounter = new client.Counter({
  name: 'service_requests_total',
  help: 'Total number of requests processed',
  labelNames: ['service', 'endpoint', 'method', 'status_code']
});

export const eventCounter = new client.Counter({
  name: 'service_events_total',
  help: 'Total number of events processed',
  labelNames: ['service', 'event_type']
});

export const activeConnections = new client.Gauge({
  name: 'service_active_connections',
  help: 'Current number of active connections',
  labelNames: ['service']
});

export const requestLatency = new client.Histogram({
  name: 'service_request_latency_seconds',
  help: 'Request latency in seconds',
  labelNames: ['service', 'endpoint', 'method'],
  buckets: [0.005, 0.01, 0.05, 0.1, 0.5, 1, 5]
});

export const requestLatencySummary = new client.Summary({
  name: 'service_request_latency_summary_seconds',
  help: 'Latency summary',
  labelNames: ['service', 'endpoint', 'method'],
  percentiles: [0.5, 0.9, 0.99]
});

export function createMetricsRouter(serviceName: string): express.Router {
  const router = express.Router();
  
  router.get('/metrics', async (_req, res) => {
    res.set('Content-Type', register.contentType);
    res.end(await register.metrics());
  });
  
  return router;
}

export function createMetricsServer(serviceName: string, port: number): express.Application {
  const app = express();
  logger.info({
    at: 'metrics#createMetricsServer',
    message: `Creating metrics server for ${serviceName} on port ${port}`,
  });
  app.get('/metrics', async (_req: express.Request, res: express.Response) => {
    res.set('Content-Type', register.contentType);
    res.end(await register.metrics());
  });
  app.get('/health', (_req, res) => {
    res.json({ ok: true });
  });
  
  app.listen(port, '0.0.0.0', () => {
    logger.info({
      at: 'metrics#createMetricsServer',
      message: `Metrics server for ${serviceName} listening on port ${port}`,
    });
  });
  
  return app;
}

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

