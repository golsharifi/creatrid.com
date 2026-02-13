"use client";

import Link from "next/link";
import { useTranslation } from "react-i18next";

export function Footer() {
  const { t } = useTranslation();

  return (
    <footer className="border-t border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/50">
      <div className="mx-auto max-w-6xl px-4 py-12 sm:px-6">
        <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
          {/* Product */}
          <div>
            <h3 className="text-sm font-semibold">{t("footer.product")}</h3>
            <ul className="mt-3 flex flex-col gap-2 text-sm text-zinc-500 dark:text-zinc-400">
              <li>
                <Link href="/dashboard" className="transition-colors hover:text-zinc-900 dark:hover:text-zinc-100">
                  {t("footer.dashboard")}
                </Link>
              </li>
              <li>
                <Link href="/discover" className="transition-colors hover:text-zinc-900 dark:hover:text-zinc-100">
                  {t("footer.discover")}
                </Link>
              </li>
              <li>
                <Link href="/marketplace" className="transition-colors hover:text-zinc-900 dark:hover:text-zinc-100">
                  {t("footer.marketplace")}
                </Link>
              </li>
              <li>
                <Link href="/pricing" className="transition-colors hover:text-zinc-900 dark:hover:text-zinc-100">
                  {t("footer.pricing")}
                </Link>
              </li>
            </ul>
          </div>

          {/* Resources */}
          <div>
            <h3 className="text-sm font-semibold">{t("footer.resources")}</h3>
            <ul className="mt-3 flex flex-col gap-2 text-sm text-zinc-500 dark:text-zinc-400">
              <li>
                <Link href="/blog" className="transition-colors hover:text-zinc-900 dark:hover:text-zinc-100">
                  {t("footer.blog")}
                </Link>
              </li>
              <li>
                <Link href="/api-keys" className="transition-colors hover:text-zinc-900 dark:hover:text-zinc-100">
                  {t("footer.apiDocs")}
                </Link>
              </li>
              <li>
                <Link href="/help" className="transition-colors hover:text-zinc-900 dark:hover:text-zinc-100">
                  {t("footer.help")}
                </Link>
              </li>
            </ul>
          </div>

          {/* Legal */}
          <div>
            <h3 className="text-sm font-semibold">{t("footer.legal")}</h3>
            <ul className="mt-3 flex flex-col gap-2 text-sm text-zinc-500 dark:text-zinc-400">
              <li>
                <Link href="/terms" className="transition-colors hover:text-zinc-900 dark:hover:text-zinc-100">
                  {t("footer.terms")}
                </Link>
              </li>
              <li>
                <Link href="/privacy" className="transition-colors hover:text-zinc-900 dark:hover:text-zinc-100">
                  {t("footer.privacy")}
                </Link>
              </li>
            </ul>
          </div>

          {/* Company */}
          <div>
            <h3 className="text-sm font-semibold">Creatrid</h3>
            <p className="mt-3 text-sm text-zinc-500 dark:text-zinc-400">
              {t("footer.tagline")}
            </p>
          </div>
        </div>

        <div className="mt-10 border-t border-zinc-200 pt-6 text-center text-sm text-zinc-500 dark:border-zinc-800 dark:text-zinc-400">
          {t("footer.copyright", { year: new Date().getFullYear() })}
        </div>
      </div>
    </footer>
  );
}
