package subagents

import (
	"fmt"

	"github.com/vnniciusg/web_search_agent/web_search_agent/internal/constants"
	"github.com/vnniciusg/web_search_agent/web_search_agent/internal/prompts"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
)

func GetFinalReportAgent(m model.LLM) (agent.Agent, error) {
	agent, err := llmagent.New(llmagent.Config{
		Name:        "FinalPRDGenerator",
		Model:       m,
		Description: "Generates the final polished Product Requirements Document.",
		Instruction: prompts.FinalReportPrompt,
		OutputKey:   constants.StateFinalReport,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create final report agent: %v", err)
	}
	return agent, nil
}
