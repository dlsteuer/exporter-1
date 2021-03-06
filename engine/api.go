package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func apiCall(path string) ([]byte, error) {
	url := fmt.Sprintf("https://engine.battlesnake.io/%s", path)
	client := http.Client{}

	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Got non 200 from engine: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getFrames(gameID string, offset int, limit int) ([]*GameFrame, error) {
	path := fmt.Sprintf("games/%s/frames?offset=%d&limit=%d", gameID, offset, limit)
	body, err := apiCall(path)
	if err != nil {
		return nil, err
	}

	response := gameFramesResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Frames, nil
}

// GetGame is commented
func GetGame(gameID string) (*Game, error) {
	path := fmt.Sprintf("games/%s", gameID)
	body, err := apiCall(path)
	if err != nil {
		return nil, err
	}

	response := gameStatusResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Game, nil
}

// GetGameFrame is commented
func GetGameFrame(gameID string, frameNum int) (*GameFrame, error) {
	gameFrames, err := getFrames(gameID, frameNum, 1)
	if err != nil {
		return nil, err
	}

	return gameFrames[0], nil
}

// GetGameFrames is commented
func GetGameFrames(gameID string, offset, frameCount int) ([]*GameFrame, error) {
	var gameFrames []*GameFrame

	batchSize := 100
	for {
		if frameCount > 0 && frameCount < batchSize {
			batchSize = frameCount
		}
		newFrames, err := getFrames(gameID, offset, batchSize)
		if err != nil {
			return nil, err
		}

		gameFrames = append(gameFrames, newFrames...)
		if len(newFrames) < batchSize || len(newFrames) == frameCount {
			break
		}
		offset += batchSize
		frameCount -= batchSize
	}

	return gameFrames, nil
}
