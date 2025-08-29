package models

type AxisInfosResponse struct {
	Name             string  `json:"name"`
	Position         float64 `json:"position"`
	LoadPercent      float64 `json:"load_percent"`
	ServoTemperature float64 `json:"servo_temperature"`
	CoderTemperature float64 `json:"coder_temperature"`
	PowerConsumption float64 `json:"power_consumption"`
}

type AlarmsResponse struct{}
type ProgramResponse struct {
	ProgramName   string `json:"program_name"`
	ProgramNumber int    `json:"program_number"`
	GCodeLine     string `json:"g_code_line"`
}

type SpindleInfosResponse struct{}

type MachineDataResponse struct {
	MachineId string `json:"machine_id"`
	Timestamp int64  `json:"timestamp"`

	IsEnabled bool `json:"is_enabled"`

	IsEmergency     bool `json:"is_emergency"`
	EmergencyStatus bool `json:"emergency_status"`

	MachineState string `json:"machine_state"`
	ProgramMode  string `json:"program_mode"`

	HasAlarms   bool             `json:"has_alarms"`
	AlarmStatus string           `json:"alarm_status"`
	Alarms      []AlarmsResponse `json:"alarms"`

	AxisMovementStatus string              `json:"axis_movement_status"`
	AxisInfos          []AxisInfosResponse `json:"axis_infos"`

	FeedRate     float64 `json:"feed_rate"`
	FeedOverride uint32  `json:"feed_override"`

	PartsCount float64 `json:"parts_count"`

	PowerOnTime   string `json:"power_on_time"`
	OperatingTime string `json:"operating_time"`
	CycleTime     string `json:"cycle_time"`
	CuttingTime   string `json:"cutting_time"`

	CurrentProgram ProgramResponse        `json:"current_program"`
	SpindleInfos   []SpindleInfosResponse `json:"spindle_infos"`

	CountourFeedRate float64 `json:"countour_feed_rate"`
	JogOverride      float64 `json:"jog_override"`
}
