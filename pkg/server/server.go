package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ethanmidgley/storage-bucket/pkg/auth"
	"github.com/ethanmidgley/storage-bucket/pkg/config"
	"github.com/ethanmidgley/storage-bucket/pkg/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Start will create the gin server and will run it
func Start(ctx context.Context) {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.SetTrustedProxies(nil)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.Conf.Yaml.ControlPlane.AllowedOrigins

	r.Use(cors.New(corsConfig))

	r.POST("/generate", auth.IsHTAuthenticated(), controllers.CreateKeys)
	r.POST("/upload", auth.IsAuthenticated(), controllers.Upload)
	r.POST("/export", auth.IsAuthenticated(), controllers.Export)
	r.POST("/delete", auth.IsAuthenticated(), controllers.Delete)

	files := r.Group("/files")
	files.Use(auth.IsAuthenticated())
	files.StaticFS("", http.Dir(config.Conf.PathPrefix+"/"+config.Conf.Yaml.Bucket.Location))

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Conf.Yaml.ControlPlane.Host, config.Conf.Yaml.ControlPlane.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Control panel listening: %s\n", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down control plane")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced shutdown: ", err)
	}

	log.Println("Server gracefully shutdown")

}
