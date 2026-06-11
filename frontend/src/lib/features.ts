/**
 * Front-end feature flags.
 *
 * TRADE is intentionally OFF by default. The client requested a tradeable
 * crypto-ID / NFT / Polymarket-compatible wallet, but the most recent backend
 * commit deliberately reframed creator tokens as *non-transferable* support
 * points for SEC compliance. Tradeable, value-bearing reputation instruments
 * are a securities/CFTC question, not just a feature toggle.
 *
 * The Trade tab and wallet are built behind this flag so nothing ships to the
 * public until the client has explicit legal sign-off. Enable per-environment
 * via NEXT_PUBLIC_ENABLE_TRADE=true once that clearance exists.
 */
export const features = {
  trade: process.env.NEXT_PUBLIC_ENABLE_TRADE === "true",
} as const;

export type FeatureName = keyof typeof features;
