package subagents

import (
	"fmt"

	"github.com/vnniciusg/web_search_agent/web_search_agent/internal/constants"
	"github.com/vnniciusg/web_search_agent/web_search_agent/internal/prompts"
	"github.com/vnniciusg/web_search_agent/web_search_agent/internal/tools"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/loopagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
)

func GetRefinementLoopAgent(m model.LLM, googleSearch tool.Tool) (agent.Agent, error) {
	reviewerAgent, err := getReviewerAgent(m)
	if err != nil {
		return nil, err
	}

	exitAgent, err := getExitAgent(m)
	if err != nil {
		return nil, err
	}

	refinerAgent, err := getRefinerAgent(m, googleSearch, exitAgent)
	if err != nil {
		return nil, err
	}

	refinementLoop, err := loopagent.New(loopagent.Config{
		AgentConfig: agent.Config{
			Name:        "RefinementLoop",
			Description: "Iteratively reviews and refines research until quality standards are met or exit is called.",
			SubAgents: []agent.Agent{
				reviewerAgent,
				refinerAgent,
			},
		},
		MaxIterations: 3,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create refinement loop: %v", err)
	}

	return refinementLoop, nil
}

func getReviewerAgent(m model.LLM) (agent.Agent, error) {
	agent, err := llmagent.New(llmagent.Config{
		Name:        "RequirementsReviewer",
		Model:       m,
		Description: "Reviews requirements document quality and determines if refinement is needed.",
		Instruction: prompts.RequirementsReviewerPrompt,
		OutputKey:   constants.StateReviewerFeedback,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create reviewer agent: %v", err)
	}
	return agent, nil
}

func getRefinerAgent(m model.LLM, googleSearch tool.Tool, exitAgent agent.Agent) (agent.Agent, error) {
	searchAgent, err := getSearchAgent(m, googleSearch)
	if err != nil {
		return nil, err
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:        "RequirementsRefiner",
		Model:       m,
		Description: "Refines requirements based on reviewer feedback. Can call SearchAgent for more info or ExitAgent when complete.",
		Instruction: prompts.RequirementsRefinerPrompt,
		Tools: []tool.Tool{
			agenttool.New(searchAgent, nil),
			agenttool.New(exitAgent, nil),
		},
		OutputKey: constants.StateResearchFindings,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create refiner agent: %v", err)
	}
	return agent, nil
}

func getSearchAgent(m model.LLM, googleSearch tool.Tool) (agent.Agent, error) {
	agent, err := llmagent.New(llmagent.Config{
		Name:        "SearchAgent",
		Model:       m,
		Description: "Specialist agent that performs Google searches to find additional information.",
		Instruction: prompts.SearchAgentPrompt,
		Tools:       []tool.Tool{googleSearch},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create search agent: %v", err)
	}
	return agent, nil
}

func getExitAgent(m model.LLM) (agent.Agent, error) {
	exitLoopTool, err := tools.GetExitLoopTool()
	if err != nil {
		return nil, err
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:        "ExitAgent",
		Model:       m,
		Description: "Loop control agent that terminates the refinement process.",
		Instruction: prompts.ExitAgentPrompt,
		Tools:       []tool.Tool{exitLoopTool},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create exit agent: %v", err)
	}
	return agent, nil
}
