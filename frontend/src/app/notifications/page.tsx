"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState, Suspense } from "react";
import { api } from "@/lib/api";
import { Bell, CheckCheck, ChevronLeft, ChevronRight } from "@/components/icons";
import { useTranslation } from "react-i18next";

type Notification = {
  id: string;
  type: string;
  title: string;
  message: string;
  isRead: boolean;
  createdAt: string;
};

function relativeTime(dateStr: string): string {
  const now = Date.now();
  const then = new Date(dateStr).getTime();
  const diffSec = Math.floor((now - then) / 1000);
  if (diffSec < 60) return "just now";
  const diffMin = Math.floor(diffSec / 60);
  if (diffMin < 60) return `${diffMin}m ago`;
  const diffHr = Math.floor(diffMin / 60);
  if (diffHr < 24) return `${diffHr}h ago`;
  const diffDay = Math.floor(diffHr / 24);
  return `${diffDay}d ago`;
}

const LIMIT = 20;

function NotificationsContent() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(0);
  const [loadingData, setLoadingData] = useState(true);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
  }, [user, loading, router]);

  useEffect(() => {
    if (!user) return;
    setLoadingData(true);
    api.notifications.list(LIMIT, page * LIMIT).then((result) => {
      if (result.data) {
        setNotifications((result.data as any).notifications ?? []);
        setTotal((result.data as any).total ?? 0);
      }
      setLoadingData(false);
    });
  }, [user, page]);

  const handleMarkRead = async (id: string) => {
    await api.notifications.markRead(id);
    setNotifications((prev) => prev.map((n) => (n.id === id ? { ...n, isRead: true } : n)));
  };

  const handleMarkAllRead = async () => {
    await api.notifications.markAllRead();
    setNotifications((prev) => prev.map((n) => ({ ...n, isRead: true })));
  };

  const totalPages = Math.ceil(total / LIMIT);
  const hasUnread = notifications.some((n) => !n.isRead);

  if (loading || !user) return null;

  return (
    <div className="mx-auto max-w-3xl px-4 py-12 sm:px-6">
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-zinc-900 dark:text-zinc-100">{t("notifications.pageTitle")}</h1>
          <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">{t("notifications.pageDescription")}</p>
        </div>
        {hasUnread && (
          <button onClick={handleMarkAllRead} className="flex items-center gap-2 rounded-lg border border-zinc-300 px-3 py-2 text-sm font-medium text-zinc-700 hover:bg-zinc-50 dark:border-zinc-700 dark:text-zinc-300 dark:hover:bg-zinc-800">
            <CheckCheck className="h-4 w-4" />
            {t("notifications.markAllRead")}
          </button>
        )}
      </div>

      {loadingData ? (
        <div className="flex items-center justify-center py-16">
          <div className="h-6 w-6 animate-spin rounded-full border-2 border-zinc-300 border-t-zinc-900 dark:border-zinc-600 dark:border-t-zinc-100" />
        </div>
      ) : notifications.length === 0 ? (
        <div className="rounded-xl border border-zinc-200 bg-zinc-50 py-16 text-center dark:border-zinc-800 dark:bg-zinc-900">
          <Bell className="mx-auto h-10 w-10 text-zinc-300 dark:text-zinc-600" />
          <p className="mt-4 text-sm font-medium text-zinc-500 dark:text-zinc-400">{t("notifications.emptyState")}</p>
        </div>
      ) : (
        <div className="space-y-2">
          {notifications.map((n) => (
            <button
              key={n.id}
              onClick={() => { if (!n.isRead) handleMarkRead(n.id); }}
              className={`flex w-full items-start gap-4 rounded-xl border px-5 py-4 text-left transition-colors ${!n.isRead ? "border-blue-100 bg-blue-50/50 hover:bg-blue-50 dark:border-blue-900/50 dark:bg-blue-950/20" : "border-zinc-200 bg-white hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-950 dark:hover:bg-zinc-900"}`}
            >
              <div className={`mt-0.5 flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full ${!n.isRead ? "bg-blue-100 text-blue-600 dark:bg-blue-900/50 dark:text-blue-400" : "bg-zinc-100 text-zinc-500 dark:bg-zinc-800 dark:text-zinc-400"}`}>
                <Bell className="h-4 w-4" />
              </div>
              <div className="min-w-0 flex-1">
                <div className="flex items-start justify-between gap-2">
                  <p className={`text-sm ${!n.isRead ? "font-semibold text-zinc-900 dark:text-zinc-100" : "font-medium text-zinc-700 dark:text-zinc-300"}`}>{n.title}</p>
                  <span className="flex-shrink-0 text-xs text-zinc-400 dark:text-zinc-500">{relativeTime(n.createdAt)}</span>
                </div>
                <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">{n.message}</p>
              </div>
              {!n.isRead && <div className="mt-2 flex-shrink-0"><div className="h-2.5 w-2.5 rounded-full bg-blue-500" /></div>}
            </button>
          ))}
        </div>
      )}

      {totalPages > 1 && (
        <div className="mt-8 flex items-center justify-between">
          <button onClick={() => setPage((p) => Math.max(0, p - 1))} disabled={page === 0} className="flex items-center gap-1 rounded-lg border border-zinc-300 px-3 py-2 text-sm font-medium text-zinc-700 hover:bg-zinc-50 disabled:opacity-40 dark:border-zinc-700 dark:text-zinc-300">
            <ChevronLeft className="h-4 w-4" />
            {t("common.previous")}
          </button>
          <span className="text-sm text-zinc-500 dark:text-zinc-400">{t("common.page", { current: page + 1, total: totalPages })}</span>
          <button onClick={() => setPage((p) => Math.min(totalPages - 1, p + 1))} disabled={page >= totalPages - 1} className="flex items-center gap-1 rounded-lg border border-zinc-300 px-3 py-2 text-sm font-medium text-zinc-700 hover:bg-zinc-50 disabled:opacity-40 dark:border-zinc-700 dark:text-zinc-300">
            {t("common.next")}
            <ChevronRight className="h-4 w-4" />
          </button>
        </div>
      )}
    </div>
  );
}

export default function NotificationsPage() {
  return (
    <Suspense fallback={null}>
      <NotificationsContent />
    </Suspense>
  );
}
