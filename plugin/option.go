package plugin

// Option configures a OpenAI client
type Option func(*client)

func WithToken(token string) Option {
	return func(c *client) {
		c.token = token
	}
}
