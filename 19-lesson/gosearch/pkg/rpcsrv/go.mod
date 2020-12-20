module gosearch/pkg/rpcsrv

go 1.15

replace gosearch/pkg/engine => ../engine

replace gosearch/pkg/crawler => ../crawler

replace gosearch/pkg/index => ../index

replace gosearch/pkg/storage => ../storage

replace gosearch/pkg/storage/mem => ../storage/mem

replace gosearch/pkg/storage/btree => ../storage/btree

require (
	gosearch/pkg/crawler v0.0.0-00010101000000-000000000000
	gosearch/pkg/engine v0.0.0-00010101000000-000000000000
)
