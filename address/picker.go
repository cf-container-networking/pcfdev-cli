package address

import (
	"fmt"

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

func (p *Picker) SelectAvailableIP(reusableInterfaces []*network.Interface) (ip string, err error) {
	allInterfaces, err := p.Network.Interfaces()
	if err != nil {
		return "", err
	}

	for _, subnetIP := range allowedSubnets {
		if p.nonReusableInterfaceExists(subnetIP, reusableInterfaces, allInterfaces) {
			continue
		}

		matchingAddrs := p.addrsInSet(subnetIP, reusableInterfaces)
		if len(matchingAddrs) == 1 {
			inUse, err := p.Driver.IsInterfaceInUse(matchingAddrs[0].Name)
			if err != nil {
				return "", err
			}

			if inUse {
				continue
			}
			return matchingAddrs[0].IP, nil
		}
		if len(matchingAddrs) > 1 {
			continue
		}

		return subnetIP, nil
	}

	return "", fmt.Errorf("all allowed network interfaces are currently taken")
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
