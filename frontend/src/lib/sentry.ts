const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

let initialized = false;

export function initSentry() {
  if (initialized) return;
  initialized = true;

  window.addEventListener("error", (event) => {
    reportError("error", event.message, event.error?.stack, window.location.href);
  });

  window.addEventListener("unhandledrejection", (event) => {
    const message = event.reason?.message || String(event.reason);
    const stack = event.reason?.stack;
    reportError("error", message, stack, window.location.href);
  });
}

export function captureException(error: unknown) {
  const err = error instanceof Error ? error : new Error(String(error));
  reportError("error", err.message, err.stack, window.location.href);
}

function reportError(level: string, message: string, stack?: string, url?: string) {
  try {
    const body = JSON.stringify({
      source: "frontend",
      level,
      message,
      stack: stack || null,
      url: url || null,
    });
    if (navigator.sendBeacon) {
      navigator.sendBeacon(`${API_URL}/api/errors`, body);
    } else {
      fetch(`${API_URL}/api/errors`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body,
        keepalive: true,
      }).catch(() => {});
    }
  } catch {
    // Silently fail — don't cause more errors from error reporting
  }
}
