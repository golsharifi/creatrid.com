"use client";

import { useState } from "react";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { useTranslation } from "react-i18next";
import {
  Shield,
  Users,
  CheckCircle,
  ArrowRight,
  BarChart3,
  Globe,
  Lock,
  Eye,
  Link2,
  Zap,
  ChevronDown,
} from "@/components/icons";

export default function Home() {
  const { user } = useAuth();
  const { t } = useTranslation();
  const [openFaq, setOpenFaq] = useState<number | null>(null);

  const toggleFaq = (index: number) => {
    setOpenFaq(openFaq === index ? null : index);
  };

  return (
    <div className="flex flex-col">
      {/* ── Hero ── */}
      <section className="mx-auto flex w-full max-w-6xl flex-col items-center gap-6 px-4 py-20 text-center sm:px-6 sm:py-32">
        <div className="inline-flex items-center gap-2 rounded-full border border-zinc-200 bg-zinc-50 px-4 py-1.5 text-sm text-zinc-600 dark:border-zinc-800 dark:bg-zinc-900 dark:text-zinc-400">
          <Shield className="h-4 w-4" />
          {t("landing.heroBadge")}
        </div>

        <h1 className="max-w-4xl text-4xl font-bold tracking-tight sm:text-5xl lg:text-6xl">
          {t("landing.heroHeadline")}
          <br />
          <span className="text-zinc-400 dark:text-zinc-500">
            {t("landing.heroHeadlineAccent")}
          </span>
        </h1>

        <p className="max-w-2xl text-lg text-zinc-600 dark:text-zinc-400">
          {t("landing.heroSubtext")}
        </p>

        <div className="flex flex-col gap-3 sm:flex-row">
          {user ? (
            <Link
              href="/dashboard"
              className="inline-flex items-center gap-2 rounded-lg bg-zinc-900 px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-zinc-800 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              {t("landing.ctaDashboard")}
              <ArrowRight className="h-4 w-4" />
            </Link>
          ) : (
            <Link
              href="/sign-in"
              className="inline-flex items-center gap-2 rounded-lg bg-zinc-900 px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-zinc-800 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              {t("landing.ctaGetPassport")}
              <ArrowRight className="h-4 w-4" />
            </Link>
          )}
          <a
            href="#how-it-works"
            className="inline-flex items-center gap-2 rounded-lg border border-zinc-200 px-6 py-3 text-sm font-medium text-zinc-700 transition-colors hover:bg-zinc-50 dark:border-zinc-800 dark:text-zinc-300 dark:hover:bg-zinc-900"
          >
            {t("landing.ctaHowItWorks")}
          </a>
        </div>
      </section>

      {/* ── Social Proof Bar ── */}
      <section className="border-y border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/50">
        <div className="mx-auto flex max-w-6xl flex-col items-center gap-6 px-4 py-10 sm:px-6 md:flex-row md:justify-between">
          <p className="text-sm font-medium text-zinc-500 dark:text-zinc-400">
            {t("landing.socialProofLabel")}
          </p>
          <div className="flex flex-wrap items-center justify-center gap-8 md:gap-12">
            <div className="flex flex-col items-center gap-1">
              <span className="text-2xl font-bold">{t("landing.statCreators")}</span>
              <span className="text-xs text-zinc-500 dark:text-zinc-400">
                {t("landing.statCreatorsLabel")}
              </span>
            </div>
            <div className="flex flex-col items-center gap-1">
              <span className="text-2xl font-bold">{t("landing.statPlatforms")}</span>
              <span className="text-xs text-zinc-500 dark:text-zinc-400">
                {t("landing.statPlatformsLabel")}
              </span>
            </div>
            <div className="flex flex-col items-center gap-1">
              <span className="text-2xl font-bold">{t("landing.statConnections")}</span>
              <span className="text-xs text-zinc-500 dark:text-zinc-400">
                {t("landing.statConnectionsLabel")}
              </span>
            </div>
          </div>
        </div>
      </section>

      {/* ── How it Works ── */}
      <section id="how-it-works" className="mx-auto w-full max-w-6xl px-4 py-20 sm:px-6 sm:py-28">
        <div className="mb-14 text-center">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
            {t("landing.howTitle")}
          </h2>
          <p className="mx-auto mt-4 max-w-2xl text-zinc-600 dark:text-zinc-400">
            {t("landing.howSubtitle")}
          </p>
        </div>

        <div className="grid gap-10 sm:grid-cols-3">
          {/* Step 1 */}
          <div className="flex flex-col items-center text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-zinc-900 text-lg font-bold text-white dark:bg-zinc-100 dark:text-zinc-900">
              1
            </div>
            <h3 className="mt-5 text-lg font-semibold">{t("landing.step1Title")}</h3>
            <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
              {t("landing.step1Desc")}
            </p>
          </div>

          {/* Step 2 */}
          <div className="flex flex-col items-center text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-zinc-900 text-lg font-bold text-white dark:bg-zinc-100 dark:text-zinc-900">
              2
            </div>
            <h3 className="mt-5 text-lg font-semibold">{t("landing.step2Title")}</h3>
            <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
              {t("landing.step2Desc")}
            </p>
          </div>

          {/* Step 3 */}
          <div className="flex flex-col items-center text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-zinc-900 text-lg font-bold text-white dark:bg-zinc-100 dark:text-zinc-900">
              3
            </div>
            <h3 className="mt-5 text-lg font-semibold">{t("landing.step3Title")}</h3>
            <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
              {t("landing.step3Desc")}
            </p>
          </div>
        </div>
      </section>

      {/* ── Features Grid ── */}
      <section className="border-t border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/50">
        <div className="mx-auto max-w-6xl px-4 py-20 sm:px-6 sm:py-28">
          <div className="mb-14 text-center">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
              {t("landing.featuresTitle")}
            </h2>
            <p className="mx-auto mt-4 max-w-2xl text-zinc-600 dark:text-zinc-400">
              {t("landing.featuresSubtitle")}
            </p>
          </div>

          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {/* Verified Identity */}
            <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900">
                <Shield className="h-5 w-5" />
              </div>
              <h3 className="mt-4 font-semibold">{t("landing.feature1Title")}</h3>
              <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                {t("landing.feature1Desc")}
              </p>
            </div>

            {/* Creator Score */}
            <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900">
                <Zap className="h-5 w-5" />
              </div>
              <h3 className="mt-4 font-semibold">{t("landing.feature2Title")}</h3>
              <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                {t("landing.feature2Desc")}
              </p>
            </div>

            {/* Cross-Platform Profiles */}
            <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900">
                <Globe className="h-5 w-5" />
              </div>
              <h3 className="mt-4 font-semibold">{t("landing.feature3Title")}</h3>
              <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                {t("landing.feature3Desc")}
              </p>
            </div>

            {/* Brand Discovery */}
            <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900">
                <Eye className="h-5 w-5" />
              </div>
              <h3 className="mt-4 font-semibold">{t("landing.feature4Title")}</h3>
              <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                {t("landing.feature4Desc")}
              </p>
            </div>

            {/* Real-time Analytics */}
            <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900">
                <BarChart3 className="h-5 w-5" />
              </div>
              <h3 className="mt-4 font-semibold">{t("landing.feature5Title")}</h3>
              <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                {t("landing.feature5Desc")}
              </p>
            </div>

            {/* Privacy First */}
            <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900">
                <Lock className="h-5 w-5" />
              </div>
              <h3 className="mt-4 font-semibold">{t("landing.feature6Title")}</h3>
              <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                {t("landing.feature6Desc")}
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* ── Pricing Preview ── */}
      <section className="mx-auto w-full max-w-6xl px-4 py-20 sm:px-6 sm:py-28">
        <div className="mb-14 text-center">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
            {t("landing.pricingTitle")}
          </h2>
          <p className="mx-auto mt-4 max-w-2xl text-zinc-600 dark:text-zinc-400">
            {t("landing.pricingSubtitle")}
          </p>
        </div>

        <div className="grid gap-6 sm:grid-cols-3">
          <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
            <h3 className="text-lg font-semibold">{t("pricing.free")}</h3>
            <div className="mt-2">
              <span className="text-3xl font-bold">{t("pricing.freePrice")}</span>
              <span className="text-sm text-zinc-500">{t("pricing.month")}</span>
            </div>
            <p className="mt-3 text-sm text-zinc-500 dark:text-zinc-400">
              {t("landing.pricingFreeDesc")}
            </p>
          </div>

          <div className="relative rounded-xl border-2 border-zinc-900 bg-white p-6 dark:border-zinc-100 dark:bg-zinc-950">
            <div className="absolute -top-3 left-1/2 -translate-x-1/2">
              <span className="rounded-full bg-zinc-900 px-3 py-1 text-xs font-medium text-white dark:bg-zinc-100 dark:text-zinc-900">
                {t("pricing.mostPopular")}
              </span>
            </div>
            <h3 className="text-lg font-semibold">{t("pricing.pro")}</h3>
            <div className="mt-2">
              <span className="text-3xl font-bold">{t("pricing.proPrice")}</span>
              <span className="text-sm text-zinc-500">{t("pricing.month")}</span>
            </div>
            <p className="mt-3 text-sm text-zinc-500 dark:text-zinc-400">
              {t("landing.pricingProDesc")}
            </p>
          </div>

          <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
            <h3 className="text-lg font-semibold">{t("pricing.business")}</h3>
            <div className="mt-2">
              <span className="text-3xl font-bold">{t("pricing.businessPrice")}</span>
              <span className="text-sm text-zinc-500">{t("pricing.month")}</span>
            </div>
            <p className="mt-3 text-sm text-zinc-500 dark:text-zinc-400">
              {t("landing.pricingBusinessDesc")}
            </p>
          </div>
        </div>

        <div className="mt-8 text-center">
          <Link
            href="/pricing"
            className="inline-flex items-center gap-2 text-sm font-medium text-zinc-600 transition-colors hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100"
          >
            {t("landing.pricingViewAll")}
            <ArrowRight className="h-4 w-4" />
          </Link>
        </div>
      </section>

      {/* ── Testimonials ── */}
      <section className="border-t border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/50">
        <div className="mx-auto max-w-6xl px-4 py-20 sm:px-6 sm:py-28">
          <div className="mb-14 text-center">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
              {t("landing.testimonialsTitle")}
            </h2>
          </div>

          <div className="grid gap-6 sm:grid-cols-3">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950"
              >
                <p className="text-sm text-zinc-600 dark:text-zinc-400">
                  &ldquo;{t(`landing.testimonial${i}Quote`)}&rdquo;
                </p>
                <div className="mt-4 flex items-center gap-3">
                  <div className="flex h-9 w-9 items-center justify-center rounded-full bg-zinc-100 text-sm font-bold text-zinc-500 dark:bg-zinc-800">
                    {t(`landing.testimonial${i}Name`)[0]}
                  </div>
                  <div>
                    <p className="text-sm font-medium">
                      {t(`landing.testimonial${i}Name`)}
                    </p>
                    <p className="text-xs text-zinc-500">
                      {t(`landing.testimonial${i}Role`)}
                    </p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Platform Logos ── */}
      <section className="mx-auto w-full max-w-6xl px-4 py-20 sm:px-6">
        <div className="mb-10 text-center">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
            {t("landing.platformsTitle")}
          </h2>
          <p className="mx-auto mt-4 max-w-2xl text-zinc-600 dark:text-zinc-400">
            {t("landing.platformsSubtitle")}
          </p>
        </div>

        <div className="flex flex-wrap items-center justify-center gap-6 sm:gap-8">
          {[
            { abbr: "YT", name: "YouTube" },
            { abbr: "GH", name: "GitHub" },
            { abbr: "\ud835\udd4f", name: "Twitter / X" },
            { abbr: "in", name: "LinkedIn" },
            { abbr: "IG", name: "Instagram" },
          ].map((platform) => (
            <div
              key={platform.name}
              className="flex flex-col items-center gap-2"
            >
              <div className="flex h-14 w-14 items-center justify-center rounded-xl border border-zinc-200 bg-white text-lg font-bold text-zinc-700 dark:border-zinc-800 dark:bg-zinc-950 dark:text-zinc-300">
                {platform.abbr}
              </div>
              <span className="text-xs text-zinc-500 dark:text-zinc-400">
                {platform.name}
              </span>
            </div>
          ))}
        </div>
      </section>

      {/* ── FAQ ── */}
      <section className="border-t border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/50">
        <div className="mx-auto max-w-3xl px-4 py-20 sm:px-6 sm:py-28">
          <div className="mb-14 text-center">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
              {t("landing.faqTitle")}
            </h2>
            <p className="mx-auto mt-4 max-w-2xl text-zinc-600 dark:text-zinc-400">
              {t("landing.faqSubtitle")}
            </p>
          </div>

          <div className="flex flex-col gap-3">
            {[
              { q: t("landing.faq1Q"), a: t("landing.faq1A") },
              { q: t("landing.faq2Q"), a: t("landing.faq2A") },
              { q: t("landing.faq3Q"), a: t("landing.faq3A") },
              { q: t("landing.faq4Q"), a: t("landing.faq4A") },
            ].map((item, index) => (
              <div
                key={index}
                className="rounded-xl border border-zinc-200 bg-white dark:border-zinc-800 dark:bg-zinc-950"
              >
                <button
                  onClick={() => toggleFaq(index)}
                  className="flex w-full items-center justify-between px-6 py-4 text-left"
                >
                  <span className="font-medium">{item.q}</span>
                  <ChevronDown
                    className={`h-5 w-5 shrink-0 text-zinc-400 transition-transform duration-200 ${
                      openFaq === index ? "rotate-180" : ""
                    }`}
                  />
                </button>
                {openFaq === index && (
                  <div className="border-t border-zinc-100 px-6 py-4 text-sm text-zinc-600 dark:border-zinc-800 dark:text-zinc-400">
                    {item.a}
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Final CTA ── */}
      <section className="mx-auto w-full max-w-6xl px-4 py-20 text-center sm:px-6 sm:py-28">
        <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
          {t("landing.ctaFinalHeadline")}
        </h2>
        <p className="mx-auto mt-4 max-w-xl text-zinc-600 dark:text-zinc-400">
          {t("landing.ctaFinalSubtext")}
        </p>
        <div className="mt-8">
          {user ? (
            <Link
              href="/dashboard"
              className="inline-flex items-center gap-2 rounded-lg bg-zinc-900 px-8 py-3.5 text-sm font-medium text-white transition-colors hover:bg-zinc-800 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              {t("landing.ctaDashboard")}
              <ArrowRight className="h-4 w-4" />
            </Link>
          ) : (
            <Link
              href="/sign-in"
              className="inline-flex items-center gap-2 rounded-lg bg-zinc-900 px-8 py-3.5 text-sm font-medium text-white transition-colors hover:bg-zinc-800 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              {t("landing.ctaGetPassport")}
              <ArrowRight className="h-4 w-4" />
            </Link>
          )}
        </div>
      </section>
    </div>
  );
}
