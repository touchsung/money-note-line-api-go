# Money Note Line API (Go)
This is a Go language-based API that allows users to interact with the Money Note Line application. With this API, developers can programmatically access and manipulate user data in the Money Note Line application, such as retrieving account information, getting transaction history, and creating transactions.

## Installation
To use this API, you will need to have Go installed on your machine. Once Go is installed, you can install the API by running the following command:

```
go get github.com/Touchsung/money-note-line-api-go
```
This will install the API and its dependencies into your Go environment.

## Usage
To use the API in your Go application, you will need to import it:

```
import "github.com/Touchsung/money-note-line-api-go"
```
From there, you can use the API's functions to interact with the Money Note Line application. For example, to retrieve a user's account information, you can use the following code:

## Environment Variables

This project uses the following environment variables to store sensitive information and configuration options:

1. CHANNEL_ACCESS_TOKEN: The access token required to authenticate the LINE Messaging API client.
2. CHANNEL_SECRET: The channel secret required to authenticate the LINE Messaging API client.
3. WIT_AI_TOKEN: The access token required to authenticate the Wit.ai API client.
4. DB_URL: The connection string required to connect to the PostgreSQL database.

To use this project, make sure to set these environment variables in your local development environment or in the deployment environment where the application will run.

## Live Application
If you want to see the application in action on LINE, you can add the LINE ID `@900ggjgm` as a friend and send messages to the account. Please note that the application currently only supports Thai language. The application is running on a live server, so you can test its functionality in a real environment.

Please also note that the application may not respond immediately, as it may be busy serving other users. Also, keep in mind that the application may be updated or modified at any time, so the features and behavior of the application may change without notice.

If you encounter any issues or have any feedback about the application, please don't hesitate to contact us at `jettapat.th@gmail.com` and we'll do our best to address your concerns.

## Contributing
Contributions to this project are welcome! If you would like to contribute, please fork the project, make your changes, and submit a pull request. Please ensure that your code follows the project's coding standards and includes tests where applicable.

## License
This API is licensed under the MIT License. Please see the LICENSE file for more information.

