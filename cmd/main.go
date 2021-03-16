package main

import (
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/calendar"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/hitokoto"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/random"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/tips"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/weather"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/weibolisten"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/server"
)

func main() {
	server.Start()
}
