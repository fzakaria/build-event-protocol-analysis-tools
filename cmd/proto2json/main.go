package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	bes "github.com/fzakaria/build-event-protocol-analysis-tools/genproto/build_event_stream"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Check if a filename is provided as a command-line argument
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <bep_file.pb>", os.Args[0])
	}
	filePath := os.Args[1]

	// Configure protojson Marshaler options
	// UseProtoNames ensures that the JSON output uses the original snake_case field names
	// from the .proto definition, which is common for BEP JSON.
	// EmitUnpopulated ensures that fields with default values are still included in the JSON.
	// Multiline and Indent make the JSON output more human-readable.
	jsonMarshaler := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
		Multiline:       true, // Output multi-line JSON
		Indent:          "  ", // Indent with two spaces
	}

	// Open the specified BEP file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", filePath, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		// 1. Read the length of the next protobuf message (varint encoded)
		//    Bazel BEP files stream messages by prefixing each with its length.
		msgLen, err := binary.ReadUvarint(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error reading message length: %v", err)
		}

		if msgLen == 0 {
			continue
		}

		// Limit message size to avoid excessive memory allocation (e.g., 100MB)
		// Adjust this limit as needed for your typical BEP event sizes.
		const maxMessageSize = 100 * 1024 * 1024
		if msgLen > maxMessageSize {
			log.Fatalf("Message length %d exceeds maximum allowed size %d", msgLen, maxMessageSize)
		}

		// 2. Read the protobuf message itself
		msgBytes := make([]byte, msgLen)
		_, err = io.ReadFull(reader, msgBytes)
		if err != nil {
			if err == io.EOF {
				fmt.Fprintln(os.Stderr, "EOF reached unexpectedly while reading message body. The BEP stream might be truncated.")
			}
			log.Fatalf("Error reading message bytes: %v", err)
		}

		// 3. Unmarshal the bytes into a BuildEvent protobuf message
		var event bes.BuildEvent
		err = proto.Unmarshal(msgBytes, &event)
		if err != nil {
			log.Printf("Error unmarshaling BuildEvent proto (length %d): %v. Skipping event.", msgLen, err)
			// You might want to dump msgBytes here for debugging, e.g., hex.Dump(msgBytes)
			continue
		}

		// 4. Marshal the BuildEvent message to JSON
		jsonBytes, err := jsonMarshaler.Marshal(&event)
		if err != nil {
			log.Printf("Error marshaling BuildEvent to JSON: %v. Skipping event: %s", err, event.String())
			continue
		}

		fmt.Print(string(jsonBytes)) // Use Print instead of Println for multi-line JSON when not part of an array
		fmt.Println()                // Add a newline to separate each JSON object if printing line-delimited JSON
	}

}
