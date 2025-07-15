package pathways

type createPathwayDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type updatePathwayDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type pathwayDto struct {
	ID          string
	Name        string
	Description string
}

type createPathwayResponseDto struct {
	ID     string `json:"pathway_id"`
	Status string `json:"status"`
}
