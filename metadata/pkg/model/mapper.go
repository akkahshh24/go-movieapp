package model

import (
	"github.com/akkahshh24/movieapp/gen"
)

// MetadataToProto converts a Metadata struct into a generated proto counterpart.
func (m *Metadata) ToProto() *gen.Metadata {
	return &gen.Metadata{
		Id:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Director:    m.Director,
	}
}

// ProtoToMetadata converts a generated proto counterpart into a Metadata struct.
func ProtoToMetadata(m *gen.Metadata) *Metadata {
	return &Metadata{
		ID:          m.Id,
		Title:       m.Title,
		Description: m.Description,
		Director:    m.Director,
	}
}
