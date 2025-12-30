# Changelog

All notable changes to this project will be documented in this file.

The format is based on **Keep a Changelog** and adheres to **Semantic Versioning**.

---

## [Unreleased]

### Added

* Foundation for proxy rotation and session management (design-ready)
* Platform-isolated service architecture for OSINT scraping

### Changed

* Refactored API routing to support multiple social media platforms

---

## [0.2.0] - 2025-12-30

### Added

* Instagram module with isolated routes, handlers, services, and models
* Makefile with build, run, swagger, docker, and k8s commands

### Changed

* Migrated from single-platform (Snapchat-only) design to multi-platform architecture
* Improved project structure following Go clean architecture principles

### Fixed

* Swagger generation path alignment with `cmd/api/main.go`
* Consistent API versioning and route grouping

---

## [0.1.0] - 2024-12-20

### Added

* Initial project setup with Go and Gin
* Snapchat ping endpoint
* Basic health check endpoint
* Swagger UI integration
* Docker support

---

## Versioning

This project follows **Semantic Versioning**:

* **MAJOR** version when you make incompatible API changes
* **MINOR** version when you add functionality in a backward-compatible manner
* **PATCH** version when you make backward-compatible bug fixes

---

Maintained by **Er Vijay Raghuwanshi**
