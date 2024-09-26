## Task List

### 1. Provision a Cloud Storage Bucket using Infrastructure as Code (IaC)
We used Terraform to deploy an S3 bucket in the `us-west-2` region.

**Steps:**
1. Navigate to the Terraform directory:
    ```bash
    cd terraform/s3
    ```
2. Initialize Terraform:
    ```bash
    terraform init
    ```
3. Plan the deployment:
    ```bash
    terraform plan
    ```
4. Apply the deployment:
    ```bash
    terraform apply
    ```

### 2. Make an Endpoint `/update_airport_image` to Update an Airport’s Image
We used the AWS SDK to connect with AWS and utilized a pre-signed URL to upload images to the bucket.

### 3. Containerize the Go Application
To build and run the containerized application, use the following commands:

**Build the Docker image:**
```bash
docker build -t bd-airports:latest .
```

**Run the Docker container:**
```bash
docker run -p 8080:8080 bd-airports:latest
```

### 4. Prepare a Deployment and Service Resource to Deploy in Kubernetes
To deploy the application in Kubernetes, apply the manifest file:
```bash
kubectl apply -f manifest.yml
```

## Getting Started

### Prerequisites
- Docker
- Kubernetes
- Terraform
- AWS CLI

---

Feel free to let me know if you need any more adjustments!
