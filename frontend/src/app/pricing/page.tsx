"use client";

import { useAuth } from "@/lib/auth-context";
import { api } from "@/lib/api";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useTranslation } from "react-i18next";

export default function PricingPage() {
  const { user } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();
  const [loading, setLoading] = useState<string | null>(null);

  const handleCheckout = async (plan: string) => {
    if (!user) {
      router.push("/sign-in");
      return;
    }
    setLoading(plan);
    const res = await api.billing.checkout(plan);
    setLoading(null);
    if (res.data?.url) {
      window.location.href = res.data.url;
    }
  };

  const freeFeatures: string[] = t("pricing.freeFeatures", {
    returnObjects: true,
  }) as unknown as string[];
  const proFeatures: string[] = t("pricing.proFeatures", {
    returnObjects: true,
  }) as unknown as string[];
  const businessFeatures: string[] = t("pricing.businessFeatures", {
    returnObjects: true,
  }) as unknown as string[];

  return (
    <div className="mx-auto max-w-5xl px-4 py-16 sm:px-6">
      <div className="text-center">
        <h1 className="text-3xl font-bold tracking-tight sm:text-4xl">
          {t("pricing.title")}
        </h1>
        <p className="mt-3 text-lg text-zinc-500 dark:text-zinc-400">
          {t("pricing.subtitle")}
        </p>
      </div>

      <div className="mt-12 grid gap-6 sm:grid-cols-3">
        {/* Free Plan */}
        <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
          <h3 className="text-lg font-semibold">{t("pricing.free")}</h3>
          <div className="mt-3">
            <span className="text-3xl font-bold">{t("pricing.freePrice")}</span>
            <span className="text-sm text-zinc-500">{t("pricing.month")}</span>
          </div>
          <ul className="mt-6 space-y-3">
            {freeFeatures.map((feature: string, i: number) => (
              <li key={i} className="flex items-start gap-2 text-sm">
                <svg
                  className="mt-0.5 h-4 w-4 shrink-0 text-zinc-400"
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
                <span className="text-zinc-600 dark:text-zinc-400">
                  {feature}
                </span>
              </li>
            ))}
          </ul>
          <button
            onClick={() => (user ? router.push("/dashboard") : router.push("/sign-in"))}
            className="mt-8 w-full rounded-lg border border-zinc-300 px-4 py-2.5 text-sm font-medium transition hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
          >
            {t("pricing.getStarted")}
          </button>
        </div>

        {/* Pro Plan */}
        <div className="relative rounded-xl border-2 border-zinc-900 bg-white p-6 dark:border-zinc-100 dark:bg-zinc-900">
          <div className="absolute -top-3 left-1/2 -translate-x-1/2">
            <span className="rounded-full bg-zinc-900 px-3 py-1 text-xs font-medium text-white dark:bg-zinc-100 dark:text-zinc-900">
              {t("pricing.mostPopular")}
            </span>
          </div>
          <h3 className="text-lg font-semibold">{t("pricing.pro")}</h3>
          <div className="mt-3">
            <span className="text-3xl font-bold">{t("pricing.proPrice")}</span>
            <span className="text-sm text-zinc-500">{t("pricing.month")}</span>
          </div>
          <ul className="mt-6 space-y-3">
            {proFeatures.map((feature: string, i: number) => (
              <li key={i} className="flex items-start gap-2 text-sm">
                <svg
                  className="mt-0.5 h-4 w-4 shrink-0 text-zinc-900 dark:text-zinc-100"
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
                <span className="text-zinc-600 dark:text-zinc-400">
                  {feature}
                </span>
              </li>
            ))}
          </ul>
          <button
            onClick={() => handleCheckout("pro")}
            disabled={loading === "pro"}
            className="mt-8 w-full rounded-lg bg-zinc-900 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-zinc-800 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
          >
            {loading === "pro" ? "..." : t("pricing.upgradeToPro")}
          </button>
        </div>

        {/* Business Plan */}
        <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
          <h3 className="text-lg font-semibold">{t("pricing.business")}</h3>
          <div className="mt-3">
            <span className="text-3xl font-bold">
              {t("pricing.businessPrice")}
            </span>
            <span className="text-sm text-zinc-500">{t("pricing.month")}</span>
          </div>
          <ul className="mt-6 space-y-3">
            {businessFeatures.map((feature: string, i: number) => (
              <li key={i} className="flex items-start gap-2 text-sm">
                <svg
                  className="mt-0.5 h-4 w-4 shrink-0 text-zinc-400"
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
                <span className="text-zinc-600 dark:text-zinc-400">
                  {feature}
                </span>
              </li>
            ))}
          </ul>
          <a
            href="mailto:sales@creatrid.com"
            className="mt-8 block w-full rounded-lg border border-zinc-300 px-4 py-2.5 text-center text-sm font-medium transition hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
          >
            {t("pricing.contactSales")}
          </a>
        </div>
      </div>
    </div>
  );
}
