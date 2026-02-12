const es = {
  // Common
  common: {
    loading: "Cargando...",
    save: "Guardar",
    cancel: "Cancelar",
    previous: "Anterior",
    next: "Siguiente",
    creator: "Creador",
    noData: "\u2014",
    page: "Pagina {{current}} de {{total}}",
    creatorsFound: "{{count}} creadores encontrados",
  },

  // Header / Navigation
  header: {
    signIn: "Iniciar sesion",
    dashboard: "Panel",
    connections: "Conexiones",
    discover: "Descubrir",
    analytics: "Estadisticas",
    settings: "Ajustes",
    logout: "Cerrar sesion",
    toggleTheme: "Cambiar tema",
    openMenu: "Abrir menu",
    closeMenu: "Cerrar menu",
    language: "Idioma",
  },

  // Footer
  footer: {
    copyright: "\u00a9 {{year}} Creatrid",
    terms: "Terminos",
    privacy: "Privacidad",
  },

  // Sign In
  signIn: {
    title: "Inicia sesion en Creatrid",
    subtitle: "Obtiene tu identidad verificada de creador",
    continueWithGoogle: "Continuar con Google",
  },

  // Onboarding
  onboarding: {
    title: "Bienvenido a Creatrid",
    subtitle: "Elige un nombre de usuario para tu perfil publico",
    displayName: "Nombre visible",
    username: "Nombre de usuario",
    usernamePrefix: "creatrid.com/",
    usernamePlaceholder: "tunombre",
    usernameHint: "Solo letras, numeros, guiones y guiones bajos. De 3 a 30 caracteres.",
    settingUp: "Configurando...",
    createPassport: "Crear mi Pasaporte",
  },

  // Dashboard
  dashboard: {
    welcomeBack: "Bienvenido de nuevo, {{name}}",
    managePassport: "Gestiona tu Creator Passport desde aqui.",
    profile: "Perfil",
    profileComplete: "{{complete}}/{{total}} completado",
    connections: "Conexiones",
    connectSocial: "Conecta tus cuentas sociales",
    creatorScore: "Creator Score",
    scoreBasedOn: "Basado en perfil, conexiones y alcance",
    scoreUnlock: "Completa tu perfil para desbloquear",
    outOf100: "/ 100",
    profileViews: "Visitas al perfil",
    todayAndWeek: "{{today}} hoy \u00b7 {{week}} esta semana",
    linkClicks: "Clics en enlaces",
    noClicks: "Sin clics aun",
    engagement: "Interaccion",
    clickThroughRate: "Tasa de clics",
    yourPublicProfile: "Tu perfil publico",
    shareLink: "Comparte este enlace con marcas y colaboradores.",
    preview: "Vista previa",
    qrCode: "Codigo QR",
    scanQR: "Escanea este codigo QR para abrir tu Creator Passport.",
    qrPerfectFor: "Perfecto para tarjetas de visita, portfolios y eventos.",
    verifyEmailTitle: "Verifica tu email",
    verifyEmailDesc: "Verifica tu email para ganar 10 puntos extra en tu Creator Score.",
    verifyEmail: "Enviar email de verificacion",
    sending: "Enviando...",
    verificationSent: "Email de verificacion enviado!",
    emailVerifiedSuccess: "Tu email ha sido verificado! Tu Creator Score ha sido actualizado.",
    getWidget: "Widget integrable",
    getWidgetDesc: "Agrega tu insignia de Creator Passport a tu sitio web.",
  },

  // Settings
  settings: {
    title: "Ajustes del perfil",
    subtitle: "Gestiona tu perfil de Creator Passport.",
    profileSection: "Perfil",
    uploading: "Subiendo...",
    profilePhoto: "Foto de perfil",
    photoHint: "JPEG, PNG o WebP. Maximo 5 MB.",
    imageErrorType: "Solo se permiten imagenes JPEG, PNG y WebP",
    imageErrorSize: "La imagen debe pesar menos de 5 MB",
    displayName: "Nombre visible",
    username: "Nombre de usuario",
    usernamePrefix: "creatrid.com/",
    bio: "Biografia",
    bioPlaceholder: "Cuentale al mundo sobre ti...",
    bioCharCount: "{{count}}/500 caracteres",
    themeSection: "Tema del perfil",
    themeDescription: "Elige un tema de color para tu perfil publico.",
    customLinksSection: "Enlaces personalizados",
    customLinksDescription: "Agrega enlaces a tu portfolio, sitio web o redes sociales.",
    addLink: "Agregar enlace",
    noLinksYet: "Sin enlaces personalizados aun. Agrega uno para mostrarlo en tu perfil.",
    linkTitle: "Titulo del enlace",
    linkUrl: "https://...",
    emailSection: "Notificaciones por email",
    emailDescription: "Elige que emails deseas recibir.",
    emailWelcome: "Email de bienvenida",
    emailWelcomeDesc: "Se envia al crear tu cuenta",
    emailConnection: "Alertas de conexion",
    emailConnectionDesc: "Se envia al conectar una nueva plataforma social",
    emailDigest: "Resumen semanal",
    emailDigestDesc: "Resumen semanal de las visitas y clics de tu perfil",
    emailCollab: "Solicitudes de colaboracion",
    emailCollabDesc: "Notificaciones sobre solicitudes de colaboracion",
    exportSection: "Exportar datos",
    exportDescription: "Descarga una copia de tus datos de perfil en formato JSON.",
    exporting: "Exportando...",
    exportButton: "Exportar mis datos",
    saving: "Guardando...",
    saveChanges: "Guardar cambios",
    dangerZone: "Zona de peligro",
    dangerDescription: "Una vez que elimines tu cuenta, todos tus datos seran eliminados permanentemente. Esta accion no se puede deshacer.",
    deleteAccount: "Eliminar cuenta",
    deleting: "Eliminando...",
    confirmDelete: "Si, eliminar mi cuenta",
    themeDefault: "Predeterminado",
    themeOcean: "Oceano",
    themeSunset: "Atardecer",
    themeForest: "Bosque",
    themeMidnight: "Medianoche",
    themeRose: "Rosa",
  },

  // Connections
  connections: {
    title: "Cuentas conectadas",
    subtitle: "Conecta tus cuentas sociales para aumentar tu Creator Score.",
    successConnected: "{{platform}} conectado correctamente!",
    connectionFailed: "Error de conexion: {{error}}",
    followers: "{{val}} seguidores",
    notConnected: "No conectado",
    connect: "Conectar",
    disconnect: "Desconectar",
    disconnecting: "...",
    comingSoon: "Proximamente",
  },

  // Discover
  discover: {
    title: "Descubrir creadores",
    subtitle: "Encuentra creadores con quienes colaborar segun su puntuacion, plataformas y alcance.",
    searchPlaceholder: "Buscar creadores...",
    allPlatforms: "Todas las plataformas",
    anyScore: "Cualquier puntuacion",
    loadingCreators: "Cargando creadores...",
    noCreatorsFound: "No se encontraron creadores con esos filtros.",
    score: "Puntuacion: {{score}}",
    connectionsCount: "{{count}} conexiones",
    requestSent: "Solicitud enviada",
    sending: "Enviando...",
    collaborate: "Colaborar",
  },

  // Profile (Public)
  profile: {
    loadingProfile: "Cargando perfil...",
    userNotFound: "Usuario no encontrado",
    profileNotExist: "Este perfil de creador no existe.",
    creatorScore: "Creator Score: {{score}}",
    connectedPlatforms: "Plataformas conectadas",
    latestVideo: "Ultimo video",
    topRepositories: "Repositorios destacados",
    links: "Enlaces",
    scanToView: "Escanea para ver este perfil",
    verifiedOnCreatrid: "Verificado en Creatrid",
  },

  // Collaborations
  collaborations: {
    title: "Colaboraciones",
    subtitle: "Gestiona tus solicitudes de colaboracion.",
    discoverCreators: "Descubrir creadores",
    inbox: "Bandeja de entrada",
    sent: "Enviados",
    noRequestsYet: "Aun no hay solicitudes de colaboracion.",
    noSentRequests: "No hay solicitudes enviadas.",
    discoverToCollab: "Descubre creadores",
    discoverToCollabSuffix: " con quienes colaborar.",
    accept: "Aceptar",
    decline: "Rechazar",
    pending: "Pendiente",
    accepted: "Aceptada",
    declined: "Rechazada",
  },

  // Admin
  admin: {
    title: "Panel de administracion",
    subtitle: "Vista general de la plataforma y gestion de usuarios.",
    totalUsers: "Total de usuarios",
    onboarded: "{{count}} registrados",
    verifiedUsers: "Usuarios verificados",
    totalConnections: "Total de conexiones",
    profileViews: "Visitas al perfil",
    linkClicks: "Clics en enlaces",
    avgScore: "Puntuacion media",
    acrossAllUsers: "De todos los usuarios",
    usersCount: "Usuarios ({{count}})",
    tableUser: "Usuario",
    tableUsername: "Nombre de usuario",
    tableScore: "Puntuacion",
    tableConnections: "Conexiones",
    tableStatus: "Estado",
    tableActions: "Acciones",
    verified: "Verificado",
    notOnboarded: "No registrado",
    verify: "Verificar",
    unverify: "Quitar verificacion",
  },

  // Landing Page
  landing: {
    // Hero
    heroBadge: "Identidad de creador verificada",
    heroHeadline: "Tu Creator Passport.",
    heroHeadlineAccent: "Un enlace. Todas las plataformas.",
    heroSubtext:
      "Construye una identidad digital verificada que demuestre quien eres en todas las plataformas. Conecta tus cuentas, gana un Creator Score y comparte un solo enlace con marcas y colaboradores.",
    ctaDashboard: "Ir al panel",
    ctaGetPassport: "Obtiene tu Creator Passport",
    ctaHowItWorks: "Descubre como funciona",

    // Social Proof Bar
    socialProofLabel: "La confianza de creadores en todo el mundo",
    statCreators: "1.000+",
    statCreatorsLabel: "Creadores",
    statPlatforms: "7",
    statPlatformsLabel: "Plataformas",
    statConnections: "10.000+",
    statConnectionsLabel: "Conexiones",

    // How it Works
    howTitle: "Como funciona",
    howSubtitle:
      "Tres sencillos pasos para una identidad de creador verificada en la que marcas y colaboradores pueden confiar.",
    step1Title: "Conecta tus cuentas",
    step1Desc:
      "Vincula tus perfiles de YouTube, GitHub, Twitter, LinkedIn, Instagram, Dribbble y Behance con un solo clic mediante OAuth.",
    step2Title: "Construye tu Creator Score",
    step2Desc:
      "Nuestro motor de puntuacion evalua la completitud de tu perfil, conexiones verificadas y alcance de audiencia en una escala de 0 a 100.",
    step3Title: "Comparte tu perfil verificado",
    step3Desc:
      "Obtiene un enlace publico y un codigo QR que muestra tu identidad verificada, conexiones y reputacion.",

    // Features Grid
    featuresTitle: "Todo lo que necesitas para demostrar quien eres",
    featuresSubtitle:
      "Un conjunto completo de herramientas para construir, gestionar y compartir tu identidad verificada de creador.",
    feature1Title: "Identidad verificada",
    feature1Desc:
      "Conecta tus cuentas sociales y demuestra que eres quien dices ser. Sin falsificaciones ni suplantadores.",
    feature2Title: "Creator Score",
    feature2Desc:
      "Una puntuacion de reputacion de 0 a 100 basada en tu perfil, conexiones verificadas y alcance de audiencia en todas las plataformas.",
    feature3Title: "Perfiles multiplataforma",
    feature3Desc:
      "Un perfil atractivo que reune todas tus plataformas, estadisticas y contenido en un solo enlace.",
    feature4Title: "Descubrimiento por marcas",
    feature4Desc:
      "Se descubierto por marcas y agencias que buscan creadores verificados con quienes colaborar.",
    feature5Title: "Estadisticas en tiempo real",
    feature5Desc:
      "Rastrea visitas al perfil, clics en enlaces e interacciones en tiempo real para saber quien esta prestando atencion.",
    feature6Title: "Privacidad ante todo",
    feature6Desc:
      "Tu controlas lo que se comparte. Elige que plataformas e informacion aparecen en tu perfil publico.",

    // Platform Logos
    platformsTitle: "Conecta todas tus plataformas",
    platformsSubtitle:
      "Creatrid soporta siete plataformas principales hoy, con mas en camino.",

    // FAQ
    faqTitle: "Preguntas frecuentes",
    faqSubtitle: "Todo lo que necesitas saber sobre Creatrid.",
    faq1Q: "Que es Creatrid?",
    faq1A:
      "Creatrid es una plataforma de identidad digital verificada para creadores. Conectas tus cuentas sociales, construyes un Creator Score y compartes un solo enlace de perfil publico que demuestra que eres autentico.",
    faq2Q: "Como se calcula el Creator Score?",
    faq2A:
      "Tu Creator Score (0\u2013100) se basa en cuatro factores: completitud del perfil (20 puntos), email verificado (10 puntos), numero de plataformas conectadas (hasta 50 puntos) y un bonus logaritmico por alcance de audiencia (hasta 20 puntos).",
    faq3Q: "Es gratis?",
    faq3A:
      "Si. Creatrid es completamente gratis para creadores. Inicia sesion con Google, conecta tus cuentas y empieza a construir tu perfil verificado sin costo alguno.",
    faq4Q: "Pueden las marcas verificar creadores?",
    faq4A:
      "Por supuesto. Las marcas y agencias pueden ver el perfil publico de cualquier creador para comprobar conexiones verificadas, Creator Score y contenido vinculado, sin necesidad de iniciar sesion.",

    // Final CTA
    ctaFinalHeadline: "Listo para demostrar quien eres?",
    ctaFinalSubtext:
      "Unete a miles de creadores que usan Creatrid para generar confianza, ser descubiertos y conseguir colaboraciones.",
  },

  // API Keys
  apiKeys: {
    title: "Claves API",
    subtitle: "Gestiona las claves API para acceso programatico a la API de verificacion.",
    createKey: "Crear clave API",
    creating: "Creando...",
    keyName: "Nombre de la clave",
    keyNamePlaceholder: "p. ej. Mi integracion de marca",
    noKeysYet: "Aun no hay claves API. Crea una para empezar a usar la API de verificacion.",
    prefix: "Prefijo",
    lastUsed: "Ultimo uso",
    never: "Nunca",
    created: "Creada",
    revoke: "Revocar",
    revoking: "...",
    newKeyTitle: "Tu nueva clave API",
    newKeyWarning: "Guarda esta clave ahora. No se mostrara de nuevo.",
    copied: "Copiado!",
    copy: "Copiar",
    docsTitle: "Documentacion de la API",
    docsBaseUrl: "URL base",
    docsAuth: "Autenticacion",
    docsAuthDesc: "Incluye tu clave API en la cabecera Authorization:",
    docsEndpoints: "Endpoints",
    docsVerifyTitle: "Verificacion completa",
    docsVerifyDesc: "Obtiene el perfil completo del creador, conexiones y puntuacion.",
    docsScoreTitle: "Solo puntuacion",
    docsScoreDesc: "Consulta rapida de la puntuacion y estado de verificacion de un creador.",
    docsSearchTitle: "Buscar creadores",
    docsSearchDesc: "Busca y filtra creadores por nombre, puntuacion o plataforma.",
    docsRateLimit: "Limite de velocidad: 100 solicitudes por minuto por clave.",
  },

  // Analytics
  analytics: {
    title: "Estadisticas",
    subtitle: "Rastrea el rendimiento de tu perfil y la interaccion de tu audiencia.",
    totalViews: "Visitas totales",
    todayViews: "{{count}} hoy",
    totalClicks: "Clics totales",
    clickThrough: "Tasa de clics",
    weekViews: "Visitas esta semana",
    viewsOverTime: "Visitas a lo largo del tiempo",
    clicksOverTime: "Clics a lo largo del tiempo",
    topReferrers: "Principales referentes",
    viewsByHour: "Visitas por hora del dia",
    clicksByType: "Clics por tipo",
    noData: "Aun no hay datos de estadisticas. Comparte tu perfil para empezar a rastrear.",
    last30Days: "Ultimos 30 dias",
  },

  // Pricing
  pricing: {
    title: "Precios",
    subtitle: "Elige el plan que se adapte a tus necesidades.",
    free: "Gratis",
    pro: "Pro",
    business: "Business",
    month: "/mes",
    freePrice: "$0",
    proPrice: "$10",
    businessPrice: "$50",
    getStarted: "Comenzar",
    upgradeToPro: "Mejorar a Pro",
    contactSales: "Contactar ventas",
    currentPlan: "Plan actual",
    features: "Caracteristicas",
    freeFeatures: [
      "3 conexiones sociales",
      "Perfil basico",
      "Creator Score",
      "Pagina de perfil publico",
    ],
    proFeatures: [
      "Conexiones ilimitadas",
      "Estadisticas avanzadas",
      "Temas personalizados",
      "Claves API (1.000 sol./mes)",
      "Widget integrable",
      "Soporte prioritario",
    ],
    businessFeatures: [
      "Todo lo de Pro",
      "API de verificacion masiva (10K sol./mes)",
      "Panel de marca",
      "Listas de creadores guardadas",
      "Widget de marca blanca",
    ],
    mostPopular: "Mas popular",
  },

  // Billing
  billing: {
    title: "Facturacion",
    subtitle: "Gestiona tu suscripcion y facturacion.",
    currentPlan: "Plan actual",
    freePlan: "Gratis",
    proPlan: "Pro",
    businessPlan: "Business",
    manage: "Gestionar suscripcion",
    upgrade: "Mejorar",
    successTitle: "Suscripcion activada!",
    successDesc:
      "Gracias por mejorar tu plan. Tus nuevas funciones ya estan disponibles.",
    canceledTitle: "Pago cancelado",
    canceledDesc: "No se realizaron cambios en tu suscripcion.",
    planFeatures: "Tu plan incluye:",
  },

  // Widget
  widget: {
    title: "Widget integrable",
    subtitle: "Agrega tu insignia de Creator Passport a tu sitio web o portfolio.",
    preview: "Vista previa",
    embedCode: "Codigo de insercion",
    htmlEmbed: "Insercion HTML",
    markdownBadge: "Insignia Markdown",
    directLink: "Enlace directo",
    copied: "Copiado!",
    copy: "Copiar",
  },
};

export default es;
