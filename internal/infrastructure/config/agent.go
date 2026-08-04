package config

type Agent struct {
	Name  string
	model Model
	skill Skill
}

type Model struct {
	BaseUrl   string
	ApiKey    string
	ModelName string
}

type Skill struct {
	Location string
}

func (agent Agent) GetModel() Model {
	return agent.model
}

func (agent Agent) GetSkill() Skill {
	return agent.skill
}
