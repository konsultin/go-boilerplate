FROM alpine:latest

# Install curl
RUN apk add --no-cache curl

# Copy the crontab file
COPY crontab /etc/crontabs/root
COPY run.sh /usr/local/bin/run
RUN chmod +x /usr/local/bin/run

# Start cron daemon
# -f: Foreground
# -d 8: Log level (8 = default)
CMD ["crond", "-f", "-d", "8"]
