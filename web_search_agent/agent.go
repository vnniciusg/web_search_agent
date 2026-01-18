package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vnniciusg/web_search_agent/web_search_agent/internal/constants"
	"github.com/vnniciusg/web_search_agent/web_search_agent/internal/subagents"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	ctx := context.Background()

	if err := runAgent(ctx); err != nil {
		log.Fatalf("Agent execution failed: %v", err)
	}
}

func runAgent(ctx context.Context) error {
	model, err := gemini.NewModel(ctx, constants.ModelName, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return fmt.Errorf("failed to create model: %v", err)
	}

	googleSearch := geminitool.GoogleSearch{}

	rootAgent, err := getRootAgent(model, googleSearch)
	if err != nil {
		return err
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(rootAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		return fmt.Errorf("run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}

	return nil
}

func getRootAgent(m model.LLM, googleSearch tool.Tool) (agent.Agent, error) {
	parallelResearchAgent, err := subagents.GetParallelResearchAgent(m, googleSearch)
	if err != nil {
		return nil, err
	}

	synthesisAgent, err := subagents.GetSynthesisAgent(m)
	if err != nil {
		return nil, err
	}

	refinementLoop, err := subagents.GetRefinementLoopAgent(m, googleSearch)
	if err != nil {
		return nil, err
	}

	finalReportAgent, err := subagents.GetFinalReportAgent(m)
	if err != nil {
		return nil, err
	}

	rootAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        constants.AppName,
			Description: "Comprehensive criticism research agent that uses parallel search and iterative refinement to find all criticisms about a topic.",
			SubAgents: []agent.Agent{
				parallelResearchAgent,
				synthesisAgent,
				refinementLoop,
				finalReportAgent,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create root agent: %v", err)
	}

	return rootAgent, nil
}
