# Strategic Insight Analyst

This project is a full-stack web application designed for Omara Technologies. It allows users to upload documents, ask questions about them, and receive intelligent answers. The application is built with a Go backend and a Next.js frontend.

## About The Project

The Strategic Insight Analyst is a powerful tool for document analysis. Users can upload PDF documents, and the system will process and store them. Then, users can engage in a chat-like interface to ask questions related to the content of the uploaded documents. The backend is powered by Google's Gemini large language model to provide insightful and accurate answers.

Key Features:

- **User Authentication:** Secure user registration and login.
- **Document Upload:** Upload PDF documents for analysis.
- **Chat Interface:** Ask questions about your documents and get intelligent answers.
- **Document Management:** View and manage your uploaded documents.

## Tech Stack

This project is a monorepo and uses the following technologies:

### Frontend

- [Next.js](https://nextjs.org/) - React framework for server-side rendering and static site generation.
- [TypeScript](https://www.typescriptlang.org/) - Typed superset of JavaScript.
- [Tailwind CSS](https://tailwindcss.com/) - A utility-first CSS framework.
- [Axios](https://axios-http.com/) - Promise-based HTTP client for the browser and Node.js.
- [Firebase](https://firebase.google.com/) - Used for frontend authentication.
- [Shadcn/ui](https://ui.shadcn.com/) - Re-usable components built using Radix UI and Tailwind CSS.

### Backend

- [Go](https://golang.org/) - A statically typed, compiled programming language.
- [PostgreSQL with pgvector](https://github.com/pgvector/pgvector) - A PostgreSQL extension for vector similarity search.
- [Docker](https://www.docker.com/) - For containerizing the application.
- [Google Cloud Storage](https://cloud.google.com/storage) - For storing uploaded documents.
- [Google Gemini](https://ai.google.dev/) - Large language model for question answering.

## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

- [Node.js](https://nodejs.org/en/) (v20 or later)
- [npm](https://www.npmjs.com/)
- [Go](https://golang.org/doc/install) (v1.24 or later)
- [Docker](https://www.docker.com/get-started)

### Installation

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/Dinesh-Gautam/omara-assignment.git
    cd your_repository
    ```

2.  **Install NPM packages for the frontend:**

    ```sh
    cd apps/frontend
    npm install
    cd ../..
    ```

3.  **Install NPM packages for the root:**

    ```sh
    npm install
    ```

4.  **Set up environment variables for the backend:**

    Navigate to the `apps/backend` directory and create a `.env` file by copying the example:

    ```sh
    cd apps/backend
    cp .env.example .env
    ```

    You will need to fill in the `.env` file with your own credentials. See the `.env.example` file for a list of all required variables.

5.  **Set up environment variables for the frontend:**

    Navigate to the `apps/frontend` directory and create a `.env.local` file by copying the example:

    ```sh
    cd apps/frontend
    cp .env.example .env.local
    ```

    You will need to fill in the `.env.local` file with your Firebase project's configuration.

### Firebase Setup

This project uses Firebase for authentication. To set it up for local development, you will need to do the following:

1.  **Create a Firebase project:** If you don't have one already, create a new project in the [Firebase Console](https://console.firebase.google.com/).

2.  **Enable Authentication Providers:** In your Firebase project, go to the "Authentication" section and enable the following sign-in providers:

    - Email/Password
    - Google
    - Anonymous

3.  **Get Firebase Configuration:** In your Firebase project settings, find your web app's configuration object. You will need to copy these values into the `apps/frontend/.env.local` file.

4.  **Generate a Service Account Key:**
    - In your Firebase project, go to "Project settings" > "Service accounts".
    - Click "Generate new private key". This will download a JSON file.
    - Rename this file to `serviceAccountKey.json`.
    - Place this file in the `apps/backend/secrets/` directory. **Important:** This directory is included in the `.gitignore` file to prevent the key from being committed to version control.

## Usage

To run the application, navigate to the root of the project and run the following command:

```sh
npm run dev
```

This will start both the frontend and backend services concurrently.

- The **frontend** will be available at `http://localhost:3000`.
- The **backend** will be available at `http://localhost:8080`.

The backend is run using Docker Compose, which will also start a PostgreSQL database container. The `--watch` flag enables hot-reloading for the backend.

## Project Structure

The project is organized as a monorepo with two main applications:

- `apps/backend`: The Go backend application.
- `apps/frontend`: The Next.js frontend application.

Each application has its own set of dependencies and configuration files. The root `package.json` contains scripts for running both applications together.
