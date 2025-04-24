package model

type Tool struct {
	Id            string `json:"id,omitempty"`
	Type          string `json:"type,omitempty"` // when splicing claude tools stream messages, it is empty
	DisplayWidth  int    `json:"display_width,omitempty"`
	DisplayHeight int    `json:"display_height,omitempty"`
	Environment   string `json:"environment,omitempty"`

	Function *Function `json:"function,omitempty"`
}

type Function struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`       // when splicing claude tools stream messages, it is empty
	Parameters  any    `json:"parameters,omitempty"` // request
	Arguments   any    `json:"arguments,omitempty"`  // response
}
