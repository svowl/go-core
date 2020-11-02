module go-core/6-lesson/cmd/engine

go 1.15

replace go-core/6-lesson/pkg/spider v0.0.0 => ../../pkg/spider

replace go-core/6-lesson/pkg/spider/memspider v0.0.0 => ../../pkg/spider/memspider

replace go-core/6-lesson/pkg/index v0.0.0 => ../../pkg/index

replace go-core/6-lesson/pkg/btree => ../../pkg/btree

replace go-core/6-lesson/pkg/storage => ../../pkg/storage

replace go-core/6-lesson/pkg/storage/memstorage => ../../pkg/storage/memstorage

require (
	go-core/6-lesson/pkg/index v0.0.0
	go-core/6-lesson/pkg/spider v0.0.0
	go-core/6-lesson/pkg/spider/memspider v0.0.0
	go-core/6-lesson/pkg/storage v0.0.0-00010101000000-000000000000
	go-core/6-lesson/pkg/storage/memstorage v0.0.0-00010101000000-000000000000
)
