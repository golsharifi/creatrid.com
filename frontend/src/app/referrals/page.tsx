"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import { api } from "@/lib/api";
import { useTranslation } from "react-i18next";
import { Gift, Copy, Check, Users, Trophy } from "lucide-react";

export default function ReferralsPage() {
  const { user, loading: authLoading } = useAuth();
  const { t } = useTranslation();
  const [code, setCode] = useState("");
  const [referrals, setReferrals] = useState<any[]>([]);
  const [stats, setStats] = useState<{ totalReferred: number; totalBonusEarned: number } | null>(null);
  const [loading, setLoading] = useState(true);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    if (!user) return;
    (async () => {
      const [codeRes, listRes] = await Promise.all([
        api.referrals.code(),
        api.referrals.list(),
      ]);
      if (codeRes.data) setCode(codeRes.data.code);
      if (listRes.data) {
        setReferrals(listRes.data.referrals || []);
        setStats(listRes.data.stats);
      }
      setLoading(false);
    })();
  }, [user]);

  const referralLink = typeof window !== "undefined"
    ? `${window.location.origin}/sign-in?ref=${code}`
    : "";

  const handleCopy = () => {
    navigator.clipboard.writeText(referralLink);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  if (authLoading || loading) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-2 border-zinc-300 border-t-zinc-900 dark:border-zinc-700 dark:border-t-zinc-100" />
      </div>
    );
  }

  if (!user) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-16 text-center">
        <p className="text-zinc-500">{t("common.signInRequired")}</p>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-4xl px-4 py-8 sm:px-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">{t("referrals.title")}</h1>
        <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">{t("referrals.subtitle")}</p>
      </div>

      {/* Referral Link */}
      <div className="mb-8 rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
        <h3 className="mb-4 text-lg font-semibold text-zinc-900 dark:text-zinc-100">{t("referrals.yourLink")}</h3>
        <div className="flex gap-2">
          <input
            type="text"
            readOnly
            value={referralLink}
            className="flex-1 rounded-lg border border-zinc-300 bg-zinc-50 px-3 py-2 text-sm text-zinc-700 dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-300"
          />
          <button
            onClick={handleCopy}
            className="flex items-center gap-2 rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white hover:bg-zinc-800 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
          >
            {copied ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
            {copied ? t("common.copied") : t("common.copy")}
          </button>
        </div>
        <p className="mt-2 text-xs text-zinc-400">{t("referrals.linkHint")}</p>
      </div>

      {/* Stats */}
      <div className="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
          <div className="flex items-center gap-3">
            <Users className="h-8 w-8 text-blue-500" />
            <div>
              <p className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">{stats?.totalReferred || 0}</p>
              <p className="text-sm text-zinc-500 dark:text-zinc-400">{t("referrals.totalReferred")}</p>
            </div>
          </div>
        </div>
        <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
          <div className="flex items-center gap-3">
            <Trophy className="h-8 w-8 text-amber-500" />
            <div>
              <p className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">+{stats?.totalBonusEarned || 0}</p>
              <p className="text-sm text-zinc-500 dark:text-zinc-400">{t("referrals.bonusEarned")}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Referral List */}
      <div className="rounded-xl border border-zinc-200 bg-white dark:border-zinc-800 dark:bg-zinc-900">
        <div className="border-b border-zinc-200 px-6 py-4 dark:border-zinc-800">
          <h3 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">{t("referrals.referred")}</h3>
        </div>
        {referrals.length === 0 ? (
          <div className="p-12 text-center">
            <Gift className="mx-auto mb-4 h-12 w-12 text-zinc-300 dark:text-zinc-600" />
            <p className="text-zinc-500 dark:text-zinc-400">{t("referrals.empty")}</p>
          </div>
        ) : (
          <div className="divide-y divide-zinc-200 dark:divide-zinc-800">
            {referrals.map((r: any) => (
              <div key={r.id} className="flex items-center justify-between px-6 py-4">
                <span className="text-sm text-zinc-700 dark:text-zinc-300">{r.referredUserName || "User"}</span>
                <div className="flex items-center gap-4">
                  <span className="text-sm font-medium text-emerald-600 dark:text-emerald-400">+{r.rewardValue} score</span>
                  <span className="text-xs text-zinc-400">{new Date(r.createdAt).toLocaleDateString()}</span>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
