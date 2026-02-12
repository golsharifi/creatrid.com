"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { OnboardingForm } from "./onboarding-form";
import { useTranslation } from "react-i18next";

export default function OnboardingPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
    if (!loading && user?.onboarded) router.push("/dashboard");
  }, [user, loading, router]);

  if (loading || !user) return null;

  return (
    <div className="flex flex-1 items-center justify-center px-4 py-24">
      <div className="w-full max-w-md space-y-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold">{t("onboarding.title")}</h1>
          <p className="mt-2 text-sm text-zinc-500">
            {t("onboarding.subtitle")}
          </p>
        </div>
        <OnboardingForm name={user.name || ""} />
      </div>
    </div>
  );
}
