"use client";

import { useAuth } from "@/lib/auth-context";
import { api } from "@/lib/api";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState, Suspense } from "react";
import { useTranslation } from "react-i18next";

function BillingContent() {
  const { user, loading: authLoading } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();
  const { t } = useTranslation();

  const [plan, setPlan] = useState<string>("free");
  const [status, setStatus] = useState<string>("active");
  const [periodEnd, setPeriodEnd] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [checkoutLoading, setCheckoutLoading] = useState<string | null>(null);
  const [portalLoading, setPortalLoading] = useState(false);

  const success = searchParams.get("success") === "true";
  const canceled = searchParams.get("canceled") === "true";

  useEffect(() => {
    if (!authLoading && !user) router.push("/sign-in");
  }, [user, authLoading, router]);

  useEffect(() => {
    if (!user) return;
    (async () => {
      const res = await api.billing.subscription();
      if (res.data) {
        setPlan(res.data.plan);
        setStatus(res.data.status);
        setPeriodEnd(res.data.currentPeriodEnd);
      }
      setLoading(false);
    })();
  }, [user]);

  const handleCheckout = async (targetPlan: string) => {
    setCheckoutLoading(targetPlan);
    const res = await api.billing.checkout(targetPlan);
    setCheckoutLoading(null);
    if (res.data?.url) {
      window.location.href = res.data.url;
    }
  };

  const handlePortal = async () => {
    setPortalLoading(true);
    const res = await api.billing.portal();
    setPortalLoading(false);
    if (res.data?.url) {
      window.location.href = res.data.url;
    }
  };

  if (authLoading || !user) return null;

  const planLabel =
    plan === "business"
      ? t("billing.businessPlan")
      : plan === "pro"
        ? t("billing.proPlan")
        : t("billing.freePlan");

  const planFeatures: Record<string, string[]> = {
    free: [
      "3 social connections",
      "Basic profile",
      "Creator Score",
      "Public profile page",
    ],
    pro: [
      "Unlimited connections",
      "Advanced analytics",
      "Custom themes",
      "API keys (1,000 req/mo)",
      "Embeddable widget",
      "Priority support",
    ],
    business: [
      "Everything in Pro",
      "Bulk verification API (10K req/mo)",
      "Brand dashboard",
      "Saved creator lists",
      "White-label embed",
    ],
  };

  const features = planFeatures[plan] || planFeatures.free;

  return (
    <div className="mx-auto max-w-2xl px-4 py-12 sm:px-6">
      <h1 className="text-2xl font-bold">{t("billing.title")}</h1>
      <p className="mt-1 text-sm text-zinc-500">{t("billing.subtitle")}</p>

      {/* Success / Canceled banners */}
      {success && (
        <div className="mt-6 rounded-lg border border-green-200 bg-green-50 p-4 dark:border-green-900 dark:bg-green-950">
          <p className="font-medium text-green-800 dark:text-green-200">
            {t("billing.successTitle")}
          </p>
          <p className="mt-1 text-sm text-green-700 dark:text-green-300">
            {t("billing.successDesc")}
          </p>
        </div>
      )}
      {canceled && (
        <div className="mt-6 rounded-lg border border-yellow-200 bg-yellow-50 p-4 dark:border-yellow-900 dark:bg-yellow-950">
          <p className="font-medium text-yellow-800 dark:text-yellow-200">
            {t("billing.canceledTitle")}
          </p>
          <p className="mt-1 text-sm text-yellow-700 dark:text-yellow-300">
            {t("billing.canceledDesc")}
          </p>
        </div>
      )}

      {/* Current plan card */}
      <div className="mt-8 rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
        {loading ? (
          <p className="text-sm text-zinc-500">{t("common.loading")}</p>
        ) : (
          <>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-zinc-500">
                  {t("billing.currentPlan")}
                </p>
                <div className="mt-1 flex items-center gap-2">
                  <span className="text-xl font-bold">{planLabel}</span>
                  {plan !== "free" && (
                    <span
                      className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                        status === "active"
                          ? "bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300"
                          : status === "past_due"
                            ? "bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300"
                            : "bg-zinc-100 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400"
                      }`}
                    >
                      {status}
                    </span>
                  )}
                </div>
              </div>
              {plan !== "free" && (
                <button
                  onClick={handlePortal}
                  disabled={portalLoading}
                  className="rounded-lg border border-zinc-300 px-4 py-2 text-sm font-medium transition hover:bg-zinc-50 disabled:opacity-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
                >
                  {portalLoading ? "..." : t("billing.manage")}
                </button>
              )}
            </div>

            {periodEnd && plan !== "free" && (
              <p className="mt-3 text-xs text-zinc-500">
                Current period ends:{" "}
                {new Date(periodEnd).toLocaleDateString()}
              </p>
            )}

            {/* Plan features */}
            <div className="mt-6">
              <p className="text-sm font-medium text-zinc-600 dark:text-zinc-400">
                {t("billing.planFeatures")}
              </p>
              <ul className="mt-3 space-y-2">
                {features.map((feature, i) => (
                  <li
                    key={i}
                    className="flex items-center gap-2 text-sm text-zinc-600 dark:text-zinc-400"
                  >
                    <svg
                      className="h-4 w-4 shrink-0 text-zinc-400"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                      strokeWidth={2}
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M5 13l4 4L19 7"
                      />
                    </svg>
                    {feature}
                  </li>
                ))}
              </ul>
            </div>

            {/* Upgrade buttons for free plan */}
            {plan === "free" && (
              <div className="mt-6 flex gap-3">
                <button
                  onClick={() => handleCheckout("pro")}
                  disabled={checkoutLoading === "pro"}
                  className="flex-1 rounded-lg bg-zinc-900 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-zinc-800 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
                >
                  {checkoutLoading === "pro"
                    ? "..."
                    : `${t("billing.upgrade")} to Pro`}
                </button>
                <button
                  onClick={() => handleCheckout("business")}
                  disabled={checkoutLoading === "business"}
                  className="flex-1 rounded-lg border border-zinc-300 px-4 py-2.5 text-sm font-medium transition hover:bg-zinc-50 disabled:opacity-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
                >
                  {checkoutLoading === "business"
                    ? "..."
                    : `${t("billing.upgrade")} to Business`}
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}

export default function BillingPage() {
  return (
    <Suspense fallback={null}>
      <BillingContent />
    </Suspense>
  );
}
