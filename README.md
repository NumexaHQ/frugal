<a>
  <h1 align="center">Numexa</h1>
    <p align="center">
      <b>🚀 An open-source Cost & Resource Optimization Platform for LLMs. Be Frugal! 💰</b>
    </p>
</a>

[![](https://img.shields.io/badge/Visit%20Us-app.numexa.io-brightgreen)](https://app.numexa.io)
[![](https://img.shields.io/badge/Join%20our%20community-Discord-blue)](https://discord.gg/mVBMKVCv)
[![](https://img.shields.io/badge/View%20Documentation-Docs-yellow)](https://docs.numexa.io/)

## Introduction

Numexa is an AI-driven cost and resource optimization tool designed to enhance operational efficiency. It achieves this by leveraging contextual insights derived from usage metrics. Numexa employs cutting-edge techniques such as intelligent caching and data retrieval, harnessing the power of vector databases to streamline operations. Explore how Numexa can revolutionize your resource management and cost-saving endeavors.

## Features

- 📝 Model agnostic functionality records unlimited requests from various providers like OpenAI, Cohere, Anthropic and more.

- 📋 Model management

- 🔔 Alerting & Notification with predefined policies, like error rate, threshold, cost, etc.

- 💾 Caching, Custom Rate Limits, and Retries,

- 📊 Track costs and latencies by users, applications, and endpoints

- 🔜 (Coming soon) Intellegient caching and data retrieval

- 🔜 (Coming soon) Cost and resource optimization



# Development Setup 💻

### Prerequisites
Before you begin, ensure you have the following installed on your system:

- **Git**: [Installation Guide](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- **Docker**: [Installation Guide](https://docs.docker.com/get-docker/)
- **Docker Compose**: [Installation Guide](https://docs.docker.com/compose/install/)
- **Make**: [Installation Guide](https://www.gnu.org/software/make/)

### Getting Started

1. **Clone the Repository:**

   ```bash
   git clone <repository_url>
   cd <repository_directory>

2. **Build and Start the Services:**
Run the following commands to build and start the project services
   ```bash
   make all
   docker compose -f docker-compose.dev.yaml up -d
3. **Verify Services:**
After running the above commands, your project services should be up and running. You can verify this by checking the logs

# Community 🤝
Join our [#Discord](https://discord.gg/mVBMKVCv) or drop email at hello@numexa.io
