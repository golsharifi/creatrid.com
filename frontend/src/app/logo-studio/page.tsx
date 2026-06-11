"use client";

import { useState } from "react";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { useTranslation } from "react-i18next";
import { api } from "@/lib/api";
import { AiMarkdown } from "@/components/ai-markdown";
import { AiQuotaBadge } from "@/components/ai-quota-badge";
import {
  Image as ImageIcon,
  Download,
  ArrowRight,
  FileText,
  Edit3,
  Zap,
} from "@/components/icons";

type LogoConcept = { name: string; rationale: string; svg: string };

const ACCENT = "var(--tab-studio)";

const STUDIO_TABS = [
  { key: "logos", label: "Logo Concepts", Icon: ImageIcon },
  { key: "copy", label: "Marketing Copy", Icon: FileText },
  { key: "refine", label: "Refine Content", Icon: Edit3 },
] as const;

type StudioTab = (typeof STUDIO_TABS)[number]["key"];

export default function LogoStudioPage() {
  const { user } = useAuth();
  const { t } = useTranslation();
  const [tab, setTab] = useState<StudioTab>("logos");

  // Logo generation state
  const [brief, setBrief] = useState("");
  const [concepts, setConcepts] = useState<LogoConcept[]>([]);
  const [logoLoading, setLogoLoading] = useState(false);
  const [logoError, setLogoError] = useState("");

  // Copy generation state
  const [copyBrief, setCopyBrief] = useState("");
  const [copyResult, setCopyResult] = useState("");
  const [copyLoading, setCopyLoading] = useState(false);
  const [copyError, setCopyError] = useState("");

  // Refine state
  const [refineText, setRefineText] = useState("");
  const [refineInstruction, setRefineInstruction] = useState("");
  const [refineResult, setRefineResult] = useState("");
  const [refineLoading, setRefineLoading] = useState(false);
  const [refineError, setRefineError] = useState("");

  // Quota meter refresh + watermark status
  const [quotaKey, setQuotaKey] = useState(0);
  const [watermarkStatus, setWatermarkStatus] = useState<string | null>(null);

  async function generateLogos() {
    if (!brief.trim() || logoLoading) return;
    setLogoLoading(true);
    setLogoError("");
    const res = await api.ai.generateLogos(brief.trim());
    setLogoLoading(false);
    if (res.error || !res.data) {
      setLogoError(res.error || "Something went wrong");
      return;
    }
    setConcepts(res.data.concepts);
    setQuotaKey((k) => k + 1);
  }

  async function generateCopy() {
    if (!copyBrief.trim() || copyLoading) return;
    setCopyLoading(true);
    setCopyError("");
    const res = await api.ai.generateCopy(copyBrief.trim());
    setCopyLoading(false);
    if (res.error || !res.data) {
      setCopyError(res.error || "Something went wrong");
      return;
    }
    setCopyResult(res.data.text);
    setQuotaKey((k) => k + 1);
  }

  async function runRefine() {
    if (!refineText.trim() || refineLoading) return;
    setRefineLoading(true);
    setRefineError("");
    const res = await api.ai.refineText(
      refineText.trim(),
      refineInstruction.trim() || "Improve clarity and impact."
    );
    setRefineLoading(false);
    if (res.error || !res.data) {
      setRefineError(res.error || "Something went wrong");
      return;
    }
    setRefineResult(res.data.text);
    setQuotaKey((k) => k + 1);
  }

  // Rasterize an SVG concept to a 512px PNG via canvas.
  async function rasterize(concept: LogoConcept): Promise<Blob> {
    const svgBlob = new Blob([concept.svg], { type: "image/svg+xml" });
    const url = URL.createObjectURL(svgBlob);
    const img = new Image();
    await new Promise<void>((resolve, reject) => {
      img.onload = () => resolve();
      img.onerror = () => reject(new Error("SVG render failed"));
      img.src = url;
    });
    const canvas = document.createElement("canvas");
    canvas.width = 512;
    canvas.height = 512;
    const ctx = canvas.getContext("2d");
    if (!ctx) throw new Error("Canvas unavailable");
    ctx.drawImage(img, 0, 0, 512, 512);
    URL.revokeObjectURL(url);
    return new Promise((resolve, reject) =>
      canvas.toBlob((b) => (b ? resolve(b) : reject(new Error("PNG export failed"))), "image/png")
    );
  }

  // Save a concept as the user's vault watermark (stamped on image uploads).
  async function setAsWatermark(concept: LogoConcept) {
    setWatermarkStatus("Setting watermark…");
    try {
      const png = await rasterize(concept);
      const res = await api.users.uploadWatermark(png);
      setWatermarkStatus(
        res.error ? `Failed: ${res.error}` : `"${concept.name}" is now your watermark ✓`
      );
    } catch (e) {
      setWatermarkStatus(`Failed: ${e instanceof Error ? e.message : "unknown error"}`);
    }
  }

  // Anchor a concept on-chain as the creator's passport marker: the PNG goes
  // into the Vault as a public item, then the existing blockchain anchoring
  // engine stamps its hash on-chain — timestamped proof this mark is theirs.
  async function anchorAsMarker(concept: LogoConcept) {
    setWatermarkStatus("Anchoring your passport marker on-chain…");
    try {
      const png = await rasterize(concept);
      const upload = await api.content.upload(
        png,
        `Passport Marker — ${concept.name}`,
        "My official logo mark, anchored on-chain as my passport identifier.",
        ["passport-marker", "logo"],
        true,
        { filename: "passport-marker.png" }
      );
      if (upload.error || !upload.data?.item?.id) {
        throw new Error(upload.error || "Vault upload failed");
      }
      const anchor = await api.blockchain.anchor(upload.data.item.id);
      if (anchor.error) {
        setWatermarkStatus(
          `Marker saved to your Vault, but anchoring is unavailable: ${anchor.error}`
        );
        return;
      }
      setWatermarkStatus(
        `"${concept.name}" anchored on-chain as your passport marker ✓ — see it in your Vault`
      );
    } catch (e) {
      setWatermarkStatus(`Failed: ${e instanceof Error ? e.message : "unknown error"}`);
    }
  }

  function downloadSVG(concept: LogoConcept) {
    const blob = new Blob([concept.svg], { type: "image/svg+xml" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${concept.name.toLowerCase().replace(/[^a-z0-9]+/g, "-")}.svg`;
    a.click();
    URL.revokeObjectURL(url);
  }

  if (!user) {
    return (
      <div className="mx-auto max-w-2xl px-4 py-24 text-center sm:px-6">
        <div
          className="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl text-white shadow-lg"
          style={{ background: ACCENT }}
        >
          <ImageIcon className="h-8 w-8" />
        </div>
        <h1 className="mt-6 text-3xl font-bold tracking-tight">Logo Studio</h1>
        <p className="mx-auto mt-4 max-w-xl text-zinc-600 dark:text-zinc-400">
          Generate logo concepts, marketing copy, and refined content with your
          AI brand assistant. Sign in to start designing your mark.
        </p>
        <Link
          href="/sign-in"
          className="mt-7 inline-flex items-center gap-2 rounded-lg px-6 py-3 text-sm font-medium text-white transition-opacity hover:opacity-90"
          style={{ background: ACCENT }}
        >
          Sign in to use the Studio
          <ArrowRight className="h-4 w-4" />
        </Link>
      </div>
    );
  }

  return (
    <div className="mx-auto w-full max-w-5xl px-4 py-12 sm:px-6">
      <div className="flex items-center gap-3">
        <div
          className="flex h-12 w-12 items-center justify-center rounded-xl text-white shadow"
          style={{ background: ACCENT }}
        >
          <ImageIcon className="h-6 w-6" />
        </div>
        <div>
          <h1 className="text-2xl font-bold tracking-tight">{t("tabs.studioTitle")}</h1>
          <p className="text-sm text-zinc-500 dark:text-zinc-400">
            {t("tabs.studioSubtitle")}
          </p>
        </div>
        <div className="ml-auto">
          <AiQuotaBadge refreshKey={quotaKey} />
        </div>
      </div>

      {/* Studio tabs */}
      <div className="mt-8 flex flex-wrap gap-2">
        {STUDIO_TABS.map(({ key, label, Icon }) => (
          <button
            key={key}
            onClick={() => setTab(key)}
            className={`flex items-center gap-2 rounded-full px-4 py-2 text-sm font-medium transition-colors ${
              tab === key
                ? "text-white"
                : "bg-zinc-100 text-zinc-600 hover:bg-zinc-200 dark:bg-zinc-800 dark:text-zinc-400 dark:hover:bg-zinc-700"
            }`}
            style={tab === key ? { background: ACCENT } : undefined}
          >
            <Icon className="h-4 w-4" />
            {label}
          </button>
        ))}
      </div>

      {/* ── Logo concepts ── */}
      {tab === "logos" && (
        <div className="mt-8">
          <label className="text-sm font-medium">
            Describe your brand
            <textarea
              value={brief}
              onChange={(e) => setBrief(e.target.value)}
              rows={3}
              placeholder="e.g. Moody synthwave music producer. Dark, neon, retro-futuristic. The word 'NOVA' or an abstract wave mark."
              className="mt-2 w-full rounded-xl border border-zinc-200 bg-white p-3 text-sm outline-none focus:border-zinc-400 dark:border-zinc-800 dark:bg-zinc-950 dark:focus:border-zinc-600"
            />
          </label>
          <button
            onClick={generateLogos}
            disabled={logoLoading || !brief.trim()}
            className="mt-3 inline-flex items-center gap-2 rounded-lg px-5 py-2.5 text-sm font-medium text-white transition-opacity hover:opacity-90 disabled:opacity-50"
            style={{ background: ACCENT }}
          >
            <Zap className="h-4 w-4" />
            {logoLoading ? "Designing… (can take a minute)" : "Generate 4 concepts"}
          </button>
          {logoError && (
            <p className="mt-3 text-sm text-red-600 dark:text-red-400">{logoError}</p>
          )}
          {watermarkStatus && (
            <p className="mt-3 text-sm text-zinc-600 dark:text-zinc-400">{watermarkStatus}</p>
          )}

          {concepts.length > 0 && (
            <div className="mt-8 grid gap-6 sm:grid-cols-2">
              {concepts.map((c, i) => (
                <div
                  key={i}
                  className="rounded-2xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-950"
                >
                  <div
                    className="mx-auto flex h-40 w-40 items-center justify-center [&_svg]:h-full [&_svg]:w-full"
                    // SVG is validated server-side: shape-only allowlist, no
                    // scripts/handlers/external refs (see backend/internal/ai).
                    dangerouslySetInnerHTML={{ __html: c.svg }}
                  />
                  <h3 className="mt-4 font-semibold">{c.name}</h3>
                  <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">
                    {c.rationale}
                  </p>
                  <div className="mt-3 flex flex-wrap items-center gap-4">
                    <button
                      onClick={() => downloadSVG(c)}
                      className="inline-flex items-center gap-1.5 text-sm font-medium"
                      style={{ color: ACCENT }}
                    >
                      <Download className="h-4 w-4" />
                      Download SVG
                    </button>
                    <button
                      onClick={() => setAsWatermark(c)}
                      className="inline-flex items-center gap-1.5 text-sm font-medium text-zinc-600 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100"
                    >
                      <Zap className="h-4 w-4" />
                      Use as watermark
                    </button>
                    <button
                      onClick={() => anchorAsMarker(c)}
                      className="inline-flex items-center gap-1.5 text-sm font-medium text-zinc-600 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100"
                      title="Save to your Vault and stamp its hash on-chain — your NFT-style passport marker"
                    >
                      ⛓ Anchor as passport marker
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* ── Marketing copy ── */}
      {tab === "copy" && (
        <div className="mt-8">
          <label className="text-sm font-medium">
            What do you need written?
            <textarea
              value={copyBrief}
              onChange={(e) => setCopyBrief(e.target.value)}
              rows={3}
              placeholder="e.g. A punchy profile bio, a tagline, and a launch post for my new photography preset pack."
              className="mt-2 w-full rounded-xl border border-zinc-200 bg-white p-3 text-sm outline-none focus:border-zinc-400 dark:border-zinc-800 dark:bg-zinc-950 dark:focus:border-zinc-600"
            />
          </label>
          <button
            onClick={generateCopy}
            disabled={copyLoading || !copyBrief.trim()}
            className="mt-3 inline-flex items-center gap-2 rounded-lg px-5 py-2.5 text-sm font-medium text-white transition-opacity hover:opacity-90 disabled:opacity-50"
            style={{ background: ACCENT }}
          >
            <Zap className="h-4 w-4" />
            {copyLoading ? "Writing…" : "Generate copy"}
          </button>
          {copyError && (
            <p className="mt-3 text-sm text-red-600 dark:text-red-400">{copyError}</p>
          )}
          {copyResult && (
            <div className="mt-6 rounded-2xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
              <AiMarkdown text={copyResult} />
            </div>
          )}
        </div>
      )}

      {/* ── Refine ── */}
      {tab === "refine" && (
        <div className="mt-8">
          <label className="text-sm font-medium">
            Paste your content
            <textarea
              value={refineText}
              onChange={(e) => setRefineText(e.target.value)}
              rows={6}
              placeholder="Paste a bio, caption, description, or script you want polished…"
              className="mt-2 w-full rounded-xl border border-zinc-200 bg-white p-3 text-sm outline-none focus:border-zinc-400 dark:border-zinc-800 dark:bg-zinc-950 dark:focus:border-zinc-600"
            />
          </label>
          <label className="mt-3 block text-sm font-medium">
            Instruction (optional)
            <input
              value={refineInstruction}
              onChange={(e) => setRefineInstruction(e.target.value)}
              placeholder="e.g. Make it shorter and more confident"
              className="mt-2 w-full rounded-xl border border-zinc-200 bg-white p-3 text-sm outline-none focus:border-zinc-400 dark:border-zinc-800 dark:bg-zinc-950 dark:focus:border-zinc-600"
            />
          </label>
          <button
            onClick={runRefine}
            disabled={refineLoading || !refineText.trim()}
            className="mt-3 inline-flex items-center gap-2 rounded-lg px-5 py-2.5 text-sm font-medium text-white transition-opacity hover:opacity-90 disabled:opacity-50"
            style={{ background: ACCENT }}
          >
            <Zap className="h-4 w-4" />
            {refineLoading ? "Refining…" : "Refine"}
          </button>
          {refineError && (
            <p className="mt-3 text-sm text-red-600 dark:text-red-400">{refineError}</p>
          )}
          {refineResult && (
            <div className="mt-6 rounded-2xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
              <AiMarkdown text={refineResult} />
            </div>
          )}
        </div>
      )}

      <p className="mt-10 rounded-xl border border-dashed border-zinc-300 bg-zinc-50 p-4 text-center text-xs text-zinc-500 dark:border-zinc-700 dark:bg-zinc-900/50 dark:text-zinc-400">
        Next up for the Studio: stamp a generated logo as the watermark on your Vault
        files and anchor it on-chain as your passport marker.
      </p>
    </div>
  );
}
