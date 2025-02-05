import { viteBundler } from "@vuepress/bundler-vite";
import { markdownContainerPlugin } from "@vuepress/plugin-markdown-container";
import { defineUserConfig } from "vuepress";
import { hopeTheme } from "vuepress-theme-hope";

export default defineUserConfig({
  bundler: viteBundler({}),
  theme: hopeTheme({
    markdown: {
      tabs: true,
    },
    plugins: {},
    sidebar: [
      {
        text: "Home",
        link: "/",
        icon: "material-symbols:home-outline-rounded",
      },
      {
        text: "Contexts",
        link: "/supercontext/",
        icon: "material-symbols:contextual-token-outline-rounded",
        collapsible: false,
        expanded: true,
        children: [
          {
            text: "PostgreSQL (Bun)",
            link: "/supercontext/postgres_bun",
            icon: "devicon-plain:postgresql-wordmark",
          },
        ],
      },
    ],
  }),

  base: "/configurator/",

  plugins: [markdownContainerPlugin({})],

  lang: "en-US",
  title: "Configurator - A-Novel Kit",
});
