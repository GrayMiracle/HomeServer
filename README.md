# HomeServer

**Access your files from your drive anywhere.**

Self-hosted personal cloud that allows you to securely access and manage any file from your drive.  

## How it works
1) Set up your config file according to the docker image
2) Run the docker image
3) Create an account as prompted
4) Set up HTTPS serving frontend
5) Sign in and access files through web interface

All operations are performed directly on host machine

## Architecture
- nginx serves the React frontend, terminates HTTPS connections, and reverse proxies API requests to the Go backend.
- The Go backend handles authentication, file operations, media streaming, and activity logging.
- SQLite persists account information, device sessions, and activity history.
- Docker Compose manages the frontend and backend containers, with the backend bound to localhost rather than exposed publicly.
- File operations use a configured directory mounted from the host machine.

## Getting Started

### Prerequisites
- Docker and Docker Compose
- A free [No-IP](https://www.noip.com) account with a hostname
- Ports 80 and 443 accessible (UPnP handles this automatically if enabled on your router)


### Setup

1. Clone the repo:
```bash
git clone https://github.com/GrayMiracle/HomeServer.git
cd HomeServer
```
2. Create `backend/config.json`:
```json
{
    "port": "3164",
    "default_folder": "/app/files",
    "domain": "yourhostname.ddns.net",
    "IsDevelopment": false,
    "noip_username": "your-ddns-key",
    "noip_password": "your-ddns-password",
    "frontend_url": "https://yourhostname.ddns.net"
}
```
3. Create a `.env` file in the repo root: `HOMESERVER_FILES_PATH=/path/to/your/files`

4. Update `frontend/nginx.conf` with your hostname

5. Run `docker compose up --build`

6. On first run with no terminal (Docker), open `http://localhost:3164/setup` to create your account. On a local terminal run, the setup wizard starts automatically.

---


## Engineering decisions

**Self-hosted:** Dynamic DNS integration through No-IP and automatic UPnP port mapping for ports 80 and 443 support access outside the local network without a VPN. Port mappings are also removed during a normal application shutdown.

**Authentication:** Signed JWT access tokens with a 15-minute lifetime are stored in HttpOnly, Secure, SameSite=Strict cookies. Passwords and account PINs are hashed using bcrypt, with account and session information persisted in SQLite.

**Auto Port Management:** UPnP opens and closes router ports on server start/stop to minimize manual configuration

**Reverse Proxy:** nginx serves the React frontend and routes API requests to the Go backend, keeping the backend from being directly exposed to the public internet. It also handles HTTPS termination using Let's Encrypt TLS certificates.

**Device Sessions:** Bcrypt hash refresh tokens support 30-day sessions and individual session revocation. 

**Access Controls:** Auth middleware validates JWTs before allowing access to file operations and activity history, and failed login attempts are subject to IP rate limiting.

**File Management:** Files are accessed, uploaded, deleted, and downloaded directly from the host machine, without any additional service.

**Media streaming:** Go's http.ServeContent handles HTTP byte-range requests, enabling HTTP 206 Partial Content responses for browser-based video and audio seeking without requiring full-file downloads.

**Activity logging:** File uploads, downloads, views, streams, deletions, authentication events, and PIN verification attempts generate activity records containing device and request information.

**Containerization:** Docker Compose runs the Go backend and nginx-served React frontend separately. Multi-stage Dockerfiles compile the Go application and Vite frontend before copying build artifacts into smaller runtime images.

## Tech stack

| Component | Technology | Use |
|---|---|---|
| Backend | Go, net/http | REST endpoints and file handling
| Frontend | React, TypeScript, Vite | Browser-based file interface
| Styling | Tailwind CSS | UI styling
| Database | SQLite | Accounts, sessions, and activity history
| Authentication | JWT, bcrypt | Session and credential verification
| Web Server | nginx | Frontend hosting and API routing
| Containerization | Docker | Containerized app deployment
| Networking | No-IP, UPnP | Dynamic hostname and router configuration
| HTTPS | Let's Encrypt / Certbot | TLS certificate management
