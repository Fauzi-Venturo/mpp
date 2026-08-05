package dto

// DisplayQuery picks the building whose TV is asking. One instansi = one screen.
type DisplayQuery struct {
	InstansiID string `form:"instansi_id" binding:"required,uuid"`
}
