package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ethanmidgley/storage-bucket/pkg/config"
	"github.com/ethanmidgley/storage-bucket/pkg/server"
	"github.com/ethanmidgley/storage-bucket/pkg/storage"
)

func main() {

	pathPrefix := flag.String("pathPrefix", ".", "a prefix to add to all the paths in the code ( do not leave a trailing slash )")
	flag.Parse()

	// Load the config from the yaml file
	log.Println("Loading bucket.yaml file")
	conf, err := config.Load(*pathPrefix)
	if err != nil {
		log.Panic(err)
	}
	config.Conf = conf

	// Make sure the bucket is available
	log.Println("Checking for active bucket file")
	if !storage.CheckBucket() {
		log.Println("No bucket found. Creating one now")
		err := storage.CreateBucket()
		if err != nil {
			log.Panic("Failed to create bucket")
		}
	} else {
		log.Println("Bucket found")
	}

	log.Println("Starting control plane server")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		cancel()
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		server.Start(ctx)
		wg.Done()
	}()

	wg.Wait()

}
