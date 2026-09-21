package api

import (
    "net/http"
    "strings"
    "homeserver/internal/db"
    "golang.org/x/crypto/bcrypt"
)

type setupRequest struct {
    Username string
    Password string
    Pin      string
}

func SetupHandler(w http.ResponseWriter, r *http.Request, needsSetup *bool) {
    if !*needsSetup {
        http.Error(w, "Setup already complete", http.StatusForbidden)
        return
    }

    if r.Method == http.MethodGet {
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(`<!DOCTYPE html>
        <html>
        <head><title>HomeCloud Setup</title></head>
        <body>
            <h1>HomeCloud First Time Setup</h1>
            <form method="POST" action="/setup">
                <label>Username</label><br>
                <input type="text" name="username"><br><br>
                <label>Password</label><br>
                <input type="password" name="password"><br><br>
                <label>PIN</label><br>
                <input type="password" name="pin"><br><br>
                <button type="submit">Create Account</button>
            </form>
        </body>
        </html>`))
        return
    }

    if r.Method == http.MethodPost {
        r.ParseForm()
        req := setupRequest{
            Username: strings.TrimSpace(r.FormValue("username")),
            Password: strings.TrimSpace(r.FormValue("password")),
            Pin:      strings.TrimSpace(r.FormValue("pin")),
        }

        if len(req.Username) < 3 {
            http.Error(w, "Username must be at least 3 characters", http.StatusBadRequest)
            return
        }
        if len(req.Password) < 8 {
            http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
            return
        }
        if len(req.Pin) < 4 {
            http.Error(w, "PIN must be at least 4 digits", http.StatusBadRequest)
            return
        }

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
        if err != nil {
            http.Error(w, "Server error", http.StatusInternalServerError)
            return
        }

        hashedPin, err := bcrypt.GenerateFromPassword([]byte(req.Pin), 12)
        if err != nil {
            http.Error(w, "Server error", http.StatusInternalServerError)
            return
        }

        if err := db.CreateUser(req.Username, string(hashedPassword), string(hashedPin)); err != nil {
            http.Error(w, "Could not create user", http.StatusInternalServerError)
            return
        }

        *needsSetup = false
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(`<!DOCTYPE html>
        <html>
        <head><title>HomeCloud Setup</title></head>
        <body>
            <h1>Account created!</h1>
            <p>Open your HomeCloud app to log in.</p>
        </body>
        </html>`))
    }
}