"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { useTranslation } from "react-i18next";
import { api } from "@/lib/api";
import { getAmbientEngine, type AmbientMode } from "@/lib/ambient-audio";
import { Music, ArrowRight, Upload } from "@/components/icons";

const ACCENT = "var(--tab-music)";

type AudioItem = {
  id: string;
  title: string;
  description?: string;
  contentType: string;
  isPublic: boolean;
  isEncrypted?: boolean;
  createdAt: string;
};

export default function MusicPage() {
  const { user } = useAuth();
  const { t } = useTranslation();

  // Ambient engine state (shared with the header toggle)
  const [playing, setPlaying] = useState(false);
  const [mode, setMode] = useState<AmbientMode>("rain");
  const [volume, setVolume] = useState(0.5);

  // Vault audio tracks
  const [tracks, setTracks] = useState<AudioItem[]>([]);
  const [tracksLoading, setTracksLoading] = useState(true);
  const [nowPlaying, setNowPlaying] = useState<string | null>(null);

  // Signature track: the public Vault track that plays on your public profile
  const [signatureId, setSignatureId] = useState<string | null>(null);
  const [signatureMsg, setSignatureMsg] = useState("");

  async function toggleSignature(track: AudioItem) {
    const clearing = signatureId === track.id;
    setSignatureMsg("");
    const res = await api.users.setSignatureTrack(clearing ? null : track.id);
    if (res.error) {
      setSignatureMsg(res.error);
      return;
    }
    setSignatureId(clearing ? null : track.id);
    setSignatureMsg(
      clearing
        ? "Signature track removed."
        : `"${track.title}" now plays on your public profile ✓`
    );
  }

  useEffect(() => {
    return getAmbientEngine().subscribe((s) => {
      setPlaying(s.playing);
      setMode(s.mode);
      setVolume(s.volume);
    });
  }, []);

  useEffect(() => {
    if (!user) return;
    let cancelled = false;
    (async () => {
      // Vault list caps at 100 per page; one page is plenty for a track picker.
      const res = await api.content.list(100, 0);
      if (cancelled) return;
      setTracksLoading(false);
      if (res.data) {
        setTracks(
          (res.data.items || []).filter((i: AudioItem) =>
            i.contentType?.startsWith("audio")
          )
        );
      }
      if (user.username) {
        const profile = await api.users.publicProfile(user.username);
        if (!cancelled && profile.data?.signatureTrack) {
          setSignatureId(profile.data.signatureTrack.id);
        }
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [user]);

  if (!user) {
    return (
      <div className="mx-auto max-w-2xl px-4 py-24 text-center sm:px-6">
        <div
          className="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl text-white shadow-lg"
          style={{ background: ACCENT }}
        >
          <Music className="h-8 w-8" />
        </div>
        <h1 className="mt-6 text-3xl font-bold tracking-tight">Ambient & Music</h1>
        <p className="mx-auto mt-4 max-w-xl text-zinc-600 dark:text-zinc-400">
          Moody rain and storm soundscapes, plus a player for your own uploaded
          tracks. Sign in to set the mood.
        </p>
        <Link
          href="/sign-in"
          className="mt-7 inline-flex items-center gap-2 rounded-lg px-6 py-3 text-sm font-medium text-white transition-opacity hover:opacity-90"
          style={{ background: ACCENT }}
        >
          Sign in
          <ArrowRight className="h-4 w-4" />
        </Link>
      </div>
    );
  }

  const engine = getAmbientEngine();

  return (
    <div className="mx-auto w-full max-w-4xl px-4 py-12 sm:px-6">
      <div className="flex items-center gap-3">
        <div
          className="flex h-12 w-12 items-center justify-center rounded-xl text-white shadow"
          style={{ background: ACCENT }}
        >
          <Music className="h-6 w-6" />
        </div>
        <div>
          <h1 className="text-2xl font-bold tracking-tight">{t("tabs.musicTitle")}</h1>
          <p className="text-sm text-zinc-500 dark:text-zinc-400">
            {t("tabs.musicSubtitle")}
          </p>
        </div>
      </div>

      {/* ── Ambient soundscape ── */}
      <div className="mt-8 rounded-2xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <h2 className="text-lg font-semibold">Ambient soundscape</h2>
            <p className="text-sm text-zinc-500 dark:text-zinc-400">
              Synthesized live in your browser — moody rain, distant thunder.
            </p>
          </div>
          <button
            onClick={() => engine.toggle()}
            className={`rounded-full px-6 py-2.5 text-sm font-semibold text-white transition-opacity hover:opacity-90`}
            style={{ background: playing ? "#3f3f46" : ACCENT }}
          >
            {playing ? t("tabs.turnOff") : t("tabs.turnOn")}
          </button>
        </div>

        <div className="mt-5 flex flex-wrap items-center gap-6">
          <div className="flex gap-2">
            {(["rain", "storm"] as const).map((m) => (
              <button
                key={m}
                onClick={() => engine.setMode(m)}
                className={`rounded-full px-4 py-1.5 text-sm font-medium capitalize transition-colors ${
                  mode === m
                    ? "text-white"
                    : "bg-zinc-100 text-zinc-600 hover:bg-zinc-200 dark:bg-zinc-800 dark:text-zinc-400"
                }`}
                style={mode === m ? { background: ACCENT } : undefined}
              >
                {m === "rain" ? "🌧 Rain" : "⛈ Storm"}
              </button>
            ))}
          </div>
          <label className="flex items-center gap-3 text-sm text-zinc-500 dark:text-zinc-400">
            Volume
            <input
              type="range"
              min={0}
              max={1}
              step={0.05}
              value={volume}
              onChange={(e) => engine.setVolume(Number(e.target.value))}
              className="w-36 accent-violet-500"
            />
          </label>
        </div>
      </div>

      {/* ── Your tracks ── */}
      <div className="mt-6 rounded-2xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <h2 className="text-lg font-semibold">{t("tabs.yourTracks")}</h2>
            <p className="text-sm text-zinc-500 dark:text-zinc-400">
              Audio files from your Vault — play them here, no streaming
              licenses needed because they&apos;re yours.
            </p>
          </div>
          <Link
            href="/vault/upload"
            className="inline-flex items-center gap-2 rounded-lg border border-zinc-200 px-4 py-2 text-sm font-medium text-zinc-700 transition-colors hover:bg-zinc-50 dark:border-zinc-800 dark:text-zinc-300 dark:hover:bg-zinc-900"
          >
            <Upload className="h-4 w-4" />
            Upload audio
          </Link>
        </div>

        {signatureMsg && (
          <p className="mt-3 text-sm" style={{ color: ACCENT }}>
            {signatureMsg}
          </p>
        )}

        {tracksLoading ? (
          <p className="mt-5 text-sm text-zinc-500">Loading…</p>
        ) : tracks.length === 0 ? (
          <p className="mt-5 text-sm text-zinc-500 dark:text-zinc-400">
            No audio in your Vault yet. Upload an MP3 to play it here.
          </p>
        ) : (
          <ul className="mt-5 flex flex-col gap-4">
            {tracks.map((t) => (
              <li
                key={t.id}
                className="rounded-xl border border-zinc-200 p-4 dark:border-zinc-800"
              >
                <div className="flex items-center justify-between gap-3">
                  <div className="min-w-0">
                    <p className="truncate font-medium">{t.title}</p>
                    {t.description && (
                      <p className="truncate text-sm text-zinc-500 dark:text-zinc-400">
                        {t.description}
                      </p>
                    )}
                  </div>
                  <div className="flex shrink-0 items-center gap-3">
                    {!t.isEncrypted && (
                      <button
                        onClick={() => toggleSignature(t)}
                        className={`text-sm font-medium ${
                          signatureId === t.id
                            ? ""
                            : "text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100"
                        }`}
                        style={signatureId === t.id ? { color: ACCENT } : undefined}
                        title={
                          t.isPublic
                            ? "Plays on your public profile"
                            : "Track must be public to be your signature"
                        }
                      >
                        {signatureId === t.id ? "★ Signature" : "☆ Set as signature"}
                      </button>
                    )}
                    <Link
                      href={`/vault/item?id=${t.id}`}
                      className="text-sm font-medium"
                      style={{ color: ACCENT }}
                    >
                      View in Vault
                    </Link>
                  </div>
                </div>
                {nowPlaying === t.id ? (
                  <audio
                    controls
                    autoPlay
                    className="mt-3 w-full"
                    src={api.content.download(t.id)}
                  />
                ) : (
                  <button
                    onClick={() => setNowPlaying(t.id)}
                    className="mt-3 inline-flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium text-white transition-opacity hover:opacity-90"
                    style={{ background: ACCENT }}
                  >
                    ▶ Play
                  </button>
                )}
              </li>
            ))}
          </ul>
        )}
      </div>

      <p className="mt-8 rounded-xl border border-dashed border-zinc-300 bg-zinc-50 p-4 text-center text-xs text-zinc-500 dark:border-zinc-700 dark:bg-zinc-900/50 dark:text-zinc-400">
        Your ★ signature track plays for visitors on your public passport page (it
        must be a public track). Third-party streaming (Spotify-style catalogs) needs
        licensing and stays out of v1.
      </p>
    </div>
  );
}
