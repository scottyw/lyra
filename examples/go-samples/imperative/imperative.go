package main

import (
	"encoding/json"
	"fmt"
	"github.com/lyraproj/servicesdk/lang/go/lyra"
	"os/exec"
)

func runContainer(image string, in map[string]string) (map[string]string, error) {
	// Marshal the args into JSON
	input, err := json.Marshal(&in)
	if err != nil {
		return nil, err
	}

	inputS := fmt.Sprintf("INPUT=%s ", input)
	cmd := exec.Command("docker", "run", "-t", "--env", inputS, image)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON output from the run
	out := map[string]string{}
	err = json.Unmarshal(output, &out)
	return out, err
}

type containerResult struct {
	Result string
}

func main() {

	// Workflow input is from Hiera and a constant (declared by annotations in the In struct).
	lyra.Serve(`imperative_go`, nil, &lyra.Workflow{
		Parameters: struct {
			Image string `lookup:"nebula.image"`
		}{},
		Return: containerResult{},
		Steps: map[string]lyra.Step{

			`run`: &lyra.Action{
				Do: func(input struct {
					Image string
				}) containerResult {
					out, err := runContainer(input.Image, map[string]string{"message": "howdy"})
					if err != nil {
						panic(err)
					}
					result := out["hello"]
					return containerResult{Result: result}
				}}}})

}
