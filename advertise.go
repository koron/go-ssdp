package ssdp

// Advertiser is a server to advertise a service.
type Advertiser struct {
}

// Advertise starts advertisement of service.
func Advertise(st, usn, location, server string, maxAge int) (*Advertiser, error) {
	// TODO:
	return nil, nil
}

// Close stops advertisement.
func (a *Advertiser) Close() error {
	// TODO:
	return nil
}

// Alive announces ssdp:alive message.
func (a *Advertiser) Alive() error {
	// TODO:
	return nil
}

// Bye announces ssdp:byebye message.
func (a *Advertiser) Bye() error {
	// TODO:
	return nil
}
