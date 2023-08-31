package app

import (
	"AvitoInternship/config"
	"AvitoInternship/internal/controllers"
	"AvitoInternship/internal/controllers/report_repository"
	"AvitoInternship/internal/controllers/repository"
	"AvitoInternship/internal/controllers/router"
	"AvitoInternship/internal/usecase"
)

type Test struct {
	id      int    `json:"id"`
	message string `json:"message"`
}

func Run(cfg *config.Config) {

	logger := controllers.NewLogger()
	pg := repository.NewPostgresDB(cfg, logger)
	gateway := router.New(logger, cfg.Port)
	objDB := report_repository.NewReportMaker(logger, "reports")

	userCreator := usecase.NewUserCreator(pg, logger)
	userSegChecker := usecase.NewUserSegmentsChecker(pg, logger)
	segmentCreator := usecase.NewSegmentCreator(pg, logger)
	segmentDeleter := usecase.NewSegmentDeleter(pg, logger)
	userSegUpdater := usecase.NewUserSegmentsUpdater(pg, logger)
	deadlineSetter := usecase.NewSegmentDeadlineSetter(pg, logger)
	reportCreator := usecase.NewReportCreator(cfg.SiteUrl, pg, objDB, logger)
	reportSender := usecase.NewReportSender(cfg.SiteUrl, "reports", pg, logger)

	gateway.Get("/user/new", userCreator.Run)
	gateway.Get("/user/:id", userSegChecker.Run)
	gateway.Post("/user/update", userSegUpdater.Run)
	gateway.Post("/segment/new", segmentCreator.Run)
	gateway.Post("/segment/delete", segmentDeleter.Run)
	gateway.Post("/deadline", deadlineSetter.Run)
	gateway.Post("/user/history", reportCreator.Run)
	gateway.GetFile("/reports/:link", reportSender.Run)

	gateway.Run()
}
