package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"runtime"
	"runtime/pprof"

	_ "net/http/pprof"

	"os"
	"os/signal"
	"time"
	"watcharis/go-migrate-lotto-history-els/handlers"
	"watcharis/go-migrate-lotto-history-els/repository/cache"
	database "watcharis/go-migrate-lotto-history-els/repository/db"
	lottoHistoryELS "watcharis/go-migrate-lotto-history-els/repository/elasticsearch"
	"watcharis/go-migrate-lotto-history-els/router"
	"watcharis/go-migrate-lotto-history-els/services"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var cpuprofile = flag.String("cpuprofile", "cpu.prof", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "mem.prof", "write memory profile to `file`")

func main() {

	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}).WithAttrs([]slog.Attr{
		slog.String("app_version", "1.0.0"),
	}))
	slog.SetDefault(logger)

	// initGoProfiling(ctx)

	db := initDatabaseReplica(ctx)
	slog.InfoContext(ctx, "connect MYSQL database success", slog.String("status_db", "success"))

	elasticSearch := initElasticsearch()
	slog.InfoContext(ctx, "connect Elasticsearch success", slog.String("status_els", "success"))

	redis := initRedis(ctx)
	slog.InfoContext(ctx, "connect Redis success", slog.String("status_redis", "success"))

	lottoRepository := database.NewLottoRepository(db)
	lottoHistoryRepository := database.NewLottoHistoryRepository(db)
	rewardRepository := database.NewRewardRepository(db)
	timetableRepository := database.NewTimetableRepository(db)
	lottoHistoryElasticSearch := lottoHistoryELS.NewLottoHistoryElasticSearch(elasticSearch)
	cache := cache.NewCache(redis)

	lottoHistoryService := services.NewLottoHistoryServices(
		logger,
		cache,
		lottoHistoryRepository,
		lottoHistoryElasticSearch,
		timetableRepository,
		lottoRepository,
		rewardRepository)

	lottoHistoryHandlers := handlers.NewResearchElasticAndDatabase(lottoHistoryService)

	// init Echo server
	e := echo.New()
	route := router.InitRouter(e, lottoHistoryHandlers)

	s := http.Server{
		Addr:    ":80",
		Handler: route,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func initDatabaseReplica(ctx context.Context) *gorm.DB {

	dsn := fmt.Sprintf(
		"%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local",
		"admin",
		"root",
		"(127.0.0.1:3306)",
		"go-migrate-to-els-db",
	)
	slog.Info("Connecting to database replica")

	dial := mysql.Open(dsn)
	db, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		log.Panic(ctx, err.Error())
	}

	// set up config connection Pools
	sqlDB, err := db.DB()
	if err != nil {
		log.Panic(ctx, err.Error())
	}
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Duration(5) * time.Minute)

	dbClient, err := db.DB()
	if err != nil {
		log.Panic(err)
	}

	if err := dbClient.Ping(); err != nil {
		log.Panic(err)
	}

	slog.Info("Database replica is running!")
	return db
}

func initElasticsearch() *elasticsearch.Client {
	// log.Println("Connecting to elasticsearch merchant")
	slog.Info("Connecting to elasticsearch merchant")

	es7, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Username:  "lottousr",
		Password:  "P@ssw0rd",
	})

	if err != nil {
		slog.Error("Cannot init elastic merchant client err", slog.Any("err", err))
	}

	info, err := es7.Info()
	if err != nil {
		log.Panic("Cannot get elastic merchant info")
	}

	if info.StatusCode != http.StatusOK {
		log.Panic("invalid status code connect els")
	}

	slog.Info("Connect to elasticsearch merchant success", slog.Int("http_status", info.StatusCode))
	return es7
}

func initRedis(ctx context.Context) *redis.Client {
	redisAddr := net.JoinHostPort("localhost", "6379")

	client := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		DB:           1,
		MinIdleConns: 8,
		PoolSize:     16,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Panicf("Cannot ping redis err : %+v", err)
	}

	return client
}

func initGoProfiling(ctx context.Context) {
	slog.InfoContext(ctx, "start profiling process ...")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// ... rest of the program ...

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
