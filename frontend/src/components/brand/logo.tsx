"use client";

/**
 * Creatrid brand mark — a "globe orb" passport seal.
 *
 * Inspired by the Earmark blueprint's globe motif but rebuilt around our own
 * cyan→azure→violet gradient with a gold orbital ring (the "passport seal").
 * Used top-left in the header alongside the wordmark.
 */

export function OrbMark({ className = "h-8 w-8" }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 48 48"
      className={className}
      role="img"
      aria-label="Creatrid globe orb"
    >
      <defs>
        <radialGradient id="orbCore" cx="38%" cy="32%" r="75%">
          <stop offset="0%" stopColor="#67e8f9" />
          <stop offset="45%" stopColor="var(--brand-azure)" />
          <stop offset="100%" stopColor="var(--brand-violet)" />
        </radialGradient>
        <linearGradient id="orbRing" x1="0" y1="0" x2="1" y2="1">
          <stop offset="0%" stopColor="var(--brand-gold)" />
          <stop offset="100%" stopColor="#e0a93b" />
        </linearGradient>
      </defs>

      {/* gold orbital ring (passport seal) */}
      <ellipse
        cx="24"
        cy="24"
        rx="22"
        ry="9"
        fill="none"
        stroke="url(#orbRing)"
        strokeWidth="2"
        transform="rotate(-24 24 24)"
      />

      {/* globe orb */}
      <circle cx="24" cy="24" r="14" fill="url(#orbCore)" />

      {/* meridians / latitudes — the "global ID" graticule */}
      <g
        fill="none"
        stroke="#ffffff"
        strokeOpacity="0.55"
        strokeWidth="1"
      >
        <ellipse cx="24" cy="24" rx="6" ry="14" />
        <ellipse cx="24" cy="24" rx="13" ry="5.5" />
        <line x1="10" y1="24" x2="38" y2="24" strokeOpacity="0.4" />
      </g>
    </svg>
  );
}

export function Logo({ className = "" }: { className?: string }) {
  return (
    <span className={`inline-flex items-center gap-2 ${className}`}>
      <OrbMark className="h-8 w-8 shrink-0" />
      <span className="flex flex-col leading-none">
        <span className="text-lg font-extrabold tracking-tight">creatrid</span>
        <span className="text-[9px] font-semibold uppercase tracking-[0.18em] text-zinc-500 dark:text-zinc-400">
          Passport ID
        </span>
      </span>
    </span>
  );
}
