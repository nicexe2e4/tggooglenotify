# tggooglenotify

A simple Go application that parses your Google duty calendar and sends a notification to Telegram with the name of the duty person.

## Prerequisites

To work, you need to create a service account through the Google Console and copy the credentials in JSON format.

## Environment Variables

You need to fill in the following environment variables in the Dockerfile:

- `TELEGRAM_TOKEN` - Telegram bot token
- `CHAT_ID` - ID of the chat where the notification should be sent
- `CALENDAR_ID` - ID of the calendar from which to parse events
- `CRON_EXPRESSION` - Cron expression specifying when to run the parsing
- `GOOGLE_APPLICATION_CREDENTIALS` - Path to the service account credentials JSON file
- `USER_TAGS` - match name in cal Events with tg login in JSON format

## Usage

1. **Create a service account** in the Google Console and download the credentials JSON file.

2. **Build the Docker image** with the necessary environment variables:

3. **Run the Docker container**

## Notes

- Ensure that the `GOOGLE_APPLICATION_CREDENTIALS` path is correctly set to the location of your service account credentials JSON file.
- The cron expression `0 9 * * 1-5` schedules the task to run at 9 AM from Monday to Friday.
- `USER_TAGS` example: '{"Name in Events":"@tg_login","Ivan":"@ivan"}'
## License

This project is licensed under the MIT License.
