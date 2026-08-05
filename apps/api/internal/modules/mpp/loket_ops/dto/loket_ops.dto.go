package dto

// CallNextRequest asks for the next waiting ticket of a service.
//
// LoketID is optional: an operator app sends its own loket (pull, the shape in
// docs/04-api/rest-endpoints.md:49), while leaving it empty lets the server pick
// the eligible loket that has been idle the longest (BR-12 / FR-QUE-02).
type CallNextRequest struct {
	JenisLayananID string `json:"jenis_layanan_id" binding:"required,uuid"`
	LoketID        string `json:"loket_id" binding:"omitempty,uuid"`
}
