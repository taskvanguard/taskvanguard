# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

- Feature: Add support for exporting tasks to workflow automation plattform (n8n) 
- Feature: Add support for more LLM Apis
- Feature: Add support for notifications
- Fix: Guide cmd does not assign goals to tasks that are generated
- Fix: init.go may has no access to .taskrc global var

## [0.2.4] - 2025-06-20

- Improved --help contents

## [0.2.3] - 2025-06-20

### Added

- Goal Management via `vanguard goals`
- Backup of TaskWarrior Config

### Changed

- Updated Readme
- Support for go install by moving main.go -> cmd/vanguard/
- Increased Output verbosity when backupping files

### Fixed

- Sanitizing LLM response in analysis mode



## [0.1.0] - 2025-06-12

### Added
- Initial public release. Core functionality implemented.
