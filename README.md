# Receipt Processor

This is a simple web service that processes receipts and calculates points based on a set of rules.

## Technologies Used

1. Go (Golang): Main programming language.
2. Gorilla Mux: HTTP router for Go used for handling routes.
3. UUID: Used to generate unique receipt IDs.
4. sync.Map: To store receipt data in memory (in a concurrent-safe way).

## Points Calculation Rules
The points awarded for a receipt are based on the following rules:

- 1 point for every alphanumeric character in the retailer name.
- 50 points if the total is a round dollar amount (e.g., 10.00).
- 25 points if the total is a multiple of 0.25.
- 5 points for every two items on the receipt.
- For any item where the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. Add that to the points.
- 6 points if the day in the purchase date is odd.
- 10 points if the purchase time is between 2:00 PM and 4:00 PM.


## Edge Cases Handled
- Invalid Receipt Format: If the receipt JSON is malformed, the server returns a 400 Bad Request.
- No Items in Receipt: If the receipt contains no items, a 400 Bad Request is returned.
- Invalid Price or Total: If the price or total values cannot be parsed as floats or are less than or equal to 0, no points are awarded for these specific values.
- Empty Item Descriptions: If any item has an empty description, it is skipped.
- Invalid Date or Time Formats: The service returns 400 Bad Request if the date or time is not in valid formats.

## Install Dependencies

The following dependencies are required for the project:

- Gorilla Mux for routing.
- Google UUID for generating unique receipt IDs.

Install these dependencies using go get:

```bash 
go get -u github.com/google/uuid
go get -u github.com/gorilla/mux
```



## Running Locally

To run this service locally using Docker, follow these steps:

1. Build the Docker image:
    ```bash
    docker build -t receipt-processor .
    ```

2. Run the Docker container:
    ```bash
    docker run -p 8080:8080 receipt-processor
    ```

The application will be available at `http://localhost:8080`.

## Endpoints

- `POST /receipts/process`: Submits a receipt for processing.
    - Request body should be a JSON object representing the receipt.
    - Returns a JSON object with a unique `id`.

- `GET /receipts/{id}/points`: Retrieves the points awarded for the receipt with the given `id`.

## Testing

You can test the application with tools like Postman or curl.

### Example curl commands :

1. Process a Receipt:

```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    { "shortDescription": "Mountain Dew 12PK", "price": "6.49" },
    { "shortDescription": "Emils Cheese Pizza", "price": "12.25" }
  ],
  "total": "18.74"
}' http://localhost:8080/receipts/process
```

2. Get Points for a Receipt (replace {id} with the actual receipt ID from the previous response):

```bash
curl http://localhost:8080/receipts/{id}/points
```