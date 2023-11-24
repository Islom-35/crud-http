package app

import (
	"fmt"
	"golang-project-template/internal/repository/psql"
	"golang-project-template/internal/service"
	"golang-project-template/internal/transport/rest"
	"golang-project-template/pkg/database"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "grpc-server",
	Run: func(cmd *cobra.Command, args []string) {
		// Application entrypoint...

		db, err := database.OpenDatabaseConnection()
		if err != nil {
			log.Fatal(err)
		}

		// init deps
		booksRepo := psql.NewBook(db)
		booksService := service.NewBooks(booksRepo)
		handler := rest.NewHandler(booksService)

		srv := &http.Server{
			Addr:    ":5005",
			Handler: handler.InitRouter(),
		}
		log.Println("Server Started at", time.Now().Format(time.RFC3339))

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
		

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
