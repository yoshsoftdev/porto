# Portfolio — Backend & Integration Engineer

## Profile

**Backend & Integration Engineer** dengan pengalaman membangun backend system, third-party API integration, automation, hardware integration, dan AI deployment.

Berpengalaman menggunakan **Golang, Python, SQL, PostgreSQL, MySQL, Redis, NATS, WebSocket, Docker, Linux**, serta berbagai API dan protokol seperti **Meta WhatsApp Cloud API, Hikvision API, Telegram Bot API, Modbus TCP, dan SNMP**.

Saya terbiasa mengerjakan sistem dari sisi backend hingga integrasi dengan perangkat dan layanan eksternal. Dalam pengembangan software, saya menggunakan AI sebagai development assistant untuk membantu implementasi, debugging, dokumentasi, dan eksplorasi solusi, sementara architecture, system flow, requirements, testing, dan keputusan teknis tetap saya tentukan dan validasi sendiri.

---

# Core Strengths

- Backend Development
- Third-party API Integration
- System Integration
- REST API
- WebSocket / Realtime System
- Message Broker
- Hardware & Protocol Integration
- Database Design
- Automation
- Edge Computing
- AI / Computer Vision Deployment
- Linux & Docker
- Full-stack Development

---

# Technical Skills

## Programming Languages

- Golang
- Python
- JavaScript
- SQL

## Backend

- Gin
- REST API
- WebSocket
- Tortoise ORM
- JWT Authentication

## Frontend

- React
- Next.js
- Vite

## Database

- PostgreSQL
- MySQL
- SQLite

## Infrastructure

- Linux
- Docker
- Nginx
- systemd
- Linux VM
- Debian ARM64

## Messaging & Cache

- NATS
- Redis
- WebSocket

## API & Integration

- Meta WhatsApp Cloud API
- WhatsApp Business API / WABA
- Hikvision API
- Telegram Bot API
- Modbus TCP
- SNMP

## AI / Computer Vision

- YOLO
- YOLO-Pose
- ByteTrack
- FastReID
- People Counting
- Object Detection
- Person Re-Identification
- Image Annotation
- Edge AI Deployment

---

# Projects

## 1. Omni-chan — WhatsApp Omnichannel

**Period:** 2026  
**Client:** PT Digital Otoma Solusi  
**Status:** Prototype  
**Role:** Full-stack Engineer  
**Team:** 1  
**Lead:** Yes

### Overview

Omni-chan is a WhatsApp omnichannel system designed to centralize conversations from multiple WhatsApp numbers into a single Team Inbox.

Previously, team members handled customers from different WhatsApp numbers, causing customers to be transferred between numbers and making conversations difficult to manage.

The system provides a centralized interface for customer service teams to manage conversations.

### Architecture

```text
WhatsApp / WABA
       │
       ▼
Go Webhook Service
       │
       ▼
      NATS
       │
       ▼
   Go Backend
    │      │
    ▼      ▼
PostgreSQL Redis
       │
       ▼
   WebSocket
       │
       ▼
 Next.js Team Inbox
```

### Technology

- Golang
- Gin
- Next.js
- PostgreSQL
- Redis
- NATS
- WebSocket
- JWT
- Linux
- Meta WhatsApp Cloud API / WABA

### Main Features

- WhatsApp conversation management
- Team Inbox
- Realtime messaging
- Text messages
- Template messages
- Image messages
- Video messages
- PDF messages
- Interactive messages
- Broadcast
- Campaign management
- Message status
- WhatsApp account management

### My Contribution

I handled the project end-to-end, including:

- Backend architecture
- Frontend development
- Database design
- WhatsApp API integration
- Webhook implementation
- NATS integration
- Redis integration
- WebSocket realtime communication
- Campaign and broadcast flow
- Testing and debugging

### Why Webhook Was Separated

The webhook service was separated from the main backend because Meta requires a publicly accessible webhook endpoint.

Keeping it separate allows the webhook service to remain lightweight while the main backend handles business logic.

NATS is used to distribute webhook data so the main backend and other consumers do not need to perform heavy processing directly inside the webhook request.

### Redis Usage

Redis is used as a cache layer to reduce repeated reads from PostgreSQL.

Frequently requested data can be served from Redis instead of repeatedly querying the main database.

### Challenge

The main challenge was learning and integrating the Meta WhatsApp API, including WABA and Co-EX requirements.

### Solution

I used Meta documentation, implementation examples, testing, and AI-assisted discussions to understand the API behavior and integration requirements.

### Highlight

Built a working omnichannel architecture combining:

**External API + Webhook + Message Broker + Backend + Database + Cache + WebSocket + Frontend**

---

# 2. Immersion Controller

**Period:** 2025  
**Client:** Universitas Mulawarman  
**Status:** Production  
**Role:** Full-stack Engineer  
**Team:** 2  
**Lead:** Yes

### Overview

Immersion Controller is a server cooling control system designed to automate immersion cooling based on server temperature.

The system monitors temperature and controls an immersion pump according to configurable conditions.

### Architecture

```text
Server Temperature
       │
       ▼
    Golang
       │
       ▼
Threshold / Logic
       │
       ▼
   Modbus TCP
       │
       ▼
 Immersion Pump
```

The system also uses Redis for state/data and a Vite frontend for monitoring and configuration.

### Technology

- Golang
- Vite
- Redis
- SQLite
- Modbus TCP
- SNMP
- Debian ARM64

### Main Features

- Temperature monitoring
- Configurable threshold
- Pump ON/OFF control
- Modbus TCP communication
- Configurable retry
- Error logging
- Web-based monitoring
- Flexible control flow

### My Contribution

I handled the software side end-to-end:

- Backend
- Frontend
- Database
- Modbus integration
- Control logic
- Error handling
- Deployment

The hardware assembly was handled by another team member.

### Engineering Challenge

One of the main challenges was distinguishing between software problems and incorrect hardware/field installation.

Modbus communication could produce unexpected values when the hardware installation or field configuration was incorrect.

### Solution

I worked with the hardware/field team to verify the correct installation and communication flow before determining whether the issue originated from software or hardware.

### Architecture Evolution

Initially, Modbus commands were executed sequentially.

Later, the requirement changed to support more flexible conditional logic.

I redesigned the backend, database, and frontend so the workflow could be configured more dynamically, with an interface inspired by workflow tools such as n8n.

### Deployment

The application runs on a **NanoPi with Debian ARM64**.

---

# 3. Visitor Management System

**Period:** 2023–2024  
**Client:** Kemenkeu  
**Status:** Completed  
**Role:** Integration Engineer  
**Team:** 6  
**Lead:** No

### Overview

Visitor Management System was developed to control and audit visitor and employee access to restricted areas.

The system helps manage access to areas such as server rooms.

### Technology

- React
- Python
- Tortoise ORM
- MySQL
- Docker
- Linux
- Hikvision API

### My Contribution

My main responsibilities included:

- Testing Hikvision API
- Identifying suitable API endpoints
- Documenting API behavior
- Implementing visitor access request/ticketing functionality
- Integrating the application with Hikvision devices

### Hikvision Integration

The system communicates with more than **30 Hikvision devices**.

A synchronization process runs every day at **23:00 server time**.

The backend acts as the authoritative source.

The synchronization process:

```text
Backend Database
       │
       ▼
GET person data
from Hikvision devices
       │
       ▼
Compare with backend
       │
       ├── Same → No action
       │
       └── Different
              │
              ▼
          POST update
              │
              ▼
       Hikvision Device
```

If synchronization fails, the system retries a configured number of times. If it continues to fail, the error is recorded in the logs.

### Challenge

This was one of my first real-world software projects.

One of the challenges was learning Tortoise ORM while working on a real production-oriented system.

### Solution

I learned the framework through documentation, built a sandbox implementation, and then applied the knowledge to the actual system.

### Highlight

The project gave me practical experience integrating a web application with physical access-control devices through an external API.

---

# 4. Asset Management System

**Period:** 2023  
**Client:** Balai Kota DKI Jakarta  
**Status:** Completed / Production  
**Role:** Full-stack Engineer

### Overview

An asset management system for managing datacenter infrastructure and equipment.

The previous system was unreliable and made it difficult to track equipment placement, rack capacity, maintenance, and infrastructure relationships.

### Technology

- React
- Python
- Tortoise ORM
- MySQL
- Docker
- Linux

### Integration

The application integrates with the **Balai Kota DKI employee API** for employee authentication.

### Main Features

- Asset management
- Rack management
- PDU management
- Server placement
- Rack capacity tracking
- Maintenance ticketing
- Asset relationship management
- Employee login

### Asset Relationship

The system models infrastructure relationships such as:

```text
Room
 │
 └── Rack
      │
      └── Position
           │
           └── Server
```

### Rack Capacity Logic

When a rack is created, the system generates U-position records.

For example, if a server requires **2U**, the system searches for two adjacent available positions.

### My Contribution

- Database relationship design
- Backend development
- Frontend development
- Asset management logic
- Rack/PDU relationship
- Authentication integration

### Challenge

The most challenging part was designing the relationships between assets, racks, positions, and infrastructure components.

### Solution

I researched existing approaches, discussed the design with the team, and iterated on the database structure until the relationship model was practical for the application.

### Impact

The client migrated to a more stable and reliable system.

---

# 5. Visibel Edge — AI People Counting

**Period:** 2025–2026  
**Client:** Confidential  
**Status:** Production / Ongoing  
**Role:** Integration Engineer  
**Users:** Retail & Mall

### Overview

Visibel Edge is an AI-based people counting system designed to extract visitor information from existing CCTV infrastructure.

The goal is to make existing CCTV systems more useful by turning video streams into visitor data.

### Use Case

Retail stores and malls can use visitor data to understand visitor traffic and evaluate performance.

### My Contribution

My primary contribution is focused on:

- AI deployment
- Dataset annotation
- Annotation tooling
- Edge deployment
- System integration

### Computer Vision Pipeline

The development pipeline includes:

```text
CCTV
 │
 ▼
YOLO-Pose
 │
 ▼
ByteTrack
 │
 ▼
Person Tracklet
 │
 ▼
FastReID
 │
 ▼
Embedding
 │
 ▼
Cosine Similarity
 │
 ▼
Global Identity
```

### AI Components

#### YOLO-Pose

Used for person detection and pose/keypoint information.

#### ByteTrack

Used to maintain tracking identity between frames and create stable tracklets.

#### FastReID

Used to extract person embeddings for Re-Identification.

#### Gallery Matching

Person embeddings can be compared using cosine similarity to determine whether observations may belong to the same person.

### Annotation

I also worked on annotation and created an internal annotation tool.

Annotation includes:

- Bounding box
- Identity ID
- Sorting
- Occlusion handling
- Overlap
- Partial body
- Edge-of-frame cases

### Engineering Challenges

Real-world CCTV conditions introduce challenges such as:

- Different camera angles
- Front/back person views
- Different lighting conditions
- Occlusion
- Similar clothing/attributes
- Tracklet fragmentation
- ID switching
- People stopping for long periods

The system therefore requires more than simply running object detection on every frame.

### Deployment

The system is designed around existing CCTV infrastructure rather than requiring replacement of existing cameras.

---

# 6. Auto-bon — Petty Cash Management

**Period:** 2025–2026  
**Client:** Internal / Personal Office  
**Status:** Production  
**Team:** 1  
**Role:** Full-stack Engineer  
**Lead:** Yes

### Overview

Auto-bon is a petty cash and reimbursement management system built to replace a spreadsheet-based workflow.

The main goal was to allow reimbursement requests and petty cash operations to be performed from a phone through Telegram.

### Architecture

```text
Telegram
   │
   ▼
Telegram Bot
   │
   ▼
Golang Backend
   │
   ├── PostgreSQL
   │
   └── Receipt Storage
```

A Next.js internal web application is used for management and reporting.

### Technology

- Golang
- Next.js
- PostgreSQL
- Telegram Bot API
- Linux VM

### Telegram Commands

```text
/start
/help
/cancel
/rembes
/pay
```

### Reimbursement Flow

```text
Employee
   │
   ▼
Request Reimbursement
   │
   ▼
Admin Review
   │
   ▼
Payment
   │
   ▼
Proof / Receipt
   │
   ▼
Ledger
```

### Main Features

- Reimbursement request
- Payment recording
- Receipt upload
- Petty cash ledger
- Income/expense tracking
- Web management
- Excel export

### My Contribution

I built the system end-to-end:

- Telegram bot
- Backend
- Database
- Web application
- Reimbursement workflow
- Ledger
- Migration
- Deployment

### Challenge

The main challenge was migrating existing spreadsheet data into the new database.

### Solution

I created a simple migration program, tested it using a small dataset, and then performed the complete migration.

### Highlight

This project demonstrates my ability to identify a repetitive manual workflow and turn it into a small internal product.

---

# 7. Battery Monitoring — Proof of Concept

**Period:** 2025  
**Client:** PT International Chemical Industry  
**Status:** POC Completed — Not Continued to Production  
**Role:** Full-stack Engineer

### Overview

A proof-of-concept battery monitoring system developed for a resource-constrained edge device.

### Hardware

- NanoPi
- 1 GB RAM
- 8 GB storage

### Original Technology

- Python
- React
- MySQL
- Docker

### My Contribution

- Backend
- Frontend
- Database
- Deployment
- POC implementation

### Outcome

The POC was completed, but the client did not continue the project into production.

The exact reason for the project not continuing is not included in this portfolio data.

### Engineering Retrospective

If I designed the same system today, I would simplify the architecture to better match the hardware constraints.

Possible stack:

```text
Golang
   │
   ├── SQLite
   ├── Redis
   └── WebSocket

Vite
   │
   ▼
Frontend

systemd + Nginx
   │
   ▼
Deployment
```

### Lesson

The project reinforced an important engineering principle:

> Architecture and technology choices should match the constraints of the target hardware.

A technically valid stack can still be a poor choice when its operational overhead is too large for the available resources.

---

# Engineering Approach

## 1. Understand the System Before Coding

I prefer to understand:

- Requirements
- Data flow
- External dependencies
- Database structure
- Failure scenarios
- Deployment environment

before implementing the system.

---

## 2. Integration First

Many of my projects involve systems that need to communicate with something outside the application itself.

Examples:

```text
Software ↔ Hikvision
Software ↔ Meta WhatsApp
Software ↔ Telegram
Software ↔ Modbus Device
Software ↔ SNMP
Software ↔ CCTV
```

Because of this, I pay particular attention to:

- API behavior
- Protocol limitations
- Authentication
- Error handling
- Retry mechanism
- Synchronization
- Data consistency

---

## 3. Design Around Constraints

I don't consider a technology good simply because it is popular.

The technology should match:

- Hardware resources
- Expected traffic
- Complexity
- Deployment environment
- Maintenance requirements

The Battery Monitoring POC was an important example of this lesson.

---

## 4. Separate Responsibilities

For larger systems, I prefer separating responsibilities so individual components remain easier to maintain.

For example, Omni-chan separates:

```text
Webhook
   ↓
Message Distribution
   ↓
Business Logic
   ↓
Database / Cache
   ↓
Realtime Communication
   ↓
Frontend
```

---

# AI-Assisted Development

I use AI as a **development assistant**, not as a replacement for engineering decisions.

I generally use AI to help with:

- Implementation
- Debugging
- Documentation
- Exploring alternatives
- Understanding unfamiliar libraries
- Generating boilerplate
- Reviewing approaches

I remain responsible for:

- Requirements
- Architecture
- Database design
- System flow
- Integration decisions
- Testing
- Debugging validation
- Final implementation

My approach is:

```text
Understand
    ↓
Design
    ↓
Ask AI / Research
    ↓
Implement
    ↓
Test
    ↓
Debug
    ↓
Validate
```

---

# What I Enjoy Building

I am particularly interested in systems that sit between multiple technologies.

Examples:

```text
API
 │
 ├── Backend
 │
 ├── Database
 │
 ├── Message Broker
 │
 ├── Hardware
 │
 ├── AI
 │
 └── Frontend
```

I enjoy solving problems where the challenge is not only writing code, but making different systems communicate reliably.

---

# Career Positioning

## Target Role

**Backend & Integration Engineer**

## Strongest Areas

1. Backend Development
2. API Integration
3. System Integration
4. Automation
5. Hardware Integration
6. Database Design
7. Realtime Systems
8. Edge Computing
9. AI Deployment

---

# Selected Projects

| Project | Main Focus | Role |
|---|---|---|
| Omni-chan | WhatsApp, Backend, Integration, Realtime | Full-stack / Lead |
| Immersion Controller | Automation, Modbus, Edge | Full-stack / Lead |
| Visitor Management | Hikvision API, Access Control | Integration Engineer |
| Asset Management | Database & Asset System | Full-stack |
| Visibel Edge | Computer Vision, Edge AI | Integration Engineer |
| Auto-bon | Telegram, Automation, Product | Full-stack / Lead |
| Battery Monitoring | Edge POC, Architecture | Full-stack |

---

# Portfolio Contact

**Name:** Yo Shandy Hermawan

**Email:** yosh.softdev@gmail.com

**GitHub:** https://github.com/yoshsoftdev

**LinkedIn:** linkedin.com/in/yo-shandy-h-80850012a

