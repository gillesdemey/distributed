package main

import (
	"github.com/hashicorp/serf/serf"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	Agent  *serf.Serf
	Config *serf.Config
	done chan os.Signal
}

func (s *Server) StopGraceful() {
	if err := s.Agent.Leave(); err != nil {
		log.Fatal(err)
	}

	if err := s.Agent.Shutdown(); err != nil {
		log.Fatal(err)
	}
}

func CreateServer(seed string) (*Server, error) {
	signals := make(chan os.Signal, 1)
	config := serf.DefaultConfig()

	agent, err := serf.Create(config)
	if err != nil {
		return nil, err
	}

	if _, err = agent.Join([]string{seed}, false); err != nil {
		return nil, err
	}

	return &Server{
		Agent:  agent,
		Config: config,
		done: signals,
	}, nil
}

func main() {
	server, err := CreateServer("0.0.0.0:7947")
	if err != nil {
		log.Fatalf("Failed to create Serf server: %s\n", err)
	}

	signal.Notify(server.done, os.Interrupt, syscall.SIGTERM)
	<-server.done
	server.StopGraceful()
	os.Exit(0)
}
