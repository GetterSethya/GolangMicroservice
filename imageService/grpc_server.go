package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GetterSethya/imageProto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	GrpcListenAddr string
	Cfg            AppConfig
	Server         *grpc.Server
	NetListener    net.Listener
	imageProto.UnimplementedUserServer
}

func NewGrpcServer(cfg AppConfig, grpclistenAddr string) *GrpcServer {

	listen, err := net.Listen("tcp", grpclistenAddr)
	if err != nil {
		log.Fatalf("Failed to start grpc userService server:%v", err)
	}

	return &GrpcServer{
		GrpcListenAddr: grpclistenAddr,
		Cfg:            cfg,
		Server:         grpc.NewServer(),
		NetListener:    listen,
	}
}

func (s *GrpcServer) RunGrpc() {
	imageProto.RegisterUserServer(s.Server, s)
	log.Println("Grpc imageService is running on port:", s.GrpcListenAddr)

	if err := s.Server.Serve(s.NetListener); err != nil {
		log.Fatalf("Failed serve imageService grpc server:%v", err)
	}
}

func (s *GrpcServer) CreateImage(ctx context.Context, req *imageProto.CreateImageReq) (*imageProto.ImageResp, error) {

	log.Println("hit Create image grpc")

	imageBytes := req.GetImageFile()
	if len(imageBytes) == 0 {
		return nil, fmt.Errorf("Missing image from request")
	}

	fileExt := strings.ToLower(filepath.Ext(req.GetFileName()))
	if !(fileExt == ".jpeg" || fileExt == ".jpg") {
		return nil, fmt.Errorf("Invalid file type: %s, only jpg/jpeg supported", fileExt)
	}

	originalDir := filepath.Join("data", "original")
	thumbDir := filepath.Join("data", "thumbnail")

	err := os.MkdirAll(originalDir, os.ModePerm)
	if err != nil {
		log.Println("Error when creating original image directory:", err)
		return nil, fmt.Errorf("Something went wrong")
	}

	err = os.MkdirAll(thumbDir, os.ModePerm)
	if err != nil {
		log.Println("Error when creating thumbnail image directory:", err)
		return nil, fmt.Errorf("Something went wrong")
	}

	timestamp := time.Now().Unix()
	uniqueId := uuid.NewString()
	stamp := fmt.Sprintf("%d-%s%s", timestamp, uniqueId, fileExt)

	// save original file
	oriFilename := filepath.Join("data", "original", stamp)
	oriFileData, err := os.Create(oriFilename)
	if err != nil {
		log.Println("Error when creating original file object:", err)
		return nil, fmt.Errorf("Something went wrong")
	}

	defer oriFileData.Close()

	fileBuff := bytes.NewBuffer(imageBytes)

	// copy thumbnail
	_, err = io.Copy(oriFileData, fileBuff)
	if err != nil {
		log.Println("Error when duplicating file object:", err)
		return nil, fmt.Errorf("Something went wrong")
	}

	thumbFilename := filepath.Join("data", "thumbnail", stamp)
	thumbFileData, err := os.Create(thumbFilename)
	if err != nil {
		log.Println("Error when creating thumbnail file object")
		return nil, fmt.Errorf("Something went wrong")
	}

	defer thumbFileData.Close()

	oriFileData.Seek(0, 0)
	img, _, err := image.Decode(oriFileData)
	if err != nil {
		log.Println("Error when decoding original image file:", err)
		return nil, fmt.Errorf("Something went wrong")
	}

	// cek apakah imagenya kecil dari ukuran thumbnail, jika tidak kita resize ke ukuran thumbnail
	if img.Bounds().Dx() > 512 {
		thumb := resize(img, 512)
		err = jpeg.Encode(thumbFileData, thumb, nil)
		if err != nil {
			log.Println("Error when scaling image:", err)
			return nil, fmt.Errorf("Invalid image/image is not supported")
		}
	}

	_, err = oriFileData.Seek(0, 0)
	if err != nil {
		log.Println("Error when seeking original image file:", err)
		return nil, fmt.Errorf("Something went wrong")
	}

	_, err = io.Copy(thumbFileData, oriFileData)
	if err != nil {
		log.Println("Error when duplicating thumbnail file object:", err)
		return nil, fmt.Errorf("Something went wrong")
	}

	respThumb := fmt.Sprintf("http://%s%s/v1/image/thumbnail"+stamp, s.Cfg.Host, s.Cfg.Port)
	respOri := fmt.Sprintf("http://%s%s/v1/image/original"+stamp, s.Cfg.Host, s.Cfg.Port)

	resp := &imageProto.ImageResp{
		Thumbnail: respThumb,
		Original:  respOri,
	}

	return resp, nil
}
