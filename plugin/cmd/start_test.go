package cmd_test

import (
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pcfdev-cli/config"
	"github.com/pivotal-cf/pcfdev-cli/plugin/cmd"
	"github.com/pivotal-cf/pcfdev-cli/plugin/cmd/mocks"
	"github.com/pivotal-cf/pcfdev-cli/vm"
	vmMocks "github.com/pivotal-cf/pcfdev-cli/vm/mocks"
)

var _ = Describe("StartCmd", func() {
	var (
		startCmd      *cmd.StartCmd
		mockCtrl      *gomock.Controller
		mockVMBuilder *mocks.MockVMBuilder
		mockVBox      *mocks.MockVBox
		mockVM        *vmMocks.MockVM
		mockCmd       *mocks.MockCmd
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockVMBuilder = mocks.NewMockVMBuilder(mockCtrl)
		mockVBox = mocks.NewMockVBox(mockCtrl)
		mockVM = vmMocks.NewMockVM(mockCtrl)
		mockCmd = mocks.NewMockCmd(mockCtrl)
		startCmd = &cmd.StartCmd{
			VBox:      mockVBox,
			VMBuilder: mockVMBuilder,
			Config: &config.Config{
				DefaultVMName: "some-default-vm-name",
			},
			Opts:        &vm.StartOpts{},
			DownloadCmd: mockCmd,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Parse", func() {
		Context("when flags are passed", func() {
			It("should set start options", func() {
				Expect(startCmd.Parse([]string{
					"-c", "2",
					"-m", "3456",
					"-n",
					"-o", "some-ova-path",
					"-r", "some-private-registry,some-other-private-registry",
					"-s", "some-service,some-other-service",
				})).To(Succeed())

				Expect(startCmd.Opts.CPUs).To(Equal(2))
				Expect(startCmd.Opts.Memory).To(Equal(uint64(3456)))
				Expect(startCmd.Opts.NoProvision).To(BeTrue())
				Expect(startCmd.Opts.OVAPath).To(Equal("some-ova-path"))
				Expect(startCmd.Opts.Registries).To(Equal("some-private-registry,some-other-private-registry"))
				Expect(startCmd.Opts.Services).To(Equal("some-service,some-other-service"))
			})
		})

		Context("when no flags are passed", func() {
			It("should set start options", func() {
				Expect(startCmd.Parse([]string{})).To(Succeed())
				Expect(startCmd.Opts.CPUs).To(Equal(0))
				Expect(startCmd.Opts.Memory).To(Equal(uint64(0)))
				Expect(startCmd.Opts.NoProvision).To(BeFalse())
				Expect(startCmd.Opts.OVAPath).To(Equal(""))
				Expect(startCmd.Opts.Registries).To(Equal(""))
				Expect(startCmd.Opts.Services).To(Equal(""))
			})
		})

		Context("when an unknown flag is passed", func() {
			It("should return an error", func() {
				Expect(startCmd.Parse(
					[]string{"-b", "some-bad-flag"})).NotTo(Succeed())
			})
		})

		Context("when an unknown argument is passed", func() {
			It("should return an error", func() {
				Expect(startCmd.Parse(
					[]string{"some-bad-argument"})).NotTo(Succeed())
			})
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			startCmd.Parse([]string{})
		})

		Context("when starting the default ova", func() {
			It("should validate start options and starts the VM", func() {
				startOpts := &vm.StartOpts{
					Memory: uint64(3456),
					CPUs:   2,
				}
				startCmd.Opts = startOpts
				gomock.InOrder(
					mockVBox.EXPECT().GetVMName().Return("", nil),
					mockVMBuilder.EXPECT().VM("some-default-vm-name").Return(mockVM, nil),
					mockVM.EXPECT().VerifyStartOpts(startOpts),
					mockCmd.EXPECT().Run(),
					mockVM.EXPECT().Start(startOpts),
				)

				Expect(startCmd.Run()).To(Succeed())
			})

			Context("when there is an old vm present", func() {
				It("should tell the user to destroy pcfdev", func() {
					mockVBox.EXPECT().GetVMName().Return("some-old-vm-name", nil)

					Expect(startCmd.Run()).To(MatchError("old version of PCF Dev already running, please run `cf dev destroy` to continue"))
				})
			})

			Context("when there is an error getting the VM name", func() {
				It("should return the error", func() {
					mockVBox.EXPECT().GetVMName().Return("", errors.New("some-error"))

					Expect(startCmd.Run()).To(MatchError("some-error"))
				})
			})

			Context("when it fails to get VM", func() {
				It("should return an error", func() {
					gomock.InOrder(
						mockVBox.EXPECT().GetVMName().Return("", nil),
						mockVMBuilder.EXPECT().VM("some-default-vm-name").Return(nil, errors.New("some-error")),
					)
					Expect(startCmd.Run()).To(MatchError("some-error"))
				})
			})

			Context("when verifying start options fails", func() {
				It("should return an error", func() {
					gomock.InOrder(
						mockVBox.EXPECT().GetVMName().Return("", nil),
						mockVMBuilder.EXPECT().VM("some-default-vm-name").Return(mockVM, nil),
						mockVM.EXPECT().VerifyStartOpts(&vm.StartOpts{}).Return(errors.New("some-error")),
					)

					Expect(startCmd.Run()).To(MatchError("some-error"))
				})
			})

			Context("when the OVA fails to download", func() {
				It("should print an error message", func() {
					gomock.InOrder(
						mockVBox.EXPECT().GetVMName().Return("", nil),
						mockVMBuilder.EXPECT().VM("some-default-vm-name").Return(mockVM, nil),
						mockVM.EXPECT().VerifyStartOpts(&vm.StartOpts{}),
						mockCmd.EXPECT().Run().Return(errors.New("some-error")),
					)

					Expect(startCmd.Run()).To(MatchError("some-error"))
				})
			})

			Context("when it fails to start VM", func() {
				It("should return an error", func() {
					gomock.InOrder(
						mockVBox.EXPECT().GetVMName().Return("", nil),
						mockVMBuilder.EXPECT().VM("some-default-vm-name").Return(mockVM, nil),
						mockVM.EXPECT().VerifyStartOpts(&vm.StartOpts{}),
						mockCmd.EXPECT().Run(),
						mockVM.EXPECT().Start(&vm.StartOpts{}).Return(errors.New("some-error")),
					)

					Expect(startCmd.Run()).To(MatchError("some-error"))
				})
			})
		})

		Context("when starting a custom ova", func() {
			It("should start the custom ova", func() {
				startOpts := &vm.StartOpts{
					OVAPath: "some-custom-ova",
				}
				startCmd.Opts = startOpts
				gomock.InOrder(
					mockVBox.EXPECT().GetVMName().Return("", nil),
					mockVMBuilder.EXPECT().VM("pcfdev-custom").Return(mockVM, nil),
					mockVM.EXPECT().VerifyStartOpts(startOpts),
					mockVM.EXPECT().Start(startOpts),
				)

				Expect(startCmd.Run()).To(Succeed())
			})

			Context("when the custom VM is already present and OVAPath is not set", func() {
				It("should start the custom VM", func() {
					gomock.InOrder(
						mockVBox.EXPECT().GetVMName().Return("pcfdev-custom", nil),
						mockVMBuilder.EXPECT().VM("pcfdev-custom").Return(mockVM, nil),
						mockVM.EXPECT().VerifyStartOpts(&vm.StartOpts{}),
						mockVM.EXPECT().Start(&vm.StartOpts{}),
					)
					Expect(startCmd.Run()).To(Succeed())
				})
			})

			Context("when the custom VM is already present and OVAPath is set", func() {
				It("should start the custom OVA", func() {
					startOpts := &vm.StartOpts{
						OVAPath: "some-custom-ova",
					}
					startCmd.Opts = startOpts
					gomock.InOrder(
						mockVBox.EXPECT().GetVMName().Return("pcfdev-custom", nil),
						mockVMBuilder.EXPECT().VM("pcfdev-custom").Return(mockVM, nil),
						mockVM.EXPECT().VerifyStartOpts(startOpts),
						mockVM.EXPECT().Start(startOpts),
					)
					Expect(startCmd.Run()).To(Succeed())
				})
			})

			Context("when the default VM is present", func() {
				It("should return an error", func() {
					startOpts := &vm.StartOpts{
						OVAPath: "some-custom-ova",
					}
					startCmd.Opts = startOpts
					mockVBox.EXPECT().GetVMName().Return("some-default-vm-name", nil)
					Expect(startCmd.Run()).To(MatchError("you must destroy your existing VM to use a custom OVA"))
				})
			})

			Context("when an old VM is present", func() {
				It("should return an error", func() {
					startOpts := &vm.StartOpts{
						OVAPath: "some-custom-ova",
					}
					startCmd.Opts = startOpts
					mockVBox.EXPECT().GetVMName().Return("some-old-vm-name", nil)
					Expect(startCmd.Run()).To(MatchError("you must destroy your existing VM to use a custom OVA"))
				})
			})
		})

		Context("when the provision option is specified", func() {
			It("should provision the VM without starting it", func() {
				startCmd.Parse([]string{"-p"})

				gomock.InOrder(
					mockVBox.EXPECT().GetVMName().Return("", nil),
					mockVMBuilder.EXPECT().VM("some-default-vm-name").Return(mockVM, nil),
					mockVM.EXPECT().Provision(),
				)

				Expect(startCmd.Run()).To(Succeed())
			})
		})

		Context("when provisioning fails", func() {
			It("return an error", func() {
				startCmd.Parse([]string{"-p"})

				gomock.InOrder(
					mockVBox.EXPECT().GetVMName().Return("", nil),
					mockVMBuilder.EXPECT().VM("some-default-vm-name").Return(mockVM, nil),
					mockVM.EXPECT().Provision().Return(errors.New("some-error")),
				)

				Expect(startCmd.Run()).To(MatchError("some-error"))
			})
		})
	})
})