require("dotenv").config();
const webpack = require("webpack");

module.exports = function override(config, env) {
  config.plugins.push(
    new webpack.ProvidePlugin({
      Buffer: ["buffer", "Buffer"],
    }),
    new webpack.EnvironmentPlugin([
      "REACT_APP_CHAIN_ID",
      "REACT_APP_CHAIN_NAME",
      "REACT_APP_RPC_URL",
      "REACT_APP_REST_URL",
    ])
  );
  config.resolve.fallback = {
    buffer: false,
    crypto: false,
    events: false,
    path: false,
    stream: false,
    string_decoder: false,
  };
  return config;
};
