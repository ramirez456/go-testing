package controller

import (
	"catching-pokemons/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestGetPokemonFromPokeApiSuccess(t *testing.T) {
	c := require.New(t)

	pokemon, err := GetPokemonFromPokeApi("4")
	c.NoError(err)

	body, err := ioutil.ReadFile("samples/poke_api_read.json")
	c.NoError(err)

	var expected models.PokeApiPokemonResponse

	err = json.Unmarshal([]byte(body), &expected)
	c.NoError(err)
	c.Equal(expected, pokemon)

}

func TestGetPokemonFromPokeApiSuccessWithMock(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "4"
	body, err := ioutil.ReadFile("samples/poke_api_response.json")
	c.NoError(err)

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)

	httpmock.RegisterResponder("GET", request, httpmock.NewStringResponder(200, string(body)))

	pokemon, err := GetPokemonFromPokeApi(id)
	c.NoError(err)

	var expected models.PokeApiPokemonResponse

	err = json.Unmarshal([]byte(body), &expected)
	c.NoError(err)
	c.Equal(expected, pokemon)

}

func TestGetPokeApiIntenalServerError(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	id := "4"
	body, err := ioutil.ReadFile("samples/poke_api_response.json")
	c.NoError(err)

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)

	httpmock.RegisterResponder("GET", request, httpmock.NewStringResponder(500, string(body)))

	_, err = GetPokemonFromPokeApi(id)
	c.NotNil(err)
	c.EqualError(ErrPokeApiFailure, err.Error())

}

func TestGetPokemosFromPokeApiNotFoundError(t *testing.T) {

	c := require.New(t)
	httpmock.Activate()

	defer httpmock.DeactivateAndReset()

	id := "bulbasaur"

	request := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)

	httpmock.RegisterResponder("GET", request, httpmock.NewStringResponder(404, ""))

	_, err := GetPokemonFromPokeApi(id)
	c.NotNil(err)
	c.EqualError(ErrPokemonNotFound, err.Error())

}

func TestGetPokemon(t *testing.T) {
	c := require.New(t)
	r, err := http.NewRequest("GET", "/pokemon/{id}", nil)
	c.NoError(err)

	w := httptest.NewRecorder()

	vars := map[string]string{
		"id": "4",
	}

	r = mux.SetURLVars(r, vars)

	GetPokemon(w, r)

	expected, err := ioutil.ReadFile("samples/api_response.json")
	c.NoError(err)

	var expectedPokemon models.Pokemon

	err = json.Unmarshal([]byte(expected), &expectedPokemon)
	c.NoError(err)

	var actualPokemon models.Pokemon

	err = json.Unmarshal([]byte(w.Body.String()), &actualPokemon)

	c.Equal(http.StatusOK, w.Code)
	c.Equal(expectedPokemon, actualPokemon)
}
