"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { useTranslation } from "react-i18next";
import { Eye, MousePointerClick, TrendingUp, CalendarDays, Download } from "lucide-react";
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import type { ReactNode } from "react";

type AnalyticsData = {
  totalViews: number;
  viewsToday: number;
  viewsThisWeek: number;
  totalClicks: number;
  clicksByType: Record<string, number>;
  viewsByDay: { date: string; count: number }[];
  clicksByDay: { date: string; count: number }[];
  topReferrers: { referrer: string; count: number }[];
  viewsByHour: { hour: number; count: number }[];
};

function formatDate(dateStr: string) {
  const d = new Date(dateStr);
  return d.toLocaleDateString("en-US", { month: "short", day: "numeric" });
}

function formatDateLabel(label: ReactNode): string {
  return formatDate(String(label));
}

function formatHour(hour: number) {
  if (hour === 0) return "12am";
  if (hour === 12) return "12pm";
  return hour < 12 ? `${hour}am` : `${hour - 12}pm`;
}

function formatHourLabel(label: ReactNode): string {
  return formatHour(Number(label));
}

export default function AnalyticsPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();
  const [data, setData] = useState<AnalyticsData | null>(null);
  const [fetching, setFetching] = useState(true);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
    if (!loading && user && !user.onboarded) router.push("/onboarding");
  }, [user, loading, router]);

  useEffect(() => {
    if (user) {
      api.analytics.summary().then((result) => {
        if (result.data) setData(result.data);
        setFetching(false);
      });
    }
  }, [user]);

  if (loading || !user) return null;

  if (fetching) {
    return (
      <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
        <p className="text-zinc-500">{t("common.loading")}</p>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
        <p className="text-zinc-500">{t("analytics.noData")}</p>
      </div>
    );
  }

  const ctr =
    data.totalViews > 0
      ? ((data.totalClicks / data.totalViews) * 100).toFixed(1)
      : "0.0";

  return (
    <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("analytics.title")}</h1>
          <p className="mt-1 text-zinc-500">{t("analytics.subtitle")}</p>
        </div>
        <a
          href={api.analyticsExport.profileUrl()}
          className="flex items-center gap-2 rounded-lg border border-zinc-200 px-3 py-2 text-sm font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
        >
          <Download className="h-4 w-4" />
          {t("analytics.exportCSV")}
        </a>
      </div>

      {/* Summary Cards */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <div className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
          <div className="mb-3 flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
            <Eye className="h-5 w-5" />
          </div>
          <p className="text-sm text-zinc-500">{t("analytics.totalViews")}</p>
          <p className="mt-1 text-2xl font-bold">{data.totalViews}</p>
          <p className="mt-1 text-xs text-zinc-400">
            {t("analytics.todayViews", { count: data.viewsToday })}
          </p>
        </div>

        <div className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
          <div className="mb-3 flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
            <MousePointerClick className="h-5 w-5" />
          </div>
          <p className="text-sm text-zinc-500">{t("analytics.totalClicks")}</p>
          <p className="mt-1 text-2xl font-bold">{data.totalClicks}</p>
        </div>

        <div className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
          <div className="mb-3 flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
            <TrendingUp className="h-5 w-5" />
          </div>
          <p className="text-sm text-zinc-500">{t("analytics.clickThrough")}</p>
          <p className="mt-1 text-2xl font-bold">
            {ctr}
            <span className="text-sm font-normal text-zinc-400">%</span>
          </p>
        </div>

        <div className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
          <div className="mb-3 flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
            <CalendarDays className="h-5 w-5" />
          </div>
          <p className="text-sm text-zinc-500">{t("analytics.weekViews")}</p>
          <p className="mt-1 text-2xl font-bold">{data.viewsThisWeek}</p>
        </div>
      </div>

      {/* Charts Grid */}
      <div className="mt-8 grid gap-6 lg:grid-cols-2">
        {/* Views Over Time */}
        <div className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
          <h3 className="mb-1 font-semibold">{t("analytics.viewsOverTime")}</h3>
          <p className="mb-4 text-xs text-zinc-400">{t("analytics.last30Days")}</p>
          {data.viewsByDay.length > 0 ? (
            <ResponsiveContainer width="100%" height={260}>
              <LineChart data={data.viewsByDay}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-zinc-200 dark:stroke-zinc-700" />
                <XAxis
                  dataKey="date"
                  tickFormatter={formatDate}
                  tick={{ fontSize: 11 }}
                  className="text-zinc-500"
                />
                <YAxis
                  allowDecimals={false}
                  tick={{ fontSize: 11 }}
                  className="text-zinc-500"
                />
                <Tooltip
                  labelFormatter={formatDateLabel}
                  contentStyle={{
                    backgroundColor: "var(--tooltip-bg, #fff)",
                    border: "1px solid var(--tooltip-border, #e4e4e7)",
                    borderRadius: "8px",
                    fontSize: "12px",
                  }}
                />
                <Line
                  type="monotone"
                  dataKey="count"
                  stroke="#3f3f46"
                  strokeWidth={2}
                  dot={false}
                  activeDot={{ r: 4 }}
                />
              </LineChart>
            </ResponsiveContainer>
          ) : (
            <p className="py-12 text-center text-sm text-zinc-400">
              {t("analytics.noData")}
            </p>
          )}
        </div>

        {/* Clicks Over Time */}
        <div className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
          <h3 className="mb-1 font-semibold">{t("analytics.clicksOverTime")}</h3>
          <p className="mb-4 text-xs text-zinc-400">{t("analytics.last30Days")}</p>
          {data.clicksByDay.length > 0 ? (
            <ResponsiveContainer width="100%" height={260}>
              <LineChart data={data.clicksByDay}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-zinc-200 dark:stroke-zinc-700" />
                <XAxis
                  dataKey="date"
                  tickFormatter={formatDate}
                  tick={{ fontSize: 11 }}
                  className="text-zinc-500"
                />
                <YAxis
                  allowDecimals={false}
                  tick={{ fontSize: 11 }}
                  className="text-zinc-500"
                />
                <Tooltip
                  labelFormatter={formatDateLabel}
                  contentStyle={{
                    backgroundColor: "var(--tooltip-bg, #fff)",
                    border: "1px solid var(--tooltip-border, #e4e4e7)",
                    borderRadius: "8px",
                    fontSize: "12px",
                  }}
                />
                <Line
                  type="monotone"
                  dataKey="count"
                  stroke="#3f3f46"
                  strokeWidth={2}
                  dot={false}
                  activeDot={{ r: 4 }}
                />
              </LineChart>
            </ResponsiveContainer>
          ) : (
            <p className="py-12 text-center text-sm text-zinc-400">
              {t("analytics.noData")}
            </p>
          )}
        </div>

        {/* Top Referrers */}
        <div className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
          <h3 className="mb-1 font-semibold">{t("analytics.topReferrers")}</h3>
          <p className="mb-4 text-xs text-zinc-400">{t("analytics.last30Days")}</p>
          {data.topReferrers.length > 0 ? (
            <ResponsiveContainer width="100%" height={260}>
              <BarChart
                data={data.topReferrers}
                layout="vertical"
                margin={{ left: 20 }}
              >
                <CartesianGrid strokeDasharray="3 3" className="stroke-zinc-200 dark:stroke-zinc-700" />
                <XAxis
                  type="number"
                  allowDecimals={false}
                  tick={{ fontSize: 11 }}
                  className="text-zinc-500"
                />
                <YAxis
                  type="category"
                  dataKey="referrer"
                  tick={{ fontSize: 11 }}
                  width={100}
                  className="text-zinc-500"
                />
                <Tooltip
                  contentStyle={{
                    backgroundColor: "var(--tooltip-bg, #fff)",
                    border: "1px solid var(--tooltip-border, #e4e4e7)",
                    borderRadius: "8px",
                    fontSize: "12px",
                  }}
                />
                <Bar dataKey="count" fill="#52525b" radius={[0, 4, 4, 0]} />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <p className="py-12 text-center text-sm text-zinc-400">
              {t("analytics.noData")}
            </p>
          )}
        </div>

        {/* Views by Hour */}
        <div className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
          <h3 className="mb-1 font-semibold">{t("analytics.viewsByHour")}</h3>
          <p className="mb-4 text-xs text-zinc-400">{t("analytics.last30Days")}</p>
          {data.viewsByHour.length > 0 ? (
            <ResponsiveContainer width="100%" height={260}>
              <BarChart data={data.viewsByHour}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-zinc-200 dark:stroke-zinc-700" />
                <XAxis
                  dataKey="hour"
                  tickFormatter={formatHour}
                  tick={{ fontSize: 11 }}
                  className="text-zinc-500"
                />
                <YAxis
                  allowDecimals={false}
                  tick={{ fontSize: 11 }}
                  className="text-zinc-500"
                />
                <Tooltip
                  labelFormatter={formatHourLabel}
                  contentStyle={{
                    backgroundColor: "var(--tooltip-bg, #fff)",
                    border: "1px solid var(--tooltip-border, #e4e4e7)",
                    borderRadius: "8px",
                    fontSize: "12px",
                  }}
                />
                <Bar dataKey="count" fill="#52525b" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <p className="py-12 text-center text-sm text-zinc-400">
              {t("analytics.noData")}
            </p>
          )}
        </div>
      </div>

      {/* Clicks by Type */}
      <div className="mt-6 rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
        <h3 className="mb-4 font-semibold">{t("analytics.clicksByType")}</h3>
        {Object.keys(data.clicksByType).length > 0 ? (
          <div className="space-y-3">
            {Object.entries(data.clicksByType)
              .sort(([, a], [, b]) => b - a)
              .map(([type, count]) => {
                const max = Math.max(
                  ...Object.values(data.clicksByType)
                );
                const pct = max > 0 ? (count / max) * 100 : 0;
                return (
                  <div key={type} className="flex items-center gap-4">
                    <span className="w-24 text-sm capitalize text-zinc-600 dark:text-zinc-400">
                      {type}
                    </span>
                    <div className="flex-1">
                      <div className="h-6 rounded-md bg-zinc-100 dark:bg-zinc-800">
                        <div
                          className="flex h-6 items-center rounded-md bg-zinc-300 px-2 text-xs font-medium dark:bg-zinc-600"
                          style={{ width: `${Math.max(pct, 8)}%` }}
                        >
                          {count}
                        </div>
                      </div>
                    </div>
                  </div>
                );
              })}
          </div>
        ) : (
          <p className="py-6 text-center text-sm text-zinc-400">
            {t("analytics.noData")}
          </p>
        )}
      </div>
    </div>
  );
}
