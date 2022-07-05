# Freshdesk Line Integration Gateway API

# Version
1. Go 1.17.*
2. MongoDB 4.4

## Prerequisite
1. Docker
2. Docker Compose

## Setup
1. Clone this repository
2. Go to repository directory
3. Config environment
    1. Rename .env.example to .env using the command `cp .env.example .env`
    2. Enter application port (APP_PORT)
    3. Config MongoDB on variable with prefix MONGO_
    4. Config SMTP Gmail
    5. Config Domain WEBHOOK on Ngrok
4. Start container using the command `docker-compose up -d`
5. Access to application `http://localhost:${APP_PORT}` and you will see welcome message

## Command function for user
1. Header {"x-signature": "<x-signature:string>"}
2. Body => JSON (Content-Type: application/json)
    1. Create USER => {"name":"<name:string>","plan":"<userPlan:int>"}
    2. Delete USER => {"userID":"<userID:string>"}
    3. Regenerate Activate Key => {"userID":"<userID:string>"}
    4. Change Plan => {"activateKey":"<ActivateKey:string>","plan":"<userPlan:int>"}
    5. Change status Activate Key => {"activateKey":"<ActivateKey:string>","status":"<statusKey:bool>"}