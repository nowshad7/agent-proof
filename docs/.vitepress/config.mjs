import { defineConfig } from 'vitepress'

const base = '/agent-proof/'
const site = 'https://nowshad7.github.io/agent-proof/'

export default defineConfig({
  base,
  lang: 'en-US',
  title: 'AgentProof',
  description:
    'Reproducible verification for agent-written changes. Explicit acceptance contracts, Go/JUnit evidence, and portable commit-specific receipts.',
  cleanUrls: true,
  lastUpdated: true,
  ignoreDeadLinks: false,

  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: `${base}logo.svg` }],
    ['meta', { name: 'theme-color', content: '#0d9488' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:title', content: 'AgentProof' }],
    [
      'meta',
      {
        property: 'og:description',
        content:
          'Reproducible verification for agent-written changes. Explicit requirements, real test evidence, portable receipts.',
      },
    ],
    ['meta', { property: 'og:url', content: site }],
    ['meta', { property: 'og:image', content: `${site}social-preview.svg` }],
    ['meta', { name: 'twitter:card', content: 'summary_large_image' }],
  ],

  sitemap: { hostname: site },

  themeConfig: {
    logo: '/logo.svg',
    siteTitle: 'AgentProof',

    search: { provider: 'local' },

    nav: [
      { text: 'Guide', link: '/guide/introduction' },
      { text: 'Reference', link: '/reference/cli' },
      { text: 'Learn', link: '/learn/agent-verification' },
      {
        text: 'Preview 0.1',
        items: [
          { text: 'Changelog', link: '/community/changelog' },
          { text: 'Roadmap', link: '/community/roadmap' },
          { text: 'Contributing', link: '/community/contributing' },
        ],
      },
    ],

    sidebar: [
      {
        text: 'Start',
        collapsed: false,
        items: [
          { text: 'Introduction', link: '/guide/introduction' },
          { text: 'Quickstart', link: '/guide/quickstart' },
          { text: 'Sample receipt', link: '/guide/sample-receipt' },
        ],
      },
      {
        text: 'The contract',
        collapsed: false,
        items: [
          { text: 'Configuration', link: '/guide/configuration' },
          { text: 'Evidence states', link: '/guide/evidence' },
          { text: 'Requirement mapping', link: '/guide/requirements' },
        ],
      },
      {
        text: 'Framework guides',
        collapsed: false,
        items: [
          { text: 'Go test JSON', link: '/guides/go' },
          { text: 'JUnit XML', link: '/guides/junit' },
        ],
      },
      {
        text: 'Run & review',
        collapsed: false,
        items: [
          { text: 'Trusted local execution', link: '/guide/execution' },
          { text: 'Agent workflows', link: '/guide/agents' },
          { text: 'GitHub Actions / CI', link: '/guide/ci' },
          { text: 'Reports', link: '/guide/reports' },
          { text: 'Verification badge', link: '/guide/badge' },
        ],
      },
      {
        text: 'Reference',
        collapsed: false,
        items: [
          { text: 'CLI', link: '/reference/cli' },
          { text: 'Contract schema', link: '/reference/contract' },
          { text: 'Receipt format', link: '/reference/receipt' },
          { text: 'Troubleshooting', link: '/reference/troubleshooting' },
        ],
      },
      {
        text: 'Trust',
        collapsed: false,
        items: [{ text: 'Security model', link: '/guide/security' }],
      },
      {
        text: 'Learn',
        collapsed: true,
        items: [
          { text: 'Verifying AI-generated code', link: '/learn/agent-verification' },
          { text: 'Detecting missing evidence', link: '/learn/missing-evidence' },
          { text: 'Commit-specific receipts', link: '/learn/commit-receipts' },
        ],
      },
      {
        text: 'Community',
        collapsed: true,
        items: [
          { text: 'Contributing', link: '/community/contributing' },
          { text: 'Roadmap', link: '/community/roadmap' },
          { text: 'Changelog', link: '/community/changelog' },
        ],
      },
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/nowshad7/agent-proof' },
    ],

    editLink: {
      pattern:
        'https://github.com/nowshad7/agent-proof/edit/main/docs/:path',
      text: 'Edit this page on GitHub',
    },

    footer: {
      message:
        'Developer preview. Receipts are unsigned and prove only the declared checks. Apache-2.0 licensed.',
      copyright: 'Released under the Apache-2.0 License.',
    },
  },
})
