package dto

type NotificationDTO struct {
	UserId        string `json:"receiver_id"`
	ReservationId string `json:"reservation_id"`
	Status        string `json:"status"`
}
