package main

// Database Initialization
const DBSocket = "0.0.0.0:5432"
const DBUser = "postgres"
const DBPass = "pp"
const DBName = "postgres"

// Web Server Default Config
const ListeningSocket = "0.0.0.0:8001"
const ShutdownGraceSeconds = 3

// HTTP Endpoints
const EndpointLive = "/live"
const EndpointNewUser = "/newuser"
const EndpointLogin = "/login"

// Token generation
const TokenSecret = "secret"
const TokenExpirationSeconds = 3600

// Login Form Field names
const LoginFormUser = "code"
const LoginFormPass = "pass"

// New User Form Field names
const NewUserFormName = "name"
const NewUserFormUserCode = "code"
const NewUserFormPass = "pass"
const NewUserFormEmail = "email"

// New User Creation
const SaltLength = 8

// User Security Groups
const SecuGrpAdmin = -1
const SecuGrpUser = 1
