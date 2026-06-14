package cmd

import (
	"fmt"
	"os"
	"scheduler-app/config"
	"scheduler-app/core"
	"scheduler-app/infra/db"
	"scheduler-app/repo"
	"scheduler-app/rest"
	"scheduler-app/rest/middleware"

	TskHandler "scheduler-app/rest/handlers/tasks"
	usrHandler "scheduler-app/rest/handlers/user"

	"scheduler-app/task"
	"scheduler-app/user"
)

func Serve() {
	cnf := config.GetConfig()

	dbCon, err := db.NewConnection(cnf.DB)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer dbCon.Close()

	err = db.MigrateDB(dbCon, "./migrations")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	broker := core.NewBroker()
	go broker.Start()

	// Repos
	userRepo := repo.NewUserRepo(dbCon)
	taskRepo := repo.NewTaskRepo(dbCon)

	// Scheduler
	scheduler := core.NewScheduler(taskRepo, broker)
	scheduler.Start()
	defer scheduler.Stop()

	// Domains
	userSvc := user.NewService(userRepo)
	taskSvc := task.NewService(taskRepo)

	middlewares := middleware.NewMiddlewares(cnf)

	userHandler := usrHandler.NewHandler(cnf, userSvc)
	TaskHandler := TskHandler.NewHandler(taskSvc, broker, middlewares)

	server := rest.NewServer(
		cnf,
		userHandler,
		TaskHandler,
	)
	server.Start()

}
