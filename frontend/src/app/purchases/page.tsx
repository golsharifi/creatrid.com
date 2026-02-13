"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState, useCallback } from "react";
import { api } from "@/lib/api";
import Link from "next/link";
import {
  Download,
  ShoppingBag,
  Image,
  Video,
  Music,
  FileText,
  File,
} from "lucide-react";
import { useTranslation } from "react-i18next";

type Purchase = {
  id: string;
  contentId: string;
  contentTitle: string;
  contentType: string;
  licenseType: string;
  priceCents: number;
  purchasedAt: string;
  creatorName: string | null;
  creatorUsername: string | null;
};

function getContentCategory(contentType: string): string {
  if (contentType.startsWith("image/")) return "image";
  if (contentType.startsWith("video/")) return "video";
  if (contentType.startsWith("audio/")) return "audio";
  if (
    contentType.startsWith("text/") ||
    contentType === "application/pdf" ||
    contentType === "application/msword" ||
    contentType.includes("document")
  )
    return "text";
  return "other";
}

function ContentTypeIcon({ contentType, className }: { contentType: string; className?: string }) {
  const category = getContentCategory(contentType);
  switch (category) {
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

export default function PurchasesPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();
  const [purchases, setPurchases] = useState<Purchase[]>([]);
  const [loadingPurchases, setLoadingPurchases] = useState(true);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
  }, [user, loading, router]);

  const fetchPurchases = useCallback(async () => {
    setLoadingPurchases(true);
    const result = await api.licenses.purchases();
    if (result.data) {
      setPurchases(result.data.purchases || []);
    }
    setLoadingPurchases(false);
  }, []);

  useEffect(() => {
    if (user) fetchPurchases();
  }, [user, fetchPurchases]);

  if (loading || !user) return null;

  return (
    <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">
          {t("purchases.title")}
        </h1>
        <p className="mt-1 text-zinc-500">
          {t("purchases.subtitle")}
        </p>
      </div>

      {loadingPurchases ? (
        <p className="py-12 text-center text-zinc-400">{t("purchases.loading")}</p>
      ) : purchases.length === 0 ? (
        <div className="py-16 text-center">
          <ShoppingBag className="mx-auto h-12 w-12 text-zinc-300 dark:text-zinc-600" />
          <p className="mt-4 text-zinc-500">{t("purchases.emptyState")}</p>
          <Link
            href="/marketplace"
            className="mt-4 inline-flex items-center gap-2 rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-700 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-300"
          >
            {t("purchases.browseMarketplace")}
          </Link>
        </div>
      ) : (
        <div className="space-y-3">
          {purchases.map((purchase) => (
            <div
              key={purchase.id}
              className="flex items-center justify-between rounded-xl border border-zinc-200 bg-white p-4 dark:border-zinc-800 dark:bg-zinc-900"
            >
              <div className="flex items-center gap-4">
                <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-800">
                  <ContentTypeIcon
                    contentType={purchase.contentType}
                    className="h-6 w-6 text-zinc-500"
                  />
                </div>
                <div>
                  <h3 className="text-sm font-semibold text-zinc-900 dark:text-zinc-100">
                    {purchase.contentTitle}
                  </h3>
                  <div className="mt-0.5 flex items-center gap-3 text-xs text-zinc-500">
                    <span className="capitalize">
                      {purchase.licenseType.replace("_", " ")} {t("purchases.license")}
                    </span>
                    <span>${(purchase.priceCents / 100).toFixed(2)}</span>
                    <span>{new Date(purchase.purchasedAt).toLocaleDateString()}</span>
                  </div>
                  {(purchase.creatorName || purchase.creatorUsername) && (
                    <p className="mt-0.5 text-xs text-zinc-400">
                      {t("purchases.by")}{" "}
                      {purchase.creatorUsername ? (
                        <Link
                          href={`/profile?u=${purchase.creatorUsername}`}
                          className="hover:text-zinc-600 dark:hover:text-zinc-300"
                        >
                          {purchase.creatorName || `@${purchase.creatorUsername}`}
                        </Link>
                      ) : (
                        purchase.creatorName
                      )}
                    </p>
                  )}
                </div>
              </div>
              <a
                href={api.content.download(purchase.contentId)}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-1.5 rounded-lg border border-zinc-300 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
              >
                <Download className="h-3.5 w-3.5" />
                {t("purchases.download")}
              </a>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
