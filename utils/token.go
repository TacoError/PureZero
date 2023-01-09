package utils

import (
	"encoding/json"
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"golang.org/x/oauth2"
	"io/ioutil"
)

func TokenSource() oauth2.TokenSource {
	check := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	token := new(oauth2.Token)
	tokenData, err := ioutil.ReadFile("token.tok")
	if err == nil {
		_ = json.Unmarshal(tokenData, token)
		fmt.Println("Found cached token.")
	} else {
		fmt.Println("No cached token, requesting...")
		token, err = auth.RequestLiveToken()
		check(err)
	}
	src := auth.RefreshTokenSource(token)
	_, err = src.Token()
	if err != nil {
		token, err = auth.RequestLiveToken()
		check(err)
		src = auth.RefreshTokenSource(token)
	} else {
		fmt.Println("Token good, continuing...")
	}
	tok, _ := src.Token()
	b, _ := json.Marshal(tok)
	_ = ioutil.WriteFile("token.tok", b, 0644)
	return src
}
