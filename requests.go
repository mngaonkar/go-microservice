package main

type SendRecommendRequest struct {
	Name string `json:"name"`
}

type GetRecommendResponse struct {
	Name string `json:"name"`
}
