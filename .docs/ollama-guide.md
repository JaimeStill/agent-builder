# Ollama Guide

- [Conventions](#conventions)
- [Version](#version)
- [Pull a Model](#pull-a-model)
- [List Running Models](#list-running-models)
- [List Local Models](#list-local-models)
- [Show Model Information](#show-model-information)
- [Delete a Model](#delete-a-model)
- [Generate a Completion](#generate-a-completion)
- [Generate a Chat Completion](#generate-a-chat-completion)
- [Generate Embeddings](#generate-embeddings)

Ollama is a local inference server that runs large language models on your hardware without requiring cloud services. It provides a REST API compatible with OpenAI's format, enabling chat completions, text generation, and model management through HTTP endpoints. The server handles model downloading, quantization, and GPU acceleration automatically, supporting models like Llama, Mistral, and CodeLlama. Ollama enables private, offline AI applications with full control over data and model execution.

## Conventions
[Back to Top](#ollama-guide)

- Models follow a `model:tag` format where `model` can have an optional namespace (`example/model`). Examples: `orca-mini:3b-q8_0` and `llama3.2:3b`.
- All durations are returned in nanoseconds.
- Certain endpoints stream responses as JSON objects. Streaming can be disabled by providing `{"stream": false }` for these endpoints.

## Version
[Back to Top](#ollama-guide)

```sh
GET /api/version
```

Retrieve the Ollama version

```sh
curl http://localhost:11434/api/version
```

**Response**

```json
{
  "version": "0.10.1"
}
```

## Pull a Model
[Back to Top](#ollama-guide)

```sh
POST /api/pull
```

Download a model from the ollama library. Cancelled pulls are resumed from where they left off, and multiple calls will share the same download progress.

**Parameters**
- `model`: name of the model to pull
- `insecure`: (optional) allow insecure connections to the library
  - only use this if you are pulling from your own library during development
- `stream`: (optional) if `false` the response will be returned as a single response object rather than a stream of objects

```sh
curl http://localhost:11434/api/pull -d '{
    "model": "llama3.2:3b"
}'
```

**Response**

If `stream` is not specified, or set to `true`, a stream of JSON objects is returned.

The first object is the manifest:

```json
{
  "status": "pulling manifest"
}
```

Then there is a series of downloading responses. Until any of the download is completed, the `completed` key may not be included. The number of files to be downloaded depends on the number of layers specified in the manifest.

```json
{
  "status": "pulling digestname",
  "digest": "digestname",
  "total": 2142590208,
  "completed": 241970
}
```

After all the files are downloaded, the final responses are:

```json
{
    "status": "verifying sha256 digest"
}
{
    "status": "writing manifest"
}
{
    "status": "removing any unused layers"
}
{
    "status": "success"
}
```

If `stream` is set to false, then the response is a single JSON object:

```json
{
  "status": "success"
}
```

## List Running Models
[Back to Top](#ollama-guide)

```sh
GET /api/ps
```

List models that are currently loaded into memory.

```sh
curl http://localhost:11434/api/ps
```

**Response**

```json
{
  "models": [
    {
      "name": "mistral:latest",
      "model": "mistral:latest",
      "size": 5137025024,
      "digest": "2ae6f6dd7a3dd734790bbbf58b8909a606e0e7e97e94b7604e0aa7ae4490e6d8",
      "details": {
        "parent_model": "",
        "format": "gguf",
        "family": "llama",
        "families": [
          "llama"
        ],
        "parameter_size": "7.2B",
        "quantization_level": "Q4_0"
      },
      "expires_at": "2024-06-04T14:38:31.83753-07:00",
      "size_vram": 5137025024
    }
  ]
}
```

## List Local Models
[Back to Top](#ollama-guide)

```sh
GET /api/tags
```

List models that are available locally.

```sh
curl http://localhost:11434/api/tags
```

```json
{
  "models": [
    {
      "name": "deepseek-r1:latest",
      "model": "deepseek-r1:latest",
      "modified_at": "2025-05-10T08:06:48.639712648-07:00",
      "size": 4683075271,
      "digest": "0a8c266910232fd3291e71e5ba1e058cc5af9d411192cf88b6d30e92b6e73163",
      "details": {
        "parent_model": "",
        "format": "gguf",
        "family": "qwen2",
        "families": [
          "qwen2"
        ],
        "parameter_size": "7.6B",
        "quantization_level": "Q4_K_M"
      }
    },
    {
      "name": "llama3.2:latest",
      "model": "llama3.2:latest",
      "modified_at": "2025-05-04T17:37:44.706015396-07:00",
      "size": 2019393189,
      "digest": "a80c4f17acd55265feec403c7aef86be0c25983ab279d83f3bcd3abbcb5b8b72",
      "details": {
        "parent_model": "",
        "format": "gguf",
        "family": "llama",
        "families": [
          "llama"
        ],
        "parameter_size": "3.2B",
        "quantization_level": "Q4_K_M"
      }
    }
  ]
}
```

## Show Model Information
[Back to Top](#ollama-guide)

```sh
POST /api/show
```

Show information about a model including details, Modelfile, template, parameters, license, and system prompt.

**Parameters**

- `model`: name of the model to show
- `verbose`: (optional): if set to `true`, returns full data for verbose response fields

```sh
curl http://localhost:11434/api/show -d '{
  "model": "llava"
}'
```

**Response**

```json
{
  "modelfile": "# Modelfile generated by \"ollama show\"\n# To build a new Modelfile based on this one, replace the FROM line with:\n# FROM llava:latest\n\nFROM /Users/matt/.ollama/models/blobs/sha256:200765e1283640ffbd013184bf496e261032fa75b99498a9613be4e94d63ad52\nTEMPLATE \"\"\"{{ .System }}\nUSER: {{ .Prompt }}\nASSISTANT: \"\"\"\nPARAMETER num_ctx 4096\nPARAMETER stop \"\u003c/s\u003e\"\nPARAMETER stop \"USER:\"\nPARAMETER stop \"ASSISTANT:\"",
  "parameters": "num_keep                       24\nstop                           \"<|start_header_id|>\"\nstop                           \"<|end_header_id|>\"\nstop                           \"<|eot_id|>\"",
  "template": "{{ if .System }}<|start_header_id|>system<|end_header_id|>\n\n{{ .System }}<|eot_id|>{{ end }}{{ if .Prompt }}<|start_header_id|>user<|end_header_id|>\n\n{{ .Prompt }}<|eot_id|>{{ end }}<|start_header_id|>assistant<|end_header_id|>\n\n{{ .Response }}<|eot_id|>",
  "details": {
    "parent_model": "",
    "format": "gguf",
    "family": "llama",
    "families": [
      "llama"
    ],
    "parameter_size": "8.0B",
    "quantization_level": "Q4_0"
  },
  "model_info": {
    "general.architecture": "llama",
    "general.file_type": 2,
    "general.parameter_count": 8030261248,
    "general.quantization_version": 2,
    "llama.attention.head_count": 32,
    "llama.attention.head_count_kv": 8,
    "llama.attention.layer_norm_rms_epsilon": 0.00001,
    "llama.block_count": 32,
    "llama.context_length": 8192,
    "llama.embedding_length": 4096,
    "llama.feed_forward_length": 14336,
    "llama.rope.dimension_count": 128,
    "llama.rope.freq_base": 500000,
    "llama.vocab_size": 128256,
    "tokenizer.ggml.bos_token_id": 128000,
    "tokenizer.ggml.eos_token_id": 128009,
    "tokenizer.ggml.merges": [],            // populates if `verbose=true`
    "tokenizer.ggml.model": "gpt2",
    "tokenizer.ggml.pre": "llama-bpe",
    "tokenizer.ggml.token_type": [],        // populates if `verbose=true`
    "tokenizer.ggml.tokens": []             // populates if `verbose=true`
  },
  "capabilities": [
    "completion",
    "vision"
  ],
}
```

## Delete a Model
[Back to Top](#ollama-guide)

```sh
DELETE /api/delete
```

Delete a model and its data.

**Parameters**
- `model`: model name to delete

```sh
curl -X DELETE http://localhost:11434/api/delete -d '{
    "model": "llama3.2:3b"
}'
```

**Response**

Returns a `200 OK` if successful or a `404 Not Found` if the model doesn't exist.

## Generate a Completion
[Back to Top](#ollama-guide)

```sh
POST /api/generate
```

Generate a response for a given prompt with a provided model. This is a streaming endpoint, so there will be a series of responses. The final response object will include statistics and additional data from the request.

| Parameter    | Description                                                                                                                                                                 | Required? |
| ------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- |
| `model`      | The model name.                                                                                                                                                             | Y         |
| `prompt`     | The prompt to generate a response for.                                                                                                                                      | N         |
| `suffix`     | The text after the model response.                                                                                                                                          | N         |
| `images`     | A list of base64-encoded images (for multimodel models such as `llava`)                                                                                                     | N         |
| `think`      | (For thinking models) should the model think before responding?                                                                                                             | N         |
| `format`     | The format to return a response in. Format can be `json` or a JSON schema.                                                                                                  | N         |
| `options`    | Additional model paramters listed in the documentation for the Modelfile such as `temperature`.                                                                             | N         |
| `system`     | System message to use (override what is defined in the [Modelfile](https://github.com/ollama/ollama/blob/main/docs/modelfile.md#format)).                                   | N         |
| `template`   | The propmt template use (overrides what is defined in the [Modelfile](https://github.com/ollama/ollama/blob/main/docs/modelfile.md#format)).                                | N         |
| `stream`     | If `false`, the response will be returned as a single response object rather than a stream of objects.                                                                      | N         |
| `raw`        | If `true`, no formatting will be applied to the prompt. You may choose to use the `raw` parameter if you are specifying a full templated prompt in your request to the API. | N         |
| `keep_alive` | Controls how long the model will stay loaded into memory following the request (default: `5m`).                                                                             | N         |

The final response in the stream includes additional data about the generation:

- `total_duration`: time spent generating the response
- `load_duration`: time spent in nanoseconds loading the model
- `prompt_eval_count`: number of tokens in the prompt
- `prompt_eval_duration`: time spent in nanoseconds evaluating the prompt
- `eval_count`: number of tokens in the response
- `eval_duration`: time in nanoseconds spent generating the response
- `context`: an encoding of the conversation used in this response (this can be sent in the next request to keep a conversational memory).
- `response`: empty if the response was streamed (if not, this will contain the full response).

To calculate how fast the response is generated in tokens per second (token/s): `eval_count / eval_duration * 10^9`.

### Request (Streaming)

```sh
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "prompt": "Why is the sky blue?"
}'
```

**First Response**

```json
{
  "model": "llama3.2",
  "created_at": "2023-08-04T08:52:19.385406455-07:00",
  "response": "The",
  "done": false
}
```

**Last Response**

```json
{
  "model": "llama3.2",
  "created_at": "2023-08-04T19:22:45.499127Z",
  "response": "",
  "done": true,
  "context": [1, 2, 3],
  "total_duration": 10706818083,
  "load_duration": 6338219291,
  "prompt_eval_count": 26,
  "prompt_eval_duration": 130079000,
  "eval_count": 259,
  "eval_duration": 4232710000
}
```

### Request (No Streaming)

```sh
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "prompt": "Why is the sky blue?",
  "stream": false
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2023-08-04T19:22:45.499127Z",
  "response": "The sky is blue because it is the color of the sky.",
  "done": true,
  "context": [1, 2, 3],
  "total_duration": 5043500667,
  "load_duration": 5025959,
  "prompt_eval_count": 26,
  "prompt_eval_duration": 325953000,
  "eval_count": 290,
  "eval_duration": 4709213000
}
```

### Request (With Suffix)

```sh
curl http://localhost:11434/api/generate -d '{
  "model": "codellama:code",
  "prompt": "def compute_gcd(a, b):",
  "suffix": "    return result",
  "options": {
    "temperature": 0
  },
  "stream": false
}'
```

**Response**

```json
{
  "model": "codellama:code",
  "created_at": "2024-07-22T20:47:51.147561Z",
  "response": "\n  if a == 0:\n    return b\n  else:\n    return compute_gcd(b % a, a)\n\ndef compute_lcm(a, b):\n  result = (a * b) / compute_gcd(a, b)\n",
  "done": true,
  "done_reason": "stop",
  "context": [...],
  "total_duration": 1162761250,
  "load_duration": 6683708,
  "prompt_eval_count": 17,
  "prompt_eval_duration": 201222000,
  "eval_count": 63,
  "eval_duration": 953997000
}
```

### Request (Structured Outputs)

```sh
curl -X POST http://localhost:11434/api/generate -H "Content-Type: application/json" -d '{
  "model": "llama3.1:8b",
  "prompt": "Ollama is 22 years old and is busy saving the world. Respond using JSON",
  "stream": false,
  "format": {
    "type": "object",
    "properties": {
      "age": {
        "type": "integer"
      },
      "available": {
        "type": "boolean"
      }
    },
    "required": [
      "age",
      "available"
    ]
  }
}'
```

**Response**

```json
{
  "model": "llama3.1:8b",
  "created_at": "2024-12-06T00:48:09.983619Z",
  "response": "{\n  \"age\": 22,\n  \"available\": true\n}",
  "done": true,
  "done_reason": "stop",
  "context": [1, 2, 3],
  "total_duration": 1075509083,
  "load_duration": 567678166,
  "prompt_eval_count": 28,
  "prompt_eval_duration": 236000000,
  "eval_count": 16,
  "eval_duration": 269000000
}
```

### Request (JSON Mode)

> [!IMPORTANT]
> When `format` is set to `json`, the output will always be a well-formed JSON object. It's important to instruct the model to respond in JSON.

```sh
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "prompt": "What color is the sky at different times of the day? Respond using JSON",
  "format": "json",
  "stream": false
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2023-11-09T21:07:55.186497Z",
  "response": "{\n\"morning\": {\n\"color\": \"blue\"\n},\n\"noon\": {\n\"color\": \"blue-gray\"\n},\n\"afternoon\": {\n\"color\": \"warm gray\"\n},\n\"evening\": {\n\"color\": \"orange\"\n}\n}\n",
  "done": true,
  "context": [1, 2, 3],
  "total_duration": 4648158584,
  "load_duration": 4071084,
  "prompt_eval_count": 36,
  "prompt_eval_duration": 439038000,
  "eval_count": 180,
  "eval_duration": 4196918000
}
```

The value of `response` will be a string containing JSON similar to:

```json
{
  "morning": {
    "color": "blue"
  },
  "noon": {
    "color": "blue-gray"
  },
  "afternoon": {
    "color": "warm gray"
  },
  "evening": {
    "color": "orange"
  }
}
```

### Request (With Images)

To submit images to multimodal models such as `llava` or `bakllava`, provide a list of base64-encoded `images`:

```sh
curl http://localhost:11434/api/generate -d '{
  "model": "llava",
  "prompt":"What is in this picture?",
  "stream": false,
  "images": ["iVBORw0KGgoA..."]
}'
```

**Response**

```json
{
  "model": "llava",
  "created_at": "2023-11-03T15:36:02.583064Z",
  "response": "A happy cartoon character, which is cute and cheerful.",
  "done": true,
  "context": [1, 2, 3],
  "total_duration": 2938432250,
  "load_duration": 2559292,
  "prompt_eval_count": 1,
  "prompt_eval_duration": 2195557000,
  "eval_count": 44,
  "eval_duration": 736432000
}
```

### Request (Raw Mode)

In some cases, you may wish to bypass the templating system and provide a full prompt. In this case, you can use the `raw` parameter to disable templating. Also note that raw mode will not return a context.

```sh
curl http://localhost:11434/api/generate -d '{
  "model": "mistral:7b",
  "prompt": "[INST] why is the sky blue? [/INST]",
  "raw": true,
  "stream": false
}'
```

**Response**

```json
{
  "model": "mistral:7b",
  "created_at": "2025-08-20T14:11:51.567090766Z",
  "response": " The sky appears blue because of a process called Rayleigh scattering. When sunlight, which is white light composed of different colors (red, orange, yellow, green, blue, indigo, and violet), travels through Earth's atmosphere, the smaller molecules of gases like nitrogen and oxygen in the air scatter short-wavelength light, such as blue and violet light, to a far greater extent than longer wavelengths, such as red light.\n\nHowever, human eyes are more sensitive to blue light and less sensitive to violet light. Additionally, some of the violet light gets absorbed by the ozone layer in Earth's atmosphere. The combined effect is that we perceive the sky as blue during a clear day.",
  "done": true,
  "done_reason": "stop",
  "total_duration": 4048232469,
  "load_duration": 1722704234,
  "prompt_eval_count": 11,
  "prompt_eval_duration": 148014742,
  "eval_count": 155,
  "eval_duration": 2176660570
}
```

### Request (Reproducible Outputs)

For reproducible outputs, set `seed` to a number:

```sh
curl http://localhost:11434/api/generate -d '{
  "model": "mistral",
  "prompt": "Why is the sky blue?",
  "options": {
    "seed": 123
  }
}'
```

**Response**

```json
{
  "model": "mistral",
  "created_at": "2023-11-03T15:36:02.583064Z",
  "response": " The sky appears blue because of a phenomenon called Rayleigh scattering.",
  "done": true,
  "total_duration": 8493852375,
  "load_duration": 6589624375,
  "prompt_eval_count": 14,
  "prompt_eval_duration": 119039000,
  "eval_count": 110,
  "eval_duration": 1779061000
}
```

### Request (With Options)

If you want to set custom options for the model at runtime rather than in the Modelfile, you can do so with the `options` parameter. This example sets every available option, but you can set any of them individually and omit the ones you do not want to override.

```sh
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "prompt": "Why is the sky blue?",
  "stream": false,
  "options": {
    "num_keep": 5,
    "seed": 42,
    "num_predict": 100,
    "top_k": 20,
    "top_p": 0.9,
    "min_p": 0.0,
    "typical_p": 0.7,
    "repeat_last_n": 33,
    "temperature": 0.8,
    "repeat_penalty": 1.2,
    "presence_penalty": 1.5,
    "frequency_penalty": 1.0,
    "penalize_newline": true,
    "stop": ["\n", "user:"],
    "numa": false,
    "num_ctx": 1024,
    "num_batch": 2,
    "num_gpu": 1,
    "main_gpu": 0,
    "use_mmap": true,
    "num_thread": 8
  }
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2023-08-04T19:22:45.499127Z",
  "response": "The sky is blue because it is the color of the sky.",
  "done": true,
  "context": [1, 2, 3],
  "total_duration": 4935886791,
  "load_duration": 534986708,
  "prompt_eval_count": 26,
  "prompt_eval_duration": 107345000,
  "eval_count": 237,
  "eval_duration": 4289432000
}
```

### Load a Model

If an empty prompt is provided, the model will be loaded into memory.

```sh
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2"
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2023-12-18T19:52:07.071755Z",
  "response": "",
  "done": true
}
```

### Unload a Model

If an empty prompt is provided and the `keep_alive` parameter is set to `0`, a model will be unloaded from memory.

```sh
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "keep_alive": 0
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2024-09-12T03:54:03.516566Z",
  "response": "",
  "done": true,
  "done_reason": "unload"
}
```

## Generate a Chat Completion
[Back to Top](#ollama-guide)

```sh
POST /api/chat
```

Generate the next message in a chat with a provided model. This is a streaming endpoint, so there will be a series of responses. Streaming can be disabled using `"stream": false`. The final response object will include statistics and additional data from the request.

| Parameter    | Description                                                                                                                                                             | Required |
| ------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| `model`      | The model name.                                                                                                                                                         | Y        |
| `messages`   | The messages of the chat (this can be used to keep a chat memory).                                                                                                      | N        |
| `tools`      | List of tools in JSON for the model to use if supported.                                                                                                                | N        |
| `think`      | (For thinking models) should the model think before responding?                                                                                                         | N        |
| `format`     | The format to return a response in. Format can be `json` or a JSON schema.                                                                                              | N        |
| `options`    | Additional model parameters listed in the documentation for the [Modelfile](https://github.com/ollama/ollama/blob/main/docs/modelfile.md#format) such as `temperature`. | N        |
| `stream`     | If `false`, the response will be returned as a single response object, rather than a stream of objects.                                                                 | N        |
| `keep_alive` | Controls how long the model will stay loaded into memory following the request (default: `5m`).                                                                         | N        |

The `message` object has the following fields:

- `role`: the role of the message, either `system`, `user`, `assistant`, or `tool`.
- `content`: the content of the message.
- `thinking`: (for thinking models) the model's thinking process.
- `images`: (optional) a list of images to include in the message (for multimodal models such as `llava`).
- `tool_calls`: (optional) a list of tools in JSON that the model wants to use.
- `tool_name`: (optonal) the name of the tool that was executed to inform the model of the result.

Tool calling is supported by providing a list of tools in the `tools` parameter. The model will generate a response that includes a list of tool calls. Models can also explain the result of the tool call in the response. See [models with tool calling capabilities](https://ollama.com/search?c=tool).

### Chat Request (Streaming)

```sh
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": [
    {
      "role": "user",
      "content": "why is the sky blue?"
    }
  ]
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2023-08-04T08:52:19.385406455-07:00",
  "message": {
    "role": "assistant",
    "content": "The",
    "images": null
  },
  "done": false
}
```

**Final Response**

```json
{
  "model": "llama3.2",
  "created_at": "2023-08-04T19:22:45.499127Z",
  "message": {
    "role": "assistant",
    "content": ""
  },
  "done": true,
  "total_duration": 4883583458,
  "load_duration": 1334875,
  "prompt_eval_count": 26,
  "prompt_eval_duration": 342546000,
  "eval_count": 282,
  "eval_duration": 4535599000
}
```

### Chat Request (Streaming with Tools)

```sh
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": [
    {
      "role": "user",
      "content": "what is the weather in tokyo?"
    }
  ],
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "get_weather",
        "description": "Get the weather in a given city",
        "parameters": {
          "type": "object",
          "properties": {
            "city": {
              "type": "string",
              "description": "The city to get the weather for"
            }
          },
          "required": ["city"]
        }
      }
    }
  ],
  "stream": true
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2025-07-07T20:22:19.184789Z",
  "message": {
    "role": "assistant",
    "content": "",
    "tool_calls": [
      {
        "function": {
          "name": "get_weather",
          "arguments": {
            "city": "Tokyo"
          }
        }
      }
    ]
  },
  "done": false
}
```

**Final Response**

```json
{
  "model": "llama3.2",
  "created_at": "2025-07-07T20:22:19.19314Z",
  "message": {
    "role": "assistant",
    "content": "The current weather in Tokyo is quite cool and wet. It's 18°C with light rain falling, so you'd definitely want to bring an umbrella if you're heading out. The humidity is quite high at 78%, and there's a moderate wind of 12 km/h coming from the northeast. It's typical autumn weather for Tokyo - perfect for staying cozy indoors with a warm drink!"
  },
  "done_reason": "stop",
  "done": true,
  "total_duration": 182242375,
  "load_duration": 41295167,
  "prompt_eval_count": 169,
  "prompt_eval_duration": 24573166,
  "eval_count": 15,
  "eval_duration": 115959084
}
```

### Chat Request (No Streaming)

```sh
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": [
    {
      "role": "user",
      "content": "why is the sky blue?"
    }
  ],
  "stream": false
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2023-12-12T14:13:43.416799Z",
  "message": {
    "role": "assistant",
    "content": "The sky is blue because it is the color of the sky."
  },
  "done": true,
  "total_duration": 5191566416,
  "load_duration": 2154458,
  "prompt_eval_count": 26,
  "prompt_eval_duration": 383809000,
  "eval_count": 298,
  "eval_duration": 4799921000
}
```

### Chat Request (No Streaming with Tools)

```sh
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": [
    {
      "role": "user",
      "content": "what is the weather in tokyo?"
    }
  ],
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "get_weather",
        "description": "Get the weather in a given city",
        "parameters": {
          "type": "object",
          "properties": {
            "city": {
              "type": "string",
              "description": "The city to get the weather for"
            }
          },
          "required": ["city"]
        }
      }
    }
  ],
  "stream": false
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2025-07-07T20:32:53.844124Z",
  "message": {
    "role": "assistant",
    "content": "The current weather in Tokyo is quite cool and wet. It's 18°C with light rain falling, so you'd definitely want to bring an umbrella if you're heading out. The humidity is quite high at 78%, and there's a moderate wind of 12 km/h coming from the northeast. It's typical autumn weather for Tokyo - perfect for staying cozy indoors with a warm drink!",
    "tool_calls": [
      {
        "function": {
          "name": "get_weather",
          "arguments": {
            "city": "Tokyo"
          }
        }
      }
    ]
  },
  "done_reason": "stop",
  "done": true,
  "total_duration": 3244883583,
  "load_duration": 2969184542,
  "prompt_eval_count": 169,
  "prompt_eval_duration": 141656333,
  "eval_count": 18,
  "eval_duration": 133293625
}
```

### Chat Request (Structured Outputs)

```sh
curl -X POST http://localhost:11434/api/chat -H "Content-Type: application/json" -d '{
  "model": "llama3.1",
  "messages": [{"role": "user", "content": "Ollama is 22 years old and busy saving the world. Return a JSON object with the age and availability."}],
  "stream": false,
  "format": {
    "type": "object",
    "properties": {
      "age": {
        "type": "integer"
      },
      "available": {
        "type": "boolean"
      }
    },
    "required": [
      "age",
      "available"
    ]
  },
  "options": {
    "temperature": 0
  }
}'
```

**Response**

```json
{
  "model": "llama3.1",
  "created_at": "2024-12-06T00:46:58.265747Z",
  "message": {
    "role": "assistant",
    "content": "{\"age\": 22, \"available\": false}"
  },
  "done_reason": "stop",
  "done": true,
  "total_duration": 2254970291,
  "load_duration": 574751416,
  "prompt_eval_count": 34,
  "prompt_eval_duration": 1502000000,
  "eval_count": 12,
  "eval_duration": 175000000
}
```

### Chat Request (With History)

Send a chat message with a conversation history. You can use this same approach to start the conversation using multi-shot or chain-of-thought prompting.

```sh
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": [
    {
      "role": "user",
      "content": "why is the sky blue?"
    },
    {
      "role": "assistant",
      "content": "due to rayleigh scattering."
    },
    {
      "role": "user",
      "content": "how is that different than mie scattering?"
    }
  ]
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2023-08-04T08:52:19.385406455-07:00",
  "message": {
    "role": "assistant",
    "content": "The"
  },
  "done": false
}
```

**Final Response**

```json
{
  "model": "llama3.2",
  "created_at": "2023-08-04T19:22:45.499127Z",
  "done": true,
  "total_duration": 8113331500,
  "load_duration": 6396458,
  "prompt_eval_count": 61,
  "prompt_eval_duration": 398801000,
  "eval_count": 468,
  "eval_duration": 7701267000
}
```

### Chat Request (With History and Tools)

```sh
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": [
    {
      "role": "user",
      "content": "what is the weather in Toronto?"
    },
    // the message from the model appended to history
    {
      "role": "assistant",
      "content": "",
      "tool_calls": [
        {
          "function": {
            "name": "get_temperature",
            "arguments": {
              "city": "Toronto"
            }
          },
        }
      ]
    },
    // the tool call result appended to history
    {
      "role": "tool",
      "content": "11 degrees celsius",
      "tool_name": "get_temperature",
    }
  ],
  "stream": false,
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "get_weather",
        "description": "Get the weather in a given city",
        "parameters": {
          "type": "object",
          "properties": {
            "city": {
              "type": "string",
              "description": "The city to get the weather for"
            }
          },
          "required": ["city"]
        }
      }
    }
  ]
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2025-07-07T20:43:37.688511Z",
  "message": {
    "role": "assistant",
    "content": "The current temperature in Toronto is 11°C."
  },
  "done_reason": "stop",
  "done": true,
  "total_duration": 890771750,
  "load_duration": 707634750,
  "prompt_eval_count": 94,
  "prompt_eval_duration": 91703208,
  "eval_count": 11,
  "eval_duration": 90282125
}
```

### Chat Request (With Images)

```sh
curl http://localhost:11434/api/chat -d '{
  "model": "llava",
  "messages": [
    {
      "role": "user",
      "content": "what is in this image?",
      "images": ["iVBORw0KGgoA..."]
    }
  ]
}'
```

**Response**

```json
{
  "model": "llava",
  "created_at": "2023-12-13T22:42:50.203334Z",
  "message": {
    "role": "assistant",
    "content": " The image features a cute, little pig with an angry facial expression. It's wearing a heart on its shirt and is waving in the air. This scene appears to be part of a drawing or sketching project.",
    "images": null
  },
  "done": true,
  "total_duration": 1668506709,
  "load_duration": 1986209,
  "prompt_eval_count": 26,
  "prompt_eval_duration": 359682000,
  "eval_count": 83,
  "eval_duration": 1303285000
}
```

### Chat Request (Reproducible Outputs)

```sh
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": [
    {
      "role": "user",
      "content": "Hello!"
    }
  ],
  "options": {
    "seed": 101,
    "temperature": 0
  }
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2023-12-12T14:13:43.416799Z",
  "message": {
    "role": "assistant",
    "content": "Hello! How are you today?"
  },
  "done": true,
  "total_duration": 5191566416,
  "load_duration": 2154458,
  "prompt_eval_count": 26,
  "prompt_eval_duration": 383809000,
  "eval_count": 298,
  "eval_duration": 4799921000
}
```

### Chat Request (With Tools)

```sh
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": [
    {
      "role": "user",
      "content": "What is the weather today in Paris?"
    }
  ],
  "stream": false,
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "get_current_weather",
        "description": "Get the current weather for a location",
        "parameters": {
          "type": "object",
          "properties": {
            "location": {
              "type": "string",
              "description": "The location to get the weather for, e.g. San Francisco, CA"
            },
            "format": {
              "type": "string",
              "description": "The format to return the weather in, e.g. 'celsius' or 'fahrenheit'",
              "enum": ["celsius", "fahrenheit"]
            }
          },
          "required": ["location", "format"]
        }
      }
    }
  ]
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2024-07-22T20:33:28.123648Z",
  "message": {
    "role": "assistant",
    "content": "Currently 68°F with light rain and broken clouds.",
    "tool_calls": [
      {
        "function": {
          "name": "get_current_weather",
          "arguments": {
            "format": "celsius",
            "location": "Paris, FR"
          }
        }
      }
    ]
  },
  "done_reason": "stop",
  "done": true,
  "total_duration": 885095291,
  "load_duration": 3753500,
  "prompt_eval_count": 122,
  "prompt_eval_duration": 328493000,
  "eval_count": 33,
  "eval_duration": 552222000
}
```

### Load a Model

If the messages array is empty, the model will be loaded into memory.

```sh
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": []
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2024-09-12T21:17:29.110811Z",
  "message": {
    "role": "assistant",
    "content": ""
  },
  "done_reason": "load",
  "done": true
}
```

### Unload a Model

If the messages array is empty and the `keep_alive` parameter is set to `0`, a model will be unloaded from memory.

```sh
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": [],
  "keep_alive": 0
}'
```

**Response**

```json
{
  "model": "llama3.2",
  "created_at": "2024-09-12T21:33:17.547535Z",
  "message": {
    "role": "assistant",
    "content": ""
  },
  "done_reason": "unload",
  "done": true
}
```

## Generate Embeddings
[Back to Top](#ollama-guide)

```sh
POST /api/embed
```

Generate embeddings from a model.

**Parameters**
- `model`: name of the model to generate embeddings from
- `input`: text or list of text to generate embeddings for

**Advanced Parameters**
- `truncate`: truncates the end of each input to fit within context length
  - returns error if `false` and the context length is exceeded
  - defaults to `true`
- `options`: additional model parameters listed in the documentation for the [Modelfile](https://github.com/ollama/ollama/blob/main/docs/modelfile.md#valid-parameters-and-values) such as `temperature`
- `keep_alive`: controls how long the model will stay loaded into memory following the request (default: `5m`)

```sh
curl http://localhost:11434/api/embed -d '{
  "model": "all-minilm",
  "input": "Why is the sky blue?"
}'
```

**Response**

```json
{
  "model": "all-minilm",
  "embeddings": [[
    0.010071029, -0.0017594862, 0.05007221, 0.04692972, 0.054916814,
    0.008599704, 0.105441414, -0.025878139, 0.12958129, 0.031952348
  ]],
  "total_duration": 14143917,
  "load_duration": 1019500,
  "prompt_eval_count": 8
}
```

### Request (Multiple Input)

```sh
curl http://localhost:11434/api/embed -d '{
  "model": "all-minilm",
  "input": ["Why is the sky blue?", "Why is the grass green?"]
}'
```

**Response**

```json
{
  "model": "all-minilm",
  "embeddings": [[
    0.010071029, -0.0017594862, 0.05007221, 0.04692972, 0.054916814,
    0.008599704, 0.105441414, -0.025878139, 0.12958129, 0.031952348
  ],[
    -0.0098027075, 0.06042469, 0.025257962, -0.006364387, 0.07272725,
    0.017194884, 0.09032035, -0.051705178, 0.09951512, 0.09072481
  ]]
}
```
