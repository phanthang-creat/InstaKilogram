package myJwt

import (
	"os"
	"strings"
	"time"
	api "ytb/src/controller/api/user"
	"ytb/src/model"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	// "github.com/twinj/uuid"
)

type MapClaims map[string]interface{}

func CreateJWT(userID, userName, role string) (model.TokenDetails, error) {
	td := model.TokenDetails{}

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["userID"] = userID
	atClaims["userName"] = userName
	atClaims["role"] = role
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return td, err
	}

	rtClams := jwt.MapClaims{}
	rtClams["userID"] = userID
	rtClams["exp"] = time.Now().Add(time.Hour * 240).Unix()
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClams)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return td, err
	}

	return td, nil
}

func RefreshToken(c echo.Context) error {
	// td := model.TokenDetails{}\
	refresh_token := strings.Split(c.Request().Header.Get("Authorization"), " ")[1]
	tkn, err := jwt.Parse(refresh_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		// log.Println("ParseToken", err)
		return c.String(401, "refresh token is expired")
	}

	if !VerifyExpires(tkn.Claims.(jwt.MapClaims)) {
		// log.Println("RefreshToken", tkn.Claims.(jwt.MapClaims))
		return c.String(401, "refresh token is expired")
	} else {
		userID := tkn.Claims.(jwt.MapClaims)["userID"].(string)
		maps, _ := api.GetUserById(userID)
		atClaims := jwt.MapClaims{}
		atClaims["userID"] = userID
		atClaims["userName"] = maps[0]["username"]
		atClaims["role"] = maps[0]["role"]
		atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
		accessToken, err := at.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return err
		}

		type NewToken struct {
			AccessToken string `json:"access_token"`
		}

		return c.JSON(200, &NewToken{AccessToken: accessToken})
	}
}

func ParseATToken(rToken string, e echo.Context) (bool, error) {
	tkn, err := jwt.Parse(rToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	// log.Println("err", err)
	if err != nil {
		// log.Printf("ParseToken %v", err)
		return false, err
	}
	e.Set("userID", tkn.Claims.(jwt.MapClaims)["userID"])
	if (tkn.Claims.(jwt.MapClaims)["userID"] == e.Param("id")) || (tkn.Claims.(jwt.MapClaims)["userName"] == e.Param("userName")) {
		e.Set("isOwner", true)
		return false, err
	} else {
		e.Set("isOwner", false)
		return false, err
	}

}

func VerifyExpires(m jwt.MapClaims) bool {
	return m["exp"].(float64) > float64(time.Now().Unix())
}
