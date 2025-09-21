package infrastructure_test

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/hibiken/asynq"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	"src/internal/config"
	"src/internal/modules/shared/domain/events"
	"src/internal/pkg/eventbus"
	"src/internal/pkg/taskqueue"
)

// Test job type
const TestJobType = "test:job"

type TestJobPayload struct {
	Message string `json:"message"`
}

var (
	redisContainer testcontainers.Container
	testConfig     *config.Config
)

func mustStartRedisContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	ctx := context.Background()

	redisContainer, err := redis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithOccurrence(1).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	host, err := redisContainer.Host(ctx)
	if err != nil {
		return redisContainer.Terminate, err
	}

	port, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return redisContainer.Terminate, err
	}

	testConfig = &config.Config{
		Port: 8080, // Default for tests
	}
	testConfig.Database.Host = "localhost"
	testConfig.Database.Port = "5432"
	testConfig.Database.Username = "postgres"
	testConfig.Database.Password = "password"
	testConfig.Database.Database = "test_db"
	testConfig.Database.Schema = "public"

	testConfig.Redis.Host = host
	testConfig.Redis.Port = port.Port()
	testConfig.Redis.Password = ""

	// Set test config
	config.SetForTests(testConfig)

	return redisContainer.Terminate, nil
}

var _ = BeforeSuite(func() {
	var err error
	teardown, err = mustStartRedisContainer()
	Expect(err).ToNot(HaveOccurred())
})

var teardown func(context.Context, ...testcontainers.TerminateOption) error

var _ = AfterSuite(func() {
	if teardown != nil {
		err := teardown(context.Background())
		if err != nil {
			fmt.Printf("could not teardown redis container: %v\n", err)
		}
	}
})

var _ = Describe("Infrastructure", func() {

	It("should publish and consume events successfully", func() {
		// Create logger
		logger := watermill.NewStdLogger(false, false)

		// Create publisher
		publisher, err := eventbus.NewPublisher(logger)
		Expect(err).ToNot(HaveOccurred())

		// Create router
		router, err := eventbus.NewRouter(logger)
		Expect(err).ToNot(HaveOccurred())

		// Channel to capture consumed events
		eventsReceived := make(chan *message.Message, 1)

		// Create a simple test handler function
		handlerFunc := func(msg *message.Message) error {
			eventsReceived <- msg
			return nil
		}

		// Add a test handler - simplified for testing
		router.AddNoPublisherHandler(
			"test_handler",
			events.NotionWebhookReceivedTopic,
			publisher.(message.Subscriber), // GoChannel implements both Publisher and Subscriber
			handlerFunc,
		)

		// Start router in background
		go func() {
			defer GinkgoRecover()
			err := router.Run(context.Background())
			if err != nil {
				Fail(fmt.Sprintf("Router failed: %v", err))
			}
		}()

		// Give router time to start
		time.Sleep(100 * time.Millisecond)

		// Publish test event
		testPayload := []byte(`{"test": "data"}`)
		msg := message.NewMessage(watermill.NewUUID(), testPayload)

		err = publisher.Publish(events.NotionWebhookReceivedTopic, msg)
		Expect(err).ToNot(HaveOccurred())

		// Wait for event to be consumed
		select {
		case receivedMsg := <-eventsReceived:
			Expect(string(receivedMsg.Payload)).To(Equal(string(testPayload)))
		case <-time.After(5 * time.Second):
			Fail("Event was not consumed within timeout")
		}
	})

	It("should enqueue and process jobs successfully", func() {
		// Create Redis client for Asynq
		redisOpt := asynq.RedisClientOpt{
			Addr:     testConfig.RedisURL(),
			Password: testConfig.Redis.Password,
		}

		// Create Asynq client
		client := taskqueue.NewClient(redisOpt)

		// Channel to capture job results
		jobProcessed := make(chan string, 1)

		// Create server with handler
		mux := asynq.NewServeMux()
		mux.HandleFunc(TestJobType, func(ctx context.Context, t *asynq.Task) error {
			var payload TestJobPayload
			if err := json.Unmarshal(t.Payload(), &payload); err != nil {
				return err
			}
			jobProcessed <- payload.Message
			return nil
		})

		server := asynq.NewServer(redisOpt, asynq.Config{
			Concurrency: 1,
			Queues: map[string]int{
				"default": 1,
			},
		})

		// Start server in background
		go func() {
			defer GinkgoRecover()
			err := server.Start(mux)
			if err != nil {
				Fail(fmt.Sprintf("Server failed: %v", err))
			}
		}()

		// Give server time to start
		time.Sleep(100 * time.Millisecond)

		// Enqueue test job
		payload := TestJobPayload{Message: "test job processed"}
		payloadBytes, err := json.Marshal(payload)
		Expect(err).ToNot(HaveOccurred())
		info, err := client.Enqueue(asynq.NewTask(TestJobType, payloadBytes))
		Expect(err).ToNot(HaveOccurred())
		Expect(info).ToNot(BeNil())

		// Wait for job to be processed
		select {
		case result := <-jobProcessed:
			Expect(result).To(Equal("test job processed"))
		case <-time.After(10 * time.Second):
			Fail("Job was not processed within timeout")
		}

		// Clean up
		server.Shutdown()
	})
})
