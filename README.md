# Puna Badminton Court Reservation System

This is a web application for reserving badminton courts at the Puna Badminton Hall. It provides a user-friendly interface for searching for available courts, making reservations, and managing existing reservations.

## Features

*   **Search for available courts:** Users can search for available courts by date and time.
*   **Make reservations:** Users can make reservations for available courts.
*   **Manage reservations:** Users can view and cancel their existing reservations.
*   **User authentication:** Users can create an account and log in to manage their reservations.

## Technologies Used

*   **Go:** The back-end of the application is written in Go.
*   **Bulma:** The front-end of the application is styled with Bulma.
*   **PostgreSQL:** The application uses a PostgreSQL database to store data.
*   **Rspack:** The front-end assets are bundled with Rspack.

## Getting Started

To get started with the application, you will need to have the following installed:

*   Go
*   PostgreSQL
*   Node.js
*   pnpm

Once you have the prerequisites installed, you can follow these steps to get the application up and running:

1.  Clone the repository:

    ```
    git clone https://github.com/your-username/puna.git
    ```

2.  Install the dependencies:

    ```
    pnpm install
    ```

3.  Create a PostgreSQL database and update the database connection string in `internal/config/config.go`.

4.  Run the database migrations:

    ```
    go run cmd/web/main.go -migrate
    ```

5.  Start the application:

    ```
    go run cmd/web/main.go
    ```

The application will now be running at `http://localhost:8080`.
