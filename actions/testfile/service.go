package myfeature

type Foo interface {
	Method()
}

type Bar interface {
}

func F() {
	type Inline interface {
	}

	var i Inline
	println(i)
}
