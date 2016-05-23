package ssdp

// Monitor monitors SSDP's alive and byebye messages.
type Monitor struct {
	alive AliveHandler
	bye   ByeHandler
}

// NewMonitor creates a new Monitor.
func NewMonitor(alive AliveHandler, bye ByeHandler) (*Monitor, error) {
	if alive == nil {
		alive = nullAlive
	}
	if bye == nil {
		bye = nullBye
	}
	// TODO:
	return &Monitor{
		alive: alive,
		bye:   bye,
	}, nil
}

// Close closes monitoring.
func (m *Monitor) Close() error {
	// TODO:
	return nil
}

// Alive represents SSDP's ssdp:alive message.
type Alive struct {
	// TODO:
}

// AliveHandler is handler of Alive message.
type AliveHandler func(*Alive)

func nullAlive(*Alive) {
	// nothing to do.
}

// Bye represents SSDP's ssdp:byebye message.
type Bye struct {
	// TODO:
}

// ByeHandler is handler of Bye message.
type ByeHandler func(*Bye)

func nullBye(*Bye) {
	// nothing to do.
}
