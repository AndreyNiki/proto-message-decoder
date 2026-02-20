package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	log "log/slog"
	"os"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

func main() {
	base64string := flag.String("base64", "", "Proto message in base64 format")
	protoFilePath := flag.String("proto-path", "", "Path to proto file containing message definition")
	messageName := flag.String("message-name", "", "Name of the message to decoded")
	help := flag.Bool("help", false, "Show help")

	flag.Parse()
	if *help || *base64string == "" || *protoFilePath == "" || *messageName == "" {
		flag.Usage()
		os.Exit(1)
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(*base64string)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}

	parser := &protoparse.Parser{}
	fds, err := parser.ParseFiles(*protoFilePath)
	if err != nil {
		log.Error("Error parse proto file", "err", err)
	}

	fd := fds[0]
	msgName := fd.GetPackage() + "." + *messageName
	msgFromProto := fd.FindMessage(msgName)
	if msgFromProto == nil {
		log.Error("Message not found", "msgName", *messageName)
	}
	msg := dynamic.NewMessage(msgFromProto)

	err = msg.Unmarshal(decodedBytes)
	if err != nil {
		log.Error("Error decode proto message", "err", err)
	}

	j, err := msg.MarshalJSONIndent()
	if err != nil {
		log.Error("Error marshal to json", "err", err)
	}

	log.Info(string(j))
}
