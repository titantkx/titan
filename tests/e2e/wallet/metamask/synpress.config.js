const { defineConfig } = require("cypress");
const synpressPlugins = require("@synthetixio/synpress/plugins");

module.exports = defineConfig({
  e2e: {
    baseUrl: "https://metamask.github.io/test-dapp/",
    specPattern: "specs",
    supportFile: "support.js",
    videosFolder: "videos",
    screenshotsFolder: "screenshots",
    video: false,
    screenshotOnRunFailure: false,
    setupNodeEvents(on, config) {
      synpressPlugins(on, config);
      return config;
    },
  },
});
