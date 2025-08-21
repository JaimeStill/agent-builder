# Agent Builder

## Tasks

Implement the following Ollama endpoints:
- Delete -> `/api/delete`
- Generate -> `/api/generate`
- Chat -> `/api/chat`
- Embed -> `/api/embed`

Will most likely embed the handlers into a sub-structure so that the ollama handlers aren't polluting the App method namespace.

### App Features

- Build Agent Package
  - PostgreSQL infrastructure
  - agent structure
  - agent logic with persistence
- Ollama Server Interface
  - list installed models
  - install model with validation
  - remove model
  - create agent from model
- Agents Interface
  - list available agents
  - create agent
  - remove agent
  - open agent
- Agent Interface
  - primarily chat driven
  - able to leverage model capabilities (vision, tools, etc.)
  - settings interface for adjusting agent values

## Commands

### Remote Into Ollama Container

```sh
docker exec -it agent-builder-ollama bash
```

## References

- [Ollama](https://ollama.com/)
  - [Docs](https://github.com/ollama/ollama/tree/main/docs)
  - [Models](https://ollama.com/search)
