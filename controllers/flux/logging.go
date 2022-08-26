/* logging has debugging and TBA logging helpers
 */

package controllers

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

// Debugging to write yaml to yaml directory at root
func saveDebugYaml(obj interface{}, basename string) {

	// This can be removed when we are happy
	d, err := yaml.Marshal(&obj)
	if err != nil {
		fmt.Printf(" 🪵 Error marshalling %s: %v", basename, err)
	} else {
		fileName := filepath.Join("yaml", basename)
		err = ioutil.WriteFile(fileName, d, 0644)
		if err != nil {
			fmt.Println(" 🪵 Unable to write data into the file")
		} else {
			fmt.Println(" 🪵 Wrote yaml to", fileName)
		}
	}
}
