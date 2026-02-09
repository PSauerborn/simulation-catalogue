import { defineAsyncComponent } from 'vue'

const articles = [
  {
    slug: 'initial-release',
    title: 'A Charged Particle in a Magnetic Field',
    date: 'January 30, 2026',
    tag: 'Initial Release',
    readTime: '15 min read',
    tags: ['Fortran', 'Vue', 'Golang'],
    summary: 'Simulation Catalogue 101. Why it exists, how its setup and what it does.',
    component: defineAsyncComponent(() => import('./InitialArticle.vue')),
  },
  {
    slug: 'boris-pusher',
    title: 'Update: Implementing the Boris Push Algorithm',
    date: 'February 8, 2026',
    tag: 'Update',
    readTime: '15 min read',
    tags: ['Physics', 'Fortran', 'Particle Dynamics'],
    summary: 'Implementation of the Boris Push Algorithm for Lorentz force integration.',
    component: defineAsyncComponent(() => import('./BorisPusherArticle.vue')),
  },
]

export function getArticles() {
  return articles
}

export function getArticleBySlug(slug) {
  return articles.find((a) => a.slug === slug) || null
}
