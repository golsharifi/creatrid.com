"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { CheckCircle2, Circle, X } from "lucide-react";
import { useTranslation } from "react-i18next";

interface ChecklistProps {
  user: {
    image: string | null;
    bio: string | null;
    emailVerified: string | null;
  };
  connectionCount: number;
  contentCount: number;
}

export function OnboardingChecklist({ user, connectionCount, contentCount }: ChecklistProps) {
  const { t } = useTranslation();
  const [dismissed, setDismissed] = useState(false);

  useEffect(() => {
    setDismissed(localStorage.getItem("onboarding-dismissed") === "true");
  }, []);

  const steps = [
    { done: !!user.image, label: t("onboarding.setPhoto"), hint: "+5 score", href: "/settings" },
    { done: !!user.bio, label: t("onboarding.writeBio"), hint: "+5 score", href: "/settings" },
    { done: connectionCount > 0, label: t("onboarding.connectAccount"), hint: "+10 score", href: "/connections" },
    { done: !!user.emailVerified, label: t("onboarding.verifyEmail"), hint: "+10 score", href: "/settings" },
    { done: contentCount > 0, label: t("onboarding.uploadContent"), hint: "+5 score", href: "/vault" },
  ];

  const completedCount = steps.filter((s) => s.done).length;
  const allDone = completedCount === steps.length;

  if (dismissed || allDone) return null;

  return (
    <div className="rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
      <div className="mb-4 flex items-center justify-between">
        <div>
          <h3 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            {t("onboarding.title")}
          </h3>
          <p className="text-sm text-zinc-500 dark:text-zinc-400">
            {completedCount}/{steps.length} {t("onboarding.completed")}
          </p>
        </div>
        <button
          onClick={() => {
            setDismissed(true);
            localStorage.setItem("onboarding-dismissed", "true");
          }}
          className="rounded-lg p-1 text-zinc-400 hover:bg-zinc-100 hover:text-zinc-600 dark:hover:bg-zinc-800 dark:hover:text-zinc-300"
        >
          <X className="h-4 w-4" />
        </button>
      </div>
      <div className="mb-4 h-2 overflow-hidden rounded-full bg-zinc-100 dark:bg-zinc-800">
        <div
          className="h-full rounded-full bg-emerald-500 transition-all"
          style={{ width: `${(completedCount / steps.length) * 100}%` }}
        />
      </div>
      <div className="space-y-3">
        {steps.map((step, i) => (
          <Link
            key={i}
            href={step.href}
            className="flex items-center gap-3 rounded-lg p-2 transition-colors hover:bg-zinc-50 dark:hover:bg-zinc-800/50"
          >
            {step.done ? (
              <CheckCircle2 className="h-5 w-5 text-emerald-500" />
            ) : (
              <Circle className="h-5 w-5 text-zinc-300 dark:text-zinc-600" />
            )}
            <span className={`flex-1 text-sm ${step.done ? "text-zinc-400 line-through dark:text-zinc-500" : "text-zinc-700 dark:text-zinc-300"}`}>
              {step.label}
            </span>
            {!step.done && (
              <span className="text-xs text-emerald-600 dark:text-emerald-400">{step.hint}</span>
            )}
          </Link>
        ))}
      </div>
    </div>
  );
}
