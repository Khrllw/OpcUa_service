package opc_custom

import (
	"github.com/awcullen/opcua/ua"
	"log"
	"reflect"
)

type ProgramPositionDataType struct {
	ProgramName    string
	BlockNumber    uint32
	BlockContent   string
	CallStackLevel uint32
}

func init() {
	// Регистрируем BinaryEncodingID
	typ := reflect.TypeOf(ProgramPositionDataType{})
	nodeID := ua.ExpandedNodeID{
		NodeID:       ua.NewNodeIDNumeric(0, 5013),
		NamespaceURI: "http://heidenhain.de/NC/",
		ServerIndex:  0,
	}

	ua.RegisterBinaryEncodingID(typ, nodeID)

	// Проверяем регистрацию
	id, ok := ua.FindBinaryEncodingIDForType(typ)
	if ok {
		log.Printf(
			"\n✅ BinaryEncodingID зарегистрирован:\n"+
				"\tType: %-30s\n"+
				"\tNodeID: %v\n"+
				"\tNamespaceURI: %-40s\n"+
				"\tServerIndex: %d\n",
			typ.Name(),
			id.NodeID,
			id.NamespaceURI,
			id.ServerIndex,
		)
	} else {
		log.Printf("❌ Не удалось найти BinaryEncodingID для типа %s\n", typ.Name())
	}
}
