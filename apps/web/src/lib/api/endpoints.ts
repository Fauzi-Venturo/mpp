// ----------------------------------------------------------------------
// Centralized backend endpoint map (mirror of src/routes/paths.ts for API URLs).
// Contract: marketplace-be/docs/api-contract/marketplace/articles.md
//
// Paths are relative (no leading slash) so they resolve against ky's
// `baseUrl` even when it carries a path prefix.

export const endpoints = {
  articles: {
    list: 'api/articles',
    details: (slug: string) => `api/articles/${encodeURIComponent(slug)}`,
    categories: 'api/article-categories',
  },
  faq: {
    list: 'api/faq',
  },
  // MPP queue domain — contract: docs/04-api/rest-endpoints.md
  booking: {
    create: 'mpp/v1/booking',
    details: (id: string) => `mpp/v1/booking/${encodeURIComponent(id)}`,
  },
  siteContent: {
    map: 'api/site-content',
  },
};
