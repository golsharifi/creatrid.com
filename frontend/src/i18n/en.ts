const en = {
  // Common
  common: {
    loading: "Loading...",
    save: "Save",
    cancel: "Cancel",
    previous: "Previous",
    next: "Next",
    creator: "Creator",
    noData: "\u2014",
    page: "Page {{current}} of {{total}}",
    creatorsFound: "{{count}} creators found",
  },

  // Header / Navigation
  header: {
    signIn: "Sign In",
    dashboard: "Dashboard",
    connections: "Connections",
    discover: "Discover",
    analytics: "Analytics",
    settings: "Settings",
    logout: "Logout",
    toggleTheme: "Toggle theme",
    openMenu: "Open menu",
    closeMenu: "Close menu",
    language: "Language",
  },

  // Footer
  footer: {
    copyright: "\u00a9 {{year}} Creatrid",
    terms: "Terms",
    privacy: "Privacy",
  },

  // Sign In
  signIn: {
    title: "Sign in to Creatrid",
    subtitle: "Get your verified creator identity",
    continueWithGoogle: "Continue with Google",
  },

  // Onboarding
  onboarding: {
    title: "Welcome to Creatrid",
    subtitle: "Choose a username for your public profile",
    displayName: "Display Name",
    username: "Username",
    usernamePrefix: "creatrid.com/",
    usernamePlaceholder: "yourname",
    usernameHint: "Letters, numbers, hyphens, and underscores only. 3-30 characters.",
    settingUp: "Setting up...",
    createPassport: "Create My Passport",
  },

  // Dashboard
  dashboard: {
    welcomeBack: "Welcome back, {{name}}",
    managePassport: "Manage your Creator Passport from here.",
    profile: "Profile",
    profileComplete: "{{complete}}/{{total}} complete",
    connections: "Connections",
    connectSocial: "Connect your social accounts",
    creatorScore: "Creator Score",
    scoreBasedOn: "Based on profile, connections & reach",
    scoreUnlock: "Complete your profile to unlock",
    outOf100: "/ 100",
    profileViews: "Profile Views",
    todayAndWeek: "{{today}} today \u00b7 {{week}} this week",
    linkClicks: "Link Clicks",
    noClicks: "No clicks yet",
    engagement: "Engagement",
    clickThroughRate: "Click-through rate",
    yourPublicProfile: "Your Public Profile",
    shareLink: "Share this link with brands and collaborators.",
    preview: "Preview",
    qrCode: "QR Code",
    scanQR: "Scan this QR code to open your Creator Passport.",
    qrPerfectFor: "Perfect for business cards, portfolios, and events.",
    verifyEmailTitle: "Verify your email",
    verifyEmailDesc: "Verify your email to earn 10 bonus points on your Creator Score.",
    verifyEmail: "Send Verification Email",
    sending: "Sending...",
    verificationSent: "Verification email sent!",
    emailVerifiedSuccess: "Your email has been verified! Your Creator Score has been updated.",
    getWidget: "Embeddable Widget",
    getWidgetDesc: "Add your Creator Passport badge to your website.",
  },

  // Settings
  settings: {
    title: "Profile Settings",
    subtitle: "Manage your Creator Passport profile.",
    profileSection: "Profile",
    uploading: "Uploading...",
    profilePhoto: "Profile Photo",
    photoHint: "JPEG, PNG, or WebP. Max 5 MB.",
    imageErrorType: "Only JPEG, PNG, and WebP images are allowed",
    imageErrorSize: "Image must be under 5 MB",
    displayName: "Display Name",
    username: "Username",
    usernamePrefix: "creatrid.com/",
    bio: "Bio",
    bioPlaceholder: "Tell the world about yourself...",
    bioCharCount: "{{count}}/500 characters",
    themeSection: "Profile Theme",
    themeDescription: "Choose a color theme for your public profile.",
    customLinksSection: "Custom Links",
    customLinksDescription: "Add links to your portfolio, website, or socials.",
    addLink: "Add Link",
    noLinksYet: "No custom links yet. Add one to show on your profile.",
    linkTitle: "Link title",
    linkUrl: "https://...",
    emailSection: "Email Notifications",
    emailDescription: "Choose which emails you'd like to receive.",
    emailWelcome: "Welcome email",
    emailWelcomeDesc: "Sent when you first create your account",
    emailConnection: "Connection alerts",
    emailConnectionDesc: "Sent when you connect a new social platform",
    emailDigest: "Weekly digest",
    emailDigestDesc: "Weekly summary of your profile views and clicks",
    emailCollab: "Collaboration requests",
    emailCollabDesc: "Notifications about collaboration requests",
    exportSection: "Export Data",
    exportDescription: "Download a copy of your profile data in JSON format.",
    exporting: "Exporting...",
    exportButton: "Export My Data",
    saving: "Saving...",
    saveChanges: "Save Changes",
    dangerZone: "Danger Zone",
    dangerDescription: "Once you delete your account, all your data will be permanently removed. This action cannot be undone.",
    deleteAccount: "Delete Account",
    deleting: "Deleting...",
    confirmDelete: "Yes, delete my account",
    themeDefault: "Default",
    themeOcean: "Ocean",
    themeSunset: "Sunset",
    themeForest: "Forest",
    themeMidnight: "Midnight",
    themeRose: "Rose",
  },

  // Connections
  connections: {
    title: "Connected Accounts",
    subtitle: "Connect your social accounts to build your Creator Score.",
    successConnected: "Successfully connected {{platform}}!",
    connectionFailed: "Connection failed: {{error}}",
    followers: "{{val}} followers",
    notConnected: "Not connected",
    connect: "Connect",
    disconnect: "Disconnect",
    disconnecting: "...",
    comingSoon: "Coming Soon",
  },

  // Discover
  discover: {
    title: "Discover Creators",
    subtitle: "Find creators to collaborate with based on their score, platforms, and reach.",
    searchPlaceholder: "Search creators...",
    allPlatforms: "All Platforms",
    anyScore: "Any Score",
    loadingCreators: "Loading creators...",
    noCreatorsFound: "No creators found matching your filters.",
    score: "Score: {{score}}",
    connectionsCount: "{{count}} connections",
    requestSent: "Request sent",
    sending: "Sending...",
    collaborate: "Collaborate",
  },

  // Profile (Public)
  profile: {
    loadingProfile: "Loading profile...",
    userNotFound: "User not found",
    profileNotExist: "This creator profile doesn't exist.",
    creatorScore: "Creator Score: {{score}}",
    connectedPlatforms: "Connected Platforms",
    latestVideo: "Latest Video",
    topRepositories: "Top Repositories",
    links: "Links",
    scanToView: "Scan to view this profile",
    verifiedOnCreatrid: "Verified on Creatrid",
  },

  // Collaborations
  collaborations: {
    title: "Collaborations",
    subtitle: "Manage collaboration requests.",
    discoverCreators: "Discover Creators",
    inbox: "Inbox",
    sent: "Sent",
    noRequestsYet: "No collaboration requests yet.",
    noSentRequests: "No sent requests.",
    discoverToCollab: "Discover creators",
    discoverToCollabSuffix: " to collaborate with.",
    accept: "Accept",
    decline: "Decline",
    pending: "Pending",
    accepted: "Accepted",
    declined: "Declined",
  },

  // Admin
  admin: {
    title: "Admin Dashboard",
    subtitle: "Platform overview and user management.",
    totalUsers: "Total Users",
    onboarded: "{{count}} onboarded",
    verifiedUsers: "Verified Users",
    totalConnections: "Total Connections",
    profileViews: "Profile Views",
    linkClicks: "Link Clicks",
    avgScore: "Avg Score",
    acrossAllUsers: "Across all users",
    usersCount: "Users ({{count}})",
    tableUser: "User",
    tableUsername: "Username",
    tableScore: "Score",
    tableConnections: "Connections",
    tableStatus: "Status",
    tableActions: "Actions",
    verified: "Verified",
    notOnboarded: "Not onboarded",
    verify: "Verify",
    unverify: "Unverify",
  },

  // Landing Page
  landing: {
    // Hero
    heroBadge: "Verified Creator Identity",
    heroHeadline: "Your Creator Passport.",
    heroHeadlineAccent: "One Link. Every Platform.",
    heroSubtext:
      "Build a verified digital identity that proves who you are across platforms. Connect your accounts, earn a Creator Score, and share one link with brands and collaborators.",
    ctaDashboard: "Go to Dashboard",
    ctaGetPassport: "Get Your Creator Passport",
    ctaHowItWorks: "See how it works",

    // Social Proof Bar
    socialProofLabel: "Trusted by creators worldwide",
    statCreators: "1,000+",
    statCreatorsLabel: "Creators",
    statPlatforms: "7",
    statPlatformsLabel: "Platforms",
    statConnections: "10,000+",
    statConnectionsLabel: "Connections",

    // How it Works
    howTitle: "How it Works",
    howSubtitle:
      "Three simple steps to a verified creator identity that brands and collaborators can trust.",
    step1Title: "Connect your accounts",
    step1Desc:
      "Link your YouTube, GitHub, Twitter, LinkedIn, Instagram, Dribbble, and Behance profiles with one-click OAuth.",
    step2Title: "Build your Creator Score",
    step2Desc:
      "Our scoring engine evaluates your profile completeness, verified connections, and audience reach on a 0\u2013100 scale.",
    step3Title: "Share your verified profile",
    step3Desc:
      "Get a single public link and QR code that showcases your verified identity, connections, and reputation.",

    // Features Grid
    featuresTitle: "Everything you need to prove who you are",
    featuresSubtitle:
      "A complete toolkit for building, managing, and sharing your verified creator identity.",
    feature1Title: "Verified Identity",
    feature1Desc:
      "Connect your social accounts and prove you are who you say you are. No fakes, no impersonators.",
    feature2Title: "Creator Score",
    feature2Desc:
      "A 0\u2013100 reputation score based on your profile, verified connections, and audience reach across platforms.",
    feature3Title: "Cross-Platform Profiles",
    feature3Desc:
      "One beautiful profile that brings together all your platforms, stats, and content in a single link.",
    feature4Title: "Brand Discovery",
    feature4Desc:
      "Get discovered by brands and agencies looking for verified creators to collaborate with.",
    feature5Title: "Real-time Analytics",
    feature5Desc:
      "Track profile views, link clicks, and engagement in real time so you know who is paying attention.",
    feature6Title: "Privacy First",
    feature6Desc:
      "You control what is shared. Choose which platforms and information appear on your public profile.",

    // Platform Logos
    platformsTitle: "Connect all your platforms",
    platformsSubtitle:
      "Creatrid supports seven major platforms today, with more on the way.",

    // FAQ
    faqTitle: "Frequently Asked Questions",
    faqSubtitle: "Everything you need to know about Creatrid.",
    faq1Q: "What is Creatrid?",
    faq1A:
      "Creatrid is a verified digital identity platform for creators. You connect your social accounts, build a Creator Score, and share a single public profile link that proves you are the real deal.",
    faq2Q: "How is the Creator Score calculated?",
    faq2A:
      "Your Creator Score (0\u2013100) is based on four factors: profile completeness (20 points), verified email (10 points), number of connected platforms (up to 50 points), and a logarithmic bonus for audience reach (up to 20 points).",
    faq3Q: "Is it free?",
    faq3A:
      "Yes. Creatrid is completely free for creators. Sign in with Google, connect your accounts, and start building your verified profile at no cost.",
    faq4Q: "Can brands verify creators?",
    faq4A:
      "Absolutely. Brands and agencies can view any creator\u2019s public profile to see verified connections, Creator Score, and linked content \u2014 no login required.",

    // Final CTA
    ctaFinalHeadline: "Ready to prove who you are?",
    ctaFinalSubtext:
      "Join thousands of creators who use Creatrid to build trust, get discovered, and land collaborations.",
  },

  // API Keys
  apiKeys: {
    title: "API Keys",
    subtitle: "Manage API keys for programmatic access to the verification API.",
    createKey: "Create API Key",
    creating: "Creating...",
    keyName: "Key Name",
    keyNamePlaceholder: "e.g. My Brand Integration",
    noKeysYet: "No API keys yet. Create one to start using the verification API.",
    prefix: "Prefix",
    lastUsed: "Last used",
    never: "Never",
    created: "Created",
    revoke: "Revoke",
    revoking: "...",
    newKeyTitle: "Your new API key",
    newKeyWarning: "Save this key now. It will not be shown again.",
    copied: "Copied!",
    copy: "Copy",
    docsTitle: "API Documentation",
    docsBaseUrl: "Base URL",
    docsAuth: "Authentication",
    docsAuthDesc: "Include your API key in the Authorization header:",
    docsEndpoints: "Endpoints",
    docsVerifyTitle: "Full verification",
    docsVerifyDesc: "Get complete creator profile, connections, and score.",
    docsScoreTitle: "Score only",
    docsScoreDesc: "Quick check for a creator's score and verification status.",
    docsSearchTitle: "Search creators",
    docsSearchDesc: "Search and filter creators by name, score, or platform.",
    docsRateLimit: "Rate limit: 100 requests per minute per key.",
  },

  // Analytics
  analytics: {
    title: "Analytics",
    subtitle: "Track your profile performance and audience engagement.",
    totalViews: "Total Views",
    todayViews: "{{count}} today",
    totalClicks: "Total Clicks",
    clickThrough: "Click-through Rate",
    weekViews: "Views This Week",
    viewsOverTime: "Views Over Time",
    clicksOverTime: "Clicks Over Time",
    topReferrers: "Top Referrers",
    viewsByHour: "Views by Hour of Day",
    clicksByType: "Clicks by Type",
    noData: "No analytics data yet. Share your profile to start tracking.",
    last30Days: "Last 30 days",
  },

  // Pricing
  pricing: {
    title: "Pricing",
    subtitle: "Choose the plan that fits your needs.",
    free: "Free",
    pro: "Pro",
    business: "Business",
    month: "/month",
    freePrice: "$0",
    proPrice: "$10",
    businessPrice: "$50",
    getStarted: "Get Started",
    upgradeToPro: "Upgrade to Pro",
    contactSales: "Contact Sales",
    currentPlan: "Current Plan",
    features: "Features",
    freeFeatures: [
      "3 social connections",
      "Basic profile",
      "Creator Score",
      "Public profile page",
    ],
    proFeatures: [
      "Unlimited connections",
      "Advanced analytics",
      "Custom themes",
      "API keys (1,000 req/mo)",
      "Embeddable widget",
      "Priority support",
    ],
    businessFeatures: [
      "Everything in Pro",
      "Bulk verification API (10K req/mo)",
      "Brand dashboard",
      "Saved creator lists",
      "White-label embed",
    ],
    mostPopular: "Most Popular",
  },

  // Billing
  billing: {
    title: "Billing",
    subtitle: "Manage your subscription and billing.",
    currentPlan: "Current Plan",
    freePlan: "Free",
    proPlan: "Pro",
    businessPlan: "Business",
    manage: "Manage Subscription",
    upgrade: "Upgrade",
    successTitle: "Subscription activated!",
    successDesc:
      "Thank you for upgrading. Your new features are now available.",
    canceledTitle: "Checkout canceled",
    canceledDesc: "No changes were made to your subscription.",
    planFeatures: "Your plan includes:",
  },

  // Widget
  widget: {
    title: "Embeddable Widget",
    subtitle: "Add your Creator Passport badge to your website or portfolio.",
    preview: "Preview",
    embedCode: "Embed Code",
    htmlEmbed: "HTML Embed",
    markdownBadge: "Markdown Badge",
    directLink: "Direct Link",
    copied: "Copied!",
    copy: "Copy",
  },
};

export default en;
