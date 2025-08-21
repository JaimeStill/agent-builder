package ollama

import "time"

type VersionResponse struct {
	Version string `json:"version"`
}

type PullRequest struct {
	Model    string `json:"model"`
	Insecure bool   `json:"insecure,omitempty"`
	Stream   *bool  `json:"stream,omitempty"`
}

func (pr *PullRequest) GetStream() bool {
	if pr.Stream == nil {
		return true
	}

	return *pr.Stream
}

type PullResponse struct {
	Status    string  `json:"status"`
	Digest    *string `json:"digest,omitempty"`
	Total     *int64  `json:"total,omitempty"`
	Completed *int64  `json:"completed,omitempty"`
}

type ModelResponse struct {
	Models []OllamaModel `json:"models"`
}

type OllamaModel struct {
	Name      string       `json:"name"`
	Model     string       `json:"model"`
	Size      int64        `json:"size"`
	Digest    string       `json:"digest"`
	ExpiresAt time.Time    `json:"expires_at"`
	SizeVRAM  int64        `json:"size_vram"`
	Details   ModelDetails `json:"details"`
}

type ModelDetails struct {
	ParentModel       string   `json:"parent_model"`
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families,omitempty"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}

type ModelTensor struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Shape []int64 `json:"shape"`
}

type ShowRequest struct {
	Model   string `json:"model"`
	Verbose bool   `json:"verbose,omitempty"`
}

type ShowResponse struct {
	License      string         `json:"license"`
	Modelfile    string         `json:"modelfile"`
	Parameters   string         `json:"parameters"`
	Template     string         `json:"template"`
	Details      ModelDetails   `json:"details"`
	ModelInfo    map[string]any `json:"model_info,omitempty"`
	Tensors      []ModelTensor  `json:"tensors,omitempty"`
	Capabilities []string       `json:"capabilities,omitempty"`
	ModifiedAt   time.Time      `json:"modified_at"`
}
