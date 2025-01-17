package main

import (
	"log"

	//"fmt"
	"io"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/jonas747/dca"
	"layeh.com/gopus"
)

var (
	args []string

	format = &audio.Format{
		NumChannels: 2,
		SampleRate:  48000,
	}
)

func main() {
	log.SetFlags(0)
	if len(os.Args) < 3 {
		log.Fatalln("usage:", os.Args[0], "[input]", "[output]")
	}

	InFile := os.Args[1]
	OutFile := os.Args[2]

	// Open the file
	inputReader, err := os.Open(InFile)
	if err != nil { // error check
		panic(err) // crash and burn
	}

	// Close the file on finish
	defer inputReader.Close()

	outputWriter, err := os.OpenFile(OutFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	defer outputWriter.Close()

	// Make a new decoder
	decoder := dca.NewDecoder(inputReader)

	opusDecoder, err := gopus.NewDecoder(
		format.SampleRate,  // sampling rate
		format.NumChannels, // channels
	)

	if err != nil {
		panic(err)
	}

	encoder := wav.NewEncoder(
		outputWriter,
		format.SampleRate, 16, format.NumChannels, 1,
	)

	defer encoder.Close()

	var intBuf = &audio.IntBuffer{
		Format:         format,
		SourceBitDepth: 16,
	}

	for {
		frame, err := decoder.OpusFrame()
		if err != nil {
			// Error happened before finishing
			if err != io.EOF {
				panic(err)
			}

			break
		}

		pcm, err := opusDecoder.Decode(frame, 960, false)
		if err != nil {
			panic(err)
		}

		ints := make([]int, len(pcm))
		for j, i := range pcm {
			ints[j] = int(i)
		}

		intBuf.Data = ints

		if err := encoder.Write(intBuf); err != nil {
			panic(err)
		}
	}
}
