package closer

type CloseFn func()

func (fn *CloseFn) Stack(stacked func()) {
	self := *fn
	*fn = func() {
		stacked()
		self()
	}
}

func (fn CloseFn) Maybe() (cancel func(), close func()) {
	do := true
	cancel = func() { do = false }
	close = func() {
		if do {
			fn()
		}
	}
	return
}
