package ssdp

func newTestMonitor(typ string, alive AliveHandler, bye ByeHandler, search SearchHandler) *Monitor {
	m := &Monitor{}
	if alive != nil {
		m.Alive = func(am *AliveMessage) {
			if am.Type == typ {
				alive(am)
			}
		}
	}
	if bye != nil {
		m.Bye = func(bm *ByeMessage) {
			if bm.Type == typ {
				bye(bm)
			}
		}
	}
	if search != nil {
		m.Search = func(sm *SearchMessage) {
			if sm.Type == typ {
				search(sm)
			}
		}
	}
	return m
}
