### wtc_system

### Status Of Deployment
[![CD](https://github.com/hawkkiller/wtc_system/actions/workflows/dev.workflow.yml/badge.svg?branch=main)](https://github.com/hawkkiller/wtc_system/actions/workflows/dev.workflow.yml)
[![Go](https://img.shields.io/badge/1.18-golang-blue)](https://github.com/golang)

#### services:
    user_service
        errors:
            wtc-000001: sys error
            wtc-000002: bad request
            wtc-000003: not found
            wtc-000004: reg error(email already registered)
            wtc-000005: invalid email or password
