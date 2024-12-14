Here’s a detailed README for the **email-signature-backend** project:

---

# Email Signature Backend

The **Email Signature Backend** is a RESTful API built with **Go** and **Fiber** that allows users to create, manage, and export professional email signatures. It supports user authentication, customizable templates, and click analytics for tracking link interactions.

## **Features**
- **User Authentication**: Secure user registration and login with JWT-based authentication.
- **Email Signature Management**: Create, retrieve, and preview email signatures.
- **Customizable Templates**: Support for multiple signature templates (e.g., Basic, Modern).
- **Export as HTML**: Export signatures as HTML for use in email clients.
- **Analytics**: Track and analyze clicks on links within the email signature.
- **Swagger Documentation**: Interactive API documentation.

---

## **Tech Stack**
- **Go**: Backend programming language.
- **Fiber**: High-performance web framework for Go.
- **PostgreSQL**: Database for storing user data, signatures, and analytics.
- **Swagger**: API documentation and testing.
- **Docker**: Containerization for deployment.

---

## **Getting Started**

### **Prerequisites**
Ensure you have the following installed:
- [Go](https://golang.org/dl/) (version 1.20+)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)
- [Swagger CLI](https://github.com/swaggo/swag)

---

### **Installation**

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-username/email-signature-backend.git
   cd email-signature-backend
   ```

2. **Set Up Environment Variables**:
   Create a `.env` file in the root directory with the following:
   ```env
   DATABASE_URL=postgres://username:password@localhost:5432/email_signature?sslmode=disable
   JWT_SECRET=your_jwt_secret_key
   ```

3. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

4. **Run Database Migrations**:
   Use the `migrate` tool to set up the database schema:
   ```bash
   migrate -path db/migrations -database ${DATABASE_URL} up
   ```

5. **Run the Application**:
   ```bash
   go run main.go
   ```

6. **Access the API**:
   - API Base URL: `http://localhost:3000/api`
   - Swagger Documentation: `http://localhost:3000/swagger/index.html`

---

### **Endpoints Overview**

#### **Authentication**
- **POST** `/api/register`: Register a new user.
- **POST** `/api/login`: Authenticate a user and generate a JWT.

#### **Signatures**
- **POST** `/api/signature`: Create a new email signature.
- **GET** `/api/signature/{id}/export`: Export a signature as HTML.
- **GET** `/api/signature/{id}/preview`: Preview a signature in the browser.

#### **Links**
- **POST** `/api/links`: Create a new link for a signature.

#### **Analytics**
- **POST** `/api/track`: Track a click on a link.
- **GET** `/api/analytics`: Retrieve click analytics for a user’s links.

---

## **Testing**

1. **Using Swagger**:
   Visit `http://localhost:3000/swagger/index.html` to interact with the API.

2. **cURL Examples**:
   **Create a Signature**:
   ```bash
   curl -X POST http://localhost:3000/api/signature \
   -H "Authorization: Bearer YOUR_JWT_TOKEN" \
   -H "Content-Type: application/json" \
   -d '{
         "template_data": {
             "name": "John Doe",
             "job_title": "Software Engineer",
             "company": "TechCorp",
             "phone": "+123456789",
             "website": "https://example.com",
             "social_links": {
                 "linkedin": "https://linkedin.com/in/johndoe",
                 "twitter": "https://twitter.com/johndoe"
             }
         }
     }'
   ```

---

## **Deployment**

### **Docker Deployment**

1. **Build Docker Image**:
   ```bash
   docker build -t email-signature-backend .
   ```

2. **Run the Container**:
   ```bash
   docker run -p 3000:3000 --env-file .env email-signature-backend
   ```

3. **Access the Application**:
   - API Base URL: `http://localhost:3000/api`
   - Swagger Documentation: `http://localhost:3000/swagger/index.html`

---

## **Project Structure**

```
email-signature-backend/
│
├── db/                   # Database migrations
│   ├── migrations/
│   └── schema.sql
├── docs/                 # Swagger documentation
├── handlers/             # API handler functions
├── middleware/           # Authentication and other middleware
├── routes/               # Route definitions
├── utils/                # Utility functions
├── main.go               # Application entry point
└── go.mod                # Go module dependencies
```

---

## **Contributing**

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Submit a pull request.

---

## **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
