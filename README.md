# Farmers Backend

A RESTful API backend for farm management built with Go.

## Features

- User authentication with JWT
- Crop and land management
- Season tracking with start and end dates
- Input tracking (seeds, fertilizer, water, labor, transport, etc.)
- Activity logging
- Cost analysis and reporting by season, input type, and annual summaries

## Technology Stack

- Go 1.25.2
- Gin web framework
- GORM with SQLite database
- JWT for authentication
- Viper for configuration management

## Project Structure

- `internal/config` - Configuration management
- `internal/db` - Database connection and migrations
- `internal/middleware` - Authentication and validation middleware
- `internal/models` - Data models (User, Crop, Land, Season, Input, Activity)
- `internal/services` - Business logic services

