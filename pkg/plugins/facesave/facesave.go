package facesave

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && (strings.HasPrefix(field.Content, ".face-save") || strings.HasSuffix(field.Content, ".jpg"))
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
	Name string `short:"n" long:"name" description:"表情包的名称,请勿携带后缀" required:"true" default:"auto"`
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
			result.Elements = append(result.Elements, message.NewText(fmt.Sprintf("表情包名称为:%v,请于一分钟之类发送一张图片", req.Name)))
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
		fileName := fmt.Sprintf("%v", cacheValue)
		print(fileName)
		if !exists || fileName == "" {
			return nil, errors.New("已经超过一分钟啦,请重新开始保持吧")
		}
		cache.Delete(key)
		v := request.Elements[0]
		field, _ := v.(*message.ImageElement)
		r, err := http.DefaultClient.Get(field.Url)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		robots, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		saveImage(fmt.Sprintf("%v", fileName), robots)
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

func saveImage(fileName string, image []byte) string {
	id, _ := uuid.GenerateUUID()
	path := fmt.Sprintf("./face/%v.jpg", id)
	ioutil.WriteFile(path, image, 0777)
	storage.Put([]byte(pluginID), []byte(fileName), []byte(path))
	return path
}

func getImage(fileName string) (*[]byte, error) {
	var filePath string
	err := storage.Get([]byte(pluginID), []byte(fileName), func(b []byte) error {
		if b != nil {
			filePath = string(b)
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
