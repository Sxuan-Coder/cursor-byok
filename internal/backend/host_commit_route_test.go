package backend

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"connectrpc.com/connect"
	"github.com/go-chi/chi/v5"

	"cursor/gen/aiserverv1/aiserverv1connect"
	"cursor/internal/backend/forwarder"
	"cursor/internal/backend/server"
	serverconfig "cursor/internal/backend/server/config"
	"cursor/internal/backend/server/upstream"
)

func newCommitRouteTestManager(t *testing.T) *serverconfig.Manager {
	t.Helper()
	dir := t.TempDir()
	store := serverconfig.NewStore(filepath.Join(dir, "config.yaml"), dir)
	manager, err := serverconfig.NewManager(context.Background(), store)
	if err != nil {
		t.Fatalf("NewManager returned error: %v", err)
	}
	return manager
}

func TestUseLocalCommitGenerationFollowsConfig(t *testing.T) {
	manager := newCommitRouteTestManager(t)
	if useLocalCommitGeneration(manager) {
		t.Fatal("default config should keep direct forwarding to the tab server")
	}
	if _, err := manager.Save(context.Background(), serverconfig.Config{CommitModelHash: "hash-commit"}); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	if !useLocalCommitGeneration(manager) {
		t.Fatal("configured CommitModelHash should switch to local generation")
	}
	if useLocalCommitGeneration(nil) {
		t.Fatal("nil manager should keep direct forwarding")
	}
}

// 镜像 host.go 中「具体 procedure 路由 + AiService 通配路由」的注册结构，
// 确保 WriteGitCommitMessage 的专用路由不会被通配路由抢先。
func TestCommitMessageRouteShadowsAiServiceWildcard(t *testing.T) {
	manager := newCommitRouteTestManager(t)
	module := &forwarder.Module{}
	handler := server.New(
		server.Use(
			server.Recover(),
			server.ServerContext(),
			server.ErrorEncoder(),
		),
		commitMessageProcedure("/aiserver.v1.AiService/WriteGitCommitMessage", "ai_write_git_commit_message", server.ConnectUnary(), upstream.Dependencies{}, module, manager),
		server.Any("/aiserver.v1.AiService/*",
			server.Name("ai_service"),
			server.HTTP(),
			server.Local(func(ctx *server.Context) error {
				_, _ = ctx.Writer.Write([]byte("wildcard"))
				return nil
			}),
		),
	)
	router, ok := handler.(chi.Router)
	if !ok {
		t.Fatalf("handler is not chi.Router: %T", handler)
	}
	routeCtx := chi.NewRouteContext()
	if !router.Match(routeCtx, http.MethodPost, "/aiserver.v1.AiService/WriteGitCommitMessage") {
		t.Fatal("commit message route should match")
	}
	if got := routeCtx.RoutePattern(); got != "/aiserver.v1.AiService/WriteGitCommitMessage" {
		t.Fatalf("matched pattern = %q, want the specific commit message route", got)
	}
	other := chi.NewRouteContext()
	router.Match(other, http.MethodPost, "/aiserver.v1.AiService/CountTokens")
	if got := other.RoutePattern(); got != "/aiserver.v1.AiService/*" {
		t.Fatalf("other AiService procedures should hit the wildcard, got %q", got)
	}
}

// 配置了 CommitModelHash 时请求必须进入本地 handler 且不触网：
// 空 body 解出空的 WriteGitCommitMessageRequest，被 "diffs are required" 拒绝。
func TestCommitMessageLocalDispatchRejectsEmptyDiffs(t *testing.T) {
	manager := newCommitRouteTestManager(t)
	if _, err := manager.Save(context.Background(), serverconfig.Config{CommitModelHash: "hash-commit"}); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	service := &forwarder.Service{}
	module := &forwarder.Module{
		LocalWriteGitCommitMessageHandler: connect.NewUnaryHandler(
			aiserverv1connect.AiServiceWriteGitCommitMessageProcedure,
			service.WriteGitCommitMessage,
		),
	}
	handler := server.New(
		server.Use(
			server.Recover(),
			server.ServerContext(),
			server.ErrorEncoder(),
		),
		commitMessageProcedure("/aiserver.v1.AiService/WriteGitCommitMessage", "ai_write_git_commit_message", server.ConnectUnary(), upstream.Dependencies{}, module, manager),
	)
	request := httptest.NewRequest(http.MethodPost, "/aiserver.v1.AiService/WriteGitCommitMessage", strings.NewReader(""))
	request.Header.Set("Content-Type", "application/proto")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d (local handler rejects empty diffs)", recorder.Code, http.StatusBadRequest)
	}
	if body := recorder.Body.String(); !strings.Contains(body, "diffs are required") {
		t.Fatalf("response body %q should contain the local handler error", body)
	}
}