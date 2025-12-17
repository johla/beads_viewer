package ui

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync/atomic"
	"time"

	"github.com/Dicklesworthstone/beads_viewer/pkg/model"
	"github.com/Dicklesworthstone/beads_viewer/pkg/search"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type semanticSearchSnapshot struct {
	Ready    bool
	Index    *search.VectorIndex
	Embedder search.Embedder
	IDs      []string
}

type SemanticSearch struct {
	snapshot atomic.Value // semanticSearchSnapshot
}

func NewSemanticSearch() *SemanticSearch {
	s := &SemanticSearch{}
	s.snapshot.Store(semanticSearchSnapshot{})
	return s
}

func (s *SemanticSearch) Snapshot() semanticSearchSnapshot {
	v := s.snapshot.Load()
	if v == nil {
		return semanticSearchSnapshot{}
	}
	return v.(semanticSearchSnapshot)
}

func (s *SemanticSearch) SetIndex(idx *search.VectorIndex, embedder search.Embedder) {
	snap := s.Snapshot()
	snap.Index = idx
	snap.Embedder = embedder
	snap.Ready = idx != nil && embedder != nil
	s.snapshot.Store(snap)
}

func (s *SemanticSearch) SetIDs(ids []string) {
	snap := s.Snapshot()
	cp := make([]string, len(ids))
	copy(cp, ids)
	snap.IDs = cp
	s.snapshot.Store(snap)
}

// Filter implements list.FilterFunc, returning ranks sorted by semantic similarity.
// When the semantic index isn't ready it falls back to list.DefaultFilter.
func (s *SemanticSearch) Filter(term string, targets []string) []list.Rank {
	if term == "" {
		// Preserve existing sort order when the user hasn't entered a query yet.
		return list.DefaultFilter(term, targets)
	}

	snap := s.Snapshot()
	if !snap.Ready || snap.Index == nil || snap.Embedder == nil {
		return list.DefaultFilter(term, targets)
	}
	if len(snap.IDs) != len(targets) {
		// If we don't have a stable ID mapping, fall back to fuzzy filtering.
		return list.DefaultFilter(term, targets)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	vecs, err := snap.Embedder.Embed(ctx, []string{term})
	if err != nil || len(vecs) != 1 {
		return list.DefaultFilter(term, targets)
	}
	q := vecs[0]

	type scored struct {
		index int
		id    string
		score float64
	}
	scoredItems := make([]scored, 0, len(snap.IDs))
	for i, id := range snap.IDs {
		entry, ok := snap.Index.Get(id)
		var score float64
		if !ok {
			// Item not in index (e.g. new issue before re-indexing).
			// Assign lowest possible score to keep it in the list but at the bottom.
			score = -2.0
		} else {
			score = dotFloat32(q, entry.Vector)
		}
		scoredItems = append(scoredItems, scored{
			index: i,
			id:    id,
			score: score,
		})
	}

	sort.Slice(scoredItems, func(i, j int) bool {
		if scoredItems[i].score == scoredItems[j].score {
			return scoredItems[i].id < scoredItems[j].id
		}
		return scoredItems[i].score > scoredItems[j].score
	})

	limit := 75
	if len(scoredItems) > limit {
		scoredItems = scoredItems[:limit]
	}
	out := make([]list.Rank, 0, len(scoredItems))
	for _, it := range scoredItems {
		out = append(out, list.Rank{Index: it.index})
	}
	return out
}

// SemanticIndexReadyMsg is emitted when the semantic index build/update completes.
type SemanticIndexReadyMsg struct {
	Embedder  search.Embedder
	Index     *search.VectorIndex
	IndexPath string
	Loaded    bool
	Stats     search.IndexSyncStats
	Error     error
}

// BuildSemanticIndexCmd builds or updates the semantic index for the given issues.
func BuildSemanticIndexCmd(issues []model.Issue) tea.Cmd {
	return func() tea.Msg {
		cfg := search.EmbeddingConfigFromEnv()
		embedder, err := search.NewEmbedderFromConfig(cfg)
		if err != nil {
			return SemanticIndexReadyMsg{Error: err}
		}

		projectDir, err := os.Getwd()
		if err != nil {
			return SemanticIndexReadyMsg{Error: err}
		}

		indexPath := search.DefaultIndexPath(projectDir, cfg)
		idx, loaded, err := search.LoadOrNewVectorIndex(indexPath, embedder.Dim())
		if err != nil {
			return SemanticIndexReadyMsg{Error: err}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		docs := search.DocumentsFromIssues(issues)
		stats, err := search.SyncVectorIndex(ctx, idx, embedder, docs, 64)
		if err != nil {
			return SemanticIndexReadyMsg{Error: err}
		}
		if !loaded || stats.Changed() {
			if err := idx.Save(indexPath); err != nil {
				return SemanticIndexReadyMsg{Error: fmt.Errorf("save semantic index: %w", err)}
			}
		}

		return SemanticIndexReadyMsg{
			Embedder:  embedder,
			Index:     idx,
			IndexPath: indexPath,
			Loaded:    loaded,
			Stats:     stats,
		}
	}
}

func dotFloat32(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var sum float64
	for i := range a {
		sum += float64(a[i]) * float64(b[i])
	}
	return sum
}
