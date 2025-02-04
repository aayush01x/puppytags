# Usage

1) Make an `.env` file of format `envformat.txt`. Then, make a `psql` database and update its database name, host, port etc in `.env` file.

2) 
    ```sh
    go run main.go
    ```


3) To update tags for a user(user10 here):

    ```sh
    curl -X POST -H "Content-Type: application/json" \
            -d '{"user_id": "user10", "tags": [2, 5, 10, 25]}' \
            http://localhost:9090/user
    ```

4) To query for best matches for user1 with pagination: 
    ```sh
    curl "http://localhost:9090/similar?user_id=user1&offset=0&limit=3"
    ```
    - `offset`: starting index 
    - `limit` : returns next top `limit` number of users starting from `offset` index.