# 媒体中心

### 预建资源文件夹

```jsx
mkdir -p /Users/elid/ljd/westdata/Media/appdata/jackett
mkdir -p /Users/elid/ljd/westdata/Media/appdata/qbittorrent
mkdir -p /Users/elid/ljd/westdata/Media/appdata/radarr
mkdir -p /Users/elid/ljd/westdata/Media/resources/tv
mkdir -p /Users/elid/ljd/westdata/Media/resources/animation
mkdir -p /Users/elid/ljd/westdata/Media/resources/downloads
mkdir -p /Users/elid/ljd/westdata/Media/resources/movie
mkdir -p /Users/elid/ljd/westdata/Media/resources/other
mkdir -p /Users/elid/ljd/westdata/Media/resources/vr
mkdir -p /Users/elid/ljd/westdata/Media/resources/game
```

### 组建安装

安装Flaresolverr

flaresolverr 主要作用：用于一些indexer无法通过jackett访问cloudflare

```jsx
docker run -d \
    --name=flaresolverr \
    -p 8191:8191 \
    -e LOG_LEVEL=info \
    -e TEST_URL=https://www.baidu.com \
    --restart unless-stopped \
    ghcr.io/flaresolverr/flaresolverr:latest
```

安装Jackett

```jsx
docker run -d
--name=jackett
-e PUID=1000
-e PGID=1000
-e TZ=Asia/Shanghai
-e AUTO_UPDATE=true
-p 9117:9117
-v /Users/elid/ljd/westdata/Media/appdata/jackett:/config
-v /Users/elid/ljd/westdata/Media/resources/downloads:/downloads
--restart unless-stopped
lscr.io/linuxserver/jackett:latest
```

安装qbittorrent

```jsx
docker run -d \
    --name=qbittorrent \
    -e PUID=1000 \
    -e PGID=1000 \
    -e TZ=Asia/Shanghai \
    -e WEBUI_PORT=8081 \
    -p 8081:8081 \
    -p 6881:6881 \
    -p 6881:6881/udp \
    -v /Users/elid/ljd/westdata/Media/appdata/qbittorrent:/config \
    -v /Users/elid/ljd/westdata/Media/resources/downloads:/downloads \
    --restart unless-stopped \
    docker.io/linuxserver/qbittorrent:latest
```

安装radarr

```jsx
docker run -d \
    --name=radarr \
    -e PUID=1000 \
    -e PGID=1000 \
    -e TZ=Asia/Shanghai \
    -p 7878:7878 \
    -v /Users/elid/ljd/westdata/Media/appdata/radarr:/config  \
    -v /Users/elid/ljd/westdata/Media/resources/tv:/tv \
    -v /Users/elid/ljd/westdata/Media/resources/animation:/animation \
    -v /Users/elid/ljd/westdata/Media/resources/movie:/movie \
    -v /Users/elid/ljd/westdata/Media/resources/other:/other \
    -v /Users/elid/ljd/westdata/Media/resources/vr:/vr \
    -v /Users/elid/ljd/westdata/Media/resources/game:/game \
    -v /Users/elid/ljd/westdata/Media/resources/downloads:/downloads \
    --restart unless-stopped \
    docker.io/linuxserver/radarr:latest
```

安装sonarr

```jsx
docker run -d \
    --name=sonarr \
    -e PUID=1000 \
    -e PGID=1000 \
    -e TZ=Asia/Shanghai \
    -p 8989:8989 \
    -v /Users/elid/ljd/westdata/Media/appdata/sonarr:/config \
    -v /Users/elid/ljd/westdata/Media/resources/tv:/tv \
    -v /Users/elid/ljd/westdata/Media/resources/animation:/animation \
    -v /Users/elid/ljd/westdata/Media/resources/movie:/movie \
    -v /Users/elid/ljd/westdata/Media/resources/other:/other \
    -v /Users/elid/ljd/westdata/Media/resources/vr:/vr \
    -v /Users/elid/ljd/westdata/Media/resources/game:/game \
    -v /Users/elid/ljd/westdata/Media/resources/downloads:/downloads \
    --restart unless-stopped \
    docker.io/linuxserver/sonarr:latest
```

安装jellyfin

```jsx
docker run -d \
    --name=jellyfin \
    -e PUID=1000 \
    -e PGID=1000 \
    -e TZ=Asia/Shanghai \
    -e JELLYFIN_PublishedServerUrl=192.168.0.105 \
    --privileged=true \
    -p 8096:8096 \
    -p 8920:8920 \
    -p 7359:7359/udp \
    -p 1900:1900/udp \
    -v /Users/elid/ljd/westdata/Media/appdata/jellyfin:/config \
    -v /Users/elid/ljd/westdata/Media/resources/tv:/tv \
    -v /Users/elid/ljd/westdata/Media/resources/animation:/animation \
    -v /Users/elid/ljd/westdata/Media/resources/movie:/movie \
    -v /Users/elid/ljd/westdata/Media/resources/other:/other \
    -v /Users/elid/ljd/westdata/Media/resources/vr:/vr \
    -v /Users/elid/ljd/westdata/Media/resources/game:/game \
    -v /Users/elid/ljd/westdata/Media/resources/downloads:/downloads \
    --restart unless-stopped \
    jellyfin/jellyfin:latest
```

docker run -d \
--name jellyfin \
--net=host \
-v /Users/elid/ljd/westdata/Media/appdata/jellyfin/config:/config  \
-v /Users/elid/ljd/westdata/Media/appdata/jellyfin/cache:/cache \
-v /Users/elid/ljd/westdata/Media/appdata/jellyfin/media:/media \
--restart always \
--privileged=true \
jellyfin/jellyfin