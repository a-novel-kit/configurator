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
    ],
  }),

  base: "/configurator/",

  plugins: [markdownContainerPlugin({})],

  lang: "en-US",
  title: "Configurator - A-Novel Kit",
});
