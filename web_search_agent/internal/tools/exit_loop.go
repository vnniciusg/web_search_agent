package tools

import (
	"fmt"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type ExitLoopArgs struct{}

type ExitLoopResult struct {
	Status string `json:"status"`
}

func exitLoop(ctx tool.Context, args ExitLoopArgs) (ExitLoopResult, error) {
	fmt.Printf("[TOOL CALL] exit loop triggered by %s - Reasearch refinement complete", ctx.AgentName())
	ctx.Actions().Escalate = true
	return ExitLoopResult{Status: "loop_exited"}, nil
}

func GetExitLoopTool() (tool.Tool, error) {
	exitLoopTool, err := functiontool.New(
		functiontool.Config{
			Name:        "exit_loop",
			Description: "Call this function ONLY when the reviewer indicates the research is complete (RESEARCH_COMPLETE). This will end the refinement loop.",
		},
		exitLoop,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create exit loop tool: %v", err)
	}
	return exitLoopTool, nil
}
