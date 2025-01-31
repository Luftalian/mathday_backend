-- +goose Up
CREATE TABLE IF NOT EXISTS events (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    organizer VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_date DATE NOT NULL,
    end_time TIME NOT NULL,
    email VARCHAR(255) NOT NULL,
    prefecture VARCHAR(255),
    event_type VARCHAR(255),
    is_online BOOLEAN DEFAULT FALSE,
    is_offline BOOLEAN DEFAULT FALSE,
    official_url VARCHAR(255),
    online_lecture_url VARCHAR(255),
    venue VARCHAR(255),
    target VARCHAR(255),
    capacity VARCHAR(255),
    description TEXT,
    tags JSON,
    speakers JSON,
    schedule JSON,
    auth_code VARCHAR(36),
    is_authenticated BOOLEAN DEFAULT FALSE
);
