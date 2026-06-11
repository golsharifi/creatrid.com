"use client";

/**
 * Ambient soundscape engine — moody rain & storms, synthesized with the Web
 * Audio API (no audio assets, nothing to license). A singleton so the header
 * toggle and the Music page control the same soundscape.
 *
 * Rain: looped filtered noise. Storm: rain + randomized low rumbling thunder
 * bursts.
 */

export type AmbientMode = "rain" | "storm";

const STORAGE_KEY = "creatrid-ambient";

type AmbientState = {
  playing: boolean;
  mode: AmbientMode;
  volume: number; // 0..1
};

type Listener = (state: AmbientState) => void;

class AmbientEngine {
  private ctx: AudioContext | null = null;
  private master: GainNode | null = null;
  private rainSource: AudioBufferSourceNode | null = null;
  private thunderTimer: ReturnType<typeof setTimeout> | null = null;
  private listeners = new Set<Listener>();

  state: AmbientState = { playing: false, mode: "rain", volume: 0.5 };

  constructor() {
    if (typeof window !== "undefined") {
      try {
        const saved = localStorage.getItem(STORAGE_KEY);
        if (saved) {
          const parsed = JSON.parse(saved) as Partial<AmbientState>;
          // Never autoplay on load — browsers block it; user re-enables with one tap.
          this.state = {
            playing: false,
            mode: parsed.mode === "storm" ? "storm" : "rain",
            volume: typeof parsed.volume === "number" ? parsed.volume : 0.5,
          };
        }
      } catch {
        // ignore corrupted storage
      }
    }
  }

  subscribe(fn: Listener): () => void {
    this.listeners.add(fn);
    fn(this.state);
    return () => this.listeners.delete(fn);
  }

  private emit() {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(this.state));
    } catch {
      // storage full/blocked — non-fatal
    }
    this.listeners.forEach((fn) => fn(this.state));
  }

  toggle() {
    if (this.state.playing) this.stop();
    else this.start();
  }

  setMode(mode: AmbientMode) {
    this.state = { ...this.state, mode };
    if (this.state.playing) {
      this.stopNodes();
      this.startNodes();
    }
    this.emit();
  }

  setVolume(volume: number) {
    this.state = { ...this.state, volume };
    if (this.master && this.ctx) {
      this.master.gain.setTargetAtTime(volume * 0.6, this.ctx.currentTime, 0.05);
    }
    this.emit();
  }

  start() {
    if (this.state.playing) return;
    this.startNodes();
    this.state = { ...this.state, playing: true };
    this.emit();
  }

  stop() {
    if (!this.state.playing) return;
    this.stopNodes();
    this.state = { ...this.state, playing: false };
    this.emit();
  }

  private ensureContext(): AudioContext {
    if (!this.ctx) {
      this.ctx = new AudioContext();
      this.master = this.ctx.createGain();
      this.master.connect(this.ctx.destination);
    }
    if (this.ctx.state === "suspended") void this.ctx.resume();
    return this.ctx;
  }

  private noiseBuffer(ctx: AudioContext, seconds: number, brown: boolean): AudioBuffer {
    const buffer = ctx.createBuffer(1, ctx.sampleRate * seconds, ctx.sampleRate);
    const data = buffer.getChannelData(0);
    let last = 0;
    for (let i = 0; i < data.length; i++) {
      const white = Math.random() * 2 - 1;
      if (brown) {
        // Brown noise: integrate white noise for a deep, rumbling character.
        last = (last + 0.02 * white) / 1.02;
        data[i] = last * 3.5;
      } else {
        data[i] = white;
      }
    }
    return buffer;
  }

  private startNodes() {
    const ctx = this.ensureContext();
    const master = this.master!;
    master.gain.value = this.state.volume * 0.6;

    // Rain bed: looping white noise through a gentle lowpass.
    const rain = ctx.createBufferSource();
    rain.buffer = this.noiseBuffer(ctx, 4, false);
    rain.loop = true;
    const rainFilter = ctx.createBiquadFilter();
    rainFilter.type = "lowpass";
    rainFilter.frequency.value = this.state.mode === "storm" ? 700 : 1100;
    const rainGain = ctx.createGain();
    rainGain.gain.value = this.state.mode === "storm" ? 0.5 : 0.4;
    rain.connect(rainFilter).connect(rainGain).connect(master);
    rain.start();
    this.rainSource = rain;

    if (this.state.mode === "storm") this.scheduleThunder();
  }

  private scheduleThunder() {
    const delay = 4000 + Math.random() * 14000;
    this.thunderTimer = setTimeout(() => {
      this.playThunder();
      if (this.state.playing && this.state.mode === "storm") this.scheduleThunder();
    }, delay);
  }

  private playThunder() {
    if (!this.ctx || !this.master) return;
    const ctx = this.ctx;
    const src = ctx.createBufferSource();
    src.buffer = this.noiseBuffer(ctx, 3, true);
    const filter = ctx.createBiquadFilter();
    filter.type = "lowpass";
    filter.frequency.value = 120;
    const gain = ctx.createGain();
    const now = ctx.currentTime;
    gain.gain.setValueAtTime(0, now);
    gain.gain.linearRampToValueAtTime(0.9, now + 0.15);
    gain.gain.exponentialRampToValueAtTime(0.001, now + 2.8);
    src.connect(filter).connect(gain).connect(this.master);
    src.start(now);
    src.stop(now + 3);
  }

  private stopNodes() {
    if (this.thunderTimer) {
      clearTimeout(this.thunderTimer);
      this.thunderTimer = null;
    }
    if (this.rainSource) {
      try {
        this.rainSource.stop();
      } catch {
        // already stopped
      }
      this.rainSource = null;
    }
  }
}

let engine: AmbientEngine | null = null;

export function getAmbientEngine(): AmbientEngine {
  if (!engine) engine = new AmbientEngine();
  return engine;
}
