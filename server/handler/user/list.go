package user

import (
	"net/http"

	"github.com/nepp-tumsat/documents-api/server/response"
)

func HandleUserList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		response.Success(writer, "Success")
	}
}