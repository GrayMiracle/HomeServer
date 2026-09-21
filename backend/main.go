// Main program
package main

// Imports
import (
	"fmt" // Printing
	"log" // Logging
	"net/http" // HTTP Server
	"homeserver/internal/api" // Custom API import
	"homeserver/internal/config" // Config package
	"homeserver/internal/db" // DB package
	"homeserver/internal/auth" // Auth package
	"homeserver/internal/noip" // No-IP package for dynamic DNS
	"homeserver/internal/upnp" // UPnP package for port forwarding
	"bufio" // Buffered I/O for user input
	"os" // OS for user input
	"os/signal" // Signal handling for graceful shutdown
	"syscall" // System calls for signal handling
	"strings" // String manipulation
	"golang.org/x/crypto/bcrypt" // bcrypt for password hashing
	"golang.org/x/term" // Terminal input for password without echo
)

// Config global var
var cfg *config.Config

// Main server
func main() {

	// Load config file
	var err error
    cfg, err = config.Load("config.json")
	if err != nil {
		log.Fatal("Could not load config file", err) // If no config file found
	}

	// Open UPnP port
	if err := upnp.Open(); err != nil {
		log.Printf("UPnP unavailable: %v", err)
	}

	// Start noip updater
	noip.Start(cfg.Domain, cfg.NoIPUsername, cfg.NoIPPassword)

	// Initialize database connection
	if err := db.Init("data/homecloud.db"); err != nil {
        log.Fatal("Could not initialize database:", err)
    }

	// Setup flag
	needsSetup := false

	exists, err := db.UserExists()
    if err != nil { log.Fatal("Could not check users:", err) }
    if !exists {
        // Check if in terminal
		if term.IsTerminal(int(os.Stdin.Fd())) {
			firstRunSetup()
		} else {
			needsSetup = true
			fmt.Println("No account found. Open http://localhost:3164/setup to create one.")
		}
    }

	// Load JWT secret
	jwtSecret, err := config.HandleSecret()
	if err != nil {
		log.Fatal("Could not load JWT secret:", err)
	}

	mux := http.NewServeMux() // Multiplexer router to serve different routes




	// ROUTES

	// Health check for sanity check
	mux.HandleFunc("/health", corsMiddleware(healthCheck)) // Runs health check fn here

	//Activity Logs
	mux.HandleFunc("/activity", corsMiddleware(auth.JWTMiddleware(jwtSecret, api.LogHandler)))


	// Authentication routes

	// setup route
	mux.HandleFunc("/setup", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		api.SetupHandler(w, r, &needsSetup)
	}))

	// login route
	mux.HandleFunc("/auth/login", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		api.LoginHandler(w, r, jwtSecret)
	}))

	// refresh route
	mux.HandleFunc("/auth/refresh", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		api.RefreshHandler(w, r, jwtSecret)
	}))

	// logout
	mux.HandleFunc("/auth/logout", corsMiddleware(auth.JWTMiddleware(jwtSecret, api.LogoutHandler)))

	// pin verification
	mux.HandleFunc("/auth/verify-pin", corsMiddleware(auth.JWTMiddleware(jwtSecret, api.VerifyPinHandler)))

	// me auth
	mux.HandleFunc("/auth/me", corsMiddleware(auth.JWTMiddleware(jwtSecret, api.MeHandler)))


	// FILE HANDLING ROUTES

	// /files route runs ListFilesHandler from api
	mux.HandleFunc("/files", corsMiddleware(auth.JWTMiddleware(jwtSecret, func(w http.ResponseWriter, r *http.Request) {
		api.ListFilesHandler(w, r, cfg.DefaultFolder)
	})))
	// /file/download route runs DownloadFileHandler from api
	mux.HandleFunc("/file/download", corsMiddleware(auth.JWTMiddleware(jwtSecret, func(w http.ResponseWriter, r *http.Request) { 
		api.DownloadFileHandler(w, r, cfg.DefaultFolder)
	})))
	// /file/delete route runs DeleteFileHandler from api
	mux.HandleFunc("/file/delete", corsMiddleware(auth.JWTMiddleware(jwtSecret, func(w http.ResponseWriter, r *http.Request) {
		api.DeleteFileHandler(w, r, cfg.DefaultFolder)
	})))
	// /file/stream to stream vid/audio file(s)
	mux.HandleFunc("/file/stream", corsMiddleware(auth.JWTMiddleware(jwtSecret, func(w http.ResponseWriter, r *http.Request) {
		api.StreamFileHandler(w, r, cfg.DefaultFolder)
	})))
	// /files/upload to upload file(s) to server 
	mux.HandleFunc("/file/upload", corsMiddleware(auth.JWTMiddleware(jwtSecret, func(w http.ResponseWriter, r *http.Request) {
		api.UploadFileHandler(w, r, cfg.DefaultFolder)
	}))) 

	// Starts indep goroutine to listen for shutdown signal syscall.SIGTERM or syscall.SIGINT
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		<-quit
		log.Println("Shutting down, closing UPnP ports...")
		upnp.Close()
		os.Exit(0)
	}()

	fmt.Printf("Home Server currently on http://localhost:%s\n", cfg.Port) // Terminal print
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux)) // Clean error handling through logs
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // http.StatusOK = 200
	w.Write([]byte("Home Server IS Running")) // Response body
}

// CORS
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS header
		origin := r.Header.Get("Origin")
		if origin == "https://"+cfg.Domain || origin == cfg.FrontendURL {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }

		// HTTP Methods that are allowed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		// Headers allowed in requests
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// Allow (cookie) credentials
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Preflight request check before actual request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK) // 200 Ok
			return
		}

		next(w, r) // Next handler
	}
}

// First run setup to create user at the start
func firstRunSetup() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("No account found. First time setup:")

    // Keep asking until valid username
    var username string
    for {
        fmt.Print("Enter username (min 3 characters): ")
        username, _ = reader.ReadString('\n')
        username = strings.TrimSpace(username)
        if len(username) >= 3 { break }
        fmt.Println("Username must be at least 3 characters, try again.")
    }

    // Keep asking until valid password
    var password string
    for {
        fmt.Print("Enter password (min 8 characters): ")
        passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
        fmt.Println()
        if err != nil { log.Fatal("Could not read password:", err) }
        password = strings.TrimSpace(string(passwordBytes))
        if len(password) >= 8 { break }
        fmt.Println("Password must be at least 8 characters, try again.")
    }

    // Keep asking until valid PIN
    var pin string
    for {
        fmt.Print("Enter PIN (min 4 digits): ")
        pinBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
        fmt.Println()
        if err != nil { log.Fatal("Could not read PIN:", err) }
        pin = strings.TrimSpace(string(pinBytes))
        if len(pin) >= 4 { break }
        fmt.Println("PIN must be at least 4 digits, try again.")
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil { log.Fatal("Could not hash password:", err) }

    hashedPin, err := bcrypt.GenerateFromPassword([]byte(pin), 12)
    if err != nil { log.Fatal("Could not hash PIN:", err) }

    if err := db.CreateUser(username, string(hashedPassword), string(hashedPin)); err != nil {
        log.Fatal("Could not create user:", err)
    }

    fmt.Println("Account created. Starting HomeCloud...")
}