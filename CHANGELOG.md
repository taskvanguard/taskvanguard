# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

- Fix: Add universal method for cleaning up LLM respones
- Change: Switch default model to gpt4.1-mini because its much better suited
- Feature: Add +ai tag for evaluation of tasks that can be done via LLM

## [0.2.7] - 2025-08-12

- Removed CLA requirement and updated documentation accordingly
- Fix: Remove unsupported temperature parameter for models newer than GPT-3.5 (thanks @rubdos)
- Fix: Parse TaskWarrior Skipped field as float64 instead of int (thanks @rubdos)
- Changed: Spot: Prompts now with [y]es (starting task)/[s]kip (skipping task)/[n]ext (adding +next tag) 

## [0.2.6] - 2025-08-11

- Feature: Spot: Provide annotations and skipped count to the LLM promt for context
- Feature: Spot: Add goals associated with an Task to the LLM prompt for context
- Bug: Spot: Fix addressing of user in prompt response
- Bug: Spot: Fix starting task if user accepts prompt

## [0.2.5] - 2025-08-10

- Feature: Improve Mass Analysis via Editor Mode
- Bug: Instruct LLM to not put punctuation at the end of Task descriptions

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
