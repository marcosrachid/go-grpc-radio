package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"

	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/gordonklaus/portaudio"
	"github.com/marcosrachid/go-grpc-radio/pkg/utils"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	for {
		files, err := ioutil.ReadDir("./audios")
		utils.Chk(err)
		randomIndex := rand.Intn(len(files))
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

			utils.Chk(binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out))
			utils.Chk(stream.Write())
			select {
			case <-sig:
				return
			default:
			}
		}
	}

}
