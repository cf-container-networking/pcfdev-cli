package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/pivotal-cf/pcfdev-cli/user"
)

type Config struct {
	DefaultVMName            string
	PCFDevHome               string
	OVADir                   string
	OVAPath                  string
	PartialOVAPath           string
	VMDir                    string
	HTTPProxy                string
	HTTPSProxy               string
	NoProxy                  string
	MinMemory                uint64
	MaxMemory                uint64
	TotalMemory              uint64
	FreeMemory               uint64
	DefaultMemory            uint64
	SpringCloudDefaultMemory uint64
	SpringCloudMinMemory     uint64
	SpringCloudMaxMemory     uint64
	DefaultCPUs              func() (int, error)
	ExpectedMD5              string
	InsecurePrivateKey       []byte
	PrivateKeyPath           string
	Version                  *Version
}

type Version struct {
	BuildVersion    string
	BuildSHA        string
	OVABuildVersion string
}

//go:generate mockgen -package mocks -destination mocks/system.go github.com/pivotal-cf/pcfdev-cli/config System
type System interface {
	TotalMemory() (uint64, error)
	FreeMemory() (uint64, error)
	PhysicalCores() (int, error)
}

func New(defaultVMName string, expectedMD5 string, insecurePrivateKey []byte, system System, version *Version) (*Config, error) {
	pcfdevHome, err := getPCFDevHome()
	if err != nil {
		return nil, err
	}
	freeMemory, err := system.FreeMemory()
	if err != nil {
		return nil, err
	}
	totalMemory, err := system.TotalMemory()
	if err != nil {
		return nil, err
	}
	minMemory := uint64(3072)
	maxMemory := uint64(4096)
	springCloudMinMemory := uint64(6144)
	springCloudMaxMemory := uint64(8192)

	return &Config{
		DefaultVMName:            defaultVMName,
		ExpectedMD5:              expectedMD5,
		PCFDevHome:               pcfdevHome,
		OVADir:                   filepath.Join(pcfdevHome, "ova"),
		VMDir:                    filepath.Join(pcfdevHome, "vms"),
		OVAPath:                  filepath.Join(pcfdevHome, "ova", defaultVMName+".ova"),
		PartialOVAPath:           filepath.Join(pcfdevHome, "ova", defaultVMName+".ova.partial"),
		HTTPProxy:                getHTTPProxy(),
		HTTPSProxy:               getHTTPSProxy(),
		NoProxy:                  getNoProxy(),
		MinMemory:                minMemory,
		MaxMemory:                maxMemory,
		TotalMemory:              totalMemory,
		FreeMemory:               freeMemory,
		DefaultMemory:            getDefaultMemory(totalMemory, minMemory, maxMemory),
		SpringCloudDefaultMemory: getDefaultMemory(totalMemory, springCloudMinMemory, springCloudMaxMemory),
		SpringCloudMinMemory:     springCloudMinMemory,
		SpringCloudMaxMemory:     springCloudMaxMemory,
		DefaultCPUs:              system.PhysicalCores,
		InsecurePrivateKey:       insecurePrivateKey,
		PrivateKeyPath:           filepath.Join(pcfdevHome, "vms", "key.pem"),
		Version:                  version,
	}, nil
}

func getPCFDevHome() (string, error) {
	if pcfdevHome := os.Getenv("PCFDEV_HOME"); pcfdevHome != "" {
		return pcfdevHome, nil
	}

	homeDir, err := user.GetHome()
	if err != nil {
		return "", fmt.Errorf("failed to find home directory: %s", err)
	}

	return filepath.Join(homeDir, ".pcfdev"), nil
}

func getHTTPProxy() string {
	if proxy := os.Getenv("HTTP_PROXY"); proxy != "" {
		return stripWhitespace(proxy)
	}
	return stripWhitespace(os.Getenv("http_proxy"))
}

func getHTTPSProxy() string {
	if proxy := os.Getenv("HTTPS_PROXY"); proxy != "" {
		return stripWhitespace(proxy)
	}
	return stripWhitespace(os.Getenv("https_proxy"))
}

func getNoProxy() string {
	if proxy := os.Getenv("NO_PROXY"); proxy != "" {
		return stripWhitespace(proxy)
	}
	return stripWhitespace(os.Getenv("no_proxy"))
}

func getDefaultMemory(totalMemory, minMemory, maxMemory uint64) uint64 {
	halfTotal := totalMemory / 2
	if halfTotal <= minMemory {
		return minMemory
	} else if halfTotal >= maxMemory {
		return maxMemory
	}
	return halfTotal
}

func stripWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}
