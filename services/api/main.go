package api

import (
    "bonex-middleware/config"
    "bonex-middleware/dao"
    "bonex-middleware/log"
    "bonex-middleware/services/api/router"
    "context"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/patrickmn/go-cache"
    "github.com/rs/cors"
    "github.com/urfave/negroni"
    "net/http"
    "time"
)

// API serves the end users requests.
type api struct {
    dao    dao.DAO
    config *config.Config
    server *http.Server
    cache  *cache.Cache
}

const (
    actionsAPIPrefix = ""
)

// NewAPI initializes a new instance of API with needed fields, but doesn't start listening,
// nor creates the router.
func New(d dao.DAO, cfg *config.Config) *api {
    // Create a cache with a default expiration time of 5 minutes, and which
    // purges expired items every 10 minutes
    c := cache.New(1*time.Minute, 10*time.Minute)

    api := &api{
        dao:    d,
        config: cfg,
        cache:  c,
    }

    return api
}

// Title returns the title.
func (this *api) Title() string {
    return "API"
}

// GracefulStop shuts down the server without interrupting any
// active connections.
func (this *api) GracefulStop(ctx context.Context) error {
    return this.server.Shutdown(ctx)
}

// Run starts the http server and binds the handlers.
func (this *api) Run() error {
    r := mux.NewRouter()

    wrapper := negroni.New()
    wrapper.Use(cors.New(cors.Options{
        AllowedOrigins:   this.config.Api.AllowedOrigins,
        AllowCredentials: true,
        AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
        AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "X-API-Key"},
    }))

    //public routes (no middleware)
    router.HandleActions(r, wrapper, actionsAPIPrefix, []*router.Route{
        {"/", "GET", this.index, nil},
        {"/merchants", "GET", this.listMerchants, nil},
        {"/merchant/{address}", "GET", this.getMerchant, nil},
        {"/merchant", "POST", this.createMerchant, nil},
        {"/image/{address}", "GET", this.getImage, nil},

        {"/subscribe", "POST", this.subscribe, nil},
        {"/unsubscribe", "POST", this.unsubscribe, nil},

        {"/subscriptions/{address}", "GET", this.getSubscriptions, nil},
        {"/subscribers/{address}", "GET", this.getSubscribers, nil},
    })

    this.server = &http.Server{Addr: fmt.Sprintf(":%d", this.config.Api.Port), Handler: r}

    log.Infof("Listening on port %d", this.config.Api.Port)
    err := this.server.ListenAndServe()
    if err != nil {
        return fmt.Errorf("Cannot run API service: %s", err.Error())
    }

    return nil
}
