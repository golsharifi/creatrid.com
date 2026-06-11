"use client";

import { useState, useEffect, useRef } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { useTheme } from "next-themes";
import { features } from "@/lib/features";
import {
  Sun,
  Moon,
  LogOut,
  User,
  Users,
  LayoutDashboard,
  Menu,
  X,
  Link2,
  Compass,
  BarChart3,
  Settings,
  Globe,
  Archive,
  ShoppingBag,
  DollarSign,
  Search,
  FolderOpen,
  Gift,
  Webhook,
  BookOpen,
  Shield,
  Zap,
  Music,
  Image as ImageIcon,
  Trophy,
  FileText,
  CreditCard,
} from "@/components/icons";
import { useTranslation } from "react-i18next";
import { NotificationBell } from "@/components/notification-bell";
import { Logo } from "@/components/brand/logo";
import { AmbientToggle } from "@/components/ambient-toggle";

const LANGUAGES = [
  { code: "en", label: "English" },
  { code: "es", label: "Español" },
  { code: "fa", label: "فارسی" },
];

// Primary "service matrix" tabs — the colored top buttons, in the client's order.
type Tab = {
  href: string;
  key: string;
  Icon: React.ComponentType<{ className?: string }>;
  accent: string; // CSS var from globals.css
  flag?: keyof typeof features;
};

const PRIMARY_TABS: Tab[] = [
  { href: "/music", key: "music", Icon: Music, accent: "var(--tab-music)" },
  { href: "/vault", key: "vault", Icon: Archive, accent: "var(--tab-vault)" },
  { href: "/logo-studio", key: "logoStudio", Icon: ImageIcon, accent: "var(--tab-studio)" },
  { href: "/vr-community", key: "vrCommunity", Icon: Trophy, accent: "var(--tab-vr)" },
  { href: "/legal-ai", key: "legalAi", Icon: FileText, accent: "var(--tab-legal)" },
  { href: "/trade", key: "trade", Icon: CreditCard, accent: "var(--tab-trade)", flag: "trade" },
];

// Secondary destinations, housed under the "more" (hamburger) menu.
const MORE_LINKS = [
  { href: "/dashboard", key: "dashboard", Icon: LayoutDashboard },
  { href: "/discover", key: "discover", Icon: Compass },
  { href: "/marketplace", key: "marketplace", Icon: ShoppingBag },
  { href: "/earnings", key: "earnings", Icon: DollarSign },
  { href: "/tokens", key: "tokens", Icon: Zap },
  { href: "/collections", key: "collections", Icon: FolderOpen },
  { href: "/connections", key: "connections", Icon: Link2 },
  { href: "/analytics", key: "analytics", Icon: BarChart3 },
  { href: "/verify", key: "verify", Icon: Shield },
  { href: "/referrals", key: "referrals", Icon: Gift },
  { href: "/webhooks", key: "webhooks", Icon: Webhook },
  { href: "/search", key: "search", Icon: Search },
  { href: "/blog", key: "blog", Icon: BookOpen },
];

export function Header() {
  const { user, logout } = useAuth();
  const { theme, setTheme } = useTheme();
  const { t, i18n } = useTranslation();
  const pathname = usePathname();
  const [mobileOpen, setMobileOpen] = useState(false);
  const [langOpen, setLangOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const langRef = useRef<HTMLDivElement>(null);

  const visibleTabs = PRIMARY_TABS.filter((tab) => !tab.flag || features[tab.flag]);

  // Close menus on outside click
  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setMobileOpen(false);
      }
      if (langRef.current && !langRef.current.contains(e.target as Node)) {
        setLangOpen(false);
      }
    }
    if (mobileOpen || langOpen) {
      document.addEventListener("mousedown", handleClickOutside);
    }
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [mobileOpen, langOpen]);

  // Close menus on escape
  useEffect(() => {
    function handleEscape(e: KeyboardEvent) {
      if (e.key === "Escape") {
        setMobileOpen(false);
        setLangOpen(false);
      }
    }
    document.addEventListener("keydown", handleEscape);
    return () => document.removeEventListener("keydown", handleEscape);
  }, []);

  // Close the menu whenever the route changes
  useEffect(() => {
    setMobileOpen(false);
  }, [pathname]);

  const currentLang = LANGUAGES.find((l) => l.code === i18n.language) || LANGUAGES[0];
  const isActive = (href: string) => pathname === href || pathname.startsWith(href + "/");

  return (
    <header className="sticky top-0 z-50 border-b border-zinc-200 bg-white/80 backdrop-blur-sm dark:border-zinc-800 dark:bg-zinc-950/80">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between gap-4 px-4 sm:px-6">
        {/* Logo — top left */}
        <Link href="/" className="shrink-0" aria-label="Creatrid home">
          <Logo />
        </Link>

        {/* Primary colored tabs — desktop */}
        {user && (
          <nav className="hidden flex-1 items-center justify-center gap-1 lg:flex">
            {visibleTabs.map((tab) => {
              const active = isActive(tab.href);
              return (
                <Link
                  key={tab.href}
                  href={tab.href}
                  className="group relative flex items-center gap-1.5 rounded-lg px-3 py-2 text-sm font-semibold transition-colors"
                  style={{ color: active ? tab.accent : undefined }}
                >
                  <span
                    className={
                      active
                        ? ""
                        : "text-zinc-500 transition-colors group-hover:text-zinc-900 dark:text-zinc-400 dark:group-hover:text-zinc-100"
                    }
                    style={active ? { color: tab.accent } : undefined}
                  >
                    <tab.Icon className="h-4 w-4" />
                  </span>
                  <span
                    className={active ? "" : "text-zinc-600 group-hover:text-zinc-900 dark:text-zinc-400 dark:group-hover:text-zinc-100"}
                  >
                    {t(`header.${tab.key}`)}
                  </span>
                  {/* active accent underline */}
                  <span
                    className="absolute inset-x-2 -bottom-px h-0.5 rounded-full transition-all"
                    style={{ background: active ? tab.accent : "transparent" }}
                  />
                </Link>
              );
            })}
          </nav>
        )}

        <div className="flex items-center gap-1.5">
          {/* Ambient sound on/off */}
          <AmbientToggle />

          {/* Language switcher */}
          <div className="relative" ref={langRef}>
            <button
              onClick={() => setLangOpen((prev) => !prev)}
              className="flex items-center gap-1 rounded-lg p-2 text-zinc-500 transition-colors hover:bg-zinc-100 hover:text-zinc-900 dark:text-zinc-400 dark:hover:bg-zinc-800 dark:hover:text-zinc-100"
              aria-label={t("header.language")}
            >
              <Globe className="h-4 w-4" />
              <span className="hidden text-xs font-medium uppercase sm:inline">{currentLang.code}</span>
            </button>
            {langOpen && (
              <div className="absolute right-0 top-full mt-1 w-36 rounded-lg border border-zinc-200 bg-white py-1 shadow-lg dark:border-zinc-800 dark:bg-zinc-950">
                {LANGUAGES.map((lang) => (
                  <button
                    key={lang.code}
                    onClick={() => {
                      i18n.changeLanguage(lang.code);
                      setLangOpen(false);
                    }}
                    className={`flex w-full items-center gap-2 px-3 py-2 text-sm transition-colors hover:bg-zinc-100 dark:hover:bg-zinc-800 ${
                      i18n.language === lang.code
                        ? "font-medium text-zinc-900 dark:text-zinc-100"
                        : "text-zinc-600 dark:text-zinc-400"
                    }`}
                  >
                    <span className="text-xs uppercase">{lang.code}</span>
                    {lang.label}
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Theme toggle */}
          <button
            onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
            className="rounded-lg p-2 text-zinc-500 transition-colors hover:bg-zinc-100 hover:text-zinc-900 dark:text-zinc-400 dark:hover:bg-zinc-800 dark:hover:text-zinc-100"
            aria-label={t("header.toggleTheme")}
          >
            <Sun className="hidden h-5 w-5 dark:block" />
            <Moon className="block h-5 w-5 dark:hidden" />
          </button>

          {user ? (
            <>
              <NotificationBell />

              {/* "More" menu — settings, profile, wallet + all secondary tabs */}
              <div className="relative" ref={menuRef}>
                <button
                  onClick={() => setMobileOpen((prev) => !prev)}
                  className="rounded-lg p-2 text-zinc-500 transition-colors hover:bg-zinc-100 hover:text-zinc-900 dark:text-zinc-400 dark:hover:bg-zinc-800 dark:hover:text-zinc-100"
                  aria-label={mobileOpen ? t("header.closeMenu") : t("header.openMenu")}
                >
                  {mobileOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
                </button>

                {mobileOpen && (
                  <div className="absolute right-0 top-full mt-2 max-h-[80vh] w-64 overflow-y-auto rounded-xl border border-zinc-200 bg-white py-2 shadow-xl dark:border-zinc-800 dark:bg-zinc-950">
                    {/* Primary tabs (shown here too, for mobile + quick access) */}
                    <p className="px-4 pb-1 pt-1 text-[10px] font-semibold uppercase tracking-wider text-zinc-400 lg:hidden">
                      {t("header.servicesLabel")}
                    </p>
                    <div className="lg:hidden">
                      {visibleTabs.map((tab) => (
                        <Link
                          key={tab.href}
                          href={tab.href}
                          onClick={() => setMobileOpen(false)}
                          className="flex items-center gap-3 px-4 py-2.5 text-sm font-medium text-zinc-700 transition-colors hover:bg-zinc-100 dark:text-zinc-300 dark:hover:bg-zinc-800"
                        >
                          <span style={{ color: tab.accent }}>
                            <tab.Icon className="h-4 w-4" />
                          </span>
                          {t(`header.${tab.key}`)}
                        </Link>
                      ))}
                      <div className="my-1 border-t border-zinc-200 dark:border-zinc-800" />
                    </div>

                    <p className="px-4 pb-1 pt-1 text-[10px] font-semibold uppercase tracking-wider text-zinc-400">
                      {t("header.moreLabel")}
                    </p>
                    {MORE_LINKS.map((link) => (
                      <Link
                        key={link.href}
                        href={link.href}
                        onClick={() => setMobileOpen(false)}
                        className="flex items-center gap-3 px-4 py-2.5 text-sm text-zinc-700 transition-colors hover:bg-zinc-100 dark:text-zinc-300 dark:hover:bg-zinc-800"
                      >
                        <link.Icon className="h-4 w-4" />
                        {t(`header.${link.key}`)}
                      </Link>
                    ))}
                    {user.role === "BRAND" && (
                      <Link
                        href="/agency"
                        onClick={() => setMobileOpen(false)}
                        className="flex items-center gap-3 px-4 py-2.5 text-sm text-zinc-700 transition-colors hover:bg-zinc-100 dark:text-zinc-300 dark:hover:bg-zinc-800"
                      >
                        <Users className="h-4 w-4" />
                        {t("header.agency")}
                      </Link>
                    )}

                    <div className="my-1 border-t border-zinc-200 dark:border-zinc-800" />
                    <Link
                      href="/settings"
                      onClick={() => setMobileOpen(false)}
                      className="flex items-center gap-3 px-4 py-2.5 text-sm text-zinc-700 transition-colors hover:bg-zinc-100 dark:text-zinc-300 dark:hover:bg-zinc-800"
                    >
                      <Settings className="h-4 w-4" />
                      {t("header.settings")}
                    </Link>
                    {user.username && (
                      <Link
                        href={`/${user.username}`}
                        onClick={() => setMobileOpen(false)}
                        className="flex items-center gap-3 px-4 py-2.5 text-sm text-zinc-700 transition-colors hover:bg-zinc-100 dark:text-zinc-300 dark:hover:bg-zinc-800"
                      >
                        <User className="h-4 w-4" />
                        {t("header.myProfile")}
                      </Link>
                    )}
                    <button
                      onClick={() => {
                        setMobileOpen(false);
                        logout();
                      }}
                      className="flex w-full items-center gap-3 px-4 py-2.5 text-sm text-zinc-700 transition-colors hover:bg-zinc-100 dark:text-zinc-300 dark:hover:bg-zinc-800"
                    >
                      <LogOut className="h-4 w-4" />
                      {t("header.logout")}
                    </button>
                  </div>
                )}
              </div>
            </>
          ) : (
            <Link
              href="/sign-in"
              className="rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              {t("header.signIn")}
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}
