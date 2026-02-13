"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect, useState, Suspense } from "react";
import { api } from "@/lib/api";
import Link from "next/link";
import { FolderOpen, Plus, Globe, Lock } from "lucide-react";
import { useTranslation } from "react-i18next";

type Collection = {
  id: string;
  title: string;
  description: string | null;
  isPublic: boolean;
  itemCount: number;
  createdAt: string;
};

function CollectionsContent() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const { t } = useTranslation();
  const [collections, setCollections] = useState<Collection[]>([]);
  const [loadingData, setLoadingData] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [isPublic, setIsPublic] = useState(true);
  const [creating, setCreating] = useState(false);

  useEffect(() => {
    if (!loading && !user) router.push("/sign-in");
  }, [user, loading, router]);

  useEffect(() => {
    if (!user) return;
    api.collections.list().then((result) => {
      if (result.data) setCollections((result.data as any).collections ?? []);
      setLoadingData(false);
    });
  }, [user]);

  const handleCreate = async () => {
    if (!title.trim()) return;
    setCreating(true);
    const result = await api.collections.create(title.trim(), description.trim() || undefined, isPublic);
    if (result.data) {
      setCollections((prev) => [(result.data as any), ...prev]);
      setTitle("");
      setDescription("");
      setShowCreate(false);
    }
    setCreating(false);
  };

  if (loading || !user) return null;

  return (
    <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-zinc-900 dark:text-zinc-100">{t("collections.title")}</h1>
          <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">{t("collections.subtitle")}</p>
        </div>
        <button onClick={() => setShowCreate(!showCreate)} className="flex items-center gap-2 rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white hover:bg-zinc-800 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200">
          <Plus className="h-4 w-4" />
          {t("collections.create")}
        </button>
      </div>

      {/* Create form */}
      {showCreate && (
        <div className="mb-8 rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-950">
          <div className="space-y-4">
            <div>
              <label className="mb-1 block text-sm font-medium text-zinc-700 dark:text-zinc-300">{t("collections.titleLabel")}</label>
              <input value={title} onChange={(e) => setTitle(e.target.value)} className="w-full rounded-lg border border-zinc-300 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-100" placeholder={t("collections.titlePlaceholder")} />
            </div>
            <div>
              <label className="mb-1 block text-sm font-medium text-zinc-700 dark:text-zinc-300">{t("collections.descriptionLabel")}</label>
              <textarea value={description} onChange={(e) => setDescription(e.target.value)} rows={2} className="w-full rounded-lg border border-zinc-300 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-100" placeholder={t("collections.descriptionPlaceholder")} />
            </div>
            <label className="flex items-center gap-2 text-sm text-zinc-700 dark:text-zinc-300">
              <input type="checkbox" checked={isPublic} onChange={(e) => setIsPublic(e.target.checked)} className="rounded" />
              {t("collections.publicLabel")}
            </label>
            <button onClick={handleCreate} disabled={creating || !title.trim()} className="rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white hover:bg-zinc-800 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200">
              {creating ? t("common.loading") : t("collections.createButton")}
            </button>
          </div>
        </div>
      )}

      {/* Collections list */}
      {loadingData ? (
        <div className="flex items-center justify-center py-16">
          <div className="h-6 w-6 animate-spin rounded-full border-2 border-zinc-300 border-t-zinc-900 dark:border-zinc-600 dark:border-t-zinc-100" />
        </div>
      ) : collections.length === 0 ? (
        <div className="rounded-xl border border-zinc-200 bg-zinc-50 py-16 text-center dark:border-zinc-800 dark:bg-zinc-900">
          <FolderOpen className="mx-auto h-10 w-10 text-zinc-300 dark:text-zinc-600" />
          <p className="mt-4 text-sm font-medium text-zinc-500 dark:text-zinc-400">{t("collections.empty")}</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {collections.map((coll) => (
            <Link key={coll.id} href={`/collections/detail?id=${coll.id}`} className="group rounded-xl border border-zinc-200 bg-white p-5 transition-colors hover:border-zinc-300 dark:border-zinc-800 dark:bg-zinc-950 dark:hover:border-zinc-700">
              <div className="flex items-start justify-between">
                <div>
                  <h3 className="font-semibold text-zinc-900 group-hover:text-zinc-700 dark:text-zinc-100 dark:group-hover:text-zinc-300">{coll.title}</h3>
                  {coll.description && <p className="mt-1 line-clamp-2 text-sm text-zinc-500 dark:text-zinc-400">{coll.description}</p>}
                </div>
                {coll.isPublic ? <Globe className="h-4 w-4 text-green-500" /> : <Lock className="h-4 w-4 text-zinc-400" />}
              </div>
              <div className="mt-3 flex items-center gap-3 text-xs text-zinc-500 dark:text-zinc-400">
                <span>{t("collections.itemCount", { count: coll.itemCount })}</span>
                <span>{new Date(coll.createdAt).toLocaleDateString()}</span>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}

export default function CollectionsPage() {
  return (
    <Suspense fallback={null}>
      <CollectionsContent />
    </Suspense>
  );
}
