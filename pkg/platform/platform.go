package platform

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// OS represents the operating system
type OS string

const (
	OSMacOS   OS = "macos"
	OSLinux   OS = "linux"
	OSUnknown OS = "unknown"
)

// Arch represents the CPU architecture
type Arch string

const (
	ArchAMD64   Arch = "amd64"
	ArchARM64   Arch = "arm64"
	ArchUnknown Arch = "unknown"
)

// Platform contains information about the current system
type Platform struct {
	OS           OS
	Arch         Arch
	OSVersion    string
	IsAppleSilicon bool
}

// Detect returns information about the current platform
func Detect() (*Platform, error) {
	p := &Platform{
		OS:   detectOS(),
		Arch: detectArch(),
	}

	if p.OS == OSUnknown {
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if p.Arch == ArchUnknown {
		return nil, fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}

	p.IsAppleSilicon = p.OS == OSMacOS && p.Arch == ArchARM64
	p.OSVersion = detectOSVersion(p.OS)

	return p, nil
}

func detectOS() OS {
	switch runtime.GOOS {
	case "darwin":
		return OSMacOS
	case "linux":
		return OSLinux
	default:
		return OSUnknown
	}
}

func detectArch() Arch {
	switch runtime.GOARCH {
	case "amd64":
		return ArchAMD64
	case "arm64":
		return ArchARM64
	default:
		return ArchUnknown
	}
}

func detectOSVersion(os OS) string {
	switch os {
	case OSMacOS:
		out, err := exec.Command("sw_vers", "-productVersion").Output()
		if err != nil {
			return "unknown"
		}
		return strings.TrimSpace(string(out))
	case OSLinux:
		out, err := exec.Command("uname", "-r").Output()
		if err != nil {
			return "unknown"
		}
		return strings.TrimSpace(string(out))
	default:
		return "unknown"
	}
}

// String returns a human-readable description of the platform
func (p *Platform) String() string {
	osName := string(p.OS)
	if p.OS == OSMacOS {
		osName = "macOS"
		if p.IsAppleSilicon {
			osName += " (Apple Silicon)"
		} else {
			osName += " (Intel)"
		}
	} else if p.OS == OSLinux {
		osName = "Linux"
	}
	return fmt.Sprintf("%s %s (%s)", osName, p.OSVersion, p.Arch)
}
