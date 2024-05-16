FROM golang:1.22 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /tggooglenotify
FROM ubuntu:22.04
RUN apt update && apt install -y ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
RUN ln -sf /usr/share/zoneinfo/Europe/Moscow /etc/localtime && echo "Europe/Moscow" > /etc/timezone
WORKDIR /
COPY --from=builder /tggooglenotify /tggooglenotify
COPY credentials.json /credentials.json

# arg
ARG TELEGRAM_TOKEN=
ARG CHAT_ID=
ARG CALENDAR_ID=
ARG CRON_EXPRESSION='0 9 * * 1-5'
ARG GOOGLE_APPLICATION_CREDENTIALS=/credentials.json
ARG USER_TAGS=''

# env
ENV TELEGRAM_TOKEN=${TELEGRAM_TOKEN}
ENV CHAT_ID=${CHAT_ID}
ENV CALENDAR_ID=${CALENDAR_ID}
ENV CRON_EXPRESSION=${CRON_EXPRESSION}
ENV GOOGLE_APPLICATION_CREDENTIALS=${GOOGLE_APPLICATION_CREDENTIALS}
ENV USER_TAGS=${USER_TAGS}

CMD ["/tggooglenotify"]
