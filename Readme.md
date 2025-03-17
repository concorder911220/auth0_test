## Setup

1. Clone the repository:
    ```sh
    git clone https://github.com/concorder911220/auth0_test.git
    cd auth0_test
    ```

2. Create a `.env` file in the root directory with the following content:
    ```env
    GOOGLE_CLIENT_ID=your-google-client-id
    GOOGLE_CLIENT_SECRET=your-google-client-secret
    DATABASE_URL="postgres://username:password@localhost:5432/user_service?sslmode=disable"
    JWT_SECRET_KEY="your-jwt-secret-key"
    REDIRECT_URL="http://localhost:8080/auth/callback"
    ```

3. Install dependencies:
    ```sh
    go mod tidy
    ```

4. Run Project:
    ```sh
    go run cmd/main.go
    ```

The application will start on [http://localhost:8080](http://localhost:8080).

## Running with Docker

1. Build and run the application using Docker:
    ```sh
    docker-compose build
    docker-compose up
    ```

The application will start on [http://localhost:8080](http://localhost:8080).

## API Endpoints

### Public Endpoints
- **GET** `/auth/login` - Initiates the OAuth login process.
- **GET** `/auth/callback` - Handles the OAuth callback.

### Protected Endpoints
- **POST** `/users` - Creates a new user.
- **GET** `/users/:id` - Retrieves a user by ID.
- **PUT** `/users/:id` - Updates a user by ID.
- **DELETE** `/users/:id` - Deletes a user by ID.

## Postman Collection
A Postman collection is provided in the `postman_collection.json` file. You can import this collection into Postman to test the API endpoints.