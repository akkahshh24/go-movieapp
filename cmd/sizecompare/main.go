package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/akkahshh24/movieapp/gen"
	"github.com/akkahshh24/movieapp/metadata/pkg/model"
	"google.golang.org/protobuf/proto"
)

var metadata = &model.Metadata{
	ID:          "123",
	Title:       "The Movie 2",
	Description: "Sequel of the legendary The Movie",
	Director:    "Foo Bars",
}

var genMetadata = &gen.Metadata{
	Id:          "123",
	Title:       "The Movie 2",
	Description: "Sequel of the legendary The Movie",
	Director:    "Foo Bars",
}

func main() {
	jsonBytes, err := serializeToJSON(metadata)
	if err != nil {
		panic(err)
	}

	xmlBytes, err := serializeToXML(metadata)
	if err != nil {
		panic(err)
	}

	protoBytes, err := serializeToProto(genMetadata)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON size:\t%dB\n", len(jsonBytes))   // JSON size:      106B
	fmt.Printf("XML size:\t%dB\n", len(xmlBytes))     // XML size:       148B
	fmt.Printf("Proto size:\t%dB\n", len(protoBytes)) // Proto size:     63B
}

func serializeToJSON(m *model.Metadata) ([]byte, error) {
	return json.Marshal(m)
}

func serializeToXML(m *model.Metadata) ([]byte, error) {
	return xml.Marshal(m)
}

func serializeToProto(m *gen.Metadata) ([]byte, error) {
	return proto.Marshal(m)
}
