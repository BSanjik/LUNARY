// Точка входа, запускает сервер API Gateway
package main

import (
	"time"

	"api-gateway/internal"
	"api-gateway/internal/config"

	"github.com/BSanjik/LUNARY/services/common/shared"
)

func main() {
	shared.InitLogger("api-gateway", "info")

	cfg, err := config.LoadConfig()
	if err != nil {
		shared.LogError(err, "Не удалось загрузить конфигурацию")
		return
	}

	healthManager := shared.NewHealthManager("api-gateway", "1.0.0")

	srv := internal.NewServer(cfg, healthManager)

	graceful := shared.NewGracefulShutdown(srv, 30*time.Second)

	shared.LogInfo("API Gateway запущен на " + cfg.ListnerAddr)
	graceful.Start()
}
