module gosearch/pkg/api

go 1.15

replace gosearch/pkg/engine => ../engine
replace gosearch/pkg/index => ../index
replace gosearch/pkg/storage => ../storage
replace gosearch/pkg/storage/mem => ../storage/mem
replace gosearch/pkg/storage/btree => ../storage/btree
replace gosearch/pkg/crawler => ../crawler

require (
	github.com/gorilla/mux v1.8.0
	gosearch/pkg/crawler v0.0.0-00010101000000-000000000000
	gosearch/pkg/engine v0.0.0-00010101000000-000000000000
	gosearch/pkg/index v0.0.0
	gosearch/pkg/storage v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/btree v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/mem v0.0.0-00010101000000-000000000000
)
