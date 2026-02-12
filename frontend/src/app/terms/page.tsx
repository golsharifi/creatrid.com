"use client";

export default function TermsPage() {
  return (
    <div className="mx-auto max-w-3xl px-4 py-16 sm:px-6">
      <h1 className="text-3xl font-bold tracking-tight">Terms of Service</h1>
      <p className="mt-2 text-sm text-zinc-500">Last updated: February 2026</p>

      <div className="mt-10 space-y-8 text-sm leading-relaxed text-zinc-600 dark:text-zinc-400">
        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            1. Acceptance of Terms
          </h2>
          <p className="mt-2">
            By accessing or using Creatrid (&quot;the Service&quot;), you agree to
            be bound by these Terms of Service. If you do not agree, do not use
            the Service.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            2. Description of Service
          </h2>
          <p className="mt-2">
            Creatrid is a verified digital identity platform for content creators.
            The Service allows you to connect social media accounts, build a
            Creator Score, and share a public profile with brands and
            collaborators.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            3. Account Registration
          </h2>
          <p className="mt-2">
            You must sign in using a valid Google account. You are responsible for
            maintaining the security of your account and all activities under it.
            You must provide accurate information and keep your profile up to date.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            4. User Conduct
          </h2>
          <p className="mt-2">You agree not to:</p>
          <ul className="mt-2 list-disc space-y-1 pl-5">
            <li>Impersonate another person or misrepresent your identity</li>
            <li>Use the Service for any unlawful or fraudulent purpose</li>
            <li>Attempt to gain unauthorized access to other accounts or systems</li>
            <li>Upload malicious content or interfere with the Service</li>
            <li>Violate any applicable laws or regulations</li>
          </ul>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            5. Social Account Connections
          </h2>
          <p className="mt-2">
            When you connect a social account, we request read-only access to your
            public profile information (username, display name, avatar, follower
            count). We do not post on your behalf or modify your social accounts.
            You may disconnect any connected account at any time.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            6. Creator Score
          </h2>
          <p className="mt-2">
            Your Creator Score is calculated based on profile completeness, email
            verification, connected platforms, and follower counts. The score is
            for informational purposes and does not constitute an endorsement or
            guarantee of any kind.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            7. Intellectual Property
          </h2>
          <p className="mt-2">
            You retain ownership of all content you provide. By using the Service,
            you grant Creatrid a limited license to display your public profile
            information as part of the platform. The Creatrid name, logo, and
            Service design are our intellectual property.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            8. Account Termination
          </h2>
          <p className="mt-2">
            You may delete your account at any time from the Settings page. Upon
            deletion, all your data including connections, analytics, and
            collaboration history will be permanently removed. We reserve the right
            to suspend or terminate accounts that violate these terms.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            9. Limitation of Liability
          </h2>
          <p className="mt-2">
            The Service is provided &quot;as is&quot; without warranties of any kind.
            Creatrid shall not be liable for any indirect, incidental, or
            consequential damages arising from your use of the Service.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            10. Changes to Terms
          </h2>
          <p className="mt-2">
            We may update these terms from time to time. Continued use of the
            Service after changes constitutes acceptance of the updated terms.
          </p>
        </section>

        <section>
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
            11. Contact
          </h2>
          <p className="mt-2">
            For questions about these terms, contact us at{" "}
            <a
              href="mailto:support@creatrid.com"
              className="text-zinc-900 underline dark:text-zinc-100"
            >
              support@creatrid.com
            </a>
            .
          </p>
        </section>
      </div>
    </div>
  );
}
