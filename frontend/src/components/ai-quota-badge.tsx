"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";

/**
 * Monthly AI generation quota meter. Mount on AI pages; call the returned
 * refresh handler (via `refreshKey` bumping) after each generation.
 */
export function AiQuotaBadge({ refreshKey = 0 }: { refreshKey?: number }) {
  const [quota, setQuota] = useState<{ used: number; limit: number; plan: string } | null>(null);

  const load = useCallback(async () => {
    const res = await api.ai.quota();
    if (res.data) setQuota(res.data);
  }, []);

  useEffect(() => {
    load();
  }, [load, refreshKey]);

  if (!quota) return null;

  const remaining = Math.max(0, quota.limit - quota.used);
  const low = remaining <= 3;

  return (
    <div className="flex items-center gap-2 text-xs text-zinc-500 dark:text-zinc-400">
      <span
        className={`inline-flex items-center gap-1.5 rounded-full px-3 py-1 font-medium ${
          low
            ? "bg-amber-50 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400"
            : "bg-zinc-100 dark:bg-zinc-800"
        }`}
      >
        {remaining} of {quota.limit} generations left this month
      </span>
      {quota.plan === "free" && low && (
        <Link href="/pricing" className="font-medium underline underline-offset-2">
          Upgrade for more
        </Link>
      )}
    </div>
  );
}
