const routes = [
  // Splash screen (entry point)
  {
    path: '/',
    component: () => import('pages/SplashPage.vue'),
  },

  // Main catalogue with layout
  {
    path: '/simulations',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      {
        path: '',
        component: () => import('pages/CataloguePage.vue'),
      },
    ],
  },

  // Blog
  {
    path: '/blog',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      {
        path: '',
        component: () => import('pages/BlogPage.vue'),
      },
      {
        path: ':slug',
        component: () => import('pages/BlogArticlePage.vue'),
      },
    ],
  },

  // Always leave this as last one
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
]

export default routes
