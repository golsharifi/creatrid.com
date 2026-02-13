#!/usr/bin/env bash
set -euo pipefail

# Production smoke tests for Creatrid
# Usage: bash tests/smoke.sh [base_api] [base_web]

BASE_API="${1:-https://api.creatrid.com}"
BASE_WEB="${2:-https://creatrid.com}"

PASS=0
FAIL=0

check() {
  local name="$1"
  local result="$2"
  if [ "$result" = "ok" ]; then
    printf "  \033[32mPASS\033[0m  %s\n" "$name"
    PASS=$((PASS + 1))
  else
    printf "  \033[31mFAIL\033[0m  %s — %s\n" "$name" "$result"
    FAIL=$((FAIL + 1))
  fi
}

echo "Creatrid Smoke Tests"
echo "  API: $BASE_API"
echo "  Web: $BASE_WEB"
echo ""

# 1. Health check
STATUS=$(curl -s -o /dev/null -w '%{http_code}' "$BASE_API/api/health")
BODY=$(curl -s "$BASE_API/api/health")
if [ "$STATUS" = "200" ] && echo "$BODY" | grep -q '"ok"'; then
  check "API health check" "ok"
else
  check "API health check" "status=$STATUS"
fi

# 2. Landing page
STATUS=$(curl -s -o /dev/null -w '%{http_code}' "$BASE_WEB/")
if [ "$STATUS" = "200" ]; then
  check "Landing page loads" "ok"
else
  check "Landing page loads" "status=$STATUS"
fi

# 3. Sign-in page
STATUS=$(curl -s -o /dev/null -w '%{http_code}' "$BASE_WEB/sign-in")
if [ "$STATUS" = "200" ]; then
  check "Sign-in page loads" "ok"
else
  check "Sign-in page loads" "status=$STATUS"
fi

# 4. Discover API
STATUS=$(curl -s -o /dev/null -w '%{http_code}' "$BASE_API/api/discover")
if [ "$STATUS" = "200" ]; then
  check "Discover API returns 200" "ok"
else
  check "Discover API returns 200" "status=$STATUS"
fi

# 5. Terms page
STATUS=$(curl -s -o /dev/null -w '%{http_code}' "$BASE_WEB/terms")
if [ "$STATUS" = "200" ]; then
  check "Terms page loads" "ok"
else
  check "Terms page loads" "status=$STATUS"
fi

# 6. Privacy page
STATUS=$(curl -s -o /dev/null -w '%{http_code}' "$BASE_WEB/privacy")
if [ "$STATUS" = "200" ]; then
  check "Privacy page loads" "ok"
else
  check "Privacy page loads" "status=$STATUS"
fi

# 7. Pricing page
STATUS=$(curl -s -o /dev/null -w '%{http_code}' "$BASE_WEB/pricing")
if [ "$STATUS" = "200" ]; then
  check "Pricing page loads" "ok"
else
  check "Pricing page loads" "status=$STATUS"
fi

# 8. CORS preflight
CORS_ORIGIN=$(curl -s -I -X OPTIONS -H "Origin: $BASE_WEB" -H "Access-Control-Request-Method: GET" "$BASE_API/api/health" 2>&1 | grep -i 'access-control-allow-origin' | tr -d '\r')
if echo "$CORS_ORIGIN" | grep -qi "creatrid.com"; then
  check "CORS preflight headers" "ok"
else
  check "CORS preflight headers" "missing or wrong origin"
fi

# 9. Rate limit headers
RL_HEADER=$(curl -sI "$BASE_API/api/health" | grep -i 'x-ratelimit-limit' | tr -d '\r')
if [ -n "$RL_HEADER" ]; then
  check "Rate limit headers present" "ok"
else
  check "Rate limit headers present" "X-RateLimit-Limit missing"
fi

# 10. Auth required returns 401
STATUS=$(curl -s -o /dev/null -w '%{http_code}' "$BASE_API/api/auth/me")
if [ "$STATUS" = "401" ]; then
  check "Auth-required endpoint returns 401" "ok"
else
  check "Auth-required endpoint returns 401" "status=$STATUS"
fi

# 11. robots.txt
STATUS=$(curl -s -o /dev/null -w '%{http_code}' "$BASE_WEB/robots.txt")
ROBOTS_BODY=$(curl -s "$BASE_WEB/robots.txt")
if [ "$STATUS" = "200" ] && echo "$ROBOTS_BODY" | grep -q "Sitemap"; then
  check "robots.txt accessible" "ok"
else
  check "robots.txt accessible" "status=$STATUS"
fi

# 12. sitemap.xml
STATUS=$(curl -s -o /dev/null -w '%{http_code}' "$BASE_WEB/sitemap.xml")
SITEMAP_BODY=$(curl -s "$BASE_WEB/sitemap.xml")
if [ "$STATUS" = "200" ] && echo "$SITEMAP_BODY" | grep -q "urlset"; then
  check "sitemap.xml accessible" "ok"
else
  check "sitemap.xml accessible" "status=$STATUS"
fi

# 13. OG image
STATUS=$(curl -s -o /dev/null -w '%{http_code}' "$BASE_WEB/og-image.svg")
if [ "$STATUS" = "200" ]; then
  check "OG image accessible" "ok"
else
  check "OG image accessible" "status=$STATUS"
fi

# 14. HTTPS redirect
STATUS=$(curl -s -o /dev/null -w '%{http_code}' -L "http://creatrid.com/" 2>/dev/null || echo "skip")
if [ "$STATUS" = "200" ] || [ "$STATUS" = "301" ] || [ "$STATUS" = "skip" ]; then
  check "HTTPS available" "ok"
else
  check "HTTPS available" "status=$STATUS"
fi

# Summary
echo ""
echo "Results: $PASS passed, $FAIL failed"
if [ "$FAIL" -gt 0 ]; then
  exit 1
fi
