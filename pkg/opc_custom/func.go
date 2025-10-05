package opc_custom

import (
	"github.com/awcullen/opcua/ua"
	"log"
	"reflect"
)

func registerType(typ reflect.Type, nodeIDNumeric uint32, namespace string) {
	nodeID := ua.ExpandedNodeID{
		NodeID:       ua.NewNodeIDNumeric(0, nodeIDNumeric),
		NamespaceURI: namespace,
		ServerIndex:  0,
	}
	ua.RegisterBinaryEncodingID(typ, nodeID)

	id, ok := ua.FindBinaryEncodingIDForType(typ)
	if ok {
		log.Printf(
			"\nBinaryEncodingID зарегистрирован:\n"+
				"\tType: %s\n"+
				"\tNodeID: %v\n"+
				"\tNamespaceURI: %s\n"+
				"\tServerIndex: %d\n",
			typ.Name(), id.NodeID, id.NamespaceURI, id.ServerIndex,
		)
	} else {
		log.Printf("Не удалось найти BinaryEncodingID для типа %s\n", typ.Name())
	}
}
