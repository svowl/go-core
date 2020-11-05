module go-core/7-lesson/cmd/engine

go 1.15

replace go-core/7-lesson/cmd/builder v0.0.0 => ../builder

replace go-core/7-lesson/cmd/search v0.0.0 => ../search

replace go-core/7-lesson/pkg/spider v0.0.0 => ../../pkg/spider

replace go-core/7-lesson/pkg/spider/mem v0.0.0 => ../../pkg/spider/mem

replace go-core/7-lesson/pkg/index v0.0.0 => ../../pkg/index

replace go-core/7-lesson/pkg/btree => ../../pkg/btree

replace go-core/7-lesson/pkg/storage => ../../pkg/storage

replace go-core/7-lesson/pkg/storage/file => ../../pkg/storage/file

replace go-core/7-lesson/pkg/storage/mem => ../../pkg/storage/mem

require (
	go-core/7-lesson/cmd/builder v0.0.0
	go-core/7-lesson/cmd/search v0.0.0
	go-core/7-lesson/pkg/index v0.0.0
	go-core/7-lesson/pkg/spider v0.0.0
	go-core/7-lesson/pkg/spider/mem v0.0.0
	go-core/7-lesson/pkg/storage v0.0.0-00010101000000-000000000000
	go-core/7-lesson/pkg/storage/file v0.0.0-00010101000000-000000000000
	go-core/7-lesson/pkg/storage/mem v0.0.0-00010101000000-000000000000
)
