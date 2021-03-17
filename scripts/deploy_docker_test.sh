git pull
docker build . -t 1179551960sdniu/miraigo:test
docker stop miraigo
docker rm miraigo
docker run -d --restart=always --name=miraigo -e TZ=Asia/Shanghai -v /docker_data/miraigo/:/data 1179551960sdniu/miraigo:test
