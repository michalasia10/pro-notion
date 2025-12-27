package eventbus

import (
	"errors"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	redisstream "github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/redis/go-redis/v9"
)

// Transport selects the Watermill backend.
type Transport string

const (
	TransportGoChannel Transport = "gochannel"
	TransportRedis     Transport = "redis"
)

// Config carries settings for building publishers/subscribers.
type Config struct {
	Transport       Transport
	RedisOptions    redis.Options
	ConsumerGroup   string
	ConsumerTimeout time.Duration
	ClaimInterval   time.Duration
}

// PubSub bundles publisher/subscriber to keep transports consistent.
type PubSub struct {
	Publisher  message.Publisher
	Subscriber message.Subscriber
	Close      func() error
}

func newPubSubInMemory(logger watermill.LoggerAdapter, cfg Config) (*PubSub, error) {
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	return &PubSub{
		Publisher:  pubsub,
		Subscriber: pubsub,
		Close:      func() error { return nil },
	}, nil
}

func newPubRedis(logger watermill.LoggerAdapter, cfg Config) (*PubSub, error) {
	pubClient := redis.NewClient(&cfg.RedisOptions)
	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client:        pubClient,
		Marshaller:    redisstream.DefaultMarshallerUnmarshaller{},
		DefaultMaxlen: 0,
	}, logger)
	if err != nil {
		return nil, err
	}

	subClient := redis.NewClient(&cfg.RedisOptions)
	subscriber, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:                 subClient,
		Unmarshaller:           redisstream.DefaultMarshallerUnmarshaller{},
		ConsumerGroup:          cfg.ConsumerGroup,
		ConsumerTimeout:        cfg.ConsumerTimeout,
		ClaimInterval:          cfg.ClaimInterval,
		ClaimBatchSize:         redisstream.DefaultClaimBatchSize,
		MaxIdleTime:            redisstream.DefaultMaxIdleTime,
		CheckConsumersInterval: redisstream.DefaultCheckConsumersInterval,
	}, logger)
	if err != nil {
		_ = publisher.Close()
		_ = pubClient.Close()
		return nil, err
	}
	return &PubSub{
		Publisher:  publisher,
		Subscriber: subscriber,
		Close: func() error {
			return errors.Join(publisher.Close(), subscriber.Close())
		},
	}, nil
}

// NewPubSub builds a publisher/subscriber pair for the configured transport.
func NewPubSub(logger watermill.LoggerAdapter, cfg Config) (*PubSub, error) {
	if cfg.Transport == "" {
		cfg.Transport = TransportRedis
	}
	switch cfg.Transport {
	case TransportGoChannel:
		return newPubSubInMemory(logger, cfg)
	case TransportRedis:
		return newPubRedis(logger, cfg)
	default:
		return nil, errors.New("unsupported event bus transport")
	}
}

// NewRouter creates a new message router
func NewRouter(logger watermill.LoggerAdapter) (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}
	return router, nil
}
