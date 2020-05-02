package commons

// Context
const CTXTimeOutSecs = 5

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
const EndpointNewTeam = "/newteam"
const EndpointLogin = "/login"

// Token generation
const TokenSecret = "secret"
const TokenExpirationSeconds = 3600

// Login Form Field names
const LoginFormUserCode = "code"
const LoginFormPass = "pass"

// New User Form fieldnames
const NewUserFormTeamID = "teamid"
const NewUserFormUserCode = "code"
const NewUserFormPass = "pass"
const NewUserFormGroupID = "groupid"

// New User Creation
const SaltLength = 8

// New Team Form fieldnames
const NewTeamFormName = "name"
const NewTeamFormDescription = "desc"
const NewTeamFormCODE = "code"

// User Security Groups
const SecuAppAdmin = 4
const SecuGroupAdmin = 3
const SecuTeamAdmin = 2
const SecuGrpUser = 1
