# novnc support audio
> My solution is to use FFmpeg + JSMpeg

### DEMO
```
docker run -d --name vnc -p 8888:8888 -e "ENABLE_VNC=true" -e "ENABLE_VNC_AUDIO=true" wuyu191287278/vnc:chrome_91.0
```

### Install ffmpeg
```
sudo apt-get -qqy install ffmpeg
```

### Start Pulse Audio
```
pulseaudio --start --exit-idle-time=-1
```

### VNC,UDP proxy

* Proxy [main.go](./main.go)
* noVNC and jsmpeg [./static](./static)

### Run proxy 
```
go run main.go --static ./static --vncAddress localhost:5900 --udpAddress :1234 
```

### Use FFMPEG to capture audio transfer to UDP protocol
```
ffmpeg -f alsa -i pulse -f mpegts -codec:a mp2 udp://localhost:1234
```

### Access to web sites
```
http://localhost:8888
```