Simple web server that follows the boot.dev course assigments.

Use the [http](https://pkg.go.dev/net/http) package to create a simple web server.

Add a [middleware](https://developer.mozilla.org/en-US/docs/Glossary/Middleware) function
that adds [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) headers to the response.

Play with stateful handlers.

Play with [Chi](https://go-chi.io/#/README) router.

A .env file is used to store the secret use to sign the JWT token.
This file is added to .gitignore for security reasons.
You can create your own .env file with the following content:
You can use ```openssl rand -base64 64``` to generate a secret key.

```
# .env
JWT_SECRET=your-secret-key
```





