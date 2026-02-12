"use client";

import { useEffect, useState, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import { QRCodeSVG } from "qrcode.react";
import { api } from "@/lib/api";
import { Shield, CheckCircle, ExternalLink, QrCode } from "lucide-react";
import {
  CopyLinkButton,
  ShareTwitterButton,
  ShareLinkedInButton,
} from "@/components/share-buttons";
import type { PublicUser, Connection } from "@/lib/types";

const PROFILE_BASE_URL =
  process.env.NEXT_PUBLIC_PROFILE_URL || "https://creatrid.com";

const THEME_STYLES: Record<string, { accent: string; badge: string }> = {
  default: { accent: "bg-zinc-900 dark:bg-zinc-100", badge: "border-zinc-200 dark:border-zinc-800" },
  ocean: { accent: "bg-blue-600", badge: "border-blue-200 dark:border-blue-800" },
  sunset: { accent: "bg-orange-500", badge: "border-orange-200 dark:border-orange-800" },
  forest: { accent: "bg-emerald-600", badge: "border-emerald-200 dark:border-emerald-800" },
  midnight: { accent: "bg-indigo-700", badge: "border-indigo-200 dark:border-indigo-800" },
  rose: { accent: "bg-pink-500", badge: "border-pink-200 dark:border-pink-800" },
};

function ProfileContent() {
  const searchParams = useSearchParams();
  const username = searchParams.get("u");
  const [user, setUser] = useState<PublicUser | null>(null);
  const [connections, setConnections] = useState<Connection[]>([]);
  const [notFound, setNotFound] = useState(false);
  const [loading, setLoading] = useState(true);
  const [showQR, setShowQR] = useState(false);

  const profileUrl = username
    ? `${PROFILE_BASE_URL}/profile?u=${username}`
    : "";

  useEffect(() => {
    if (!username) {
      setNotFound(true);
      setLoading(false);
      return;
    }
    api.users.publicProfile(username).then((result) => {
      if (result.error) {
        setNotFound(true);
      } else if (result.data) {
        setUser(result.data.user);
      }
      setLoading(false);
    });
    api.users.publicConnections(username).then((result) => {
      if (result.data) {
        setConnections(result.data.connections || []);
      }
    });
    // Track profile view
    api.analytics.trackView(username);
  }, [username]);

  if (loading) {
    return (
      <div className="flex flex-1 items-center justify-center py-24">
        <p className="text-zinc-400">Loading profile...</p>
      </div>
    );
  }

  if (notFound || !user) {
    return (
      <div className="flex flex-1 items-center justify-center py-24">
        <div className="text-center">
          <h1 className="text-2xl font-bold">User not found</h1>
          <p className="mt-2 text-zinc-500">
            This creator profile doesn&apos;t exist.
          </p>
        </div>
      </div>
    );
  }

  const theme = THEME_STYLES[user.theme] || THEME_STYLES.default;

  return (
    <div className="mx-auto max-w-2xl px-4 py-16 sm:px-6">
      {/* Theme accent bar */}
      {user.theme && user.theme !== "default" && (
        <div className={`-mt-16 mb-8 h-2 rounded-full ${theme.accent}`} />
      )}
      <div className="flex flex-col items-center text-center">
        {/* Avatar */}
        <div className="relative">
          {user.image ? (
            <img
              src={user.image}
              alt={user.name || "Creator"}
              className="h-24 w-24 rounded-full object-cover"
            />
          ) : (
            <div className="flex h-24 w-24 items-center justify-center rounded-full bg-zinc-100 text-3xl font-bold text-zinc-400 dark:bg-zinc-800">
              {(user.name || "?")[0].toUpperCase()}
            </div>
          )}
          {user.isVerified && (
            <div className="absolute -bottom-1 -right-1 rounded-full bg-white p-0.5 dark:bg-zinc-950">
              <CheckCircle className="h-6 w-6 text-blue-500" />
            </div>
          )}
        </div>

        {/* Name & Username */}
        <h1 className="mt-6 text-2xl font-bold">{user.name}</h1>
        <p className="text-zinc-500">@{user.username}</p>

        {/* Bio */}
        {user.bio && (
          <p className="mt-4 max-w-md text-zinc-600 dark:text-zinc-400">
            {user.bio}
          </p>
        )}

        {/* Creator Score */}
        {user.creatorScore !== null && (
          <div className="mt-6 flex items-center gap-2 rounded-full border border-zinc-200 px-4 py-2 dark:border-zinc-800">
            <Shield className="h-4 w-4" />
            <span className="text-sm font-medium">
              Creator Score: {user.creatorScore}
            </span>
          </div>
        )}

        {/* Connected Platforms */}
        {connections.length > 0 && (
          <div className="mt-8 w-full max-w-sm">
            <h2 className="mb-3 text-sm font-semibold text-zinc-500">
              Connected Platforms
            </h2>
            <div className="space-y-2">
              {connections.map((conn) => (
                <a
                  key={conn.platform}
                  href={conn.profileUrl || "#"}
                  target="_blank"
                  rel="noopener noreferrer"
                  onClick={() => username && api.analytics.trackClick(username, "connection", conn.platform)}
                  className="flex items-center gap-3 rounded-lg border border-zinc-200 p-3 transition-colors hover:bg-zinc-50 dark:border-zinc-800 dark:hover:bg-zinc-900"
                >
                  {conn.avatarUrl ? (
                    <img
                      src={conn.avatarUrl}
                      alt=""
                      className="h-8 w-8 rounded-full object-cover"
                    />
                  ) : (
                    <div className="flex h-8 w-8 items-center justify-center rounded-full bg-zinc-100 text-xs font-bold dark:bg-zinc-800">
                      {conn.platform[0].toUpperCase()}
                    </div>
                  )}
                  <div className="flex-1 text-left">
                    <p className="text-sm font-medium capitalize">
                      {conn.platform}
                    </p>
                    <p className="text-xs text-zinc-500">
                      {conn.username || conn.displayName}
                    </p>
                  </div>
                  {conn.followerCount !== null && (
                    <span className="text-xs text-zinc-400">
                      {conn.followerCount.toLocaleString()}
                    </span>
                  )}
                  <ExternalLink className="h-3.5 w-3.5 text-zinc-300" />
                </a>
              ))}
            </div>
          </div>
        )}

        {/* Custom Links */}
        {user.customLinks && user.customLinks.length > 0 && (
          <div className="mt-8 w-full max-w-sm">
            <h2 className="mb-3 text-sm font-semibold text-zinc-500">
              Links
            </h2>
            <div className="space-y-2">
              {user.customLinks.map((link, i) => (
                <a
                  key={i}
                  href={link.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  onClick={() => username && api.analytics.trackClick(username, "custom", link.url)}
                  className={`flex items-center justify-between rounded-lg border p-3 transition-colors hover:bg-zinc-50 dark:hover:bg-zinc-900 ${theme.badge}`}
                >
                  <span className="text-sm font-medium">{link.title}</span>
                  <ExternalLink className="h-3.5 w-3.5 text-zinc-300" />
                </a>
              ))}
            </div>
          </div>
        )}

        {/* Share Section */}
        <div className="mt-8 w-full max-w-sm">
          <div className="flex flex-wrap items-center justify-center gap-2">
            <CopyLinkButton url={profileUrl} />
            <ShareTwitterButton
              url={profileUrl}
              text={`Check out ${user.name}'s Creator Passport on Creatrid`}
            />
            <ShareLinkedInButton url={profileUrl} />
            <button
              onClick={() => setShowQR(!showQR)}
              className="flex items-center gap-1.5 rounded-lg border border-zinc-200 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
            >
              <QrCode className="h-3.5 w-3.5" />
              QR Code
            </button>
          </div>

          {/* QR Code */}
          {showQR && (
            <div className="mt-4 flex flex-col items-center">
              <div className="rounded-xl border border-zinc-200 bg-white p-4 dark:border-zinc-800 dark:bg-zinc-900">
                <QRCodeSVG
                  value={profileUrl}
                  size={180}
                  level="M"
                  marginSize={2}
                />
              </div>
              <p className="mt-2 text-xs text-zinc-400">
                Scan to view this profile
              </p>
            </div>
          )}
        </div>

        {/* Creatrid Badge */}
        <div className="mt-8 flex items-center gap-2 text-xs text-zinc-400">
          <Shield className="h-3 w-3" />
          Verified on Creatrid
        </div>
      </div>
    </div>
  );
}

export default function PublicProfilePage() {
  return (
    <Suspense
      fallback={
        <div className="flex flex-1 items-center justify-center py-24">
          <p className="text-zinc-400">Loading profile...</p>
        </div>
      }
    >
      <ProfileContent />
    </Suspense>
  );
}
