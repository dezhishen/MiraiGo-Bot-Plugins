package facesave

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/cache"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/command"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/storage"
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

type ImageInfo struct {
	Filename string `json:"filename"`
	Size     int32  `json:"size"`
	Width    int32  `json:"width"`
	Height   int32  `json:"height"`
	Url      string `json:"url"`
	Md5      []byte `json:"md5"`
	Data     []byte `json:"data"`
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
			imageInfo, err := getImage(faceKey)
			if err != nil || imageInfo == nil {
				return nil, nil
			}
			if plugins.GroupMessage == request.MessageType {
				// imageElement, err := request.QQClient.UploadGroupImage(request.GroupCode, bytes.NewReader(*image))
				// if err != nil {
				// 	return nil, err
				// }
				imageElement := message.NewGroupImage(
					"",
					imageInfo.Md5,
					0,
					imageInfo.Size,
					imageInfo.Width,
					imageInfo.Height,
					2000,
				)
				result.Elements = append(result.Elements, imageElement)
			} else {
				// imageElement, err := request.QQClient.UploadPrivateImage(request.Sender.Uin, bytes.NewReader(*image))
				// if err != nil {
				// 	return nil, err
				// }
				imageElement := message.NewGroupImage(
					"",
					imageInfo.Md5,
					0,
					imageInfo.Size,
					imageInfo.Width,
					imageInfo.Height,
					2000,
				)
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
		// println(fmt.Sprintf("\nFilename	%v\nSize	%v\nWidth    %v\nHeight   %v\nUrl      %v\nMd5      %v\nData     %v\n",
		// 	field.Filename,
		// 	field.Size,
		// 	field.Width,
		// 	field.Height,
		// 	field.Url,
		// 	field.Md5,
		// 	field.Data,
		// ))
		// println("url:  " + field.Url)
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
		_, err := saveImage(field, fileName)
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

func saveImage(file *message.ImageElement, fileName string) ([]byte, error) {
	// id, _ := uuid.GenerateUUID()
	// path := fmt.Sprintf("./face/%v.jpg", id)
	// // run shell `wget URL -O filepath`
	// cmd := exec.Command("wget", url, "-O", path)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Run()
	// if err != nil {
	// 	return "", err
	// }
	fileInfo := &ImageInfo{
		Filename: file.Filename,
		Size:     file.Size,   //int32  `json:"size"`
		Width:    file.Width,  //int32  `json:"width"`
		Height:   file.Height, //int32  `json:"height"`
		Url:      file.Url,    //string `json:"url"`
		Md5:      file.Md5,    //[]byte `json:"md5"`
		Data:     file.Data,   //[]byte `json:"data"`
	}
	jsonBytes, _ := json.Marshal(fileInfo)
	storage.Put([]byte(pluginID), []byte(fileName), jsonBytes)
	return jsonBytes, nil
}

func getImage(fileName string) (*ImageInfo, error) {
	var result ImageInfo
	err := storage.Get([]byte(pluginID), []byte(fileName), func(b []byte) error {
		if b != nil {
			err := json.Unmarshal(b, &result)
			return err
		}
		return errors.New("图片不存在")
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
	// ok, _ := pathExists(filePath)
	// if ok {
	// 	scaleFilePath := strings.ReplaceAll(filePath, ".jpg", fmt.Sprintf("_%v.jpg", width))
	// 	ok, _ := pathExists(scaleFilePath)
	// 	if !ok {
	// 		datatype, err := imgtype.Get(filePath)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		if datatype == `image/jpeg` || datatype == `image/png` {
	// 			file, err := os.Open(filePath)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			defer file.Close()
	// 			src, _, err := image.Decode(file)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			bound := src.Bounds()
	// 			dx := bound.Dx()
	// 			scaleFile, _ := os.Create(scaleFilePath)
	// 			defer scaleFile.Close()
	// 			if dx > width {
	// 				dy := bound.Dy()
	// 				dst := image.NewRGBA(image.Rect(0, 0, width, width*dy/dx))
	// 				err = graphics.Scale(dst, src)
	// 				if err != nil {
	// 					return nil, err
	// 				}
	// 				err = jpeg.Encode(scaleFile, dst, &jpeg.Options{Quality: 100})
	// 				if err != nil {
	// 					return nil, err
	// 				}
	// 			} else {
	// 				io.Copy(scaleFile, file)
	// 			}
	// 		} else {
	// 			file, err := os.Open(filePath)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			defer file.Close()
	// 			scaleFile, _ := os.Create(scaleFilePath)
	// 			defer scaleFile.Close()
	// 			io.Copy(scaleFile, file)
	// 		}

	// 	}
	// 	file, err := os.Open(scaleFilePath)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	defer file.Close()
	// 	content, err := ioutil.ReadAll(file)
	// 	return &content, err
	// }
	// return nil, errors.New("图片不存在")
}
