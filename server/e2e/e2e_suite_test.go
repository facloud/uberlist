package client_test

import (
	"fmt"
	"os/exec"
	"syscall"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var (
	serverBinPath  string
	serverSess     *gexec.Session
	serverEndpoint string
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)

	SynchronizedBeforeSuite(func() []byte {
		serverBinPath, err := gexec.Build(
			"github.com/glestaris/uberlist-server/cmd/uberlist-server",
		)
		Expect(err).NotTo(HaveOccurred())

		return []byte(serverBinPath)
	}, func(data []byte) {
		serverBinPath = string(data)
	})

	BeforeEach(func() {
		port := 5000 + GinkgoParallelNode()
		serverEndpoint = fmt.Sprintf("0.0.0.0:%d", port)
		serverSess = startServer("address", serverEndpoint)
	})

	AfterEach(func() {
		stopServer(serverSess)
	})

	RunSpecs(t, "Client Suite")
}

func startServer(args ...string) *gexec.Session {
	sess, err := gexec.Start(
		exec.Command(serverBinPath, args...),
		GinkgoWriter,
		GinkgoWriter,
	)
	Expect(err).NotTo(HaveOccurred())

	return sess
}

func stopServer(serverSession *gexec.Session) {
	serverSession.Signal(syscall.SIGTERM)
}
