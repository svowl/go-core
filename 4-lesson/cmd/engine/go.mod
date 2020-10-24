module go-core/4-lesson/cmd/engine

go 1.15

replace go-core/4-lesson/pkg/spider v0.0.0 => ../../pkg/spider

replace go-core/4-lesson/pkg/fakespider v0.0.0 => ../../pkg/fakespider

replace go-core/4-lesson/pkg/index v0.0.0 => ../../pkg/index

require (
	go-core/4-lesson/pkg/fakespider v0.0.0
	go-core/4-lesson/pkg/index v0.0.0
	go-core/4-lesson/pkg/spider v0.0.0
)
