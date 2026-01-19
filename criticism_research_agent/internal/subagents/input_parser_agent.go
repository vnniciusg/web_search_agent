package subagents

import (
	"github.com/vnniciusg/criticism_research_agent/criticism_research_agent/internal/constants"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
)

func GetInputParserAgent(m model.LLM) (agent.Agent, error) {
	agent, err := llmagent.New(llmagent.Config{
		Name:        "InputParser",
		Model:       m,
		Description: "Parses and extract the research topic from user input.",
		Instruction: `Extract the main topic or product category from the user's message.                                                                                                                                                                                  
		Output ONLY the topic name, nothing else. For example:                                                                                                                                                                                                                     
		- If user says "research project management tools", output: project management tools                                                                                                                                                                                       
		- If user says "analyze food delivery apps", output: food delivery apps`,
		OutputKey: constants.StateTopic,
	})

	return agent, err
}
