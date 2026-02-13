"use client";

import { useSearchParams } from "next/navigation";
import { useEffect, useState, useCallback, Suspense } from "react";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth-context";
import Link from "next/link";
import {
  ArrowLeft,
  Image,
  Video,
  Music,
  FileText,
  File,
  CheckCircle,
  XCircle,
  ShoppingCart,
  User,
} from "@/components/icons";
import { useTranslation } from "react-i18next";

type MarketplaceDetail = {
  content: {
    id: string;
    title: string;
    description: string;
    contentType: string;
    tags: string[];
    thumbnailUrl?: string;
    hashSha256: string;
    createdAt: string;
  };
  offerings: {
    id: string;
    licenseType: string;
    priceCents: number;
    currency: string;
    termsText: string;
    isActive: boolean;
  }[];
  creatorName: string | null;
  creatorUsername: string | null;
  creatorImage: string | null;
};

function ContentTypeIcon({ contentType, className }: { contentType: string; className?: string }) {
  switch (contentType) {
    case "image":
      return <Image className={className} />;
    case "video":
      return <Video className={className} />;
    case "audio":
      return <Music className={className} />;
    case "text":
      return <FileText className={className} />;
    default:
      return <File className={className} />;
  }
}

function MarketplaceDetailContent() {
  const searchParams = useSearchParams();
  const { user } = useAuth();
  const { t } = useTranslation();
  const id = searchParams.get("id");

  const [item, setItem] = useState<MarketplaceDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [buyingId, setBuyingId] = useState<string | null>(null);

  const success = searchParams.get("success") === "true";
  const canceled = searchParams.get("canceled") === "true";

  const fetchItem = useCallback(async () => {
    if (!id) return;
    setLoading(true);
    const result = await api.marketplace.detail(id);
    if (result.data) {
      setItem(result.data);
    } else {
      setError(result.error || "Not found");
    }
    setLoading(false);
  }, [id]);

  useEffect(() => {
    if (id) fetchItem();
  }, [id, fetchItem]);

  const handleBuy = async (offeringId: string) => {
    setBuyingId(offeringId);
    const result = await api.licenses.checkout(offeringId);
    if (result.data?.url) {
      window.location.href = result.data.url;
    } else {
      setBuyingId(null);
    }
  };

  if (!id) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
        <p className="text-center text-red-500">No content ID specified</p>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
        <p className="text-center text-zinc-400">{t("marketplace.detail.loading")}</p>
      </div>
    );
  }

  if (error || !item) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
        <p className="text-center text-red-500">{error || t("marketplace.detail.notFound")}</p>
        <div className="mt-4 text-center">
          <Link
            href="/marketplace"
            className="text-sm text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300"
          >
            {t("marketplace.detail.backToMarketplace")}
          </Link>
        </div>
      </div>
    );
  }

  const content = item.content;
  const activeOfferings = item.offerings.filter((o) => o.isActive);

  return (
    <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
      {/* Back link */}
      <Link
        href="/marketplace"
        className="mb-6 inline-flex items-center gap-1 text-sm text-zinc-500 transition-colors hover:text-zinc-700 dark:hover:text-zinc-300"
      >
        <ArrowLeft className="h-4 w-4" />
        {t("marketplace.detail.backToMarketplace")}
      </Link>

      {/* Success / Canceled banners */}
      {success && (
        <div className="mb-6 flex items-center gap-3 rounded-xl border border-emerald-200 bg-emerald-50 p-4 dark:border-emerald-800 dark:bg-emerald-950">
          <CheckCircle className="h-5 w-5 text-emerald-600" />
          <p className="text-sm font-medium text-emerald-800 dark:text-emerald-200">
            {t("marketplace.detail.purchaseSuccess")}
          </p>
        </div>
      )}
      {canceled && (
        <div className="mb-6 flex items-center gap-3 rounded-xl border border-amber-200 bg-amber-50 p-4 dark:border-amber-800 dark:bg-amber-950">
          <XCircle className="h-5 w-5 text-amber-600" />
          <p className="text-sm font-medium text-amber-800 dark:text-amber-200">
            {t("marketplace.detail.purchaseCanceled")}
          </p>
        </div>
      )}

      <div className="grid gap-8 lg:grid-cols-3">
        {/* Main content */}
        <div className="lg:col-span-2">
          {/* Content Preview */}
          <div className="mb-6 overflow-hidden rounded-xl border border-zinc-200 bg-white dark:border-zinc-800 dark:bg-zinc-900">
            <div className="relative flex items-center justify-center bg-zinc-50 p-8 dark:bg-zinc-800">
              {content.thumbnailUrl && content.contentType === "image" ? (
                <img
                  src={content.thumbnailUrl}
                  alt={content.title}
                  className="max-h-96 rounded-lg object-contain"
                />
              ) : (
                <div className="text-center">
                  <ContentTypeIcon
                    contentType={content.contentType}
                    className="mx-auto h-16 w-16 text-zinc-400"
                  />
                  <p className="mt-2 text-sm capitalize text-zinc-500">{content.contentType}</p>
                </div>
              )}
              <div className="absolute bottom-3 right-3 rounded-lg bg-zinc-900/70 px-2.5 py-1 text-xs text-white">
                {t("marketplace.detail.previewOnly")}
              </div>
            </div>
          </div>

          {/* Title & Description */}
          <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
            <h1 className="text-xl font-bold text-zinc-900 dark:text-zinc-100">{content.title}</h1>
            {content.description && (
              <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">{content.description}</p>
            )}
            {content.tags && content.tags.length > 0 && (
              <div className="mt-3 flex flex-wrap gap-1.5">
                {content.tags.map((tag) => (
                  <span
                    key={tag}
                    className="rounded-md bg-zinc-100 px-2 py-0.5 text-xs text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-4">
          {/* Creator Info */}
          <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900">
            <h3 className="text-sm font-semibold text-zinc-900 dark:text-zinc-100">
              {t("marketplace.detail.creator")}
            </h3>
            <div className="mt-3 flex items-center gap-3">
              {item.creatorImage ? (
                <img
                  src={item.creatorImage}
                  alt=""
                  className="h-10 w-10 rounded-full object-cover"
                />
              ) : (
                <div className="flex h-10 w-10 items-center justify-center rounded-full bg-zinc-100 dark:bg-zinc-800">
                  <User className="h-5 w-5 text-zinc-400" />
                </div>
              )}
              <div>
                <p className="text-sm font-medium text-zinc-900 dark:text-zinc-100">
                  {item.creatorName || t("common.creator")}
                </p>
                {item.creatorUsername && (
                  <Link
                    href={`/profile?u=${item.creatorUsername}`}
                    className="text-xs text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300"
                  >
                    @{item.creatorUsername}
                  </Link>
                )}
              </div>
            </div>
          </div>

          {/* License Options */}
          <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900">
            <h3 className="text-sm font-semibold text-zinc-900 dark:text-zinc-100">
              {t("marketplace.detail.licenseOptions")}
            </h3>

            {activeOfferings.length === 0 ? (
              <p className="mt-3 text-sm text-zinc-500">
                {t("marketplace.detail.noLicensesAvailable")}
              </p>
            ) : (
              <div className="mt-3 space-y-3">
                {activeOfferings.map((offering) => (
                  <div
                    key={offering.id}
                    className="rounded-lg border border-zinc-200 p-3 dark:border-zinc-700"
                  >
                    <div className="flex items-center justify-between">
                      <span className="text-sm font-medium capitalize text-zinc-900 dark:text-zinc-100">
                        {offering.licenseType.replace("_", " ")}
                      </span>
                      <span className="text-sm font-bold text-zinc-900 dark:text-zinc-100">
                        ${(offering.priceCents / 100).toFixed(2)}
                      </span>
                    </div>
                    {offering.termsText && (
                      <p className="mt-1 text-xs text-zinc-500">{offering.termsText}</p>
                    )}
                    <button
                      onClick={() => handleBuy(offering.id)}
                      disabled={buyingId === offering.id || !user}
                      className="mt-2 flex w-full items-center justify-center gap-1.5 rounded-lg bg-zinc-900 px-3 py-1.5 text-xs font-medium text-white transition-colors hover:bg-zinc-700 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-300"
                    >
                      <ShoppingCart className="h-3.5 w-3.5" />
                      {buyingId === offering.id
                        ? t("marketplace.detail.processing")
                        : t("marketplace.detail.buy")}
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Content Type */}
          <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900">
            <h3 className="text-sm font-semibold text-zinc-900 dark:text-zinc-100">
              {t("marketplace.detail.contentType")}
            </h3>
            <div className="mt-3 flex items-center gap-2">
              <ContentTypeIcon
                contentType={content.contentType}
                className="h-5 w-5 text-zinc-500"
              />
              <span className="text-sm capitalize text-zinc-600 dark:text-zinc-400">
                {content.contentType}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function MarketplaceDetailPage() {
  return (
    <Suspense fallback={null}>
      <MarketplaceDetailContent />
    </Suspense>
  );
}
