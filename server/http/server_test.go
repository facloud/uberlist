package http_test

import (
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/glestaris/uberlist-server/http"
	"github.com/glestaris/uberlist-server/http/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	var (
		fakeBackend        *fakes.FakeBackend
		fakeMessageEncoder *fakes.FakeMessageEncoder

		logger *logrus.Logger

		server *http.Server
	)

	BeforeEach(func() {
		fakeMessageEncoder = new(fakes.FakeMessageEncoder)
		fakeBackend = new(fakes.FakeBackend)

		logger, _ = test.NewNullLogger()

		server = http.NewServer(logger, fakeBackend, fakeMessageEncoder)
	})

	Context("when the server is serving", func() {
		var (
			serverCh     chan struct{}
			endpointAddr net.IP
			endpointPort uint16
		)

		BeforeEach(func() {
			serverCh = make(chan struct{})

			endpointAddr = net.ParseIP("127.0.0.1")
			endpointPort = uint16(8000 + GinkgoParallelNode())
			go func() {
				defer GinkgoRecover()

				Expect(server.Serve(logger, endpointAddr, endpointPort)).To(Succeed())
				close(serverCh)
			}()
			Eventually(func() bool {
				return server.IsListening(logger)
			}).Should(BeTrue())
		})

		AfterEach(func() {
			Expect(server.Close(logger)).To(Succeed())
			Eventually(serverCh).Should(BeClosed())
		})

		Describe("Serve", func() {
			It("returns an error", func() {
				Expect(
					server.Serve(logger, endpointAddr, endpointPort),
				).To(HaveOccurred())
			})
		})

		Describe("Close", func() {
			It("succeeds when closing the same server twice", func() {
				Expect(server.Close(logger)).To(Succeed())
				Expect(server.Close(logger)).To(Succeed())
			})

			It("stops the listener", func() {
				Expect(server.IsListening(logger)).To(BeTrue())
				Expect(server.Close(logger)).To(Succeed())
				Expect(server.IsListening(logger)).To(BeFalse())
			})
		})
	})

	Describe("Close", func() {
		It("succeeds when the server is not serving", func() {
			Expect(server.Close(logger)).To(Succeed())
		})
	})
})
