package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/nepp-tumsat/documents-api/infrastructure"
	"github.com/nepp-tumsat/documents-api/infrastructure/persistence"
	"github.com/nepp-tumsat/documents-api/server/json/reads"
	"github.com/nepp-tumsat/documents-api/server/json/writes"
	"github.com/nepp-tumsat/documents-api/server/response"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/xerrors"
)

func HandleAuthSignIn() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		var requestBody reads.AuthSignInRequest
		err := json.NewDecoder(request.Body).Decode(&requestBody)
		if err != nil {
			log.Printf("%+v\n", xerrors.Errorf("Error in json: %v", err))
			response.BadRequest(writer, "Can't decode of json")
			return
		}

		authRepo := persistence.NewAuthDB(infrastructure.DB)

		userAuth, err := authRepo.SelectUserAuthByEmail(requestBody.Email)
		if err != nil {
			log.Printf("%+v\n", xerrors.Errorf("Error in repository: %v", err))
			return
		}

		err = passwordVerify(requestBody.Password, userAuth.Hash)
		if err != nil {
			log.Printf("%+v\n", xerrors.Errorf("Error in request: %v", err))
			response.BadRequest(writer, "Can't verify of password")
			return
		}

		token, err := uuid.NewRandom()
		if err != nil {
			log.Printf("%+v\n", xerrors.Errorf("Error in uuid: %v", err))
			return
		}

		err = authRepo.InsertToAuthTokens(userAuth.UserID, token.String())
		if err != nil {
			log.Printf("%+v\n", xerrors.Errorf("Error in repository: %v", err))
			return
		}

		userName, err := authRepo.SelectUserNameByUserID(userAuth.UserID)
		if err != nil {
			log.Printf("%+v\n", xerrors.Errorf("Error in repository: %v", err))
			return
		}

		response.Success(writer, writes.AuthSignInResponse{UserName: userName, Token: token.String()})
	}

}

func passwordVerify(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
