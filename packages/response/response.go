package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Data    any    `json:"data"`
	Success int    `json:"success"`
	Message string `json:"message"`
}

func HttpResponse[T any](w http.ResponseWriter, data T, success int, message string, status int) {
	var response Response

	response.Data = data
	response.Success = success
	response.Message = message

	if success < 1 {
		fmt.Println(message)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
