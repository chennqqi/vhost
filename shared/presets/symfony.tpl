    location / {
        try_files $uri /app.php$is_args$args;
    }
    