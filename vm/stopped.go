package vm

import (
	"fmt"
	"os"
	"time"
)

type Stopped struct {
	Name    string
	Domain  string
	IP      string
	SSHPort string

	VBox                VBox
	SSH                 SSH
	UI                  UI
	RequirementsChecker RequirementsChecker
}

func (s *Stopped) Stop() error {
	s.UI.Say("PCF Dev is stopped")
	return nil
}

func (s *Stopped) Start() error {
	if err := s.RequirementsChecker.Check(); err != nil {
		if !s.UI.Confirm("Less than 3 GB of memory detected, continue (y/N): ") {
			s.UI.Say("Exiting...")
			return nil
		}
	}

	s.UI.Say("Starting VM...")
	if err := s.VBox.StartVM(s.Name, s.IP, s.SSHPort, s.Domain); err != nil {
		return &StartVMError{err}
	}

	s.UI.Say("Provisioning VM...")
	provisionCommand := fmt.Sprintf("sudo /var/pcfdev/run %s %s '$2a$04$EpJtIJ8w6hfCwbKYBkn3t.GCY18Pk6s7yN66y37fSJlLuDuMkdHtS'", s.Domain, s.IP)
	if err := s.SSH.RunSSHCommand(provisionCommand, s.SSHPort, 2*time.Minute, os.Stdout, os.Stderr); err != nil {
		return &ProvisionVMError{err}
	}

	s.UI.Say("PCF Dev is now running")
	return nil
}

func (s *Stopped) Status() {
	s.UI.Say("Stopped")
}

func (s *Stopped) Destroy() error {
	return s.VBox.DestroyVM(s.Name)
}

func (s *Stopped) Suspend() error {
	s.UI.Say("Your VM is currently stopped and cannot be suspended.")
	return nil
}

func (s *Stopped) Resume() error {
	s.UI.Say("Your VM is currently stopped. Only a suspended VM can be resumed.")
	return nil
}