"use client";

import { useAuth } from "@/lib/auth-context";
import { useEffect, useState, useRef, useCallback } from "react";
import { api } from "@/lib/api";
import Link from "next/link";
import { Bell, Check, CheckCheck } from "@/components/icons";
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

export function NotificationBell() {
  const { user } = useAuth();
  const { t } = useTranslation();
  const [unreadCount, setUnreadCount] = useState(0);
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [open, setOpen] = useState(false);
  const [loadingList, setLoadingList] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!user) return;

    // Initial fetch of unread count
    const fetchCount = () => {
      api.notifications.unreadCount().then((r) => {
        if (r.data) setUnreadCount((r.data as any).count ?? 0);
      });
    };
    fetchCount();

    // Try SSE for real-time updates, fallback to polling
    const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    let fallbackInterval: ReturnType<typeof setInterval> | null = null;
    let es: EventSource | null = null;

    try {
      es = new EventSource(`${API_URL}/api/notifications/stream`, {
        withCredentials: true,
      });

      es.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          if (data.unreadCount != null) {
            setUnreadCount(data.unreadCount);
          } else {
            setUnreadCount((prev) => prev + 1);
          }
          if (data.notification) {
            setNotifications((prev) => [data.notification, ...prev].slice(0, 5));
          }
        } catch {
          setUnreadCount((prev) => prev + 1);
        }
      };

      es.onerror = () => {
        // Close SSE and fallback to polling
        if (es) {
          es.close();
          es = null;
        }
        if (!fallbackInterval) {
          fallbackInterval = setInterval(fetchCount, 30000);
        }
      };
    } catch {
      // SSE not supported, use polling
      fallbackInterval = setInterval(fetchCount, 30000);
    }

    return () => {
      if (es) es.close();
      if (fallbackInterval) clearInterval(fallbackInterval);
    };
  }, [user]);

  const fetchNotifications = useCallback(async () => {
    setLoadingList(true);
    const result = await api.notifications.list(5, 0);
    if (result.data) setNotifications((result.data as any).notifications ?? []);
    setLoadingList(false);
  }, []);

  useEffect(() => {
    if (open) fetchNotifications();
  }, [open, fetchNotifications]);

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) setOpen(false);
    }
    if (open) document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [open]);

  const handleMarkRead = async (id: string) => {
    await api.notifications.markRead(id);
    setNotifications((prev) => prev.map((n) => (n.id === id ? { ...n, isRead: true } : n)));
    setUnreadCount((prev) => Math.max(0, prev - 1));
  };

  const handleMarkAllRead = async () => {
    await api.notifications.markAllRead();
    setNotifications((prev) => prev.map((n) => ({ ...n, isRead: true })));
    setUnreadCount(0);
  };

  if (!user) return null;

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        onClick={() => setOpen((prev) => !prev)}
        className="relative rounded-lg p-2 text-zinc-500 transition-colors hover:bg-zinc-100 hover:text-zinc-900 dark:text-zinc-400 dark:hover:bg-zinc-800 dark:hover:text-zinc-100"
        aria-label={t("notifications.bell")}
      >
        <Bell className="h-5 w-5" />
        {unreadCount > 0 && (
          <span className="absolute -right-0.5 -top-0.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-red-500 px-1 text-[10px] font-bold text-white">
            {unreadCount > 99 ? "99+" : unreadCount}
          </span>
        )}
      </button>

      {open && (
        <div className="absolute right-0 top-full mt-2 w-80 rounded-xl border border-zinc-200 bg-white shadow-lg dark:border-zinc-800 dark:bg-zinc-950">
          <div className="flex items-center justify-between border-b border-zinc-200 px-4 py-3 dark:border-zinc-800">
            <h3 className="text-sm font-semibold text-zinc-900 dark:text-zinc-100">{t("notifications.title")}</h3>
            {unreadCount > 0 && (
              <button onClick={handleMarkAllRead} className="flex items-center gap-1 text-xs text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100">
                <CheckCheck className="h-3.5 w-3.5" />
                {t("notifications.markAllRead")}
              </button>
            )}
          </div>
          <div className="max-h-80 overflow-y-auto">
            {loadingList ? (
              <div className="flex items-center justify-center py-8">
                <div className="h-5 w-5 animate-spin rounded-full border-2 border-zinc-300 border-t-zinc-900 dark:border-zinc-600 dark:border-t-zinc-100" />
              </div>
            ) : notifications.length === 0 ? (
              <div className="py-8 text-center text-sm text-zinc-500 dark:text-zinc-400">{t("notifications.empty")}</div>
            ) : (
              notifications.map((n) => (
                <button
                  key={n.id}
                  onClick={() => { if (!n.isRead) handleMarkRead(n.id); }}
                  className={`flex w-full items-start gap-3 px-4 py-3 text-left transition-colors hover:bg-zinc-50 dark:hover:bg-zinc-900 ${!n.isRead ? "bg-zinc-50/50 dark:bg-zinc-900/50" : ""}`}
                >
                  <div className="mt-1.5 flex-shrink-0">
                    {!n.isRead ? <div className="h-2 w-2 rounded-full bg-blue-500" /> : <div className="h-2 w-2" />}
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className={`text-sm ${!n.isRead ? "font-semibold text-zinc-900 dark:text-zinc-100" : "text-zinc-700 dark:text-zinc-300"}`}>{n.title}</p>
                    <p className="mt-0.5 line-clamp-2 text-xs text-zinc-500 dark:text-zinc-400">{n.message}</p>
                    <p className="mt-1 text-[10px] text-zinc-400 dark:text-zinc-500">{relativeTime(n.createdAt)}</p>
                  </div>
                  {!n.isRead && <Check className="mt-1 h-3.5 w-3.5 flex-shrink-0 text-zinc-400" />}
                </button>
              ))
            )}
          </div>
          <div className="border-t border-zinc-200 px-4 py-2 dark:border-zinc-800">
            <Link href="/notifications" onClick={() => setOpen(false)} className="block text-center text-xs font-medium text-zinc-600 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100">
              {t("notifications.viewAll")}
            </Link>
          </div>
        </div>
      )}
    </div>
  );
}
