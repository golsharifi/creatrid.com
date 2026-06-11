"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth-context";
import type { ContentComment } from "@/lib/types";
import { Trash2 } from "@/components/icons";

/**
 * Community comments on public vault content. Posting on someone else's work
 * credits the owner arena points (reputation through engagement).
 */
export function ContentComments({ contentId }: { contentId: string }) {
  const { user } = useAuth();
  const [comments, setComments] = useState<ContentComment[]>([]);
  const [total, setTotal] = useState(0);
  const [body, setBody] = useState("");
  const [posting, setPosting] = useState(false);
  const [error, setError] = useState("");
  const [likes, setLikes] = useState(0);
  const [liked, setLiked] = useState(false);
  const [likeBusy, setLikeBusy] = useState(false);

  const load = useCallback(async () => {
    const res = await api.comments.list(contentId);
    if (res.data) {
      setComments(res.data.comments);
      setTotal(res.data.total);
    }
    if (user) {
      const likeRes = await api.likes.mine(contentId);
      if (likeRes.data) {
        setLikes(likeRes.data.likes);
        setLiked(likeRes.data.liked);
      }
    } else {
      const likeRes = await api.likes.stats(contentId);
      if (likeRes.data) setLikes(likeRes.data.likes);
    }
  }, [contentId, user]);

  useEffect(() => {
    load();
  }, [load]);

  async function toggleLike() {
    if (!user || likeBusy) return;
    setLikeBusy(true);
    const res = liked ? await api.likes.unlike(contentId) : await api.likes.like(contentId);
    setLikeBusy(false);
    if (res.data) {
      setLiked(res.data.liked);
      setLikes(res.data.likes);
    }
  }

  async function post() {
    const text = body.trim();
    if (!text || posting) return;
    setPosting(true);
    setError("");
    const res = await api.comments.create(contentId, text);
    setPosting(false);
    if (res.error) {
      setError(res.error);
      return;
    }
    setBody("");
    load();
  }

  async function remove(id: string) {
    const res = await api.comments.delete(id);
    if (!res.error) load();
  }

  return (
    <div className="mt-8 rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">
          Comments {total > 0 && <span className="text-sm font-normal text-zinc-500">({total})</span>}
        </h2>
        <button
          onClick={toggleLike}
          disabled={!user || likeBusy}
          className={`inline-flex items-center gap-1.5 rounded-full px-3 py-1.5 text-sm font-medium transition-colors ${
            liked
              ? "bg-rose-50 text-rose-600 dark:bg-rose-900/30 dark:text-rose-400"
              : "bg-zinc-100 text-zinc-600 hover:bg-zinc-200 dark:bg-zinc-800 dark:text-zinc-400 dark:hover:bg-zinc-700"
          } disabled:opacity-60`}
          title={user ? (liked ? "Unlike" : "Like this work (+1 reputation to the creator)") : "Sign in to like"}
        >
          {liked ? "♥" : "♡"} {likes}
        </button>
      </div>

      {user ? (
        <div className="mt-4">
          <textarea
            value={body}
            onChange={(e) => setBody(e.target.value)}
            rows={2}
            maxLength={2000}
            placeholder="Share your thoughts on this work…"
            className="w-full rounded-lg border border-zinc-200 bg-white p-3 text-sm outline-none focus:border-zinc-400 dark:border-zinc-700 dark:bg-zinc-950 dark:focus:border-zinc-600"
          />
          <div className="mt-2 flex items-center gap-3">
            <button
              onClick={post}
              disabled={posting || !body.trim()}
              className="rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-zinc-800 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              {posting ? "Posting…" : "Post comment"}
            </button>
            {error && <p className="text-sm text-red-600 dark:text-red-400">{error}</p>}
          </div>
        </div>
      ) : (
        <p className="mt-3 text-sm text-zinc-500 dark:text-zinc-400">
          <Link href="/sign-in" className="font-medium underline underline-offset-2">
            Sign in
          </Link>{" "}
          to join the conversation.
        </p>
      )}

      <ul className="mt-5 flex flex-col gap-4">
        {comments.length === 0 && (
          <li className="text-sm text-zinc-500 dark:text-zinc-400">
            No comments yet — be the first.
          </li>
        )}
        {comments.map((c) => (
          <li key={c.id} className="flex items-start gap-3">
            {c.image ? (
              <img src={c.image} alt="" className="h-8 w-8 rounded-full object-cover" />
            ) : (
              <div className="flex h-8 w-8 items-center justify-center rounded-full bg-zinc-100 text-xs font-bold text-zinc-500 dark:bg-zinc-800">
                {(c.name || c.username || "?")[0].toUpperCase()}
              </div>
            )}
            <div className="min-w-0 flex-1">
              <p className="text-sm">
                {c.username ? (
                  <Link
                    href={`/profile?u=${c.username}`}
                    className="font-medium hover:underline"
                  >
                    {c.name || `@${c.username}`}
                  </Link>
                ) : (
                  <span className="font-medium">{c.name || "Creator"}</span>
                )}
                <span className="ml-2 text-xs text-zinc-400">
                  {new Date(c.createdAt).toLocaleDateString()}
                </span>
              </p>
              <p className="mt-0.5 whitespace-pre-wrap text-sm text-zinc-600 dark:text-zinc-400">
                {c.body}
              </p>
            </div>
            {user && user.id === c.userId && (
              <button
                onClick={() => remove(c.id)}
                className="shrink-0 rounded p-1 text-zinc-400 transition-colors hover:text-red-500"
                aria-label="Delete comment"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
}
