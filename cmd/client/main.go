package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/marcosrachid/go-grpc-radio/pkg/utils"
)

const sampleRate = 44100
const seconds = 2

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()
	buffer := make([]float32, sampleRate*seconds)

	stream, err := portaudio.OpenDefaultStream(0, 1, sampleRate, len(buffer), func(out []float32) {
		resp, err := http.Get("http://localhost:8080/audio")
		utils.Chk(err)
		body, _ := ioutil.ReadAll(resp.Body)
		responseReader := bytes.NewReader(body)
		binary.Read(responseReader, binary.BigEndian, &buffer)
		for i := range out {
			out[i] = buffer[i]
		}
	})
	utils.Chk(err)
	utils.Chk(stream.Start())
	time.Sleep(time.Second * 40)
	utils.Chk(stream.Stop())
	defer stream.Close()

	if err != nil {
		fmt.Println(err)
	}

}
