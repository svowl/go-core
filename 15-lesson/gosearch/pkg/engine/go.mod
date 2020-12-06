module gosearch/pkg/engine

go 1.15

replace gosearch/pkg/crawler => ../../pkg/crawler
replace gosearch/pkg/index => ../../pkg/index
replace gosearch/pkg/storage => ../../pkg/storage
replace gosearch/pkg/storage/btree => ../../pkg/storage/btree
replace gosearch/pkg/storage/mem => ../../pkg/storage/mem

require (
	golang.org/x/text v0.3.4
	gosearch/pkg/crawler v0.0.0-00010101000000-000000000000
	gosearch/pkg/index v0.0.0
	gosearch/pkg/storage v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/btree v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/mem v0.0.0-00010101000000-000000000000
)
