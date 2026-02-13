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
  Archive,
  ShoppingBag,
  Code,
  CreditCard,
  Gift,
  DollarSign,
  Key,
  Webhook,
  FolderOpen,
  TrendingUp,
} from "@/components/icons";
import { blogPosts } from "@/app/blog/data";

type IconComponent = React.ComponentType<{ className?: string }>;

export default function Home() {
  const { user } = useAuth();
  const { t } = useTranslation();
  const [openFaq, setOpenFaq] = useState<number | null>(null);
  const [activeTab, setActiveTab] = useState(0);

  const toggleFaq = (index: number) => {
    setOpenFaq(openFaq === index ? null : index);
  };

  const featureTabs: { label: string; features: { icon: IconComponent; title: string; desc: string }[] }[] = [
    {
      label: t("landing.tabIdentity"),
      features: [
        { icon: Shield, title: t("landing.featIdentityTitle"), desc: t("landing.featIdentityDesc") },
        { icon: Zap, title: t("landing.featScoreTitle"), desc: t("landing.featScoreDesc") },
        { icon: Lock, title: t("landing.feat2faTitle"), desc: t("landing.feat2faDesc") },
        { icon: Code, title: t("landing.featWidgetTitle"), desc: t("landing.featWidgetDesc") },
      ],
    },
    {
      label: t("landing.tabContent"),
      features: [
        { icon: Archive, title: t("landing.featVaultTitle"), desc: t("landing.featVaultDesc") },
        { icon: Link2, title: t("landing.featBlockchainTitle"), desc: t("landing.featBlockchainDesc") },
        { icon: CheckCircle, title: t("landing.featLicensingTitle"), desc: t("landing.featLicensingDesc") },
        { icon: FolderOpen, title: t("landing.featCollectionsTitle"), desc: t("landing.featCollectionsDesc") },
      ],
    },
    {
      label: t("landing.tabMonetization"),
      features: [
        { icon: ShoppingBag, title: t("landing.featMarketplaceTitle"), desc: t("landing.featMarketplaceDesc") },
        { icon: CreditCard, title: t("landing.featTokensTitle"), desc: t("landing.featTokensDesc") },
        { icon: Gift, title: t("landing.featSubscriptionsTitle"), desc: t("landing.featSubscriptionsDesc") },
        { icon: DollarSign, title: t("landing.featEarningsTitle"), desc: t("landing.featEarningsDesc") },
      ],
    },
    {
      label: t("landing.tabDiscovery"),
      features: [
        { icon: Eye, title: t("landing.featDiscoveryTitle"), desc: t("landing.featDiscoveryDesc") },
        { icon: Users, title: t("landing.featCollabTitle"), desc: t("landing.featCollabDesc") },
        { icon: TrendingUp, title: t("landing.featReferralsTitle"), desc: t("landing.featReferralsDesc") },
        { icon: Globe, title: t("landing.featAgencyTitle"), desc: t("landing.featAgencyDesc") },
      ],
    },
    {
      label: t("landing.tabAnalytics"),
      features: [
        { icon: BarChart3, title: t("landing.featProfileAnalyticsTitle"), desc: t("landing.featProfileAnalyticsDesc") },
        { icon: BarChart3, title: t("landing.featContentAnalyticsTitle"), desc: t("landing.featContentAnalyticsDesc") },
        { icon: Key, title: t("landing.featApiTitle"), desc: t("landing.featApiDesc") },
        { icon: Webhook, title: t("landing.featWebhooksTitle"), desc: t("landing.featWebhooksDesc") },
      ],
    },
  ];

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
            {[
              { value: t("landing.statCreators"), label: t("landing.statCreatorsLabel") },
              { value: t("landing.statPlatforms"), label: t("landing.statPlatformsLabel") },
              { value: t("landing.statConnections"), label: t("landing.statConnectionsLabel") },
              { value: t("landing.statContent"), label: t("landing.statContentLabel") },
              { value: t("landing.statAnchors"), label: t("landing.statAnchorsLabel") },
            ].map((stat) => (
              <div key={stat.label} className="flex flex-col items-center gap-1">
                <span className="text-2xl font-bold">{stat.value}</span>
                <span className="text-xs text-zinc-500 dark:text-zinc-400">
                  {stat.label}
                </span>
              </div>
            ))}
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

        <div className="grid gap-10 sm:grid-cols-2 lg:grid-cols-4">
          {[1, 2, 3, 4].map((step) => (
            <div key={step} className="flex flex-col items-center text-center">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-zinc-900 text-lg font-bold text-white dark:bg-zinc-100 dark:text-zinc-900">
                {step}
              </div>
              <h3 className="mt-5 text-lg font-semibold">{t(`landing.step${step}Title`)}</h3>
              <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                {t(`landing.step${step}Desc`)}
              </p>
            </div>
          ))}
        </div>
      </section>

      {/* ── Tabbed Features ── */}
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

          {/* Tab buttons */}
          <div className="mb-10 flex flex-wrap justify-center gap-2">
            {featureTabs.map((tab, i) => (
              <button
                key={i}
                onClick={() => setActiveTab(i)}
                className={`rounded-full px-4 py-2 text-sm font-medium transition-colors ${
                  activeTab === i
                    ? "bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900"
                    : "bg-zinc-100 text-zinc-600 hover:bg-zinc-200 dark:bg-zinc-800 dark:text-zinc-400 dark:hover:bg-zinc-700"
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>

          {/* Feature cards */}
          <div className="grid gap-6 sm:grid-cols-2">
            {featureTabs[activeTab].features.map((feature, i) => {
              const Icon = feature.icon;
              return (
                <div
                  key={i}
                  className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950"
                >
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900">
                    <Icon className="h-5 w-5" />
                  </div>
                  <h3 className="mt-4 font-semibold">{feature.title}</h3>
                  <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                    {feature.desc}
                  </p>
                </div>
              );
            })}
          </div>
        </div>
      </section>

      {/* ── Blockchain Highlight ── */}
      <section className="bg-zinc-900 dark:bg-zinc-950">
        <div className="mx-auto max-w-6xl px-4 py-20 sm:px-6 sm:py-28">
          <div className="flex flex-col items-center text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-zinc-800 dark:bg-zinc-900">
              <Link2 className="h-6 w-6 text-white" />
            </div>
            <h2 className="mt-6 text-3xl font-bold tracking-tight text-white sm:text-4xl">
              {t("landing.blockchainTitle")}
            </h2>
            <p className="mx-auto mt-4 max-w-2xl text-zinc-400">
              {t("landing.blockchainSubtitle")}
            </p>
            <div className="mt-10 grid w-full max-w-2xl gap-6 sm:grid-cols-3">
              {[
                { stat: t("landing.blockchainStat1"), desc: t("landing.blockchainStat1Desc") },
                { stat: t("landing.blockchainStat2"), desc: t("landing.blockchainStat2Desc") },
                { stat: t("landing.blockchainStat3"), desc: t("landing.blockchainStat3Desc") },
              ].map((item) => (
                <div key={item.stat} className="rounded-lg border border-zinc-800 bg-zinc-800/50 p-4">
                  <p className="text-lg font-semibold text-white">{item.stat}</p>
                  <p className="mt-1 text-xs text-zinc-400">{item.desc}</p>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* ── Monetization Highlight ── */}
      <section className="mx-auto w-full max-w-6xl px-4 py-20 sm:px-6 sm:py-28">
        <div className="mb-14 text-center">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
            {t("landing.monetizationTitle")}
          </h2>
          <p className="mx-auto mt-4 max-w-2xl text-zinc-600 dark:text-zinc-400">
            {t("landing.monetizationSubtitle")}
          </p>
        </div>

        <div className="grid gap-6 sm:grid-cols-3">
          {[
            {
              icon: CheckCircle,
              title: t("landing.monetizationLicensingTitle"),
              desc: t("landing.monetizationLicensingDesc"),
            },
            {
              icon: CreditCard,
              title: t("landing.monetizationTokensTitle"),
              desc: t("landing.monetizationTokensDesc"),
            },
            {
              icon: Gift,
              title: t("landing.monetizationSubscriptionsTitle"),
              desc: t("landing.monetizationSubscriptionsDesc"),
            },
          ].map((item) => {
            const Icon = item.icon;
            return (
              <div
                key={item.title}
                className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950"
              >
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900">
                  <Icon className="h-5 w-5" />
                </div>
                <h3 className="mt-4 font-semibold">{item.title}</h3>
                <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                  {item.desc}
                </p>
              </div>
            );
          })}
        </div>
      </section>

      {/* ── Use Cases ── */}
      <section className="border-t border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/50">
        <div className="mx-auto max-w-6xl px-4 py-20 sm:px-6 sm:py-28">
          <div className="mb-14 text-center">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
              {t("landing.useCasesTitle")}
            </h2>
            <p className="mx-auto mt-4 max-w-2xl text-zinc-600 dark:text-zinc-400">
              {t("landing.useCasesSubtitle")}
            </p>
          </div>

          <div className="grid gap-6 sm:grid-cols-3">
            <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
              <h3 className="text-lg font-semibold">{t("landing.useCaseCreatorsTitle")}</h3>
              <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                {t("landing.useCaseCreatorsDesc")}
              </p>
            </div>
            <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
              <h3 className="text-lg font-semibold">{t("landing.useCaseBrandsTitle")}</h3>
              <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                {t("landing.useCaseBrandsDesc")}
              </p>
            </div>
            <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
              <h3 className="text-lg font-semibold">{t("landing.useCaseAgenciesTitle")}</h3>
              <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                {t("landing.useCaseAgenciesDesc")}
              </p>
            </div>
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
            { abbr: "Be", name: "Behance" },
            { abbr: "Dr", name: "Dribbble" },
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

          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
            {[1, 2, 3, 4].map((i) => (
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
            {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map((i) => (
              <div
                key={i}
                className="rounded-xl border border-zinc-200 bg-white dark:border-zinc-800 dark:bg-zinc-950"
              >
                <button
                  onClick={() => toggleFaq(i - 1)}
                  className="flex w-full items-center justify-between px-6 py-4 text-left"
                >
                  <span className="font-medium">{t(`landing.faq${i}Q`)}</span>
                  <ChevronDown
                    className={`h-5 w-5 shrink-0 text-zinc-400 transition-transform duration-200 ${
                      openFaq === i - 1 ? "rotate-180" : ""
                    }`}
                  />
                </button>
                {openFaq === i - 1 && (
                  <div className="border-t border-zinc-100 px-6 py-4 text-sm text-zinc-600 dark:border-zinc-800 dark:text-zinc-400">
                    {t(`landing.faq${i}A`)}
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Blog Highlights ── */}
      <section className="border-t border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/50">
        <div className="mx-auto max-w-6xl px-4 py-20 sm:px-6 sm:py-28">
          <div className="mb-14 text-center">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
              {t("landing.blogHighlightsTitle")}
            </h2>
            <p className="mx-auto mt-4 max-w-2xl text-zinc-600 dark:text-zinc-400">
              {t("landing.blogHighlightsSubtitle")}
            </p>
          </div>

          <div className="grid gap-6 sm:grid-cols-3">
            {blogPosts.slice(0, 3).map((post) => (
              <Link
                key={post.slug}
                href={`/blog/${post.slug}`}
                className="group rounded-xl border border-zinc-200 bg-white p-6 transition-colors hover:border-zinc-300 dark:border-zinc-800 dark:bg-zinc-950 dark:hover:border-zinc-700"
              >
                <p className="text-xs text-zinc-500 dark:text-zinc-400">
                  {new Date(post.publishedAt).toLocaleDateString("en-US", {
                    year: "numeric",
                    month: "long",
                    day: "numeric",
                  })}{" "}
                  &middot; {post.readTime} min read
                </p>
                <h3 className="mt-2 font-semibold group-hover:underline">
                  {post.title}
                </h3>
                <p className="mt-2 text-sm text-zinc-600 line-clamp-2 dark:text-zinc-400">
                  {post.excerpt}
                </p>
              </Link>
            ))}
          </div>

          <div className="mt-8 text-center">
            <Link
              href="/blog"
              className="inline-flex items-center gap-2 text-sm font-medium text-zinc-600 transition-colors hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100"
            >
              {t("landing.blogViewAll")}
              <ArrowRight className="h-4 w-4" />
            </Link>
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
