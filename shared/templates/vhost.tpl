server {
    listen       {{ .Port }};
    server_name  {{ .Domain }}{{ .DomainSuffix }};
    root   {{ .ProjectPath }}/www;

    access_log {{ .ProjectPath }}/log/access.log;
    error_log {{ .ProjectPath }}/log/error.log;

    listen 443 http2 ssl;
    
    ssl on;
    ssl_certificate {{ .ProjectPath }}/etc/cert.pem;
    ssl_certificate_key {{ .ProjectPath }}/etc/key.pem;

{{ .Preset }}
}
