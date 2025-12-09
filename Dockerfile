# ---- Builder Image (Go + TypeScript) ----
    FROM node:20-alpine AS ts-builder
    RUN npm install -g typescript
    WORKDIR app/
    # Copy TypeScript source and compile
    COPY runtime/ ./
    RUN tsc
    
    # ---- Go Builder ----
    FROM golang:1.24-alpine AS go-builder
    WORKDIR /app
    
    # Install dependencies
    COPY go.* ./
    COPY internal/ ./internal
    COPY --from=ts-builder /app/dist/engine.js ./internal/scripting/engine.js
    RUN go build -o tbdmud ./internal/cmd/Main.go
    
    # ---- Final Minimal Image ----
    FROM alpine:latest
    WORKDIR /server
    
    # Copy built Go binary and frontend files
    COPY --from=go-builder /app/tbdmud /server/tbdmud
    
    ENV TELNET_PORT=4000
    ENV HTTP_PORT=8080
    ENV WORLD=/opt/world

    VOLUME /opt/world

    # Expose and run
    CMD ["/server/tbdmud"]