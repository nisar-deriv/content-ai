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
- **report-generator**: A dummy service that keeps running to allow manual report generation.
- **report**: The service that is used to manually generate the report.

## Usage

### To Manually Generate the Report

To manually generate the weekly report, run the `report` service:

```sh
docker-compose run report
