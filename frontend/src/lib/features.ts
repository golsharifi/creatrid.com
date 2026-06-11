/**
 * Front-end feature flags.
 *
 * TRADE_EXTERNAL — the Trade tab itself is live with the closed-loop sticker
 * exchange (stickers ↔ arena points, fully in-platform, no securities
 * exposure). This flag gates only the EXTERNAL portion the client requested:
 * transferable crypto IDs, public NFT markets, Polymarket compatibility.
 * That requires securities/CFTC legal sign-off — the most recent backend
 * compliance work deliberately made creator tokens non-transferable (commit
 * 8da40a1). The client is obtaining clearance; enable per-environment with
 * NEXT_PUBLIC_ENABLE_TRADE=true once she has it in writing.
 *
 * STREAMING — the online-stations player on /music is fully built, but
 * streaming third-party music catalogs requires public-performance /
 * streaming licenses (e.g. PROs, label deals). The client is obtaining the
 * licenses herself; once licensed stream URLs are configured, enable with
 * NEXT_PUBLIC_ENABLE_STREAMING=true.
 */
export const features = {
  tradeExternal: process.env.NEXT_PUBLIC_ENABLE_TRADE === "true",
  streaming: process.env.NEXT_PUBLIC_ENABLE_STREAMING === "true",
} as const;

export type FeatureName = keyof typeof features;
