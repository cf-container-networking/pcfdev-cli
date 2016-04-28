package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var context struct {
	cfHome string
}

var _ = BeforeSuite(func() {
	ifconfig := exec.Command("ifconfig")
	session, err := gexec.Start(ifconfig, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Expect(session.Wait().Out.Contents()).NotTo(ContainSubstring("192.168.11.1"))

	uninstallCommand := exec.Command("cf", "uninstall-plugin", "pcfdev")
	session, err = gexec.Start(uninstallCommand, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "10s").Should(gexec.Exit())

	pluginPath, err := gexec.Build("github.com/pivotal-cf/pcfdev-cli")
	Expect(err).NotTo(HaveOccurred())
	installCommand := exec.Command("cf", "install-plugin", "-f", pluginPath)
	session, err = gexec.Start(installCommand, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "1m").Should(gexec.Exit(0))
})

var _ = AfterSuite(func() {
	uninstallCommand := exec.Command("cf", "uninstall-plugin", "pcfdev")
	session, err := gexec.Start(uninstallCommand, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "10s").Should(gexec.Exit(0))
})

var _ = Describe("pcfdev", func() {
	Context("pivnet api token is set in environment", func() {
		BeforeEach(func() {
			err := os.RemoveAll(filepath.Join(os.Getenv("HOME"), ".pcfdev"))
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			output, err := exec.Command("VBoxManage", "showvminfo", "pcfdev-2016-03-29_1728", "--machinereadable").Output()
			if err != nil {
				return
			}

			regex := regexp.MustCompile(`hostonlyadapter2="(.*)"`)
			vboxnet := regex.FindStringSubmatch(string(output))[1]

			exec.Command("VBoxManage", "controlvm", "pcfdev-2016-03-29_1728", "poweroff").Run()
			exec.Command("VBoxManage", "unregistervm", "pcfdev-2016-03-29_1728", "--delete").Run()
			exec.Command("VBoxManage", "hostonlyif", "remove", vboxnet).Run()
		})

		It("should start, stop, and destroy a virtualbox instance", func() {
			pcfdevCommand := exec.Command("cf", "dev", "start")
			session, err := gexec.Start(pcfdevCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "1h").Should(gexec.Exit(0))
			Expect(session).To(gbytes.Say("PCF Dev is now running"))
			Expect(isVMRunning()).To(BeTrue())

			// rerunning start has no effect
			restartCommand := exec.Command("cf", "dev", "start")
			session, err = gexec.Start(restartCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "3m").Should(gexec.Exit(0))
			Expect(session).To(gbytes.Say("PCF Dev is running"))
			Expect(isVMRunning()).To(BeTrue())

			Eventually(cf("login", "-a", "api.local.pcfdev.io", "-u", "admin", "-p", "admin", "--skip-ssl-validation"), 5*time.Second).Should(gexec.Exit(0))
			Eventually(cf("push", "app", "-o", "cloudfoundry/lattice-app"), 2*time.Minute).Should(gexec.Exit(0))

			pcfdevCommand = exec.Command("cf", "dev", "stop")
			session, err = gexec.Start(pcfdevCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "10m").Should(gexec.Exit(0))
			Expect(session).To(gbytes.Say("PCF Dev is now stopped"))
			Expect(isVMRunning()).NotTo(BeTrue())

			pcfdevCommand = exec.Command("cf", "dev", "destroy")
			session, err = gexec.Start(pcfdevCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "10m").Should(gexec.Exit(0))
			Expect(session).To(gbytes.Say("PCF Dev VM has been destroyed"))

			// rerunning destroy has no effect
			redestroyCommand := exec.Command("cf", "dev", "destroy")
			session, err = gexec.Start(redestroyCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "10m").Should(gexec.Exit(0))
			Expect(session).To(gbytes.Say("PCF Dev VM has not been created"))

			// can start and push app after running destroy
			pcfdevCommand = exec.Command("cf", "dev", "start")
			session, err = gexec.Start(pcfdevCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "1h").Should(gexec.Exit(0))
			Expect(session).To(gbytes.Say("PCF Dev is now running"))
			Expect(isVMRunning()).To(BeTrue())

			Eventually(cf("login", "-a", "api.local.pcfdev.io", "-u", "admin", "-p", "admin", "--skip-ssl-validation"), 5*time.Second).Should(gexec.Exit(0))
			Eventually(cf("push", "app", "-o", "cloudfoundry/lattice-app"), 2*time.Minute).Should(gexec.Exit(0))
		})

		It("should respond to pcfdev alias", func() {
			pcfdevCommand := exec.Command("cf", "pcfdev")
			session, err := gexec.Start(pcfdevCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
			Expect(session).To(gbytes.Say(`Usage: cf dev download\|start\|status\|stop\|destroy`))
		})

		It("should download a VM without importing it", func() {
			pcfdevCommand := exec.Command("cf", "dev", "download")
			session, err := gexec.Start(pcfdevCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "1h").Should(gexec.Exit(0))

			_, err = os.Stat(filepath.Join(os.Getenv("HOME"), ".pcfdev", "pcfdev.ova"))
			Expect(err).NotTo(HaveOccurred())

			listVmsCommand := exec.Command("VBoxManage", "list", "vms")
			session, err = gexec.Start(listVmsCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))
			Expect(session).NotTo(gbytes.Say("pcfdev-2016-03-29_1728"))

			pcfdevCommand = exec.Command("cf", "dev", "download")
			session, err = gexec.Start(pcfdevCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "3m").Should(gexec.Exit(0))
		})
	})
})

func loadEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		Fail("missing "+name, 1)
	}
	return value
}

func cf(args ...string) *gexec.Session {
	command := exec.Command("cf", args...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return session
}

func isVMRunning() bool {
	vmStatus, err := exec.Command("VBoxManage", "showvminfo", "pcfdev-2016-03-29_1728", "--machinereadable").Output()
	Expect(err).NotTo(HaveOccurred())
	return strings.Contains(string(vmStatus), `VMState="running"`)
}
