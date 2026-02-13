"use client";

import { useEffect, useState, useRef, useCallback } from "react";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth-context";
import Link from "next/link";
import { Search, CheckCircle, Link2, Shield, ChevronLeft, ChevronRight } from "@/components/icons";
import { useTranslation } from "react-i18next";

type Creator = {
  id: string;
  name: string | null;
  username: string | null;
  image: string | null;
  bio: string | null;
  creatorScore: number | null;
  isVerified: boolean;
  connections: number;
};

const PLATFORMS = [
  { value: "", labelKey: "discover.allPlatforms" },
  { value: "youtube", label: "YouTube" },
  { value: "github", label: "GitHub" },
  { value: "twitter", label: "Twitter/X" },
  { value: "linkedin", label: "LinkedIn" },
  { value: "instagram", label: "Instagram" },
  { value: "dribbble", label: "Dribbble" },
  { value: "behance", label: "Behance" },
];

export default function DiscoverPage() {
  const { user } = useAuth();
  const { t } = useTranslation();
  const [creators, setCreators] = useState<Creator[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(0);
  const [platform, setPlatform] = useState("");
  const [minScore, setMinScore] = useState(0);
  const [searchQuery, setSearchQuery] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [sendingTo, setSendingTo] = useState<string | null>(null);
  const [sentTo, setSentTo] = useState<Set<string>>(new Set());
  const debounceRef = useRef<ReturnType<typeof setTimeout>>(undefined);

  const handleSearch = useCallback((value: string) => {
    setSearchQuery(value);
    clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => {
      setDebouncedSearch(value);
      setPage(0);
    }, 300);
  }, []);

  useEffect(() => {
    setLoading(true);
    api.discover
      .list({ limit: 20, offset: page * 20, minScore, platform, q: debouncedSearch || undefined })
      .then((r) => {
        if (r.data) {
          setCreators(r.data.creators);
          setTotal(r.data.total);
        }
        setLoading(false);
      });
  }, [page, platform, minScore, debouncedSearch]);

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [page]);

  async function sendCollabRequest(creatorId: string) {
    if (!user) return;
    setSendingTo(creatorId);
    const result = await api.collaborations.send(creatorId, "");
    if (result.data) {
      setSentTo((prev) => new Set(prev).add(creatorId));
    }
    setSendingTo(null);
  }

  return (
    <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold">{t("discover.title")}</h1>
        <p className="mt-1 text-zinc-500">
          {t("discover.subtitle")}
        </p>
      </div>

      {/* Filters */}
      <div className="mb-6 flex flex-wrap items-center gap-3">
        <div className="flex items-center gap-2 rounded-lg border border-zinc-200 px-3 py-2 dark:border-zinc-700">
          <Search className="h-4 w-4 text-zinc-400" />
          <input
            type="text"
            placeholder={t("discover.searchPlaceholder")}
            value={searchQuery}
            onChange={(e) => handleSearch(e.target.value)}
            className="bg-transparent text-sm placeholder:text-zinc-400 focus:outline-none"
          />
        </div>
        <div className="flex items-center gap-2 rounded-lg border border-zinc-200 px-3 py-2 dark:border-zinc-700">
          <select
            value={platform}
            onChange={(e) => { setPlatform(e.target.value); setPage(0); }}
            className="bg-transparent text-sm focus:outline-none dark:bg-zinc-800 dark:text-zinc-100"
          >
            {PLATFORMS.map((p) => (
              <option key={p.value} value={p.value} className="dark:bg-zinc-800">
                {p.labelKey ? t(p.labelKey) : p.label}
              </option>
            ))}
          </select>
        </div>
        <div className="flex items-center gap-2 rounded-lg border border-zinc-200 px-3 py-2 dark:border-zinc-700">
          <Shield className="h-4 w-4 text-zinc-400" />
          <select
            value={minScore}
            onChange={(e) => { setMinScore(Number(e.target.value)); setPage(0); }}
            className="bg-transparent text-sm focus:outline-none dark:bg-zinc-800 dark:text-zinc-100"
          >
            <option value={0} className="dark:bg-zinc-800">{t("discover.anyScore")}</option>
            <option value={20} className="dark:bg-zinc-800">20+</option>
            <option value={40} className="dark:bg-zinc-800">40+</option>
            <option value={60} className="dark:bg-zinc-800">60+</option>
            <option value={80} className="dark:bg-zinc-800">80+</option>
          </select>
        </div>
        <span className="text-sm text-zinc-500">{t("common.creatorsFound", { count: total })}</span>
      </div>

      {/* Creator Grid */}
      {loading ? (
        <p className="py-12 text-center text-zinc-400">{t("discover.loadingCreators")}</p>
      ) : creators.length === 0 ? (
        <p className="py-12 text-center text-zinc-400">{t("discover.noCreatorsFound")}</p>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {creators.map((c) => (
            <div
              key={c.id}
              className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800"
            >
              <div className="flex items-start gap-3">
                {c.image ? (
                  <img src={c.image} alt="" className="h-12 w-12 rounded-full object-cover" />
                ) : (
                  <div className="flex h-12 w-12 items-center justify-center rounded-full bg-zinc-100 text-lg font-bold text-zinc-400 dark:bg-zinc-800">
                    {(c.name || "?")[0].toUpperCase()}
                  </div>
                )}
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-1.5">
                    <Link
                      href={`/profile?u=${c.username}`}
                      className="truncate font-semibold hover:underline"
                    >
                      {c.name || t("common.creator")}
                    </Link>
                    {c.isVerified && <CheckCircle className="h-4 w-4 shrink-0 text-blue-500" />}
                  </div>
                  {c.username && (
                    <p className="text-sm text-zinc-500">@{c.username}</p>
                  )}
                </div>
              </div>

              {c.bio && (
                <p className="mt-3 line-clamp-2 text-sm text-zinc-600 dark:text-zinc-400">
                  {c.bio}
                </p>
              )}

              <div className="mt-3 flex items-center gap-4 text-xs text-zinc-500">
                {c.creatorScore !== null && (
                  <span className="flex items-center gap-1">
                    <Shield className="h-3 w-3" />
                    {t("discover.score", { score: c.creatorScore })}
                  </span>
                )}
                <span className="flex items-center gap-1">
                  <Link2 className="h-3 w-3" />
                  {t("discover.connectionsCount", { count: c.connections })}
                </span>
              </div>

              {user && user.id !== c.id && (
                <div className="mt-4">
                  {sentTo.has(c.id) ? (
                    <span className="text-xs text-zinc-400">{t("discover.requestSent")}</span>
                  ) : (
                    <button
                      onClick={() => sendCollabRequest(c.id)}
                      disabled={sendingTo === c.id}
                      className="rounded-lg bg-zinc-900 px-3 py-1.5 text-xs font-medium text-white transition-colors hover:bg-zinc-700 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-300"
                    >
                      {sendingTo === c.id ? t("discover.sending") : t("discover.collaborate")}
                    </button>
                  )}
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Pagination */}
      {total > 20 && (
        <div className="mt-6 flex items-center justify-center gap-4">
          <button
            onClick={() => setPage((p) => Math.max(0, p - 1))}
            disabled={page === 0}
            className="flex items-center gap-1 rounded-lg border border-zinc-200 px-3 py-1.5 text-sm font-medium disabled:opacity-50 dark:border-zinc-700"
          >
            <ChevronLeft className="h-4 w-4" /> {t("common.previous")}
          </button>
          <span className="text-sm text-zinc-500">
            {t("common.page", { current: page + 1, total: Math.ceil(total / 20) })}
          </span>
          <button
            onClick={() => setPage((p) => p + 1)}
            disabled={(page + 1) * 20 >= total}
            className="flex items-center gap-1 rounded-lg border border-zinc-200 px-3 py-1.5 text-sm font-medium disabled:opacity-50 dark:border-zinc-700"
          >
            {t("common.next")} <ChevronRight className="h-4 w-4" />
          </button>
        </div>
      )}
    </div>
  );
}
