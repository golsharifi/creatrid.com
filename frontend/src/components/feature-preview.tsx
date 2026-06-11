"use client";

import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { ArrowRight, CheckCircle } from "@/components/icons";

type IconComponent = React.ComponentType<{ className?: string }>;

export type PreviewSection = {
  heading: string;
  items: string[];
};

/**
 * Branded "preview" scaffold for the new service-matrix tabs.
 *
 * These tabs (Music, Logo Studio, VR Community, Legal AI, Trade) are routed and
 * branded now so the client can see the full navigation, with each feature's
 * scope described. The interactive build for each lands in its own phase.
 */
export function FeaturePreview({
  accent,
  Icon,
  title,
  tagline,
  intro,
  sections,
  statusLabel = "In development",
  note,
}: {
  accent: string;
  Icon: IconComponent;
  title: string;
  tagline: string;
  intro: string;
  sections: PreviewSection[];
  statusLabel?: string;
  note?: string;
}) {
  const { user } = useAuth();

  return (
    <div className="mx-auto w-full max-w-5xl px-4 py-16 sm:px-6 sm:py-24">
      {/* Hero */}
      <div className="flex flex-col items-center text-center">
        <div
          className="flex h-16 w-16 items-center justify-center rounded-2xl text-white shadow-lg"
          style={{ background: accent }}
        >
          <Icon className="h-8 w-8" />
        </div>
        <span
          className="mt-5 inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-semibold"
          style={{ background: `color-mix(in srgb, ${accent} 14%, transparent)`, color: accent }}
        >
          <span className="h-1.5 w-1.5 rounded-full" style={{ background: accent }} />
          {statusLabel}
        </span>
        <h1 className="mt-4 text-4xl font-bold tracking-tight sm:text-5xl">{title}</h1>
        <p className="mt-3 text-lg font-medium" style={{ color: accent }}>
          {tagline}
        </p>
        <p className="mx-auto mt-4 max-w-2xl text-zinc-600 dark:text-zinc-400">{intro}</p>

        {!user && (
          <Link
            href="/sign-in"
            className="mt-7 inline-flex items-center gap-2 rounded-lg px-6 py-3 text-sm font-medium text-white transition-opacity hover:opacity-90"
            style={{ background: accent }}
          >
            Get your Passport to unlock
            <ArrowRight className="h-4 w-4" />
          </Link>
        )}
      </div>

      {/* Feature sections */}
      <div className="mt-16 grid gap-6 sm:grid-cols-2">
        {sections.map((section) => (
          <div
            key={section.heading}
            className="rounded-2xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950"
          >
            <h2 className="text-lg font-semibold">{section.heading}</h2>
            <ul className="mt-4 flex flex-col gap-3">
              {section.items.map((item) => (
                <li key={item} className="flex items-start gap-2.5 text-sm text-zinc-600 dark:text-zinc-400">
                  <span className="mt-0.5 shrink-0" style={{ color: accent }}>
                    <CheckCircle className="h-4 w-4" />
                  </span>
                  <span>{item}</span>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>

      {note && (
        <p className="mx-auto mt-10 max-w-2xl rounded-xl border border-dashed border-zinc-300 bg-zinc-50 p-4 text-center text-sm text-zinc-500 dark:border-zinc-700 dark:bg-zinc-900/50 dark:text-zinc-400">
          {note}
        </p>
      )}
    </div>
  );
}
