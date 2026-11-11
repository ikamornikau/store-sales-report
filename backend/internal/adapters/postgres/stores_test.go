package postgres

import "context"

func (t *test) TestGetStores() {
	t.insertStore("S-1001", "Harbour", "CAD", "America/St_Johns")
	t.insertStore("S-1010", "Deli", "USD", "America/New_York")

	stores, err := t.postgres.GetStores(context.Background())
	t.Require().NoError(err)
	t.Require().ElementsMatch([]Store{
		{ID: "S-1001", Name: "Harbour", Currency: "CAD", Timezone: "America/St_Johns"},
		{ID: "S-1010", Name: "Deli", Currency: "USD", Timezone: "America/New_York"},
	}, stores)
}

func (t *test) TestGetStoreByID() {
	t.insertStore("S-1001", "Harbour", "CAD", "America/St_Johns")

	store, err := t.postgres.GetStoreByID(context.Background(), "S-1001")
	t.Require().NoError(err)
	t.Require().Equal(Store{
		ID: "S-1001", Name: "Harbour", Currency: "CAD", Timezone: "America/St_Johns",
	}, store)
}

func (t *test) TestGetStoreByIDNotFound() {
	_, err := t.postgres.GetStoreByID(context.Background(), "S-404")
	t.Require().ErrorIs(err, ErrNotFound)
}
