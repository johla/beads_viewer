package analysis

import (
	"testing"

	"github.com/Dicklesworthstone/beads_viewer/pkg/model"
)

// Graph: square A-B-C-D (cycle) + leaf E attached to C.
// Undirected view: A-B, B-C, C-D, D-A, C-E.
// Expectations:
// - Core numbers: A,B,C,D >=2; E lower.
// - Articulation: C (removing disconnects E).
func TestInsightsIncludesCoreAndArticulation(t *testing.T) {
	issues := []model.Issue{
		{ID: "A", Status: model.StatusOpen},
		{ID: "B", Status: model.StatusOpen, Dependencies: []*model.Dependency{{IssueID: "B", DependsOnID: "A", Type: model.DepBlocks}}},
		{ID: "C", Status: model.StatusOpen, Dependencies: []*model.Dependency{{IssueID: "C", DependsOnID: "B", Type: model.DepBlocks}}},
		{ID: "D", Status: model.StatusOpen, Dependencies: []*model.Dependency{{IssueID: "D", DependsOnID: "C", Type: model.DepBlocks}}},
		{ID: "A2", Status: model.StatusOpen, Dependencies: []*model.Dependency{{IssueID: "A2", DependsOnID: "D", Type: model.DepBlocks}, {IssueID: "A2", DependsOnID: "A", Type: model.DepBlocks}}}, // closes cycle
		{ID: "E", Status: model.StatusOpen, Dependencies: []*model.Dependency{{IssueID: "E", DependsOnID: "C", Type: model.DepBlocks}}},
	}

	an := NewAnalyzer(issues)
	stats := an.Analyze()
	ins := stats.GenerateInsights(10)

	if len(ins.Cores) == 0 {
		t.Fatalf("expected cores list populated; core map=%v", stats.CoreNumber())
	}
	if ins.Cores[0].Value < ins.Cores[len(ins.Cores)-1].Value {
		t.Fatalf("cores not sorted desc: %#v", ins.Cores)
	}
	foundC := false
	for _, id := range ins.Articulation {
		if id == "C" {
			foundC = true
		}
	}
	if !foundC {
		t.Fatalf("expected articulation points to include C, got %v", ins.Articulation)
	}
}
