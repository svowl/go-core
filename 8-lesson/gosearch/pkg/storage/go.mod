module gosearch/pkg/storage

go 1.15

replace gosearch/pkg/crawler => ../../pkg/crawler

replace gosearch/pkg/storage/mem => ./mem

replace gosearch/pkg/storage/btree => ./btree

require (
	gosearch/pkg/crawler v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/btree v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/mem v0.0.0-00010101000000-000000000000
)
