package main

import (
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/calendar"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/haimage"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/hitokoto"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/jrrp"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/lpl"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/mc"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/pixiv"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/random"

	// _ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/tiangou"

	// _ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/segmentation"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/thecat"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/thedog"

	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/tips"

	// _ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/vader"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/caihongpi"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/dujitang"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/weather"

	// _ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/weibolisten"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/dictionary"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/sovietjokes"
	_ "github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/plugins/translate"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/server"
)

func main() {
	server.Start()
}
