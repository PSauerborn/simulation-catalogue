<template>
  <q-page class="blog-article-page">
    <header class="page-header">
      <div class="header-content">
        <q-btn
          flat
          no-caps
          icon="eva-arrow-back-outline"
          label="Back to Blog"
          to="/blog"
          class="back-btn"
        />
        <div class="header-text" v-if="article">
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
          <h1 class="page-title">{{ article.title }}</h1>
          <p class="page-subtitle">{{ article.summary }}</p>
          <div class="article-tags" v-if="article.tags?.length">
            <q-badge v-for="tag in article.tags" :key="tag" outline color="grey-6" :label="tag" />
          </div>
        </div>
      </div>
    </header>

    <div class="article-container" v-if="article">
      <article class="glass-card article-content">
        <component :is="article.component" />
      </article>
    </div>

    <div class="not-found" v-else>
      <q-icon name="eva-file-remove-outline" size="80px" color="grey-6" />
      <h3>Article Not Found</h3>
      <p>The article you're looking for doesn't exist.</p>
      <q-btn flat no-caps label="Back to Blog" to="/blog" color="primary" />
    </div>
  </q-page>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useMeta } from 'quasar'
import { getArticleBySlug } from 'src/articles'

const route = useRoute()
const article = computed(() => getArticleBySlug(route.params.slug))

// Dynamic SEO Meta Tags based on article
useMeta(() => {
  const baseUrl = 'https://simulation-catalogue.s31-software.com'
  const currentArticle = article.value

  if (!currentArticle) {
    return {
      title: 'Article Not Found',
      titleTemplate: (title) => `${title} | Simulation Catalogue`,
    }
  }

  const articleUrl = `${baseUrl}/blog/${currentArticle.slug}`
  const keywords = currentArticle.tags?.join(', ') || 'physics, fortran, simulations'

  return {
    title: currentArticle.title,
    titleTemplate: (title) => `${title} | Simulation Catalogue`,
    meta: {
      description: {
        name: 'description',
        content: currentArticle.summary,
      },
      keywords: {
        name: 'keywords',
        content: keywords,
      },
      author: {
        name: 'author',
        content: 'Simulation Catalogue',
      },
      // Open Graph
      ogTitle: { property: 'og:title', content: `${currentArticle.title} | Simulation Catalogue` },
      ogDescription: { property: 'og:description', content: currentArticle.summary },
      ogType: { property: 'og:type', content: 'article' },
      ogUrl: { property: 'og:url', content: articleUrl },
      ogSiteName: { property: 'og:site_name', content: 'Simulation Catalogue' },
      ogArticlePublishedTime: { property: 'article:published_time', content: currentArticle.date },
      ogArticleTag: { property: 'article:tag', content: currentArticle.tag },
      // Twitter Card
      twitterCard: { name: 'twitter:card', content: 'summary_large_image' },
      twitterTitle: { name: 'twitter:title', content: currentArticle.title },
      twitterDescription: { name: 'twitter:description', content: currentArticle.summary },
    },
    link: {
      canonical: { rel: 'canonical', href: articleUrl },
    },
  }
})
</script>

<style lang="scss" scoped>
@import '../css/quasar.variables.scss';

.blog-article-page {
  min-height: 100vh;
  padding-bottom: 48px;
}

.page-header {
  padding: 32px 24px 24px;
  background: linear-gradient(180deg, rgba(48, 209, 88, 0.08) 0%, transparent 100%);

  .header-content {
    max-width: 1400px;
    margin: 0 auto;
  }

  .back-btn {
    color: #a1a1a6;
    margin-bottom: 24px;
    font-size: 0.85rem;

    &:hover {
      color: #f5f5f7;
    }
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
    margin-top: 16px;
  }

  .page-title {
    font-size: 2.5rem;
    font-weight: 700;
    margin: 0 0 12px 0;
    letter-spacing: -0.02em;
    color: #f5f5f7;

    @media (max-width: 768px) {
      font-size: 1.75rem;
    }
  }

  .page-subtitle {
    font-size: 1.1rem;
    color: #a1a1a6;
    margin: 0;
    max-width: 700px;
    line-height: 1.6;

    @media (max-width: 768px) {
      font-size: 0.95rem;
    }
  }
}

.article-container {
  max-width: 1400px;
  margin: 32px auto 0;
  padding: 0 24px;

  @media (max-width: 600px) {
    padding: 0 12px;
    margin-top: 24px;
  }
}

.article-content {
  background: $glass-bg;
  border: 1px solid $glass-border;
  border-radius: $border-radius-lg;
  padding: 40px;
  max-width: 800px;
  margin: 0 auto;

  @media (max-width: 600px) {
    padding: 24px 16px;
  }

  :deep(.article-body) {
    p {
      font-size: 1rem;
      color: #d1d1d6;
      line-height: 1.7;
      margin: 0 0 16px 0;

      &:last-child {
        margin-bottom: 0;
      }
    }

    ul {
      color: #d1d1d6;
      padding-left: 24px;
      margin: 0 0 16px 0;

      li {
        margin-bottom: 8px;
        line-height: 1.6;
      }
    }

    ol {
      font-size: 1rem;
      color: #d1d1d6;
      padding-left: 24px;
      margin: 0 0 16px 0;
      line-height: 1.7;

      li {
        margin-bottom: 12px;
        line-height: 1.7;
      }
    }

    .inline-link {
      color: $primary;
      text-decoration: none;
      font-weight: 500;

      &:hover {
        text-decoration: underline;
      }
    }
  }
}

.not-found {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
  text-align: center;
  padding: 48px 24px;

  h3 {
    color: #f5f5f7;
    margin: 24px 0 8px 0;
  }

  p {
    color: #a1a1a6;
    max-width: 400px;
  }
}
</style>
