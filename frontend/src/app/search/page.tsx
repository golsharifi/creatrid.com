"use client";

import { useState, useEffect, Suspense } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { api } from "@/lib/api";
import Link from "next/link";
import { Search as SearchIcon, Image, User, ChevronLeft, ChevronRight } from "lucide-react";
import { useTranslation } from "react-i18next";

type ContentResult = {
  id: string;
  title: string;
  description: string | null;
  contentType: string;
  thumbnailUrl: string | null;
  creatorName: string | null;
  creatorUsername: string | null;
};

type CreatorResult = {
  id: string;
  name: string | null;
  username: string | null;
  bio: string | null;
  image: string | null;
  creatorScore: number | null;
  isVerified: boolean;
};

const LIMIT = 20;

function SearchContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const { t } = useTranslation();

  const [query, setQuery] = useState(searchParams.get("q") ?? "");
  const [tab, setTab] = useState<"all" | "content" | "creators">(
    (searchParams.get("type") as any) ?? "all"
  );
  const [content, setContent] = useState<ContentResult[]>([]);
  const [contentTotal, setContentTotal] = useState(0);
  const [creators, setCreators] = useState<CreatorResult[]>([]);
  const [creatorsTotal, setCreatorsTotal] = useState(0);
  const [page, setPage] = useState(0);
  const [loading, setLoading] = useState(false);
  const [searched, setSearched] = useState(false);

  const doSearch = async (q: string, type: string, p: number) => {
    if (!q.trim()) return;
    setLoading(true);
    setSearched(true);
    const result = await api.search.query(q.trim(), type === "all" ? undefined : type, LIMIT, p * LIMIT);
    if (result.data) {
      const data = result.data as any;
      setContent(data.content ?? []);
      setContentTotal(data.contentTotal ?? 0);
      setCreators(data.creators ?? []);
      setCreatorsTotal(data.creatorsTotal ?? 0);
    }
    setLoading(false);
  };

  useEffect(() => {
    const q = searchParams.get("q");
    if (q) {
      setQuery(q);
      doSearch(q, tab, 0);
    }
  }, []);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(0);
    doSearch(query, tab, 0);
    router.replace(`/search?q=${encodeURIComponent(query)}&type=${tab}`);
  };

  const handleTabChange = (newTab: "all" | "content" | "creators") => {
    setTab(newTab);
    setPage(0);
    doSearch(query, newTab, 0);
  };

  const total = tab === "content" ? contentTotal : tab === "creators" ? creatorsTotal : contentTotal + creatorsTotal;
  const totalPages = Math.ceil(total / LIMIT);

  return (
    <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold tracking-tight text-zinc-900 dark:text-zinc-100">{t("search.title")}</h1>
        <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">{t("search.subtitle")}</p>
      </div>

      {/* Search form */}
      <form onSubmit={handleSearch} className="mb-6 flex gap-2">
        <div className="relative flex-1">
          <SearchIcon className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-zinc-400" />
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="w-full rounded-lg border border-zinc-300 bg-white py-2 pl-10 pr-4 text-sm dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-100"
            placeholder={t("search.placeholder")}
          />
        </div>
        <button type="submit" className="rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white hover:bg-zinc-800 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200">
          {t("search.searchButton")}
        </button>
      </form>

      {/* Tabs */}
      <div className="mb-6 flex gap-1 rounded-lg bg-zinc-100 p-1 dark:bg-zinc-900">
        {(["all", "content", "creators"] as const).map((t2) => (
          <button
            key={t2}
            onClick={() => handleTabChange(t2)}
            className={`flex-1 rounded-md px-3 py-1.5 text-sm font-medium transition-colors ${tab === t2 ? "bg-white text-zinc-900 shadow-sm dark:bg-zinc-800 dark:text-zinc-100" : "text-zinc-500 hover:text-zinc-700 dark:text-zinc-400 dark:hover:text-zinc-300"}`}
          >
            {t(`search.tab${t2.charAt(0).toUpperCase() + t2.slice(1)}`)}
          </button>
        ))}
      </div>

      {/* Results */}
      {loading ? (
        <div className="flex items-center justify-center py-16">
          <div className="h-6 w-6 animate-spin rounded-full border-2 border-zinc-300 border-t-zinc-900 dark:border-zinc-600 dark:border-t-zinc-100" />
        </div>
      ) : !searched ? null : (
        <div className="space-y-6">
          {/* Content results */}
          {(tab === "all" || tab === "content") && content.length > 0 && (
            <div>
              {tab === "all" && <h2 className="mb-3 text-sm font-semibold text-zinc-500 dark:text-zinc-400">{t("search.contentResults")} ({contentTotal})</h2>}
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                {content.map((item) => (
                  <Link key={item.id} href={`/marketplace/item?id=${item.id}`} className="flex gap-3 rounded-xl border border-zinc-200 bg-white p-4 transition-colors hover:border-zinc-300 dark:border-zinc-800 dark:bg-zinc-950 dark:hover:border-zinc-700">
                    {item.thumbnailUrl ? (
                      <img src={item.thumbnailUrl} alt={item.title} className="h-16 w-16 flex-shrink-0 rounded-lg object-cover" />
                    ) : (
                      <div className="flex h-16 w-16 flex-shrink-0 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-900">
                        <Image className="h-5 w-5 text-zinc-400" />
                      </div>
                    )}
                    <div className="min-w-0">
                      <p className="font-medium text-zinc-900 dark:text-zinc-100">{item.title}</p>
                      {item.description && <p className="mt-0.5 line-clamp-1 text-xs text-zinc-500 dark:text-zinc-400">{item.description}</p>}
                      {item.creatorUsername && <p className="mt-1 text-xs text-zinc-400">@{item.creatorUsername}</p>}
                    </div>
                  </Link>
                ))}
              </div>
            </div>
          )}

          {/* Creator results */}
          {(tab === "all" || tab === "creators") && creators.length > 0 && (
            <div>
              {tab === "all" && <h2 className="mb-3 text-sm font-semibold text-zinc-500 dark:text-zinc-400">{t("search.creatorResults")} ({creatorsTotal})</h2>}
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                {creators.map((creator) => (
                  <Link key={creator.id} href={`/profile?u=${creator.username}`} className="flex items-center gap-3 rounded-xl border border-zinc-200 bg-white p-4 transition-colors hover:border-zinc-300 dark:border-zinc-800 dark:bg-zinc-950 dark:hover:border-zinc-700">
                    {creator.image ? (
                      <img src={creator.image} alt={creator.name ?? ""} className="h-12 w-12 flex-shrink-0 rounded-full object-cover" />
                    ) : (
                      <div className="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full bg-zinc-100 dark:bg-zinc-900">
                        <User className="h-5 w-5 text-zinc-400" />
                      </div>
                    )}
                    <div className="min-w-0">
                      <p className="font-medium text-zinc-900 dark:text-zinc-100">{creator.name ?? creator.username}</p>
                      <p className="text-xs text-zinc-500 dark:text-zinc-400">@{creator.username}{creator.creatorScore != null ? ` · Score: ${creator.creatorScore}` : ""}</p>
                    </div>
                  </Link>
                ))}
              </div>
            </div>
          )}

          {content.length === 0 && creators.length === 0 && (
            <div className="rounded-xl border border-zinc-200 bg-zinc-50 py-16 text-center dark:border-zinc-800 dark:bg-zinc-900">
              <SearchIcon className="mx-auto h-10 w-10 text-zinc-300 dark:text-zinc-600" />
              <p className="mt-4 text-sm text-zinc-500 dark:text-zinc-400">{t("search.noResults")}</p>
            </div>
          )}
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="mt-8 flex items-center justify-between">
          <button onClick={() => { setPage((p) => Math.max(0, p - 1)); doSearch(query, tab, Math.max(0, page - 1)); }} disabled={page === 0} className="flex items-center gap-1 rounded-lg border border-zinc-300 px-3 py-2 text-sm font-medium text-zinc-700 hover:bg-zinc-50 disabled:opacity-40 dark:border-zinc-700 dark:text-zinc-300">
            <ChevronLeft className="h-4 w-4" />
            {t("common.previous")}
          </button>
          <span className="text-sm text-zinc-500 dark:text-zinc-400">{t("common.page", { current: page + 1, total: totalPages })}</span>
          <button onClick={() => { setPage((p) => Math.min(totalPages - 1, p + 1)); doSearch(query, tab, Math.min(totalPages - 1, page + 1)); }} disabled={page >= totalPages - 1} className="flex items-center gap-1 rounded-lg border border-zinc-300 px-3 py-2 text-sm font-medium text-zinc-700 hover:bg-zinc-50 disabled:opacity-40 dark:border-zinc-700 dark:text-zinc-300">
            {t("common.next")}
            <ChevronRight className="h-4 w-4" />
          </button>
        </div>
      )}
    </div>
  );
}

export default function SearchPage() {
  return (
    <Suspense fallback={null}>
      <SearchContent />
    </Suspense>
  );
}
