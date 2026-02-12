"use client";

import { useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth-context";
import { Plus, Trash2, Camera } from "lucide-react";
import type { CustomLink } from "@/lib/types";

const THEMES = [
  { key: "default", label: "Default", colors: "bg-zinc-900 dark:bg-zinc-100" },
  { key: "ocean", label: "Ocean", colors: "bg-blue-600" },
  { key: "sunset", label: "Sunset", colors: "bg-orange-500" },
  { key: "forest", label: "Forest", colors: "bg-emerald-600" },
  { key: "midnight", label: "Midnight", colors: "bg-indigo-700" },
  { key: "rose", label: "Rose", colors: "bg-pink-500" },
];

interface SettingsFormProps {
  name: string;
  username: string;
  bio: string;
  image: string;
  theme: string;
  customLinks: CustomLink[];
}

export function SettingsForm({
  name,
  username,
  bio,
  image,
  theme,
  customLinks,
}: SettingsFormProps) {
  const router = useRouter();
  const { refresh } = useAuth();
  const [formData, setFormData] = useState({ name, username, bio, theme });
  const [links, setLinks] = useState<CustomLink[]>(customLinks || []);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [imagePreview, setImagePreview] = useState(image || "");
  const [uploading, setUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  async function handleImageUpload(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!["image/jpeg", "image/png", "image/webp"].includes(file.type)) {
      setError("Only JPEG, PNG, and WebP images are allowed");
      return;
    }
    if (file.size > 5 * 1024 * 1024) {
      setError("Image must be under 5 MB");
      return;
    }

    setUploading(true);
    setError("");

    const result = await api.users.uploadImage(file);
    if (result.error) {
      setError(result.error);
    } else if (result.data) {
      setImagePreview(result.data.image);
      await refresh();
    }
    setUploading(false);
  }

  const addLink = () => {
    if (links.length >= 10) return;
    setLinks([...links, { title: "", url: "" }]);
  };

  const removeLink = (index: number) => {
    setLinks(links.filter((_, i) => i !== index));
  };

  const updateLink = (index: number, field: "title" | "url", value: string) => {
    const updated = [...links];
    updated[index] = { ...updated[index], [field]: value };
    setLinks(updated);
  };

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);

    // Filter out empty links
    const validLinks = links.filter((l) => l.title.trim() && l.url.trim());

    const result = await api.users.updateProfile({
      name: formData.name.trim(),
      username: formData.username.toLowerCase().trim(),
      bio: formData.bio.trim(),
      theme: formData.theme,
      customLinks: validLinks,
    });

    if (result.error) {
      setError(result.error);
      setLoading(false);
    } else {
      await refresh();
      router.push("/dashboard");
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-8">
      {/* Profile Section */}
      <div className="space-y-6">
        <h2 className="text-lg font-semibold">Profile</h2>

        {/* Avatar Upload */}
        <div className="flex items-center gap-4">
          <div className="relative">
            {imagePreview ? (
              <img
                src={imagePreview}
                alt="Profile"
                className="h-20 w-20 rounded-full object-cover"
              />
            ) : (
              <div className="flex h-20 w-20 items-center justify-center rounded-full bg-zinc-100 text-2xl font-bold text-zinc-400 dark:bg-zinc-800">
                {(name || "?")[0].toUpperCase()}
              </div>
            )}
            <button
              type="button"
              onClick={() => fileInputRef.current?.click()}
              disabled={uploading}
              className="absolute -bottom-1 -right-1 rounded-full border-2 border-white bg-zinc-900 p-1.5 text-white transition-colors hover:bg-zinc-700 dark:border-zinc-900 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-300"
            >
              <Camera className="h-3.5 w-3.5" />
            </button>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/jpeg,image/png,image/webp"
              onChange={handleImageUpload}
              className="hidden"
            />
          </div>
          <div className="text-sm">
            <p className="font-medium">{uploading ? "Uploading..." : "Profile Photo"}</p>
            <p className="text-xs text-zinc-500">JPEG, PNG, or WebP. Max 5 MB.</p>
          </div>
        </div>

        <div>
          <label
            htmlFor="name"
            className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
          >
            Display Name
          </label>
          <input
            id="name"
            type="text"
            value={formData.name}
            onChange={(e) =>
              setFormData({ ...formData, name: e.target.value })
            }
            className="mt-1 block w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm focus:border-zinc-400 focus:outline-none focus:ring-1 focus:ring-zinc-400 dark:border-zinc-700 dark:bg-zinc-800"
            required
          />
        </div>

        <div>
          <label
            htmlFor="username"
            className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
          >
            Username
          </label>
          <div className="mt-1 flex items-center rounded-lg border border-zinc-200 bg-white dark:border-zinc-700 dark:bg-zinc-800">
            <span className="pl-3 text-sm text-zinc-400">creatrid.com/</span>
            <input
              id="username"
              type="text"
              value={formData.username}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  username: e.target.value.replace(/[^a-zA-Z0-9_-]/g, ""),
                })
              }
              className="block w-full bg-transparent px-1 py-2 text-sm focus:outline-none"
              required
              minLength={3}
              maxLength={30}
            />
          </div>
        </div>

        <div>
          <label
            htmlFor="bio"
            className="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
          >
            Bio
          </label>
          <textarea
            id="bio"
            value={formData.bio}
            onChange={(e) =>
              setFormData({ ...formData, bio: e.target.value })
            }
            rows={4}
            maxLength={500}
            className="mt-1 block w-full rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm focus:border-zinc-400 focus:outline-none focus:ring-1 focus:ring-zinc-400 dark:border-zinc-700 dark:bg-zinc-800"
            placeholder="Tell the world about yourself..."
          />
          <p className="mt-1 text-xs text-zinc-400">
            {formData.bio.length}/500 characters
          </p>
        </div>
      </div>

      {/* Theme Section */}
      <div className="space-y-4">
        <h2 className="text-lg font-semibold">Profile Theme</h2>
        <p className="text-sm text-zinc-500">
          Choose a color theme for your public profile.
        </p>
        <div className="grid grid-cols-3 gap-3 sm:grid-cols-6">
          {THEMES.map((t) => (
            <button
              key={t.key}
              type="button"
              onClick={() => setFormData({ ...formData, theme: t.key })}
              className={`flex flex-col items-center gap-2 rounded-lg border p-3 transition-colors ${
                formData.theme === t.key
                  ? "border-zinc-900 dark:border-zinc-100"
                  : "border-zinc-200 dark:border-zinc-700"
              }`}
            >
              <div className={`h-8 w-8 rounded-full ${t.colors}`} />
              <span className="text-xs font-medium">{t.label}</span>
            </button>
          ))}
        </div>
      </div>

      {/* Custom Links Section */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-lg font-semibold">Custom Links</h2>
            <p className="text-sm text-zinc-500">
              Add links to your portfolio, website, or socials.
            </p>
          </div>
          {links.length < 10 && (
            <button
              type="button"
              onClick={addLink}
              className="flex items-center gap-1.5 rounded-lg border border-zinc-200 px-3 py-1.5 text-xs font-medium transition-colors hover:bg-zinc-50 dark:border-zinc-700 dark:hover:bg-zinc-800"
            >
              <Plus className="h-3.5 w-3.5" />
              Add Link
            </button>
          )}
        </div>

        {links.length === 0 && (
          <p className="rounded-lg border border-dashed border-zinc-200 p-6 text-center text-sm text-zinc-400 dark:border-zinc-700">
            No custom links yet. Add one to show on your profile.
          </p>
        )}

        <div className="space-y-3">
          {links.map((link, index) => (
            <div
              key={index}
              className="flex items-start gap-2 rounded-lg border border-zinc-200 p-3 dark:border-zinc-700"
            >
              <div className="flex-1 space-y-2">
                <input
                  type="text"
                  placeholder="Link title"
                  value={link.title}
                  onChange={(e) => updateLink(index, "title", e.target.value)}
                  maxLength={50}
                  className="block w-full rounded-md border border-zinc-200 bg-white px-2.5 py-1.5 text-sm focus:border-zinc-400 focus:outline-none dark:border-zinc-600 dark:bg-zinc-800"
                />
                <input
                  type="url"
                  placeholder="https://..."
                  value={link.url}
                  onChange={(e) => updateLink(index, "url", e.target.value)}
                  className="block w-full rounded-md border border-zinc-200 bg-white px-2.5 py-1.5 text-sm focus:border-zinc-400 focus:outline-none dark:border-zinc-600 dark:bg-zinc-800"
                />
              </div>
              <button
                type="button"
                onClick={() => removeLink(index)}
                className="mt-1 rounded p-1 text-zinc-400 hover:text-red-500"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            </div>
          ))}
        </div>
      </div>

      {error && <p className="text-sm text-red-500">{error}</p>}

      <button
        type="submit"
        disabled={loading}
        className="rounded-lg bg-zinc-900 px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-200"
      >
        {loading ? "Saving..." : "Save Changes"}
      </button>
    </form>
  );
}
