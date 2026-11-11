![Backend](https://github.com/ikamornikau/store-sales-report/actions/workflows/backend.yaml/badge.svg)
![Frontend](https://github.com/ikamornikau/store-sales-report/actions/workflows/frontend.yaml/badge.svg)

# store-sales-report

Sales report for one store: a Go API and a Vue screen that can run on its own or inside another Vue 3 app.

## Compose

```sh
docker compose up --build
```

| Service | Address |
| --- | --- |
| UI | http://localhost:8081 |
| API | http://localhost:8080 |
| Postgres | `localhost:5432`, user `store`, password `store`, database `store_sales_report` |

Compose applies migrations and loads the seed before the API starts. The UI container waits until `GET /ready` on the API succeeds. The UI proxies `/api/` to the API.

## Tests

### Backend

Postgres from Compose must be listening on `localhost:5432` with user `store` and database `store_sales_report`. The suite creates and drops `store_sales_report_test` and applies migrations itself.

```sh
docker compose up -d postgres
cd backend
go test ./...
```

Go version is in `backend/go.mod` (1.27.3).

### Frontend

```sh
cd frontend
pnpm install
npm test
```

Node 22 and pnpm 12.16.1, as in CI. `pnpm test:unit` starts Vitest in watch mode.

## API

The contract is the OpenAPI 3.0.3 document [backend/docs/openapi.yaml](backend/docs/openapi.yaml). Its server URL is `http://localhost:8080`.

## Embed the report screen

`ReportScreen` is published from `frontend` as an ES module. `vue`, `chart.js`, and `vue-chartjs` stay outside the bundle.

```sh
cd frontend
pnpm install
pnpm build:lib
```

The build writes `frontend/dist-lib`. Point the host app at this package and install the same peers:

```json
{
  "dependencies": {
    "store-sales-report": "file:../store-sales-report/frontend",
    "chart.js": "^4.6.1",
    "vue": "^3.5.42",
    "vue-chartjs": "^5.3.4"
  }
}
```

```ts
import { ReportScreen } from 'store-sales-report'
import 'store-sales-report/style.css'
```

```vue
<ReportScreen api-base-url="https://reports.example" />
```

`api-base-url` is the origin of this API. Pass `locale` (`en`, `es`, or `fr`) when the host owns the language, and `fetcher` (`typeof fetch`) when the host owns the HTTP client. The screen emits `update:locale`.
