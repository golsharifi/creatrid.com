"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState, Suspense } from "react";
import { api } from "@/lib/api";
import Link from "next/link";
import { ArrowLeft, Trash2, Image, Video, Music, FileText } from "@/components/icons";
import { useTranslation } from "react-i18next";

type CollectionDetail = {
  id: string;
  userId: string;
  title: string;
  description: string | null;
  isPublic: boolean;
  itemCount: number;
};

type CollectionItem = {
  contentId: string;
  title: string;
  contentType: string;
  thumbnailUrl: string | null;
  position: number;
};

function typeIcon(type: string) {
  switch (type) {
    case "image": return <Image className="h-5 w-5" />;
    case "video": return <Video className="h-5 w-5" />;
    case "audio": return <Music className="h-5 w-5" />;
    default: return <FileText className="h-5 w-5" />;
  }
}

function CollectionDetailContent() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();
  const { t } = useTranslation();
  const id = searchParams.get("id");

  const [collection, setCollection] = useState<CollectionDetail | null>(null);
  const [items, setItems] = useState<CollectionItem[]>([]);
  const [loadingData, setLoadingData] = useState(true);

  useEffect(() => {
    if (!id) return;
    api.collections.get(id).then((result) => {
      if (result.data) {
        const data = result.data as any;
        setCollection(data.collection);
        setItems(data.items ?? []);
      }
      setLoadingData(false);
    });
  }, [id]);

  const isOwner = user && collection && user.id === collection.userId;

  const handleRemoveItem = async (contentId: string) => {
    if (!id) return;
    await api.collections.removeItem(id, contentId);
    setItems((prev) => prev.filter((i) => i.contentId !== contentId));
  };

  if (loadingData) {
    return (
      <div className="flex items-center justify-center py-32">
        <div className="h-6 w-6 animate-spin rounded-full border-2 border-zinc-300 border-t-zinc-900 dark:border-zinc-600 dark:border-t-zinc-100" />
      </div>
    );
  }

  if (!collection) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-16 text-center">
        <p className="text-zinc-500 dark:text-zinc-400">{t("collections.notFound")}</p>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-4xl px-4 py-12 sm:px-6">
      <Link href="/collections" className="mb-6 inline-flex items-center gap-1 text-sm text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100">
        <ArrowLeft className="h-4 w-4" />
        {t("collections.backToCollections")}
      </Link>

      <div className="mb-8">
        <h1 className="text-2xl font-bold tracking-tight text-zinc-900 dark:text-zinc-100">{collection.title}</h1>
        {collection.description && <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">{collection.description}</p>}
      </div>

      {items.length === 0 ? (
        <div className="rounded-xl border border-zinc-200 bg-zinc-50 py-16 text-center dark:border-zinc-800 dark:bg-zinc-900">
          <p className="text-sm text-zinc-500 dark:text-zinc-400">{t("collections.noItems")}</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {items.map((item) => (
            <div key={item.contentId} className="group relative rounded-xl border border-zinc-200 bg-white p-4 dark:border-zinc-800 dark:bg-zinc-950">
              {item.thumbnailUrl ? (
                <img src={item.thumbnailUrl} alt={item.title} className="mb-3 h-32 w-full rounded-lg object-cover" />
              ) : (
                <div className="mb-3 flex h-32 items-center justify-center rounded-lg bg-zinc-100 dark:bg-zinc-900">
                  {typeIcon(item.contentType)}
                </div>
              )}
              <p className="text-sm font-medium text-zinc-900 dark:text-zinc-100">{item.title}</p>
              <p className="text-xs text-zinc-500 dark:text-zinc-400">{item.contentType}</p>

              {isOwner && (
                <button onClick={() => handleRemoveItem(item.contentId)} className="absolute right-2 top-2 rounded-lg bg-white/80 p-1.5 text-zinc-400 opacity-0 transition-opacity hover:text-red-500 group-hover:opacity-100 dark:bg-zinc-950/80">
                  <Trash2 className="h-4 w-4" />
                </button>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default function CollectionDetailPage() {
  return (
    <Suspense fallback={null}>
      <CollectionDetailContent />
    </Suspense>
  );
}
