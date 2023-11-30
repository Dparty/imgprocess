package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"

	"imgprocess/utils/https"
	"imgprocess/utils/imageProcess"
	"imgprocess/utils/methods"
)

var (
	/*addr     = flag.String("addr", ":8080", "TCP address to listen to")*/
	/*host     = flag.String("host", "test", "TCP address to listen to")*/
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

func main() {

	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}
	if err := fasthttp.ListenAndServe(addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

var addr string
var host string
var defaultConfig image.Config
var defaultType string
var defaultImageMap map[string]*defaultImg
var defaultSize int
var defaultImage image.Image

type defaultImg struct {
	defaultImage  image.Image
	defaultConfig image.Config
	defaultType   string
	defaultSize   int
}

func init() {
	//defaulte img
	flag.Parse()
	fmt.Printf("args=%s, num=%d\n", flag.Args(), flag.NArg())
	addr = ":" + flag.Args()[0]
	host = flag.Args()[1]

	files, err := ioutil.ReadDir("./default_img/")
	if err != nil {
		log.Fatal(err)
	}

	defaultImageMap = make(map[string]*defaultImg)
	for _, file := range files {
		defaultImage, defaultConfig, defaultType, defaultSize, _ = getImageFromFilePath("./default_img/" + file.Name())
		fileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		defaultImageMap[fileName] = &defaultImg{
			defaultImage:  defaultImage,
			defaultConfig: defaultConfig,
			defaultType:   defaultType,
			defaultSize:   defaultSize,
		}
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	host = flag.Args()[1]
	host = host + string(ctx.Path())
	fmt.Printf("host : %s\n\n", host)

	var scale float64
	query := string(ctx.QueryArgs().Peek("q"))
	defaultQuery := string(ctx.QueryArgs().Peek("default"))
	methodStr := string(ctx.QueryArgs().Peek("method"))
	max := string(ctx.QueryArgs().Peek("max"))
	maxNum, _ := strconv.Atoi(max)
	maxNum = maxNum * 1000

	resp := https.DoRequest(host)
	bodyBytes := resp.Body()

	var respHeader fasthttp.ResponseHeader
	var contentType string
	var img image.Image
	var im image.Config

	if resp.StatusCode() == 200 {
		respHeader = resp.Header
		contentType = string(respHeader.ContentType())
		fmt.Printf("Header %q\n", respHeader.ContentType())

		if query != "" {
			var num string = query
			scale, _ = strconv.ParseFloat(num, 64)
		} else if methodStr == "UB" && len(bodyBytes) > maxNum {
			scale = methods.UpperBound(len(bodyBytes), maxNum)
		} else {
			scale = 1
		}
		img, _, _ = image.Decode(bytes.NewReader(bodyBytes))
		im, _, _ = image.DecodeConfig(bytes.NewReader(bodyBytes))
		send_s3 := imageProcess.CompressImg(img, im, scale, contentType)
		fmt.Printf("原來大小%d 壓縮大小%d\n", len(bodyBytes), len(send_s3))
		https.Response(ctx, send_s3, contentType)
	} else if resp.StatusCode() != 200 && defaultQuery != "" {
		contentType = "image/" + defaultImageMap[defaultQuery].defaultType
		if query != "" {
			var num string = query
			scale, _ = strconv.ParseFloat(num, 64)
		} else if methodStr == "UB" && defaultImageMap[defaultQuery].defaultSize > maxNum {
			scale = methods.UpperBound(defaultImageMap[defaultQuery].defaultSize, maxNum)
		} else {
			scale = 1
		}
		img = defaultImageMap[defaultQuery].defaultImage
		im = defaultImageMap[defaultQuery].defaultConfig
		send_s3 := imageProcess.CompressImg(img, im, scale, contentType)
		https.Response(ctx, send_s3, contentType)
	} else {
		log.Println("resp.StatusCode = ", resp.StatusCode())
		log.Println("err")
		log.Println(fasthttp.StatusUnsupportedMediaType)
		ctx.Error("unsupported path", fasthttp.StatusUnsupportedMediaType)
	}
	fmt.Printf("size %d %d\n\n", im.Width, im.Height)
}

func getImageFromFilePath(filePath string) (image.Image, image.Config, string, int, error) {
	f, err := os.Open(filePath)
	imgBytes, _ := ioutil.ReadAll(f)
	if err != nil {
		return nil, image.Config{}, "", 0, err
	}
	img, imageType, err := image.Decode(bytes.NewReader(imgBytes))
	im, _, _ := image.DecodeConfig(bytes.NewReader(imgBytes))
	fmt.Printf("The file is %d bytes long\n", len(imgBytes))
	defer f.Close()
	return img, im, imageType, len(imgBytes), err
}
