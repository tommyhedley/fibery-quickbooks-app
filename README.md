# Fibery-quickbooks-app
Fibery-quickbooks-app is a custom integration app for fibery.io. It pulls implemented datatypes from QuickBooks online and converts them to Fibery schema and data. Full, delta (Change Data Capture in QuickBooks), and webhook sync options are all available based on the methods available for each datatype in the Quickbooks API. Currently there is no sync-back/2-way sync option, but that may be implemented in the future using Fibery custome integration actions or the Fibery 2-way sync API when available.


[Fibery Custom Integration API Docs](https://the.fibery.io/@public/User_Guide/Guide/Integrations-API-267)

[QuickBooks Online API Docs](https://developer.intuit.com/app/developer/qbo/docs/get-started)

## Environment Variables
Fibery-quickbooks-app uses a number of enviroment variables to spefic QuickBooks API parameters. For testing I use a .env file and for production you options will vary depending on how it is hosted. They are as follows:

#### Mode
> [!Note]
> Options: "production" or "sandbox"

MODE="sandbox"

#### Server Configuration
PORT="8080"
LOGGER_LEVEL="info"
LOGGER_STYLE="json"

#### Quickbooks Configuration 
MINOR_VERSION="75"
SCOPE="com.intuit.quickbooks.accounting openid profile email phone address"
WEBHOOK_TOKEN=""

#### Quickbooks Token Refresh Before Expiration Time (In Seconds)
TOKEN_REFRESH_BEFORE_EXPIRATION="600"

#### Quickbooks Discover URLS
DISCOVERY_SANDBOX_ENDPOINT="https://developer.api.intuit.com/.well-known/openid_sandbox_configuration"
DISCOVERY_PRODUCTION_ENDPOINT="https://developer.api.intuit.com/.well-known/openid_configuration"

#### Quickbooks API Endpoints
SANDBOX_ENDPOINT="https://sandbox-quickbooks.api.intuit.com"
PRODUCTION_ENDPONT="https://quickbooks.api.intuit.com"

#### Quickbooks Client IDs and Secrets
OAUTH_CLIENT_ID_SANDBOX=""
OAUTH_CLIENT_SECRET_SANDBOX=""
OAUTH_CLIENT_ID_PRODUCTION=""
OAUTH_CLIENT_SECRET_PRODUCTION=""

## Data Types
> [!Note]
> This app does not comprehensivley implement all possible datatypes. Please feel free to fork if you would like to implement more types.
