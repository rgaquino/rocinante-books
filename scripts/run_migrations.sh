#!/usr/bin/env bash

migrate -path "./migrations" -database $1 up