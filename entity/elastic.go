package entity

type ESCellResult struct {
	Cell     int    `json:"cell"`
	CellType string `json:"cell_type"`
	Status   string `json:"status"`
	Message  []any  `json:"msg"`
}
