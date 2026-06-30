package handler

type BodyDTO[T any] struct {
	Body T
}

type BodyInputDTO[T any] struct {
	Body T
}

type ActionDTO struct {
	OK bool `json:"ok"`
}
