git pull
docker build . -t dezhishen/miraigo:test
docker stop miraigo
docker rm miraigo
docker run --user $(id -u) -d -e BOT_BAIDU_FANYI_KEY=`echo $BOT_BAIDU_FANYI_KEY` -e BOT_BAIDU_FANYI_ID=`echo $BOT_BAIDU_FANYI_ID` -e BOT_PIXIV_TOKEN="eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJTZG5pdSIsInV1aWQiOiJjNTQ3OGY3MGUxY2E0NmVjODJiOTNiZGVlOWNiYjU2MSIsImlhdCI6MTYxNzA5Mjk0MSwiYWNjb3VudCI6IntcImVtYWlsXCI6XCIxMTc5NTUxOTYwQHFxLmNvbVwiLFwiZ2VuZGVyXCI6LTEsXCJoYXNQcm9uXCI6MCxcImlkXCI6ODE2LFwicGFzc1dvcmRcIjpcIjFkZmViY2NjNmJjNzI4ZTc1MGExZjQ1MjhlY2Q2NjcxXCIsXCJzdGF0dXNcIjowLFwidXNlck5hbWVcIjpcIlNkbml1XCJ9IiwianRpIjoiODE2In0.VgJxTsHwXLMRApBXEW0WOuicsoc3WLLT6BrcN4lmeDk" -e BOT_DJT_KEY=9bc86e70f1c8f0d9 -e BOT_CHP_KEY=ce3552aa350641f0 --restart=always --name=miraigo -e TZ=Asia/Shanghai -v /docker_data/miraigo/:/data dezhishen/miraigo:test
