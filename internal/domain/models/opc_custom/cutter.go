package opc_custom

import (
	"github.com/awcullen/opcua/ua"
	"log"
	"reflect"
)

type CutterLocationDataType struct {
	Position                 float64
	PositionEngineeringUnits ua.EUInformation
	CoordinateName           string
}

func init() {
	// Регистрируем BinaryEncodingID
	typ := reflect.TypeOf(CutterLocationDataType{})
	nodeID := ua.ExpandedNodeID{
		NodeID:       ua.NewNodeIDNumeric(0, 5010),
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
