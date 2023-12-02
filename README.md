# Embedded System Course IoT Server

Comes with an embedded MQTT broker.

- UI: `https://{host}:8080`
- Broker: `{host}:8083`

## Build

With Docker:

```bash
$ docker build -t {Your Repo}:latest .
```

## How to use

Run with Docker Compose:

```bash
$ docker compose up -d
```

Usage:

1. Go to the UI
2. Register a room
3. Get the room ID from `Room ID` column in the UI dashboard
4. Publish messages to the broker to update the room status

```bash
$ mqtt pub -t room_events/{Room ID} -m '{"timestamp": "2023-12-02T15:11:02+07:00","status": "OCCUPIED"}' -h localhost -p 8083
```

```json
{
  "timestamp": "2023-12-02T15:11:02+07:00",
  "status": "OCCUPIED"
}
```

Status:

- `OCCUPIED`: Room is occupied and has people in it
- `EMPTY`: Room is empty

5. View the changes from the UI dashboard
