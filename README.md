# Chat Application
Real time chat application developed by the requirements listed in the follwing [document](https://docs.google.com/document/d/1Kyotim0S9Ef-qnKuZaKJKtZyWa_GkHcHOo1s2wieWJM/edit#heading=h.er3tlss7yggl)

## Technologies Used

- **Backend**: Go
- **Database**: SQL (PostgreSQL)
- **Frontend**: ReactJS
- **Containerization**: Docker
- **Orchestration**: Docker Compose

## Architecture

The application consists of several services, each running in its own Docker container:

- `frontend`: The frontend service, built with React and served using NGINX.
- `backend`: The backend service, built with Go.
- `postgres`: The PostgreSQL database service.
- `rabbitmq`: The RabbitMQ service, used for real-time messaging.

## How to Run

1. Ensure Docker and Docker Compose are installed on your machine.
    - Visit the Docker website and download the appropriate Docker Desktop installer
       - For Windows: [Download Docker Desktop for Windows](https://hub.docker.com/editions/community/docker-ce-desktop-windows/)
       - For Mac: [Download Docker Desktop for Mac](https://hub.docker.com/editions/community/docker-ce-desktop-mac/)
    - After installation, you can verify that Docker is installed correctly by opening a terminal and running the following command:
      ```bash
      docker --version
      ```

2. Clone the repository
    ```bash
    git clone https://github.com/johnnyCake1/chat-app
    ```
    and navigate to the project directory
    ```bash
    cd chat-app
    ```

3. Run the following command to build and start the application:
    ```bash
    docker-compose up --build
    ```

4. Once the application is running, you can access the frontend at `http://localhost`. For the development environment, the backend is accessible at `http://localhost:8080`. All the requests to http://localhost will be proxied to the backend service by NGINX running in the frontend container.
