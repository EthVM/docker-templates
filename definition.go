package main

import (
	"encoding/json"
	"fmt"
)

var definitionSchemaData = []byte(`{
	"title": "Definition",
	"type": "object",
	"properties": {
		"vars": {
            "type": "object",
			"properties": {
				"include": {
                    "type": "array",
                    "default" : []
				},
				"global": {
                    "type": "object",
                    "default" : {}
				}
			}
		},
		"templates": {
			"type": "array",
			"items": {
				"type": "object",
				"required": ["src", "dest"],
				"properties": {
					"src": {
						"type": "string"
					},
					"dest": {
						"type": "string"
					},
					"local_vars": {
						"type": "object"
					},
					"include_vars": {
						"type": "array",
						"items": {
							"type": "string"
						}
					}
				}
			}
		}
	},
	"required": ["templates"]
}`)

var includeVarsSchemaData = []byte(`{
    "title": "IncludeVars",
    "type": "object",
    "properties": {
        "vars": {
            "type": "object"
        }
    },
    "required": ["vars"]
}`)

type template struct {
	Src         string                  `toml:"src"          json:"src"`
	Dest        string                  `toml:"dest"         json:"dest"`
	LocalVars   *map[string]interface{} `toml:"local_vars"   json:"local_vars"   default:"{}"`
	IncludeVars *[]string               `toml:"include_vars" json:"include_vars" default:"[]"`
}

func (t *template) MarshalJSON() []byte {
	j, err := json.Marshal(t)
	if err != nil {
		l.Fatalf("Problem marshalling to JSON: %s", err)
	}
	return j
}

func (t *template) String() string {
	j, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		l.Fatalf("Problem marshalling to JSON: %s", err)
	}
	return fmt.Sprintf("%s", j)
}

type vars struct {
	Include *[]string               `toml:"include" json:"include" default:"[]"`
	Global  *map[string]interface{} `toml:"global"  json:"global"  default:"{}"`
}

type includeVars struct {
	Vars map[string]interface{} `toml:"vars" json:"vars"`
}

func (v *includeVars) MarshalJSON() []byte {
	j, err := json.Marshal(v)
	if err != nil {
		l.Fatalf("Problem marshalling to JSON: %s", err)
	}
	return j
}

func (v *includeVars) String() string {
	j, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		l.Fatalf("Problem marshalling to JSON: %s", err)
	}
	return fmt.Sprintf("%s", j)
}

type definition struct {
	Vars      vars       `toml:"vars"      json:"vars"`
	Templates []template `toml:"templates" json:"templates"`
}

func (d *definition) String() string {
	j, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		l.Fatalf("Problem marshalling to JSON: %s", err)
	}
	return fmt.Sprintf("%s", j)
}

func (d *definition) MarshalJSON() []byte {
	j, err := json.Marshal(d)
	if err != nil {
		l.Fatalf("Problem marshalling to JSON: %s", err)
	}
	return j
}
