import path from 'node:path'
import { defineConfig } from '@rspack/cli';
import { rspack } from "@rspack/core";
import { purgeCSSPlugin } from '@fullhuman/postcss-purgecss';
import browserslist from 'browserslist';
import { glob } from 'glob';
import { fileURLToPath } from 'node:url';

const BASE_DIR = './frontend'
const ASSET_PATH = process.env.ASSET_PATH || '/assets/';
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const extraJsBundles = (await glob(`${BASE_DIR}/scripts/**/*.js`, { dotRelative: true }))
  .reduce((acc, p) => ({ ...acc, [path.basename(p, ".js")]: p }), {});

export default defineConfig({
  entry: {
    main: [
      `${BASE_DIR}/src/index.js`,
    ],
    ...extraJsBundles
  },
  module: {
    rules: [
      {
        test: /\.(sass|scss)$/,
        use: [
          {
            loader: 'postcss-loader',
            options: {
              postcssOptions: {
                plugins: [
                  purgeCSSPlugin({
                    content: ['./templates/**/*.tmpl', `${BASE_DIR}/**/*.js`],
                    variables: true,
                    safelist: {
                      greedy: [/datepicker/],
                    },
                  }),
                ],
              },
            },
          },
          {
            loader: 'sass-loader',
            options: {
              // recommended by Rspack docs
              api: 'modern-compiler',
              implementation: await import('sass-embedded'),
            }
          },
        ],
        type: 'css/auto',
      },
      {
        test: /\.js$/,
        use: {
          loader: 'builtin:swc-loader',
          options: {
            env: {
              targets: browserslist(),
            },
          },
        },
      },
    ],
  },
  plugins: [
    new rspack.CopyRspackPlugin({
      patterns: [{ from: `${BASE_DIR}/images`, to: 'images' }],
    }),
  ],
  output: {
    path: path.resolve(__dirname, 'static'),
    filename: 'js/[name].js',
    cssFilename: 'css/[name].css',
    publicPath: ASSET_PATH,
  },
  optimization: {
    minimizer: [
      new rspack.SwcJsMinimizerRspackPlugin(),
      new rspack.LightningCssMinimizerRspackPlugin(),
    ]
  },
  experiments: {
    css: true,
  },
  devServer: {
    proxy: [
      {
        context: ['/'],
        target: "http://localhost:8080",
      }
    ],
    // onListening: consider starting the backend server
  }
});