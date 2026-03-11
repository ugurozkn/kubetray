package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/ugurozkn/kubetray/pkg/config"
	"github.com/ugurozkn/kubetray/pkg/k8s"
	"github.com/ugurozkn/kubetray/pkg/platform"
)

// JSON-RPC types

type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any         `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCP types

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type initResult struct {
	ProtocolVersion string     `json:"protocolVersion"`
	ServerInfo      serverInfo `json:"serverInfo"`
	Capabilities    any        `json:"capabilities"`
}

type tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"inputSchema"`
}

type toolsResult struct {
	Tools []tool `json:"tools"`
}

type toolCallParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}

type contentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type callResult struct {
	Content []contentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// Server

type Server struct {
	version string
}

func NewServer(version string) *Server {
	return &Server{version: version}
}

func (s *Server) Run() error {
	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		line = []byte(strings.TrimSpace(string(line)))
		if len(line) == 0 {
			continue
		}

		var req request
		if err := json.Unmarshal(line, &req); err != nil {
			continue
		}

		resp := s.handle(req)
		if resp == nil {
			continue // notification, no response
		}

		out, _ := json.Marshal(resp)
		fmt.Fprintf(writer, "%s\n", out)
	}
}

func (s *Server) handle(req request) *response {
	switch req.Method {
	case "initialize":
		return s.ok(req.ID, initResult{
			ProtocolVersion: "2024-11-05",
			ServerInfo:      serverInfo{Name: "kubetray", Version: s.version},
			Capabilities:    map[string]any{"tools": map[string]any{}},
		})

	case "notifications/initialized":
		return nil

	case "tools/list":
		return s.ok(req.ID, toolsResult{Tools: s.tools()})

	case "tools/call":
		var params toolCallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return s.err(req.ID, -32602, "invalid params")
		}
		return s.callTool(req.ID, params)

	case "ping":
		return s.ok(req.ID, map[string]any{})

	default:
		return s.err(req.ID, -32601, "method not found: "+req.Method)
	}
}

func (s *Server) tools() []tool {
	return []tool{
		{
			Name:        "cluster_start",
			Description: "Start a local Kubernetes cluster using k3s via k3d. Creates a new cluster or starts an existing one.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"cpus":   map[string]any{"type": "integer", "description": "CPU allocation (default: 2)"},
					"memory": map[string]any{"type": "string", "description": "Memory allocation, e.g. 2G, 4G, 8G (default: 2G)"},
				},
			},
		},
		{
			Name:        "cluster_stop",
			Description: "Stop the running Kubernetes cluster without deleting it. Data and configuration are preserved.",
			InputSchema: map[string]any{"type": "object", "properties": map[string]any{}},
		},
		{
			Name:        "cluster_clean",
			Description: "Completely delete the Kubernetes cluster and all associated data. This is irreversible.",
			InputSchema: map[string]any{"type": "object", "properties": map[string]any{}},
		},
		{
			Name:        "cluster_status",
			Description: "Check if the Kubernetes cluster is running and get basic info.",
			InputSchema: map[string]any{"type": "object", "properties": map[string]any{}},
		},
	}
}

func (s *Server) callTool(id json.RawMessage, params toolCallParams) *response {
	switch params.Name {
	case "cluster_start":
		return s.ok(id, s.doStart(params.Arguments))
	case "cluster_stop":
		return s.ok(id, s.doStop())
	case "cluster_clean":
		return s.ok(id, s.doClean())
	case "cluster_status":
		return s.ok(id, s.doStatus())
	default:
		return s.ok(id, textError("unknown tool: "+params.Name))
	}
}

func (s *Server) doStart(args map[string]any) callResult {
	plat, err := platform.Detect()
	if err != nil {
		return textError("platform detection failed: " + err.Error())
	}

	cfg, err := config.Load()
	if err != nil {
		return textError("failed to load config: " + err.Error())
	}
	_ = cfg.EnsureDirectories()

	clusterMgr := k8s.NewClusterManager(cfg, plat)
	if clusterMgr.IsRunning() {
		return textOK("Cluster '%s' is already running. Kubeconfig: %s", cfg.ClusterName, cfg.KubeconfigPath())
	}

	cpus := cfg.DefaultCPUs
	memory := cfg.DefaultMemory
	if v, ok := args["cpus"]; ok {
		if f, ok := v.(float64); ok {
			cpus = int(f)
		}
	}
	if v, ok := args["memory"]; ok {
		if str, ok := v.(string); ok && str != "" {
			memory = str
		}
	}

	if err := clusterMgr.StartCluster(cpus, memory); err != nil {
		return textError("failed to start cluster: " + err.Error())
	}

	return textOK("Cluster '%s' started (%d CPUs, %s RAM). Kubeconfig: %s", cfg.ClusterName, cpus, memory, cfg.KubeconfigPath())
}

func (s *Server) doStop() callResult {
	plat, err := platform.Detect()
	if err != nil {
		return textError("platform detection failed: " + err.Error())
	}

	cfg, err := config.Load()
	if err != nil {
		return textError("failed to load config: " + err.Error())
	}

	clusterMgr := k8s.NewClusterManager(cfg, plat)
	if !clusterMgr.IsRunning() {
		return textOK("Cluster '%s' is not running.", cfg.ClusterName)
	}

	if err := clusterMgr.StopCluster(); err != nil {
		return textError("failed to stop cluster: " + err.Error())
	}

	return textOK("Cluster '%s' stopped. Data preserved.", cfg.ClusterName)
}

func (s *Server) doClean() callResult {
	plat, err := platform.Detect()
	if err != nil {
		return textError("platform detection failed: " + err.Error())
	}

	cfg, err := config.Load()
	if err != nil {
		return textError("failed to load config: " + err.Error())
	}

	clusterMgr := k8s.NewClusterManager(cfg, plat)
	if err := clusterMgr.DeleteCluster(); err != nil {
		return textError("failed to delete cluster: " + err.Error())
	}

	_ = os.RemoveAll(cfg.DataDir)

	return textOK("Cluster '%s' deleted and %s removed.", cfg.ClusterName, cfg.DataDir)
}

func (s *Server) doStatus() callResult {
	plat, err := platform.Detect()
	if err != nil {
		return textError("platform detection failed: " + err.Error())
	}

	cfg, err := config.Load()
	if err != nil {
		return textError("failed to load config: " + err.Error())
	}

	clusterMgr := k8s.NewClusterManager(cfg, plat)
	running := clusterMgr.IsRunning()

	if !running {
		return textOK("Cluster '%s' is not running.", cfg.ClusterName)
	}

	// Get node info if running
	info := fmt.Sprintf("Cluster '%s' is running.\nKubeconfig: %s", cfg.ClusterName, cfg.KubeconfigPath())

	kubeconfigPath := cfg.KubeconfigPath()
	cmd := exec.Command("kubectl", "--kubeconfig", kubeconfigPath, "get", "nodes", "--no-headers")
	if out, err := cmd.Output(); err == nil {
		info += "\n\nNodes:\n" + strings.TrimSpace(string(out))
	}

	cmd = exec.Command("kubectl", "--kubeconfig", kubeconfigPath, "get", "pods", "-A", "--no-headers")
	if out, err := cmd.Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		info += fmt.Sprintf("\n\nPods: %d running", len(lines))
	}

	return callResult{Content: []contentItem{{Type: "text", Text: info}}}
}

// helpers

func (s *Server) ok(id json.RawMessage, result any) *response {
	return &response{JSONRPC: "2.0", ID: id, Result: result}
}

func (s *Server) err(id json.RawMessage, code int, msg string) *response {
	return &response{JSONRPC: "2.0", ID: id, Error: &rpcError{Code: code, Message: msg}}
}

func textOK(format string, args ...any) callResult {
	return callResult{Content: []contentItem{{Type: "text", Text: fmt.Sprintf(format, args...)}}}
}

func textError(msg string) callResult {
	return callResult{Content: []contentItem{{Type: "text", Text: msg}}, IsError: true}
}
