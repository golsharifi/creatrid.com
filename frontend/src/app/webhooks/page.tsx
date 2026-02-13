"use client";

import { useState, useEffect, useCallback } from "react";
import { useAuth } from "@/lib/auth-context";
import { api } from "@/lib/api";
import { useTranslation } from "react-i18next";
import { Webhook, Plus, Trash2, ToggleLeft, ToggleRight, ChevronDown, ChevronUp } from "lucide-react";

const AVAILABLE_EVENTS = [
  { value: "license.sold", label: "License Sold" },
  { value: "content.uploaded", label: "Content Uploaded" },
  { value: "collaboration.received", label: "Collaboration Received" },
  { value: "payout.completed", label: "Payout Completed" },
];

export default function WebhooksPage() {
  const { user, loading: authLoading } = useAuth();
  const { t } = useTranslation();
  const [endpoints, setEndpoints] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [newUrl, setNewUrl] = useState("");
  const [newEvents, setNewEvents] = useState<string[]>([]);
  const [creating, setCreating] = useState(false);
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [deliveries, setDeliveries] = useState<any[]>([]);

  const loadEndpoints = useCallback(async () => {
    const res = await api.webhooks.list();
    if (res.data) setEndpoints(res.data.endpoints || []);
    setLoading(false);
  }, []);

  useEffect(() => {
    if (user) loadEndpoints();
  }, [user, loadEndpoints]);

  const handleCreate = async () => {
    if (!newUrl || newEvents.length === 0) return;
    setCreating(true);
    const res = await api.webhooks.create(newUrl, newEvents);
    if (res.data) {
      setShowCreate(false);
      setNewUrl("");
      setNewEvents([]);
      loadEndpoints();
    }
    setCreating(false);
  };

  const handleDelete = async (id: string) => {
    await api.webhooks.delete(id);
    loadEndpoints();
  };

  const handleToggle = async (id: string, isActive: boolean) => {
    await api.webhooks.update(id, { isActive: !isActive });
    loadEndpoints();
  };

  const handleExpand = async (id: string) => {
    if (expandedId === id) {
      setExpandedId(null);
      return;
    }
    setExpandedId(id);
    const res = await api.webhooks.deliveries(id);
    if (res.data) setDeliveries(res.data.deliveries || []);
  };

  if (authLoading || loading) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-2 border-zinc-300 border-t-zinc-900 dark:border-zinc-700 dark:border-t-zinc-100" />
      </div>
    );
  }

  if (!user) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-16 text-center">
        <p className="text-zinc-500">{t("common.signInRequired")}</p>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-4xl px-4 py-8 sm:px-6">
      <div className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-100">{t("webhooks.title")}</h1>
          <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">{t("webhooks.subtitle")}</p>
        </div>
        <button
          onClick={() => setShowCreate(!showCreate)}
          className="flex items-center gap-2 rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white hover:bg-zinc-800 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
        >
          <Plus className="h-4 w-4" />
          {t("webhooks.create")}
        </button>
      </div>

      {showCreate && (
        <div className="mb-6 rounded-xl border border-zinc-200 bg-white p-6 dark:border-zinc-800 dark:bg-zinc-900">
          <h3 className="mb-4 text-lg font-semibold text-zinc-900 dark:text-zinc-100">{t("webhooks.newEndpoint")}</h3>
          <div className="mb-4">
            <label className="mb-1 block text-sm font-medium text-zinc-700 dark:text-zinc-300">URL</label>
            <input
              type="url"
              value={newUrl}
              onChange={(e) => setNewUrl(e.target.value)}
              placeholder="https://example.com/webhook"
              className="w-full rounded-lg border border-zinc-300 bg-white px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-100"
            />
          </div>
          <div className="mb-4">
            <label className="mb-2 block text-sm font-medium text-zinc-700 dark:text-zinc-300">{t("webhooks.events")}</label>
            <div className="space-y-2">
              {AVAILABLE_EVENTS.map((event) => (
                <label key={event.value} className="flex items-center gap-2 text-sm text-zinc-700 dark:text-zinc-300">
                  <input
                    type="checkbox"
                    checked={newEvents.includes(event.value)}
                    onChange={(e) => {
                      if (e.target.checked) setNewEvents([...newEvents, event.value]);
                      else setNewEvents(newEvents.filter((ev) => ev !== event.value));
                    }}
                    className="rounded border-zinc-300 dark:border-zinc-600"
                  />
                  {event.label}
                </label>
              ))}
            </div>
          </div>
          <button
            onClick={handleCreate}
            disabled={creating || !newUrl || newEvents.length === 0}
            className="rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white hover:bg-zinc-800 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
          >
            {creating ? t("common.saving") : t("webhooks.save")}
          </button>
        </div>
      )}

      {endpoints.length === 0 ? (
        <div className="rounded-xl border border-zinc-200 bg-white p-12 text-center dark:border-zinc-800 dark:bg-zinc-900">
          <Webhook className="mx-auto mb-4 h-12 w-12 text-zinc-300 dark:text-zinc-600" />
          <p className="text-zinc-500 dark:text-zinc-400">{t("webhooks.empty")}</p>
        </div>
      ) : (
        <div className="space-y-4">
          {endpoints.map((ep) => (
            <div key={ep.id} className="rounded-xl border border-zinc-200 bg-white dark:border-zinc-800 dark:bg-zinc-900">
              <div className="flex items-center gap-4 p-4">
                <button onClick={() => handleToggle(ep.id, ep.isActive)} className="text-zinc-500 dark:text-zinc-400">
                  {ep.isActive ? <ToggleRight className="h-6 w-6 text-emerald-500" /> : <ToggleLeft className="h-6 w-6" />}
                </button>
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm font-medium text-zinc-900 dark:text-zinc-100">{ep.url}</p>
                  <p className="text-xs text-zinc-500 dark:text-zinc-400">{(ep.events || []).join(", ")}</p>
                </div>
                <button
                  onClick={() => handleExpand(ep.id)}
                  className="rounded-lg p-2 text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800"
                >
                  {expandedId === ep.id ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                </button>
                <button
                  onClick={() => handleDelete(ep.id)}
                  className="rounded-lg p-2 text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20"
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              </div>
              {expandedId === ep.id && (
                <div className="border-t border-zinc-200 p-4 dark:border-zinc-800">
                  <h4 className="mb-2 text-sm font-medium text-zinc-700 dark:text-zinc-300">{t("webhooks.recentDeliveries")}</h4>
                  {deliveries.length === 0 ? (
                    <p className="text-sm text-zinc-400">{t("webhooks.noDeliveries")}</p>
                  ) : (
                    <div className="space-y-2">
                      {deliveries.slice(0, 10).map((d: any) => (
                        <div key={d.id} className="flex items-center justify-between rounded-lg bg-zinc-50 px-3 py-2 text-sm dark:bg-zinc-800">
                          <span className="text-zinc-700 dark:text-zinc-300">{d.eventType}</span>
                          <span className={d.responseStatus >= 200 && d.responseStatus < 300 ? "text-emerald-600" : "text-red-500"}>
                            {d.responseStatus || "pending"}
                          </span>
                          <span className="text-xs text-zinc-400">{new Date(d.createdAt).toLocaleString()}</span>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
