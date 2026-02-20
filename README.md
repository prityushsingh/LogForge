# LogForge 🔥

LogForge is a scalable centralized logging backend built in Go.

## ✨ Features

* Asynchronous log ingestion using goroutines and channels
* Persistent storage with MongoDB
* Query APIs with filtering and time-range support
* Indexed queries for performance
* Designed for scalability and reliability

## 🧱 Architecture

Client → HTTP API → Channel → Worker → MongoDB

## 🚀 Getting Started

### Run locally

```bash
go run cmd/server/main.go
```

### Send a log

```bash
POST /logs
```

### Query logs

```bash
GET /logs
```

## 🛠 Tech Stack

* Go + Gin
* MongoDB
* Concurrent pipeline design

## 🎯 Purpose

Built to learn and demonstrate backend infrastructure design and distributed system concepts.
