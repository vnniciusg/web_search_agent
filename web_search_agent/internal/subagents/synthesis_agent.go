package subagents

import (
	"fmt"

	"github.com/vnniciusg/web_search_agent/web_search_agent/internal/constants"
	"github.com/vnniciusg/web_search_agent/web_search_agent/internal/prompts"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
)

func GetSynthesisAgent(m model.LLM) (agent.Agent, error) {
	agent, err := llmagent.New(llmagent.Config{
		Name:        "RequirementsSynthesizer",
		Model:       m,
		Description: "Synthesizes findings from all parallel researchers into a prioritized product requirements document with KPIs.",
		Instruction: prompts.RequirementsSynthesizerPrompt,
		OutputKey:   constants.StateResearchFindings,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create synthesis agent: %v", err)
	}
	return agent, nil
}
