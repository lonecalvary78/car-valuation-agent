package config

type Agent struct {
	Name  string
	model Model
	skill Skill
}

type Model struct {
	BaseUrl      string
	ApiKey       string
	ModelName    string
	AdvanceSetup AdvancedSetup
}

type AdvancedSetup struct {
	Temperature      float32
	TopP             float32
	TopK             float32
	PresencePenalty  float32
	FrequencyPenalty float32
	MaximumTokens    int32
}

func (advancedSetup AdvancedSetup) HasAdvancedSetup() bool {
	return advancedSetup.Temperature > 0.0 &&
		advancedSetup.TopK > 0.0 &&
		advancedSetup.TopP > 0.0 &&
		advancedSetup.PresencePenalty > 0.0 &&
		advancedSetup.FrequencyPenalty > 0.0 &&
		advancedSetup.MaximumTokens > 0
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
