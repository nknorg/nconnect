import colors from 'vuetify/es5/util/colors'

export default {
  // Target: https://go.nuxtjs.dev/config-target
  target: 'static',

  // Global page headers: https://go.nuxtjs.dev/config-head
  head: {
    titleTemplate: '%s - nconnect-web',
    title: 'nconnect-web',
    meta: [
      {charset: 'utf-8'},
      {name: 'viewport', content: 'width=device-width, initial-scale=1'},
      {hid: 'description', name: 'description', content: ''},
      {name: 'format-detection', content: 'telephone=no'}
    ],
    link: [
      {rel: 'icon', type: 'image/x-icon', href: '/favicon.ico'}
    ],
    script: [

    ]
  },

  // Global CSS: https://go.nuxtjs.dev/config-css
  css: [],

  // Plugins to run before rendering page: https://go.nuxtjs.dev/config-plugins
  plugins: [
    '~/plugins/i18n'
  ],

  // Auto import components: https://go.nuxtjs.dev/config-components
  components: true,

  // Modules for dev and build (recommended): https://go.nuxtjs.dev/config-modules
  buildModules: [
    // https://go.nuxtjs.dev/vuetify
    '@nuxtjs/vuetify',
  ],

  // Modules: https://go.nuxtjs.dev/config-modules
  modules: [
    // https://go.nuxtjs.dev/axios
    '@nuxtjs/axios',
    // https://go.nuxtjs.dev/pwa
    '@nuxtjs/pwa',
  ],
  i18n: {
    // Options
    // vue-i18n configuration
    vueI18n: {
      fallbackLocale: 'en',
    },
    parsePages: false,
    // If true, vue-i18n-loader is added to Nuxt's Webpack config
    vueI18nLoader: false,

    // List of locales supported by your app
    // This can either be an array of codes: ['en', 'fr', 'es']
    // Or an array of objects for more complex configurations:
    // [
    //   { code: 'en', iso: 'en-US', file: 'en.js' },
    //   { code: 'fr', iso: 'fr-FR', file: 'fr.js' },
    //   { code: 'es', iso: 'es-ES', file: 'es.js' }
    // ]
    locales: [
      {code: 'en', name: 'English', iso: 'en-US', file: 'en.json'},
      {code: 'zh', name: '简体中文', iso: 'zh', file: 'zh-CN.json'},
      {code: 'zh-TW', name: '繁体中文', iso: 'zh-Hant', file: 'zh-TW.json'},
    ],

    // The app's default locale, URLs for this locale won't have a prefix if
    // strategy is prefix_except_default
    defaultLocale: 'en',

    // Separator used to generated routes name for each locale, you shouldn't
    // need to change this
    // routesNameSeparator: '___',

    // Suffix added to generated routes name for default locale if strategy is prefix_and_default,
    // you shouldn't need to change this
    // defaultLocaleRouteNameSuffix: 'default',

    // Routes generation strategy, can be set to one of the following:
    // - 'prefix_except_default': add locale prefix for every locale except default
    // - 'prefix': add locale prefix for every locale
    // - 'prefix_and_default': add locale prefix for every locale and default
    // strategy: 'prefix_except_default',

    // Wether or not the translations should be lazy-loaded, if this is enabled,
    // you MUST configure langDir option, and locales must be an array of objects,
    // each containing a file key
    lazy: true,

    // Directory that contains translations files when lazy-loading messages,
    // this CAN NOT be empty if lazy-loading is enabled
    langDir: 'locales/',

    // Set this to a path to which you want to redirect users accessing root URL (/)
    // rootRedirect: null,

    // Enable browser language detection to automatically redirect user
    // to their preferred language as they visit your app for the first time
    // Set to false to disable
    detectBrowserLanguage: {
      // If enabled, a cookie is set once a user has been redirected to his
      // preferred language to prevent subsequent redirections
      // Set to false to redirect every time
      useCookie: true,
      // Cookie name
      cookieKey: 'i18n_redirected',
      // Set to always redirect to value stored in the cookie, not just once
      alwaysRedirect: true,
      // If no locale for the browsers locale is a match, use this one as a fallback
      fallbackLocale: 'en'
    },

    // Set this to true if you're using different domains for each language
    // If enabled, no prefix is added to your routes and you MUST configure locales
    // as an array of objects, each containing a domain key
    // differentDomains: false,

    // If using different domains, set this to true to get hostname from X-Forwared-Host
    // HTTP header instead of window.location
    // forwardedHost: false,

    // If true, SEO metadata is generated for routes that have i18n enabled
    // Set to false to disable app-wide
    // seo: true,

    // Base URL to use as prefix for alternate URLs in hreflang tags
    // baseUrl: '',

    // By default a store module is registered and kept in sync with the
    // app's i18n current state
    // Set to false to disable
    // vuex: {
    // Module namespace
    // moduleName: 'i18n',

    // Mutations config
    // mutations: {
    // Mutation to commit to store current locale, set to false to disable
    // setLocale: 'I18N_SET_LOCALE',

    // Mutation to commit to store current message, set to false to disable
    // setMessages: 'I18N_SET_MESSAGES'
    // },

    // PreserveState from server
    // preserveState: false
    // },

    // By default, custom routes are extracted from page files using acorn parsing,
    // set this to false to disable this
    // parsePages: true,

    // If parsePages option is disabled, the module will look for custom routes in
    // the pages option, refer to the "Routing" section for usage
    // pages: {
    //   'inspire':{
    //     en:'/locales/en/inspire.js',
    //     zh:'/locales/zh/inspire'
    //   }
    // },

    // By default, custom paths will be encoded using encodeURI method.
    // This does not work with regexp: "/foo/:slug-:id(\\d+)". If you want to use
    // regexp in the path, then set this option to false, and make sure you process
    // path encoding yourself.
    encodePaths: true,

    // Called right before app's locale changes
    // beforeLanguageSwitch: () => null,

    // Called after app's locale has changed
    // onLanguageSwitched: () => null
  },

  // Axios module configuration: https://go.nuxtjs.dev/config-axios
  axios: {
    // Workaround to avoid enforcing hard-coded localhost:3000: https://github.com/nuxt-community/axios-module/issues/308
    baseURL: '/',
  },

  // PWA module configuration: https://go.nuxtjs.dev/pwa
  pwa: {
    workbox: false,
    manifest: {
      lang: 'en'
    }
  },

  // Vuetify module configuration: https://go.nuxtjs.dev/config-vuetify
  vuetify: {
    customVariables: ['~/assets/variables.scss'],
    theme: {
      dark: true,
      themes: {
        dark: {
          primary: colors.blue.darken2,
          accent: colors.grey.darken3,
          secondary: colors.amber.darken3,
          info: colors.teal.lighten1,
          warning: colors.amber.base,
          error: colors.deepOrange.accent4,
          success: colors.green.accent3
        }
      }
    }
  },

  router: {
    base: process.env.BASE_URL || '',
    extendRoutes (routes, resolve) {
      routes.push({
        name: 'index',
        path: '/index.html',
        component: resolve(__dirname, 'pages/index.vue'),
      });
    },
  },

  generate: {
    subFolders: false
  },

  // Build Configuration: https://go.nuxtjs.dev/config-build
  build: {}
}
