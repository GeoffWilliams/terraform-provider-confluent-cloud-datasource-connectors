package provider

import (
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type SchemaRegistryRestClientFactory struct {
	userAgent  string
	maxRetries *int
}

type KafkaRestClientFactory struct {
	userAgent  string
	maxRetries *int
}

type RetryableClientFactoryOption = func(c *RetryableClientFactory)

type RetryableClientFactory struct {
	maxRetries *int
}

func WithMaxRetries(maxRetries int) RetryableClientFactoryOption {
	return func(c *RetryableClientFactory) {
		c.maxRetries = &maxRetries
	}
}

func NewRetryableClientFactory(opts ...RetryableClientFactoryOption) *RetryableClientFactory {
	c := &RetryableClientFactory{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// CreateRetryableClient creates retryable HTTP client that performs automatic retries with exponential backoff for 429
// and 5** (except 501) errors. Otherwise, the response is returned and left to the caller to interpret.
func (f RetryableClientFactory) CreateRetryableClient() *http.Client {
	// Implicitly using default retry configuration
	// under the assumption is it's OK to spend retrying a single HTTP call around 15 seconds in total: 1 + 2 + 4 + 8
	// An exponential backoff equation: https://github.com/hashicorp/go-retryablehttp/blob/master/client.go#L493
	// retryWaitMax = math.Pow(2, float64(attemptNum)) * float64(retryWaitMin)
	// defaultRetryWaitMin = 1 * time.Second
	// defaultRetryWaitMax = 30 * time.Second
	// defaultRetryMax     = 4

	retryClient := retryablehttp.NewClient()

	if f.maxRetries != nil {
		retryClient.RetryMax = *f.maxRetries
	}

	return retryClient.StandardClient()
}
