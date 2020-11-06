module go-core/search/cmd/engine

go 1.15

replace go-core/3-lesson/pkg/spider v0.0.0 => ../../pkg/spider
replace go-core/3-lesson/pkg/spider/mem v0.0.0 => ../../pkg/spider/mem

require (
	go-core/3-lesson/pkg/spider v0.0.0
	go-core/3-lesson/pkg/spider/mem v0.0.0
	golang.org/x/net v0.0.0-20200925080053-05aa5d4ee321 // indirect
)
