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
  logger.info({
    at: 'metrics#createMetricsServer',
    message: `Metrics server for ${serviceName} on port ${port}`,
  });
  app.get('/health', (_req, res) => {
    res.json({ ok: true });
  });
  
  app.listen(port, () => {
    logger.info({
      at: 'metrics#createMetricsServer',
      message: `Metrics server for ${serviceName} listening on port ${port}`,
    });
  });
  
  return app;
}

export function trackEvent(serviceName: string, eventType: string) {
  eventCounter.labels(serviceName, eventType).inc();
}

export function setActiveConnections(serviceName: string, value: number) {
  activeConnections.labels(serviceName).set(value);
}

