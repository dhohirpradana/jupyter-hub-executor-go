FROM nginx

WORKDIR /app

COPY dpapp /app
COPY .env /app

EXPOSE 9090

CMD ./dpapp