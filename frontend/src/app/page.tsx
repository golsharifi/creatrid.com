"use client";

import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { Shield, Users, CheckCircle, ArrowRight } from "lucide-react";

export default function Home() {
  const { user } = useAuth();

  return (
    <div className="flex flex-col">
      {/* Hero */}
      <section className="mx-auto flex w-full max-w-6xl flex-col items-center gap-8 px-4 py-24 text-center sm:px-6 sm:py-32">
        <div className="inline-flex items-center gap-2 rounded-full border border-zinc-200 bg-zinc-50 px-4 py-1.5 text-sm text-zinc-600 dark:border-zinc-800 dark:bg-zinc-900 dark:text-zinc-400">
          <Shield className="h-4 w-4" />
          Verified Creator Identity
        </div>

        <h1 className="max-w-3xl text-4xl font-bold tracking-tight sm:text-6xl">
          Prove who you are.
          <br />
          <span className="text-zinc-400 dark:text-zinc-500">
            Own your reputation.
          </span>
        </h1>

        <p className="max-w-xl text-lg text-zinc-600 dark:text-zinc-400">
          Your unified creator profile — verified across platforms. Connect your
          accounts, build your reputation score, and let brands know you&apos;re
          the real deal.
        </p>

        <div className="flex flex-col gap-3 sm:flex-row">
          {user ? (
            <Link
              href="/dashboard"
              className="inline-flex items-center gap-2 rounded-lg bg-zinc-900 px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-zinc-800 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              Go to Dashboard
              <ArrowRight className="h-4 w-4" />
            </Link>
          ) : (
            <Link
              href="/sign-in"
              className="inline-flex items-center gap-2 rounded-lg bg-zinc-900 px-6 py-3 text-sm font-medium text-white transition-colors hover:bg-zinc-800 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              Get Your Creator Passport
              <ArrowRight className="h-4 w-4" />
            </Link>
          )}
        </div>
      </section>

      {/* Features */}
      <section className="border-t border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/50">
        <div className="mx-auto grid max-w-6xl gap-8 px-4 py-24 sm:grid-cols-3 sm:px-6">
          <div className="flex flex-col gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900">
              <Shield className="h-5 w-5" />
            </div>
            <h3 className="text-lg font-semibold">Verified Identity</h3>
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              Connect your social accounts and prove you are who you say you
              are. No fakes, no impersonators.
            </p>
          </div>

          <div className="flex flex-col gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900">
              <Users className="h-5 w-5" />
            </div>
            <h3 className="text-lg font-semibold">Creator Score</h3>
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              A reputation score based on your verified presence across
              platforms. The more you connect, the stronger your profile.
            </p>
          </div>

          <div className="flex flex-col gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900">
              <CheckCircle className="h-5 w-5" />
            </div>
            <h3 className="text-lg font-semibold">Trusted by Brands</h3>
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              Brands and agencies verify creators before hiring. Your Creatrid
              profile is your proof of legitimacy.
            </p>
          </div>
        </div>
      </section>
    </div>
  );
}
