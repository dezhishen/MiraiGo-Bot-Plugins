package facesave

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/cache"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/command"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/storage"
	"github.com/go-basic/uuid"
)

// Plugin Random插件
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

var pluginID = "face-save"

// PluginInfo PluginInfo
func (p Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   pluginID,
		Name: "表情包收集",
	}
}

// IsFireEvent 是否触发
func (p Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		return true
	} else if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Image {
		var key string
		if msg.MessageType == plugins.GroupMessage {
			key = fmt.Sprintf("%v:%v", msg.MessageType, msg.GroupCode)
		} else {
			key = fmt.Sprintf("%v:%v", msg.MessageType, msg.Sender.Uin)
		}
		_, exists := cache.Get(key)
		return exists
	}
	return false
}

type FaceSaveReq struct {
	Name string `short:"n" long:"name" description:"表情的名称" required:"true"`
}

// OnMessageEvent OnMessageEvent
func (p Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{}
	var key string
	if request.MessageType == plugins.GroupMessage {
		key = fmt.Sprintf("%v:%v", request.MessageType, request.GroupCode)
	} else {
		key = fmt.Sprintf("%v:%v", request.MessageType, request.Sender.Uin)
	}
	if len(request.Elements) == 1 && request.Elements[0].Type() == message.Text {
		//标记
		v := request.Elements[0]
		field, _ := v.(*message.TextElement)
		context := field.Content
		if strings.HasPrefix(context, ".face-save") {
			req := FaceSaveReq{}
			_, err := command.Parse(".face-save", &req, strings.Split(context, " "))
			if err != nil {
				return nil, err
			}
			cache.Set(key, req.Name, 1*time.Minute)
			result.Elements = append(result.Elements, message.NewText(fmt.Sprintf("表情名称为:%v,请于一分钟之内发送一张图片", req.Name)))
		} else {
			faceKey := strings.TrimSpace(context)
			image, err := getImage(faceKey)
			if err != nil || image == nil {
				return nil, nil
			}
			if plugins.GroupMessage == request.MessageType {
				imageElement, err := request.QQClient.UploadGroupImage(request.GroupCode, bytes.NewReader(*image))
				if err != nil {
					return nil, err
				}
				result.Elements = append(result.Elements, imageElement)
			} else {
				imageElement, err := request.QQClient.UploadPrivateImage(request.Sender.Uin, bytes.NewReader(*image))
				if err != nil {
					return nil, err
				}
				result.Elements = append(result.Elements, imageElement)
			}
		}
	} else if len(request.Elements) == 1 && request.Elements[0].Type() == message.Image {
		cacheValue, exists := cache.Get(key)
		if !exists {
			return nil, errors.New("已经超过一分钟啦,请重新开始保持吧")
		}
		fileName := cacheValue.(string)
		if fileName == "" {
			return nil, errors.New("已经超过一分钟啦,请重新开始保持吧")
		}
		cache.Delete(key)
		v := request.Elements[0]
		field, _ := v.(*message.ImageElement)
		println("url:  " + field.Url)
		// reqest, _ := http.NewRequest("GET", field.Url, nil)
		// // Accept: image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8
		// // Accept-Encoding: gzip, deflate, br
		// // Accept-Language: zh-CN
		// // Cache-Control: no-cache
		// // Connection: keep-alive
		// // Host: gchat.qpic.cn
		// // Pragma: no-cache
		// // Sec-Fetch-Dest: image
		// // Sec-Fetch-Mode: no-cors
		// // Sec-Fetch-Site: cross-site
		// // User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) electron-qq/1.4.7 Chrome/89.0.4389.128 Electron/12.0.7 Safari/537.36
		// reqest.Header.Add("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
		// reqest.Header.Add("Accept-Encoding", "gzip, deflate, br")
		// reqest.Header.Add("Cache-Control", "no-cache")
		// reqest.Header.Add("Connection", "keep-alive")
		// reqest.Header.Add("Host", "gchat.qpic.cn")
		// reqest.Header.Add("Pragma", "no-cache")
		// reqest.Header.Add("Sec-Fetch-Dest", "image")
		// reqest.Header.Add("Sec-Fetch-Mode", "no-cors")
		// reqest.Header.Add("Sec-Fetch-Site", "cross-site")
		// reqest.Host = "gchat.qpic.cn"
		// r, err := http.DefaultClient.Do(reqest)
		// if err != nil {
		// 	return nil, err
		// }
		// defer r.Body.Close()
		// robots, err := ioutil.ReadAll(r.Body)
		// if err != nil {
		// 	return nil, err
		// }
		_, err := saveImage(field.Url, fileName)
		if err != nil {
			return nil, err
		}
		result.Elements = append(result.Elements, message.NewText(fmt.Sprintf("保存成功,发送[%v]试试吧", fileName)))
	}
	return result, nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
	exists, _ := pathExists("./face")
	if !exists {
		os.Mkdir("./face", 0777)
	}
}

func saveImage(url, fileName string) (string, error) {
	id, _ := uuid.GenerateUUID()
	path := fmt.Sprintf("./face/%v.jpg", id)
	// run shell `wget URL -O filepath`
	cmd := exec.Command("wget", url, "-O", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	storage.Put([]byte(pluginID), []byte(fileName), []byte(path))
	return path, nil
}

func getImage(fileName string) (*[]byte, error) {
	var filePath string
	err := storage.Get([]byte(pluginID), []byte(fileName), func(b []byte) error {
		if b != nil {
			filePath = string(b)
			return nil
		}
		return errors.New("图片不存在")
	})
	if err != nil {
		return nil, err
	}
	ok, _ := pathExists(filePath)
	if ok {
		file, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		return &content, err
	}
	return nil, errors.New("图片不存在")
}
