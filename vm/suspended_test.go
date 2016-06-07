package vm_test

import (
	"errors"

	"github.com/golang/mock/gomock"
	conf "github.com/pivotal-cf/pcfdev-cli/config"
	"github.com/pivotal-cf/pcfdev-cli/vm"
	"github.com/pivotal-cf/pcfdev-cli/vm/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Suspended", func() {
	var (
		mockCtrl    *gomock.Controller
		mockUI      *mocks.MockUI
		mockVBox    *mocks.MockVBox
		suspendedVM vm.Suspended
		config      *conf.VMConfig
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUI = mocks.NewMockUI(mockCtrl)
		mockVBox = mocks.NewMockVBox(mockCtrl)
		config = &conf.VMConfig{}

		suspendedVM = vm.Suspended{
			Name:    "some-vm",
			Domain:  "some-domain",
			IP:      "some-ip",
			SSHPort: "some-port",
			Config:  config,

			VBox: mockVBox,
			UI:   mockUI,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Suspend", func() {
		It("should say a message", func() {
			mockUI.EXPECT().Say("Your VM is suspended.")
			Expect(suspendedVM.Suspend()).To(Succeed())
		})
	})

	Describe("Stop", func() {
		It("should say a message", func() {
			mockUI.EXPECT().Say("Your VM is currently suspended. You must resume your VM with `cf dev resume` to shut it down.")
			Expect(suspendedVM.Stop()).To(Succeed())
		})
	})

	Describe("Start", func() {
		It("should start vm", func() {
			gomock.InOrder(
				mockUI.EXPECT().Say("Resuming VM..."),
				mockVBox.EXPECT().ResumeVM("some-vm").Return(nil),
			)

			Expect(suspendedVM.Start()).To(Succeed())
		})

		Context("when starting the vm fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockUI.EXPECT().Say("Resuming VM..."),
					mockVBox.EXPECT().ResumeVM("some-vm").Return(errors.New("some-error")),
				)

				Expect(suspendedVM.Start()).To(MatchError("failed to resume VM: some-error"))
			})
		})
	})

	Describe("Resume", func() {
		It("should start vm", func() {
			gomock.InOrder(
				mockUI.EXPECT().Say("Resuming VM..."),
				mockVBox.EXPECT().ResumeVM("some-vm").Return(nil),
			)

			Expect(suspendedVM.Resume()).To(Succeed())
		})

		Context("when starting the vm fails", func() {
			It("should return an error", func() {
				gomock.InOrder(
					mockUI.EXPECT().Say("Resuming VM..."),
					mockVBox.EXPECT().ResumeVM("some-vm").Return(errors.New("some-error")),
				)

				Expect(suspendedVM.Resume()).To(MatchError("failed to resume VM: some-error"))
			})
		})
	})

	Describe("Status", func() {
		It("should say Suspended", func() {
			mockUI.EXPECT().Say("Suspended")

			suspendedVM.Status()
		})
	})

	Describe("Destroy", func() {
		It("should destroy the vm", func() {
			mockVBox.EXPECT().DestroyVM("some-vm").Return(nil)

			Expect(suspendedVM.Destroy()).To(Succeed())
		})
	})

	Describe("Config", func() {
		It("should return the config", func() {
			Expect(suspendedVM.GetConfig()).To(BeIdenticalTo(config))
		})
	})
})
