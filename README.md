# Content Writer AI

This project provides a content writer AI that processes weekly updates from different teams, enhances the text using NLP models, stores the updates, and compiles a final report. The project supports both OpenAI GPT-3/4 and a local Ollama server for text enhancement.

## Prerequisites

- Docker
- Docker Compose

## Setup

1. **Clone the repository:**

    ```sh
    git clone https://github.com/nisar-deriv/content-ai.git
    cd content-writer-ai
    ```

2. **Set Environment Variables:**

    Set your OpenAI API key and the `USE_OLLAMA` environment variable.

    ```sh
    export OPENAI_API_KEY=your_openai_api_key
    export USE_OLLAMA=true  # or false
    ```

3. **Build and Run the Docker Containers:**

    ```sh
    docker-compose up --build
    ```

## Services

- **web**: Handles incoming requests and processes them.

## Usage

### To Manually Fetch Updates From Slack channels

To manually Fetch Updates from slack, send a GET request to the `/fetch-updates` endpoint:

```sh
curl http://localhost:8080/fetch-updates

### To Manually Generate the Report

To manually generate the weekly report, send a GET request to the `/generate-report-ai` endpoint:

```sh
curl http://localhost:8080/generate-report-ai