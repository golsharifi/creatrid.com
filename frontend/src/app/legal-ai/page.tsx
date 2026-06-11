"use client";

import { useState } from "react";
import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { useTranslation } from "react-i18next";
import { api } from "@/lib/api";
import { AiMarkdown } from "@/components/ai-markdown";
import { AiQuotaBadge } from "@/components/ai-quota-badge";
import { FileText, ArrowRight, Zap, Shield } from "@/components/icons";

const ACCENT = "var(--tab-legal)";

// Built-in guided prompts, as requested in the blueprint ("built in ai prompts guiding").
const GUIDED_PROMPTS = [
  {
    label: "Copyright basics",
    prompt:
      "I'm an independent creator. Explain in plain language what copyright I automatically own in the content I make, what registration adds, and the 3 most common ways creators lose or weaken their rights.",
  },
  {
    label: "Draft a DMCA takedown",
    prompt:
      "Someone re-uploaded my work without permission. Draft a complete DMCA takedown notice I can send to the hosting platform, with bracketed placeholders for my details, the infringing URL, and the original work.",
  },
  {
    label: "License my work",
    prompt:
      "I want to license my digital content (images/audio/video) to buyers while keeping ownership. Explain the difference between exclusive and non-exclusive licenses and draft a simple non-exclusive license agreement template.",
  },
  {
    label: "Collab agreement",
    prompt:
      "I'm collaborating with another creator on a joint project. Draft a short collaboration agreement covering ownership split, revenue split, credit, and what happens if one of us leaves.",
  },
  {
    label: "Cease & desist",
    prompt:
      "Draft a firm but professional cease-and-desist letter for someone impersonating my creator brand on social media, with bracketed placeholders.",
  },
  {
    label: "NDA template",
    prompt:
      "Draft a simple mutual NDA I can use before sharing unreleased work with a potential brand partner, with bracketed placeholders.",
  },
];

export default function LegalAIPage() {
  const { user } = useAuth();
  const { t } = useTranslation();
  const [question, setQuestion] = useState("");
  const [answer, setAnswer] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [quotaKey, setQuotaKey] = useState(0);

  async function ask(q: string) {
    if (!q.trim() || loading) return;
    setQuestion(q);
    setLoading(true);
    setError("");
    setAnswer("");
    const res = await api.ai.legal(q.trim());
    setLoading(false);
    if (res.error || !res.data) {
      setError(res.error || "Something went wrong");
      return;
    }
    setAnswer(res.data.text);
    setQuotaKey((k) => k + 1);
  }

  if (!user) {
    return (
      <div className="mx-auto max-w-2xl px-4 py-24 text-center sm:px-6">
        <div
          className="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl text-white shadow-lg"
          style={{ background: ACCENT }}
        >
          <FileText className="h-8 w-8" />
        </div>
        <h1 className="mt-6 text-3xl font-bold tracking-tight">Legal AI Assistant</h1>
        <p className="mx-auto mt-4 max-w-xl text-zinc-600 dark:text-zinc-400">
          Guided help with copyright, licensing, takedowns, and contracts — in
          plain language. Sign in to ask your first question.
        </p>
        <Link
          href="/sign-in"
          className="mt-7 inline-flex items-center gap-2 rounded-lg px-6 py-3 text-sm font-medium text-white transition-opacity hover:opacity-90"
          style={{ background: ACCENT }}
        >
          Sign in to get started
          <ArrowRight className="h-4 w-4" />
        </Link>
      </div>
    );
  }

  return (
    <div className="mx-auto w-full max-w-4xl px-4 py-12 sm:px-6">
      <div className="flex items-center gap-3">
        <div
          className="flex h-12 w-12 items-center justify-center rounded-xl text-white shadow"
          style={{ background: ACCENT }}
        >
          <FileText className="h-6 w-6" />
        </div>
        <div>
          <h1 className="text-2xl font-bold tracking-tight">{t("tabs.legalTitle")}</h1>
          <p className="text-sm text-zinc-500 dark:text-zinc-400">
            {t("tabs.legalSubtitle")}
          </p>
        </div>
        <div className="ml-auto">
          <AiQuotaBadge refreshKey={quotaKey} />
        </div>
      </div>

      {/* Guided prompts */}
      <p className="mt-8 text-sm font-medium text-zinc-500 dark:text-zinc-400">
        Start with a guided prompt:
      </p>
      <div className="mt-3 flex flex-wrap gap-2">
        {GUIDED_PROMPTS.map((g) => (
          <button
            key={g.label}
            onClick={() => ask(g.prompt)}
            disabled={loading}
            className="rounded-full border border-zinc-200 bg-white px-4 py-2 text-sm font-medium text-zinc-700 transition-colors hover:border-zinc-300 hover:bg-zinc-50 disabled:opacity-50 dark:border-zinc-800 dark:bg-zinc-950 dark:text-zinc-300 dark:hover:bg-zinc-900"
          >
            {g.label}
          </button>
        ))}
      </div>

      {/* Free-form question */}
      <div className="mt-6">
        <textarea
          value={question}
          onChange={(e) => setQuestion(e.target.value)}
          rows={3}
          placeholder="Or ask your own question… e.g. A brand wants to use my video in their ad — what should the license say?"
          className="w-full rounded-xl border border-zinc-200 bg-white p-3 text-sm outline-none focus:border-zinc-400 dark:border-zinc-800 dark:bg-zinc-950 dark:focus:border-zinc-600"
        />
        <button
          onClick={() => ask(question)}
          disabled={loading || !question.trim()}
          className="mt-2 inline-flex items-center gap-2 rounded-lg px-5 py-2.5 text-sm font-medium text-white transition-opacity hover:opacity-90 disabled:opacity-50"
          style={{ background: ACCENT }}
        >
          <Zap className="h-4 w-4" />
          {loading ? "Thinking…" : "Ask"}
        </button>
      </div>

      {error && <p className="mt-4 text-sm text-red-600 dark:text-red-400">{error}</p>}

      {answer && (
        <div className="mt-8 rounded-2xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
          <AiMarkdown text={answer} />
        </div>
      )}

      <div className="mt-10 flex items-start gap-3 rounded-xl border border-dashed border-zinc-300 bg-zinc-50 p-4 text-xs text-zinc-500 dark:border-zinc-700 dark:bg-zinc-900/50 dark:text-zinc-400">
        <Shield className="mt-0.5 h-4 w-4 shrink-0" />
        <p>
          This assistant provides general legal information and document drafts — not
          legal advice, and no attorney–client relationship is created. Laws vary by
          jurisdiction; consult a licensed attorney before acting. Tip: your Vault&apos;s
          blockchain anchors are timestamped proof of authorship you can attach as
          evidence.
        </p>
      </div>
    </div>
  );
}
