# Notification API Documentation

## Overview

The Notification API provides email and SMS notification capabilities using SMTP for emails and Africa's Talking for SMS messages. This API allows you to send individual or bulk notifications to customers and users.

## Configuration

### SMTP Configuration (Email)
```yaml
smtp:
  host: smtp.gmail.com          # SMTP server hostname
  port: 587                     # SMTP server port
  username: your-email@gmail.com # SMTP username
  password: your-app-password    # SMTP password/app password
  from: your-email@gmail.com     # Default sender email
  tls: true                     # Enable TLS encryption
```

### Africa's Talking Configuration (SMS)
```yaml
at:
  api_key: your_api_key         # Africa's Talking API key
  username: your_username       # Africa's Talking username
  base_url: https://api.africastalking.com # API base URL (optional)
```

## REST API Endpoints

All notification endpoints are available under `/api/notifications/` and require authentication.

### Send Email

**Endpoint:** `POST /api/notifications/email`

**Description:** Send an email notification to a single recipient.

**Request Body:**
```json
{
  "to": "recipient@example.com",
  "subject": "Email Subject",
  "body": "Plain text email body",
  "html_body": "<p>HTML email body (optional)</p>"
}
```

**Response:**
```json
{
  "message": "Email sent successfully",
  "to": "recipient@example.com"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/notifications/email \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "to": "customer@example.com",
    "subject": "Order Confirmation",
    "body": "Thank you for your order!",
    "html_body": "<h1>Thank you for your order!</h1><p>Your order has been confirmed.</p>"
  }'
```

### Send SMS

**Endpoint:** `POST /api/notifications/sms`

**Description:** Send an SMS notification to a single phone number.

**Request Body:**
```json
{
  "phone_number": "+1234567890",
  "message": "Your SMS message here"
}
```

**Response:**
```json
{
  "message": "SMS sent successfully",
  "phone_number": "+1234567890"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/notifications/sms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "phone_number": "+1234567890",
    "message": "Your order has been confirmed!"
  }'
```

### Send Bulk Email

**Endpoint:** `POST /api/notifications/bulk/email`

**Description:** Send the same email to multiple recipients.

**Request Body:**
```json
{
  "recipients": [
    "customer1@example.com",
    "customer2@example.com",
    "customer3@example.com"
  ],
  "subject": "Newsletter Subject",
  "body": "Newsletter content",
  "html_body": "<h1>Newsletter</h1><p>Content here</p>"
}
```

**Response:**
```json
{
  "message": "Bulk emails sent successfully",
  "recipients": [
    "customer1@example.com",
    "customer2@example.com",
    "customer3@example.com"
  ],
  "count": 3
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/notifications/bulk/email \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "recipients": ["customer1@example.com", "customer2@example.com"],
    "subject": "Weekly Newsletter",
    "body": "Check out our latest products!",
    "html_body": "<h2>Weekly Newsletter</h2><p>Check out our latest products!</p>"
  }'
```

### Send Bulk SMS

**Endpoint:** `POST /api/notifications/bulk/sms`

**Description:** Send the same SMS to multiple phone numbers.

**Request Body:**
```json
{
  "phone_numbers": [
    "+1234567890",
    "+0987654321",
    "+1122334455"
  ],
  "message": "Your bulk SMS message"
}
```

**Response:**
```json
{
  "message": "Bulk SMS sent successfully",
  "phone_numbers": [
    "+1234567890",
    "+0987654321",
    "+1122334455"
  ],
  "count": 3
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/notifications/bulk/sms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "phone_numbers": ["+1234567890", "+0987654321"],
    "message": "Flash sale! 50% off everything!"
  }'
```

### Send Generic Notification

**Endpoint:** `POST /api/notifications/send`

**Description:** Send a notification based on type (email or SMS).

**Request Body:**
```json
{
  "type": "email",  // or "sms"
  "to": "recipient@example.com",  // for email
  "phone_number": "+1234567890",  // for SMS
  "subject": "Email Subject",     // for email
  "body": "Email body",          // for email
  "html_body": "<p>HTML body</p>", // for email (optional)
  "message": "SMS message"       // for SMS
}
```

**Response:**
```json
{
  "message": "Notification sent successfully",
  "type": "email"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/notifications/send \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "type": "email",
    "to": "customer@example.com",
    "subject": "Welcome!",
    "body": "Welcome to our service!"
  }'
```

## GraphQL Integration

The notification service is also available in GraphQL resolvers and can be used to send notifications when orders are created, updated, or status changes occur.

### Order Confirmation Notifications

When an order is created, the system can automatically send:
- **Email confirmation** with order details, items, and total
- **SMS confirmation** with order number and total

### Order Status Updates

When order status changes, the system can send:
- **SMS notifications** for status updates (shipped, delivered, etc.)
- **Email notifications** for important status changes

## Phone Number Format

Phone numbers should be in international format:
- **Valid:** `+1234567890`, `+44123456789`, `+254712345678`
- **Invalid:** `1234567890`, `0712345678`

The system automatically validates phone numbers and returns appropriate errors for invalid formats.

## Email Validation

Email addresses are validated using standard email format validation:
- **Valid:** `user@example.com`, `test.email+tag@domain.co.uk`
- **Invalid:** `invalid-email`, `@domain.com`, `user@`

## Error Handling

### Common Error Responses

**400 Bad Request:**
```json
{
  "error": "Missing required fields: to, subject, body"
}
```

**401 Unauthorized:**
```json
{
  "error": "Missing Authorization header"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Failed to send email: SMTP authentication failed"
}
```

### SMTP Errors

- **Authentication failed:** Check username/password
- **Connection refused:** Check host/port configuration
- **TLS errors:** Verify TLS settings and certificates

### SMS Errors

- **Invalid phone number:** Check phone number format
- **API key invalid:** Verify Africa's Talking credentials
- **Insufficient credits:** Check Africa's Talking account balance

## Rate Limiting

- **Email:** No built-in rate limiting (depends on SMTP provider)
- **SMS:** Limited by Africa's Talking account limits and pricing

## Security Considerations

1. **SMTP Credentials:** Use app-specific passwords for Gmail/Google Workspace
2. **API Keys:** Store Africa's Talking API keys securely
3. **Authentication:** All endpoints require valid JWT/OIDC tokens
4. **Input Validation:** All inputs are validated before sending

## Testing

### Test Email
```bash
curl -X POST http://localhost:8080/api/notifications/email \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "to": "test@example.com",
    "subject": "Test Email",
    "body": "This is a test email from the notification API"
  }'
```

### Test SMS
```bash
curl -X POST http://localhost:8080/api/notifications/sms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "phone_number": "+1234567890",
    "message": "Test SMS from notification API"
  }'
```

## Integration Examples

### Order Confirmation Flow
```javascript
// When an order is created
const order = await createOrder(orderData);

// Send email confirmation
await fetch('/api/notifications/email', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    to: order.customer.email,
    subject: `Order Confirmation - ${order.orderNumber}`,
    body: `Thank you for your order! Order #${order.orderNumber}`,
    html_body: generateOrderConfirmationHTML(order)
  })
});

// Send SMS confirmation
await fetch('/api/notifications/sms', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    phone_number: order.customer.phone,
    message: `Order #${order.orderNumber} confirmed! Total: $${order.totalAmount}`
  })
});
```

### Newsletter Campaign
```javascript
// Send newsletter to all customers
const customers = await getAllCustomers();
const recipients = customers.map(c => c.email);

await fetch('/api/notifications/bulk/email', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    recipients: recipients,
    subject: 'Weekly Newsletter - New Products!',
    body: 'Check out our latest products and special offers!',
    html_body: generateNewsletterHTML()
  })
});
```

## Monitoring and Logging

The notification service logs:
- **Success:** Email/SMS sent successfully
- **Errors:** Failed deliveries with error details
- **Validation:** Invalid email addresses or phone numbers
- **Authentication:** SMTP and API authentication issues

Check server logs for notification delivery status and any errors that occur during the sending process.
