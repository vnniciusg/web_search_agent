package tools

import (
	"fmt"

	"google.golang.org/adk/tool"
)

type ExitLoopArgs struct{}

type ExitLoopResult struct {
	Status string `json:"status"`
}

func ExitLoop(ctx tool.Context, args ExitLoopArgs) (ExitLoopResult, error) {
	fmt.Printf("[TOOL CALL] exit loop triggered by %s - Reasearch refinement complete", ctx.AgentName())
	ctx.Actions().Escalate = true
	return ExitLoopResult{Status: "loop_exited"}, nil
}
