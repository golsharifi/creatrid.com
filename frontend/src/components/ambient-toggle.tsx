"use client";

import { useEffect, useState } from "react";
import { getAmbientEngine } from "@/lib/ambient-audio";
import { Music } from "@/components/icons";

/**
 * Header on/off toggle for the ambient soundscape (rain & storms), as
 * requested in the blueprint. Lights up in the Music tab accent when playing.
 */
export function AmbientToggle() {
  const [playing, setPlaying] = useState(false);

  useEffect(() => {
    return getAmbientEngine().subscribe((s) => setPlaying(s.playing));
  }, []);

  return (
    <button
      onClick={() => getAmbientEngine().toggle()}
      className={`relative rounded-lg p-2 transition-colors ${
        playing
          ? ""
          : "text-zinc-500 hover:bg-zinc-100 hover:text-zinc-900 dark:text-zinc-400 dark:hover:bg-zinc-800 dark:hover:text-zinc-100"
      }`}
      style={playing ? { color: "var(--tab-music)" } : undefined}
      aria-label={playing ? "Turn ambient sound off" : "Turn ambient sound on"}
      aria-pressed={playing}
      title={playing ? "Ambient sound: on" : "Ambient sound: off"}
    >
      <Music className="h-5 w-5" />
      {playing && (
        <span
          className="absolute right-1 top-1 h-1.5 w-1.5 animate-pulse rounded-full"
          style={{ background: "var(--tab-music)" }}
        />
      )}
    </button>
  );
}
