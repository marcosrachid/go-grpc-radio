package pb

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/gordonklaus/portaudio"
	"github.com/marcosrachid/go-grpc-radio/pkg/utils"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const (
	ADDR string = "localhost"
	PORT int    = 4000
)

type StreamServer struct{}

func NewServer() *StreamerService {
	server := &StreamServer{}
	return NewStreamerService(server)
}

func (s *StreamServer) Audio(empty *emptypb.Empty, a Streamer_AudioServer) error {
	files, err := ioutil.ReadDir("./audios")
	utils.Chk(err)

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(files))
	fmt.Println("random: ", randomIndex)
	file := files[randomIndex]

	fmt.Println("Playing: ", file.Name())
	// create mpg123 decoder instance
	decoder, err := mpg123.NewDecoder("")
	utils.Chk(err)

	utils.Chk(decoder.Open("./audios/" + file.Name()))
	defer decoder.Close()

	// get audio format information
	rate, channels, _ := decoder.GetFormat()

	// make sure output format does not change
	decoder.FormatNone()
	decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]int16, 8192)
	stream, err := portaudio.OpenDefaultStream(0, channels, float64(rate), len(out), &out)
	utils.Chk(err)
	defer stream.Close()

	utils.Chk(stream.Start())
	defer stream.Stop()
	for {
		audio := make([]byte, 2*len(out))
		_, err = decoder.Read(audio)
		if err == mpg123.EOF {
			break
		}
		utils.Chk(err)

		a.Send(&Data{
			Sequence: int32(randomIndex + 1),
			Filename: file.Name(),
			Rate:     rate,
			Channels: int64(channels),
			Data:     audio,
		})
	}
	return nil
}
