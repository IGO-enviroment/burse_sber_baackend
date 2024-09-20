package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Read() (Settings, error) {
	fileBytes, err := ioutil.ReadFile("./.config/" + "local" + ".json")
	if err != nil {
		return Settings{}, err
	}

	var result Settings
	if err := json.Unmarshal(fileBytes, &result); err != nil {
		return Settings{}, fmt.Errorf("unmarshalling error: %v", err)
	}

	return result, nil
}
