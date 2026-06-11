"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { useTranslation } from "react-i18next";
import { api } from "@/lib/api";
import { Trophy, ArrowRight, CheckCircle, Gift } from "@/components/icons";

const ACCENT = "var(--tab-vr)";

type Sticker = { id: string; name: string; emoji: string; rarity: string; description: string };
type OwnedSticker = Sticker & { count: number };
type Achievement = { key: string; name: string; emoji: string };
type LeaderRow = {
  userId: string;
  username: string | null;
  name: string | null;
  image: string | null;
  points: number;
  stickers: number;
  isVerified: boolean;
};

function rarityClasses(rarity: string, owned: boolean): string {
  if (!owned) return "border-dashed border-zinc-300 opacity-40 dark:border-zinc-700";
  switch (rarity) {
    case "legendary":
      return "border-amber-300 bg-amber-50 dark:border-amber-700 dark:bg-amber-900/30";
    case "rare":
      return "border-purple-300 bg-purple-50 dark:border-purple-700 dark:bg-purple-900/30";
    case "uncommon":
      return "border-blue-300 bg-blue-50 dark:border-blue-700 dark:bg-blue-900/30";
    default:
      return "border-zinc-200 bg-zinc-50 dark:border-zinc-700 dark:bg-zinc-800/50";
  }
}

function CooldownTimer({ until, onDone }: { until: string; onDone: () => void }) {
  const [left, setLeft] = useState(() => new Date(until).getTime() - Date.now());
  useEffect(() => {
    const t = setInterval(() => {
      const ms = new Date(until).getTime() - Date.now();
      setLeft(ms);
      if (ms <= 0) {
        clearInterval(t);
        onDone();
      }
    }, 1000);
    return () => clearInterval(t);
  }, [until, onDone]);
  if (left <= 0) return null;
  const h = Math.floor(left / 3_600_000);
  const m = Math.floor((left % 3_600_000) / 60_000);
  const s = Math.floor((left % 60_000) / 1000);
  return (
    <span className="font-mono">
      {h > 0 ? `${h}h ` : ""}
      {m}m {s}s
    </span>
  );
}

export default function VRCommunityPage() {
  const { user } = useAuth();
  const { t } = useTranslation();

  const [points, setPoints] = useState(0);
  const [owned, setOwned] = useState<OwnedSticker[]>([]);
  const [catalog, setCatalog] = useState<Sticker[]>([]);
  const [achievementsList, setAchievementsList] = useState<Achievement[]>([]);
  const [nextExplore, setNextExplore] = useState<string | null>(null);
  const [leaders, setLeaders] = useState<LeaderRow[]>([]);
  const [exploring, setExploring] = useState(false);
  const [drop, setDrop] = useState<{ sticker: Sticker; points: number } | null>(null);
  const [error, setError] = useState("");

  // Gift dialog
  const [giftSticker, setGiftSticker] = useState<OwnedSticker | null>(null);
  const [giftTo, setGiftTo] = useState("");
  const [gifting, setGifting] = useState(false);
  const [giftMsg, setGiftMsg] = useState("");

  const loadMe = useCallback(async () => {
    const res = await api.arena.me();
    if (res.data) {
      setPoints(res.data.points);
      setOwned(res.data.stickers as OwnedSticker[]);
      setCatalog(res.data.catalog);
      setAchievementsList(res.data.achievements);
      setNextExplore(res.data.nextExplore);
    }
  }, []);

  useEffect(() => {
    api.arena.leaderboard().then((res) => {
      if (res.data) setLeaders(res.data.leaderboard);
    });
    if (user) loadMe();
  }, [user, loadMe]);

  async function explore() {
    if (exploring) return;
    setExploring(true);
    setError("");
    setDrop(null);
    const res = await api.arena.explore();
    setExploring(false);
    if (res.error || !res.data) {
      setError(res.error || "Explore failed");
      return;
    }
    setDrop({ sticker: res.data.sticker, points: res.data.points });
    setNextExplore(res.data.nextExplore);
    loadMe();
  }

  async function sendGift() {
    if (!giftSticker || !giftTo.trim() || gifting) return;
    setGifting(true);
    setGiftMsg("");
    const res = await api.arena.gift(giftTo.trim().replace(/^@/, ""), giftSticker.id);
    setGifting(false);
    if (res.error) {
      setGiftMsg(res.error);
      return;
    }
    setGiftMsg(`Sent ${giftSticker.emoji} ${giftSticker.name} to @${giftTo.trim().replace(/^@/, "")} ✓`);
    setGiftSticker(null);
    setGiftTo("");
    loadMe();
  }

  const ownedMap = new Map(owned.map((s) => [s.id, s]));

  return (
    <div className="mx-auto w-full max-w-5xl px-4 py-12 sm:px-6">
      <div className="flex flex-wrap items-center gap-3">
        <div
          className="flex h-12 w-12 items-center justify-center rounded-xl text-white shadow"
          style={{ background: ACCENT }}
        >
          <Trophy className="h-6 w-6" />
        </div>
        <div>
          <h1 className="text-2xl font-bold tracking-tight">{t("tabs.arenaTitle")}</h1>
          <p className="text-sm text-zinc-500 dark:text-zinc-400">
            {t("tabs.arenaSubtitle")}
          </p>
        </div>
        {user && (
          <div className="ml-auto flex items-center gap-3">
            <Link
              href="/trade"
              className="text-sm font-medium text-zinc-500 underline-offset-2 hover:underline dark:text-zinc-400"
            >
              Wallet & history →
            </Link>
            <div className="rounded-full px-4 py-2 text-sm font-bold text-white" style={{ background: ACCENT }}>
              {points} pts
            </div>
          </div>
        )}
      </div>

      {!user ? (
        <div className="mt-12 rounded-2xl border border-zinc-200 bg-white p-10 text-center dark:border-zinc-800 dark:bg-zinc-950">
          <p className="mx-auto max-w-md text-zinc-600 dark:text-zinc-400">
            Sign in to explore the farm, earn stickers for your passport, and gift
            them to other creators.
          </p>
          <Link
            href="/sign-in"
            className="mt-5 inline-flex items-center gap-2 rounded-lg px-6 py-3 text-sm font-medium text-white transition-opacity hover:opacity-90"
            style={{ background: ACCENT }}
          >
            Join the arena
            <ArrowRight className="h-4 w-4" />
          </Link>
        </div>
      ) : (
        <>
          {/* ── The farm (explore action) ── */}
          <div className="mt-8 rounded-2xl border border-zinc-200 bg-white p-6 text-center dark:border-zinc-800 dark:bg-zinc-950">
            <p className="text-5xl">🌾🐓🐑🐖🌾</p>
            <h2 className="mt-3 text-lg font-semibold">The Farm</h2>
            <p className="mx-auto mt-1 max-w-md text-sm text-zinc-500 dark:text-zinc-400">
              Explore every few hours to find a random sticker and earn points.
              Rarer finds, bigger points.
            </p>
            {nextExplore && new Date(nextExplore) > new Date() ? (
              <p className="mt-4 text-sm text-zinc-500 dark:text-zinc-400">
                The farm is regrowing —{" "}
                <CooldownTimer until={nextExplore} onDone={() => setNextExplore(null)} /> until
                your next explore.
              </p>
            ) : (
              <button
                onClick={explore}
                disabled={exploring}
                className="mt-4 rounded-full px-8 py-3 text-sm font-semibold text-white transition-opacity hover:opacity-90 disabled:opacity-50"
                style={{ background: ACCENT }}
              >
                {exploring ? t("tabs.exploring") : `🔍 ${t("tabs.explore")}`}
              </button>
            )}
            {drop && (
              <div className="mx-auto mt-5 max-w-sm rounded-xl border-2 p-4"
                style={{ borderColor: ACCENT }}>
                <p className="text-4xl">{drop.sticker.emoji}</p>
                <p className="mt-1 font-semibold">
                  You found {drop.sticker.name}!{" "}
                  <span className="text-xs uppercase tracking-wide text-zinc-500">
                    {drop.sticker.rarity}
                  </span>
                </p>
                <p className="text-sm text-zinc-500 dark:text-zinc-400">{drop.sticker.description}</p>
                <p className="mt-1 text-sm font-bold" style={{ color: ACCENT }}>
                  +{drop.points} points
                </p>
              </div>
            )}
            {error && <p className="mt-3 text-sm text-red-600 dark:text-red-400">{error}</p>}
          </div>

          {/* ── Sticker collection ── */}
          <div className="mt-6 rounded-2xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold">
                {t("tabs.stickers")}{" "}
                <span className="text-sm font-normal text-zinc-500">
                  ({owned.length}/{catalog.length})
                </span>
              </h2>
            </div>
            <div className="mt-4 grid grid-cols-3 gap-3 sm:grid-cols-4 md:grid-cols-6">
              {catalog.map((s) => {
                const mine = ownedMap.get(s.id);
                return (
                  <div
                    key={s.id}
                    className={`relative flex flex-col items-center rounded-xl border p-3 text-center ${rarityClasses(s.rarity, !!mine)}`}
                    title={mine ? s.description : `Not found yet (${s.rarity})`}
                  >
                    <span className="text-3xl">{mine ? s.emoji : "❔"}</span>
                    <span className="mt-1 text-xs font-medium">{s.name}</span>
                    <span className="text-[10px] uppercase tracking-wide text-zinc-400">
                      {s.rarity}
                    </span>
                    {mine && mine.count > 1 && (
                      <span className="absolute -right-1.5 -top-1.5 rounded-full bg-zinc-900 px-1.5 text-[10px] font-bold text-white dark:bg-zinc-100 dark:text-zinc-900">
                        ×{mine.count}
                      </span>
                    )}
                    {mine && (
                      <button
                        onClick={() => {
                          setGiftSticker(mine);
                          setGiftMsg("");
                        }}
                        className="mt-1.5 inline-flex items-center gap-1 text-[11px] font-medium text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-100"
                      >
                        <Gift className="h-3 w-3" />
                        Gift
                      </button>
                    )}
                  </div>
                );
              })}
            </div>

            {/* Gift dialog (inline) */}
            {giftSticker && (
              <div className="mt-4 rounded-xl border border-zinc-200 p-4 dark:border-zinc-700">
                <p className="text-sm font-medium">
                  Gift {giftSticker.emoji} {giftSticker.name} to:
                </p>
                <div className="mt-2 flex flex-wrap gap-2">
                  <input
                    value={giftTo}
                    onChange={(e) => setGiftTo(e.target.value)}
                    placeholder="@username"
                    className="flex-1 rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-900"
                  />
                  <button
                    onClick={sendGift}
                    disabled={gifting || !giftTo.trim()}
                    className="rounded-lg px-4 py-2 text-sm font-medium text-white transition-opacity hover:opacity-90 disabled:opacity-50"
                    style={{ background: ACCENT }}
                  >
                    {gifting ? "Sending…" : "Send"}
                  </button>
                  <button
                    onClick={() => setGiftSticker(null)}
                    className="rounded-lg border border-zinc-200 px-4 py-2 text-sm dark:border-zinc-700"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            )}
            {giftMsg && (
              <p className="mt-3 text-sm" style={{ color: ACCENT }}>
                {giftMsg}
              </p>
            )}
          </div>

          {/* ── Achievements ── */}
          {achievementsList.length > 0 && (
            <div className="mt-6 rounded-2xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
              <h2 className="text-lg font-semibold">{t("tabs.achievements")}</h2>
              <div className="mt-3 flex flex-wrap gap-2">
                {achievementsList.map((a) => (
                  <span
                    key={a.key}
                    className="inline-flex items-center gap-1.5 rounded-full bg-zinc-100 px-3 py-1.5 text-sm font-medium dark:bg-zinc-800"
                  >
                    {a.emoji} {a.name}
                  </span>
                ))}
              </div>
            </div>
          )}
        </>
      )}

      {/* ── Leaderboard (public) ── */}
      <div className="mt-6 rounded-2xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
        <h2 className="text-lg font-semibold">{t("tabs.leaderboard")}</h2>
        {leaders.length === 0 ? (
          <p className="mt-3 text-sm text-zinc-500 dark:text-zinc-400">
            No explorers yet — be the first on the board.
          </p>
        ) : (
          <ol className="mt-4 flex flex-col gap-2">
            {leaders.map((l, i) => (
              <li
                key={l.userId}
                className="flex items-center gap-3 rounded-xl border border-zinc-100 px-4 py-2.5 dark:border-zinc-800"
              >
                <span className="w-6 text-sm font-bold text-zinc-400">
                  {i === 0 ? "🥇" : i === 1 ? "🥈" : i === 2 ? "🥉" : i + 1}
                </span>
                {l.image ? (
                  <img src={l.image} alt="" className="h-8 w-8 rounded-full object-cover" />
                ) : (
                  <div className="flex h-8 w-8 items-center justify-center rounded-full bg-zinc-100 text-xs font-bold text-zinc-500 dark:bg-zinc-800">
                    {(l.name || l.username || "?")[0].toUpperCase()}
                  </div>
                )}
                <Link
                  href={`/profile?u=${l.username}`}
                  className="flex min-w-0 items-center gap-1.5 font-medium hover:underline"
                >
                  <span className="truncate">{l.name || `@${l.username}`}</span>
                  {l.isVerified && <CheckCircle className="h-4 w-4 shrink-0 text-blue-500" />}
                </Link>
                <span className="ml-auto text-sm text-zinc-500">{l.stickers} stickers</span>
                <span className="w-20 text-right text-sm font-bold" style={{ color: ACCENT }}>
                  {l.points} pts
                </span>
              </li>
            ))}
          </ol>
        )}
      </div>

      <p className="mt-8 rounded-xl border border-dashed border-zinc-300 bg-zinc-50 p-4 text-center text-xs text-zinc-500 dark:border-zinc-700 dark:bg-zinc-900/50 dark:text-zinc-400">
        Points and stickers carry in-platform status only — they are not cashable or
        externally tradeable. Engagement on your public work (comments) also feeds
        your points. A richer 3D arena builds on these same mechanics later.
      </p>
    </div>
  );
}
