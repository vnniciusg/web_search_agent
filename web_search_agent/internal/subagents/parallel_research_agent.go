package subagents

import (
	"fmt"

	"github.com/vnniciusg/web_search_agent/web_search_agent/internal/constants"
	"github.com/vnniciusg/web_search_agent/web_search_agent/internal/prompts"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/parallelagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
)

func GetParallelResearchAgent(m model.LLM, googleSearch tool.Tool) (agent.Agent, error) {
	userPainPointsResearcher, err := getUserPainPointsResearcher(m, googleSearch)
	if err != nil {
		return nil, err
	}

	competitorResearcher, err := getCompetitorResearcher(m, googleSearch)
	if err != nil {
		return nil, err
	}

	technicalResearcher, err := getTechnicalResearcher(m, googleSearch)
	if err != nil {
		return nil, err
	}

	marketGapsResearcher, err := getMarketGapsResearcher(m, googleSearch)
	if err != nil {
		return nil, err
	}

	parallelResearchAgent, err := parallelagent.New(parallelagent.Config{
		AgentConfig: agent.Config{
			Name:        "ParallelRequirementsResearch",
			Description: "Executes multiple specialized researchers in parallel to gather comprehensive product requirements insights.",
			SubAgents: []agent.Agent{
				userPainPointsResearcher,
				competitorResearcher,
				technicalResearcher,
				marketGapsResearcher,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create parallel research agent: %v", err)
	}

	return parallelResearchAgent, nil
}

func getUserPainPointsResearcher(m model.LLM, googleSearch tool.Tool) (agent.Agent, error) {
	agent, err := llmagent.New(llmagent.Config{
		Name:        "UserPainPointsResearcher",
		Model:       m,
		Description: "Searches for user complaints, pain points, and unmet needs that translate into product requirements.",
		Instruction: prompts.UserPainPointsResearcherPrompt,
		Tools:       []tool.Tool{googleSearch},
		OutputKey:   constants.StateUserPainPoints,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user pain points researcher: %v", err)
	}
	return agent, nil
}

func getCompetitorResearcher(m model.LLM, googleSearch tool.Tool) (agent.Agent, error) {
	agent, err := llmagent.New(llmagent.Config{
		Name:        "CompetitorWeaknessResearcher",
		Model:       m,
		Description: "Analyzes competitor weaknesses and gaps to identify differentiation opportunities.",
		Instruction: prompts.CompetitorWeaknessResearcherPrompt,
		Tools:       []tool.Tool{googleSearch},
		OutputKey:   constants.StateCompetitorFindings,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create competitor researcher: %v", err)
	}
	return agent, nil
}

func getTechnicalResearcher(m model.LLM, googleSearch tool.Tool) (agent.Agent, error) {
	agent, err := llmagent.New(llmagent.Config{
		Name:        "TechnicalLimitationsResearcher",
		Model:       m,
		Description: "Identifies technical limitations and constraints that inform technical requirements and KPIs.",
		Instruction: prompts.TechnicalLimitationsResearcherPrompt,
		Tools:       []tool.Tool{googleSearch},
		OutputKey:   constants.StateTechnicalFindings,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create technical researcher: %v", err)
	}
	return agent, nil
}

func getMarketGapsResearcher(m model.LLM, googleSearch tool.Tool) (agent.Agent, error) {
	agent, err := llmagent.New(llmagent.Config{
		Name:        "MarketGapsResearcher",
		Model:       m,
		Description: "Identifies underserved market segments and emerging opportunities.",
		Instruction: prompts.MarketGapsResearcherPrompt,
		Tools:       []tool.Tool{googleSearch},
		OutputKey:   constants.StateMarketGapsFindings,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create market gaps researcher: %v", err)
	}
	return agent, nil
}
