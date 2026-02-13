"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState, useCallback, Suspense } from "react";
import { api } from "@/lib/api";
import Link from "next/link";
import {
  ArrowLeft,
  Download,
  Trash2,
  Image,
  Video,
  Music,
  FileText,
  File,
  Lock,
  Globe,
  Shield,
  Plus,
  X,
  Edit3,
  Check,
  CheckCircle,
  Clock,
  AlertTriangle,
  ExternalLink,
} from "@/components/icons";
import { useTranslation } from "react-i18next";

type ContentDetail = {
  id: string;
  title: string;
  description: string;
  contentType: string;
  mimeType: string;
  fileSize: number;
  isPublic: boolean;
  tags: string[];
  fileUrl?: string;
  thumbnailUrl?: string;
  hashSha256: string;
  createdAt: string;
};

type LicenseOffering = {
  id: string;
  licenseType: string;
  priceCents: number;
  termsText: string;
  isActive: boolean;
  createdAt: string;
};

const LICENSE_TYPES = [
  { value: "personal", label: "Personal" },
  { value: "commercial", label: "Commercial" },
  { value: "editorial", label: "Editorial" },
  { value: "ai_training", label: "AI Training" },
];

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

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`;
}

function VaultDetailContent() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();
  const { t } = useTranslation();
  const id = searchParams.get("id");

  const [item, setItem] = useState<ContentDetail | null>(null);
  const [loadingItem, setLoadingItem] = useState(true);
  const [error, setError] = useState("");

  // Edit state
  const [editing, setEditing] = useState(false);
  const [editTitle, setEditTitle] = useState("");
  const [editDescription, setEditDescription] = useState("");
  const [editTags, setEditTags] = useState("");
  const [saving, setSaving] = useState(false);

  // Delete state
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [deleting, setDeleting] = useState(false);

  // License state
  const [offerings, setOfferings] = useState<LicenseOffering[]>([]);
  const [showAddLicense, setShowAddLicense] = useState(false);
  const [newLicenseType, setNewLicenseType] = useState("personal");
  const [newLicensePrice, setNewLicensePrice] = useState("");
  const [newLicenseTerms, setNewLicenseTerms] = useState("");
  const [addingLicense, setAddingLicense] = useState(false);

  // Proof state
  const [proof, setProof] = useState<{
    id: string;
    hashSha256: string;
    createdAt: string;
    title: string;
  } | null>(null);

  // Blockchain anchor state
  const [anchor, setAnchor] = useState<{
    id: string;
    contentId: string;
    contentHash: string;
    txHash: string | null;
    chain: string;
    blockNumber: number | null;
    contractAddress: string | null;
    anchorStatus: string;
    createdAt: string;
    confirmedAt: string | null;
  } | null>(null);
  const [anchoring, setAnchoring] = useState(false);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
  }, [user, loading, router]);

  const fetchItem = useCallback(async () => {
    if (!id) return;
    setLoadingItem(true);
    const result = await api.content.get(id);
    if (result.data) {
      setItem(result.data);
      setEditTitle(result.data.title);
      setEditDescription(result.data.description || "");
      setEditTags((result.data.tags || []).join(", "));
    } else {
      setError(result.error || "Not found");
    }
    setLoadingItem(false);
  }, [id]);

  const fetchOfferings = useCallback(async () => {
    if (!id) return;
    const result = await api.licenses.list(id);
    if (result.data) {
      setOfferings(result.data.offerings || []);
    }
  }, [id]);

  const fetchProof = useCallback(async () => {
    if (!id) return;
    const result = await api.content.proof(id);
    if (result.data) {
      setProof(result.data);
    }
  }, [id]);

  const fetchAnchor = useCallback(async () => {
    if (!id) return;
    const result = await api.blockchain.getAnchor(id);
    if (result.data) {
      setAnchor(result.data.anchor);
    }
  }, [id]);

  useEffect(() => {
    if (user && id) {
      fetchItem();
      fetchOfferings();
      fetchProof();
      fetchAnchor();
    }
  }, [user, id, fetchItem, fetchOfferings, fetchProof, fetchAnchor]);

  const handleSave = async () => {
    if (!item) return;
    setSaving(true);
    const tags = editTags
      .split(",")
      .map((tag) => tag.trim())
      .filter(Boolean);
    const result = await api.content.update(item.id, {
      title: editTitle.trim(),
      description: editDescription.trim(),
      tags,
    });
    if (result.data?.success) {
      setItem({ ...item, title: editTitle.trim(), description: editDescription.trim(), tags });
      setEditing(false);
    }
    setSaving(false);
  };

  const handleTogglePublic = async () => {
    if (!item) return;
    const result = await api.content.update(item.id, { isPublic: !item.isPublic });
    if (result.data?.success) {
      setItem({ ...item, isPublic: !item.isPublic });
    }
  };

  const handleDelete = async () => {
    if (!item) return;
    setDeleting(true);
    const result = await api.content.delete(item.id);
    if (result.data) {
      router.push("/vault");
    } else {
      setDeleting(false);
      setShowDeleteConfirm(false);
    }
  };

  const handleAddLicense = async () => {
    if (!id || !newLicensePrice) return;
    setAddingLicense(true);
    const priceCents = Math.round(parseFloat(newLicensePrice) * 100);
    const result = await api.licenses.create(
      id,
      newLicenseType,
      priceCents,
      newLicenseTerms.trim() || undefined
    );
    if (result.data) {
      await fetchOfferings();
      setShowAddLicense(false);
      setNewLicenseType("personal");
      setNewLicensePrice("");
      setNewLicenseTerms("");
    }
    setAddingLicense(false);
  };

  const handleDeleteLicense = async (offeringId: string) => {
    await api.licenses.delete(offeringId);
    setOfferings((prev) => prev.filter((o) => o.id !== offeringId));
  };

  const handleAnchor = async () => {
    if (!id) return;
    setAnchoring(true);
    const result = await api.blockchain.anchor(id);
    if (result.data) {
      setAnchor(result.data.anchor);
    }
    setAnchoring(false);
  };

  if (loading || !user) return null;

  if (!id) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
        <p className="text-center text-red-500">No content ID specified</p>
      </div>
    );
  }

  if (loadingItem) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
        <p className="text-center text-zinc-400">{t("vault.detail.loading")}</p>
      </div>
    );
  }

  if (error || !item) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
        <p className="text-center text-red-500">{error || t("vault.detail.notFound")}</p>
        <div className="mt-4 text-center">
          <Link href="/vault" className="text-sm text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300">
            {t("vault.detail.backToVault")}
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
      {/* Back link */}
      <Link
        href="/vault"
        className="mb-6 inline-flex items-center gap-1 text-sm text-zinc-500 transition-colors hover:text-zinc-700 dark:hover:text-zinc-300"
      >
        <ArrowLeft className="h-4 w-4" />
        {t("vault.detail.backToVault")}
      </Link>

      <div className="grid gap-8 lg:grid-cols-3">
        {/* Main content - left column */}
        <div className="lg:col-span-2">
          {/* Content Preview */}
          <div className="mb-6 overflow-hidden rounded-xl border border-zinc-200 bg-white dark:border-zinc-800 dark:bg-zinc-900">
            <div className="flex items-center justify-center bg-zinc-50 p-8 dark:bg-zinc-800">
              {item.contentType === "image" && item.thumbnailUrl ? (
                <img
                  src={item.thumbnailUrl}
                  alt={item.title}
                  className="max-h-96 rounded-lg object-contain"
                />
              ) : item.contentType === "audio" ? (
                <div className="w-full px-4">
                  <Music className="mx-auto mb-4 h-16 w-16 text-zinc-400" />
                </div>
              ) : (
                <div className="text-center">
                  <ContentTypeIcon
                    contentType={item.contentType}
                    className="mx-auto h-16 w-16 text-zinc-400"
                  />
                  <p className="mt-2 text-sm capitalize text-zinc-500">{item.contentType}</p>
                </div>
              )}
            </div>
          </div>

          {/* Title and Description */}
          <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
            {editing ? (
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
                    {t("vault.detail.titleLabel")}
                  </label>
                  <input
                    type="text"
                    value={editTitle}
                    onChange={(e) => setEditTitle(e.target.value)}
                    className="mt-1 w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm text-zinc-900 focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-100"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
                    {t("vault.detail.descriptionLabel")}
                  </label>
                  <textarea
                    value={editDescription}
                    onChange={(e) => setEditDescription(e.target.value)}
                    rows={3}
                    className="mt-1 w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm text-zinc-900 focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-100"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
                    {t("vault.detail.tagsLabel")}
                  </label>
                  <input
                    type="text"
                    value={editTags}
                    onChange={(e) => setEditTags(e.target.value)}
                    placeholder={t("vault.upload.tagsPlaceholder")}
                    className="mt-1 w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm text-zinc-900 focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-100"
                  />
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={handleSave}
                    disabled={saving}
                    className="flex items-center gap-1.5 rounded-lg bg-zinc-900 px-3 py-1.5 text-xs font-medium text-white transition-colors hover:bg-zinc-700 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-300"
                  >
                    <Check className="h-3.5 w-3.5" />
                    {saving ? t("vault.detail.saving") : t("vault.detail.save")}
                  </button>
                  <button
                    onClick={() => {
                      setEditing(false);
                      setEditTitle(item.title);
                      setEditDescription(item.description || "");
                      setEditTags((item.tags || []).join(", "));
                    }}
                    className="rounded-lg border border-zinc-300 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
                  >
                    {t("vault.detail.cancel")}
                  </button>
                </div>
              </div>
            ) : (
              <>
                <div className="flex items-start justify-between">
                  <h1 className="text-xl font-bold text-zinc-900 dark:text-zinc-100">
                    {item.title}
                  </h1>
                  <button
                    onClick={() => setEditing(true)}
                    className="rounded-lg p-1.5 text-zinc-400 transition-colors hover:bg-zinc-100 hover:text-zinc-600 dark:hover:bg-zinc-800"
                  >
                    <Edit3 className="h-4 w-4" />
                  </button>
                </div>
                {item.description && (
                  <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                    {item.description}
                  </p>
                )}
                {item.tags && item.tags.length > 0 && (
                  <div className="mt-3 flex flex-wrap gap-1.5">
                    {item.tags.map((tag) => (
                      <span
                        key={tag}
                        className="rounded-md bg-zinc-100 px-2 py-0.5 text-xs text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                )}
              </>
            )}
          </div>

          {/* License Management */}
          <div className="mt-6 rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
                {t("vault.detail.licenseManagement")}
              </h2>
              <button
                onClick={() => setShowAddLicense(!showAddLicense)}
                className="flex items-center gap-1.5 rounded-lg border border-zinc-300 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
              >
                <Plus className="h-3.5 w-3.5" />
                {t("vault.detail.addLicense")}
              </button>
            </div>

            {/* Add License Form */}
            {showAddLicense && (
              <div className="mt-4 space-y-3 rounded-lg border border-zinc-200 bg-zinc-50 p-4 dark:border-zinc-700 dark:bg-zinc-800">
                <div>
                  <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
                    {t("vault.detail.licenseType")}
                  </label>
                  <select
                    value={newLicenseType}
                    onChange={(e) => setNewLicenseType(e.target.value)}
                    className="mt-1 w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-100"
                  >
                    {LICENSE_TYPES.map((lt) => (
                      <option key={lt.value} value={lt.value}>
                        {lt.label}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
                    {t("vault.detail.price")}
                  </label>
                  <div className="relative mt-1">
                    <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-zinc-400">
                      $
                    </span>
                    <input
                      type="number"
                      step="0.01"
                      min="0"
                      value={newLicensePrice}
                      onChange={(e) => setNewLicensePrice(e.target.value)}
                      placeholder="0.00"
                      className="w-full rounded-lg border border-zinc-200 bg-white py-2 pl-7 pr-3 text-sm focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-100"
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
                    {t("vault.detail.terms")}
                  </label>
                  <textarea
                    value={newLicenseTerms}
                    onChange={(e) => setNewLicenseTerms(e.target.value)}
                    rows={2}
                    placeholder={t("vault.detail.termsPlaceholder")}
                    className="mt-1 w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm focus:border-zinc-400 focus:outline-none dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-100"
                  />
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={handleAddLicense}
                    disabled={addingLicense || !newLicensePrice}
                    className="rounded-lg bg-zinc-900 px-3 py-1.5 text-xs font-medium text-white transition-colors hover:bg-zinc-700 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-300"
                  >
                    {addingLicense ? t("vault.detail.adding") : t("vault.detail.addLicenseButton")}
                  </button>
                  <button
                    onClick={() => setShowAddLicense(false)}
                    className="rounded-lg border border-zinc-300 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
                  >
                    {t("vault.detail.cancel")}
                  </button>
                </div>
              </div>
            )}

            {/* Existing Offerings */}
            {offerings.length === 0 ? (
              <p className="mt-4 text-sm text-zinc-500">{t("vault.detail.noLicenses")}</p>
            ) : (
              <div className="mt-4 space-y-3">
                {offerings.map((offering) => (
                  <div
                    key={offering.id}
                    className="flex items-center justify-between rounded-lg border border-zinc-200 p-3 dark:border-zinc-700"
                  >
                    <div>
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium capitalize text-zinc-900 dark:text-zinc-100">
                          {offering.licenseType.replace("_", " ")}
                        </span>
                        <span className="text-sm font-semibold text-zinc-900 dark:text-zinc-100">
                          ${(offering.priceCents / 100).toFixed(2)}
                        </span>
                      </div>
                      {offering.termsText && (
                        <p className="mt-0.5 text-xs text-zinc-500">{offering.termsText}</p>
                      )}
                    </div>
                    <button
                      onClick={() => handleDeleteLicense(offering.id)}
                      className="rounded-lg p-1.5 text-zinc-400 transition-colors hover:bg-red-50 hover:text-red-500 dark:hover:bg-red-950"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Proof of Ownership */}
          <div className="mt-6 rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
            <div className="flex items-center gap-2">
              <Shield className="h-5 w-5 text-zinc-500" />
              <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
                {t("vault.detail.proofOfOwnership")}
              </h2>
            </div>
            {proof ? (
              <div className="mt-4 space-y-2 text-sm">
                <div>
                  <span className="font-medium text-zinc-700 dark:text-zinc-300">
                    {t("vault.detail.sha256")}:
                  </span>
                  <code className="ml-2 break-all rounded bg-zinc-100 px-2 py-0.5 text-xs dark:bg-zinc-800">
                    {proof.hashSha256}
                  </code>
                </div>
                <div>
                  <span className="font-medium text-zinc-700 dark:text-zinc-300">
                    {t("vault.detail.registeredAt")}:
                  </span>
                  <span className="ml-2 text-zinc-600 dark:text-zinc-400">
                    {new Date(proof.createdAt).toLocaleString()}
                  </span>
                </div>
              </div>
            ) : (
              <div className="mt-4 space-y-2 text-sm">
                <div>
                  <span className="font-medium text-zinc-700 dark:text-zinc-300">
                    {t("vault.detail.sha256")}:
                  </span>
                  <code className="ml-2 break-all rounded bg-zinc-100 px-2 py-0.5 text-xs dark:bg-zinc-800">
                    {item.hashSha256}
                  </code>
                </div>
                <div>
                  <span className="font-medium text-zinc-700 dark:text-zinc-300">
                    {t("vault.detail.uploadedAt")}:
                  </span>
                  <span className="ml-2 text-zinc-600 dark:text-zinc-400">
                    {new Date(item.createdAt).toLocaleString()}
                  </span>
                </div>
              </div>
            )}
          </div>

          {/* Blockchain Anchor */}
          <div className="mt-6 rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
            <div className="flex items-center gap-2">
              <Shield className="h-5 w-5 text-emerald-500" />
              <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
                {t("blockchain.chain")}
              </h2>
            </div>

            {anchor ? (
              <div className="mt-4">
                {/* Status badge */}
                <div className="mb-4">
                  {anchor.anchorStatus === "confirmed" ? (
                    <span className="inline-flex items-center gap-1.5 rounded-full bg-emerald-50 px-3 py-1 text-sm font-medium text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300">
                      <CheckCircle className="h-4 w-4" />
                      {t("blockchain.verified")}
                    </span>
                  ) : anchor.anchorStatus === "pending" ? (
                    <span className="inline-flex items-center gap-1.5 rounded-full bg-amber-50 px-3 py-1 text-sm font-medium text-amber-700 dark:bg-amber-950 dark:text-amber-300">
                      <Clock className="h-4 w-4" />
                      {t("blockchain.pending")}
                    </span>
                  ) : (
                    <span className="inline-flex items-center gap-1.5 rounded-full bg-red-50 px-3 py-1 text-sm font-medium text-red-700 dark:bg-red-950 dark:text-red-300">
                      <AlertTriangle className="h-4 w-4" />
                      {t("blockchain.failed")}
                    </span>
                  )}
                </div>

                {/* Anchor details */}
                <div className="space-y-2 text-sm">
                  {anchor.txHash && (
                    <div>
                      <span className="font-medium text-zinc-700 dark:text-zinc-300">
                        {t("blockchain.txHash")}:
                      </span>
                      <div className="mt-1 flex items-center gap-2">
                        <code className="break-all rounded bg-zinc-100 px-2 py-0.5 text-xs dark:bg-zinc-800">
                          {anchor.txHash}
                        </code>
                        {anchor.chain === "polygon" && !anchor.txHash.startsWith("0x000000") && (
                          <a
                            href={`https://polygonscan.com/tx/${anchor.txHash}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="flex-shrink-0 text-emerald-600 hover:text-emerald-700"
                          >
                            <ExternalLink className="h-3.5 w-3.5" />
                          </a>
                        )}
                      </div>
                    </div>
                  )}
                  <div className="flex gap-6">
                    <div>
                      <span className="font-medium text-zinc-700 dark:text-zinc-300">
                        {t("blockchain.chain")}:
                      </span>
                      <span className="ml-2 capitalize text-zinc-600 dark:text-zinc-400">
                        {anchor.chain}
                      </span>
                    </div>
                    {anchor.blockNumber && (
                      <div>
                        <span className="font-medium text-zinc-700 dark:text-zinc-300">
                          {t("blockchain.blockNumber")}:
                        </span>
                        <span className="ml-2 text-zinc-600 dark:text-zinc-400">
                          {anchor.blockNumber.toLocaleString()}
                        </span>
                      </div>
                    )}
                  </div>
                  {anchor.confirmedAt && (
                    <div>
                      <span className="font-medium text-zinc-700 dark:text-zinc-300">
                        {t("blockchain.anchoredAt")}:
                      </span>
                      <span className="ml-2 text-zinc-600 dark:text-zinc-400">
                        {new Date(anchor.confirmedAt).toLocaleString()}
                      </span>
                    </div>
                  )}
                </div>

                {/* Simulated notice */}
                {anchor.contractAddress === "0x0000000000000000000000000000000000000000" && (
                  <div className="mt-3 rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-700 dark:border-amber-800 dark:bg-amber-950 dark:text-amber-300">
                    {t("blockchain.simulated")}
                  </div>
                )}

                {/* Verify link */}
                <div className="mt-3">
                  <Link
                    href={`/verify?hash=${anchor.contentHash}`}
                    className="text-xs font-medium text-emerald-600 hover:text-emerald-700 dark:text-emerald-400"
                  >
                    {t("blockchain.proofCertificate")} &rarr;
                  </Link>
                </div>
              </div>
            ) : (
              <div className="mt-4">
                <p className="mb-3 text-sm text-zinc-500">
                  Anchor this content&apos;s SHA-256 hash on the blockchain for immutable proof of existence.
                </p>
                <button
                  onClick={handleAnchor}
                  disabled={anchoring}
                  className="flex items-center gap-2 rounded-lg bg-emerald-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-emerald-700 disabled:opacity-50 dark:bg-emerald-500 dark:hover:bg-emerald-600"
                >
                  <Shield className="h-4 w-4" />
                  {anchoring ? t("blockchain.anchoring") : t("blockchain.anchor")}
                </button>
              </div>
            )}
          </div>
        </div>

        {/* Sidebar - right column */}
        <div className="space-y-4">
          {/* Metadata card */}
          <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900">
            <h3 className="text-sm font-semibold text-zinc-900 dark:text-zinc-100">
              {t("vault.detail.details")}
            </h3>
            <dl className="mt-3 space-y-2.5 text-sm">
              <div className="flex justify-between">
                <dt className="text-zinc-500">{t("vault.detail.type")}</dt>
                <dd className="font-medium capitalize text-zinc-900 dark:text-zinc-100">{item.contentType}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-zinc-500">{t("vault.detail.size")}</dt>
                <dd className="font-medium text-zinc-900 dark:text-zinc-100">
                  {formatFileSize(item.fileSize)}
                </dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-zinc-500">{t("vault.detail.created")}</dt>
                <dd className="font-medium text-zinc-900 dark:text-zinc-100">
                  {new Date(item.createdAt).toLocaleDateString()}
                </dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-zinc-500">{t("vault.detail.visibility")}</dt>
                <dd className="flex items-center gap-1 font-medium text-zinc-900 dark:text-zinc-100">
                  {item.isPublic ? (
                    <>
                      <Globe className="h-3.5 w-3.5 text-green-500" />
                      {t("vault.detail.public")}
                    </>
                  ) : (
                    <>
                      <Lock className="h-3.5 w-3.5 text-zinc-400" />
                      {t("vault.detail.private")}
                    </>
                  )}
                </dd>
              </div>
            </dl>
          </div>

          {/* Actions */}
          <div className="space-y-2">
            <button
              onClick={handleTogglePublic}
              className="flex w-full items-center justify-center gap-2 rounded-lg border border-zinc-300 px-4 py-2 text-sm font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
            >
              {item.isPublic ? (
                <>
                  <Lock className="h-4 w-4" />
                  {t("vault.detail.makePrivate")}
                </>
              ) : (
                <>
                  <Globe className="h-4 w-4" />
                  {t("vault.detail.makePublic")}
                </>
              )}
            </button>

            <a
              href={`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/content/${item.id}/download`}
              className="flex w-full items-center justify-center gap-2 rounded-lg border border-zinc-300 px-4 py-2 text-sm font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
            >
              <Download className="h-4 w-4" />
              {t("vault.detail.download")}
            </a>

            {showDeleteConfirm ? (
              <div className="rounded-lg border border-red-200 bg-red-50 p-3 dark:border-red-800 dark:bg-red-950">
                <p className="text-sm text-red-700 dark:text-red-300">
                  {t("vault.detail.deleteConfirm")}
                </p>
                <div className="mt-2 flex gap-2">
                  <button
                    onClick={handleDelete}
                    disabled={deleting}
                    className="rounded-lg bg-red-600 px-3 py-1.5 text-xs font-medium text-white transition-colors hover:bg-red-700 disabled:opacity-50"
                  >
                    {deleting ? t("vault.detail.deleting") : t("vault.detail.confirmDelete")}
                  </button>
                  <button
                    onClick={() => setShowDeleteConfirm(false)}
                    className="rounded-lg border border-zinc-300 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
                  >
                    {t("vault.detail.cancel")}
                  </button>
                </div>
              </div>
            ) : (
              <button
                onClick={() => setShowDeleteConfirm(true)}
                className="flex w-full items-center justify-center gap-2 rounded-lg border border-red-200 px-4 py-2 text-sm font-medium text-red-600 transition-colors hover:bg-red-50 dark:border-red-800 dark:text-red-400 dark:hover:bg-red-950"
              >
                <Trash2 className="h-4 w-4" />
                {t("vault.detail.delete")}
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default function VaultDetailPage() {
  return (
    <Suspense fallback={null}>
      <VaultDetailContent />
    </Suspense>
  );
}
