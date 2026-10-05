// module.go 负责把 forwarder service 装配成 legacy HTTP/Connect handler。
package forwarder

import (
	"net/http"

	"connectrpc.com/connect"

	modeladapter "cursor/internal/backend/agent/model"
	"cursor/gen/aiserverv1/aiserverv1connect"
)

type Module struct {
	Service                              *Service
	LocalBidiHandler                     http.Handler
	LocalRunSSE                          http.Handler
	LocalWriteGitCommitMessageHandler    http.Handler
	AiHandler                            http.Handler
	RepositoryServiceHandler             http.Handler
	UploadServiceHandler                 http.Handler
}

// NewModule 创建 forwarder 模块，并导出本地 Bidi / RunSSE 处理器。
func NewModule(historyRoot string, channelService modeladapter.ChannelResolver) *Module {
	service := NewService(historyRoot, channelService)
	legacyBidiAppendProcedure := "/aiserver.v1.BidiService/BidiAppend"
	legacyRunSSEProcedure := "/agent.v1.AgentService/RunSSE"
	return &Module{
		Service:                           service,
		LocalBidiHandler:                  connect.NewUnaryHandler(legacyBidiAppendProcedure, service.BidiAppend),
		LocalRunSSE:                       NewLegacyRunSSEHandler(legacyRunSSEProcedure, service.RunSSE),
		LocalWriteGitCommitMessageHandler: connect.NewUnaryHandler(aiserverv1connect.AiServiceWriteGitCommitMessageProcedure, service.WriteGitCommitMessage),
		AiHandler:                         newAIHandler(service),
		RepositoryServiceHandler:          newRepositoryServiceHandler(service),
		UploadServiceHandler:              newUploadServiceHandler(service),
	}
}
