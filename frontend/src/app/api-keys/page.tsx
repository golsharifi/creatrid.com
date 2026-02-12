"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState, useCallback } from "react";
import { api } from "@/lib/api";
import { Key, Trash2, Copy, Check, BookOpen } from "lucide-react";
import { useTranslation } from "react-i18next";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

interface APIKey {
  id: string;
  name: string;
  keyPrefix: string;
  lastUsedAt: string | null;
  createdAt: string;
}

export default function APIKeysPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();
  const [keys, setKeys] = useState<APIKey[]>([]);
  const [loadingKeys, setLoadingKeys] = useState(true);
  const [creating, setCreating] = useState(false);
  const [newKeyName, setNewKeyName] = useState("");
  const [newKey, setNewKey] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [revoking, setRevoking] = useState<string | null>(null);

  const fetchKeys = useCallback(async () => {
    const result = await api.apiKeys.list();
    if (result.data) {
      setKeys(result.data.keys || []);
    }
    setLoadingKeys(false);
  }, []);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
    if (user) fetchKeys();
  }, [user, loading, router, fetchKeys]);

  const handleCreate = async () => {
    if (!newKeyName.trim()) return;
    setCreating(true);
    const result = await api.apiKeys.create(newKeyName.trim());
    if (result.data) {
      setNewKey(result.data.key);
      setNewKeyName("");
      fetchKeys();
    }
    setCreating(false);
  };

  const handleCopy = async (text: string) => {
    await navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleRevoke = async (id: string) => {
    setRevoking(id);
    const result = await api.apiKeys.delete(id);
    if (result.data?.success) {
      setKeys((prev) => prev.filter((k) => k.id !== id));
    }
    setRevoking(null);
  };

  if (loading || !user) return null;

  return (
    <div className="mx-auto max-w-3xl px-4 py-12 sm:px-6">
      <div className="flex items-center gap-3">
        <Key className="h-6 w-6 text-zinc-400" />
        <div>
          <h1 className="text-2xl font-bold">{t("apiKeys.title")}</h1>
          <p className="mt-1 text-sm text-zinc-500">
            {t("apiKeys.subtitle")}
          </p>
        </div>
      </div>

      {/* Create new key */}
      <div className="mt-8 rounded-xl border border-zinc-200 p-5 dark:border-zinc-800">
        <h2 className="text-sm font-semibold">{t("apiKeys.createKey")}</h2>
        <div className="mt-3 flex gap-3">
          <input
            type="text"
            value={newKeyName}
            onChange={(e) => setNewKeyName(e.target.value)}
            placeholder={t("apiKeys.keyNamePlaceholder")}
            maxLength={100}
            className="flex-1 rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm outline-none focus:border-zinc-400 dark:border-zinc-700 dark:bg-zinc-900 dark:focus:border-zinc-500"
            onKeyDown={(e) => {
              if (e.key === "Enter") handleCreate();
            }}
          />
          <button
            onClick={handleCreate}
            disabled={creating || !newKeyName.trim()}
            className="rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
          >
            {creating ? t("apiKeys.creating") : t("apiKeys.createKey")}
          </button>
        </div>
      </div>

      {/* Newly created key banner */}
      {newKey && (
        <div className="mt-4 rounded-xl border border-green-200 bg-green-50 p-5 dark:border-green-800 dark:bg-green-950">
          <h3 className="text-sm font-semibold text-green-800 dark:text-green-200">
            {t("apiKeys.newKeyTitle")}
          </h3>
          <div className="mt-2 flex items-center gap-2">
            <code className="flex-1 break-all rounded-lg bg-green-100 px-3 py-2 text-xs font-mono text-green-900 dark:bg-green-900 dark:text-green-100">
              {newKey}
            </code>
            <button
              onClick={() => handleCopy(newKey)}
              className="rounded-lg border border-green-300 px-3 py-2 text-xs font-medium text-green-700 transition-colors hover:bg-green-100 dark:border-green-700 dark:text-green-300 dark:hover:bg-green-900"
            >
              {copied ? (
                <span className="flex items-center gap-1">
                  <Check className="h-3 w-3" /> {t("apiKeys.copied")}
                </span>
              ) : (
                <span className="flex items-center gap-1">
                  <Copy className="h-3 w-3" /> {t("apiKeys.copy")}
                </span>
              )}
            </button>
          </div>
          <p className="mt-2 text-xs text-green-700 dark:text-green-400">
            {t("apiKeys.newKeyWarning")}
          </p>
        </div>
      )}

      {/* Existing keys list */}
      <div className="mt-8">
        {loadingKeys ? (
          <p className="text-sm text-zinc-400">{t("common.loading")}</p>
        ) : keys.length === 0 ? (
          <p className="text-sm text-zinc-400">{t("apiKeys.noKeysYet")}</p>
        ) : (
          <div className="space-y-3">
            {keys.map((key) => (
              <div
                key={key.id}
                className="flex items-center justify-between rounded-xl border border-zinc-200 p-4 dark:border-zinc-800"
              >
                <div className="min-w-0 flex-1">
                  <p className="text-sm font-medium">{key.name}</p>
                  <div className="mt-1 flex flex-wrap gap-x-4 gap-y-1 text-xs text-zinc-400">
                    <span>
                      {t("apiKeys.prefix")}: <code className="font-mono">{key.keyPrefix}...</code>
                    </span>
                    <span>
                      {t("apiKeys.lastUsed")}:{" "}
                      {key.lastUsedAt
                        ? new Date(key.lastUsedAt).toLocaleDateString()
                        : t("apiKeys.never")}
                    </span>
                    <span>
                      {t("apiKeys.created")}:{" "}
                      {new Date(key.createdAt).toLocaleDateString()}
                    </span>
                  </div>
                </div>
                <button
                  onClick={() => handleRevoke(key.id)}
                  disabled={revoking === key.id}
                  className="ml-3 rounded-lg border border-red-200 p-2 text-red-600 transition-colors hover:bg-red-50 disabled:opacity-50 dark:border-red-800 dark:text-red-400 dark:hover:bg-red-950"
                  title={t("apiKeys.revoke")}
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* API Documentation */}
      <div className="mt-12 rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
        <div className="flex items-center gap-2">
          <BookOpen className="h-5 w-5 text-zinc-400" />
          <h2 className="text-lg font-semibold">{t("apiKeys.docsTitle")}</h2>
        </div>

        <div className="mt-6 space-y-6">
          {/* Base URL */}
          <div>
            <h3 className="text-sm font-semibold text-zinc-600 dark:text-zinc-400">
              {t("apiKeys.docsBaseUrl")}
            </h3>
            <code className="mt-1 block rounded-lg bg-zinc-100 px-3 py-2 text-xs font-mono dark:bg-zinc-800">
              {API_URL}
            </code>
          </div>

          {/* Authentication */}
          <div>
            <h3 className="text-sm font-semibold text-zinc-600 dark:text-zinc-400">
              {t("apiKeys.docsAuth")}
            </h3>
            <p className="mt-1 text-xs text-zinc-500">{t("apiKeys.docsAuthDesc")}</p>
            <code className="mt-2 block rounded-lg bg-zinc-100 px-3 py-2 text-xs font-mono dark:bg-zinc-800">
              Authorization: Bearer crk_your_api_key_here
            </code>
          </div>

          {/* Endpoints */}
          <div>
            <h3 className="text-sm font-semibold text-zinc-600 dark:text-zinc-400">
              {t("apiKeys.docsEndpoints")}
            </h3>
            <div className="mt-3 space-y-4">
              {/* Verify endpoint */}
              <div className="rounded-lg border border-zinc-100 p-4 dark:border-zinc-700">
                <div className="flex items-center gap-2">
                  <span className="rounded bg-green-100 px-2 py-0.5 text-xs font-bold text-green-700 dark:bg-green-900 dark:text-green-300">
                    GET
                  </span>
                  <code className="text-xs font-mono">/api/v1/verify/&#123;username&#125;</code>
                </div>
                <p className="mt-1 text-xs text-zinc-500">{t("apiKeys.docsVerifyDesc")}</p>
                <pre className="mt-2 overflow-x-auto rounded-lg bg-zinc-50 p-3 text-xs font-mono text-zinc-700 dark:bg-zinc-900 dark:text-zinc-300">
{`curl -H "Authorization: Bearer crk_..." \\
  ${API_URL}/api/v1/verify/johndoe`}
                </pre>
              </div>

              {/* Score endpoint */}
              <div className="rounded-lg border border-zinc-100 p-4 dark:border-zinc-700">
                <div className="flex items-center gap-2">
                  <span className="rounded bg-green-100 px-2 py-0.5 text-xs font-bold text-green-700 dark:bg-green-900 dark:text-green-300">
                    GET
                  </span>
                  <code className="text-xs font-mono">/api/v1/verify/&#123;username&#125;/score</code>
                </div>
                <p className="mt-1 text-xs text-zinc-500">{t("apiKeys.docsScoreDesc")}</p>
                <pre className="mt-2 overflow-x-auto rounded-lg bg-zinc-50 p-3 text-xs font-mono text-zinc-700 dark:bg-zinc-900 dark:text-zinc-300">
{`curl -H "Authorization: Bearer crk_..." \\
  ${API_URL}/api/v1/verify/johndoe/score`}
                </pre>
              </div>

              {/* Search endpoint */}
              <div className="rounded-lg border border-zinc-100 p-4 dark:border-zinc-700">
                <div className="flex items-center gap-2">
                  <span className="rounded bg-green-100 px-2 py-0.5 text-xs font-bold text-green-700 dark:bg-green-900 dark:text-green-300">
                    GET
                  </span>
                  <code className="text-xs font-mono">/api/v1/search</code>
                </div>
                <p className="mt-1 text-xs text-zinc-500">{t("apiKeys.docsSearchDesc")}</p>
                <pre className="mt-2 overflow-x-auto rounded-lg bg-zinc-50 p-3 text-xs font-mono text-zinc-700 dark:bg-zinc-900 dark:text-zinc-300">
{`curl -H "Authorization: Bearer crk_..." \\
  "${API_URL}/api/v1/search?q=design&minScore=50&limit=10"`}
                </pre>
              </div>
            </div>
          </div>

          <p className="text-xs text-zinc-400">{t("apiKeys.docsRateLimit")}</p>
        </div>
      </div>
    </div>
  );
}
