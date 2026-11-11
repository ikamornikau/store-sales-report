# frontend

This template should help get you started developing with Vue 3 in Vite.

## Recommended IDE Setup

[VS Code](https://code.visualstudio.com/) + [Vue (Official)](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Recommended Browser Setup

- Chromium-based browsers (Chrome, Edge, Brave, etc.):
  - [Vue.js devtools](https://chromewebstore.google.com/detail/vuejs-devtools/nhdogjmejiglipccpnnnanhbledajbpd)
  - [Turn on Custom Object Formatter in Chrome DevTools](http://bit.ly/object-formatters)
- Firefox:
  - [Vue.js devtools](https://addons.mozilla.org/en-US/firefox/addon/vue-js-devtools/)
  - [Turn on Custom Object Formatter in Firefox DevTools](https://fxdx.dev/firefox-devtools-custom-object-formatters/)

## Type Support for `.vue` Imports in TS

TypeScript cannot handle type information for `.vue` imports by default, so we replace the `tsc` CLI with `vue-tsc` for type checking. In editors, we need [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) to make the TypeScript language service aware of `.vue` types.

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Project Setup

```sh
pnpm install
```

### Compile and Hot-Reload for Development

```sh
pnpm dev
```

### Type-Check, Compile and Minify for Production

```sh
pnpm build
```

### Build the embeddable library

`ReportScreen` is published from this package as an ES module. `vue`, `chart.js`, and `vue-chartjs` stay outside the bundle. The app image is unchanged: `pnpm build` still writes the SPA to `dist`.

```sh
pnpm build:lib
```

In a Vue 3 host:

```ts
import { ReportScreen } from 'store-sales-report'
import 'store-sales-report/style.css'
```

```vue
<ReportScreen api-base-url="https://reports.example" />
```

Pass `locale` (`en`, `es`, or `fr`) and a `fetcher` (`typeof fetch`) when the host owns the language or the HTTP client.

### Run Unit Tests with [Vitest](https://vitest.dev/)

```sh
npm test
```

### Lint with [ESLint](https://eslint.org/)

```sh
pnpm lint
```
