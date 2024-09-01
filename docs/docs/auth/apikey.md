# API Key Authentication

This document describes the API Key authentication mechanism implemented in our API.

## Overview

API Key authentication is a simple and straightforward method for authenticating API requests. It involves including a pre-shared key in the request headers.

## Configuration

The following environment variable needs to be set:

- `VALID_API_KEY`: The valid API key that clients should use

## Usage

To authenticate API requests, include the API key in the `X-API-Key` header: