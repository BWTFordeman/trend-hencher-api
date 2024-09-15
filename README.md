# trend-hencher-api


## Table of Contents
- [Installation](#installation)

## Installation
To install this project, follow these steps:

1. Clone the repo
```bash
git clone https://github.com/BWTFordeman/tred-hencher-api.git
```
2. Setup environment variables
### Local variable
    $env:ENVIRONMENT="local"
### Credentials variable
    $env:GOOGLE_APPLICATION_CREDENTIALS="location\service-account-file.json"
### Cloud project ID
    $env:GOOGLE_CLOUD_PROJECT="id"     

3. Run the project
 ```bash
go run main.go
```