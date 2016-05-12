package http_test

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/glestaris/uberlist-server/backend"
	"github.com/glestaris/uberlist-server/http"
	"github.com/glestaris/uberlist-server/http/fakes"

	. "github.com/glestaris/uberlist-server/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Roundtrip", func() {
	var (
		fakeBackend        *fakes.FakeBackend
		fakeMessageEncoder *fakes.FakeMessageEncoder
		logger             *logrus.Logger

		client *http.Client
		server *http.Server

		endpointAddr net.IP
		endpointPort uint16
		serverCh     chan struct{}
	)

	BeforeEach(func() {
		var err error

		fakeMessageEncoder = new(fakes.FakeMessageEncoder)
		fakeBackend = new(fakes.FakeBackend)
		logger, _ = test.NewNullLogger()

		server = http.NewServer(logger, fakeBackend, fakeMessageEncoder)
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

		client, err = http.NewClient(
			logger,
			fmt.Sprintf("ws://%s:%d", endpointAddr, endpointPort),
			fakeMessageEncoder,
		)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(server.Close(logger)).To(Succeed())
		Eventually(serverCh).Should(BeClosed())
	})

	Describe("Server", func() {
		It("receives task messages", func() {
			msg := &fakes.FakeMessage{
				Id: "banana",
			}

			Expect(client.SendMessage(logger, msg)).To(Succeed())

			Eventually(func() []backend.Message {
				msgs := []backend.Message{}
				for _, e := range fakeBackend.Messages {
					msgs = append(msgs, e.Msg)
				}
				return msgs
			}).Should(Equal([]backend.Message{msg}))
		})

		Context("when applying the message fails", func() {
			BeforeEach(func() {
				fakeBackend.SendMessageReturns = errors.New("Hello world")
			})

			It("responds with the error", func() {
				msg := &fakes.FakeMessage{}

				Expect(client.SendMessage(logger, msg)).To(MatchError("Hello world"))
			})
		})

		Describe("Close", func() {
			It("stops receiving connections", func() {
				Expect(server.Close(logger)).To(Succeed())

				oldDialTimeout := http.DialTimeout
				http.DialTimeout = time.Millisecond * 200
				defer func() {
					http.DialTimeout = oldDialTimeout
				}()

				_, err := http.NewClient(
					logger,
					fmt.Sprintf("ws://%s:%d", endpointAddr, endpointPort),
					fakeMessageEncoder,
				)
				Expect(err).To(HaveOccurred())
			})

			It("closes existing connections", func() {
				Expect(server.Close(logger)).To(Succeed())

				msg := &fakes.FakeMessage{}
				Expect(client.SendMessage(logger, msg)).To(
					MatchErrorType(http.ConnectionDroppedError("")),
				)
			})
		})
	})

	Describe("Client", func() {
		Describe("SendMessage", func() {
			It("waits for the server to process the message", func() {
				msg := &fakes.FakeMessage{}

				Expect(fakeBackend.Messages).To(HaveLen(0))
				Expect(client.SendMessage(logger, msg)).To(Succeed())
				Expect(fakeBackend.Messages).To(HaveLen(1))
			})

			Context("when encoding the message fails", func() {
				var encodeError error

				BeforeEach(func() {
					encodeError = errors.New("My clementine went bad")
					fakeMessageEncoder.EncodeError = encodeError
				})

				It("returns the error", func() {
					msg := &fakes.FakeMessage{}
					Expect(client.SendMessage(logger, msg)).To(MatchError(encodeError))
				})
			})
		})
	})
})
