package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"os"

	"github.com/fzakaria/build-event-protocol-analysis-tools/converter"
	bes "github.com/fzakaria/build-event-protocol-analysis-tools/genproto/build_event_stream"
	"github.com/parquet-go/parquet-go"
	"google.golang.org/protobuf/proto"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <bep_file.pb>", os.Args[0])
	}
	filePath := os.Args[1]
	bepFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open BEP file: %v", err)
	}
	defer bepFile.Close()

	f, err := os.OpenFile("output.parquet", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open parquet file: %v", err)
	}
	defer f.Close()

	writer := parquet.NewWriter(f)

	reader := bufio.NewReader(bepFile)
	for {
		length, err := binary.ReadUvarint(reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading message length: %v", err)
		}
		if length == 0 {
			log.Println("Read zero length message, skipping.")
			continue
		}
		msgBytes := make([]byte, length)
		if _, err := io.ReadFull(reader, msgBytes); err != nil {
			log.Fatalf("Error reading message bytes: %v", err)
		}
		var event bes.BuildEvent
		if err := proto.Unmarshal(msgBytes, &event); err != nil {
			log.Printf("Failed to unmarshal BuildEvent: %v. Skipping.", err)
			continue
		}

		// Create the Parquet row based on the event type
		if row, err := converter.Convert(&event); err != nil {
			log.Printf("Failed to convert event: %v\n", err)
		} else {
			if err := writer.Write(row); err != nil {
				log.Fatalf("Failed to write row to Parquet file: %v", err)
			}
		}
	}

	if err := writer.Close(); err != nil {
		log.Fatalf("Parquet writer close error: %v", err)
	}

}
