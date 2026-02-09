<template>
  <q-page class="blog-page">
    <!-- Header Section -->
    <header class="page-header">
      <div class="header-content">
        <div class="header-text">
          <h1 class="page-title"><span class="gradient-text">Blog</span></h1>
          <p class="page-subtitle">
            Articles on physics simulations, Fortran, and computational science
          </p>
        </div>
      </div>
    </header>

    <!-- Search Bar -->
    <div class="search-bar">
      <q-input
        v-model="searchQuery"
        placeholder="Search articles..."
        filled
        dense
        dark
        class="search-input"
      >
        <template #prepend>
          <q-icon name="eva-search-outline" color="grey" />
        </template>
        <template #append v-if="searchQuery">
          <q-icon name="eva-close-outline" @click="searchQuery = ''" class="cursor-pointer" />
        </template>
      </q-input>
    </div>

    <!-- Article List -->
    <div class="blog-container">
      <div v-if="filteredArticles.length === 0" class="empty-state">
        <q-icon name="eva-file-text-outline" size="64px" color="grey-6" />
        <h3>No Articles Found</h3>
        <p>No articles match "{{ searchQuery }}". Try a different search.</p>
      </div>

      <router-link
        v-for="article in filteredArticles"
        :key="article.slug"
        :to="`/blog/${article.slug}`"
        class="article-card glass-card"
      >
        <div class="article-meta">
          <span class="article-date">
            <q-icon name="eva-calendar-outline" size="16px" />
            {{ article.date }}
          </span>
          <span class="article-read-time">
            <q-icon name="eva-clock-outline" size="16px" />
            {{ article.readTime }}
          </span>
          <q-badge color="primary" outline :label="article.tag" />
        </div>
        <h2 class="article-title">{{ article.title }}</h2>
        <p class="article-summary">{{ article.summary }}</p>
        <div class="article-tags" v-if="article.tags?.length">
          <q-badge
            v-for="tag in article.tags"
            :key="tag"
            outline
            color="grey-6"
            :label="tag"
            class="article-tag"
          />
        </div>
        <span class="read-more">
          Read article
          <q-icon name="eva-arrow-forward-outline" size="16px" />
        </span>
      </router-link>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useMeta } from 'quasar'
import { getArticles } from 'src/articles'

// SEO Meta Tags
useMeta({
  title: 'Blog',
  titleTemplate: (title) => `${title} | Simulation Catalogue`,
  meta: {
    description: {
      name: 'description',
      content:
        'Articles on physics simulations, Fortran programming, and computational science. Learn about particle dynamics, numerical integration, and more.',
    },
    keywords: {
      name: 'keywords',
      content:
        'physics blog, fortran tutorials, computational physics, boris pusher, velocity verlet, numerical methods, particle simulation',
    },
    // Open Graph
    ogTitle: { property: 'og:title', content: 'Blog | Simulation Catalogue' },
    ogDescription: {
      property: 'og:description',
      content: 'Articles on physics simulations, Fortran programming, and computational science.',
    },
    ogType: { property: 'og:type', content: 'blog' },
    ogUrl: { property: 'og:url', content: 'https://simulation-catalogue.s31-software.com/blog' },
    // Twitter Card
    twitterCard: { name: 'twitter:card', content: 'summary' },
    twitterTitle: { name: 'twitter:title', content: 'Blog | Simulation Catalogue' },
    twitterDescription: {
      name: 'twitter:description',
      content: 'Articles on physics simulations, Fortran programming, and computational science.',
    },
  },
  link: {
    canonical: { rel: 'canonical', href: 'https://simulation-catalogue.s31-software.com/blog' },
  },
})

const articles = getArticles()
const searchQuery = ref('')

const filteredArticles = computed(() => {
  if (!searchQuery.value) return articles

  const query = searchQuery.value.toLowerCase()
  return articles.filter((article) => {
    return (
      article.title.toLowerCase().includes(query) ||
      article.summary.toLowerCase().includes(query) ||
      article.tag.toLowerCase().includes(query) ||
      article.tags?.some((t) => t.toLowerCase().includes(query))
    )
  })
})
</script>

<style lang="scss" scoped>
@import '../css/quasar.variables.scss';

.blog-page {
  min-height: 100vh;
  padding-bottom: 48px;
}

.page-header {
  padding: 48px 24px 24px;
  background: linear-gradient(180deg, rgba(48, 209, 88, 0.08) 0%, transparent 100%);

  .header-content {
    max-width: 1400px;
    margin: 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    flex-wrap: wrap;
    gap: 24px;
  }

  .header-text {
    .page-title {
      font-size: 2.5rem;
      font-weight: 700;
      margin: 0 0 8px 0;
      letter-spacing: -0.02em;

      @media (max-width: 768px) {
        font-size: 1.75rem;
      }
    }

    .page-subtitle {
      font-size: 1.1rem;
      color: #a1a1a6;
      margin: 0;

      @media (max-width: 768px) {
        font-size: 0.95rem;
      }
    }
  }
}

.search-bar {
  max-width: 1400px;
  margin: 24px auto 0;
  padding: 0 24px;
  display: flex;
  justify-content: center;

  .search-input {
    width: 100%;
    max-width: 800px;
  }

  @media (max-width: 600px) {
    padding: 0 12px;
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 200px;
  text-align: center;
  padding: 48px 24px;
  width: 100%;

  h3 {
    color: #f5f5f7;
    margin: 16px 0 8px 0;
    font-size: 1.25rem;
  }

  p {
    color: #a1a1a6;
    max-width: 400px;
    margin: 0;
  }
}

.blog-container {
  max-width: 1400px;
  margin: 32px auto 0;
  padding: 0 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;

  @media (max-width: 600px) {
    padding: 0 12px;
    margin-top: 24px;
  }
}

.article-card {
  display: block;
  background: $glass-bg;
  border: 1px solid $glass-border;
  border-radius: $border-radius-lg;
  padding: 32px;
  width: 100%;
  max-width: 800px;
  text-decoration: none;
  transition:
    border-color 0.2s ease,
    background 0.2s ease;

  &:hover {
    border-color: rgba(255, 255, 255, 0.2);
    background: rgba(255, 255, 255, 0.07);
  }

  @media (max-width: 600px) {
    padding: 24px 16px;
  }

  .article-meta {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;

    .article-date,
    .article-read-time {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 0.85rem;
      color: #a1a1a6;
    }
  }

  .article-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 16px;
  }

  .article-title {
    font-size: 1.5rem;
    font-weight: 700;
    color: #f5f5f7;
    margin: 0 0 12px 0;
    letter-spacing: -0.01em;

    @media (max-width: 600px) {
      font-size: 1.25rem;
    }
  }

  .article-summary {
    font-size: 1rem;
    color: #a1a1a6;
    line-height: 1.6;
    margin: 0 0 16px 0;
  }

  .read-more {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-size: 0.9rem;
    font-weight: 500;
    color: $primary;
  }
}
</style>
