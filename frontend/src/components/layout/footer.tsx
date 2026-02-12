"use client";

import Link from "next/link";
import { useTranslation } from "react-i18next";

export function Footer() {
  const { t } = useTranslation();

  return (
    <footer className="border-t border-zinc-200 dark:border-zinc-800">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4 text-sm text-zinc-500 sm:px-6">
        <p>{t("footer.copyright", { year: new Date().getFullYear() })}</p>
        <div className="flex items-center gap-4">
          <Link href="/terms" className="transition-colors hover:text-zinc-900 dark:hover:text-zinc-100">
            {t("footer.terms")}
          </Link>
          <Link href="/privacy" className="transition-colors hover:text-zinc-900 dark:hover:text-zinc-100">
            {t("footer.privacy")}
          </Link>
        </div>
      </div>
    </footer>
  );
}
