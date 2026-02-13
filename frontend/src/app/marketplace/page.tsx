"use client";

import { useEffect, useState, useRef, useCallback } from "react";
import { api } from "@/lib/api";
import Link from "next/link";
import {
  Search,
  Image,
  Video,
  Music,
  FileText,
  File,
  ChevronLeft,
  ChevronRight,
} from "@/components/icons";
import { useTranslation } from "react-i18next";

type MarketplaceItem = {
  id: string;
  title: string;
  contentType: string;
  creatorName: string | null;
  creatorUsername: string | null;
  lowestPriceCents: number | null;
  tags: string[];
  thumbnailUrl?: string;
};

const FILTERS = [
  { key: "", labelKey: "marketplace.filterAll" },
  { key: "image", labelKey: "marketplace.filterImages", icon: Image },
  { key: "video", labelKey: "marketplace.filterVideo", icon: Video },
  { key: "audio", labelKey: "marketplace.filterAudio", icon: Music },
  { key: "text", labelKey: "marketplace.filterText", icon: FileText },
  { key: "other", labelKey: "marketplace.filterOther", icon: File },
];

const SORT_OPTIONS = [
  { value: "newest", labelKey: "marketplace.sortNewest" },
  { value: "price_asc", labelKey: "marketplace.sortPriceLow" },
  { value: "price_desc", labelKey: "marketplace.sortPriceHigh" },
];

function getContentCategory(contentType: string): string {
  if (contentType.startsWith("image/")) return "image";
  if (contentType.startsWith("video/")) return "video";
  if (contentType.startsWith("audio/")) return "audio";
  if (
    contentType.startsWith("text/") ||
    contentType === "application/pdf" ||
    contentType === "application/msword" ||
    contentType.includes("document")
  )
    return "text";
  return "other";
}

function ContentTypeIcon({ contentType, className }: { contentType: string; className?: string }) {
  const category = getContentCategory(contentType);
  switch (category) {
    case "image":
      return <Image className={className} />;
    case "video":
      return <Video className={className} />;
    case "audio":
      return <Music className={className} />;
    case "text":
      return <FileText className={className} />;
    default:
      return <File className={className} />;
  }
}

const LIMIT = 20;

export default function MarketplacePage() {
  const { t } = useTranslation();
  const [items, setItems] = useState<MarketplaceItem[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(0);
  const [typeFilter, setTypeFilter] = useState("");
  const [sort, setSort] = useState("newest");
  const [searchQuery, setSearchQuery] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [loading, setLoading] = useState(true);
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
    api.marketplace
      .browse({
        type: typeFilter || undefined,
        q: debouncedSearch || undefined,
        sort,
        limit: LIMIT,
        offset: page * LIMIT,
      })
      .then((r) => {
        if (r.data) {
          setItems(r.data.items || []);
          setTotal(r.data.total || 0);
        }
        setLoading(false);
      });
  }, [page, typeFilter, sort, debouncedSearch]);

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [page]);

  return (
    <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">
          {t("marketplace.title")}
        </h1>
        <p className="mt-1 text-zinc-500">
          {t("marketplace.subtitle")}
        </p>
      </div>

      {/* Search and Filters */}
      <div className="mb-6 flex flex-wrap items-center gap-3">
        <div className="flex items-center gap-2 rounded-lg border border-zinc-200 px-3 py-2 dark:border-zinc-700">
          <Search className="h-4 w-4 text-zinc-400" />
          <input
            type="text"
            placeholder={t("marketplace.searchPlaceholder")}
            value={searchQuery}
            onChange={(e) => handleSearch(e.target.value)}
            className="bg-transparent text-sm placeholder:text-zinc-400 focus:outline-none dark:text-zinc-100"
          />
        </div>

        {/* Sort */}
        <div className="flex items-center gap-2 rounded-lg border border-zinc-200 px-3 py-2 dark:border-zinc-700">
          <select
            value={sort}
            onChange={(e) => {
              setSort(e.target.value);
              setPage(0);
            }}
            className="bg-transparent text-sm focus:outline-none dark:bg-zinc-800 dark:text-zinc-100"
          >
            {SORT_OPTIONS.map((option) => (
              <option key={option.value} value={option.value} className="dark:bg-zinc-800">
                {t(option.labelKey)}
              </option>
            ))}
          </select>
        </div>

        <span className="text-sm text-zinc-500">
          {t("marketplace.resultsFound", { count: total })}
        </span>
      </div>

      {/* Type Filters */}
      <div className="mb-6 flex flex-wrap gap-2">
        {FILTERS.map((f) => (
          <button
            key={f.key}
            onClick={() => {
              setTypeFilter(f.key);
              setPage(0);
            }}
            className={`flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm font-medium transition-colors ${
              typeFilter === f.key
                ? "bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900"
                : "border border-zinc-200 text-zinc-600 hover:bg-zinc-50 dark:border-zinc-700 dark:text-zinc-400 dark:hover:bg-zinc-800"
            }`}
          >
            {f.icon && <f.icon className="h-3.5 w-3.5" />}
            {t(f.labelKey)}
          </button>
        ))}
      </div>

      {/* Content Grid */}
      {loading ? (
        <p className="py-12 text-center text-zinc-400">{t("marketplace.loading")}</p>
      ) : items.length === 0 ? (
        <div className="py-16 text-center">
          <File className="mx-auto h-12 w-12 text-zinc-300 dark:text-zinc-600" />
          <p className="mt-4 text-zinc-500">{t("marketplace.noResults")}</p>
        </div>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {items.map((item) => (
            <Link
              key={item.id}
              href={`/marketplace/item?id=${item.id}`}
              className="group rounded-xl border border-zinc-200 bg-white p-4 transition-colors hover:border-zinc-300 dark:border-zinc-800 dark:bg-zinc-900 dark:hover:border-zinc-700"
            >
              {/* Thumbnail / Icon */}
              <div className="mb-3 flex h-32 items-center justify-center rounded-lg bg-zinc-50 dark:bg-zinc-800">
                {item.thumbnailUrl && getContentCategory(item.contentType) === "image" ? (
                  <img
                    src={item.thumbnailUrl}
                    alt={item.title}
                    className="h-full w-full rounded-lg object-cover"
                  />
                ) : (
                  <ContentTypeIcon
                    contentType={item.contentType}
                    className="h-10 w-10 text-zinc-400"
                  />
                )}
              </div>

              {/* Info */}
              <h3 className="truncate text-sm font-semibold text-zinc-900 dark:text-zinc-100">
                {item.title}
              </h3>

              <p className="mt-1 text-xs text-zinc-500">
                {item.creatorName || item.creatorUsername || t("marketplace.unknownCreator")}
              </p>

              <div className="mt-2 flex items-center justify-between">
                <span className="rounded-md bg-zinc-100 px-2 py-0.5 text-xs font-medium text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400">
                  {getContentCategory(item.contentType)}
                </span>
                {item.lowestPriceCents !== null ? (
                  <span className="text-sm font-semibold text-zinc-900 dark:text-zinc-100">
                    {t("marketplace.fromPrice", {
                      price: (item.lowestPriceCents / 100).toFixed(2),
                    })}
                  </span>
                ) : (
                  <span className="text-xs text-zinc-400">{t("marketplace.noPrice")}</span>
                )}
              </div>
            </Link>
          ))}
        </div>
      )}

      {/* Pagination */}
      {total > LIMIT && (
        <div className="mt-6 flex items-center justify-center gap-4">
          <button
            onClick={() => setPage((p) => Math.max(0, p - 1))}
            disabled={page === 0}
            className="flex items-center gap-1 rounded-lg border border-zinc-200 px-3 py-1.5 text-sm font-medium disabled:opacity-50 dark:border-zinc-700"
          >
            <ChevronLeft className="h-4 w-4" /> {t("common.previous")}
          </button>
          <span className="text-sm text-zinc-500">
            {t("common.page", { current: page + 1, total: Math.ceil(total / LIMIT) })}
          </span>
          <button
            onClick={() => setPage((p) => p + 1)}
            disabled={(page + 1) * LIMIT >= total}
            className="flex items-center gap-1 rounded-lg border border-zinc-200 px-3 py-1.5 text-sm font-medium disabled:opacity-50 dark:border-zinc-700"
          >
            {t("common.next")} <ChevronRight className="h-4 w-4" />
          </button>
        </div>
      )}
    </div>
  );
}
