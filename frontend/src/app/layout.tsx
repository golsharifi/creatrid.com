import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import { Providers } from "@/components/providers";
import { Header } from "@/components/layout/header";
import { Footer } from "@/components/layout/footer";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  metadataBase: new URL("https://creatrid.com"),
  title: {
    default: "Creatrid — Verified Creator Identity",
    template: "%s | Creatrid",
  },
  description:
    "Your verified digital identity as a creator. Connect social accounts, build a Creator Score, and share one link with brands and collaborators.",
  openGraph: {
    title: "Creatrid — Verified Creator Identity",
    description:
      "Connect your social accounts, build a Creator Score, and share a verified public profile. One link for every platform.",
    url: "https://creatrid.com",
    siteName: "Creatrid",
    type: "website",
    images: [{ url: "/og-image.svg", width: 1200, height: 630, alt: "Creatrid — Verified Creator Identity" }],
  },
  twitter: {
    card: "summary_large_image",
    title: "Creatrid — Verified Creator Identity",
    description:
      "Connect your social accounts, build a Creator Score, and share a verified public profile.",
    images: ["/og-image.svg"],
  },
  robots: { index: true, follow: true },
  alternates: { canonical: "https://creatrid.com" },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <link rel="manifest" href="/manifest.json" />
        <meta name="theme-color" content="#09090b" />
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
        <link rel="apple-touch-icon" href="/icons/icon.svg" />
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{
            __html: JSON.stringify({
              "@context": "https://schema.org",
              "@type": "WebApplication",
              name: "Creatrid",
              url: "https://creatrid.com",
              description: "Verified digital identity platform for creators. Connect social accounts, build a Creator Score, and share one link.",
              applicationCategory: "SocialNetworkingApplication",
              operatingSystem: "Web",
            }),
          }}
        />
      </head>
      <body
        className={`${geistSans.variable} ${geistMono.variable} font-sans antialiased`}
      >
        <Providers>
          <div className="flex min-h-screen flex-col">
            <Header />
            <main className="flex-1">{children}</main>
            <Footer />
          </div>
        </Providers>
      </body>
    </html>
  );
}
