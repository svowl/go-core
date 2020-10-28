module go-core/5-lesson/cmd/engine

go 1.15

replace go-core/5-lesson/pkg/spider v0.0.0 => ../../pkg/spider

replace go-core/5-lesson/pkg/index v0.0.0 => ../../pkg/index

replace go-core/5-lesson/pkg/btree => ../../pkg/btree

require (
	go-core/5-lesson/pkg/index v0.0.0
	go-core/5-lesson/pkg/spider v0.0.0
)
