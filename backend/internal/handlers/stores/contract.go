package stores

import "github.com/ikamornikau/store-sales-report/backend/internal/adapters/postgres"

type Store struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Timezone string `json:"timezone"`
}

func storesToAPI(stores []postgres.Store) []Store {
	apiStores := make([]Store, 0, len(stores))
	for _, store := range stores {
		apiStores = append(apiStores, Store(store))
	}

	return apiStores
}
