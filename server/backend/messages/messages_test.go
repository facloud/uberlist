package messages_test

import (
	"github.com/glestaris/uberlist-server"
	"github.com/glestaris/uberlist-server/backend"
	"github.com/glestaris/uberlist-server/backend/messages"
	"github.com/glestaris/uberlist-server/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Messages", func() {
	var store backend.Store

	BeforeEach(func() {
		store = storage.NewLocalStore()
	})

	Describe("NewMessage", func() {
		It("successfully adds a task", func() {
			task := uberlist.Task{Title: "Hello world"}
			msg := messages.NewTaskMessage{NewTask: task}

			Expect(msg.Apply(store)).To(Succeed())

			tasks, err := store.OrderedTasks()
			Expect(err).NotTo(HaveOccurred())
			Expect(tasks).To(HaveLen(1))
			Expect(tasks[0].Title).To(Equal(task.Title))
		})
	})

	Describe("UpdateMessage", func() {
		It("successfully updates a task", func() {
			task := uberlist.Task{Title: "Hello world"}
			addMsg := messages.NewTaskMessage{NewTask: task}

			Expect(addMsg.Apply(store)).To(Succeed())

			tasks, err := store.OrderedTasks()
			Expect(err).NotTo(HaveOccurred())
			storedTask := tasks[0]

			storedTask.Title = "New hello world"
			updateMsg := messages.UpdateTaskMessage{
				Task: storedTask,
			}

			Expect(updateMsg.Apply(store)).To(Succeed())

			tasks, err = store.OrderedTasks()
			Expect(err).NotTo(HaveOccurred())
			Expect(tasks).To(HaveLen(1))
			Expect(tasks[0].Title).To(Equal(storedTask.Title))
		})
	})
})
