package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-cmp/cmp"
)

func verify_modules_names(modules_names []string) {
	if len(modules_names) > 0 {
		for _, module_name := range modules_names {
			// TODO: Compile instead
			matched, err := regexp.MatchString("[a-zA-Z0-9\\-\\_]+", module_name)
			if err != nil {
				log.Fatal(err)
			}

			if !matched {
				log.Fatal(module_name + " is invalid")
			}
		}
	}
}

func main() {
	var add_module_names string
	var remove_module_names string
	flag.StringVar(&add_module_names, "add-module", "", "Modules to append. Separated by space.")
	flag.StringVar(&remove_module_names, "remove-module", "", "Modules to remove. Separated by space.")
	flag.Parse()

	if len(add_module_names) > 0 && len(remove_module_names) > 0 {
		log.Fatal("add-module and remove-module are mutual exclusive.")
	} else if len(add_module_names) == 0 && len(remove_module_names) == 0 {
		log.Fatal("Either of add-module or remove-module must be set.")
	}

	var add_modules []string
	var remove_modules []string
	if len(add_module_names) > 0 {
		add_modules = strings.Split(add_module_names, " ")
		verify_modules_names(add_modules)
	} else if len(remove_module_names) > 0 {
		remove_modules = strings.Split(remove_module_names, " ")
		verify_modules_names(remove_modules)
	}

	defaultConfig := "/etc/mkinitcpio.conf"
	data, err := os.ReadFile(defaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	config_string := bytes.NewBuffer(data).String()
	config_lines := strings.Split(config_string, "\n")

	modules_index := -1
	for index, line := range config_lines {
		if strings.HasPrefix(line, "MODULES=(") {
			modules_index = index
		}
	}

	if modules_index > -1 {
		// TODO: Consider regex instead
		modules_string := strings.Replace(config_lines[modules_index], "MODULES=", "", 1)
		modules_string = strings.TrimLeft(modules_string, "(")
		modules_string = strings.TrimRight(modules_string, ")")

		new_modules := strings.Split(modules_string, " ")

		if len(add_module_names) > 0 {
			new_modules = append(new_modules, add_modules...)
		} else if len(remove_module_names) > 0 {
			// TODO: Find out how to update the slice
		}

		new_modules_string := strings.Join(new_modules, " ")
		new_modules_string = "MODULES=(" + new_modules_string + ")"

		config_lines[modules_index] = new_modules_string
		new_config_string := strings.Join(config_lines, "\n")

		if diff := cmp.Diff(config_string, new_config_string); diff != "" {
			log.Printf("Diff (-old +new):\n%s", diff)
		}
		
		// TODO: Require confirmation
		// temp_file := "/tmp/mkinitcpio.conf.generated"
		// os.WriteFile(temp_file, []byte(config_string), 0666)
	}
}
