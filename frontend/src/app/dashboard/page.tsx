"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState, Suspense } from "react";
import { api } from "@/lib/api";
import Link from "next/link";
import { User, Link2, BarChart3, ArrowRight, QrCode, Eye, MousePointerClick, Mail, CheckCircle, Code, Archive, Users, Upload, Share2, TrendingUp, Activity, Inbox } from "@/components/icons";
import { OnboardingChecklist } from "@/components/onboarding-checklist";
import { TierBadge } from "@/components/tier-badge";
import { QRCodeSVG } from "qrcode.react";
import { CopyLinkButton, ShareTwitterButton, ShareLinkedInButton } from "@/components/share-buttons";
import { AreaChart, Area, ResponsiveContainer, Tooltip, XAxis } from "recharts";
import { useTranslation } from "react-i18next";

const PROFILE_BASE_URL =
  process.env.NEXT_PUBLIC_PROFILE_URL || "https://creatrid.com";

function DashboardContent() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();
  const { t } = useTranslation();

  const [connectionCount, setConnectionCount] = useState(0);
  const [contentCount, setContentCount] = useState(0);
  const [showQR, setShowQR] = useState(false);
  const [sendingVerification, setSendingVerification] = useState(false);
  const [verificationSent, setVerificationSent] = useState(false);
  const justVerified = searchParams.get("verified") === "true";
  const [analytics, setAnalytics] = useState<{
    totalViews: number;
    viewsToday: number;
    viewsThisWeek: number;
    totalClicks: number;
    clicksByType: Record<string, number>;
    viewsByDay: { date: string; count: number }[];
  } | null>(null);
  const [recommendations, setRecommendations] = useState<any[]>([]);
  const [pendingCollabCount, setPendingCollabCount] = useState(0);
  const [recentNotifications, setRecentNotifications] = useState<any[]>([]);
  const [copiedProfile, setCopiedProfile] = useState(false);
  const [connections, setConnections] = useState<any[]>([]);

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
          setConnections(result.data.connections || []);
        }
      });
      api.analytics.summary().then((result) => {
        if (result.data) {
          setAnalytics(result.data);
        }
      });
      api.content.list(1, 0).then((result) => {
        if (result.data) {
          setContentCount(result.data.total || 0);
        }
      });
      api.recommendations.list().then((result) => {
        if (result.data) {
          setRecommendations(result.data.creators || []);
        }
      });
      api.collaborations.inbox().then((result) => {
        if (result.data) {
          const pending = (result.data.requests || []).filter(
            (r: any) => r.status === "pending"
          );
          setPendingCollabCount(pending.length);
        }
      });
      api.notifications.list(5, 0).then((result) => {
        if (result.data) {
          setRecentNotifications(result.data.notifications || []);
        }
      });
    }
  }, [user]);

  const handleSendVerification = async () => {
    setSendingVerification(true);
    const result = await api.auth.sendVerificationEmail();
    setSendingVerification(false);
    if (result.data) setVerificationSent(true);
  };

  if (loading || !user) return null;

  const profileComplete =
    [user.name, user.username, user.bio, user.image].filter(Boolean).length;
  const profileTotal = 4;

  return (
    <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
      {/* Email verified success banner */}
      {justVerified && (
        <div className="mb-6 flex items-center gap-3 rounded-xl border border-emerald-200 bg-emerald-50 p-4 dark:border-emerald-800 dark:bg-emerald-950">
          <CheckCircle className="h-5 w-5 text-emerald-600" />
          <p className="text-sm font-medium text-emerald-800 dark:text-emerald-200">
            {t("dashboard.emailVerifiedSuccess")}
          </p>
        </div>
      )}

      {/* Email verification banner */}
      {!user.emailVerified && !justVerified && (
        <div className="mb-6 flex flex-col gap-3 rounded-xl border border-amber-200 bg-amber-50 p-4 dark:border-amber-800 dark:bg-amber-950 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-3">
            <Mail className="h-5 w-5 text-amber-600" />
            <div>
              <p className="text-sm font-medium text-amber-800 dark:text-amber-200">
                {t("dashboard.verifyEmailTitle")}
              </p>
              <p className="text-xs text-amber-600 dark:text-amber-400">
                {t("dashboard.verifyEmailDesc")}
              </p>
            </div>
          </div>
          <button
            onClick={handleSendVerification}
            disabled={sendingVerification || verificationSent}
            className="rounded-lg bg-amber-600 px-4 py-1.5 text-xs font-medium text-white transition-colors hover:bg-amber-700 disabled:opacity-50"
          >
            {verificationSent
              ? t("dashboard.verificationSent")
              : sendingVerification
                ? t("dashboard.sending")
                : t("dashboard.verifyEmail")}
          </button>
        </div>
      )}

      <div className="mb-8">
        <h1 className="text-2xl font-bold">
          {t("dashboard.welcomeBack", { name: user.name || t("common.creator") })}
        </h1>
        <p className="mt-1 text-zinc-500">
          {t("dashboard.managePassport")}
        </p>
      </div>

      {/* Onboarding Checklist */}
      <div className="mb-6">
        <OnboardingChecklist
          user={{
            image: user.image || null,
            bio: user.bio || null,
            emailVerified: user.emailVerified || null,
          }}
          connectionCount={connectionCount}
          contentCount={contentCount}
        />
      </div>

      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {/* Content Vault Card */}
        <Link
          href="/vault"
          className="group rounded-xl border border-zinc-200 p-6 transition-colors hover:border-zinc-300 dark:border-zinc-800 dark:hover:border-zinc-700"
        >
          <div className="mb-4 flex items-center justify-between">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
              <Archive className="h-5 w-5" />
            </div>
            <ArrowRight className="h-4 w-4 text-zinc-400 transition-transform group-hover:translate-x-1" />
          </div>
          <h3 className="font-semibold">{t("dashboard.contentVault")}</h3>
          <p className="mt-1 text-sm text-zinc-500">
            {t("dashboard.contentVaultDesc")}
          </p>
          <p className="mt-3 text-2xl font-bold">{contentCount}</p>
        </Link>

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
          <h3 className="font-semibold">{t("dashboard.profile")}</h3>
          <p className="mt-1 text-sm text-zinc-500">
            {t("dashboard.profileComplete", { complete: profileComplete, total: profileTotal })}
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
          <h3 className="font-semibold">{t("dashboard.connections")}</h3>
          <p className="mt-1 text-sm text-zinc-500">
            {t("dashboard.connectSocial")}
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
          <h3 className="font-semibold">{t("dashboard.creatorScore")}</h3>
          <p className="mt-1 text-sm text-zinc-500">
            {user.creatorScore !== null
              ? t("dashboard.scoreBasedOn")
              : t("dashboard.scoreUnlock")}
          </p>
          <div className="mt-3 flex items-center gap-2">
            <span
              className={`text-2xl font-bold ${
                user.creatorScore !== null
                  ? "text-zinc-900 dark:text-zinc-100"
                  : "text-zinc-300 dark:text-zinc-600"
              }`}
            >
              {user.creatorScore ?? t("common.noData")}
              {user.creatorScore !== null && (
                <span className="text-sm font-normal text-zinc-400"> {t("dashboard.outOf100")}</span>
              )}
            </span>
            {user.creatorTier && (
              <TierBadge tier={user.creatorTier} />
            )}
          </div>
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
          <Link
            href="/analytics"
            className="group rounded-xl border border-zinc-200 p-6 transition-colors hover:border-zinc-300 dark:border-zinc-800 dark:hover:border-zinc-700"
          >
            <div className="mb-2 flex items-center justify-between">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
                <Eye className="h-5 w-5" />
              </div>
              <ArrowRight className="h-4 w-4 text-zinc-400 transition-transform group-hover:translate-x-1" />
            </div>
            <h3 className="font-semibold">{t("dashboard.profileViews")}</h3>
            <p className="mt-2 text-2xl font-bold">{analytics.totalViews}</p>
            <p className="mt-1 text-xs text-zinc-500">
              {t("dashboard.todayAndWeek", { today: analytics.viewsToday, week: analytics.viewsThisWeek })}
            </p>
          </Link>
          <div className="rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
            <div className="mb-2 flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
              <MousePointerClick className="h-5 w-5" />
            </div>
            <h3 className="font-semibold">{t("dashboard.linkClicks")}</h3>
            <p className="mt-2 text-2xl font-bold">{analytics.totalClicks}</p>
            <div className="mt-1 flex flex-wrap gap-2 text-xs text-zinc-500">
              {Object.entries(analytics.clicksByType).map(([type, count]) => (
                <span key={type} className="capitalize">
                  {type}: {count}
                </span>
              ))}
              {Object.keys(analytics.clicksByType).length === 0 && (
                <span>{t("dashboard.noClicks")}</span>
              )}
            </div>
          </div>
          <div className="rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
            <div className="mb-2 flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
              <BarChart3 className="h-5 w-5" />
            </div>
            <h3 className="font-semibold">{t("dashboard.engagement")}</h3>
            <p className="mt-2 text-2xl font-bold">
              {analytics.totalViews > 0
                ? ((analytics.totalClicks / analytics.totalViews) * 100).toFixed(1)
                : "0.0"}
              <span className="text-sm font-normal text-zinc-400">%</span>
            </p>
            <p className="mt-1 text-xs text-zinc-500">{t("dashboard.clickThroughRate")}</p>
          </div>
        </div>
      )}

      {/* Mini Views Trend + Score Breakdown */}
      {analytics && (
        <div className="mt-6 grid gap-6 sm:grid-cols-2">
          {/* Mini Views Trend Chart */}
          <div className="rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
            <h3 className="mb-4 text-sm font-semibold">{t("dashboard.viewsTrend")}</h3>
            {analytics.viewsByDay && analytics.viewsByDay.length > 0 ? (
              <ResponsiveContainer width="100%" height={120}>
                <AreaChart data={analytics.viewsByDay.slice(-7)}>
                  <defs>
                    <linearGradient id="viewsGradient" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#18181b" stopOpacity={0.3} />
                      <stop offset="95%" stopColor="#18181b" stopOpacity={0} />
                    </linearGradient>
                  </defs>
                  <XAxis
                    dataKey="date"
                    tick={{ fontSize: 10, fill: "#71717a" }}
                    tickFormatter={(d: string) => new Date(d).toLocaleDateString(undefined, { weekday: "short" })}
                    axisLine={false}
                    tickLine={false}
                  />
                  <Tooltip
                    labelFormatter={(d) => new Date(String(d)).toLocaleDateString()}
                    contentStyle={{ fontSize: 12, borderRadius: 8, border: "1px solid #e4e4e7" }}
                  />
                  <Area
                    type="monotone"
                    dataKey="count"
                    stroke="#18181b"
                    fill="url(#viewsGradient)"
                    strokeWidth={2}
                  />
                </AreaChart>
              </ResponsiveContainer>
            ) : (
              <p className="text-xs text-zinc-400">{t("analytics.noData")}</p>
            )}
          </div>

          {/* Creator Score Breakdown */}
          <div className="rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
            <h3 className="mb-4 text-sm font-semibold">{t("dashboard.scoreBreakdown")}</h3>
            <ScoreBreakdown user={user} connectionCount={connectionCount} connections={connections} />
          </div>
        </div>
      )}

      {/* Pending Collaborations + Quick Actions */}
      <div className="mt-6 grid gap-6 sm:grid-cols-2">
        {/* Pending Collaborations Card */}
        <div className="rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
          <div className="mb-2 flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
            <Inbox className="h-5 w-5" />
          </div>
          <h3 className="font-semibold">{t("dashboard.pendingCollabs")}</h3>
          <p className="mt-2 text-2xl font-bold">{pendingCollabCount}</p>
          <Link
            href="/collaborations"
            className="mt-3 inline-flex items-center gap-1 text-xs font-medium text-zinc-600 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100"
          >
            {t("dashboard.viewInbox")}
            <ArrowRight className="h-3 w-3" />
          </Link>
        </div>

        {/* Quick Actions */}
        <div className="rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
          <h3 className="mb-4 font-semibold">{t("dashboard.quickActions")}</h3>
          <div className="grid grid-cols-2 gap-3">
            <Link
              href="/vault/upload"
              className="flex items-center gap-2 rounded-lg border border-zinc-200 px-3 py-2.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
            >
              <Upload className="h-4 w-4 text-zinc-500" />
              {t("dashboard.uploadContent")}
            </Link>
            <Link
              href="/connections"
              className="flex items-center gap-2 rounded-lg border border-zinc-200 px-3 py-2.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
            >
              <Link2 className="h-4 w-4 text-zinc-500" />
              {t("dashboard.connectPlatform")}
            </Link>
            <button
              onClick={() => {
                if (profileUrl) {
                  navigator.clipboard.writeText(profileUrl);
                  setCopiedProfile(true);
                  setTimeout(() => setCopiedProfile(false), 2000);
                }
              }}
              className="flex items-center gap-2 rounded-lg border border-zinc-200 px-3 py-2.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
            >
              <Share2 className="h-4 w-4 text-zinc-500" />
              {copiedProfile ? t("dashboard.profileCopied") : t("dashboard.shareProfile")}
            </button>
            <Link
              href="/analytics"
              className="flex items-center gap-2 rounded-lg border border-zinc-200 px-3 py-2.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
            >
              <TrendingUp className="h-4 w-4 text-zinc-500" />
              {t("dashboard.viewAnalytics")}
            </Link>
          </div>
        </div>
      </div>

      {/* Recent Activity Feed */}
      <div className="mt-6 rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
        <div className="mb-4 flex items-center gap-2">
          <Activity className="h-4 w-4 text-zinc-500" />
          <h3 className="font-semibold">{t("dashboard.recentActivity")}</h3>
        </div>
        {recentNotifications.length > 0 ? (
          <ul className="space-y-3">
            {recentNotifications.slice(0, 5).map((n: any) => (
              <li key={n.id} className="flex items-start gap-3 text-sm">
                <div className={`mt-1 h-2 w-2 rounded-full ${n.readAt ? "bg-zinc-300 dark:bg-zinc-600" : "bg-zinc-900 dark:bg-zinc-100"}`} />
                <div>
                  <p className="text-zinc-700 dark:text-zinc-300">{n.title || n.message}</p>
                  <p className="text-xs text-zinc-400">{new Date(n.createdAt).toLocaleString()}</p>
                </div>
              </li>
            ))}
          </ul>
        ) : (
          <p className="text-sm text-zinc-400">{t("dashboard.noRecentActivity")}</p>
        )}
      </div>

      {/* Widget Card */}
      {user.username && (
        <div className="mt-6">
          <Link
            href="/widget"
            className="group flex items-center justify-between rounded-xl border border-zinc-200 p-4 transition-colors hover:border-zinc-300 dark:border-zinc-800 dark:hover:border-zinc-700"
          >
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
                <Code className="h-5 w-5" />
              </div>
              <div>
                <h3 className="text-sm font-semibold">{t("dashboard.getWidget")}</h3>
                <p className="text-xs text-zinc-500">{t("dashboard.getWidgetDesc")}</p>
              </div>
            </div>
            <ArrowRight className="h-4 w-4 text-zinc-400 transition-transform group-hover:translate-x-1" />
          </Link>
        </div>
      )}

      {/* Recommended Creators */}
      {recommendations.length > 0 && (
        <div className="mt-6">
          <h3 className="mb-4 text-lg font-semibold">{t("dashboard.recommendedCreators")}</h3>
          <div className="flex gap-4 overflow-x-auto pb-2">
            {recommendations.map((creator: any) => (
              <Link
                key={creator.id}
                href={`/profile?u=${creator.username}`}
                className="flex min-w-[200px] flex-col items-center rounded-xl border border-zinc-200 p-4 transition-colors hover:border-zinc-300 dark:border-zinc-800 dark:hover:border-zinc-700"
              >
                {creator.image ? (
                  <img src={creator.image} alt="" className="h-12 w-12 rounded-full object-cover" />
                ) : (
                  <div className="flex h-12 w-12 items-center justify-center rounded-full bg-zinc-100 text-lg font-bold dark:bg-zinc-800">
                    {(creator.name || "?")[0].toUpperCase()}
                  </div>
                )}
                <p className="mt-2 text-sm font-medium">{creator.name}</p>
                <p className="text-xs text-zinc-500">@{creator.username}</p>
                {creator.creatorTier && <TierBadge tier={creator.creatorTier} size="sm" />}
                <p className="mt-1 text-xs text-zinc-400">
                  {creator.sharedPlatforms} {t("dashboard.sharedPlatforms")}
                </p>
              </Link>
            ))}
          </div>
        </div>
      )}

      {/* Public Profile Link */}
      {user.username && (
        <div className="mt-8 rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
          <div className="flex items-start justify-between">
            <div>
              <h3 className="font-semibold">{t("dashboard.yourPublicProfile")}</h3>
              <p className="mt-1 text-sm text-zinc-500">
                {t("dashboard.shareLink")}
              </p>
            </div>
            <Link
              href={`/profile?u=${user.username}`}
              className="rounded-lg border border-zinc-200 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
            >
              {t("dashboard.preview")}
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
              {t("dashboard.qrCode")}
            </button>
          </div>
          {showQR && (
            <div className="mt-4 flex items-center gap-4">
              <div className="rounded-lg border border-zinc-200 bg-white p-3 dark:border-zinc-800 dark:bg-zinc-900">
                <QRCodeSVG value={profileUrl} size={120} level="M" marginSize={2} />
              </div>
              <div className="text-sm text-zinc-500">
                <p>{t("dashboard.scanQR")}</p>
                <p className="mt-1 text-xs text-zinc-400">
                  {t("dashboard.qrPerfectFor")}
                </p>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

function ScoreBreakdown({ user, connectionCount, connections }: { user: any; connectionCount: number; connections: any[] }) {
  const { t } = useTranslation();

  // Profile completeness: has avatar + bio = 20, one missing = 10, both missing = 0
  const hasAvatar = !!user.image;
  const hasBio = !!user.bio;
  const profilePts = (hasAvatar ? 10 : 0) + (hasBio ? 10 : 0);

  // Email verified: 10 pts
  const emailPts = user.emailVerified ? 10 : 0;

  // Connections: 10 per connection, max 50
  const connectionsPts = Math.min(connectionCount * 10, 50);

  // Followers bonus: logarithmic, max 20
  const totalFollowers = connections.reduce((sum: number, c: any) => sum + (c.followerCount || 0), 0);
  const followersPts = totalFollowers > 0 ? Math.min(Math.round(Math.log10(totalFollowers) * 5), 20) : 0;

  const bars = [
    { label: t("dashboard.profileScore"), value: profilePts, max: 20 },
    { label: t("dashboard.emailScore"), value: emailPts, max: 10 },
    { label: t("dashboard.connectionsScore"), value: connectionsPts, max: 50 },
    { label: t("dashboard.followersScore"), value: followersPts, max: 20 },
  ];

  return (
    <div className="space-y-3">
      {bars.map((bar) => (
        <div key={bar.label}>
          <div className="flex items-center justify-between text-xs">
            <span className="text-zinc-600 dark:text-zinc-400">{bar.label}</span>
            <span className="font-medium">{bar.value}/{bar.max}</span>
          </div>
          <div className="mt-1 h-2 rounded-full bg-zinc-100 dark:bg-zinc-800">
            <div
              className="h-2 rounded-full bg-zinc-900 transition-all dark:bg-zinc-100"
              style={{ width: `${(bar.value / bar.max) * 100}%` }}
            />
          </div>
        </div>
      ))}
    </div>
  );
}

export default function DashboardPage() {
  return (
    <Suspense fallback={null}>
      <DashboardContent />
    </Suspense>
  );
}
