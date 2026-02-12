"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import Link from "next/link";
import { User, Link2, BarChart3, ArrowRight, QrCode, Eye, MousePointerClick } from "lucide-react";
import { QRCodeSVG } from "qrcode.react";
import { CopyLinkButton, ShareTwitterButton, ShareLinkedInButton } from "@/components/share-buttons";

const PROFILE_BASE_URL =
  process.env.NEXT_PUBLIC_PROFILE_URL || "https://creatrid.com";

export default function DashboardPage() {
  const { user, loading } = useAuth();
  const router = useRouter();

  const [connectionCount, setConnectionCount] = useState(0);
  const [showQR, setShowQR] = useState(false);
  const [analytics, setAnalytics] = useState<{
    totalViews: number;
    viewsToday: number;
    viewsThisWeek: number;
    totalClicks: number;
    clicksByType: Record<string, number>;
  } | null>(null);

  const profileUrl = user?.username
    ? `${PROFILE_BASE_URL}/profile?u=${user.username}`
    : "";

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
    if (!loading && user && !user.onboarded) router.push("/onboarding");
  }, [user, loading, router]);

  useEffect(() => {
    if (user) {
      api.connections.list().then((result) => {
        if (result.data) {
          setConnectionCount(result.data.connections?.length || 0);
        }
      });
      api.analytics.summary().then((result) => {
        if (result.data) {
          setAnalytics(result.data);
        }
      });
    }
  }, [user]);

  if (loading || !user) return null;

  const profileComplete =
    [user.name, user.username, user.bio, user.image].filter(Boolean).length;
  const profileTotal = 4;

  return (
    <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold">
          Welcome back, {user.name || "Creator"}
        </h1>
        <p className="mt-1 text-zinc-500">
          Manage your Creator Passport from here.
        </p>
      </div>

      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {/* Profile Card */}
        <Link
          href="/settings"
          className="group rounded-xl border border-zinc-200 p-6 transition-colors hover:border-zinc-300 dark:border-zinc-800 dark:hover:border-zinc-700"
        >
          <div className="mb-4 flex items-center justify-between">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
              <User className="h-5 w-5" />
            </div>
            <ArrowRight className="h-4 w-4 text-zinc-400 transition-transform group-hover:translate-x-1" />
          </div>
          <h3 className="font-semibold">Profile</h3>
          <p className="mt-1 text-sm text-zinc-500">
            {profileComplete}/{profileTotal} complete
          </p>
          <div className="mt-3 h-2 rounded-full bg-zinc-100 dark:bg-zinc-800">
            <div
              className="h-2 rounded-full bg-zinc-900 dark:bg-zinc-100"
              style={{
                width: `${(profileComplete / profileTotal) * 100}%`,
              }}
            />
          </div>
        </Link>

        {/* Connections Card */}
        <Link
          href="/connections"
          className="group rounded-xl border border-zinc-200 p-6 transition-colors hover:border-zinc-300 dark:border-zinc-800 dark:hover:border-zinc-700"
        >
          <div className="mb-4 flex items-center justify-between">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
              <Link2 className="h-5 w-5" />
            </div>
            <ArrowRight className="h-4 w-4 text-zinc-400 transition-transform group-hover:translate-x-1" />
          </div>
          <h3 className="font-semibold">Connections</h3>
          <p className="mt-1 text-sm text-zinc-500">
            Connect your social accounts
          </p>
          <p className="mt-3 text-2xl font-bold">{connectionCount}</p>
        </Link>

        {/* Score Card */}
        <div className="rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
          <div className="mb-4">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
              <BarChart3 className="h-5 w-5" />
            </div>
          </div>
          <h3 className="font-semibold">Creator Score</h3>
          <p className="mt-1 text-sm text-zinc-500">
            {user.creatorScore !== null
              ? "Based on profile, connections & reach"
              : "Complete your profile to unlock"}
          </p>
          <p
            className={`mt-3 text-2xl font-bold ${
              user.creatorScore !== null
                ? "text-zinc-900 dark:text-zinc-100"
                : "text-zinc-300 dark:text-zinc-600"
            }`}
          >
            {user.creatorScore ?? "—"}
            {user.creatorScore !== null && (
              <span className="text-sm font-normal text-zinc-400"> / 100</span>
            )}
          </p>
          {user.creatorScore !== null && (
            <div className="mt-3 h-2 rounded-full bg-zinc-100 dark:bg-zinc-800">
              <div
                className="h-2 rounded-full bg-zinc-900 dark:bg-zinc-100"
                style={{ width: `${user.creatorScore}%` }}
              />
            </div>
          )}
        </div>
      </div>

      {/* Analytics */}
      {analytics && (
        <div className="mt-6 grid gap-6 sm:grid-cols-3">
          <div className="rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
            <div className="mb-2 flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
              <Eye className="h-5 w-5" />
            </div>
            <h3 className="font-semibold">Profile Views</h3>
            <p className="mt-2 text-2xl font-bold">{analytics.totalViews}</p>
            <p className="mt-1 text-xs text-zinc-500">
              {analytics.viewsToday} today &middot; {analytics.viewsThisWeek} this week
            </p>
          </div>
          <div className="rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
            <div className="mb-2 flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
              <MousePointerClick className="h-5 w-5" />
            </div>
            <h3 className="font-semibold">Link Clicks</h3>
            <p className="mt-2 text-2xl font-bold">{analytics.totalClicks}</p>
            <div className="mt-1 flex flex-wrap gap-2 text-xs text-zinc-500">
              {Object.entries(analytics.clicksByType).map(([type, count]) => (
                <span key={type} className="capitalize">
                  {type}: {count}
                </span>
              ))}
              {Object.keys(analytics.clicksByType).length === 0 && (
                <span>No clicks yet</span>
              )}
            </div>
          </div>
          <div className="rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
            <div className="mb-2 flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
              <BarChart3 className="h-5 w-5" />
            </div>
            <h3 className="font-semibold">Engagement</h3>
            <p className="mt-2 text-2xl font-bold">
              {analytics.totalViews > 0
                ? ((analytics.totalClicks / analytics.totalViews) * 100).toFixed(1)
                : "0.0"}
              <span className="text-sm font-normal text-zinc-400">%</span>
            </p>
            <p className="mt-1 text-xs text-zinc-500">Click-through rate</p>
          </div>
        </div>
      )}

      {/* Public Profile Link */}
      {user.username && (
        <div className="mt-8 rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
          <div className="flex items-start justify-between">
            <div>
              <h3 className="font-semibold">Your Public Profile</h3>
              <p className="mt-1 text-sm text-zinc-500">
                Share this link with brands and collaborators.
              </p>
            </div>
            <Link
              href={`/profile?u=${user.username}`}
              className="rounded-lg border border-zinc-200 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
            >
              Preview
            </Link>
          </div>
          <div className="mt-4 flex flex-wrap items-center gap-2">
            <CopyLinkButton url={profileUrl} />
            <ShareTwitterButton
              url={profileUrl}
              text={`Check out my Creator Passport on Creatrid`}
            />
            <ShareLinkedInButton url={profileUrl} />
            <button
              onClick={() => setShowQR(!showQR)}
              className="flex items-center gap-1.5 rounded-lg border border-zinc-200 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
            >
              <QrCode className="h-3.5 w-3.5" />
              QR Code
            </button>
          </div>
          {showQR && (
            <div className="mt-4 flex items-center gap-4">
              <div className="rounded-lg border border-zinc-200 bg-white p-3 dark:border-zinc-800 dark:bg-zinc-900">
                <QRCodeSVG value={profileUrl} size={120} level="M" marginSize={2} />
              </div>
              <div className="text-sm text-zinc-500">
                <p>Scan this QR code to open your Creator Passport.</p>
                <p className="mt-1 text-xs text-zinc-400">
                  Perfect for business cards, portfolios, and events.
                </p>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
