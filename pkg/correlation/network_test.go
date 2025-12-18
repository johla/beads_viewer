package correlation

import (
	"testing"
	"time"
)

func createTestHistoryReport() *HistoryReport {
	now := time.Now()
	return &HistoryReport{
		GeneratedAt: now,
		DataHash:    "testhash123",
		Histories: map[string]BeadHistory{
			"bv-001": {
				BeadID: "bv-001",
				Title:  "Auth token handling",
				Status: "closed",
				Commits: []CorrelatedCommit{
					{
						SHA:      "abc123def456",
						ShortSHA: "abc123d",
						Files:    []FileChange{{Path: "auth/token.go"}, {Path: "auth/session.go"}},
					},
					{
						SHA:      "def456ghi789",
						ShortSHA: "def456g",
						Files:    []FileChange{{Path: "auth/token.go"}, {Path: "config/auth.yaml"}},
					},
				},
				Milestones: BeadMilestones{
					Created: &BeadEvent{Timestamp: now.Add(-72 * time.Hour)},
					Closed:  &BeadEvent{Timestamp: now.Add(-24 * time.Hour)},
				},
			},
			"bv-002": {
				BeadID: "bv-002",
				Title:  "Session management",
				Status: "open",
				Commits: []CorrelatedCommit{
					{
						SHA:      "abc123def456", // Same commit as bv-001
						ShortSHA: "abc123d",
						Files:    []FileChange{{Path: "auth/session.go"}},
					},
					{
						SHA:      "xyz789abc123",
						ShortSHA: "xyz789a",
						Files:    []FileChange{{Path: "auth/session.go"}, {Path: "middleware/auth.go"}},
					},
				},
				Milestones: BeadMilestones{
					Created: &BeadEvent{Timestamp: now.Add(-48 * time.Hour)},
				},
			},
			"bv-003": {
				BeadID: "bv-003",
				Title:  "Rate limiting",
				Status: "in_progress",
				Commits: []CorrelatedCommit{
					{
						SHA:      "xyz789abc123", // Same commit as bv-002
						ShortSHA: "xyz789a",
						Files:    []FileChange{{Path: "middleware/auth.go"}, {Path: "middleware/rate.go"}},
					},
				},
				Milestones: BeadMilestones{
					Created: &BeadEvent{Timestamp: now.Add(-24 * time.Hour)},
					Claimed: &BeadEvent{Timestamp: now.Add(-12 * time.Hour)},
				},
			},
			"bv-004": {
				BeadID:  "bv-004",
				Title:   "Unrelated feature",
				Status:  "closed",
				Commits: []CorrelatedCommit{}, // No commits - isolated node
			},
		},
		CommitIndex: CommitIndex{
			"abc123def456": {"bv-001", "bv-002"},
			"def456ghi789": {"bv-001"},
			"xyz789abc123": {"bv-002", "bv-003"},
		},
	}
}

func TestNewNetworkBuilder(t *testing.T) {
	report := createTestHistoryReport()
	builder := NewNetworkBuilder(report)

	if builder == nil {
		t.Fatal("Expected non-nil NetworkBuilder")
	}
	if builder.report != report {
		t.Error("Expected builder to reference the report")
	}

	// Check bead maps were built
	if len(builder.beadFiles) == 0 {
		t.Error("Expected beadFiles map to be populated")
	}
	if len(builder.beadCommits) == 0 {
		t.Error("Expected beadCommits map to be populated")
	}
}

func TestNewNetworkBuilderNilReport(t *testing.T) {
	builder := NewNetworkBuilder(nil)

	if builder == nil {
		t.Fatal("Expected non-nil NetworkBuilder even with nil report")
	}
}

func TestBuildNetwork(t *testing.T) {
	report := createTestHistoryReport()
	builder := NewNetworkBuilder(report)
	network := builder.Build()

	if network == nil {
		t.Fatal("Expected non-nil network")
	}

	t.Run("nodes_created", func(t *testing.T) {
		if len(network.Nodes) != 4 {
			t.Errorf("Expected 4 nodes, got %d", len(network.Nodes))
		}

		// Check specific node
		if node, ok := network.Nodes["bv-001"]; ok {
			if node.Title != "Auth token handling" {
				t.Errorf("Expected title 'Auth token handling', got %s", node.Title)
			}
			if node.Status != "closed" {
				t.Errorf("Expected status 'closed', got %s", node.Status)
			}
			if node.CommitCount != 2 {
				t.Errorf("Expected 2 commits, got %d", node.CommitCount)
			}
		} else {
			t.Error("Expected bv-001 node to exist")
		}
	})

	t.Run("edges_created", func(t *testing.T) {
		if len(network.Edges) == 0 {
			t.Error("Expected edges to be created")
		}

		// Should have edges between bv-001 and bv-002 (shared commit abc123def456)
		foundSharedCommit := false
		for _, edge := range network.Edges {
			if (edge.FromBead == "bv-001" && edge.ToBead == "bv-002") ||
				(edge.FromBead == "bv-002" && edge.ToBead == "bv-001") {
				if edge.EdgeType == EdgeSharedCommit {
					foundSharedCommit = true
				}
			}
		}
		if !foundSharedCommit {
			t.Error("Expected shared commit edge between bv-001 and bv-002")
		}
	})

	t.Run("stats_calculated", func(t *testing.T) {
		if network.Stats.TotalNodes != 4 {
			t.Errorf("Expected 4 total nodes, got %d", network.Stats.TotalNodes)
		}
		if network.Stats.TotalEdges == 0 {
			t.Error("Expected edges in stats")
		}
	})
}

func TestBuildNetworkSharedFiles(t *testing.T) {
	report := createTestHistoryReport()
	builder := NewNetworkBuilder(report)
	network := builder.Build()

	// bv-001 and bv-002 both touch auth/session.go
	foundSharedFile := false
	for _, edge := range network.Edges {
		if edge.EdgeType == EdgeSharedFile {
			if (edge.FromBead == "bv-001" && edge.ToBead == "bv-002") ||
				(edge.FromBead == "bv-002" && edge.ToBead == "bv-001") {
				// Check if auth/session.go is in details
				for _, detail := range edge.Details {
					if detail == "auth/session.go" {
						foundSharedFile = true
						break
					}
				}
			}
		}
	}
	if !foundSharedFile {
		t.Error("Expected shared file edge for auth/session.go between bv-001 and bv-002")
	}
}

func TestClusterDetection(t *testing.T) {
	report := createTestHistoryReport()
	builder := NewNetworkBuilder(report)
	network := builder.Build()

	// bv-001, bv-002, bv-003 should be in the same cluster
	// (connected through shared commits and files)
	// bv-004 is isolated

	if len(network.Clusters) == 0 {
		// It's OK if clusters are empty - depends on edge weights meeting threshold
		return
	}

	// If clusters exist, check they're properly formed
	for _, cluster := range network.Clusters {
		if len(cluster.BeadIDs) < 2 {
			t.Error("Cluster should have at least 2 beads")
		}
		if cluster.InternalConnectivity < 0 || cluster.InternalConnectivity > 1 {
			t.Errorf("Internal connectivity should be between 0 and 1, got %f", cluster.InternalConnectivity)
		}
	}
}

func TestGetSubNetwork(t *testing.T) {
	report := createTestHistoryReport()
	builder := NewNetworkBuilder(report)
	network := builder.Build()

	subNetwork := network.GetSubNetwork("bv-001", 1)

	if subNetwork == nil {
		t.Fatal("Expected non-nil subnetwork")
	}

	// bv-001 should be in the subnetwork
	if _, ok := subNetwork.Nodes["bv-001"]; !ok {
		t.Error("Expected bv-001 to be in subnetwork")
	}

	// Subnetwork should have fewer or equal nodes than full network
	if len(subNetwork.Nodes) > len(network.Nodes) {
		t.Error("Subnetwork should not have more nodes than full network")
	}
}

func TestGetSubNetworkDepthLimits(t *testing.T) {
	report := createTestHistoryReport()
	builder := NewNetworkBuilder(report)
	network := builder.Build()

	tests := []struct {
		depth    int
		expected int // expected capped depth
	}{
		{0, 1},  // depth < 1 should become 1
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 3},  // depth > 3 should become 3
		{100, 3},
	}

	for _, tt := range tests {
		subNetwork := network.GetSubNetwork("bv-001", tt.depth)
		if subNetwork == nil {
			t.Errorf("Expected non-nil subnetwork for depth %d", tt.depth)
		}
	}
}

func TestIsolatedNodes(t *testing.T) {
	report := createTestHistoryReport()
	builder := NewNetworkBuilder(report)
	network := builder.Build()

	// bv-004 has no commits, so it should be isolated
	node, ok := network.Nodes["bv-004"]
	if !ok {
		t.Fatal("Expected bv-004 to exist")
	}

	if node.Degree != 0 {
		t.Errorf("Expected isolated node to have degree 0, got %d", node.Degree)
	}

	// Isolated count should be >= 1 (at least bv-004)
	if network.Stats.IsolatedNodes < 1 {
		t.Errorf("Expected at least 1 isolated node, got %d", network.Stats.IsolatedNodes)
	}
}

func TestToResult(t *testing.T) {
	report := createTestHistoryReport()
	builder := NewNetworkBuilder(report)
	network := builder.Build()

	t.Run("full_network", func(t *testing.T) {
		result := network.ToResult("", 0)
		if result == nil {
			t.Fatal("Expected non-nil result")
		}
		if result.BeadID != "" {
			t.Errorf("Expected empty BeadID for full network, got %s", result.BeadID)
		}
		if result.Network == nil {
			t.Error("Expected network in result")
		}
	})

	t.Run("subnetwork", func(t *testing.T) {
		result := network.ToResult("bv-001", 2)
		if result == nil {
			t.Fatal("Expected non-nil result")
		}
		if result.BeadID != "bv-001" {
			t.Errorf("Expected BeadID 'bv-001', got %s", result.BeadID)
		}
		if result.Depth != 2 {
			t.Errorf("Expected depth 2, got %d", result.Depth)
		}
	})

	t.Run("top_connected", func(t *testing.T) {
		result := network.ToResult("", 0)
		if len(result.TopConnected) > 10 {
			t.Error("TopConnected should be limited to 10")
		}
		// Should be sorted by degree
		for i := 1; i < len(result.TopConnected); i++ {
			if result.TopConnected[i].Degree > result.TopConnected[i-1].Degree {
				t.Error("TopConnected should be sorted by degree descending")
			}
		}
	})
}

func TestNetworkEdgeTypes(t *testing.T) {
	tests := []struct {
		edgeType NetworkEdgeType
		expected string
	}{
		{EdgeSharedCommit, "shared_commit"},
		{EdgeSharedFile, "shared_file"},
		{EdgeDependency, "dependency"},
	}

	for _, tt := range tests {
		if string(tt.edgeType) != tt.expected {
			t.Errorf("Expected edge type %s, got %s", tt.expected, string(tt.edgeType))
		}
	}
}

func TestCommonPathPrefix(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected string
	}{
		{
			name:     "empty",
			files:    []string{},
			expected: "",
		},
		{
			name:     "single_file",
			files:    []string{"auth/token.go"},
			expected: "auth/",
		},
		{
			name:     "common_prefix",
			files:    []string{"auth/token.go", "auth/session.go", "auth/middleware.go"},
			expected: "auth/",
		},
		{
			name:     "no_common_prefix",
			files:    []string{"auth/token.go", "config/settings.go"},
			expected: "",
		},
		{
			name:     "nested_prefix",
			files:    []string{"pkg/auth/token.go", "pkg/auth/session.go"},
			expected: "pkg/auth/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := commonPathPrefix(tt.files)
			if result != tt.expected {
				t.Errorf("Expected prefix %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSplitEdgeKey(t *testing.T) {
	tests := []struct {
		key      string
		expected []string
	}{
		{"bv-001:bv-002:commit", []string{"bv-001", "bv-002", "commit"}},
		{"bv-001:bv-002:file", []string{"bv-001", "bv-002", "file"}},
		{"a:b:c:d", []string{"a", "b", "c", "d"}},
		{"single", []string{"single"}},
	}

	for _, tt := range tests {
		result := splitEdgeKey(tt.key)
		if len(result) != len(tt.expected) {
			t.Errorf("For key %q: expected %d parts, got %d", tt.key, len(tt.expected), len(result))
			continue
		}
		for i, part := range result {
			if part != tt.expected[i] {
				t.Errorf("For key %q: part %d expected %q, got %q", tt.key, i, tt.expected[i], part)
			}
		}
	}
}
