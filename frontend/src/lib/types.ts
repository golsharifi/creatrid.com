export interface CustomLink {
  title: string;
  url: string;
}

export interface EmailPrefs {
  welcome: boolean;
  connectionAlert: boolean;
  weeklyDigest: boolean;
  collaborations: boolean;
}

export interface User {
  id: string;
  name: string | null;
  email: string;
  emailVerified: string | null;
  image: string | null;
  username: string | null;
  bio: string | null;
  role: "CREATOR" | "BRAND" | "ADMIN";
  creatorScore: number | null;
  creatorTier: string;
  isVerified: boolean;
  totpEnabled: boolean;
  onboarded: boolean;
  theme: string;
  customLinks: CustomLink[];
  emailPrefs: EmailPrefs;
  createdAt: string;
  updatedAt: string;
}

export interface PublicUser {
  id: string;
  name: string | null;
  username: string | null;
  bio: string | null;
  image: string | null;
  role: "CREATOR" | "BRAND" | "ADMIN";
  creatorScore: number | null;
  creatorTier: string;
  isVerified: boolean;
  theme: string;
  customLinks: CustomLink[];
  createdAt: string;
}

export interface Connection {
  platform: string;
  username: string | null;
  displayName: string | null;
  avatarUrl: string | null;
  profileUrl: string | null;
  followerCount: number | null;
  metadata: Record<string, unknown>;
  connectedAt: string;
}
