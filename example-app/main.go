package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/cloudresourcemanager/v1"
)

type server struct {
	crm       *cloudresourcemanager.Service
	projectId string
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	project, err := s.crm.Projects.Get(s.projectId).Do()
	if err != nil {
		log.Printf("failed to get project: %v", err)
		return
	}
	fmt.Fprintf(w, "Hello from project %s (projects/%d).\n", project.ProjectId, project.ProjectNumber)
}

func main() {
	ctx := context.Background()
	crm, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		log.Printf("failed to create crm client: %s", err)
		return
	}
	s := &server{
		crm:       crm,
		projectId: os.Getenv("PROJECT_ID"),
	}
	http.Handle("/", s)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
