"use client";

import { FeaturePreview } from "@/components/feature-preview";
import { features } from "@/lib/features";
import { CreditCard, Lock } from "@/components/icons";

export default function TradePage() {
  // The full tradeable experience is built behind a flag and stays off until
  // the client has explicit securities/CFTC legal sign-off. See src/lib/features.ts.
  if (!features.trade) {
    return (
      <div className="mx-auto w-full max-w-2xl px-4 py-24 text-center sm:px-6">
        <div
          className="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl text-white shadow-lg"
          style={{ background: "var(--tab-trade)" }}
        >
          <Lock className="h-8 w-8" />
        </div>
        <h1 className="mt-6 text-3xl font-bold tracking-tight">Trade — pending legal review</h1>
        <p className="mx-auto mt-4 max-w-xl text-zinc-600 dark:text-zinc-400">
          The Trade arena — tradeable crypto-ID logos, NFT markers, and an in-app wallet — is
          fully scaffolded but intentionally disabled. Tradeable, value-bearing reputation
          instruments are a securities and CFTC question, not just a feature toggle, and creator
          tokens were deliberately reframed as non-transferable support points for compliance.
        </p>
        <p className="mx-auto mt-4 max-w-xl text-sm text-zinc-500 dark:text-zinc-500">
          It switches on per-environment via <code className="rounded bg-zinc-100 px-1.5 py-0.5 dark:bg-zinc-800">NEXT_PUBLIC_ENABLE_TRADE=true</code> once
          legal clearance exists.
        </p>
      </div>
    );
  }

  return (
    <FeaturePreview
      accent="var(--tab-trade)"
      Icon={CreditCard}
      title="Trade Arena"
      tagline="Your logo is your asset."
      intro="The on-chain universe where creator logos become crypto IDs and NFT markers, with a secure in-app wallet to trade virtual assets and points across compatible chains."
      sections={[
        {
          heading: "Crypto IDs & NFT markers",
          items: [
            "Mint your logo as an NFT passport marker",
            "Compatible with public & private blockchains",
            "Provenance-backed by your Vault anchors",
          ],
        },
        {
          heading: "In-app wallet",
          items: [
            "Trade virtual assets and arena points",
            "Transaction history and live balances",
            "Wallet shared with the VR Community Arena",
          ],
        },
      ]}
      statusLabel="Enabled — legal-cleared build"
      note="Reminder: only ship transferable/tradeable instruments with documented securities & CFTC sign-off."
    />
  );
}
