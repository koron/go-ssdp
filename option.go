package ssdp

import (
	"github.com/koron/go-ssdp/internal/multicast"
)

type config struct {
	udpConfig
	multicastConfig
	advertiseConfig
}

func opts2config(opts []Option) (cfg config, err error) {
	for _, o := range opts {
		err := o.apply(&cfg)
		if err != nil {
			return config{}, err
		}
	}
	return cfg, nil
}

type udpConfig struct {
	laddr string
	raddr string
}

func (uc udpConfig) laddrResolver() multicast.Resolver {
	if uc.laddr == "" {
		return multicast.LocalAddrResolver
	}
	return multicast.AddressResolver(uc.laddr)
}

func (uc udpConfig) raddrResolver() multicast.Resolver {
	if uc.raddr == "" {
		return multicast.RemoteAddrResolver
	}
	return multicast.AddressResolver(uc.raddr)
}

type multicastConfig struct {
	ttl   int
	sysIf bool
}

func (mc multicastConfig) options() (opts []multicast.ConnOption) {
	if mc.ttl > 0 {
		opts = append(opts, multicast.ConnTTL(mc.ttl))
	}
	if mc.sysIf {
		opts = append(opts, multicast.ConnSystemAssginedInterface())
	}
	return opts
}

type advertiseConfig struct {
	addHost bool
}

// Option is option set for SSDP API.
type Option interface {
	apply(c *config) error
}

type optionFunc func(*config) error

func (of optionFunc) apply(c *config) error {
	return of(c)
}

// LocalAddr return as Option that specify local address.
// This option works with Advertize() and Monitor only.
// Default "224.0.0.1:1900" will be used for Advertize() and Monitor()
// when omitted the option.
func LocalAddr(laddr string) Option {
	return optionFunc(func(c *config) error {
		c.udpConfig.laddr = laddr
		return nil
	})
}

// RemoteAddr return as Option that specify remote address.
// Default "239.255.255.250:1900" will be used when omitted the option.
func RemoteAddress(raddr string) Option {
	return optionFunc(func(c *config) error {
		c.udpConfig.raddr = raddr
		return nil
	})
}

// TTL returns as Option that set TTL for multicast packets.
func TTL(ttl int) Option {
	return optionFunc(func(c *config) error {
		c.multicastConfig.ttl = ttl
		return nil
	})
}

// OnlySystemInterface returns as Option that using only a system assigned
// multicast interface.
func OnlySystemInterface() Option {
	return optionFunc(func(c *config) error {
		c.multicastConfig.sysIf = true
		return nil
	})
}

// AdvertiseHost returns as Option that add HOST header to response for
// M-SEARCH requests.
// This option works with Advertise() function only.
// This is added to support SmartThings.
// See https://github.com/koron/go-ssdp/issues/30 for details.
func AdvertiseHost() Option {
	return optionFunc(func(c *config) error {
		c.advertiseConfig.addHost = true
		return nil
	})
}
