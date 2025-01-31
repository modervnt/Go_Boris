FROM golang:1.21-alpine AS builder

#install ffmpeg 
RUN apk add --no-cache ffmpeg

WORKDIR /app

COPY main.go
COPY input.MP4

RUN go build -o video-extractor .

FROM alpine:latest

RUN apk add --no-cache ffmpeg

WORKDIR /app
COPY --from=builder /app/video-extractor .
COPY --from=builder /app/input.MP4 .

# Étape 7 : Définir la commande par défaut
CMD ["./video-extractor"]

