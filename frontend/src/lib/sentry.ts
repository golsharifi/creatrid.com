import * as Sentry from "@sentry/react";

const SENTRY_DSN = process.env.NEXT_PUBLIC_SENTRY_DSN;

export function initSentry() {
  if (!SENTRY_DSN) return;

  Sentry.init({
    dsn: SENTRY_DSN,
    tracesSampleRate: 0.1,
    environment: process.env.NODE_ENV,
  });
}

export function captureException(error: unknown) {
  if (SENTRY_DSN) {
    Sentry.captureException(error);
  }
}
