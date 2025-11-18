import { logger } from '@dydxprotocol-indexer/base';
import express from 'express';
import { trackRequest, trackRequestLatency } from '../metrics';
import { ResponseWithBody } from 'src/types';

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
  request.metricsTimer = trackRequestLatency('socks', endpoint, request.method);

  response.on('finish', () => {
    // Track metrics
    const endpoint = request.path || request.url.split('?')[0];
    trackRequest('socks', endpoint, request.method, response.statusCode);
    if (request.metricsTimer) {
      request.metricsTimer();
    }
    const protocol: string = request.protocol;
    const host: string | undefined = request.get('host');
    const url: string = request.originalUrl;
    const fullUrl: string = `${protocol}://${host}${url}`;

    // Convert RegExpMatchArray | null into true/false (boolean).
    const isError: boolean = !!response.statusCode.toString().match(/^[^2]/);
    if (request.method !== 'GET') {
      logger.info({
        at: 'requestLogger#logRequest',
        request_id: (request as any).id || undefined,
        message: {
          request: {
            url: fullUrl,
            method: request.method,
            headers: request.headers,
            query: request.query,
            body: request.body,
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
