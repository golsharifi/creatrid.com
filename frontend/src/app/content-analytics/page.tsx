"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState, Suspense } from "react";
import { api } from "@/lib/api";
import { BarChart3, Eye, Download, DollarSign } from "lucide-react";
import { useTranslation } from "react-i18next";

type ContentAnalyticsItem = {
  contentId: string;
  title: string;
  totalViews: number;
  totalDownloads: number;
  revenueCents: number;
};

function ContentAnalyticsContent() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();
  const [items, setItems] = useState<ContentAnalyticsItem[]>([]);
  const [totalRevenue, setTotalRevenue] = useState(0);
  const [loadingData, setLoadingData] = useState(true);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
  }, [user, loading, router]);

  useEffect(() => {
    if (!user) return;
    api.contentAnalytics.summary().then((result) => {
      if (result.data) {
        const data = result.data as any;
        setItems(data.items ?? []);
        setTotalRevenue(data.totalRevenue ?? 0);
      }
      setLoadingData(false);
    });
  }, [user]);

  if (loading || !user) return null;

  const totalViews = items.reduce((sum, i) => sum + i.totalViews, 0);
  const totalDownloads = items.reduce((sum, i) => sum + i.totalDownloads, 0);

  return (
    <div className="mx-auto max-w-5xl px-4 py-12 sm:px-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold tracking-tight text-zinc-900 dark:text-zinc-100">{t("contentAnalytics.title")}</h1>
        <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">{t("contentAnalytics.subtitle")}</p>
      </div>

      {/* Summary cards */}
      <div className="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-3">
        <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-950">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-100 text-blue-600 dark:bg-blue-900/50 dark:text-blue-400">
              <Eye className="h-5 w-5" />
            </div>
            <div>
              <p className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">{totalViews.toLocaleString()}</p>
              <p className="text-xs text-zinc-500 dark:text-zinc-400">{t("contentAnalytics.totalViews")}</p>
            </div>
          </div>
        </div>
        <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-950">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-green-100 text-green-600 dark:bg-green-900/50 dark:text-green-400">
              <Download className="h-5 w-5" />
            </div>
            <div>
              <p className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">{totalDownloads.toLocaleString()}</p>
              <p className="text-xs text-zinc-500 dark:text-zinc-400">{t("contentAnalytics.totalDownloads")}</p>
            </div>
          </div>
        </div>
        <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-950">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-purple-100 text-purple-600 dark:bg-purple-900/50 dark:text-purple-400">
              <DollarSign className="h-5 w-5" />
            </div>
            <div>
              <p className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">${(totalRevenue / 100).toFixed(2)}</p>
              <p className="text-xs text-zinc-500 dark:text-zinc-400">{t("contentAnalytics.totalRevenue")}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Content table */}
      {loadingData ? (
        <div className="flex items-center justify-center py-16">
          <div className="h-6 w-6 animate-spin rounded-full border-2 border-zinc-300 border-t-zinc-900 dark:border-zinc-600 dark:border-t-zinc-100" />
        </div>
      ) : items.length === 0 ? (
        <div className="rounded-xl border border-zinc-200 bg-zinc-50 py-16 text-center dark:border-zinc-800 dark:bg-zinc-900">
          <BarChart3 className="mx-auto h-10 w-10 text-zinc-300 dark:text-zinc-600" />
          <p className="mt-4 text-sm font-medium text-zinc-500 dark:text-zinc-400">{t("contentAnalytics.noData")}</p>
        </div>
      ) : (
        <div className="overflow-hidden rounded-xl border border-zinc-200 dark:border-zinc-800">
          <table className="w-full text-left text-sm">
            <thead className="border-b border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900">
              <tr>
                <th className="px-4 py-3 font-medium text-zinc-600 dark:text-zinc-400">{t("contentAnalytics.content")}</th>
                <th className="px-4 py-3 text-right font-medium text-zinc-600 dark:text-zinc-400">{t("contentAnalytics.views")}</th>
                <th className="px-4 py-3 text-right font-medium text-zinc-600 dark:text-zinc-400">{t("contentAnalytics.downloads")}</th>
                <th className="px-4 py-3 text-right font-medium text-zinc-600 dark:text-zinc-400">{t("contentAnalytics.revenue")}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800">
              {items.map((item) => (
                <tr key={item.contentId} className="bg-white dark:bg-zinc-950">
                  <td className="px-4 py-3 font-medium text-zinc-900 dark:text-zinc-100">{item.title}</td>
                  <td className="px-4 py-3 text-right text-zinc-600 dark:text-zinc-400">{item.totalViews.toLocaleString()}</td>
                  <td className="px-4 py-3 text-right text-zinc-600 dark:text-zinc-400">{item.totalDownloads.toLocaleString()}</td>
                  <td className="px-4 py-3 text-right text-zinc-600 dark:text-zinc-400">${(item.revenueCents / 100).toFixed(2)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

export default function ContentAnalyticsPage() {
  return (
    <Suspense fallback={null}>
      <ContentAnalyticsContent />
    </Suspense>
  );
}
