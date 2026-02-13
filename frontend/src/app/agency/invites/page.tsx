"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { useTranslation } from "react-i18next";

type Invite = {
  id: string;
  agencyId: string;
  creatorId: string;
  status: string;
  invitedAt: string;
  joinedAt: string | null;
  agencyName: string | null;
  agencyDescription: string | null;
};

export default function AgencyInvitesPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();
  const [invites, setInvites] = useState<Invite[]>([]);
  const [responding, setResponding] = useState<string | null>(null);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
  }, [user, loading, router]);

  useEffect(() => {
    if (user) {
      api.agency.invites().then((r) => {
        if (r.data) setInvites(r.data.invites);
      });
    }
  }, [user]);

  async function handleRespond(id: string, action: "accept" | "decline") {
    setResponding(id);
    const result = await api.agency.respondToInvite(id, action);
    if (result.data) {
      setInvites((prev) => prev.filter((inv) => inv.id !== id));
    }
    setResponding(null);
  }

  if (loading || !user) return null;

  return (
    <div className="mx-auto max-w-2xl px-4 py-12 sm:px-6">
      <div className="mb-8">
        <h1 className="text-2xl font-bold">{t("agency.invites")}</h1>
        <p className="mt-1 text-zinc-500">{t("agency.invitesSubtitle")}</p>
      </div>

      {invites.length === 0 ? (
        <div className="rounded-xl border border-zinc-200 px-6 py-12 text-center dark:border-zinc-800">
          <p className="text-zinc-500">{t("agency.noInvites")}</p>
        </div>
      ) : (
        <div className="space-y-4">
          {invites.map((invite) => (
            <div
              key={invite.id}
              className="rounded-xl border border-zinc-200 p-5 dark:border-zinc-800"
            >
              <div className="flex items-start justify-between">
                <div>
                  <h3 className="font-semibold">{invite.agencyName || "Unknown Agency"}</h3>
                  {invite.agencyDescription && (
                    <p className="mt-1 text-sm text-zinc-500">{invite.agencyDescription}</p>
                  )}
                  <p className="mt-2 text-xs text-zinc-400">
                    {new Date(invite.invitedAt).toLocaleDateString()}
                  </p>
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={() => handleRespond(invite.id, "accept")}
                    disabled={responding === invite.id}
                    className="rounded-lg bg-zinc-900 px-3 py-1.5 text-xs font-medium text-white transition-colors hover:bg-zinc-800 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
                  >
                    {responding === invite.id ? "..." : t("agency.accept")}
                  </button>
                  <button
                    onClick={() => handleRespond(invite.id, "decline")}
                    disabled={responding === invite.id}
                    className="rounded-lg border border-zinc-200 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-zinc-50 disabled:opacity-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
                  >
                    {t("agency.decline")}
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
