module gosearch/pkg/index

go 1.15

replace gosearch/pkg/storage => ../../pkg/storage

replace gosearch/pkg/storage/mem => ../../pkg/storage/mem

replace gosearch/pkg/storage/btree => ../../pkg/storage/btree

replace gosearch/pkg/crawler => ../../pkg/crawler

require (
	gosearch/pkg/crawler v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/mem v0.0.0-00010101000000-000000000000
)
