# Prototype: Uploads to S3 with Pre-Signed URL

## Key Features

*   **Secure Uploads:** Generates pre-signed URLs, allowing clients to upload files directly to S3 without requiring AWS credentials.
*   **Confirmation Endpoint:** Provides an endpoint to confirm successful uploads by verifying the file's existence in S3.

## Run locally
1.  **Clone the repository:**
2.  **Set environment variables:**

    Create a `.env` file in the project root with the required environment variables:

    ```
    S3_BUCKET_NAME=your-bucket-name
    AWS_REGION=your-aws-region
    AWS_ACCESS_KEY_ID=your-access-key-id
    AWS_SECRET_ACCESS_KEY=your-secret-access-key
    CDN=your-cdn-url
    ```

3.  **Run the application:**

    ```bash
    go run main.go
    ```

    The server will start listening on port 8080.

## Usage

1.  **Prepare Upload:**

    Send a POST request to the `/prepare-upload` endpoint with the desired file extension in the request body:

    ```json
    {
      "file_extension": ".jpeg"
    }
    ```

    The server will return a JSON response with the pre-signed URL and the object key:

    ```json
    {
      "presigned_url": "https://your-bucket-name.s3.your-aws-region.amazonaws.com/uploads/your-uuid.jpeg?...",
      "key": "uploads/your-uuid.jpeg"
    }
    ```

2.  **Upload File:**

    Use the pre-signed URL to upload the file directly to S3 using a PUT request.  You'll need to set the `Content-Type` header to `image/jpeg`.

    ```bash
    curl -v -X PUT -H "Content-Type: image/jpeg" --data-binary "@/path/to/your/image.jpeg" "https://your-bucket-name.s3.your-aws-region.amazonaws.com/uploads/your-uuid.jpeg?..."
    ```

3.  **Confirm Upload:**

    Send a POST request to the `/upload-confirm` endpoint with the object key in the request body:

    ```json
    {
      "key": "uploads/your-uuid.jpeg"
    }
    ```

    The server will save the image key to database and return a JSON response indicating whether the upload was successful:

    ```json
    {
      "status": "file uploaded successfully",
      "success": true,
      "url": "https://your-cdn-url/uploads/your-uuid.jpeg"
    }
    ```

4.  **Get all uploaded images:**

    Send a GET request to the `/get-uploaded-images`:
    The server will return a JSON response with all the images:

    ```json
    {
      "images": [
        "https://your-cdn-url/uploads/your-uuid.jpeg",
        "https://your-cdn-url/uploads/your-uuid.jpeg",
        "https://your-cdn-url/uploads/your-uuid.jpeg",
      ],
      "success": true
    }
    ```