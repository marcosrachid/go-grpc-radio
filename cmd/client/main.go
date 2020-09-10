package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/marcosrachid/go-grpc-stream/internal/pb"
	"github.com/marcosrachid/go-grpc-stream/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
)

const sampleRate = 44100
const seconds = 2

func main() {
	wd, _ := os.Getwd()
	certFile := filepath.Join(wd, "ssl", "cert.pem")
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		log.Fatalf("Error creating credentials: %s\n", err)
	}

	serverAddr := fmt.Sprintf(
		"%s:%s",
		utils.GetenvDefault("ADDR", pb.ADDR),
		utils.GetenvDefault("PORT", strconv.Itoa(pb.PORT)),
	)
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Fatalf("Fail to dial: %s\n", err)
	}

	defer conn.Close()
	client := pb.NewStreamerClient(conn)

	stream, err := client.Audio(context.Background(), &emptypb.Empty{})

	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]int16, 8192)
	var portAudioStream *portaudio.Stream

	for {
		time.Sleep(50 * time.Millisecond)
		utils.CallClear()
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot receive response: ", err)
		}
		log.Printf("Now Playing: %d - %s", res.GetSequence(), res.GetFilename())

		// fmt.Println("audio data: ", res.GetData())

		if portAudioStream == nil {
			portAudioStream, err = portaudio.OpenDefaultStream(0, int(res.GetChannels()), float64(res.GetRate()), len(out), &out)
			utils.Chk(err)
			defer portAudioStream.Close()

			utils.Chk(portAudioStream.Start())
			defer portAudioStream.Stop()
		}

		utils.Chk(binary.Read(bytes.NewBuffer(res.GetData()), binary.LittleEndian, out))
		utils.Chk(portAudioStream.Write())
	}
}
