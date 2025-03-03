# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

-   Request ID middleware that adds a unique identifier to each request
    -   Generates a UUID for each request if not provided in the `X-Request-ID` header
    -   Includes the request ID in all response headers
    -   Adds the request ID to all log entries for request tracing
    -   Implements best practices for distributed tracing
-   BaseHandler for all API handlers to provide common functionality
    -   Provides access to request-specific logger with request ID
    -   Simplifies access to request context data
-   Enhanced validation error messages for better user experience
    -   Custom error messages for required fields, email validation, and length constraints
    -   Clear field names in error messages based on JSON tags

### Changed

-   Updated Logger middleware to include request ID in all log entries
-   Modified AuthHandler and TodoHandler to use request-specific logger
-   Improved error handling with consistent request ID tracking
-   Enhanced validator implementation with custom error messages for common validation rules
-   Replaced "username" field with "fullname" throughout the application
    -   Updated database schema to remove username column and add fullname column
    -   Updated user model to use fullname instead of username
    -   Updated authentication flow to use email as the unique identifier
    -   Updated JWT tokens to include fullname instead of username
    -   Updated API endpoints to use fullname in requests and responses

### Security

-   Enhanced request tracing capabilities for better debugging and audit trails
-   Improved correlation between logs and requests for security incident investigation

### Fixed

-   Fixed validation in Register and Login handlers to properly validate all fields
-   Fixed DeleteTodo handler to use the correct field name (ID instead of TodoID)

### Migration Notes

-   Existing usernames have been migrated to the fullname field
-   Email is now the only unique identifier for users
-   Login now requires email instead of username or email
