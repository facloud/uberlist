package http_test

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/glestaris/uberlist-server/http"
	"github.com/glestaris/uberlist-server/http/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var (
		fakeMessageEncoder *fakes.FakeMessageEncoder
		logger             *logrus.Logger
	)

	BeforeEach(func() {
		fakeMessageEncoder = new(fakes.FakeMessageEncoder)
	})

	Describe("NewClient", func() {
		It("returns an error when the server is not listening", func() {
			oldDialTimeout := http.DialTimeout
			http.DialTimeout = time.Millisecond * 200
			defer func() {
				http.DialTimeout = oldDialTimeout
			}()

			_, err := http.NewClient(
				logger, "ws://10.10.10.10:1555", fakeMessageEncoder,
			)
			Expect(err).NotTo(Succeed())
		})
	})
})
