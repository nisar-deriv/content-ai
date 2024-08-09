# Content Writer AI

This project provides a content writer AI that processes weekly updates from different teams, enhances the text using NLP models, stores the updates, and compiles a final report. The project supports a local Ollama server for text enhancement.

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

    Set your environment variables as needed.

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