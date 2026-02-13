"use client";

import { useState, Suspense } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { Shield } from "@/components/icons";
import { useTranslation } from "react-i18next";

function TwoFAContent() {
  const router = useRouter();
  const { t } = useTranslation();
  const [code, setCode] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [useBackup, setUseBackup] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!code.trim()) {
      setError(t("twoFA.codeRequired"));
      return;
    }

    setError("");
    setLoading(true);

    const result = await api.auth.verify2FA(code.trim());
    if (result.error) {
      setError(result.error);
      setLoading(false);
    } else {
      router.push("/dashboard");
    }
  }

  return (
    <div className="flex flex-1 items-center justify-center px-4 py-24">
      <div className="w-full max-w-sm space-y-8">
        <div className="text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-zinc-100 dark:bg-zinc-800">
            <Shield className="h-6 w-6" />
          </div>
          <h1 className="text-2xl font-bold">{t("twoFA.title")}</h1>
          <p className="mt-2 text-sm text-zinc-500">
            {t("twoFA.enterCode")}
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {useBackup ? (
            <div>
              <label
                htmlFor="backup-code"
                className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
              >
                {t("twoFA.backupCodes")}
              </label>
              <input
                id="backup-code"
                type="text"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                placeholder="00000000"
                className="mt-1 block w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm tracking-widest focus:border-zinc-400 focus:outline-none focus:ring-1 focus:ring-zinc-400 dark:border-zinc-700 dark:bg-zinc-800"
                autoFocus
              />
            </div>
          ) : (
            <div>
              <label
                htmlFor="totp-code"
                className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
              >
                {t("twoFA.enterCode")}
              </label>
              <input
                id="totp-code"
                type="text"
                inputMode="numeric"
                pattern="[0-9]*"
                maxLength={6}
                value={code}
                onChange={(e) =>
                  setCode(e.target.value.replace(/[^0-9]/g, ""))
                }
                placeholder="000000"
                className="mt-1 block w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-center text-lg tracking-[0.5em] focus:border-zinc-400 focus:outline-none focus:ring-1 focus:ring-zinc-400 dark:border-zinc-700 dark:bg-zinc-800"
                autoFocus
                autoComplete="one-time-code"
              />
            </div>
          )}

          {error && <p className="text-sm text-red-500">{error}</p>}

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-lg bg-zinc-900 px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
          >
            {loading ? t("common.loading") : t("twoFA.verify")}
          </button>

          <button
            type="button"
            onClick={() => {
              setUseBackup(!useBackup);
              setCode("");
              setError("");
            }}
            className="w-full text-center text-sm text-zinc-500 transition-colors hover:text-zinc-700 dark:hover:text-zinc-300"
          >
            {useBackup ? t("twoFA.enterCode") : t("twoFA.useBackupCode")}
          </button>
        </form>
      </div>
    </div>
  );
}

export default function TwoFAPage() {
  return (
    <Suspense fallback={null}>
      <TwoFAContent />
    </Suspense>
  );
}
