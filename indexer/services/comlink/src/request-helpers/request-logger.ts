import { logger, safeJsonStringify } from '@dydxprotocol-indexer/base';
import express from 'express';

import config from '../config';
import { trackRequest, trackRequestLatency } from '../metrics';
import { ResponseWithBody } from '../types';

// Store latency timer in request object
declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace Express {
    interface Request {
      metricsTimer?: () => void;
    }
  }
}

export default (
  request: express.Request,
  response: express.Response,
  next: express.NextFunction,
) => {
  // Start latency timer
  const endpoint = request.path || request.url.split('?')[0];
  request.metricsTimer = trackRequestLatency('comlink', endpoint, request.method);

  response.on('finish', () => {
    // Track metrics
    const endpoint = request.path || request.url.split('?')[0];
    trackRequest('comlink', endpoint, request.method, response.statusCode);
    if (request.metricsTimer) {
      request.metricsTimer();
    }
    const { protocol } : { protocol: string } = request;
    const host: string | undefined = request.get('host');
    const url: string = request.originalUrl;
    const fullUrl: string = `${protocol}://${host}${url}`;

    const isError: RegExpMatchArray | null = response.statusCode.toString().match(/^[^2]/);
    // Don't log GET requests unless configured to
    const shouldLogMethod: boolean = request.method !== 'GET' || config.LOG_GETS;
    if (shouldLogMethod || response.statusCode !== 200) {
      logger.info({
        at: 'requestLogger#logRequest',
        request_id: (request as any).id || undefined,
        message: {
          request: {
            url: fullUrl,
            method: request.method,
            headers: request.headers,
            query: request.query,
            body: safeJsonStringify(request.body),
          },
          response: {
            statusCode: response.statusCode,
            errorBody: isError && (response as ResponseWithBody).body,
            statusMessage: response.statusMessage,
            headers: response.getHeaders(),
          },
        },
      });
    }
  });

  return next();
};
