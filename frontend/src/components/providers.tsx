"use client";

import { useEffect } from "react";
import { ThemeProvider } from "next-themes";
import { AuthProvider } from "@/lib/auth-context";
import { initSentry } from "@/lib/sentry";
import { useTranslation } from "react-i18next";
import "@/i18n/config";

const RTL_LANGUAGES = ["fa"];

export function Providers({ children }: { children: React.ReactNode }) {
  const { i18n } = useTranslation();

  useEffect(() => {
    initSentry();
  }, []);

  // Update dir and lang attributes on the html element when language changes
  useEffect(() => {
    const dir = RTL_LANGUAGES.includes(i18n.language) ? "rtl" : "ltr";
    document.documentElement.dir = dir;
    document.documentElement.lang = i18n.language;
  }, [i18n.language]);

  return (
    <AuthProvider>
      <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
        {children}
      </ThemeProvider>
    </AuthProvider>
  );
}
