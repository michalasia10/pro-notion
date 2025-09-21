package events

const NotionWebhookReceivedTopic = "notion.webhook.received"

type NotionWebhookReceived struct {
	Payload []byte
}
