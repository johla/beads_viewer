package analysis

import (
	"testing"

	"github.com/Dicklesworthstone/beads_viewer/pkg/model"
)

func TestExtractKeywords(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		want        []string
	}{
		{
			name:        "simple",
			title:       "Fix login bug",
			description: "User cannot login with email",
			want:        []string{"fix", "login", "bug", "user", "cannot", "email"},
		},
		{
			name:        "stopwords",
			title:       "The issue with the login",
			description: "It is not working for some users",
			want:        []string{"issue", "login", "working", "users"},
		},
		{
			name:        "short words",
			title:       "UI fix",
			description: "Go to page",
			want:        []string{"page"}, // "UI", "fix", "Go", "to" filtered out or short? Wait, "fix" is 3 chars.
		},
		{
			name:        "deduplication",
			title:       "Login login",
			description: "Login",
			want:        []string{"login"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractKeywords(tt.title, tt.description)
			// Check if all expected keywords are present
			for _, w := range tt.want {
				found := false
				for _, g := range got {
					if g == w {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("extractKeywords() missing %q, got %v", w, got)
				}
			}
		})
	}
}

func TestJaccardSimilarity(t *testing.T) {
	tests := []struct {
		name  string
		set1  []string
		set2  []string
		want  float64
	}{
		{
			name: "identical",
			set1: []string{"a", "b", "c"},
			set2: []string{"a", "b", "c"},
			want: 1.0,
		},
		{
			name: "disjoint",
			set1: []string{"a", "b"},
			set2: []string{"c", "d"},
			want: 0.0,
		},
		{
			name: "partial overlap",
			set1: []string{"a", "b", "c"},
			set2: []string{"b", "c", "d"},
			want: 0.5, // intersection(b,c)=2, union(a,b,c,d)=4 -> 2/4 = 0.5
		},
		{
			name: "subset",
			set1: []string{"a", "b"},
			set2: []string{"a", "b", "c", "d"},
			want: 0.5, // intersection=2, union=4 -> 0.5
		},
		{
			name: "empty",
			set1: []string{},
			set2: []string{"a"},
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := jaccardSimilarity(tt.set1, tt.set2)
			if got != tt.want {
				t.Errorf("jaccardSimilarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectDuplicates(t *testing.T) {
	issues := []model.Issue{
		{ID: "1", Title: "Fix login page", Description: "Login button is broken", Status: model.StatusOpen},
		{ID: "2", Title: "Login button fix", Description: "Cannot login, button issue", Status: model.StatusOpen},
		{ID: "3", Title: "Database migration", Description: "Add users table", Status: model.StatusOpen},
	}

	config := DefaultDuplicateConfig()
	config.JaccardThreshold = 0.1 // Low threshold to ensure detection

	suggestions := DetectDuplicates(issues, config)

	if len(suggestions) == 0 {
		t.Error("expected duplicate suggestion")
	} else {
		sug := suggestions[0]
		if sug.Type != SuggestionPotentialDuplicate {
			t.Errorf("expected suggestion type %q, got %q", SuggestionPotentialDuplicate, sug.Type)
		}
		if (sug.TargetBead == "1" && sug.RelatedBead == "2") || (sug.TargetBead == "2" && sug.RelatedBead == "1") {
			// Correct pair
		} else {
			t.Errorf("expected duplicate pair 1-2, got %s-%s", sug.TargetBead, sug.RelatedBead)
		}
	}
}
