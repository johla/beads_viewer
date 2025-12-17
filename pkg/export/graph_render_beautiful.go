package export

import (
	"fmt"
	"time"
)

// generateUltimateHTML creates the enhanced HTML visualization with all features
func generateUltimateHTML(title, dataHash, graphDataJSON string, nodeCount, edgeCount int, projectName, forceGraphLib, markedLib string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s | bv Graph</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #0f0f1a;
            --bg-secondary: #1a1a2e;
            --bg-tertiary: #16213e;
            --bg-elevated: #252545;
            --bg-glass: rgba(26, 26, 46, 0.85);
            --fg: #e8e8f0;
            --fg-muted: #8888aa;
            --fg-dim: #555577;
            --purple: #a855f7;
            --purple-glow: rgba(168, 85, 247, 0.4);
            --pink: #ec4899;
            --cyan: #22d3ee;
            --green: #22c55e;
            --orange: #f97316;
            --red: #ef4444;
            --yellow: #eab308;
            --gold: #fbbf24;
            --gold-glow: rgba(251, 191, 36, 0.6);
            --shadow: 0 8px 32px rgba(0,0,0,0.4);
            --shadow-glow: 0 0 40px var(--purple-glow);
            --radius: 12px;
            --radius-lg: 16px;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background: var(--bg);
            color: var(--fg);
            height: 100vh;
            display: flex;
            flex-direction: column;
            overflow: hidden;
            font-size: 14px;
            line-height: 1.5;
        }

        /* Header */
        header {
            background: linear-gradient(135deg, var(--bg-secondary), var(--bg-tertiary));
            padding: 0.75rem 1.5rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
            border-bottom: 1px solid var(--purple);
            z-index: 100;
            box-shadow: var(--shadow);
            backdrop-filter: blur(10px);
        }
        .logo { display: flex; align-items: center; gap: 0.75rem; }
        .logo-icon {
            width: 40px; height: 40px;
            background: linear-gradient(135deg, var(--purple), var(--pink));
            border-radius: var(--radius);
            display: flex; align-items: center; justify-content: center;
            font-family: 'JetBrains Mono', monospace;
            font-weight: 700; font-size: 16px;
            box-shadow: var(--shadow-glow);
        }
        h1 { font-size: 1.25rem; font-weight: 600; letter-spacing: -0.02em; }
        h1 span {
            background: linear-gradient(90deg, var(--purple), var(--pink), var(--cyan));
            -webkit-background-clip: text; -webkit-text-fill-color: transparent;
            background-size: 200%% auto;
            animation: shimmer 3s linear infinite;
        }
        @keyframes shimmer { 0%% { background-position: 0%% center; } 100%% { background-position: 200%% center; } }

        /* Toolbar */
        .toolbar { display: flex; gap: 0.75rem; align-items: center; flex-wrap: wrap; }
        .toolbar-group {
            display: flex; gap: 0.25rem;
            padding: 0.25rem;
            background: var(--bg);
            border-radius: var(--radius);
            border: 1px solid var(--bg-elevated);
        }
        button, select {
            font-family: inherit;
            font-size: 0.8rem;
            padding: 0.5rem 0.875rem;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            transition: all 0.2s ease;
            font-weight: 500;
        }
        button { background: transparent; color: var(--fg-muted); }
        button:hover { background: var(--bg-elevated); color: var(--fg); transform: translateY(-1px); }
        button.active {
            background: linear-gradient(135deg, var(--purple), var(--pink));
            color: white;
            box-shadow: 0 4px 15px var(--purple-glow);
        }
        select {
            background: var(--bg);
            color: var(--fg);
            border: 1px solid var(--bg-elevated);
            padding-right: 2rem;
        }
        select:focus { outline: none; border-color: var(--purple); box-shadow: 0 0 0 3px var(--purple-glow); }

        /* Search */
        .search-container { position: relative; }
        .search-input {
            font-family: inherit;
            font-size: 0.875rem;
            padding: 0.5rem 1rem 0.5rem 2.5rem;
            background: var(--bg);
            color: var(--fg);
            border: 1px solid var(--bg-elevated);
            border-radius: var(--radius);
            width: 280px;
            transition: all 0.2s ease;
        }
        .search-input:focus {
            outline: none;
            border-color: var(--purple);
            box-shadow: 0 0 0 3px var(--purple-glow);
            width: 320px;
        }
        .search-icon {
            position: absolute; left: 0.875rem; top: 50%%; transform: translateY(-50%%);
            color: var(--fg-muted); font-size: 0.875rem;
        }

        /* Search Results Dropdown */
        .search-results {
            position: absolute; top: 100%%; left: 0; right: 0;
            background: var(--bg-glass); backdrop-filter: blur(20px);
            border: 1px solid var(--purple); border-radius: var(--radius);
            margin-top: 0.5rem; max-height: 400px; overflow-y: auto;
            z-index: 1000; display: none; box-shadow: var(--shadow);
        }
        .search-results.visible { display: block; }
        .search-result-item {
            padding: 0.875rem 1rem; cursor: pointer;
            border-bottom: 1px solid var(--bg-elevated);
            transition: all 0.15s ease;
        }
        .search-result-item:hover { background: var(--bg-elevated); }
        .search-result-item:last-child { border-bottom: none; }
        .search-result-id { font-family: 'JetBrains Mono', monospace; color: var(--cyan); font-size: 0.8rem; font-weight: 600; }
        .search-result-title { font-size: 0.875rem; margin-top: 0.25rem; }
        .search-result-preview {
            font-size: 0.75rem; color: var(--fg-muted); margin-top: 0.5rem;
            max-height: 60px; overflow: hidden;
        }
        .search-result-preview p { margin: 0.25rem 0; }

        /* Main */
        main { flex: 1; display: flex; overflow: hidden; position: relative; }
        #graph-wrapper { flex: 1; position: relative; }
        #graph-container {
            position: absolute; top: 0; left: 0; right: 0; bottom: 0;
            background: radial-gradient(ellipse at center, var(--bg-secondary) 0%%, var(--bg) 100%%);
        }

        /* Overlay Stats */
        .overlay-stats {
            position: absolute; top: 1rem; left: 1rem;
            background: var(--bg-glass); backdrop-filter: blur(15px);
            padding: 0.75rem 1rem; border-radius: var(--radius);
            font-size: 0.8rem; color: var(--fg-muted);
            z-index: 10; display: flex; gap: 1.25rem;
            border: 1px solid var(--bg-elevated);
            box-shadow: var(--shadow);
        }
        .overlay-stats .stat { display: flex; align-items: center; gap: 0.375rem; }
        .overlay-stats .stat-value { color: var(--cyan); font-weight: 600; font-family: 'JetBrains Mono', monospace; }

        /* Hover Panel - Full Details */
        #hover-panel {
            position: absolute; top: 1rem; left: 50%%; transform: translateX(-50%%);
            background: var(--bg-glass); backdrop-filter: blur(20px);
            border: 1px solid var(--gold); border-radius: var(--radius-lg);
            padding: 1.25rem; width: 500px; max-height: 70vh;
            overflow-y: auto; z-index: 50;
            box-shadow: 0 8px 32px rgba(0,0,0,0.5), 0 0 60px var(--gold-glow);
            display: none;
        }
        #hover-panel.visible { display: block; animation: panelIn 0.2s ease-out; }
        @keyframes panelIn { from { opacity: 0; transform: translateX(-50%%) translateY(-10px); } }
        .hover-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 1rem; }
        .hover-id { font-family: 'JetBrains Mono', monospace; font-size: 1rem; color: var(--cyan); font-weight: 600; }
        .hover-type-badge {
            padding: 0.25rem 0.625rem; border-radius: 6px;
            font-size: 0.7rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em;
        }
        .hover-title { font-size: 1.1rem; font-weight: 600; margin-bottom: 1rem; line-height: 1.4; }
        .hover-badges { display: flex; flex-wrap: wrap; gap: 0.375rem; margin-bottom: 1rem; }
        .hover-section { margin-bottom: 1rem; }
        .hover-section-title {
            font-size: 0.7rem; font-weight: 600; color: var(--purple);
            text-transform: uppercase; letter-spacing: 0.1em; margin-bottom: 0.5rem;
            padding-bottom: 0.25rem; border-bottom: 1px solid var(--bg-elevated);
        }
        .hover-content { font-size: 0.875rem; color: var(--fg); line-height: 1.6; }
        .hover-content p { margin-bottom: 0.5rem; }
        .hover-content ul, .hover-content ol { margin-left: 1.25rem; margin-bottom: 0.5rem; }
        .hover-content code {
            background: var(--bg); padding: 0.125rem 0.375rem; border-radius: 4px;
            font-family: 'JetBrains Mono', monospace; font-size: 0.8rem;
        }
        .hover-content pre {
            background: var(--bg); padding: 0.75rem; border-radius: 8px;
            overflow-x: auto; margin: 0.5rem 0;
        }
        .hover-content pre code { padding: 0; background: none; }
        .hover-meta { display: grid; grid-template-columns: repeat(2, 1fr); gap: 0.5rem; font-size: 0.8rem; }
        .hover-meta-item { display: flex; flex-direction: column; }
        .hover-meta-label { color: var(--fg-muted); font-size: 0.7rem; }
        .hover-meta-value { font-weight: 500; }
        .hover-deps { display: flex; flex-wrap: wrap; gap: 0.375rem; }
        .hover-dep-chip {
            padding: 0.25rem 0.5rem; background: var(--bg);
            border-radius: 6px; font-size: 0.75rem; font-family: 'JetBrains Mono', monospace;
            cursor: pointer; transition: all 0.15s ease;
        }
        .hover-dep-chip:hover { background: var(--purple); color: white; }
        .hover-commits { max-height: 150px; overflow-y: auto; }
        .hover-commit {
            padding: 0.5rem; background: var(--bg); border-radius: 6px;
            margin-bottom: 0.375rem; font-size: 0.75rem;
        }
        .hover-commit-sha { font-family: 'JetBrains Mono', monospace; color: var(--cyan); }
        .hover-commit-msg { margin-top: 0.25rem; color: var(--fg-muted); }
        .hover-close {
            position: absolute; top: 0.75rem; right: 0.75rem;
            background: none; border: none; color: var(--fg-muted);
            cursor: pointer; font-size: 1.25rem; padding: 0.25rem;
        }
        .hover-close:hover { color: var(--fg); }

        /* Sidebar */
        #sidebar {
            width: 340px;
            background: linear-gradient(180deg, var(--bg-secondary) 0%%, var(--bg) 100%%);
            border-left: 1px solid var(--purple);
            overflow-y: auto; padding: 1.25rem;
            display: flex; flex-direction: column; gap: 1rem;
        }
        .stats-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 0.625rem; }
        .stat-card {
            background: var(--bg-tertiary); padding: 0.875rem;
            border-radius: var(--radius); text-align: center;
            border: 1px solid var(--bg-elevated);
            transition: all 0.2s ease;
        }
        .stat-card:hover { border-color: var(--purple); transform: translateY(-2px); box-shadow: var(--shadow); }
        .stat-card .stat-value {
            font-size: 1.75rem; font-weight: 700;
            font-family: 'JetBrains Mono', monospace;
            background: linear-gradient(135deg, var(--green), var(--cyan));
            -webkit-background-clip: text; -webkit-text-fill-color: transparent;
        }
        .stat-card .stat-value.warning {
            background: linear-gradient(135deg, var(--orange), var(--red));
            -webkit-background-clip: text; -webkit-text-fill-color: transparent;
        }
        .stat-card .stat-label {
            font-size: 0.7rem; color: var(--fg-muted);
            text-transform: uppercase; letter-spacing: 0.05em; margin-top: 0.25rem;
        }

        /* Panels */
        .panel {
            background: var(--bg-tertiary); border-radius: var(--radius);
            padding: 1rem; border: 1px solid var(--bg-elevated);
        }
        .panel-title {
            font-size: 0.75rem; font-weight: 600; color: var(--purple);
            text-transform: uppercase; letter-spacing: 0.1em; margin-bottom: 0.75rem;
            display: flex; align-items: center; gap: 0.5rem;
        }
        .panel-title::before {
            content: ''; width: 4px; height: 14px;
            background: linear-gradient(180deg, var(--purple), var(--pink));
            border-radius: 2px;
        }

        /* Legend */
        .legend { display: flex; flex-wrap: wrap; gap: 0.5rem; }
        .legend-item {
            display: flex; align-items: center; gap: 0.375rem;
            font-size: 0.75rem; color: var(--fg-muted);
            padding: 0.25rem 0.5rem; background: var(--bg); border-radius: 6px;
        }
        .legend-dot {
            width: 12px; height: 12px; border-radius: 50%%;
            box-shadow: 0 0 8px currentColor;
        }

        /* Triage Panel */
        .triage-item {
            padding: 0.75rem; background: var(--bg); border-radius: 8px;
            margin-bottom: 0.5rem; cursor: pointer; transition: all 0.15s ease;
            border-left: 3px solid var(--purple);
        }
        .triage-item:hover { transform: translateX(4px); border-left-color: var(--gold); }
        .triage-item-header { display: flex; justify-content: space-between; align-items: center; }
        .triage-item-id { font-family: 'JetBrains Mono', monospace; font-size: 0.75rem; color: var(--cyan); font-weight: 600; }
        .triage-item-score { font-size: 0.7rem; color: var(--gold); font-weight: 600; }
        .triage-item-title { font-size: 0.8rem; margin-top: 0.25rem; }
        .triage-item-reason { font-size: 0.7rem; color: var(--fg-muted); margin-top: 0.375rem; }

        /* Badges */
        .badge {
            font-size: 0.65rem; padding: 0.2rem 0.5rem; border-radius: 4px;
            text-transform: uppercase; font-weight: 600; letter-spacing: 0.03em;
        }
        .badge-open { background: var(--green); color: var(--bg); }
        .badge-in_progress { background: var(--orange); color: var(--bg); }
        .badge-blocked { background: var(--red); color: white; }
        .badge-closed { background: var(--fg-dim); color: var(--bg); }
        .badge-type { background: var(--purple); color: white; }
        .badge-feature { background: linear-gradient(135deg, var(--purple), var(--pink)); color: white; }
        .badge-bug { background: var(--red); color: white; }
        .badge-task { background: var(--cyan); color: var(--bg); }
        .badge-epic { background: linear-gradient(135deg, var(--gold), var(--orange)); color: var(--bg); }
        .badge-articulation { background: linear-gradient(135deg, var(--pink), var(--purple)); color: white; animation: pulse 2s infinite; }
        .badge-critical { background: linear-gradient(135deg, var(--red), var(--orange)); color: white; }
        @keyframes pulse { 0%%, 100%% { opacity: 1; } 50%% { opacity: 0.7; } }

        /* Node Detail */
        #node-detail { display: none; }
        #node-detail.visible { display: block; }
        .detail-header { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 0.75rem; }
        .detail-id { font-size: 1rem; font-weight: 700; color: var(--cyan); font-family: 'JetBrains Mono', monospace; }
        .detail-priority {
            padding: 0.25rem 0.625rem; border-radius: 6px;
            font-size: 0.7rem; font-weight: 700; text-transform: uppercase;
        }
        .detail-name { font-size: 0.875rem; color: var(--fg); line-height: 1.5; margin-bottom: 0.625rem; }
        .detail-badges { display: flex; gap: 0.375rem; flex-wrap: wrap; margin-bottom: 0.75rem; }
        .detail-metrics { display: grid; grid-template-columns: repeat(2, 1fr); gap: 0.375rem; font-size: 0.75rem; }
        .metric-item {
            display: flex; justify-content: space-between;
            padding: 0.375rem 0; border-bottom: 1px solid var(--bg-elevated);
        }
        .metric-label { color: var(--fg-muted); }
        .metric-value { color: var(--fg); font-weight: 500; font-family: 'JetBrains Mono', monospace; }
        .metric-value.highlight { color: var(--green); }
        .no-selection { text-align: center; padding: 2rem 1rem; color: var(--fg-muted); font-size: 0.8rem; }
        .no-selection-icon { font-size: 2rem; margin-bottom: 0.625rem; opacity: 0.4; }

        /* Shortcuts */
        .keyboard-hints { font-size: 0.75rem; color: var(--fg-muted); line-height: 1.8; }
        .keyboard-hints kbd {
            display: inline-block; background: var(--bg);
            padding: 0.2rem 0.5rem; border-radius: 4px;
            margin: 0 0.125rem; border: 1px solid var(--bg-elevated);
            font-family: 'JetBrains Mono', monospace; font-size: 0.7rem;
        }

        /* Footer */
        footer {
            background: var(--bg-tertiary); padding: 0.5rem 1.25rem;
            font-size: 0.75rem; color: var(--fg-muted);
            display: flex; justify-content: space-between; align-items: center;
            border-top: 1px solid var(--bg-elevated);
        }
        footer a { color: var(--cyan); text-decoration: none; transition: color 0.15s; }
        footer a:hover { color: var(--purple); }

        /* Toast */
        .toast {
            position: fixed; bottom: 100px; left: 50%%; transform: translateX(-50%%);
            background: var(--bg-glass); backdrop-filter: blur(15px);
            border: 1px solid var(--purple); padding: 0.75rem 1.5rem;
            border-radius: var(--radius); font-size: 0.875rem;
            z-index: 1000; box-shadow: var(--shadow);
            opacity: 0; transition: all 0.3s ease;
        }
        .toast.visible { opacity: 1; transform: translateX(-50%%) translateY(-10px); }

        /* Context Menu */
        .context-menu {
            position: fixed; background: var(--bg-glass); backdrop-filter: blur(20px);
            border: 1px solid var(--purple); border-radius: var(--radius);
            padding: 0.5rem 0; z-index: 1000; min-width: 200px;
            box-shadow: var(--shadow); display: none;
        }
        .context-menu.visible { display: block; animation: menuIn 0.15s ease-out; }
        @keyframes menuIn { from { opacity: 0; transform: scale(0.95); } }
        .context-menu-item {
            padding: 0.625rem 1rem; font-size: 0.8rem; cursor: pointer;
            display: flex; align-items: center; gap: 0.625rem; transition: all 0.1s;
        }
        .context-menu-item:hover { background: var(--bg-elevated); }
        .context-menu-divider { height: 1px; background: var(--bg-elevated); margin: 0.375rem 0; }

        /* Fullscreen Button */
        .fullscreen-btn {
            position: absolute; top: 1rem; right: 1rem;
            background: var(--bg-glass); backdrop-filter: blur(10px);
            border: 1px solid var(--purple); color: var(--fg);
            padding: 0.5rem 0.75rem; border-radius: 8px;
            z-index: 10; cursor: pointer; font-size: 1rem;
            transition: all 0.2s ease;
        }
        .fullscreen-btn:hover { background: var(--bg-elevated); box-shadow: var(--shadow-glow); }

        /* Top Nodes Panel */
        .top-nodes-panel {
            position: absolute; top: 60px; right: 1rem;
            background: var(--bg-glass); backdrop-filter: blur(15px);
            border: 1px solid var(--purple); border-radius: var(--radius);
            padding: 0.75rem; z-index: 10; max-height: 250px;
            overflow-y: auto; width: 220px; display: none;
        }
        .top-nodes-panel.visible { display: block; }
        .top-node-item {
            padding: 0.5rem; font-size: 0.75rem; cursor: pointer;
            border-radius: 6px; display: flex; justify-content: space-between;
            transition: all 0.15s ease;
        }
        .top-node-item:hover { background: var(--bg-elevated); transform: translateX(4px); }
        .top-node-item .rank { color: var(--gold); font-weight: 600; }

        /* Heatmap Legend */
        .heatmap-legend {
            position: absolute; bottom: 1rem; right: 1rem;
            background: var(--bg-glass); backdrop-filter: blur(15px);
            border: 1px solid var(--purple); border-radius: var(--radius);
            padding: 0.75rem; z-index: 10; display: none;
        }
        .heatmap-legend.visible { display: block; }
        .heatmap-gradient {
            width: 140px; height: 14px;
            background: linear-gradient(90deg, var(--green), var(--yellow), var(--orange), var(--red));
            border-radius: 4px; margin-bottom: 0.375rem;
        }
        .heatmap-labels { display: flex; justify-content: space-between; font-size: 0.7rem; color: var(--fg-muted); }

        /* Scrollbar */
        ::-webkit-scrollbar { width: 8px; height: 8px; }
        ::-webkit-scrollbar-track { background: var(--bg); }
        ::-webkit-scrollbar-thumb { background: var(--bg-elevated); border-radius: 4px; }
        ::-webkit-scrollbar-thumb:hover { background: var(--purple); }
    </style>
</head>
<body>
    <header>
        <div class="logo">
            <div class="logo-icon">bv</div>
            <h1><span>%s</span> Graph</h1>
        </div>
        <div class="toolbar">
            <div class="search-container">
                <span class="search-icon">🔍</span>
                <input type="text" class="search-input" id="search-input" placeholder="Search beads... (full text)">
                <div class="search-results" id="search-results"></div>
            </div>
            <div class="toolbar-group">
                <select id="view-mode">
                    <option value="force">Force</option>
                    <option value="td">DAG ↓</option>
                    <option value="lr">DAG →</option>
                    <option value="radialout">Radial</option>
                </select>
            </div>
            <div class="toolbar-group">
                <select id="filter-status">
                    <option value="">All Status</option>
                    <option value="open">Open</option>
                    <option value="in_progress">In Progress</option>
                    <option value="blocked">Blocked</option>
                    <option value="closed">Closed</option>
                </select>
                <select id="filter-type">
                    <option value="">All Types</option>
                    <option value="feature">Feature</option>
                    <option value="bug">Bug</option>
                    <option value="task">Task</option>
                    <option value="epic">Epic</option>
                </select>
            </div>
            <div class="toolbar-group">
                <select id="size-by">
                    <option value="pagerank">Size: PageRank</option>
                    <option value="betweenness">Size: Betweenness</option>
                    <option value="critical">Size: Critical Path</option>
                    <option value="indegree">Size: In-Degree</option>
                </select>
            </div>
            <div class="toolbar-group">
                <button id="btn-heatmap" title="Toggle Heatmap">🔥</button>
                <button id="btn-triage" title="Triage View">📋</button>
                <button id="btn-top" title="Top Nodes">⭐</button>
                <button id="btn-fit" title="Fit (F)">Fit</button>
                <button id="btn-reset" title="Reset (R)">Reset</button>
            </div>
        </div>
    </header>
    <main>
        <div id="graph-wrapper">
            <div id="graph-container"></div>
            <div class="overlay-stats">
                <div class="stat"><span class="stat-value">%d</span> nodes</div>
                <div class="stat"><span class="stat-value">%d</span> edges</div>
                <div class="stat" id="stat-visible"><span class="stat-value">%d</span> visible</div>
            </div>
            <button class="fullscreen-btn" id="btn-fullscreen" title="Fullscreen (Space)">⛶</button>
            <div class="top-nodes-panel" id="top-nodes-panel"></div>
            <div class="heatmap-legend" id="heatmap-legend">
                <div class="heatmap-gradient"></div>
                <div class="heatmap-labels"><span>Low</span><span id="heatmap-metric">PageRank</span><span>High</span></div>
            </div>
            <div id="hover-panel">
                <button class="hover-close" id="hover-close">×</button>
                <div class="hover-header">
                    <span class="hover-id" id="hover-id">-</span>
                    <span class="hover-type-badge" id="hover-type-badge">-</span>
                </div>
                <div class="hover-title" id="hover-title">-</div>
                <div class="hover-badges" id="hover-badges"></div>
                <div id="hover-description" class="hover-section" style="display:none;">
                    <div class="hover-section-title">Description</div>
                    <div class="hover-content" id="hover-description-content"></div>
                </div>
                <div id="hover-design" class="hover-section" style="display:none;">
                    <div class="hover-section-title">Design</div>
                    <div class="hover-content" id="hover-design-content"></div>
                </div>
                <div id="hover-acceptance" class="hover-section" style="display:none;">
                    <div class="hover-section-title">Acceptance Criteria</div>
                    <div class="hover-content" id="hover-acceptance-content"></div>
                </div>
                <div id="hover-notes" class="hover-section" style="display:none;">
                    <div class="hover-section-title">Notes</div>
                    <div class="hover-content" id="hover-notes-content"></div>
                </div>
                <div class="hover-section">
                    <div class="hover-section-title">Metadata</div>
                    <div class="hover-meta" id="hover-meta"></div>
                </div>
                <div id="hover-blocked-by" class="hover-section" style="display:none;">
                    <div class="hover-section-title">Blocked By</div>
                    <div class="hover-deps" id="hover-blocked-by-list"></div>
                </div>
                <div id="hover-blocks" class="hover-section" style="display:none;">
                    <div class="hover-section-title">Blocks</div>
                    <div class="hover-deps" id="hover-blocks-list"></div>
                </div>
                <div id="hover-commits" class="hover-section" style="display:none;">
                    <div class="hover-section-title">Related Commits</div>
                    <div class="hover-commits" id="hover-commits-list"></div>
                </div>
                <div class="hover-section">
                    <div class="hover-section-title">Graph Metrics</div>
                    <div class="hover-meta" id="hover-metrics"></div>
                </div>
            </div>
        </div>
        <div id="sidebar">
            <div class="stats-grid">
                <div class="stat-card"><div class="stat-value" id="stat-nodes">%d</div><div class="stat-label">Nodes</div></div>
                <div class="stat-card"><div class="stat-value" id="stat-edges">%d</div><div class="stat-label">Edges</div></div>
                <div class="stat-card"><div class="stat-value" id="stat-actionable">-</div><div class="stat-label">Actionable</div></div>
                <div class="stat-card"><div class="stat-value warning" id="stat-blocked">-</div><div class="stat-label">Blocked</div></div>
                <div class="stat-card"><div class="stat-value" id="stat-critical">-</div><div class="stat-label">Critical</div></div>
                <div class="stat-card"><div class="stat-value warning" id="stat-articulation">-</div><div class="stat-label">Cut Pts</div></div>
            </div>
            <div class="panel" id="triage-panel" style="display:none;">
                <div class="panel-title">Top Recommendations</div>
                <div id="triage-list"></div>
            </div>
            <div class="panel">
                <div class="panel-title">Status Legend</div>
                <div class="legend">
                    <div class="legend-item"><div class="legend-dot" style="background:#22c55e;color:#22c55e"></div>Open</div>
                    <div class="legend-item"><div class="legend-dot" style="background:#f97316;color:#f97316"></div>In Progress</div>
                    <div class="legend-item"><div class="legend-dot" style="background:#ef4444;color:#ef4444"></div>Blocked</div>
                    <div class="legend-item"><div class="legend-dot" style="background:#555577;color:#555577"></div>Closed</div>
                </div>
            </div>
            <div class="panel">
                <div class="panel-title">Type Shapes</div>
                <div class="legend">
                    <div class="legend-item"><span style="font-size:1rem">●</span> Feature</div>
                    <div class="legend-item"><span style="font-size:1rem">▲</span> Bug</div>
                    <div class="legend-item"><span style="font-size:1rem">■</span> Task</div>
                    <div class="legend-item"><span style="font-size:1rem">◆</span> Epic</div>
                </div>
            </div>
            <div class="panel">
                <div class="panel-title">Selected Node</div>
                <div id="node-detail">
                    <div class="detail-header">
                        <div class="detail-id" id="detail-id">-</div>
                        <div class="detail-priority" id="detail-priority">P2</div>
                    </div>
                    <div class="detail-name" id="detail-name">-</div>
                    <div class="detail-badges" id="detail-badges"></div>
                    <div class="detail-metrics">
                        <div class="metric-item"><span class="metric-label">PageRank</span><span class="metric-value" id="m-pagerank">-</span></div>
                        <div class="metric-item"><span class="metric-label">Rank</span><span class="metric-value" id="m-prrank">-</span></div>
                        <div class="metric-item"><span class="metric-label">Betweenness</span><span class="metric-value" id="m-between">-</span></div>
                        <div class="metric-item"><span class="metric-label">BW Rank</span><span class="metric-value" id="m-bwrank">-</span></div>
                        <div class="metric-item"><span class="metric-label">Critical</span><span class="metric-value" id="m-critical">-</span></div>
                        <div class="metric-item"><span class="metric-label">Slack</span><span class="metric-value" id="m-slack">-</span></div>
                        <div class="metric-item"><span class="metric-label">In-Deg</span><span class="metric-value" id="m-indeg">-</span></div>
                        <div class="metric-item"><span class="metric-label">Out-Deg</span><span class="metric-value" id="m-outdeg">-</span></div>
                    </div>
                </div>
                <div class="no-selection" id="no-selection">
                    <div class="no-selection-icon">🔍</div>
                    Click a node to see details<br>
                    <small>or hover for full info</small>
                </div>
            </div>
            <div class="panel">
                <div class="panel-title">Shortcuts</div>
                <div class="keyboard-hints">
                    <kbd>F</kbd> Fit · <kbd>R</kbd> Reset · <kbd>Space</kbd> Fullscreen<br>
                    <kbd>Esc</kbd> Clear · <kbd>1-4</kbd> View modes<br>
                    <kbd>H</kbd> Heatmap · <kbd>T</kbd> Top · <kbd>G</kbd> Triage
                </div>
            </div>
        </div>
    </main>
    <footer>
        <div>Generated %s | Hash: %s</div>
        <div>Project: %s | <a href="https://github.com/Dicklesworthstone/beads_viewer">bv</a></div>
    </footer>
    <div class="toast" id="toast"></div>
    <div class="context-menu" id="context-menu">
        <div class="context-menu-item" id="ctx-focus">🎯 Focus on this node</div>
        <div class="context-menu-item" id="ctx-details">📄 Show full details</div>
        <div class="context-menu-item" id="ctx-deps">📥 Show dependencies</div>
        <div class="context-menu-item" id="ctx-dependents">📤 Show dependents</div>
        <div class="context-menu-item" id="ctx-connected">✨ Highlight connected</div>
        <div class="context-menu-divider"></div>
        <div class="context-menu-item" id="ctx-path">🛤️ Find path to...</div>
        <div class="context-menu-item" id="ctx-copy">📋 Copy ID</div>
    </div>
    <script>%s</script>
    <script>%s</script>
    <script>
const DATA = %s;
const STATUS_COLORS = { open: '#22c55e', in_progress: '#f97316', blocked: '#ef4444', closed: '#555577' };
const PRIORITY_COLORS = ['#ef4444', '#f97316', '#eab308', '#22c55e', '#555577'];
const TYPE_COLORS = { feature: '#a855f7', bug: '#ef4444', task: '#22d3ee', epic: '#fbbf24' };

// Configure marked for safe HTML rendering
marked.setOptions({ breaks: true, gfm: true });

// Stats calculation
let actionable = 0, blocked = 0, onCriticalPath = 0, articulationCount = 0;
const blockerCount = {};
DATA.links.forEach(l => blockerCount[l.source] = (blockerCount[l.source] || 0) + 1);
DATA.nodes.forEach(n => {
    n.blockerCount = blockerCount[n.id] || 0;
    if ((n.status === 'open' || n.status === 'in_progress') && n.blockerCount === 0) actionable++;
    if (n.status === 'blocked') blocked++;
    if (n.slack === 0) onCriticalPath++;
    if (n.is_articulation) articulationCount++;
});
document.getElementById('stat-actionable').textContent = actionable;
document.getElementById('stat-blocked').textContent = blocked;
document.getElementById('stat-critical').textContent = onCriticalPath;
document.getElementById('stat-articulation').textContent = articulationCount;

// Max values for sizing
const maxPR = Math.max(...DATA.nodes.map(n => n.pagerank || 0), 0.001);
const maxBW = Math.max(...DATA.nodes.map(n => n.betweenness || 0), 0.001);
const maxCP = Math.max(...DATA.nodes.map(n => n.critical_path || 0), 1);
const maxInDeg = Math.max(...DATA.nodes.map(n => n.in_degree || 0), 1);

let sizeMetric = 'pagerank', heatmapMode = false, hoveredNode = null, highlightedNodes = new Set();

function getNodeSize(n) {
    const base = 5, scale = 16;
    switch(sizeMetric) {
        case 'pagerank': return base + ((n.pagerank || 0) / maxPR) * scale;
        case 'betweenness': return base + ((n.betweenness || 0) / maxBW) * scale;
        case 'critical': return base + ((n.critical_path || 0) / maxCP) * scale;
        case 'indegree': return base + ((n.in_degree || 0) / maxInDeg) * scale;
        default: return base + ((n.pagerank || 0) / maxPR) * scale;
    }
}

function getHeatmapColor(n) {
    let val = 0, max = 1;
    switch(sizeMetric) {
        case 'pagerank': val = n.pagerank || 0; max = maxPR; break;
        case 'betweenness': val = n.betweenness || 0; max = maxBW; break;
        case 'critical': val = n.critical_path || 0; max = maxCP; break;
        case 'indegree': val = n.in_degree || 0; max = maxInDeg; break;
    }
    const ratio = val / max;
    const hue = 120 - ratio * 120; // Green to red
    return 'hsl(' + hue + ', 80%%, 50%%)';
}

// Get connected subgraph (for golden glow highlight)
function getConnectedNodes(nodeId, depth = 2) {
    const connected = new Set([nodeId]);
    const queue = [{id: nodeId, d: 0}];
    while (queue.length > 0) {
        const {id, d} = queue.shift();
        if (d >= depth) continue;
        DATA.links.forEach(l => {
            const src = typeof l.source === 'object' ? l.source.id : l.source;
            const tgt = typeof l.target === 'object' ? l.target.id : l.target;
            if (src === id && !connected.has(tgt)) { connected.add(tgt); queue.push({id: tgt, d: d+1}); }
            if (tgt === id && !connected.has(src)) { connected.add(src); queue.push({id: src, d: d+1}); }
        });
    }
    return connected;
}

const container = document.getElementById('graph-container');
const Graph = ForceGraph()(container)
    .graphData(JSON.parse(JSON.stringify(DATA)))
    .backgroundColor('transparent')
    .nodeId('id')
    .nodeLabel(null)
    .nodeColor(n => {
        if (highlightedNodes.size > 0 && !highlightedNodes.has(n.id)) return (STATUS_COLORS[n.status] || '#555577') + '20';
        if (heatmapMode) return getHeatmapColor(n);
        return STATUS_COLORS[n.status] || '#555577';
    })
    .nodeVal(n => getNodeSize(n))
    .linkColor(l => {
        const src = typeof l.source === 'object' ? l.source.id : l.source;
        const tgt = typeof l.target === 'object' ? l.target.id : l.target;
        if (highlightedNodes.size > 0) {
            if (highlightedNodes.has(src) && highlightedNodes.has(tgt)) return '#fbbf24aa';
            return '#44475a15';
        }
        return l.critical ? '#ec489980' : '#44475a40';
    })
    .linkWidth(l => {
        const src = typeof l.source === 'object' ? l.source.id : l.source;
        const tgt = typeof l.target === 'object' ? l.target.id : l.target;
        if (highlightedNodes.size > 0 && highlightedNodes.has(src) && highlightedNodes.has(tgt)) return 3;
        return l.critical ? 2 : 1;
    })
    .linkDirectionalArrowLength(5)
    .linkDirectionalArrowColor(l => {
        const src = typeof l.source === 'object' ? l.source.id : l.source;
        const tgt = typeof l.target === 'object' ? l.target.id : l.target;
        if (highlightedNodes.size > 0 && highlightedNodes.has(src) && highlightedNodes.has(tgt)) return '#fbbf24';
        return l.critical ? '#ec4899' : '#44475a';
    })
    .linkDirectionalArrowRelPos(1)
    .linkCurvature(0.1)
    .linkDirectionalParticles(l => l.critical ? 2 : 0)
    .linkDirectionalParticleSpeed(0.003)
    .linkDirectionalParticleWidth(2)
    .linkDirectionalParticleColor(() => '#ec4899')
    .d3AlphaDecay(0.02)
    .d3VelocityDecay(0.25)
    .nodeCanvasObject((node, ctx, globalScale) => {
        const x = node.x, y = node.y;
        if (x === undefined || y === undefined || !isFinite(x) || !isFinite(y)) return;
        const size = getNodeSize(node);
        const baseColor = heatmapMode ? getHeatmapColor(node) : STATUS_COLORS[node.status] || '#555577';
        const isHighlighted = highlightedNodes.size === 0 || highlightedNodes.has(node.id);
        const isHovered = hoveredNode && hoveredNode.id === node.id;
        const alpha = isHighlighted ? 1 : 0.15;

        // Golden glow for hovered node's connected subgraph
        if (isHovered || (highlightedNodes.has(node.id) && highlightedNodes.size > 0)) {
            ctx.beginPath(); ctx.arc(x, y, size + 8, 0, 2 * Math.PI);
            const g = ctx.createRadialGradient(x, y, size, x, y, size + 12);
            g.addColorStop(0, 'rgba(251, 191, 36, 0.6)'); g.addColorStop(1, 'transparent');
            ctx.fillStyle = g; ctx.fill();
        }

        // Articulation point glow
        if (node.is_articulation && isHighlighted) {
            ctx.beginPath(); ctx.arc(x, y, size + 6, 0, 2 * Math.PI);
            const g = ctx.createRadialGradient(x, y, size, x, y, size + 8);
            g.addColorStop(0, '#ec489960'); g.addColorStop(1, 'transparent');
            ctx.fillStyle = g; ctx.fill();
        }

        // Critical path indicator
        if (node.slack === 0 && isHighlighted) {
            ctx.beginPath(); ctx.arc(x, y, size + 3, 0, 2 * Math.PI);
            ctx.fillStyle = baseColor + '30'; ctx.fill();
        }

        // Priority ring
        const pColor = PRIORITY_COLORS[node.priority] || PRIORITY_COLORS[2];
        ctx.globalAlpha = alpha;
        ctx.beginPath(); ctx.arc(x, y, size + 1.5, 0, 2 * Math.PI);
        ctx.strokeStyle = pColor; ctx.lineWidth = 2; ctx.stroke();

        // Node shape based on type
        ctx.fillStyle = baseColor;
        ctx.beginPath();
        switch(node.type) {
            case 'bug': // Triangle
                ctx.moveTo(x, y - size);
                ctx.lineTo(x + size * 0.866, y + size * 0.5);
                ctx.lineTo(x - size * 0.866, y + size * 0.5);
                ctx.closePath();
                break;
            case 'task': // Square
                ctx.rect(x - size * 0.7, y - size * 0.7, size * 1.4, size * 1.4);
                break;
            case 'epic': // Diamond
                ctx.moveTo(x, y - size);
                ctx.lineTo(x + size, y);
                ctx.lineTo(x, y + size);
                ctx.lineTo(x - size, y);
                ctx.closePath();
                break;
            default: // Circle (feature)
                ctx.arc(x, y, size, 0, 2 * Math.PI);
        }
        ctx.fill();

        // Highlight effect
        const hl = ctx.createRadialGradient(x - size/3, y - size/3, 0, x, y, size);
        hl.addColorStop(0, 'rgba(255,255,255,0.2)'); hl.addColorStop(1, 'transparent');
        ctx.fillStyle = hl; ctx.fill();
        ctx.globalAlpha = 1;

        // Labels at zoom
        if (globalScale > 1.2 && isHighlighted) {
            const fontSize = Math.max(10 / globalScale, 3);
            ctx.font = fontSize + 'px Inter, sans-serif';
            ctx.textAlign = 'center'; ctx.textBaseline = 'middle';
            ctx.fillStyle = '#e8e8f0';
            ctx.fillText(node.id, x, y + size + fontSize + 2);
            if (globalScale > 2) {
                ctx.fillStyle = pColor;
                ctx.font = (fontSize * 0.8) + 'px JetBrains Mono, monospace';
                ctx.fillText('P' + node.priority, x, y);
            }
        }
    })
    .nodePointerAreaPaint((n, c, ctx) => {
        if (n.x === undefined || n.y === undefined || !isFinite(n.x) || !isFinite(n.y)) return;
        const size = getNodeSize(n) + 4;
        ctx.fillStyle = c; ctx.beginPath(); ctx.arc(n.x, n.y, size, 0, 2 * Math.PI); ctx.fill();
    })
    .onNodeClick(handleNodeClick)
    .onNodeRightClick((node, event) => { event.preventDefault(); showContextMenu(node, event); })
    .onNodeHover(handleNodeHover)
    .onBackgroundClick(() => { clearSelection(); hideContextMenu(); hideHoverPanel(); })
    .onBackgroundRightClick(() => hideContextMenu());

// Hover handling with golden glow
function handleNodeHover(node) {
    hoveredNode = node;
    container.style.cursor = node ? 'pointer' : 'grab';
    if (node) {
        highlightedNodes = getConnectedNodes(node.id, 2);
    } else {
        highlightedNodes = new Set();
    }
    Graph.nodeColor(Graph.nodeColor()); // Trigger re-render
    Graph.linkColor(Graph.linkColor());
    Graph.linkWidth(Graph.linkWidth());
}

// Show full details hover panel
function showHoverPanel(node) {
    const panel = document.getElementById('hover-panel');
    document.getElementById('hover-id').textContent = node.id;
    document.getElementById('hover-title').textContent = node.title;

    // Type badge
    const typeBadge = document.getElementById('hover-type-badge');
    typeBadge.textContent = node.type || 'task';
    typeBadge.className = 'hover-type-badge badge-' + (node.type || 'task');

    // Badges
    const badgesEl = document.getElementById('hover-badges');
    badgesEl.innerHTML = '';
    const addBadge = (cls, text) => { const b = document.createElement('span'); b.className = 'badge ' + cls; b.textContent = text; badgesEl.appendChild(b); };
    addBadge('badge-' + node.status, node.status.replace('_', ' '));
    addBadge('', 'P' + node.priority);
    if (node.is_articulation) addBadge('badge-articulation', 'Cut Vertex');
    if (node.slack === 0) addBadge('badge-critical', 'Critical Path');
    (node.labels || []).forEach(l => addBadge('badge-type', l));

    // Description
    const descSection = document.getElementById('hover-description');
    if (node.description) {
        descSection.style.display = 'block';
        document.getElementById('hover-description-content').innerHTML = marked.parse(node.description);
    } else { descSection.style.display = 'none'; }

    // Design
    const designSection = document.getElementById('hover-design');
    if (node.design) {
        designSection.style.display = 'block';
        document.getElementById('hover-design-content').innerHTML = marked.parse(node.design);
    } else { designSection.style.display = 'none'; }

    // Acceptance Criteria
    const acSection = document.getElementById('hover-acceptance');
    if (node.acceptance_criteria) {
        acSection.style.display = 'block';
        document.getElementById('hover-acceptance-content').innerHTML = marked.parse(node.acceptance_criteria);
    } else { acSection.style.display = 'none'; }

    // Notes
    const notesSection = document.getElementById('hover-notes');
    if (node.notes) {
        notesSection.style.display = 'block';
        document.getElementById('hover-notes-content').innerHTML = marked.parse(node.notes);
    } else { notesSection.style.display = 'none'; }

    // Metadata
    const metaEl = document.getElementById('hover-meta');
    metaEl.innerHTML = '';
    const addMeta = (label, value) => {
        if (!value) return;
        metaEl.innerHTML += '<div class="hover-meta-item"><span class="hover-meta-label">' + label + '</span><span class="hover-meta-value">' + value + '</span></div>';
    };
    addMeta('Assignee', node.assignee);
    addMeta('Created', node.created_at);
    addMeta('Updated', node.updated_at);
    addMeta('Due Date', node.due_date);
    if (node.closed_at) addMeta('Closed', node.closed_at);

    // Blocked By
    const blockedBySection = document.getElementById('hover-blocked-by');
    const blockedByList = document.getElementById('hover-blocked-by-list');
    if (node.blocked_by && node.blocked_by.length > 0) {
        blockedBySection.style.display = 'block';
        blockedByList.innerHTML = node.blocked_by.map(id => '<span class="hover-dep-chip" data-id="' + id + '">' + id + '</span>').join('');
    } else { blockedBySection.style.display = 'none'; }

    // Blocks
    const blocksSection = document.getElementById('hover-blocks');
    const blocksList = document.getElementById('hover-blocks-list');
    if (node.blocks && node.blocks.length > 0) {
        blocksSection.style.display = 'block';
        blocksList.innerHTML = node.blocks.map(id => '<span class="hover-dep-chip" data-id="' + id + '">' + id + '</span>').join('');
    } else { blocksSection.style.display = 'none'; }

    // Commits
    const commitsSection = document.getElementById('hover-commits');
    const commitsList = document.getElementById('hover-commits-list');
    if (node.commits && node.commits.length > 0) {
        commitsSection.style.display = 'block';
        commitsList.innerHTML = node.commits.slice(0, 5).map(c => '<div class="hover-commit"><span class="hover-commit-sha">' + c.short_sha + '</span> <span class="hover-commit-msg">' + (c.message || '').split('\\n')[0].substring(0, 60) + '</span></div>').join('');
    } else { commitsSection.style.display = 'none'; }

    // Metrics
    const metricsEl = document.getElementById('hover-metrics');
    metricsEl.innerHTML = '';
    const addMetric = (label, value) => {
        metricsEl.innerHTML += '<div class="hover-meta-item"><span class="hover-meta-label">' + label + '</span><span class="hover-meta-value">' + value + '</span></div>';
    };
    const fmt = (v, d) => (v != null && isFinite(v)) ? v.toFixed(d) : '-';
    addMetric('PageRank', fmt(node.pagerank * 100, 3) + '%%');
    addMetric('PR Rank', '#' + (node.pagerank_rank || '-'));
    addMetric('Betweenness', fmt(node.betweenness, 4));
    addMetric('BW Rank', '#' + (node.betweenness_rank || '-'));
    addMetric('Critical Path', fmt(node.critical_path, 1));
    addMetric('Slack', fmt(node.slack, 1));
    addMetric('In-Degree', node.in_degree ?? '-');
    addMetric('Out-Degree', node.out_degree ?? '-');

    panel.classList.add('visible');

    // Wire up dep chip clicks
    panel.querySelectorAll('.hover-dep-chip').forEach(chip => {
        chip.onclick = () => {
            const targetId = chip.dataset.id;
            const graphNodes = Graph.graphData().nodes;
            const target = graphNodes.find(n => n.id === targetId);
            if (target) {
                Graph.centerAt(target.x, target.y, 500);
                Graph.zoom(2.5, 500);
                setTimeout(() => showHoverPanel(target), 600);
            }
        };
    });
}

function hideHoverPanel() {
    document.getElementById('hover-panel').classList.remove('visible');
}

document.getElementById('hover-close').onclick = hideHoverPanel;

let selectedNode = null;
function selectNode(node) {
    selectedNode = node;
    showHoverPanel(node);
    document.getElementById('detail-id').textContent = node.id;
    document.getElementById('detail-name').textContent = node.title;
    const prioEl = document.getElementById('detail-priority');
    prioEl.textContent = 'P' + node.priority;
    prioEl.style.background = PRIORITY_COLORS[node.priority];
    prioEl.style.color = node.priority <= 1 ? 'white' : '#0f0f1a';
    const badgesEl = document.getElementById('detail-badges');
    badgesEl.innerHTML = '';
    const sb = document.createElement('span'); sb.className = 'badge badge-' + node.status; sb.textContent = node.status.replace('_', ' '); badgesEl.appendChild(sb);
    const tb = document.createElement('span'); tb.className = 'badge badge-' + (node.type || 'task'); tb.textContent = node.type || 'task'; badgesEl.appendChild(tb);
    const fmtSide = (v, d) => (v != null && isFinite(v)) ? v.toFixed(d) : '-';
    document.getElementById('m-pagerank').textContent = fmtSide(node.pagerank * 100, 2) + '%%';
    document.getElementById('m-prrank').textContent = '#' + (node.pagerank_rank || '-');
    document.getElementById('m-between').textContent = fmtSide(node.betweenness, 4);
    document.getElementById('m-bwrank').textContent = '#' + (node.betweenness_rank || '-');
    document.getElementById('m-critical').textContent = fmtSide(node.critical_path, 1);
    const slackEl = document.getElementById('m-slack');
    slackEl.textContent = fmtSide(node.slack, 1);
    slackEl.className = 'metric-value' + (node.slack === 0 ? ' highlight' : '');
    document.getElementById('m-indeg').textContent = node.in_degree ?? '-';
    document.getElementById('m-outdeg').textContent = node.out_degree ?? '-';
    document.getElementById('node-detail').classList.add('visible');
    document.getElementById('no-selection').style.display = 'none';
}

function clearSelection() {
    selectedNode = null;
    highlightedNodes = new Set();
    document.getElementById('node-detail').classList.remove('visible');
    document.getElementById('no-selection').style.display = 'block';
    Graph.nodeColor(Graph.nodeColor());
    Graph.linkColor(Graph.linkColor());
}

// Full-text search
let searchDebounce = null;
document.getElementById('search-input').oninput = e => {
    clearTimeout(searchDebounce);
    searchDebounce = setTimeout(() => performSearch(e.target.value), 150);
};

function performSearch(query) {
    const resultsEl = document.getElementById('search-results');
    if (!query || query.length < 2) {
        resultsEl.classList.remove('visible');
        applyFilters();
        return;
    }
    const q = query.toLowerCase();
    const matches = DATA.nodes.filter(n => {
        return n.id.toLowerCase().includes(q) ||
               n.title.toLowerCase().includes(q) ||
               (n.description || '').toLowerCase().includes(q) ||
               (n.design || '').toLowerCase().includes(q) ||
               (n.notes || '').toLowerCase().includes(q) ||
               (n.acceptance_criteria || '').toLowerCase().includes(q) ||
               (n.labels || []).some(l => l.toLowerCase().includes(q)) ||
               (n.assignee || '').toLowerCase().includes(q);
    }).slice(0, 8);

    if (matches.length === 0) {
        resultsEl.innerHTML = '<div class="search-result-item">No results found</div>';
    } else {
        resultsEl.innerHTML = matches.map(n => {
            let preview = '';
            const fields = [n.description, n.design, n.notes, n.acceptance_criteria];
            for (const f of fields) {
                if (f && f.toLowerCase().includes(q)) {
                    const idx = f.toLowerCase().indexOf(q);
                    const start = Math.max(0, idx - 30);
                    const end = Math.min(f.length, idx + q.length + 50);
                    preview = '...' + f.substring(start, end).replace(new RegExp(q, 'gi'), '<mark>$&</mark>') + '...';
                    break;
                }
            }
            return '<div class="search-result-item" data-id="' + n.id + '">' +
                   '<div class="search-result-id">' + n.id + ' <span class="badge badge-' + n.status + '">' + n.status + '</span></div>' +
                   '<div class="search-result-title">' + n.title + '</div>' +
                   (preview ? '<div class="search-result-preview">' + preview + '</div>' : '') +
                   '</div>';
        }).join('');
    }
    resultsEl.classList.add('visible');

    resultsEl.querySelectorAll('.search-result-item[data-id]').forEach(item => {
        item.onclick = () => {
            const id = item.dataset.id;
            const graphNodes = Graph.graphData().nodes;
            const node = graphNodes.find(n => n.id === id);
            if (node) {
                selectNode(node);
                Graph.centerAt(node.x, node.y, 500);
                Graph.zoom(2.5, 500);
            }
            resultsEl.classList.remove('visible');
            document.getElementById('search-input').value = '';
        };
    });
}

// Close search on click outside
document.addEventListener('click', e => {
    if (!e.target.closest('.search-container')) {
        document.getElementById('search-results').classList.remove('visible');
    }
});

// Context menu
let contextNode = null;
function showContextMenu(node, event) {
    contextNode = node;
    const menu = document.getElementById('context-menu');
    menu.style.left = event.clientX + 'px';
    menu.style.top = event.clientY + 'px';
    menu.classList.add('visible');
}
function hideContextMenu() { document.getElementById('context-menu').classList.remove('visible'); contextNode = null; }
document.getElementById('ctx-focus').onclick = () => { if (contextNode) { Graph.centerAt(contextNode.x, contextNode.y, 500); Graph.zoom(3, 500); } hideContextMenu(); };
document.getElementById('ctx-details').onclick = () => { if (contextNode) showHoverPanel(contextNode); hideContextMenu(); };
document.getElementById('ctx-deps').onclick = () => { if (contextNode) highlightDependencies(contextNode, 'deps'); hideContextMenu(); };
document.getElementById('ctx-dependents').onclick = () => { if (contextNode) highlightDependencies(contextNode, 'dependents'); hideContextMenu(); };
document.getElementById('ctx-connected').onclick = () => {
    if (contextNode) {
        highlightedNodes = getConnectedNodes(contextNode.id, 3);
        Graph.nodeColor(Graph.nodeColor());
        Graph.linkColor(Graph.linkColor());
        showToast(highlightedNodes.size + ' connected nodes highlighted');
    }
    hideContextMenu();
};
document.getElementById('ctx-copy').onclick = () => { if (contextNode) { navigator.clipboard.writeText(contextNode.id); showToast('Copied: ' + contextNode.id); } hideContextMenu(); };
document.getElementById('ctx-path').onclick = () => { showToast('Click another node to find path'); pathStartNode = contextNode; hideContextMenu(); };

let pathStartNode = null;
function findPath(startId, endId) {
    const queue = [[startId]];
    const visited = new Set([startId]);
    while (queue.length > 0) {
        const path = queue.shift();
        const current = path[path.length - 1];
        if (current === endId) return path;
        DATA.links.forEach(l => {
            const src = typeof l.source === 'object' ? l.source.id : l.source;
            const tgt = typeof l.target === 'object' ? l.target.id : l.target;
            if (src === current && !visited.has(tgt)) { visited.add(tgt); queue.push([...path, tgt]); }
            if (tgt === current && !visited.has(src)) { visited.add(src); queue.push([...path, src]); }
        });
    }
    return null;
}

function highlightPath(path) {
    highlightedNodes = new Set(path);
    Graph.nodeColor(Graph.nodeColor());
    Graph.linkColor(Graph.linkColor());
    updateVisibleCount();
    showToast('Path: ' + path.length + ' nodes');
}

function handleNodeClick(node) {
    if (pathStartNode) {
        const path = findPath(pathStartNode.id, node.id);
        if (path) highlightPath(path);
        else showToast('No path found');
        pathStartNode = null;
    } else {
        selectNode(node);
    }
}

function highlightDependencies(node, type) {
    const connected = new Set([node.id]);
    DATA.links.forEach(l => {
        const src = typeof l.source === 'object' ? l.source.id : l.source;
        const tgt = typeof l.target === 'object' ? l.target.id : l.target;
        if (type === 'deps' && src === node.id) connected.add(tgt);
        if (type === 'dependents' && tgt === node.id) connected.add(src);
    });
    highlightedNodes = connected;
    Graph.nodeColor(Graph.nodeColor());
    Graph.linkColor(Graph.linkColor());
    updateVisibleCount();
    showToast(connected.size + ' nodes shown');
}

// Filters
let statusFilter = '', typeFilter = '';
let currentVisibilityFilter = () => true;
document.getElementById('filter-status').onchange = e => { statusFilter = e.target.value; applyFilters(); };
document.getElementById('filter-type').onchange = e => { typeFilter = e.target.value; applyFilters(); };

function applyFilters() {
    const searchVal = document.getElementById('search-input').value.toLowerCase();
    currentVisibilityFilter = n => {
        const matchSearch = !searchVal || n.id.toLowerCase().includes(searchVal) || n.title.toLowerCase().includes(searchVal);
        const matchStatus = !statusFilter || n.status === statusFilter;
        const matchType = !typeFilter || n.type === typeFilter;
        return matchSearch && matchStatus && matchType;
    };
    Graph.nodeVisibility(currentVisibilityFilter);
    updateVisibleCount();
}

function updateVisibleCount() {
    const count = DATA.nodes.filter(n => currentVisibilityFilter(n)).length;
    document.getElementById('stat-visible').innerHTML = '<span class="stat-value">' + count + '</span> visible';
}

// View mode
document.getElementById('view-mode').onchange = e => {
    const mode = e.target.value;
    Graph.dagMode(mode === 'force' ? null : mode);
    setTimeout(() => Graph.zoomToFit(400), 100);
};

// Size metric
document.getElementById('size-by').onchange = e => {
    sizeMetric = e.target.value;
    document.getElementById('heatmap-metric').textContent = { pagerank: 'PageRank', betweenness: 'Betweenness', critical: 'Critical Path', indegree: 'In-Degree' }[sizeMetric];
    Graph.nodeVal(n => getNodeSize(n));
    if (heatmapMode) Graph.nodeColor(n => getHeatmapColor(n));
};

// Controls
document.getElementById('btn-fit').onclick = () => Graph.zoomToFit(400, 50);
document.getElementById('btn-reset').onclick = () => {
    document.getElementById('filter-status').value = '';
    document.getElementById('filter-type').value = '';
    document.getElementById('search-input').value = '';
    document.getElementById('view-mode').value = 'force';
    document.getElementById('size-by').value = 'pagerank';
    statusFilter = ''; typeFilter = ''; sizeMetric = 'pagerank'; heatmapMode = false;
    highlightedNodes = new Set();
    Graph.dagMode(null); Graph.nodeVisibility(() => true); Graph.nodeVal(n => getNodeSize(n));
    Graph.nodeColor(n => STATUS_COLORS[n.status] || '#555577');
    Graph.linkColor(l => l.critical ? '#ec489980' : '#44475a40');
    clearSelection(); hideHoverPanel(); Graph.zoomToFit(400, 50); updateVisibleCount();
    document.getElementById('heatmap-legend').classList.remove('visible');
    document.getElementById('top-nodes-panel').classList.remove('visible');
    document.getElementById('triage-panel').style.display = 'none';
    document.getElementById('btn-heatmap').classList.remove('active');
    document.getElementById('btn-triage').classList.remove('active');
    document.getElementById('btn-top').classList.remove('active');
};

// Heatmap toggle
document.getElementById('btn-heatmap').onclick = () => {
    heatmapMode = !heatmapMode;
    document.getElementById('btn-heatmap').classList.toggle('active', heatmapMode);
    document.getElementById('heatmap-legend').classList.toggle('visible', heatmapMode);
    Graph.nodeColor(n => heatmapMode ? getHeatmapColor(n) : STATUS_COLORS[n.status] || '#555577');
};

// Triage panel
document.getElementById('btn-triage').onclick = () => {
    const panel = document.getElementById('triage-panel');
    const btn = document.getElementById('btn-triage');
    const visible = panel.style.display === 'none';
    panel.style.display = visible ? 'block' : 'none';
    btn.classList.toggle('active', visible);
    if (visible && DATA.triage && DATA.triage.recommendations) {
        const list = document.getElementById('triage-list');
        list.innerHTML = DATA.triage.recommendations.slice(0, 5).map(r => {
            const score = (r.score != null && isFinite(r.score)) ? r.score.toFixed(2) : '-';
            const reason = (r.reasons && r.reasons.length > 0) ? r.reasons[0] : '';
            return '<div class="triage-item" data-id="' + (r.id || '') + '">' +
                '<div class="triage-item-header"><span class="triage-item-id">' + (r.id || '-') + '</span><span class="triage-item-score">' + score + '</span></div>' +
                '<div class="triage-item-title">' + (r.title || '') + '</div>' +
                '<div class="triage-item-reason">' + reason + '</div></div>';
        }).join('');
        list.querySelectorAll('.triage-item').forEach(item => {
            item.onclick = () => {
                const graphNodes = Graph.graphData().nodes;
                const node = graphNodes.find(n => n.id === item.dataset.id);
                if (node) { selectNode(node); Graph.centerAt(node.x, node.y, 500); Graph.zoom(2.5, 500); }
            };
        });
    }
};

// Top nodes panel
document.getElementById('btn-top').onclick = () => {
    const panel = document.getElementById('top-nodes-panel');
    const visible = panel.classList.toggle('visible');
    document.getElementById('btn-top').classList.toggle('active', visible);
    if (visible) {
        const sorted = [...DATA.nodes].sort((a, b) => (b.pagerank || 0) - (a.pagerank || 0)).slice(0, 10);
        panel.innerHTML = sorted.map((n, i) => '<div class="top-node-item" data-id="' + n.id + '"><span class="rank">#' + (i+1) + '</span><span>' + n.id + '</span></div>').join('');
        panel.querySelectorAll('.top-node-item').forEach(el => {
            el.onclick = () => {
                const graphNodes = Graph.graphData().nodes;
                const node = graphNodes.find(n => n.id === el.dataset.id);
                if (node) { selectNode(node); Graph.centerAt(node.x, node.y, 500); Graph.zoom(2.5, 500); }
            };
        });
    }
};

// Fullscreen
document.getElementById('btn-fullscreen').onclick = () => {
    if (!document.fullscreenElement) container.requestFullscreen();
    else document.exitFullscreen();
};

// Toast
function showToast(msg) {
    const toast = document.getElementById('toast');
    toast.textContent = msg; toast.classList.add('visible');
    setTimeout(() => toast.classList.remove('visible'), 2500);
}

// Keyboard shortcuts
document.onkeydown = e => {
    if (e.target.tagName === 'INPUT') return;
    switch(e.key.toLowerCase()) {
        case 'f': Graph.zoomToFit(400, 50); break;
        case 'r': document.getElementById('btn-reset').click(); break;
        case 'escape': clearSelection(); hideHoverPanel(); highlightedNodes = new Set(); Graph.nodeColor(Graph.nodeColor()); break;
        case ' ': e.preventDefault(); document.getElementById('btn-fullscreen').click(); break;
        case 'h': document.getElementById('btn-heatmap').click(); break;
        case 't': document.getElementById('btn-top').click(); break;
        case 'g': document.getElementById('btn-triage').click(); break;
        case '1': document.getElementById('view-mode').value = 'force'; Graph.dagMode(null); break;
        case '2': document.getElementById('view-mode').value = 'td'; Graph.dagMode('td'); break;
        case '3': document.getElementById('view-mode').value = 'lr'; Graph.dagMode('lr'); break;
        case '4': document.getElementById('view-mode').value = 'radialout'; Graph.dagMode('radialout'); break;
    }
};

// Initial fit
setTimeout(() => { Graph.zoomToFit(400, 50); updateVisibleCount(); }, 800);
    </script>
</body>
</html>`, title, title, nodeCount, edgeCount, nodeCount, nodeCount, edgeCount, timestamp, dataHash, projectName, forceGraphLib, markedLib, graphDataJSON)
}
