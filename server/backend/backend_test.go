package backend_test

import (
	"github.com/glestaris/uberlist-server"
	"github.com/glestaris/uberlist-server/backend"
	"github.com/glestaris/uberlist-server/backend/fakes"
	"github.com/glestaris/uberlist-server/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backend", func() {
	var (
		store backend.Store
		bcknd *backend.Backend
	)

	BeforeEach(func() {
		store = storage.NewLocalStore()
		bcknd = backend.NewBackend(store)
	})

	Describe("Message feed", func() {
		It("applies received message", func() {
			msg := &fakes.FakeMessage{}
			sender := backend.Client{ID: "1"}

			Expect(bcknd.SendMessage(sender, msg)).To(Succeed())

			Expect(msg.Applied).To(BeTrue())
			Expect(msg.AppliedToStore).To(Equal(store))
		})

		It("updates subscribed clients", func() {
			msg := &fakes.FakeMessage{}
			sender := backend.Client{ID: "1"}
			recipient := backend.Client{ID: "2"}

			recipientChan, err := bcknd.Subscribe(recipient)
			Expect(err).NotTo(HaveOccurred())

			Expect(bcknd.SendMessage(sender, msg)).To(Succeed())

			Eventually(recipientChan).Should(Receive())
		})

		It("does not update the calling client", func() {
			msg := &fakes.FakeMessage{}
			sender := backend.Client{ID: "1"}

			senderChan, err := bcknd.Subscribe(sender)
			Expect(err).NotTo(HaveOccurred())

			Expect(bcknd.SendMessage(sender, msg)).To(Succeed())

			Consistently(senderChan).ShouldNot(Receive())
		})

		It("closes the channel of an unsubscribed client", func() {
			client := backend.Client{ID: "1"}

			ch, err := bcknd.Subscribe(client)
			Expect(err).NotTo(HaveOccurred())

			Expect(bcknd.Unsubscribe(client)).To(Succeed())

			Eventually(ch).Should(BeClosed())
		})

		It("returns an error when trying to unsubscribe a non-existing client", func() {
			Expect(bcknd.Unsubscribe(backend.Client{ID: "banana"})).NotTo(Succeed())
		})
	})

	Describe("Tasks", func() {
		var tasks []uberlist.Task

		BeforeEach(func() {
			var err error

			tasks = []uberlist.Task{
				uberlist.Task{ID: 1},
				uberlist.Task{ID: 2},
			}

			for _, t := range tasks {
				t.ID, err = store.AddTask(t)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("returns the tasks", func() {
			retTasks, err := bcknd.Tasks()
			Expect(err).NotTo(HaveOccurred())
			Expect(retTasks).To(Equal(tasks))
		})
	})
})
