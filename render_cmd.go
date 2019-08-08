package main

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
	"github.com/imdario/mergo"
	"github.com/qri-io/jsonschema"
	"gopkg.in/urfave/cli.v1"
)

func renderCmd(c *cli.Context) error {
	if len(c.Args()) < 1 {
		l.Panicln("This command requires at least one definition file path")
	}

	definitionPath := c.Args().Get(0)
	l.Println("Definition file path:", definitionPath)

	definition := unmarshalDefinitionFile(definitionPath)

	if errs := validateDefinition(definition); errs != nil {
		errJSON, _ := json.MarshalIndent(errs, "", "  ")
		l.Fatalf("Error detected in definition file: %s", errJSON)
	}

	resolveIncludeVars(definition.Vars.Include, definition.Vars.Global)

	for _, template := range definition.Templates {
		l.Printf("Processing template: %s", template.String())
		resolveIncludeVars(template.IncludeVars, template.LocalVars)
		if err := mergo.Merge(template.LocalVars, definition.Vars.Global, mergo.WithOverride); err != nil {
			l.Fatalf("Problem merging variables: %s", err)
		}
		renderFile(template)
	}

	l.Printf("Finished rendering templates!")

	return nil
}

func unmarshalDefinitionFile(path string) *definition {
	d := &definition{}

	if _, err := toml.DecodeFile(path, d); err != nil {
		l.Fatalf("Problem decoding TOML file: %s", err)
	}

	if err := defaults.Set(d); err != nil {
		l.Fatalf("Problem setting up default values: %s", err)
	}

	return d
}

func validateDefinition(d *definition) []jsonschema.ValError {
	rs := &jsonschema.RootSchema{}
	if err := json.Unmarshal(definitionSchemaData, rs); err != nil {
		l.Fatalf("Error unmarshaling schema: %s", err)
	}

	l.Printf("Validating parsed definition file: %s", d.String())

	json := d.MarshalJSON()
	if errs, _ := rs.ValidateBytes(json); len(errs) > 0 {
		return errs
	}

	return nil
}

func validateIncludeVars(v *includeVars) []jsonschema.ValError {
	rs := &jsonschema.RootSchema{}
	if err := json.Unmarshal(includeVarsSchemaData, rs); err != nil {
		l.Fatalf("Error unmarshaling schema: %s", err)
	}

	l.Printf("Validating include vars file: %s", v.String())

	json := v.MarshalJSON()
	if errs, _ := rs.ValidateBytes(json); len(errs) > 0 {
		return errs
	}

	return nil
}

func resolveIncludeVars(include *[]string, global *map[string]interface{}) {
	if include != nil && len(*include) > 0 {

		for _, path := range *include {
			v := &includeVars{}
			if _, err := toml.DecodeFile(path, v); err != nil {
				l.Fatalf("Problem decoding TOML vars file: %s", err)
			}

			l.Printf("Loaded global vars file: %s", v.String())
			validateIncludeVars(v)

			if err := mergo.Merge(global, v.Vars, mergo.WithOverride); err != nil {
				l.Fatalf("Problem merging variables: %s", err)
			}
		}

	}
}
