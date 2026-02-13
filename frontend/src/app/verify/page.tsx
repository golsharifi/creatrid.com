"use client";

import { useSearchParams } from "next/navigation";
import { useEffect, useState, Suspense, useCallback } from "react";
import { api } from "@/lib/api";
import { Shield, Search, ExternalLink, CheckCircle, Clock, AlertTriangle } from "@/components/icons";
import { QRCodeSVG } from "qrcode.react";
import { useTranslation } from "react-i18next";

type AnchorData = {
  id: string;
  contentId: string;
  userId: string;
  contentHash: string;
  txHash: string | null;
  chain: string;
  blockNumber: number | null;
  contractAddress: string | null;
  anchorStatus: string;
  errorMessage: string | null;
  createdAt: string;
  confirmedAt: string | null;
};

type ContentData = {
  id: string;
  title: string;
  contentType: string;
  createdAt: string;
};

function StatusBadge({ status }: { status: string }) {
  const { t } = useTranslation();

  if (status === "confirmed") {
    return (
      <span className="inline-flex items-center gap-1.5 rounded-full bg-emerald-50 px-3 py-1 text-sm font-medium text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300">
        <CheckCircle className="h-4 w-4" />
        {t("blockchain.verified")}
      </span>
    );
  }
  if (status === "pending") {
    return (
      <span className="inline-flex items-center gap-1.5 rounded-full bg-amber-50 px-3 py-1 text-sm font-medium text-amber-700 dark:bg-amber-950 dark:text-amber-300">
        <Clock className="h-4 w-4" />
        {t("blockchain.pending")}
      </span>
    );
  }
  return (
    <span className="inline-flex items-center gap-1.5 rounded-full bg-red-50 px-3 py-1 text-sm font-medium text-red-700 dark:bg-red-950 dark:text-red-300">
      <AlertTriangle className="h-4 w-4" />
      {t("blockchain.failed")}
    </span>
  );
}

function VerifyContent() {
  const searchParams = useSearchParams();
  const { t } = useTranslation();
  const hashParam = searchParams.get("hash");
  const txParam = searchParams.get("tx");

  const [searchHash, setSearchHash] = useState(hashParam || "");
  const [anchor, setAnchor] = useState<AnchorData | null>(null);
  const [content, setContent] = useState<ContentData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [searched, setSearched] = useState(false);

  const doVerify = useCallback(async (hash: string) => {
    if (!hash.trim()) return;
    setLoading(true);
    setError("");
    setAnchor(null);
    setContent(null);
    setSearched(true);

    const result = await api.blockchain.verify(hash.trim());
    if (result.data) {
      setAnchor(result.data.anchor);
      setContent(result.data.content);
    } else {
      setError(result.error || t("blockchain.notFound"));
    }
    setLoading(false);
  }, [t]);

  // Auto-fetch if hash or tx is provided
  useEffect(() => {
    if (hashParam) {
      setSearchHash(hashParam);
      doVerify(hashParam);
    }
    // tx param: we don't have a client-side API for tx lookup, so just show a message
    // The hash param is the primary way to verify
  }, [hashParam, txParam, doVerify]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    doVerify(searchHash);
  };

  const verifyUrl =
    typeof window !== "undefined" && anchor
      ? `${window.location.origin}/verify?hash=${anchor.contentHash}`
      : "";

  return (
    <div className="mx-auto max-w-3xl px-4 py-12 sm:px-6">
      {/* Header */}
      <div className="mb-10 text-center">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-emerald-500 to-teal-600 shadow-lg">
          <Shield className="h-8 w-8 text-white" />
        </div>
        <h1 className="text-3xl font-bold tracking-tight text-zinc-900 dark:text-zinc-100">
          {t("blockchain.verifyTitle")}
        </h1>
        <p className="mt-2 text-zinc-500 dark:text-zinc-400">
          {t("blockchain.verifySubtitle")}
        </p>
      </div>

      {/* Search Form */}
      <form onSubmit={handleSubmit} className="mb-10">
        <div className="flex gap-3">
          <div className="relative flex-1">
            <Search className="pointer-events-none absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-zinc-400" />
            <input
              type="text"
              value={searchHash}
              onChange={(e) => setSearchHash(e.target.value)}
              placeholder={t("blockchain.searchPlaceholder")}
              className="w-full rounded-xl border border-zinc-200 bg-white py-3 pl-10 pr-4 font-mono text-sm text-zinc-900 transition-colors placeholder:text-zinc-400 focus:border-emerald-400 focus:outline-none focus:ring-2 focus:ring-emerald-100 dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-100 dark:focus:border-emerald-500 dark:focus:ring-emerald-900"
            />
          </div>
          <button
            type="submit"
            disabled={loading || !searchHash.trim()}
            className="rounded-xl bg-emerald-600 px-6 py-3 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-emerald-700 disabled:opacity-50 dark:bg-emerald-500 dark:hover:bg-emerald-600"
          >
            {loading ? t("common.loading") : t("blockchain.verifyButton")}
          </button>
        </div>
      </form>

      {/* Error State */}
      {error && searched && (
        <div className="rounded-xl border border-red-200 bg-red-50 p-6 text-center dark:border-red-800 dark:bg-red-950">
          <AlertTriangle className="mx-auto mb-3 h-8 w-8 text-red-400" />
          <p className="text-sm font-medium text-red-700 dark:text-red-300">{error}</p>
        </div>
      )}

      {/* Result */}
      {anchor && (
        <div className="space-y-6">
          {/* Status Banner */}
          <div
            className={`rounded-2xl border p-6 text-center ${
              anchor.anchorStatus === "confirmed"
                ? "border-emerald-200 bg-gradient-to-br from-emerald-50 to-teal-50 dark:border-emerald-800 dark:from-emerald-950 dark:to-teal-950"
                : anchor.anchorStatus === "pending"
                  ? "border-amber-200 bg-gradient-to-br from-amber-50 to-yellow-50 dark:border-amber-800 dark:from-amber-950 dark:to-yellow-950"
                  : "border-red-200 bg-gradient-to-br from-red-50 to-pink-50 dark:border-red-800 dark:from-red-950 dark:to-pink-950"
            }`}
          >
            {anchor.anchorStatus === "confirmed" ? (
              <CheckCircle className="mx-auto mb-3 h-12 w-12 text-emerald-500" />
            ) : anchor.anchorStatus === "pending" ? (
              <Clock className="mx-auto mb-3 h-12 w-12 text-amber-500" />
            ) : (
              <AlertTriangle className="mx-auto mb-3 h-12 w-12 text-red-500" />
            )}
            <StatusBadge status={anchor.anchorStatus} />
            {content && (
              <p className="mt-3 text-lg font-semibold text-zinc-900 dark:text-zinc-100">
                {content.title}
              </p>
            )}
          </div>

          {/* Proof Certificate */}
          <div className="rounded-2xl border border-zinc-200 bg-white p-6 shadow-sm dark:border-zinc-800 dark:bg-zinc-900">
            <h2 className="mb-4 text-lg font-semibold text-zinc-900 dark:text-zinc-100">
              {t("blockchain.proofCertificate")}
            </h2>

            <dl className="space-y-4">
              {/* Content Hash */}
              <div>
                <dt className="text-xs font-medium uppercase tracking-wider text-zinc-500">
                  {t("blockchain.contentHash")}
                </dt>
                <dd className="mt-1 break-all rounded-lg bg-zinc-50 px-3 py-2 font-mono text-xs text-zinc-700 dark:bg-zinc-800 dark:text-zinc-300">
                  {anchor.contentHash}
                </dd>
              </div>

              {/* Blockchain */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <dt className="text-xs font-medium uppercase tracking-wider text-zinc-500">
                    {t("blockchain.chain")}
                  </dt>
                  <dd className="mt-1 text-sm font-medium capitalize text-zinc-900 dark:text-zinc-100">
                    {anchor.chain}
                  </dd>
                </div>
                <div>
                  <dt className="text-xs font-medium uppercase tracking-wider text-zinc-500">
                    {t("blockchain.blockNumber")}
                  </dt>
                  <dd className="mt-1 text-sm font-medium text-zinc-900 dark:text-zinc-100">
                    {anchor.blockNumber ? anchor.blockNumber.toLocaleString() : "-"}
                  </dd>
                </div>
              </div>

              {/* Transaction Hash */}
              {anchor.txHash && (
                <div>
                  <dt className="text-xs font-medium uppercase tracking-wider text-zinc-500">
                    {t("blockchain.txHash")}
                  </dt>
                  <dd className="mt-1 flex items-center gap-2">
                    <span className="break-all rounded-lg bg-zinc-50 px-3 py-2 font-mono text-xs text-zinc-700 dark:bg-zinc-800 dark:text-zinc-300">
                      {anchor.txHash}
                    </span>
                    {anchor.chain === "polygon" && !anchor.txHash.startsWith("0x000000") && (
                      <a
                        href={`https://polygonscan.com/tx/${anchor.txHash}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex-shrink-0 text-emerald-600 hover:text-emerald-700"
                      >
                        <ExternalLink className="h-4 w-4" />
                      </a>
                    )}
                  </dd>
                </div>
              )}

              {/* Anchored At */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <dt className="text-xs font-medium uppercase tracking-wider text-zinc-500">
                    {t("blockchain.anchoredAt")}
                  </dt>
                  <dd className="mt-1 text-sm text-zinc-900 dark:text-zinc-100">
                    {anchor.confirmedAt
                      ? new Date(anchor.confirmedAt).toLocaleString()
                      : new Date(anchor.createdAt).toLocaleString()}
                  </dd>
                </div>
                {content && (
                  <div>
                    <dt className="text-xs font-medium uppercase tracking-wider text-zinc-500">
                      Content Type
                    </dt>
                    <dd className="mt-1 text-sm capitalize text-zinc-900 dark:text-zinc-100">
                      {content.contentType}
                    </dd>
                  </div>
                )}
              </div>
            </dl>

            {/* Simulated notice */}
            {anchor.contractAddress === "0x0000000000000000000000000000000000000000" && (
              <div className="mt-4 rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-700 dark:border-amber-800 dark:bg-amber-950 dark:text-amber-300">
                {t("blockchain.simulated")}
              </div>
            )}
          </div>

          {/* QR Code for sharing */}
          {verifyUrl && (
            <div className="rounded-2xl border border-zinc-200 bg-white p-6 text-center shadow-sm dark:border-zinc-800 dark:bg-zinc-900">
              <p className="mb-4 text-sm font-medium text-zinc-700 dark:text-zinc-300">
                Scan to verify this proof
              </p>
              <div className="mx-auto inline-block rounded-xl border border-zinc-200 bg-white p-3 dark:border-zinc-700">
                <QRCodeSVG value={verifyUrl} size={160} level="M" marginSize={2} />
              </div>
              <p className="mt-3 break-all text-xs text-zinc-400">{verifyUrl}</p>
            </div>
          )}
        </div>
      )}

      {/* Empty state when no search */}
      {!anchor && !error && !searched && (
        <div className="mt-8 text-center">
          <div className="mx-auto mb-4 flex h-20 w-20 items-center justify-center rounded-full bg-zinc-100 dark:bg-zinc-800">
            <Shield className="h-10 w-10 text-zinc-400" />
          </div>
          <p className="text-sm text-zinc-500">
            {t("blockchain.verifySubtitle")}
          </p>
        </div>
      )}
    </div>
  );
}

export default function VerifyPage() {
  return (
    <Suspense fallback={null}>
      <VerifyContent />
    </Suspense>
  );
}
