package eventbus

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

// NewPublisher creates a new message publisher using GoChannel for local development
func NewPublisher(logger watermill.LoggerAdapter) (message.Publisher, error) {
	pub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	return pub, nil
}

// NewRouter creates a new message router
func NewRouter(logger watermill.LoggerAdapter) (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}
	return router, nil
}
