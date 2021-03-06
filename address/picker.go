package address

import (
	"fmt"

	cfg "github.com/pivotal-cf/pcfdev-cli/config"
	"github.com/pivotal-cf/pcfdev-cli/network"
)

//go:generate mockgen -package mocks -destination mocks/network.go github.com/pivotal-cf/pcfdev-cli/address Network
type Network interface {
	Interfaces() (interfaces []*network.Interface, err error)
}

//go:generate mockgen -package mocks -destination mocks/driver.go github.com/pivotal-cf/pcfdev-cli/address Driver
type Driver interface {
	IsInterfaceInUse(interfaceName string) (inUse bool, err error)
}

type Picker struct {
	Network Network
	Driver  Driver
}

func (p *Picker) SelectAvailableInterface(reusableInterfaces []*network.Interface, config *cfg.VMConfig) (*cfg.NetworkConfig, error) {
	if config.IP != "" || config.Domain != "" {
		var subnetIP, ip, domain string
		var err error

		if config.IP != "" {
			subnetIP, err = SubnetForIP(config.IP)
			if err != nil {
				return nil, err
			}
			ip = config.IP
		} else {
			subnetIP, err = SubnetForDomain(config.Domain)
			if err != nil {
				return nil, err
			}
			ip = IPForSubnet(subnetIP)
		}

		if config.Domain != "" {
			domain = config.Domain
		} else {
			domain = DomainForIP(ip)
		}

		var networkInterface *network.Interface
		if addrs := p.addrsInSet(subnetIP, reusableInterfaces); len(addrs) > 0 {
			networkInterface = addrs[0]
		} else {
			networkInterface = &network.Interface{
				IP:     subnetIP,
				Exists: false,
			}
		}

		return &cfg.NetworkConfig{
			VMIP:      ip,
			VMDomain:  domain,
			Interface: networkInterface,
		}, nil
	}

	allInterfaces, err := p.Network.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, subnetIP := range allowedSubnets {
		if p.nonReusableInterfaceExists(subnetIP, reusableInterfaces, allInterfaces) {
			continue
		}

		matchingAddrs := p.addrsInSet(subnetIP, reusableInterfaces)
		domain := DomainForIP(IPForSubnet(subnetIP))

		switch len(matchingAddrs) {
		case 0:
			return &cfg.NetworkConfig{
				VMIP:     IPForSubnet(subnetIP),
				VMDomain: domain,
				Interface: &network.Interface{
					IP:     subnetIP,
					Exists: false,
				},
			}, nil
		case 1:
			inUse, err := p.Driver.IsInterfaceInUse(matchingAddrs[0].Name)
			if err != nil {
				return nil, err
			}

			if inUse {
				continue
			}

			return &cfg.NetworkConfig{
				VMIP:      IPForSubnet(subnetIP),
				VMDomain:  domain,
				Interface: matchingAddrs[0],
			}, nil
		}
	}

	return nil, fmt.Errorf("all allowed network interfaces are currently taken")
}

func (p *Picker) addrsInSet(ip string, set []*network.Interface) (addrs []*network.Interface) {
	addrs = make([]*network.Interface, 0, 1)
	for _, addr := range set {
		if addr.IP == ip {
			addrs = append(addrs, addr)
		}
	}

	return addrs
}

func (p *Picker) nonReusableInterfaceExists(ip string, reusableInterfaces []*network.Interface, allInterfaces []*network.Interface) bool {
	for _, iface := range allInterfaces {
		reusable := false
		for _, reusableIface := range reusableInterfaces {
			if iface.HardwareAddress == reusableIface.HardwareAddress {
				reusable = true
			}
		}

		if !reusable && ip == iface.IP {
			return true
		}
	}
	return false
}
