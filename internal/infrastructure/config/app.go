package config

type AppConfig struct {
	server Server
	agent  Agent
}

func (app *AppConfig) GetServer() Server {
	return app.server
}

func (app *AppConfig) GetAgent() Agent {
	return app.agent
}
