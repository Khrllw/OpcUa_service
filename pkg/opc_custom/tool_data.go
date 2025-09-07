package opc_custom

import "time"

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
