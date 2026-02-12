"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { SettingsForm } from "./settings-form";

export default function SettingsPage() {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
  }, [user, loading, router]);

  if (loading || !user) return null;

  return (
    <div className="mx-auto max-w-2xl px-4 py-12 sm:px-6">
      <h1 className="text-2xl font-bold">Profile Settings</h1>
      <p className="mt-1 text-sm text-zinc-500">
        Manage your Creator Passport profile.
      </p>

      <div className="mt-8">
        <SettingsForm
          name={user.name || ""}
          username={user.username || ""}
          bio={user.bio || ""}
          image={user.image || ""}
          theme={user.theme || "default"}
          customLinks={user.customLinks || []}
        />
      </div>
    </div>
  );
}
