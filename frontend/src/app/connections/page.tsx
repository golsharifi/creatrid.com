"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState, useCallback, Suspense } from "react";
import { api } from "@/lib/api";
import type { Connection } from "@/lib/types";
import { CheckCircle, ExternalLink } from "lucide-react";
import { useTranslation } from "react-i18next";

interface PlatformConfig {
  key: string;
  name: string;
  icon: string;
  available: boolean;
}

const PLATFORMS: PlatformConfig[] = [
  { key: "youtube", name: "YouTube", icon: "YT", available: true },
  { key: "github", name: "GitHub", icon: "GH", available: true },
  { key: "twitter", name: "Twitter / X", icon: "\u{1d54f}", available: true },
  { key: "linkedin", name: "LinkedIn", icon: "in", available: true },
  { key: "instagram", name: "Instagram", icon: "IG", available: true },
  { key: "behance", name: "Behance", icon: "Be", available: true },
  { key: "dribbble", name: "Dribbble", icon: "Dr", available: true },
];

function ConnectionsContent() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();
  const { t } = useTranslation();
  const [connections, setConnections] = useState<Connection[]>([]);
  const [loadingConnections, setLoadingConnections] = useState(true);
  const [disconnecting, setDisconnecting] = useState<string | null>(null);

  const connected = searchParams.get("connected");
  const error = searchParams.get("error");

  const fetchConnections = useCallback(async () => {
    const result = await api.connections.list();
    if (result.data) {
      setConnections(result.data.connections || []);
    }
    setLoadingConnections(false);
  }, []);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
    if (user) fetchConnections();
  }, [user, loading, router, fetchConnections]);

  const handleConnect = (platform: string) => {
    window.location.href = api.connections.connectUrl(platform);
  };

  const handleDisconnect = async (platform: string) => {
    setDisconnecting(platform);
    const result = await api.connections.disconnect(platform);
    if (result.data?.success) {
      setConnections((prev) => prev.filter((c) => c.platform !== platform));
    }
    setDisconnecting(null);
  };

  if (loading || !user) return null;

  const connectionMap = new Map(connections.map((c) => [c.platform, c]));

  return (
    <div className="mx-auto max-w-2xl px-4 py-12 sm:px-6">
      <h1 className="text-2xl font-bold">{t("connections.title")}</h1>
      <p className="mt-1 text-sm text-zinc-500">
        {t("connections.subtitle")}
      </p>

      {connected && (
        <div className="mt-4 rounded-lg bg-green-50 p-3 text-sm text-green-700 dark:bg-green-950 dark:text-green-300">
          <CheckCircle className="mr-2 inline h-4 w-4" />
          {t("connections.successConnected", { platform: PLATFORMS.find((p) => p.key === connected)?.name || connected })}
        </div>
      )}
      {error && (
        <div className="mt-4 rounded-lg bg-red-50 p-3 text-sm text-red-700 dark:bg-red-950 dark:text-red-300">
          {t("connections.connectionFailed", { error: error.replace(/_/g, " ") })}
        </div>
      )}

      <div className="mt-8 space-y-3">
        {PLATFORMS.map((platform) => {
          const conn = connectionMap.get(platform.key);
          const isConnected = !!conn;

          return (
            <div
              key={platform.key}
              className="flex items-center justify-between rounded-xl border border-zinc-200 p-4 dark:border-zinc-800"
            >
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-100 text-sm font-bold dark:bg-zinc-800">
                  {conn?.avatarUrl ? (
                    <img
                      src={conn.avatarUrl}
                      alt=""
                      className="h-10 w-10 rounded-lg object-cover"
                    />
                  ) : (
                    platform.icon
                  )}
                </div>
                <div>
                  <p className="text-sm font-medium">{platform.name}</p>
                  {isConnected ? (
                    <div className="flex items-center gap-2">
                      <p className="text-xs text-green-600 dark:text-green-400">
                        {conn.username || conn.displayName}
                      </p>
                      {conn.followerCount !== null && (
                        <span className="text-xs text-zinc-400">
                          {t("connections.followers", { val: conn.followerCount.toLocaleString() })}
                        </span>
                      )}
                    </div>
                  ) : (
                    <p className="text-xs text-zinc-400">{t("connections.notConnected")}</p>
                  )}
                </div>
              </div>
              <div className="flex items-center gap-2">
                {isConnected && conn.profileUrl && (
                  <a
                    href={conn.profileUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="rounded-lg p-1.5 text-zinc-400 hover:text-zinc-600"
                  >
                    <ExternalLink className="h-4 w-4" />
                  </a>
                )}
                {isConnected ? (
                  <button
                    onClick={() => handleDisconnect(platform.key)}
                    disabled={disconnecting === platform.key}
                    className="rounded-lg border border-red-200 px-3 py-1.5 text-xs font-medium text-red-600 transition-colors hover:bg-red-50 disabled:opacity-50 dark:border-red-800 dark:text-red-400 dark:hover:bg-red-950"
                  >
                    {disconnecting === platform.key ? t("connections.disconnecting") : t("connections.disconnect")}
                  </button>
                ) : platform.available ? (
                  <button
                    onClick={() => handleConnect(platform.key)}
                    className="rounded-lg border border-zinc-200 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
                  >
                    {t("connections.connect")}
                  </button>
                ) : (
                  <button
                    disabled
                    className="cursor-not-allowed rounded-lg border border-zinc-200 px-3 py-1.5 text-xs font-medium opacity-40 dark:border-zinc-700"
                  >
                    {t("connections.comingSoon")}
                  </button>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

export default function ConnectionsPage() {
  return (
    <Suspense fallback={null}>
      <ConnectionsContent />
    </Suspense>
  );
}
