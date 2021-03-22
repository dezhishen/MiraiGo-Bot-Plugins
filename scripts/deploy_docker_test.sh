git pull
docker build . -t 1179551960sdniu/miraigo:test
docker stop miraigo
docker rm miraigo
docker run -d -e BOT_DJT_KEY=9bc86e70f1c8f0d9 -e BOT_CHP_KEY=ce3552aa350641f0 --restart=always --name=miraigo -e TZ=Asia/Shanghai -v /docker_data/miraigo/:/data 1179551960sdniu/miraigo:test
