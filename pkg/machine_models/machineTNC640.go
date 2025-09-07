package machine_models

import (
	"encoding/json"
	"fmt"
	"github.com/awcullen/opcua/ua"
	"log"
	"opc_ua_service/internal/domain/models"
	"opc_ua_service/pkg/opc_custom"
	"time"
)

type ToolData struct {
	Comment                     *string   `json:"Comment,omitempty"`
	DatabaseId                  *string   `json:"DatabaseId,omitempty"`
	Name                        string    `json:"Name,omitempty"`
	ToolIndex                   int32     `json:"ToolIndex,omitempty"`
	ToolNumber                  int32     `json:"ToolNumber,omitempty"`
	Type                        string    `json:"Type,omitempty"`
	AttributeForPocket          int32     `json:"AttributeForPocket,omitempty"`
	PlcStatus                   int32     `json:"PlcStatus,omitempty"`
	LastUsage                   time.Time `json:"LastUsage,omitempty"`
	LockedStatus                int32     `json:"LockedStatus,omitempty"`
	AllowedOvertime             int32     `json:"AllowedOvertime,omitempty"`
	CurrentLifetime             int32     `json:"CurrentLifetime,omitempty"`
	MaximumLifetime             int32     `json:"MaximumLifetime,omitempty"`
	MaximumLifetimeToolCall     int32     `json:"MaximumLifetimeToolCall,omitempty"`
	UsableLength                int32     `json:"UsableLength,omitempty"`
	MaximumSpeed                *float64  `json:"MaximumSpeed,omitempty"`
	LengthOffset                int32     `json:"LengthOffset,omitempty"`
	RadiusOffset                int32     `json:"RadiusOffset,omitempty"`
	CarrierKinematics           string    `json:"CarrierKinematics,omitempty"`
	ReplacementToolNumber       *int32    `json:"ReplacementToolNumber,omitempty"`
	LengthOversize              int32     `json:"LengthOversize,omitempty"`
	Length                      int32     `json:"Length,omitempty"`
	NeckRadius                  float64   `json:"NeckRadius,omitempty"`
	Radius                      int32     `json:"Radius,omitempty"`
	CuttingDirection            int32     `json:"CuttingDirection,omitempty"`
	LengthBreakageTolerance     int32     `json:"LengthBreakageTolerance,omitempty"`
	LengthTolerance             int32     `json:"LengthTolerance,omitempty"`
	RadiusBreakageTolerance     int32     `json:"RadiusBreakageTolerance,omitempty"`
	RadiusTolerance             int32     `json:"RadiusTolerance,omitempty"`
	CutterEdgeLength            int32     `json:"CutterEdgeLength,omitempty"`
	CuttingData                 string    `json:"CuttingData,omitempty"`
	ToolEdgeMaterial            string    `json:"ToolEdgeMaterial,omitempty"`
	EdgeRadiusTolerance         int32     `json:"EdgeRadiusTolerance,omitempty"`
	Liftoff                     int32     `json:"Liftoff,omitempty"`
	ActiveFeedControlStrategy   *string   `json:"ActiveFeedControlStrategy,omitempty"`
	AfcOverloadSwitchoff        *string   `json:"AfcOverloadSwitchoff,omitempty"`
	AfcOverloadWarning          *string   `json:"AfcOverloadWarning,omitempty"`
	AfcReferencePower           *string   `json:"AfcReferencePower,omitempty"`
	PointAngle                  int32     `json:"PointAngle,omitempty"`
	RadiusAtTip                 int32     `json:"RadiusAtTip,omitempty"`
	EdgeRadiusOversize          int32     `json:"EdgeRadiusOversize,omitempty"`
	CutterEdgeRadius            int32     `json:"CutterEdgeRadius,omitempty"`
	NumberOfCutterEdges         int32     `json:"NumberOfCutterEdges,omitempty"`
	RadiusOversize              int32     `json:"RadiusOversize,omitempty"`
	MaximumPlungeAngle          int32     `json:"MaximumPlungeAngle,omitempty"`
	FrontfaceCutterWidth        int32     `json:"FrontfaceCutterWidth,omitempty"`
	EdgeRadiusCompensationTable *string   `json:"EdgeRadiusCompensationTable,omitempty"`
	ActiveChatterControl        int32     `json:"ActiveChatterControl,omitempty"`
}

type MachineData struct {
	SpeedOverride        SpeedOverride  `json:"SpeedOverride,omitempty"`
	FeedOverride         FeedOverride   `json:"FeedOverride,omitempty"`
	RapidOverride        RapidOverride  `json:"RapidOverride,omitempty"`
	OperatingMode        *int32         `json:"OperatingMode,omitempty"`
	RapidTraverseActive  *bool          `json:"RapidTraverseActive,omitempty"`
	ControlUpTime        *float64       `json:"ControlUpTime,omitempty"`
	MachineUpTime        *float64       `json:"MachineUpTime,omitempty"`
	ProgramExecutionTime *float64       `json:"ProgramExecutionTime,omitempty"`
	ExecutionState       ExecutionState `json:"ExecutionState,omitempty"`
	CurrentCall          *string        `json:"CurrentCall,omitempty"`
	ExecutionStack       []string       `json:"ExecutionStack,omitempty"`
	ActiveProgramName    *string        `json:"ActiveProgramName,omitempty"`
	Name                 *string        `json:"Name,omitempty"`
	SerialNumber         *string        `json:"SerialNumber,omitempty"`
}

type ExecutionState struct {
	CurrentState   ua.LocalizedText `json:"CurrentState,omitempty"`
	LastTransition ua.LocalizedText `json:"LastTransition,omitempty"`
}

type FeedOverride struct {
	Value            uint32           `json:"Value,omitempty"`
	EURange          []float64        `json:"EURange,omitempty"`
	EngineeringUnits ua.EUInformation `json:"EngineeringUnits,omitempty"`
}

type SpeedOverride struct {
	Value            uint32           `json:"Value,omitempty"`
	EURange          []float64        `json:"EURange,omitempty"`
	EngineeringUnits ua.EUInformation `json:"EngineeringUnits,omitempty"`
}

type RapidOverride struct {
	Value            uint32           `json:"Value,omitempty"`
	EURange          []float64        `json:"EURange,omitempty"`
	EngineeringUnits ua.EUInformation `json:"EngineeringUnits,omitempty"`
}

// Модель для станка HeidenhainTNC640
type HeidenhainTNC640Data struct {
	CutterLocation *[]opc_custom.CutterLocationDataType  `json:"cutter_location"`
	CurrentTool    opc_custom.ToolData                   `json:"current_tool"`
	Machine        MachineData                           `json:"machine_data"`
	ExecutionStack *[]opc_custom.ProgramPositionDataType `json:"execution_stack"`
	Timestamp      time.Time                             `json:"timestamp"`
}

func (m *HeidenhainTNC640Data) GetExecutionStack() ([]opc_custom.ProgramPositionDataType, error) {
	return *m.ExecutionStack, nil
}
func (m *HeidenhainTNC640Data) GetMachineID() (*string, error) {
	return m.Machine.SerialNumber, nil
}

func (m *HeidenhainTNC640Data) ConvertNodeToMachineData(nodeID string, v any) error {
	//ua.ObjectIDReadRequestEncodingDefaultXML
	switch nodeID {
	case "ns=1;i=56004": // SerialNumber
		if val, ok := v.(string); ok {
			m.Machine.SerialNumber = &val
			return nil
		}
	// -------------------------- SPEED OVERRIDE --------------------------
	case "ns=1;i=100027": // SpeedOverride
		if val, ok := v.(uint32); ok {
			m.Machine.SpeedOverride.Value = val
			return nil
		}
	case "ns=1;i=100028": // SpeedOverrideEURange
		if r, ok := v.(ua.Range); ok {
			m.Machine.SpeedOverride.EURange = []float64{r.Low, r.High}
			return nil
		}
		return fmt.Errorf("unexpected type for SpeedOverrideEURange: %T", v)
	case "ns=1;i=300003": // SpeedOverrideEngineeringUnits
		if v != nil {
			if r, ok := v.(ua.EUInformation); ok {
				m.Machine.SpeedOverride.EngineeringUnits = r
				return nil
			}
			return fmt.Errorf("unexpected type for SpeedOverrideEngineeringUnits: %T", v)
		}
		return nil
	// -------------------------- CUTTER --------------------------
	case "ns=1;i=100039": // CurrentToolName
		if v != nil {
			if val, ok := v.(string); ok {
				m.CurrentTool.Name = val
				return nil
			}
			return fmt.Errorf("unexpected type for CurrentToolName: %T", v)
		}
		return nil
	case "ns=1;i=100003": // CutterLocation
		if v != nil {
			eoSlice, ok := v.([]ua.ExtensionObject)
			if !ok {
				log.Printf("CutterLocation: expected []ua.ExtensionObject, got %T", v)
				return fmt.Errorf("unexpected type for CutterLocation: %T", v)
			}

			var result []opc_custom.CutterLocationDataType
			for _, eo := range eoSlice {
				if eo == nil {
					continue
				}
				cl, ok := eo.(opc_custom.CutterLocationDataType)
				if !ok {
					log.Printf("Unexpected type inside ExtensionObject: %T", eo)
					continue
				}
				result = append(result, cl)
			}
			m.CutterLocation = &result
		}
		return nil
	// -------------------------- FEED OVERRIDE --------------------------
	case "ns=1;i=100025": // FeedOverride
		if val, ok := v.(uint32); ok {
			m.Machine.FeedOverride.Value = val
			return nil
		}
	case "ns=1;i=100026": // FeedOverrideEURange
		if r, ok := v.(ua.Range); ok {
			m.Machine.FeedOverride.EURange = []float64{r.Low, r.High}
			return nil
		}
		return fmt.Errorf("unexpected type for FeedOverrideEURange: %T", v)
	case "ns=1;i=300002": // FeedOverrideEngineeringUnits
		if r, ok := v.(ua.EUInformation); ok {
			m.Machine.FeedOverride.EngineeringUnits = r
			return nil
		}
		return fmt.Errorf("unexpected type for FeedOverrideEngineeringUnits: %T", v)

	case "ns=1;i=100024": // OperatingMode
		if val, ok := v.(int32); ok {
			m.Machine.OperatingMode = &val
			return nil
		}

	// -------------------------- RAPID --------------------------
	case "ns=1;i=100029": // RapidOverrideValue
		if val, ok := v.(uint32); ok {
			m.Machine.RapidOverride.Value = val
			return nil
		}
	case "ns=1;i=100030": // RapidOverrideEURange
		if r, ok := v.(ua.Range); ok {
			m.Machine.RapidOverride.EURange = []float64{r.Low, r.High}
			return nil
		}
		return fmt.Errorf("unexpected type for RapidOverrideEURange: %T", v)
	case "ns=1;i=300004": // RapidOverrideEngineeringUnits
		if r, ok := v.(ua.EUInformation); ok {
			m.Machine.RapidOverride.EngineeringUnits = r
			return nil
		}
		return fmt.Errorf("unexpected type for RapidOverrideEngineeringUnits: %T", v)
	case "ns=1;i=100031": // RapidTraverseActive
		if val, ok := v.(bool); ok {
			m.Machine.RapidTraverseActive = &val
			return nil
		}

	case "ns=1;i=56031": // ControlUpTime
		if val, ok := v.(float64); ok {
			m.Machine.ControlUpTime = &val
			return nil
		}
	case "ns=1;i=56033": // MachineUpTime
		if val, ok := v.(float64); ok {
			m.Machine.MachineUpTime = &val
			return nil
		}
	case "ns=1;i=56032": // ProgramExecutionTime
		if val, ok := v.(float64); ok {
			m.Machine.ProgramExecutionTime = &val
			return nil
		}

	// -------------------------- STATE --------------------------
	case "ns=1;i=51002": // StateCurrentState
		if val, ok := v.(ua.LocalizedText); ok {
			m.Machine.ExecutionState.CurrentState = val
			return nil
		}

	case "ns=1;i=100005": // CurrentCall
		if val, ok := v.(string); ok {
			m.Machine.CurrentCall = &val
			return nil
		}
	case "ns=1;i=100006": // ExecutionStack
		if v != nil {
			eoSlice, ok := v.([]ua.ExtensionObject)
			if !ok {
				log.Printf("CutterLocation: expected []ua.ExtensionObject, got %T", v)
				return fmt.Errorf("unexpected type for CutterLocation: %T", v)
			}

			var result []opc_custom.ProgramPositionDataType
			for _, eo := range eoSlice {
				if eo == nil {
					continue
				}
				cl, ok := eo.(opc_custom.ProgramPositionDataType)
				if !ok {
					log.Printf("Unexpected type inside ExtensionObject: %T", eo)
					continue
				}
				result = append(result, cl)
			}
			m.ExecutionStack = &result
		}
		return nil
	case "ns=1;i=100022": // ActiveProgramName
		if val, ok := v.(string); ok {
			m.Machine.ActiveProgramName = &val
			return nil
		}

	// -------------------------- EXECUTION STATE --------------------------
	case "ns=1;i=100010": // ExecutionStateCurrentState
		if val, ok := v.(ua.LocalizedText); ok {
			m.Machine.ExecutionState.CurrentState = val
			return nil
		}

	case "ns=1;i=100008": // ExecutionStateLastTransition
		if v != nil {
			if val, ok := v.(ua.LocalizedText); ok {
				m.Machine.ExecutionState.LastTransition = val
			}
		} else {
			return nil
		}

	default:
		return fmt.Errorf("unsupported NodeID: %s", nodeID)
	}
	return fmt.Errorf("type mismatch for NodeID: %s", nodeID)
}

// ToResponse TODO: Все что закомменчено хз то это
func (m *HeidenhainTNC640Data) ToResponse() models.MachineDataResponse {
	resp := models.MachineDataResponse{
		MachineId: getStringOrDefault(m.Machine.SerialNumber, ""),
		Timestamp: 0,
		//IsEnabled:          false,
		//IsEmergency:        false,
		MachineState: m.Machine.ExecutionState.LastTransition.String(),
		//ProgramMode:        "",
		//AxisMovementStatus: "",
		//HasAlarms:          false,
		//AlarmStatus:        "",
		//Alarms:             nil,
		FeedOverride: m.Machine.FeedOverride.Value,
		//FeedRate: 0,
		//PartsCount:    0,
		PowerOnTime:   formatTime(m.Machine.MachineUpTime),
		OperatingTime: formatTime(m.Machine.ControlUpTime),
		//CycleTime:     0,
		CuttingTime: formatTime(m.Machine.ProgramExecutionTime),
		//SpindleInfos:
		//CountourFeedRate:
		//JogOverride:
	}
	if m.ExecutionStack != nil {
		pr_data := *m.ExecutionStack
		var last uint32 = 0
		for _, prData := range pr_data {
			if prData.CallStackLevel >= last {
				resp.CurrentProgram.ProgramNumber = int(prData.BlockNumber)
				resp.CurrentProgram.ProgramName = prData.ProgramName
				resp.CurrentProgram.GCodeLine = prData.BlockContent
			}
		}
	}
	var axisInfos []models.AxisInfosResponse

	if m.CutterLocation != nil {
		axisInfos = make([]models.AxisInfosResponse, len(*m.CutterLocation))

		for i, cutInfo := range *m.CutterLocation {
			axisInfos[i] = models.AxisInfosResponse{
				Name:     cutInfo.CoordinateName,
				Position: cutInfo.Position,
			}
		}
	}
	resp.AxisInfos = axisInfos
	return resp
}

func (m *HeidenhainTNC640Data) ToJSON() string {
	bytes, err := json.MarshalIndent(m.ToResponse(), "", "  ")
	if err != nil {
		log.Printf("Failed to marshal HeidenhainTNC640Data to JSON: %v", err)
		return ""
	}
	return string(bytes)
}

func (ci *HeidenhainTNC640Data) GetRelevantNodeIDs() []ua.NodeIDNumeric {
	// Здесь перечисляем NodeID, которые нам нужны для Heidenhain TNC640
	return []ua.NodeIDNumeric{
		ua.NewNodeIDNumeric(1, 56004), // SerialNumber

		ua.NewNodeIDNumeric(1, 100024), // OperatingMode
		// ------------------------ TOOL ------------------------
		ua.NewNodeIDNumeric(1, 100039), // CurrentToolName
		ua.NewNodeIDNumeric(1, 100003), // CutterLocation
		// ------------------------ FEED ------------------------
		ua.NewNodeIDNumeric(1, 100025), // FeedOverride
		ua.NewNodeIDNumeric(1, 100026), // FeedOverrideEURange
		ua.NewNodeIDNumeric(1, 300002), // FeedOverrideEngineeringUnits
		// ------------------------ RAPID ------------------------
		ua.NewNodeIDNumeric(1, 100029), // RapidOverride
		ua.NewNodeIDNumeric(1, 100030), // RapidOverrideEURange
		ua.NewNodeIDNumeric(1, 300004), // RapidOverrideEngineeringUnits
		ua.NewNodeIDNumeric(1, 100031), // RapidTraverseActive
		// ------------------------ SPEED ------------------------
		ua.NewNodeIDNumeric(1, 100027), // SpeedOverride
		ua.NewNodeIDNumeric(1, 100028), // SpeedOverrideEURange
		ua.NewNodeIDNumeric(1, 300003), // SpeedOverrideEngineeringUnits
		// ------------------------ TIME ------------------------
		ua.NewNodeIDNumeric(1, 56031), // ControlUpTime
		ua.NewNodeIDNumeric(1, 56033), // MachineUpTime
		ua.NewNodeIDNumeric(1, 56032), // ProgramExecutionTime
		// ---------------  PROGRAM  -------------------
		ua.NewNodeIDNumeric(1, 51002),  // CurrentState
		ua.NewNodeIDNumeric(1, 100005), // CurrentCall
		ua.NewNodeIDNumeric(1, 100006), // ExecutionStack
		ua.NewNodeIDNumeric(1, 100022), // ActiveProgramName
		// ----------------- EXECUTION STATE -------------------------
		ua.NewNodeIDNumeric(1, 100010), // ExecutionStateCurrentState
		ua.NewNodeIDNumeric(1, 100008), // ExecutionStateLastTransition
	}
}

func formatTime(ms *float64) string {
	if ms != nil {
		d := time.Duration(*ms * float64(time.Millisecond))

		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		secs := int(d.Seconds()) % 60

		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
	}
	return ""
}

// ---- вспомогательные функции ----

func getStringOrDefault(ptr *string, defaultVal string) string {
	if ptr != nil {
		return *ptr
	}
	return defaultVal
}

func getStringOrDefaultPointer(ptr interface{}, defaultVal string) string {
	if ptr != nil {
		return fmt.Sprintf("%v", ptr) // безопасно преобразуем в строку
	}
	return defaultVal
}

func getFloatOrDefault(ptr *float64, defaultVal float64) float64 {
	if ptr != nil {
		return *ptr
	}
	return defaultVal
}

func formatTimeSafe(ptr *int64) int64 {
	if ptr != nil {
		return *ptr
	}
	return 0
}
