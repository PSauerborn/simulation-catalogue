const routes = [
  // Splash screen (entry point)
  {
    path: '/',
    component: () => import('pages/SplashPage.vue'),
  },

  // Main catalogue with layout
  {
    path: '/catalogue',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      {
        path: '',
        component: () => import('pages/CataloguePage.vue'),
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
