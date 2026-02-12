"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Code, Image, Link2, Check, Copy } from "lucide-react";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
const PROFILE_BASE_URL =
  process.env.NEXT_PUBLIC_PROFILE_URL || "https://creatrid.com";

export default function WidgetPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();
  const [copiedIndex, setCopiedIndex] = useState<number | null>(null);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
    if (!loading && user && !user.onboarded) router.push("/onboarding");
  }, [user, loading, router]);

  const copyToClipboard = (text: string, index: number) => {
    navigator.clipboard.writeText(text);
    setCopiedIndex(index);
    setTimeout(() => setCopiedIndex(null), 2000);
  };

  if (loading || !user) return null;

  const username = user.username || "";

  const embedCodes = [
    {
      label: t("widget.htmlEmbed"),
      icon: Code,
      code: `<iframe src="${PROFILE_BASE_URL}/widget?u=${username}" width="220" height="80" frameborder="0"></iframe>`,
    },
    {
      label: t("widget.markdownBadge"),
      icon: Image,
      code: `![Creator Score](${API_URL}/api/widget/${username}/svg)`,
    },
    {
      label: t("widget.directLink"),
      icon: Link2,
      code: `${API_URL}/api/widget/${username}/svg`,
    },
  ];

  return (
    <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold">{t("widget.title")}</h1>
        <p className="mt-1 text-zinc-500">{t("widget.subtitle")}</p>
      </div>

      {/* Preview */}
      <div className="mb-8 rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
        <h2 className="mb-4 text-lg font-semibold">{t("widget.preview")}</h2>
        <div className="flex flex-col items-start gap-6 sm:flex-row sm:items-center">
          {/* SVG Badge Preview */}
          <div>
            <p className="mb-2 text-xs font-medium text-zinc-500">
              {t("widget.markdownBadge")}
            </p>
            <img
              src={`${API_URL}/api/widget/${username}/svg`}
              alt="Creator Score Badge"
              width={200}
              height={60}
            />
          </div>

          {/* HTML Widget Preview */}
          <div>
            <p className="mb-2 text-xs font-medium text-zinc-500">
              {t("widget.htmlEmbed")}
            </p>
            <iframe
              src={`${API_URL}/api/widget/${username}/html`}
              width={280}
              height={70}
              frameBorder={0}
              className="rounded-lg"
              title="Creator Widget Preview"
            />
          </div>
        </div>
      </div>

      {/* Embed Codes */}
      <div className="rounded-xl border border-zinc-200 p-6 dark:border-zinc-800">
        <h2 className="mb-4 text-lg font-semibold">{t("widget.embedCode")}</h2>
        <div className="space-y-4">
          {embedCodes.map((item, index) => (
            <div key={index}>
              <div className="mb-1.5 flex items-center gap-2">
                <item.icon className="h-4 w-4 text-zinc-500" />
                <span className="text-sm font-medium">{item.label}</span>
              </div>
              <div className="flex items-stretch gap-2">
                <div className="flex-1 overflow-x-auto rounded-lg border border-zinc-200 bg-zinc-50 px-3 py-2.5 font-mono text-xs text-zinc-700 dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-300">
                  <code className="whitespace-nowrap">{item.code}</code>
                </div>
                <button
                  onClick={() => copyToClipboard(item.code, index)}
                  className="flex items-center gap-1.5 rounded-lg border border-zinc-200 px-3 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
                >
                  {copiedIndex === index ? (
                    <>
                      <Check className="h-3.5 w-3.5 text-green-500" />
                      {t("widget.copied")}
                    </>
                  ) : (
                    <>
                      <Copy className="h-3.5 w-3.5" />
                      {t("widget.copy")}
                    </>
                  )}
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
