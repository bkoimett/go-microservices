# 🛠️ Hackathon Starter Kit - Go Microservices

A collection of production-ready, battle-tested Go microservices designed to supercharge your hackathon projects. Build faster, deploy easier, and win bigger.

## 🎯 Why This Kit?

After analyzing the 2026 hackathon landscape—from the EUDIS Defence Hackathon (defending airspace with drone detection) to AWS Breaking Barriers (GenAI for social impact)—I realized successful teams don't build from scratch. They build with reusable, reliable components.

**This is my personal arsenal of 26 Go microservices**, built at a pace of two per week, each one teaching me something new while creating instantly deployable functionality for any hackathon project.

## 📦 What's Inside

### Phase 1: The Fundamentals
*Every app needs these.*

| Service | Description | Tech Learnings |
|---------|-------------|----------------|
| [user-service](./user-service) | Complete auth with JWT, registration, email verification | JWT, bcrypt, PostgreSQL, RESTful design |
| [file-uploader](./file-uploader) | Image/file upload with validation, S3-compatible storage | Multipart forms, cloud SDKs, background jobs |
| [rate-limiter](./rate-limiter) | IP/API key rate limiting middleware | Middleware patterns, Redis, token bucket algorithm |
| [feature-flags](./feature-flags) | Toggle features on/off with admin API | Admin vs public APIs, A/B testing concepts |
| [feedback-service](./feedback-service) | Comments/ratings for any entity | Generic data modeling, nested structures |
| [notification-hub](./notification-hub) | Email/SMS/Webhook notifications | Adapter pattern, third-party APIs |
| [url-shortener](./url-shortener) | Short URLs with click tracking | Base62 encoding, redirects, analytics |
| [pdf-generator](./pdf-generator) | HTML/CSS/JSON to PDF conversion | Binary data in APIs, PDF libraries |

### Phase 2: AI & Data Intelligence
*Smart features for the GenAI era.*

| Service | Description | Tech Learnings |
|---------|-------------|----------------|
| [web-scraper](./web-scraper) | Extract main content from any URL | HTTP clients, goquery, concurrency |
| [ai-summarizer](./ai-summarizer) | LLM-powered text summarization | OpenAI API, prompt engineering, streaming |
| [sentiment-analysis](./sentiment-analysis) | Text sentiment scoring | ML concepts, API caching, Python interop |
| [rag-chatbot](./rag-chatbot) | Context-aware AI chatbot backend | RAG, vector databases, embeddings |
| [data-viz](./data-viz) | Generate charts from CSV/JSON | Data processing, go-echarts |
| [trending-aggregator](./trending-aggregator) | Social media/tech trends API | OAuth, cron jobs, data aggregation |

### Phase 3: Payments & E-commerce
*Winning solutions need money flows.*

| Service | Description | Tech Learnings |
|---------|-------------|----------------|
| [stripe-payments](./stripe-payments) | Stripe Payment Intent wrapper | Idempotency, webhooks, PCI basics |
| [paypal-checkout](./paypal-checkout) | PayPal integration microservice | Payment gateway patterns |
| [inventory-service](./inventory-service) | Thread-safe product inventory | Optimistic locking, transactions |
| [discount-engine](./discount-engine) | Promo code validation engine | Rules engine design, validation |

### Phase 4: Real-Time & Communication
*Live collaboration features.*

| Service | Description | Tech Learnings |
|---------|-------------|----------------|
| [chat-server](./chat-server) | WebSocket chat with rooms | Goroutines, channels, gorilla/websocket |
| [live-polls](./live-polls) | Real-time polling with SSE | Server-Sent Events, real-time push |
| [collab-board](./collab-board) | Collaborative whiteboard backend | Operational Transforms, sync |
| [reminder-service](./reminder-service) | Scheduled email/SMS reminders | Job queues, cron, timezones |

### Phase 5: Niche & Specialized
*For defence-grade and IoT solutions.*

| Service | Description | Tech Learnings |
|---------|-------------|----------------|
| [geo-fencing](./geo-fencing) | Geofence monitoring with alerts | Spatial data, geohashing |
| [iot-simulator](./iot-simulator) | Mock sensor data streaming | Time-series, data simulation |
| [api-gateway](./api-gateway) | Lightweight reverse proxy/router | Reverse proxy, middleware chaining |
| [health-dashboard](./health-dashboard) | Service health monitoring UI | html/template, service discovery |

## 🚀 Quick Start

Each service is designed to be independently run and integrated. Here's the standard pattern:

```bash
# Clone the specific service you need
git clone https://github.com/yourusername/hackathon-starter-kit.git
cd hackathon-starter-kit/user-service

# Copy and configure environment
cp .env.example .env

# Run with Docker (recommended)
docker-compose up -d

# Or run locally
go mod download
go run cmd/main.go
```

## 🔌 Integration Patterns

These microservices are built to be "plug-and-play" in three ways:

### 1. **As a Standalone Microservice**
```bash
# Service runs independently, exposes REST API
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secure123"}'
```

### 2. **Via Docker Compose**
```yaml
# In your hackathon project's docker-compose.yml
services:
  user-service:
    image: yourusername/user-service:latest
    environment:
      - DB_CONNECTION=postgresql://postgres:password@db:5432/users
    ports:
      - "8081:8080"
  
  your-app:
    build: .
    ports:
      - "3000:3000"
    depends_on:
      - user-service
```

### 3. **As a Go Package**
Some services can be imported directly:
```go
import "github.com/yourusername/hackathon-starter-kit/rate-limiter"

// Use as middleware in your own Go service
r.Use(ratelimiter.New(ratelimiter.Config{
    RequestsPerSecond: 10,
    Burst: 20,
}))
```

## 🏗️ Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   API Gateway   │────▶│  User Service   │────▶│    PostgreSQL   │
│  (Service #25)  │     │   (Service #1)  │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
         │                       │
         ▼                       ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│ Notification    │     │  Rate Limiter   │     │      Redis      │
│     Hub         │────▶│   (Service #3)  │────▶│                 │
│  (Service #6)   │     │                 │     └─────────────────┘
└─────────────────┘     └─────────────────┘
         │                       │
         ▼                       ▼
┌─────────────────┐     ┌─────────────────┐
│    Third-Party  │     │  AI Services    │
│  APIs (Email,   │     │  (Services 9-14)│
│      SMS)       │     │                 │
└─────────────────┘     └─────────────────┘
```

## 💡 How I Build: The 2-Per-Week Workflow

Each service follows a disciplined 5-day creation cycle:

| Day | Focus | Activity |
|-----|-------|----------|
| **Monday** | Plan | Define API contract, sketch DB schema |
| **Tuesday/Wednesday** | Build | Core logic with tests, use stdlib first |
| **Thursday** | Containerize | Dockerfile, docker-compose, comprehensive README |
| **Friday** | Learn & Document | Write key learnings, example API calls |

## 🛠️ Tech Stack

- **Language:** Go 1.21+
- **Frameworks:** Standard library + carefully chosen packages (gorilla/mux, sqlx)
- **Databases:** PostgreSQL, Redis, MongoDB (where appropriate)
- **Infrastructure:** Docker, Docker Compose
- **APIs:** REST (OpenAPI 3.0 documented), WebSockets, SSE
- **Testing:** Table-driven tests, integration tests with testcontainers

## 📚 Learning Journey

Each service includes a `LEARNINGS.md` file documenting:
- Key concepts discovered while building
- Design decisions and trade-offs
- Common pitfalls and how to avoid them
- Real-world applications in hackathons

Example from [user-service](./user-service/LEARNINGS.md):
> *"JWT refresh token rotation is critical for security. I learned the hard way that storing refresh tokens in an HTTP-only cookie with proper expiry prevents XSS attacks while allowing seamless re-authentication."*

## 🏁 Getting Started with Your First Service

Ready to build your own? Here's how to use this template:

1. **Fork this repo**
2. **Pick your first service** from Phase 1 (I recommend starting with `user-service`)
3. **Follow the 5-day workflow**
4. **Deploy and test** with the provided docker-compose
5. **Share your learnings** in a PR to help others

## 🤝 Contributing

Found a bug? Have an idea for a new service? Want to share your hackathon win using these tools? Open an issue or PR!

**Planned for 2026 Q2:**
- Kubernetes Helm charts
- GitHub Actions CI/CD templates
- More AI services (image generation, voice cloning)
- Blockchain utilities for Web3 hackathons

## 📄 License

MIT - Use these services freely in your hackathon projects, commercial applications, or learning journey.

---

**Built with ❤️ and Go for the 2026 hackathon season.**

*"The best hackathon project is the one you don't have to build from scratch."*
