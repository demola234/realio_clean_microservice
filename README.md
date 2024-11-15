# Realio- Go MicroService Application

Realio is a Micro-service architecture built with Go and Python, while adhering to clean architecture principles, provides a scalable and modular way to handle core functionalities like property listings, user management, search, and recommendations. This architecture allows users to view, filter, and get personalised recommendations for properties they may wish to buy. Below is an in-depth look at how to structure and implement these services:

### 1. **Microservices Overview**

This application can be broken down into multiple Micro-services, each handling a specific function. Here’s a high-level overview of key services:

### Core Microservices

1. **User Service** (Go): Manages user accounts, authentication, and profiles.
2. **Property Service** (Go): Manages property listings, including CRUD (Create, Read, Update, Delete) operations.
3. **Search Service** (Go): Provides search functionality with filters for property type, location, price range, etc.
4. **Recommendation Service** (Go): This service uses machine learning to provide property recommendations based on user preferences and past interactions.
5. **Booking and Scheduling Service** (Go): Manages property viewing appointments and scheduling.
6. **Messaging Service** (Go): Facilitates communication between buyers, sellers, and agents.
7. **Notification Service** (Go): This service sends alerts and updates to users via email, SMS, and in-app notifications.

Each service is designed to operate independently, allowing for modular deployment and scaling based on user demand.

### **Detailed Service Architecture**

### **1. User Service (Go):**

- **Responsibilities**: Manages user registration, authentication, and user profile data.
- **Database**: PostgreSQL for structured user data.
- **Endpoints**:
  - `POST /register`: Register a new user.
  - `POST /verify_user`: Verify new users OTP.
  - `POST /resend_otp`: Resend User OTP.
  - `POST /check_user_email`: Check if Users’ Emails Already Exist.
  - `POST /logout`: Logout User.
  - `POST /login`: Authenticate a user.
  - `GET /profile/{user_id}`: Retrieve user profile details.
  - `GET /refresh_token`: Refreshes access token when user token is expired.
- **Storage:** [Cloudinary](https://www.google.com/search?sca_esv=f9749d82eb8de094&sxsrf=ADLYWIIPiXLw3w7mhzepFUYAC8wlZkB_Ug:1730135298258&q=Cloudinary&spell=1&sa=X&ved=2ahUKEwjSx_aeyLGJAxVC1DgGHRqxGk8QBSgAegQICBAB) or [Amazon S3](https://aws.amazon.com/pm/serv-s3/?gclid=Cj0KCQjw7Py4BhCbARIsAMMx-_JuG1730vIV3IVqAy-un_ZoBJZmZvdVhKw6eInTkro2UJhhPsLHPDQaAsbEEALw_wcB&trk=c8974be7-bc21-436d-8108-722e8ab912e1&sc_channel=ps&ef_id=Cj0KCQjw7Py4BhCbARIsAMMx-_JuG1730vIV3IVqAy-un_ZoBJZmZvdVhKw6eInTkro2UJhhPsLHPDQaAsbEEALw_wcB:G:s&s_kwcid=AL!4422!3!645125274431!e!!g!!amazon%20s3!19574556914!145779857032) for User Image Storage.
- Implement **Role-Based Access Control (RBAC)** for users, agents, and admins.

### **2. Property Listings (Go):**

- **Responsibilities**: Manages property listings, CRUD operations, and image storage.
  - **Database**: PostgreSQL for structured user data and Kafka as a Messaging Broker.
- **Endpoints**:
  - `GET /properties`: Retrieve a list of properties.
  - `POST /properties`: Add a new property.
  - `PUT /properties/{property_id}`: Edit property Added.
  - `DELETE /properties/{property_id}`: Delete property Added.
  - `GET /properties/{property_id}`: View property details.

### **3. Search Service (Go)**

- **Responsibilities**: Provides search and filtering functionality.
- **Database**: Elasticsearch for search indexing.
- **Event Management**: Data Sync at Intervals, When properties are created, updated, or deleted in the property management service, it can publish events containing relevant property details. These events can be sent through a message broker (Kafka) to which the search service subscribes. The search service then indexes the properties in Elasticsearch based on these events.
- **Endpoints**:
  - `GET /search`: Retrieve properties based on search criteria.

### **4. Booking and Scheduling Service (Go)**

- **Responsibilities**: Manages property view scheduling and booking confirmations.
- **Database**: PostgreSQL to handle structured booking data.
- **Endpoints**:
  - `POST /bookings`: Create a booking.
  - `PUT /bookings/{booking_id}`: Modify or reschedule a booking.
  - `DELETE /bookings/{booking_id}`: Cancel a booking.
  - `GET /bookings`: Retrieve a list of bookings.
  - `GET /bookings/{booking_id}`: Edit property details.

### **5. Messaging Service (Go)**

- **Responsibilities**:
  - Manages conversations and messages between users.
  - Supports both real-time communication via **Socket.IO** and persistent message storage in **MongoDB**.
  - Exposes RESTful API endpoints for sending, retrieving, and deleting messages in conversations
- **Database**: MongoDB for storing messages and conversations.
- Realtime Connection: SocketIO for real-time
- **Endpoints**:
  - `POST /messages`: Send a message.
  - `GET /messages/{conversation_id}`: Retrieve messages in a conversation.
  - `DELETE /messages/{conversation_id}`: Delete messages in a conversation.
- **Socket IO:**
  - **`join_room`**: Users join a conversation room.
  - **`send_message`**: Sends a message to the room.
  - **`receive_message`**: Listens for new messages in the room.
- **Database Design:** Using **MongoDB** is suitable here due to its flexibility with nested data, allowing efficient storage and retrieval of messages and conversations.

### **6. Notification Service (Go)**

- **Responsibilities**: Sends notifications for bookings, messages, or updates.
- **Database**: Redis for temporary storage and PostgreSQL for persistent notification storage.
- **Endpoints**:
  - `POST /notifications`: Send a notification.
  - `GET /notifications/{user_id}`: Fetch a list of user notifications.
  - `PUT /notification/{notification_id}`: Mark Notification as Read.
  - `PUT /notifications`: Mark All Notifications as Read.

### **7. Recommendation Service (Python)**

- **Responsibilities**: The recommendation service remains in Python, leveraging machine learning libraries like **Scikit-Learn and Reinforcement Learning** for personalised recommendations. This service connects seamlessly with other Go-based services, serving as an independent Micro-service for generating property suggestions based on user history and interactions.
- **Database**: Redis for temporary storage and PostgreSQL for persistent notification storage.
- ## **Endpoints**:

### **Scalability Considerations**

- **Load Balancer**: Use **NGINX** or **AWS ELB** to balance the load between instances. \*\*Optional
- **Database Scaling**: Use **Read Replicas** for scaling reads. Sharding for large datasets. \*\*Optional
- **Caching**: Use **Redis** to cache frequently accessed data like property listings.
- **Search**: Use **Elasticsearch** for handling complex search queries.
- **Microservices**: Break down large components into microservices (e.g., separate services for User, Property, Booking).

### **Security Considerations**

- **Data Encryption**: Use HTTPS for secure communication. Encrypt sensitive data like passwords.
- **Access Control**: Implement proper authorization for different roles (buyers, sellers, agents, admins).
- **Input Validation**: Sanitize inputs to prevent SQL injection and other common vulnerabilities.
- **Rate Limiting**: Prevent abuse by implementing rate limits on critical endpoints.

### **Monitoring & Maintenance**

- Use **Prometheus** and **Grafana** for monitoring and alerting.
- Use **ELK Stack (Elasticsearch, Logstash)** for logging and visualization.
- Automate deployments with **CI/CD** tools like GitHub Actions.
- Regular backups for databases using tools like **AWS Backup** or scheduled cron jobs.

### **Third-Party Integrations**

- **Crash Analytics**: Sentry
- **Geolocation & Maps**: Google Maps API.
- **Analytics**: Google Analytics.
- **Notifications**: Firebase for push notifications.

**DATABASE TABLE:**

### **User Table:**

| Field        | Type      | Description                 |
| ------------ | --------- | --------------------------- |
| `id`         | UUID      | Primary key                 |
| `name`       | String    | User's full name            |
| `email`      | String    | User's email (unique)       |
| `password`   | String    | Hashed password             |
| `role`       | Enum      | buyer, seller, agent, admin |
| `phone`      | String    | Contact number              |
| `created_at` | Timestamp | Timestamp of registration   |
| `updated_at` | Timestamp | Timestamp of last update    |

### User Session Table Schema

| Column          | Data Type       | Description                                                    |
| --------------- | --------------- | -------------------------------------------------------------- |
| `session_id`    | UUID            | Unique identifier for each session.                            |
| `user_id`       | UUID            | Foreign key linking to the user table (identifies the user).   |
| `token`         | VARCHAR (255)   | The session token, which can be a JWT or another token format. |
| `created_at`    | TIMESTAMP       | Timestamp of when the session was created.                     |
| `expires_at`    | TIMESTAMP       | Timestamp of when the session expires.                         |
| `last_activity` | TIMESTAMP       | Tracks the last activity time for session timeout checks.      |
| `ip_address`    | VARCHAR (45)    | The IP address from which the session was initiated.           |
| `user_agent`    | VARCHAR (255)   | The user agent (browser or device info) for the session.       |
| `is_active`     | BOOLEAN         | Indicates whether the session is currently active.             |
| `revoked_at`    | TIMESTAMP       | Timestamp for when the session was revoked, if applicable.     |
| `device_info`   | JSON or VARCHAR | (Optional) Stores additional device details if needed.         |

### **Property Table:**

| Field              | Type      | Description                            |
| ------------------ | --------- | -------------------------------------- |
| `id`               | UUID      | Primary key                            |
| `title`            | String    | Property title                         |
| `description`      | Text      | Detailed description                   |
| `price`            | Decimal   | Price of the property                  |
| `type`             | Enum      | House, Apartment, Land, etc.(Optional) |
| `address`          | String    | Address details                        |
| `zip_code`         | String    | Zip code(Optional)                     |
| `owner_id`         | UUID      | Reference to the user (seller)         |
| `images`           | Array     | List of image URLs(Optional)           |
| `no_of_bed_rooms`  | Int       | Number of Bed Rooms                    |
| `no_of_bath_rooms` | Int       | Number of Bath Rooms                   |
| `no_of_toilets`    | Int       | Number of Toilets                      |
| `geo_location`     | JSON      | Latitude & longitude(Optional)         |
| `status`           | Enum      | Available, Sold, Rented, etc.          |
| `created_at`       | Timestamp | Timestamp of listing                   |
| `updatedAt`        | Timestamp | Timestamp of last update               |

### **Booking Table:**

| Field          | Type      | Description                     |
| -------------- | --------- | ------------------------------- |
| `id`           | UUID      | Primary key                     |
| `property_id`  | UUID      | Reference to Property table     |
| `user_id`      | UUID      | Reference to User table (buyer) |
| `booking_date` | Date      | Date of the booking             |
| `status`       | Enum      | Pending, Confirmed, Cancelled   |
| `created_at`   | Timestamp | Timestamp of booking            |
| `updated_at`   | Timestamp | Timestamp of last update        |

### **Message Table:**

| Field             | Type      | Description                              |
| ----------------- | --------- | ---------------------------------------- |
| `id`              | UUID      | Primary key                              |
| `conversation_id` | UUID      | Reference to Conversation table          |
| `sender_id`       | UUID      | Reference to Send on User table          |
| `content`         | String    | Message Content                          |
| `read`            | Bool      | Confirm if message has been read by User |
| `created_at`      | Timestamp | Timestamp of booking                     |
| `updated_at`      | Timestamp | Timestamp of last update                 |
| `receiver_id`     | UUID      | Reference to Receiver on User table      |

### **Conversation Table:**

| Field             | Type      | Description                              |
| ----------------- | --------- | ---------------------------------------- |
| `id`              | UUID      | Primary key                              |
| `sender_id`       | UUID      | Reference to Send on User table          |
| `last_conversion` | String    | Message Content                          |
| `read`            | Bool      | Confirm if message has been read by User |
| `created_at`      | Timestamp | Timestamp of booking                     |
| `updated_at`      | Timestamp | Timestamp of last update                 |
| `receiver_id`     | UUID      | Reference to Receiver on User table      |
