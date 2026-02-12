"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { SettingsForm } from "./settings-form";
import { useTranslation } from "react-i18next";

export default function SettingsPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
  }, [user, loading, router]);

  if (loading || !user) return null;

  return (
    <div className="mx-auto max-w-2xl px-4 py-12 sm:px-6">
      <h1 className="text-2xl font-bold">{t("settings.title")}</h1>
      <p className="mt-1 text-sm text-zinc-500">
        {t("settings.subtitle")}
      </p>

      <div className="mt-8">
        <SettingsForm
          name={user.name || ""}
          username={user.username || ""}
          bio={user.bio || ""}
          image={user.image || ""}
          theme={user.theme || "default"}
          customLinks={user.customLinks || []}
          emailPrefs={user.emailPrefs || { welcome: true, connectionAlert: true, weeklyDigest: true, collaborations: true }}
        />
      </div>
    </div>
  );
}
