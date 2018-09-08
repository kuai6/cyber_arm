package command

type Command struct {
	Name      string   `json:"name"`
	Arguments []string `json:"arguments,omitempty"`
}
