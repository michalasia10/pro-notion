package http_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"src/internal/config"
	sharedEvents "src/internal/modules/shared/domain/events"
	webhooksHTTP "src/internal/modules/webhooks/interfaces/http"
)

func testWebhookConfig() *config.Config {
	return &config.Config{
		Notion: struct {
			ClientID      string
			ClientSecret  string
			RedirectURL   string
			APIVersion    string
			WebhookSecret string
		}{
			WebhookSecret: "test-webhook-secret",
		},
	}
}

var _ = Describe("Webhook Router", func() {
	var (
		testCfg   *config.Config
		publisher message.Publisher
		router    http.Handler
		secret    string
	)

	BeforeEach(func() {
		testCfg = testWebhookConfig()
		secret = testCfg.Notion.WebhookSecret

		// Set test config
		config.SetForTests(testCfg)

		// Create mock publisher/subscriber
		pubSub := gochannel.NewGoChannel(gochannel.Config{}, watermill.NewStdLogger(false, false))
		publisher = pubSub

		router = webhooksHTTP.NewRouter(publisher)
	})

	AfterEach(func() {
		config.SetForTests(nil)
	})

	computeSignature := func(payload []byte, secret string) string {
		h := hmac.New(sha256.New, []byte(secret))
		h.Write(payload)
		return "sha256=" + hex.EncodeToString(h.Sum(nil))
	}

	Context("Webhook validation", func() {
		It("should reject requests without webhook secret configured", func() {
			// Temporarily set empty secret
			emptyCfg := &config.Config{
				Notion: struct {
					ClientID      string
					ClientSecret  string
					RedirectURL   string
					APIVersion    string
					WebhookSecret string
				}{
					WebhookSecret: "",
				},
			}
			config.SetForTests(emptyCfg)

			req := httptest.NewRequest("POST", "/notion", bytes.NewReader([]byte(`{}`)))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			Expect(rec.Body.String()).To(ContainSubstring("Webhook secret not configured"))
		})

		It("should reject requests without X-Notion-Signature header", func() {
			req := httptest.NewRequest("POST", "/notion", bytes.NewReader([]byte(`{}`)))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
			Expect(rec.Body.String()).To(ContainSubstring("MISSING_SIGNATURE: Missing webhook signature header"))
		})

		It("should reject requests with invalid signature format", func() {
			req := httptest.NewRequest("POST", "/notion", bytes.NewReader([]byte(`{}`)))
			req.Header.Set("X-Notion-Signature", "invalid-format")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
			Expect(rec.Body.String()).To(ContainSubstring("Invalid signature format"))
		})

		It("should reject requests with invalid signature", func() {
			payload := []byte(`{"test": "data"}`)
			req := httptest.NewRequest("POST", "/notion", bytes.NewReader(payload))
			req.Header.Set("X-Notion-Signature", "sha256=invalid-signature")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusUnauthorized))
			Expect(rec.Body.String()).To(ContainSubstring("INVALID_SIGNATURE: Invalid webhook signature"))
		})

		It("should accept requests with valid signature", func() {
			payload := []byte(`{"test": "data"}`)
			signature := computeSignature(payload, secret)

			req := httptest.NewRequest("POST", "/notion", bytes.NewReader(payload))
			req.Header.Set("X-Notion-Signature", signature)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("Webhook event processed successfully"))
		})
	})

	Context("Webhook Handler", func() {
		It("should handle verification request", func() {
			verificationPayload := map[string]string{
				"verification_token": "test-token-123",
			}
			payloadBytes, _ := json.Marshal(verificationPayload)
			signature := computeSignature(payloadBytes, secret)

			req := httptest.NewRequest("POST", "/notion", bytes.NewReader(payloadBytes))
			req.Header.Set("X-Notion-Signature", signature)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))

			var response map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			Expect(err).ToNot(HaveOccurred())
			Expect(response["token"]).To(Equal("test-token-123"))
			Expect(response["message"]).To(ContainSubstring("Verification token received"))
		})

		It("should publish regular webhook events", func() {
			// Create a subscriber to verify publishing
			pubSub := publisher.(*gochannel.GoChannel)
			subscriber, err := pubSub.Subscribe(context.Background(), sharedEvents.NotionWebhookReceivedTopic)
			Expect(err).ToNot(HaveOccurred())

			eventPayload := map[string]interface{}{
				"type": "page.updated",
				"data": map[string]string{
					"page_id": "test-page-id",
				},
			}
			payloadBytes, _ := json.Marshal(eventPayload)
			signature := computeSignature(payloadBytes, secret)

			req := httptest.NewRequest("POST", "/notion", bytes.NewReader(payloadBytes))
			req.Header.Set("X-Notion-Signature", signature)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(ContainSubstring("Webhook event processed successfully"))

			// Verify message was published
			select {
			case msg := <-subscriber:
				var publishedEvent sharedEvents.NotionWebhookReceived
				err := json.Unmarshal(msg.Payload, &publishedEvent)
				Expect(err).ToNot(HaveOccurred())
				Expect(publishedEvent.Payload).To(Equal(payloadBytes))
				msg.Ack()
			case <-time.After(100 * time.Millisecond):
				Fail("Expected message was not published")
			}
		})
	})
})
