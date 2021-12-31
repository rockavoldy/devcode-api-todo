package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"devcode-api-todo/repo"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sys/unix"
)

var lc = net.ListenConfig{
	Control: func(network, address string, c syscall.RawConn) error {
		var opErr error
		if err := c.Control(func(fd uintptr) {
			opErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
		}); err != nil {
			return err
		}
		return opErr
	},
}

func Logger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log.Println("from: ", os.Getpid())
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func main() {
	mysql_host := os.Getenv("MYSQL_HOST")
	if mysql_host == "" {
		mysql_host = "localhost"
	}
	mysql_user := os.Getenv("MYSQL_USER")
	if mysql_user == "" {
		mysql_user = "root"
	}
	mysql_password := os.Getenv("MYSQL_PASSWORD")
	if mysql_password == "" {
		mysql_password = "example"
	}
	mysql_dbname := os.Getenv("MYSQL_DBNAME")
	if mysql_dbname == "" {
		mysql_dbname = "teestdb"
	}

	db := repo.ConnectDB(mysql_host, mysql_user, mysql_password, mysql_dbname)
	InitRepo := repo.NewRepo(db)

	router := chi.NewRouter()
	router.Use(Logger())

	router.Mount("/activity-groups", RouterActivity(InitRepo))
	router.Mount("/todo-items", RouterTodo(InitRepo))

	l, err := lc.Listen(context.Background(), "tcp", ":3030")
	if err != nil {
		panic(err)
	}

	server := &http.Server{Handler: router}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
		shutdownCancel()
	}()

	log.Println("Listening on port :3030")
	err = server.Serve(l)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
