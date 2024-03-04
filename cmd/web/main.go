package main

import (
	"blog_project/internal/models"
	"context"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/go-playground/form"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type application struct {
	debug          bool
	logger         *slog.Logger
	blocks         models.BlockModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":3000", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://shelestov2905:p0cbZHri5qle3S6H@cluster0.qn3unxy.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	usersCollection := client.Database("BlogProject").Collection("users")
	blocksCollection := client.Database("BlogProject").Collection("blocks")

	logger.Info("Pinged your deployment. You successfully connected to MongoDB!")

	// templateCache, err := newTemplateCache()
	// if err != nil {
	// 	logger.Error(err.Error())
	// 	os.Exit(1)
	// }

	formDecoder := form.NewDecoder()

	// sessionManager := scs.New()
	// sessionManager.Store = mysqlstore.New(db)
	// sessionManager.Lifetime = 12 * time.Hour
	// sessionManager.Cookie.Secure = true

	app := &application{
		logger: logger,
		blocks: &models.BlockModel{Collection: blocksCollection},
		users:  &models.UserModel{Collection: usersCollection},
		// templateCache:  templateCache,
		formDecoder: formDecoder,
		// sessionManager: sessionManager,
	}

	// tlsConfig := &tls.Config{
	// 	CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	// }

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
		// TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server", "addr", *addr)

	// err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
