
<h1 align="center">social-network</h1>
<p align="center">
  <img src="https://img.shields.io/badge/-Svelte-black?style=flat-square&logo=svelte" alt="Svelte Badge">
  <img src="https://img.shields.io/badge/-Golang-black?style=flat-square&logo=go" alt="Golang Badge">
  <img src="https://img.shields.io/badge/-Docker-black?style=flat-square&logo=docker" alt="Docker Badge">
  <img src="https://img.shields.io/badge/-SQLite-black?style=flat-square&logo=sqlite" alt="SQLite Badge">
</p>


# Quick Start

After cloning the repository, ensure the `run.sh` script has executable permissions:

```bash
chmod +x run.sh
```

Run

```bash
./start.sh
```

Open [http://localhost:5000](http://localhost:5000)



### Default Accounts

Upon starting the services, three default accounts will be created for you:

    admin:secret
    admin2:secret
    admin3:secret

These accounts can be used to explore the application's features without the need for initial setup.

## Docker

To get started, simply run
```bash
docker-compose up
```

This command will:

1.  Build two Docker images:
    
    -   **frontend**: A lightweight image for the frontend application.
    -   **backend**: A lightweight image for the backend server.

2.  After building, the services will start and be accessible at:
    
    -   **frontend**: [localhost:5000](http://localhost:5000)
    -   **backend**: [localhost:80](http://localhost:80)




## Running Locally


### Frontend

1.  Navigate to the frontend directory:
```bash
cd frontend
```
2.  Install the necessary npm packages:
```bash
npm install
```

3.  Start the frontend development server:

```bash
npm run dev
```

### Backend

1.  Navigate to the backend directory:

```bash
cd backend
```

2.  Start the backend server:

```bash
go run .
```

With both services running, you can access the frontend at [localhost:5000](http://localhost:5000) and the backend at [localhost:80](http://localhost:80).
