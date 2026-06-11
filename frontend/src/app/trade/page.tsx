"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { api } from "@/lib/api";
import { features } from "@/lib/features";
import type { StickerListing, UserSticker } from "@/lib/types";
import { CreditCard, ArrowRight, Lock, TrendingUp } from "@/components/icons";

const ACCENT = "var(--tab-trade)";

const REASON_LABELS: Record<string, string> = {
  explore: "Farm explore",
  upload: "Work uploaded",
  comment_received: "Comment on your work",
  like_received: "Like on your work",
  like_removed: "Like removed",
  trade_buy: "Sticker purchased",
  trade_sale: "Sticker sold",
};

export default function TradePage() {
  const { user } = useAuth();

  // Wallet
  const [points, setPoints] = useState(0);
  const [events, setEvents] = useState<{ delta: number; reason: string; createdAt: string }[]>([]);

  // Exchange
  const [listings, setListings] = useState<StickerListing[]>([]);
  const [myListings, setMyListings] = useState<StickerListing[]>([]);
  const [myStickers, setMyStickers] = useState<UserSticker[]>([]);
  const [sellSticker, setSellSticker] = useState("");
  const [sellPrice, setSellPrice] = useState("25");
  const [busy, setBusy] = useState(false);
  const [msg, setMsg] = useState("");
  const [error, setError] = useState("");

  const loadAll = useCallback(async () => {
    const [floor, wallet, mine, arena] = await Promise.all([
      api.trade.listings(),
      user ? api.arena.wallet() : Promise.resolve(null),
      user ? api.trade.myListings() : Promise.resolve(null),
      user ? api.arena.me() : Promise.resolve(null),
    ]);
    if (floor.data) setListings(floor.data.listings);
    if (wallet?.data) {
      setPoints(wallet.data.points);
      setEvents(wallet.data.events);
    }
    if (mine?.data) setMyListings(mine.data.listings);
    if (arena?.data) setMyStickers(arena.data.stickers);
  }, [user]);

  useEffect(() => {
    loadAll();
  }, [loadAll]);

  async function createListing() {
    const price = parseInt(sellPrice, 10);
    if (!sellSticker || !price || busy) return;
    setBusy(true);
    setError("");
    setMsg("");
    const res = await api.trade.createListing(sellSticker, price);
    setBusy(false);
    if (res.error) {
      setError(res.error);
      return;
    }
    setMsg("Listing created — your sticker is in escrow until it sells or you cancel.");
    setSellSticker("");
    loadAll();
  }

  async function buy(l: StickerListing) {
    if (busy) return;
    setBusy(true);
    setError("");
    setMsg("");
    const res = await api.trade.buy(l.id);
    setBusy(false);
    if (res.error) {
      setError(res.error);
      return;
    }
    setMsg(`You bought ${l.emoji} ${l.name} for ${l.pricePoints} points!`);
    loadAll();
  }

  async function cancel(id: string) {
    if (busy) return;
    setBusy(true);
    const res = await api.trade.cancelListing(id);
    setBusy(false);
    if (!res.error) {
      setMsg("Listing cancelled — sticker returned to your collection.");
      loadAll();
    }
  }

  return (
    <div className="mx-auto w-full max-w-5xl px-4 py-12 sm:px-6">
      <div className="flex flex-wrap items-center gap-3">
        <div
          className="flex h-12 w-12 items-center justify-center rounded-xl text-white shadow"
          style={{ background: ACCENT }}
        >
          <CreditCard className="h-6 w-6" />
        </div>
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Trade</h1>
          <p className="text-sm text-zinc-500 dark:text-zinc-400">
            The sticker exchange — trade passport stickers for points, in-platform.
          </p>
        </div>
        {user && (
          <div
            className="ml-auto rounded-full px-4 py-2 text-sm font-bold text-white"
            style={{ background: ACCENT }}
          >
            {points} pts
          </div>
        )}
      </div>

      {(msg || error) && (
        <p
          className={`mt-4 text-sm ${error ? "text-red-600 dark:text-red-400" : ""}`}
          style={error ? undefined : { color: ACCENT }}
        >
          {error || msg}
        </p>
      )}

      {/* ── Exchange floor (public) ── */}
      <div className="mt-8 rounded-2xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
        <h2 className="flex items-center gap-2 text-lg font-semibold">
          <TrendingUp className="h-5 w-5" style={{ color: ACCENT }} />
          Open listings
        </h2>
        {listings.length === 0 ? (
          <p className="mt-3 text-sm text-zinc-500 dark:text-zinc-400">
            Nothing on the floor right now — list one of your stickers below.
          </p>
        ) : (
          <div className="mt-4 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {listings.map((l) => (
              <div
                key={l.id}
                className="flex items-center gap-3 rounded-xl border border-zinc-200 p-4 dark:border-zinc-800"
              >
                <span className="text-3xl">{l.emoji}</span>
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm font-semibold">{l.name}</p>
                  <p className="text-xs text-zinc-500">
                    <span className="uppercase">{l.rarity}</span>
                    {l.seller && <> · @{l.seller}</>}
                  </p>
                </div>
                {user && l.sellerId !== user.id ? (
                  <button
                    onClick={() => buy(l)}
                    disabled={busy}
                    className="rounded-lg px-3 py-1.5 text-xs font-semibold text-white transition-opacity hover:opacity-90 disabled:opacity-50"
                    style={{ background: ACCENT }}
                  >
                    {l.pricePoints} pts
                  </button>
                ) : (
                  <span className="text-sm font-bold" style={{ color: ACCENT }}>
                    {l.pricePoints} pts
                  </span>
                )}
              </div>
            ))}
          </div>
        )}
      </div>

      {!user ? (
        <div className="mt-6 rounded-2xl border border-zinc-200 bg-white p-8 text-center dark:border-zinc-800 dark:bg-zinc-950">
          <p className="text-zinc-600 dark:text-zinc-400">
            Sign in to trade stickers and see your wallet.
          </p>
          <Link
            href="/sign-in"
            className="mt-4 inline-flex items-center gap-2 rounded-lg px-6 py-3 text-sm font-medium text-white transition-opacity hover:opacity-90"
            style={{ background: ACCENT }}
          >
            Sign in
            <ArrowRight className="h-4 w-4" />
          </Link>
        </div>
      ) : (
        <>
          {/* ── Sell a sticker ── */}
          <div className="mt-6 rounded-2xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
            <h2 className="text-lg font-semibold">Sell a sticker</h2>
            {myStickers.length === 0 ? (
              <p className="mt-3 text-sm text-zinc-500 dark:text-zinc-400">
                You have no stickers yet —{" "}
                <Link href="/vr-community" className="font-medium underline underline-offset-2">
                  explore the farm
                </Link>{" "}
                to find some.
              </p>
            ) : (
              <div className="mt-4 flex flex-wrap items-end gap-3">
                <label className="text-sm">
                  Sticker
                  <select
                    value={sellSticker}
                    onChange={(e) => setSellSticker(e.target.value)}
                    className="mt-1 block rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
                  >
                    <option value="">Choose…</option>
                    {myStickers.map((s) => (
                      <option key={s.id} value={s.id}>
                        {s.emoji} {s.name} (×{s.count}, {s.rarity})
                      </option>
                    ))}
                  </select>
                </label>
                <label className="text-sm">
                  Price (points)
                  <input
                    type="number"
                    min={1}
                    max={100000}
                    value={sellPrice}
                    onChange={(e) => setSellPrice(e.target.value)}
                    className="mt-1 block w-28 rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
                  />
                </label>
                <button
                  onClick={createListing}
                  disabled={busy || !sellSticker}
                  className="rounded-lg px-5 py-2.5 text-sm font-medium text-white transition-opacity hover:opacity-90 disabled:opacity-50"
                  style={{ background: ACCENT }}
                >
                  List it
                </button>
              </div>
            )}

            {myListings.filter((l) => l.status === "open").length > 0 && (
              <div className="mt-5">
                <p className="text-xs font-semibold uppercase tracking-wider text-zinc-500">
                  Your open listings
                </p>
                <ul className="mt-2 flex flex-col gap-2">
                  {myListings
                    .filter((l) => l.status === "open")
                    .map((l) => (
                      <li
                        key={l.id}
                        className="flex items-center gap-3 rounded-lg border border-zinc-100 px-3 py-2 text-sm dark:border-zinc-800"
                      >
                        <span>{l.emoji}</span>
                        <span className="font-medium">{l.name}</span>
                        <span className="text-zinc-500">{l.pricePoints} pts</span>
                        <button
                          onClick={() => cancel(l.id)}
                          className="ml-auto text-xs font-medium text-zinc-500 hover:text-red-500"
                        >
                          Cancel
                        </button>
                      </li>
                    ))}
                </ul>
              </div>
            )}
          </div>

          {/* ── Wallet: balance + transaction history ── */}
          <div className="mt-6 rounded-2xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold">Wallet</h2>
              <span className="text-xl font-bold" style={{ color: ACCENT }}>
                {points} pts
              </span>
            </div>
            <p className="mt-1 text-xs text-zinc-500 dark:text-zinc-400">
              Points are reputation currency: earned from exploring, uploading work,
              and community engagement. In-platform only — never cashable.
            </p>
            {events.length === 0 ? (
              <p className="mt-4 text-sm text-zinc-500 dark:text-zinc-400">No transactions yet.</p>
            ) : (
              <ul className="mt-4 flex max-h-72 flex-col gap-1.5 overflow-y-auto">
                {events.map((e, i) => (
                  <li
                    key={i}
                    className="flex items-center justify-between rounded-lg bg-zinc-50 px-3 py-2 text-sm dark:bg-zinc-900"
                  >
                    <span className="text-zinc-600 dark:text-zinc-400">
                      {REASON_LABELS[e.reason] || e.reason}
                    </span>
                    <span className="flex items-center gap-3">
                      <span className="text-xs text-zinc-400">
                        {new Date(e.createdAt).toLocaleDateString()}
                      </span>
                      <span
                        className={`w-14 text-right font-mono font-semibold ${
                          e.delta >= 0
                            ? "text-emerald-600 dark:text-emerald-400"
                            : "text-red-600 dark:text-red-400"
                        }`}
                      >
                        {e.delta >= 0 ? `+${e.delta}` : e.delta}
                      </span>
                    </span>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </>
      )}

      {/* ── External markets: built-out vision, pending the client's licenses ── */}
      <div className="mt-6 rounded-2xl border border-dashed border-zinc-300 bg-zinc-50 p-6 dark:border-zinc-700 dark:bg-zinc-900/50">
        <h2 className="flex items-center gap-2 text-lg font-semibold">
          <Lock className="h-5 w-5 text-zinc-400" />
          External markets — pending legal clearance
        </h2>
        <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
          {features.tradeExternal
            ? "External trading is enabled for this environment."
            : "The next stage of Trade — your logo as a transferable on-chain crypto ID, public NFT markets, and Polymarket-compatible assets — is designed and ready to build out, but transferable value-bearing instruments require securities/CFTC legal clearance first."}
        </p>
        <ul className="mt-3 flex flex-col gap-1.5 text-sm text-zinc-500 dark:text-zinc-400">
          <li>• Your passport marker is already on-chain today (anchor your logo in Logo Studio).</li>
          <li>• Sticker trading and the points wallet above are fully live, in-platform.</li>
          <li>
            • Once licensing/legal sign-off is in hand, external trading switches on via{" "}
            <code className="rounded bg-zinc-200 px-1 py-0.5 text-xs dark:bg-zinc-800">
              NEXT_PUBLIC_ENABLE_TRADE=true
            </code>{" "}
            — no rebuild needed for the in-platform layer it extends.
          </li>
        </ul>
      </div>
    </div>
  );
}
