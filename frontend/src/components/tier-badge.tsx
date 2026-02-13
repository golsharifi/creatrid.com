"use client";

const TIER_CONFIG: Record<string, { label: string; color: string; bg: string }> = {
  newcomer: { label: "Newcomer", color: "text-zinc-600 dark:text-zinc-400", bg: "bg-zinc-100 dark:bg-zinc-800" },
  rising: { label: "Rising", color: "text-blue-600 dark:text-blue-400", bg: "bg-blue-50 dark:bg-blue-900/30" },
  established: { label: "Established", color: "text-purple-600 dark:text-purple-400", bg: "bg-purple-50 dark:bg-purple-900/30" },
  elite: { label: "Elite", color: "text-amber-600 dark:text-amber-400", bg: "bg-amber-50 dark:bg-amber-900/30" },
};

export function TierBadge({ tier, size = "sm" }: { tier: string; size?: "sm" | "md" }) {
  const config = TIER_CONFIG[tier] || TIER_CONFIG.newcomer;
  const sizeClasses = size === "md" ? "px-3 py-1 text-sm" : "px-2 py-0.5 text-xs";

  return (
    <span className={`inline-flex items-center gap-1 rounded-full font-medium ${config.color} ${config.bg} ${sizeClasses}`}>
      {tier === "elite" && <span>★</span>}
      {config.label}
    </span>
  );
}
