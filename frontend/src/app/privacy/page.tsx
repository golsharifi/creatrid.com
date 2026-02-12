"use client";

export default function PrivacyPage() {
  return (
    <div className="mx-auto max-w-3xl px-4 py-16 sm:px-6">
      <h1 className="text-3xl font-bold tracking-tight">Privacy Policy</h1>
      <p className="mt-2 text-sm text-zinc-500">Last updated: February 2026</p>

      <div className="mt-10 space-y-8 text-sm leading-relaxed text-zinc-600 dark:text-zinc-400">
        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            1. Information We Collect
          </h2>
          <p className="mt-2">When you use Creatrid, we collect:</p>
          <ul className="mt-2 list-disc space-y-1 pl-5">
            <li>
              <strong>Account information:</strong> Name, email address, and
              profile photo from your Google account when you sign in.
            </li>
            <li>
              <strong>Profile information:</strong> Username, bio, theme
              preference, and custom links you provide.
            </li>
            <li>
              <strong>Social connections:</strong> Public profile data from
              platforms you connect (username, display name, avatar, follower
              count, public metadata). We request read-only access.
            </li>
            <li>
              <strong>Analytics data:</strong> Profile views and link clicks
              (visitor IP addresses are stored for spam prevention and are not
              shared).
            </li>
            <li>
              <strong>Collaboration messages:</strong> Messages sent through the
              collaboration feature.
            </li>
          </ul>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            2. How We Use Your Information
          </h2>
          <ul className="mt-2 list-disc space-y-1 pl-5">
            <li>To provide and operate the Service</li>
            <li>To calculate and display your Creator Score</li>
            <li>To display your public profile to other users and brands</li>
            <li>To send email notifications (which you can opt out of)</li>
            <li>To prevent abuse and enforce our Terms of Service</li>
            <li>To improve the Service</li>
          </ul>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            3. Data Storage and Security
          </h2>
          <p className="mt-2">
            Your data is stored on Azure-hosted PostgreSQL databases in Europe.
            We use HTTPS encryption for all data in transit, httpOnly cookies for
            authentication, and follow industry-standard security practices.
            OAuth tokens for connected platforms are stored encrypted and are
            never exposed to the frontend.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            4. Data Sharing
          </h2>
          <p className="mt-2">
            We do not sell your personal data. Information is shared only in these
            cases:
          </p>
          <ul className="mt-2 list-disc space-y-1 pl-5">
            <li>
              <strong>Public profile:</strong> Your username, name, bio, creator
              score, connected platforms (without tokens), and custom links are
              visible on your public profile.
            </li>
            <li>
              <strong>Service providers:</strong> We use Azure for hosting and
              Google for authentication. These providers process data on our
              behalf under their respective privacy policies.
            </li>
            <li>
              <strong>Legal requirements:</strong> We may disclose data if
              required by law.
            </li>
          </ul>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            5. Your Rights
          </h2>
          <p className="mt-2">You have the right to:</p>
          <ul className="mt-2 list-disc space-y-1 pl-5">
            <li>
              <strong>Access:</strong> View all data we hold about you in your
              Settings page.
            </li>
            <li>
              <strong>Export:</strong> Download your profile data in JSON format.
            </li>
            <li>
              <strong>Rectification:</strong> Update your profile information at
              any time.
            </li>
            <li>
              <strong>Erasure:</strong> Delete your account and all associated
              data permanently from the Settings page.
            </li>
            <li>
              <strong>Disconnect:</strong> Remove any connected social account at
              any time.
            </li>
            <li>
              <strong>Opt out:</strong> Manage email notification preferences in
              Settings.
            </li>
          </ul>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            6. Cookies
          </h2>
          <p className="mt-2">
            We use a single httpOnly authentication cookie to maintain your
            session. We do not use tracking cookies or third-party analytics
            cookies.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            7. Data Retention
          </h2>
          <p className="mt-2">
            We retain your data for as long as your account is active. When you
            delete your account, all data is permanently removed within 30 days.
            Analytics data (profile views, link clicks) is deleted immediately
            upon account deletion.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            8. Children
          </h2>
          <p className="mt-2">
            The Service is not intended for users under 16. We do not knowingly
            collect data from children under 16.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            9. Changes to This Policy
          </h2>
          <p className="mt-2">
            We may update this policy from time to time. We will notify you of
            significant changes via email or a notice on the Service.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            10. Contact
          </h2>
          <p className="mt-2">
            For privacy-related questions, contact us at{" "}
            <a
              href="mailto:privacy@creatrid.com"
              className="text-zinc-900 underline dark:text-zinc-100"
            >
              privacy@creatrid.com
            </a>
            .
          </p>
        </section>
      </div>
    </div>
  );
}
