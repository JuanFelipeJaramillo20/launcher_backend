package dto

type CommandRequest struct {
	Command string `json:"command" binding:"required"`
}

type CommandResponse struct {
	Output string `json:"output"`
}
