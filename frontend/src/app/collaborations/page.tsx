"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import Link from "next/link";
import { Inbox, Send, CheckCircle, XCircle, Clock } from "@/components/icons";
import { useTranslation } from "react-i18next";

type IncomingRequest = {
  id: string;
  fromUserId: string;
  message: string;
  status: string;
  createdAt: string;
  fromName: string | null;
  fromUsername: string | null;
  fromImage: string | null;
};

type OutgoingRequest = {
  id: string;
  toUserId: string;
  message: string;
  status: string;
  createdAt: string;
  toName: string | null;
  toUsername: string | null;
  toImage: string | null;
};

export default function CollaborationsPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();
  const [tab, setTab] = useState<"inbox" | "outbox">("inbox");
  const [incoming, setIncoming] = useState<IncomingRequest[]>([]);
  const [outgoing, setOutgoing] = useState<OutgoingRequest[]>([]);
  const [responding, setResponding] = useState<string | null>(null);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
  }, [user, loading, router]);

  useEffect(() => {
    if (user) {
      api.collaborations.inbox().then((r) => {
        if (r.data) setIncoming(r.data.requests);
      });
      api.collaborations.outbox().then((r) => {
        if (r.data) setOutgoing(r.data.requests);
      });
    }
  }, [user]);

  async function respond(id: string, action: "accept" | "decline") {
    setResponding(id);
    const result = await api.collaborations.respond(id, action);
    if (result.data) {
      setIncoming((prev) =>
        prev.map((r) => (r.id === id ? { ...r, status: result.data!.status } : r))
      );
    }
    setResponding(null);
  }

  if (loading || !user) return null;

  const pendingCount = incoming.filter((r) => r.status === "pending").length;

  return (
    <div className="mx-auto max-w-3xl px-4 py-12 sm:px-6">
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("collaborations.title")}</h1>
          <p className="mt-1 text-zinc-500">
            {t("collaborations.subtitle")}
          </p>
        </div>
        <Link
          href="/discover"
          className="rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-700 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-300"
        >
          {t("collaborations.discoverCreators")}
        </Link>
      </div>

      {/* Tabs */}
      <div className="mb-6 flex gap-1 rounded-lg border border-zinc-200 p-1 dark:border-zinc-800">
        <button
          onClick={() => setTab("inbox")}
          className={`flex flex-1 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors ${
            tab === "inbox"
              ? "bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900"
              : "text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-100"
          }`}
        >
          <Inbox className="h-4 w-4" />
          {t("collaborations.inbox")}
          {pendingCount > 0 && (
            <span className="rounded-full bg-blue-500 px-1.5 py-0.5 text-[10px] font-bold text-white">
              {pendingCount}
            </span>
          )}
        </button>
        <button
          onClick={() => setTab("outbox")}
          className={`flex flex-1 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors ${
            tab === "outbox"
              ? "bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900"
              : "text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-100"
          }`}
        >
          <Send className="h-4 w-4" />
          {t("collaborations.sent")}
        </button>
      </div>

      {/* Inbox */}
      {tab === "inbox" && (
        <div className="space-y-3">
          {incoming.length === 0 ? (
            <p className="rounded-xl border border-dashed border-zinc-200 p-8 text-center text-sm text-zinc-400 dark:border-zinc-800">
              {t("collaborations.noRequestsYet")}
            </p>
          ) : (
            incoming.map((req) => (
              <div
                key={req.id}
                className="rounded-xl border border-zinc-200 p-4 dark:border-zinc-800"
              >
                <div className="flex items-start gap-3">
                  {req.fromImage ? (
                    <img src={req.fromImage} alt="" className="h-10 w-10 rounded-full object-cover" />
                  ) : (
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-zinc-100 text-sm font-bold dark:bg-zinc-800">
                      {(req.fromName || "?")[0].toUpperCase()}
                    </div>
                  )}
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <Link
                        href={`/profile?u=${req.fromUsername}`}
                        className="font-medium hover:underline"
                      >
                        {req.fromName || t("common.creator")}
                      </Link>
                      <StatusBadge status={req.status} />
                    </div>
                    {req.fromUsername && (
                      <p className="text-sm text-zinc-500">@{req.fromUsername}</p>
                    )}
                    {req.message && (
                      <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                        {req.message}
                      </p>
                    )}
                    <p className="mt-1 text-xs text-zinc-400">
                      {new Date(req.createdAt).toLocaleDateString()}
                    </p>
                  </div>
                  {req.status === "pending" && (
                    <div className="flex gap-2">
                      <button
                        onClick={() => respond(req.id, "accept")}
                        disabled={responding === req.id}
                        className="rounded-lg bg-zinc-900 px-3 py-1.5 text-xs font-medium text-white hover:bg-zinc-700 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900"
                      >
                        {t("collaborations.accept")}
                      </button>
                      <button
                        onClick={() => respond(req.id, "decline")}
                        disabled={responding === req.id}
                        className="rounded-lg border border-zinc-200 px-3 py-1.5 text-xs font-medium hover:bg-zinc-50 disabled:opacity-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
                      >
                        {t("collaborations.decline")}
                      </button>
                    </div>
                  )}
                </div>
              </div>
            ))
          )}
        </div>
      )}

      {/* Outbox */}
      {tab === "outbox" && (
        <div className="space-y-3">
          {outgoing.length === 0 ? (
            <p className="rounded-xl border border-dashed border-zinc-200 p-8 text-center text-sm text-zinc-400 dark:border-zinc-800">
              {t("collaborations.noSentRequests")} <Link href="/discover" className="underline">{t("collaborations.discoverToCollab")}</Link>{t("collaborations.discoverToCollabSuffix")}
            </p>
          ) : (
            outgoing.map((req) => (
              <div
                key={req.id}
                className="rounded-xl border border-zinc-200 p-4 dark:border-zinc-800"
              >
                <div className="flex items-start gap-3">
                  {req.toImage ? (
                    <img src={req.toImage} alt="" className="h-10 w-10 rounded-full object-cover" />
                  ) : (
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-zinc-100 text-sm font-bold dark:bg-zinc-800">
                      {(req.toName || "?")[0].toUpperCase()}
                    </div>
                  )}
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <Link
                        href={`/profile?u=${req.toUsername}`}
                        className="font-medium hover:underline"
                      >
                        {req.toName || t("common.creator")}
                      </Link>
                      <StatusBadge status={req.status} />
                    </div>
                    {req.toUsername && (
                      <p className="text-sm text-zinc-500">@{req.toUsername}</p>
                    )}
                    {req.message && (
                      <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">
                        {req.message}
                      </p>
                    )}
                    <p className="mt-1 text-xs text-zinc-400">
                      {new Date(req.createdAt).toLocaleDateString()}
                    </p>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      )}
    </div>
  );
}

function StatusBadge({ status }: { status: string }) {
  const { t } = useTranslation();

  if (status === "pending") {
    return (
      <span className="inline-flex items-center gap-1 rounded-full bg-amber-50 px-2 py-0.5 text-[10px] font-medium text-amber-600 dark:bg-amber-900/20 dark:text-amber-400">
        <Clock className="h-3 w-3" /> {t("collaborations.pending")}
      </span>
    );
  }
  if (status === "accepted") {
    return (
      <span className="inline-flex items-center gap-1 rounded-full bg-green-50 px-2 py-0.5 text-[10px] font-medium text-green-600 dark:bg-green-900/20 dark:text-green-400">
        <CheckCircle className="h-3 w-3" /> {t("collaborations.accepted")}
      </span>
    );
  }
  if (status === "declined") {
    return (
      <span className="inline-flex items-center gap-1 rounded-full bg-red-50 px-2 py-0.5 text-[10px] font-medium text-red-500 dark:bg-red-900/20 dark:text-red-400">
        <XCircle className="h-3 w-3" /> {t("collaborations.declined")}
      </span>
    );
  }
  return null;
}
