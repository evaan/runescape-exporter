# runescape-exporter
A Prometheus exporter for RuneScape stats (using Wise Old Man)

To run (Docker):
```
docker run -d -p 8340:8340 -e PLAYER_NAME="your rsn" evaan/runescape-exporter:latest
```

To run (Other):
```
go build
PLAYER_NAME="your rsn" ./runescape-exporter
```

Add <your-ip>:8340 to prometheus/whatever metrics software