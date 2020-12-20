module gosearch/cmd/gosearch

go 1.15

replace gosearch/pkg/crawler => ../../pkg/crawler

replace gosearch/pkg/crawler/spider => ../../pkg/crawler/spider

replace gosearch/pkg/crawler/membot => ../../pkg/crawler/membot

replace gosearch/pkg/engine => ../../pkg/engine

replace gosearch/pkg/index => ../../pkg/index

replace gosearch/pkg/storage/btree => ../../pkg/storage/btree

replace gosearch/pkg/storage => ../../pkg/storage

replace gosearch/pkg/storage/mem => ../../pkg/storage/mem

replace gosearch/pkg/webapp => ../../pkg/webapp

replace gosearch/pkg/api => ../../pkg/api

replace gosearch/pkg/rpcsrv => ../../pkg/rpcsrv

require (
	github.com/gorilla/mux v1.8.0
	gosearch/pkg/api v0.0.0-00010101000000-000000000000
	gosearch/pkg/crawler v0.0.0-00010101000000-000000000000
	gosearch/pkg/crawler/spider v0.0.0-00010101000000-000000000000
	gosearch/pkg/engine v0.0.0-00010101000000-000000000000
	gosearch/pkg/index v0.0.0
	gosearch/pkg/rpcsrv v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/btree v0.0.0-00010101000000-000000000000
	gosearch/pkg/storage/mem v0.0.0-00010101000000-000000000000
	gosearch/pkg/webapp v0.0.0-00010101000000-000000000000
)
