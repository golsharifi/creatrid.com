"use client";

import { useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth-context";
import { Plus, Trash2, Camera, Download, AlertTriangle, Shield, Copy, Check } from "@/components/icons";
import type { CustomLink, EmailPrefs } from "@/lib/types";
import { useTranslation } from "react-i18next";

const THEME_KEYS = [
  { key: "default", colors: "bg-zinc-900 dark:bg-zinc-100" },
  { key: "ocean", colors: "bg-blue-600" },
  { key: "sunset", colors: "bg-orange-500" },
  { key: "forest", colors: "bg-emerald-600" },
  { key: "midnight", colors: "bg-indigo-700" },
  { key: "rose", colors: "bg-pink-500" },
];

const THEME_LABEL_MAP: Record<string, string> = {
  default: "settings.themeDefault",
  ocean: "settings.themeOcean",
  sunset: "settings.themeSunset",
  forest: "settings.themeForest",
  midnight: "settings.themeMidnight",
  rose: "settings.themeRose",
};

interface SettingsFormProps {
  name: string;
  username: string;
  bio: string;
  image: string;
  theme: string;
  customLinks: CustomLink[];
  emailPrefs: EmailPrefs;
  totpEnabled: boolean;
}

export function SettingsForm({
  name,
  username,
  bio,
  image,
  theme,
  customLinks,
  emailPrefs: initialEmailPrefs,
  totpEnabled: initialTotpEnabled,
}: SettingsFormProps) {
  const router = useRouter();
  const { refresh } = useAuth();
  const { t } = useTranslation();
  const [formData, setFormData] = useState({ name, username, bio, theme });
  const [links, setLinks] = useState<CustomLink[]>(customLinks || []);
  const [emailPrefs, setEmailPrefs] = useState<EmailPrefs>(initialEmailPrefs || { welcome: true, connectionAlert: true, weeklyDigest: true, collaborations: true });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [exporting, setExporting] = useState(false);
  const [imagePreview, setImagePreview] = useState(image || "");
  const [uploading, setUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // 2FA state
  const [twoFAEnabled, setTwoFAEnabled] = useState(initialTotpEnabled);
  const [twoFASetup, setTwoFASetup] = useState<{ secret: string; qrUrl: string } | null>(null);
  const [twoFACode, setTwoFACode] = useState("");
  const [twoFABackupCodes, setTwoFABackupCodes] = useState<string[] | null>(null);
  const [twoFALoading, setTwoFALoading] = useState(false);
  const [twoFAError, setTwoFAError] = useState("");
  const [twoFADisableCode, setTwoFADisableCode] = useState("");
  const [showDisable2FA, setShowDisable2FA] = useState(false);
  const [backupCodesCopied, setBackupCodesCopied] = useState(false);

  async function handleImageUpload(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!["image/jpeg", "image/png", "image/webp"].includes(file.type)) {
      setError(t("settings.imageErrorType"));
      return;
    }
    if (file.size > 5 * 1024 * 1024) {
      setError(t("settings.imageErrorSize"));
      return;
    }

    setUploading(true);
    setError("");

    const result = await api.users.uploadImage(file);
    if (result.error) {
      setError(result.error);
    } else if (result.data) {
      setImagePreview(result.data.image);
      await refresh();
    }
    setUploading(false);
  }

  const addLink = () => {
    if (links.length >= 10) return;
    setLinks([...links, { title: "", url: "" }]);
  };

  const removeLink = (index: number) => {
    setLinks(links.filter((_, i) => i !== index));
  };

  const updateLink = (index: number, field: "title" | "url", value: string) => {
    const updated = [...links];
    updated[index] = { ...updated[index], [field]: value };
    setLinks(updated);
  };

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);

    // Filter out empty links
    const validLinks = links.filter((l) => l.title.trim() && l.url.trim());

    const result = await api.users.updateProfile({
      name: formData.name.trim(),
      username: formData.username.toLowerCase().trim(),
      bio: formData.bio.trim(),
      theme: formData.theme,
      customLinks: validLinks,
      emailPrefs,
    });

    if (result.error) {
      setError(result.error);
      setLoading(false);
    } else {
      await refresh();
      router.push("/dashboard");
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-8">
      {/* Profile Section */}
      <div className="space-y-6">
        <h2 className="text-lg font-semibold">{t("settings.profileSection")}</h2>

        {/* Avatar Upload */}
        <div className="flex items-center gap-4">
          <div className="relative">
            {imagePreview ? (
              <img
                src={imagePreview}
                alt="Profile"
                className="h-20 w-20 rounded-full object-cover"
              />
            ) : (
              <div className="flex h-20 w-20 items-center justify-center rounded-full bg-zinc-100 text-2xl font-bold text-zinc-400 dark:bg-zinc-800">
                {(name || "?")[0].toUpperCase()}
              </div>
            )}
            <button
              type="button"
              onClick={() => fileInputRef.current?.click()}
              disabled={uploading}
              className="absolute -bottom-1 -right-1 rounded-full border-2 border-white bg-zinc-900 p-1.5 text-white transition-colors hover:bg-zinc-700 dark:border-zinc-900 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-300"
            >
              <Camera className="h-3.5 w-3.5" />
            </button>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/jpeg,image/png,image/webp"
              onChange={handleImageUpload}
              className="hidden"
            />
          </div>
          <div className="text-sm">
            <p className="font-medium">{uploading ? t("settings.uploading") : t("settings.profilePhoto")}</p>
            <p className="text-xs text-zinc-500">{t("settings.photoHint")}</p>
          </div>
        </div>

        <div>
          <label
            htmlFor="name"
            className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
          >
            {t("settings.displayName")}
          </label>
          <input
            id="name"
            type="text"
            value={formData.name}
            onChange={(e) =>
              setFormData({ ...formData, name: e.target.value })
            }
            className="mt-1 block w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm focus:border-zinc-400 focus:outline-none focus:ring-1 focus:ring-zinc-400 dark:border-zinc-700 dark:bg-zinc-800"
            required
          />
        </div>

        <div>
          <label
            htmlFor="username"
            className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
          >
            {t("settings.username")}
          </label>
          <div className="mt-1 flex items-center rounded-lg border border-zinc-200 bg-white dark:border-zinc-700 dark:bg-zinc-800">
            <span className="pl-3 text-sm text-zinc-400">{t("settings.usernamePrefix")}</span>
            <input
              id="username"
              type="text"
              value={formData.username}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  username: e.target.value.replace(/[^a-zA-Z0-9_-]/g, ""),
                })
              }
              className="block w-full bg-transparent px-1 py-2 text-sm focus:outline-none"
              required
              minLength={3}
              maxLength={30}
            />
          </div>
        </div>

        <div>
          <label
            htmlFor="bio"
            className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
          >
            {t("settings.bio")}
          </label>
          <textarea
            id="bio"
            value={formData.bio}
            onChange={(e) =>
              setFormData({ ...formData, bio: e.target.value })
            }
            rows={4}
            maxLength={500}
            className="mt-1 block w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm focus:border-zinc-400 focus:outline-none focus:ring-1 focus:ring-zinc-400 dark:border-zinc-700 dark:bg-zinc-800"
            placeholder={t("settings.bioPlaceholder")}
          />
          <p className="mt-1 text-xs text-zinc-400">
            {t("settings.bioCharCount", { count: formData.bio.length })}
          </p>
        </div>
      </div>

      {/* Theme Section */}
      <div className="space-y-4">
        <h2 className="text-lg font-semibold">{t("settings.themeSection")}</h2>
        <p className="text-sm text-zinc-500">
          {t("settings.themeDescription")}
        </p>
        <div className="grid grid-cols-3 gap-3 sm:grid-cols-6">
          {THEME_KEYS.map((thm) => (
            <button
              key={thm.key}
              type="button"
              onClick={() => setFormData({ ...formData, theme: thm.key })}
              className={`flex flex-col items-center gap-2 rounded-lg border p-3 transition-colors ${
                formData.theme === thm.key
                  ? "border-zinc-900 dark:border-zinc-100"
                  : "border-zinc-200 dark:border-zinc-700"
              }`}
            >
              <div className={`h-8 w-8 rounded-full ${thm.colors}`} />
              <span className="text-xs font-medium">{t(THEME_LABEL_MAP[thm.key])}</span>
            </button>
          ))}
        </div>
      </div>

      {/* Custom Links Section */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-lg font-semibold">{t("settings.customLinksSection")}</h2>
            <p className="text-sm text-zinc-500">
              {t("settings.customLinksDescription")}
            </p>
          </div>
          {links.length < 10 && (
            <button
              type="button"
              onClick={addLink}
              className="flex items-center gap-1.5 rounded-lg border border-zinc-200 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
            >
              <Plus className="h-3.5 w-3.5" />
              {t("settings.addLink")}
            </button>
          )}
        </div>

        {links.length === 0 && (
          <p className="rounded-lg border border-dashed border-zinc-200 p-6 text-center text-sm text-zinc-400 dark:border-zinc-700">
            {t("settings.noLinksYet")}
          </p>
        )}

        <div className="space-y-3">
          {links.map((link, index) => (
            <div
              key={index}
              className="flex items-start gap-2 rounded-lg border border-zinc-200 p-3 dark:border-zinc-700"
            >
              <div className="flex-1 space-y-2">
                <input
                  type="text"
                  placeholder={t("settings.linkTitle")}
                  value={link.title}
                  onChange={(e) => updateLink(index, "title", e.target.value)}
                  maxLength={50}
                  className="block w-full rounded-md border border-zinc-200 bg-white px-2.5 py-1.5 text-sm focus:border-zinc-400 focus:outline-none dark:border-zinc-600 dark:bg-zinc-800"
                />
                <input
                  type="url"
                  placeholder={t("settings.linkUrl")}
                  value={link.url}
                  onChange={(e) => updateLink(index, "url", e.target.value)}
                  className="block w-full rounded-md border border-zinc-200 bg-white px-2.5 py-1.5 text-sm focus:border-zinc-400 focus:outline-none dark:border-zinc-600 dark:bg-zinc-800"
                />
              </div>
              <button
                type="button"
                onClick={() => removeLink(index)}
                className="mt-1 rounded p-1 text-zinc-400 hover:text-red-500"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            </div>
          ))}
        </div>
      </div>

      {/* Notification Preferences */}
      <div className="space-y-4">
        <h2 className="text-lg font-semibold">{t("settings.emailSection")}</h2>
        <p className="text-sm text-zinc-500">{t("settings.emailDescription")}</p>
        <div className="space-y-3">
          {([
            { key: "welcome" as const, label: t("settings.emailWelcome"), description: t("settings.emailWelcomeDesc") },
            { key: "connectionAlert" as const, label: t("settings.emailConnection"), description: t("settings.emailConnectionDesc") },
            { key: "weeklyDigest" as const, label: t("settings.emailDigest"), description: t("settings.emailDigestDesc") },
            { key: "collaborations" as const, label: t("settings.emailCollab"), description: t("settings.emailCollabDesc") },
          ]).map((pref) => (
            <label key={pref.key} className="flex items-center justify-between rounded-lg border border-zinc-200 p-3 dark:border-zinc-700">
              <div>
                <p className="text-sm font-medium">{pref.label}</p>
                <p className="text-xs text-zinc-500">{pref.description}</p>
              </div>
              <button
                type="button"
                role="switch"
                aria-checked={emailPrefs[pref.key]}
                onClick={() => setEmailPrefs({ ...emailPrefs, [pref.key]: !emailPrefs[pref.key] })}
                className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors ${
                  emailPrefs[pref.key] ? "bg-zinc-900 dark:bg-zinc-100" : "bg-zinc-200 dark:bg-zinc-700"
                }`}
              >
                <span
                  className={`pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow transition-transform dark:bg-zinc-900 ${
                    emailPrefs[pref.key] ? "translate-x-5" : "translate-x-0"
                  }`}
                />
              </button>
            </label>
          ))}
        </div>
      </div>

      {/* Security: Two-Factor Authentication */}
      <div className="space-y-4">
        <div className="flex items-center gap-2">
          <Shield className="h-5 w-5" />
          <h2 className="text-lg font-semibold">{t("twoFA.title")}</h2>
        </div>
        <p className="text-sm text-zinc-500">{t("twoFA.subtitle")}</p>

        {twoFABackupCodes ? (
          // Show backup codes after successful setup
          <div className="space-y-3">
            <p className="text-sm font-medium text-green-600 dark:text-green-400">{t("twoFA.enabled")}</p>
            <p className="text-sm text-zinc-600 dark:text-zinc-400">{t("twoFA.backupCodesDesc")}</p>
            <div className="rounded-lg border border-zinc-200 bg-zinc-50 p-4 dark:border-zinc-700 dark:bg-zinc-800/50">
              <div className="mb-2 flex items-center justify-between">
                <p className="text-sm font-medium">{t("twoFA.backupCodes")}</p>
                <button
                  type="button"
                  onClick={() => {
                    navigator.clipboard.writeText(twoFABackupCodes.join("\n"));
                    setBackupCodesCopied(true);
                    setTimeout(() => setBackupCodesCopied(false), 2000);
                  }}
                  className="flex items-center gap-1 text-xs text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300"
                >
                  {backupCodesCopied ? <Check className="h-3 w-3" /> : <Copy className="h-3 w-3" />}
                  {backupCodesCopied ? t("common.copied") : t("common.copy")}
                </button>
              </div>
              <div className="grid grid-cols-2 gap-1 font-mono text-sm">
                {twoFABackupCodes.map((code, i) => (
                  <span key={i} className="text-zinc-700 dark:text-zinc-300">{code}</span>
                ))}
              </div>
            </div>
            <button
              type="button"
              onClick={() => {
                setTwoFABackupCodes(null);
                setTwoFAEnabled(true);
                setTwoFASetup(null);
                setTwoFACode("");
                refresh();
              }}
              className="rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              {t("common.save")}
            </button>
          </div>
        ) : twoFASetup ? (
          // Show QR code and verification input
          <div className="space-y-4">
            <p className="text-sm">{t("twoFA.scanQR")}</p>
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(twoFASetup.qrUrl)}`}
              alt="2FA QR Code"
              className="mx-auto h-48 w-48 rounded-lg border border-zinc-200 dark:border-zinc-700"
              width={200}
              height={200}
            />
            <div>
              <p className="mb-1 text-xs text-zinc-500">{t("twoFA.secretKey")}</p>
              <code className="block break-all rounded bg-zinc-100 px-3 py-2 text-xs dark:bg-zinc-800">
                {twoFASetup.secret}
              </code>
            </div>
            <div>
              <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
                {t("twoFA.enterCode")}
              </label>
              <input
                type="text"
                inputMode="numeric"
                pattern="[0-9]*"
                maxLength={6}
                value={twoFACode}
                onChange={(e) => setTwoFACode(e.target.value.replace(/[^0-9]/g, ""))}
                placeholder="000000"
                className="mt-1 block w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-center text-lg tracking-[0.5em] focus:border-zinc-400 focus:outline-none focus:ring-1 focus:ring-zinc-400 dark:border-zinc-700 dark:bg-zinc-800"
              />
            </div>
            {twoFAError && <p className="text-sm text-red-500">{twoFAError}</p>}
            <div className="flex gap-2">
              <button
                type="button"
                disabled={twoFALoading || twoFACode.length !== 6}
                onClick={async () => {
                  setTwoFALoading(true);
                  setTwoFAError("");
                  const result = await api.auth.confirm2FA(twoFACode);
                  if (result.error) {
                    setTwoFAError(result.error);
                  } else if (result.data) {
                    setTwoFABackupCodes(result.data.backupCodes);
                  }
                  setTwoFALoading(false);
                }}
                className="rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
              >
                {twoFALoading ? t("common.loading") : t("twoFA.verify")}
              </button>
              <button
                type="button"
                onClick={() => {
                  setTwoFASetup(null);
                  setTwoFACode("");
                  setTwoFAError("");
                }}
                className="rounded-lg border border-zinc-200 px-4 py-2 text-sm font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
              >
                {t("common.cancel")}
              </button>
            </div>
          </div>
        ) : twoFAEnabled ? (
          // 2FA is enabled — show status and disable option
          <div className="space-y-3">
            <div className="flex items-center gap-2 rounded-lg border border-green-200 bg-green-50 p-3 dark:border-green-800 dark:bg-green-950">
              <Shield className="h-4 w-4 text-green-600 dark:text-green-400" />
              <p className="text-sm font-medium text-green-700 dark:text-green-300">{t("twoFA.enabled")}</p>
            </div>
            {showDisable2FA ? (
              <div className="space-y-3">
                <div>
                  <label className="block text-sm font-medium text-zinc-700 dark:text-zinc-300">
                    {t("twoFA.enterCode")}
                  </label>
                  <input
                    type="text"
                    inputMode="numeric"
                    pattern="[0-9]*"
                    maxLength={6}
                    value={twoFADisableCode}
                    onChange={(e) => setTwoFADisableCode(e.target.value.replace(/[^0-9]/g, ""))}
                    placeholder="000000"
                    className="mt-1 block w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-center text-lg tracking-[0.5em] focus:border-zinc-400 focus:outline-none focus:ring-1 focus:ring-zinc-400 dark:border-zinc-700 dark:bg-zinc-800"
                  />
                </div>
                {twoFAError && <p className="text-sm text-red-500">{twoFAError}</p>}
                <div className="flex gap-2">
                  <button
                    type="button"
                    disabled={twoFALoading || twoFADisableCode.length !== 6}
                    onClick={async () => {
                      setTwoFALoading(true);
                      setTwoFAError("");
                      const result = await api.auth.disable2FA(twoFADisableCode);
                      if (result.error) {
                        setTwoFAError(result.error);
                      } else {
                        setTwoFAEnabled(false);
                        setShowDisable2FA(false);
                        setTwoFADisableCode("");
                        await refresh();
                      }
                      setTwoFALoading(false);
                    }}
                    className="rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-red-700 disabled:cursor-not-allowed disabled:opacity-50"
                  >
                    {twoFALoading ? t("common.loading") : t("twoFA.disable")}
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      setShowDisable2FA(false);
                      setTwoFADisableCode("");
                      setTwoFAError("");
                    }}
                    className="rounded-lg border border-zinc-200 px-4 py-2 text-sm font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
                  >
                    {t("common.cancel")}
                  </button>
                </div>
              </div>
            ) : (
              <button
                type="button"
                onClick={() => setShowDisable2FA(true)}
                className="rounded-lg border border-red-300 px-4 py-2 text-sm font-medium text-red-600 transition-colors hover:bg-red-50 dark:border-red-800 dark:text-red-400 dark:hover:bg-red-950"
              >
                {t("twoFA.disable")}
              </button>
            )}
          </div>
        ) : (
          // 2FA is disabled — show enable button
          <div className="space-y-3">
            <p className="text-sm text-zinc-500">{t("twoFA.disabled")}</p>
            <button
              type="button"
              disabled={twoFALoading}
              onClick={async () => {
                setTwoFALoading(true);
                setTwoFAError("");
                const result = await api.auth.setup2FA();
                if (result.error) {
                  setTwoFAError(result.error);
                } else if (result.data) {
                  setTwoFASetup(result.data);
                }
                setTwoFALoading(false);
              }}
              className="rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              {twoFALoading ? t("common.loading") : t("twoFA.enable")}
            </button>
            {twoFAError && <p className="text-sm text-red-500">{twoFAError}</p>}
          </div>
        )}
      </div>

      {/* Export Data */}
      <div className="space-y-4">
        <h2 className="text-lg font-semibold">{t("settings.exportSection")}</h2>
        <p className="text-sm text-zinc-500">{t("settings.exportDescription")}</p>
        <button
          type="button"
          disabled={exporting}
          onClick={async () => {
            setExporting(true);
            const result = await api.users.exportProfile();
            if (result.data) {
              const blob = new Blob([JSON.stringify(result.data, null, 2)], { type: "application/json" });
              const url = URL.createObjectURL(blob);
              const a = document.createElement("a");
              a.href = url;
              a.download = "creatrid-profile.json";
              a.click();
              URL.revokeObjectURL(url);
            }
            setExporting(false);
          }}
          className="flex items-center gap-2 rounded-lg border border-zinc-200 px-3 py-2 text-sm font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
        >
          <Download className="h-4 w-4" />
          {exporting ? t("settings.exporting") : t("settings.exportButton")}
        </button>
      </div>

      {error && <p className="text-sm text-red-500">{error}</p>}

      <button
        type="submit"
        disabled={loading}
        className="rounded-lg bg-zinc-900 px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
      >
        {loading ? t("settings.saving") : t("settings.saveChanges")}
      </button>

      {/* Danger Zone: Delete Account */}
      <div className="space-y-4 rounded-lg border border-red-200 p-6 dark:border-red-900">
        <div className="flex items-center gap-2">
          <AlertTriangle className="h-5 w-5 text-red-500" />
          <h2 className="text-lg font-semibold text-red-600 dark:text-red-400">{t("settings.dangerZone")}</h2>
        </div>
        <p className="text-sm text-zinc-600 dark:text-zinc-400">
          {t("settings.dangerDescription")}
        </p>
        {showDeleteConfirm ? (
          <div className="flex items-center gap-3">
            <button
              type="button"
              disabled={deleting}
              onClick={async () => {
                setDeleting(true);
                const result = await api.users.deleteAccount();
                if (result.error) {
                  setError(result.error);
                  setDeleting(false);
                  setShowDeleteConfirm(false);
                } else {
                  window.location.href = "/";
                }
              }}
              className="rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-red-700"
            >
              {deleting ? t("settings.deleting") : t("settings.confirmDelete")}
            </button>
            <button
              type="button"
              onClick={() => setShowDeleteConfirm(false)}
              className="rounded-lg border border-zinc-200 px-4 py-2 text-sm font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
            >
              {t("common.cancel")}
            </button>
          </div>
        ) : (
          <button
            type="button"
            onClick={() => setShowDeleteConfirm(true)}
            className="rounded-lg border border-red-300 px-4 py-2 text-sm font-medium text-red-600 transition-colors hover:bg-red-50 dark:border-red-800 dark:text-red-400 dark:hover:bg-red-950"
          >
            {t("settings.deleteAccount")}
          </button>
        )}
      </div>
    </form>
  );
}
