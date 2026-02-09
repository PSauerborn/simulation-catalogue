import katex from 'katex'
import 'katex/dist/katex.min.css'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css'
import { marked } from 'marked'

/**
 * Highlight code using highlight.js
 * @param {string} code - The code to highlight
 * @param {string} lang - The language for syntax highlighting
 * @returns {string} HTML string with highlighted code
 */
export function highlightCode(code, lang) {
  const langClass = lang ? `hljs language-${lang}` : 'hljs'
  let highlighted
  if (lang && hljs.getLanguage(lang)) {
    highlighted = hljs.highlight(code, { language: lang }).value
  } else {
    highlighted = hljs.highlightAuto(code).value
  }
  return `<pre><code class="${langClass}">${highlighted}</code></pre>`
}

/**
 * Render LaTeX inline using KaTeX
 * @param {string} input - LaTeX string to render
 * @returns {string} HTML string with rendered LaTeX
 */
export function renderLatexInline(input) {
  return katex.renderToString(input, { throwOnError: false, displayMode: false })
}

/**
 * Render LaTeX in display mode using KaTeX
 * @param {string} input - LaTeX string to render
 * @returns {string} HTML string with rendered LaTeX in display mode
 */
export function renderLatexDisplay(input) {
  return katex.renderToString(input, { throwOnError: false, displayMode: true })
}

/**
 * Render inline markdown (bold, italic, code, links) without block elements
 * @param {string} input - Markdown string to render
 * @returns {string} HTML string with rendered inline markdown
 */
export function renderMarkdownInline(input) {
  return marked.parseInline(input)
}
