"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState, useCallback } from "react";
import { api } from "@/lib/api";
import Link from "next/link";
import {
  Upload,
  Image,
  Video,
  Music,
  FileText,
  File,
  Lock,
  Globe,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import { useTranslation } from "react-i18next";

type ContentItem = {
  id: string;
  title: string;
  description: string;
  contentType: string;
  fileSize: number;
  isPublic: boolean;
  tags: string[];
  thumbnailUrl?: string;
  createdAt: string;
};

const FILTERS = [
  { key: "all", labelKey: "vault.filterAll" },
  { key: "image", labelKey: "vault.filterImages", icon: Image },
  { key: "video", labelKey: "vault.filterVideo", icon: Video },
  { key: "audio", labelKey: "vault.filterAudio", icon: Music },
  { key: "text", labelKey: "vault.filterText", icon: FileText },
  { key: "other", labelKey: "vault.filterOther", icon: File },
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

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`;
}

const LIMIT = 20;

export default function VaultPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();
  const [items, setItems] = useState<ContentItem[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(0);
  const [filter, setFilter] = useState("all");
  const [loadingItems, setLoadingItems] = useState(true);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
  }, [user, loading, router]);

  const fetchItems = useCallback(async () => {
    setLoadingItems(true);
    const result = await api.content.list(LIMIT, page * LIMIT);
    if (result.data) {
      setItems(result.data.items || []);
      setTotal(result.data.total || 0);
    }
    setLoadingItems(false);
  }, [page]);

  useEffect(() => {
    if (user) fetchItems();
  }, [user, fetchItems]);

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [page]);

  if (loading || !user) return null;

  const filteredItems =
    filter === "all"
      ? items
      : items.filter((item) => getContentCategory(item.contentType) === filter);

  return (
    <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
      <div className="mb-8 flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">
            {t("vault.title")}
          </h1>
          <p className="mt-1 text-zinc-500">
            {t("vault.subtitle")}
          </p>
        </div>
        <Link
          href="/vault/upload"
          className="flex items-center gap-2 rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-700 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-300"
        >
          <Upload className="h-4 w-4" />
          {t("vault.uploadContent")}
        </Link>
      </div>

      {/* Filters */}
      <div className="mb-6 flex flex-wrap gap-2">
        {FILTERS.map((f) => (
          <button
            key={f.key}
            onClick={() => setFilter(f.key)}
            className={`flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm font-medium transition-colors ${
              filter === f.key
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
      {loadingItems ? (
        <p className="py-12 text-center text-zinc-400">{t("vault.loading")}</p>
      ) : filteredItems.length === 0 ? (
        <div className="py-16 text-center">
          <File className="mx-auto h-12 w-12 text-zinc-300 dark:text-zinc-600" />
          <p className="mt-4 text-zinc-500">
            {items.length === 0
              ? t("vault.emptyState")
              : t("vault.noFilterResults")}
          </p>
          {items.length === 0 && (
            <Link
              href="/vault/upload"
              className="mt-4 inline-flex items-center gap-2 rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-700 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-300"
            >
              <Upload className="h-4 w-4" />
              {t("vault.uploadFirst")}
            </Link>
          )}
        </div>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {filteredItems.map((item) => (
            <Link
              key={item.id}
              href={`/vault/item?id=${item.id}`}
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
              <div className="flex items-start justify-between gap-2">
                <h3 className="truncate text-sm font-semibold text-zinc-900 dark:text-zinc-100">
                  {item.title}
                </h3>
                {item.isPublic ? (
                  <Globe className="h-4 w-4 shrink-0 text-green-500" />
                ) : (
                  <Lock className="h-4 w-4 shrink-0 text-zinc-400" />
                )}
              </div>

              <div className="mt-2 flex items-center gap-2">
                <span className="rounded-md bg-zinc-100 px-2 py-0.5 text-xs font-medium text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400">
                  {getContentCategory(item.contentType)}
                </span>
                <span className="text-xs text-zinc-400">
                  {formatFileSize(item.fileSize)}
                </span>
              </div>

              <p className="mt-2 text-xs text-zinc-400">
                {new Date(item.createdAt).toLocaleDateString()}
              </p>
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
