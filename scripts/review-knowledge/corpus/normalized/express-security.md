# Production Best Practices: Security

The term *“production”* refers to the stage in the software lifecycle when an application or API is generally available to its end-users or consumers. In contrast, in the *“development”* stage, you’re still actively writing and testing code, and the application is not open to external access. The corresponding system environments are known as *production* and *development* environments, respectively.

Development and production environments are usually set up differently and have vastly different requirements. What’s fine in development may not be acceptable in production. For example, in a development environment you may want verbose logging of errors for debugging, while the same behavior can become a security concern in a production environment. And in development, you don’t need to worry about scalability, reliability, and performance, while those concerns become critical in production.

### Note

If you believe you have discovered a security vulnerability in Express, please see Security Policies and Procedures.

Security best practices for Express applications in production include:

- Don’t use deprecated or vulnerable versions of Express
- Use TLS
- Do not trust user input
- Use Helmet
- Reduce fingerprinting
- Use cookies securely
- Prevent brute-force attacks against authorization
- Ensure your dependencies are secure
- Additional considerations

## Don’t use deprecated or vulnerable versions of Express

Express 2.x and 3.x are no longer maintained. Security and performance issues in these versions won’t be fixed. Do not use them! If you haven’t moved to version 4, follow the migration guide or consider Commercial Support Options.

Also ensure you are not using any of the vulnerable Express versions listed on the Security updates page. If you are, update to one of the stable releases, preferably the latest.

## Use TLS

If your app deals with or transmits sensitive data, use Transport Layer Security (TLS) to secure the connection and the data. This technology encrypts data before it is sent from the client to the server, thus preventing some common (and easy) hacks. Although Ajax and POST requests might not be visibly obvious and seem “hidden” in browsers, their network traffic is vulnerable to packet sniffing and man-in-the-middle attacks.

You may be familiar with Secure Socket Layer (SSL) encryption. TLS is simply the next progression of SSL. In other words, if you were using SSL before, consider upgrading to TLS. In general, we recommend Nginx to handle TLS. For a good reference to configure TLS on Nginx (and other servers), see Recommended Server Configurations (TLSRef).

Also, a handy tool to get a free TLS certificate is Let’s Encrypt, a free, automated, and open certificate authority (CA) provided by the Internet Security Research Group (ISRG).

## Do not trust user input

For web applications, one of the most critical security requirements is proper user input validation and handling. This comes in many forms and we will not cover all of them here. Ultimately, the responsibility for validating and correctly handling the types of user input your application accepts is yours.

### Prevent open redirects

An example of potentially dangerous user input is an *open redirect*, where an application accepts a URL as user input (often in the URL query, for example `?url=https://example.com`) and uses `res.redirect` to set the `location` header and
return a 3xx status.

An application must validate that it supports redirecting to the incoming URL to avoid sending users to malicious links such as phishing websites, among other risks.

Here is an example of checking URLs before using `res.redirect` or `res.location`:

## Use Helmet

Helmet can help protect your app from some well-known web vulnerabilities by setting HTTP headers appropriately.

Helmet is a middleware function that sets security-related HTTP response headers. Helmet sets the following headers by default:

- `Content-Security-Policy`: A powerful allow-list of what can happen on your page which mitigates many attacks
- `Cross-Origin-Opener-Policy`: Helps process-isolate your page
- `Cross-Origin-Resource-Policy`: Blocks others from loading your resources cross-origin
- `Origin-Agent-Cluster`: Changes process isolation to be origin-based
- `Referrer-Policy`: Controls the- `Referer`header
- `Strict-Transport-Security`: Tells browsers to prefer HTTPS
- `X-Content-Type-Options`: Avoids MIME sniffing
- `X-DNS-Prefetch-Control`: Controls DNS prefetching
- `X-Download-Options`: Forces downloads to be saved (Internet Explorer only)
- `X-Frame-Options`: Legacy header that mitigates Clickjacking attacks
- `X-Permitted-Cross-Domain-Policies`: Controls cross-domain behavior for Adobe products, like Acrobat
- `X-Powered-By`: Info about the web server. Removed because it could be used in simple attacks
- `X-XSS-Protection`: Legacy header that tries to mitigate XSS attacks, but makes things worse, so Helmet disables it

Each header can be configured or disabled. To read more about it please go to its documentation website.

Install Helmet like any other module:

Then to use it in your code:

## Reduce fingerprinting

It can help to provide an extra layer of security to reduce the ability of attackers to determine the software that a server uses, known as “fingerprinting.” Though not a security issue itself, reducing the ability to fingerprint an application improves its overall security posture. Server software can be fingerprinted by quirks in how it responds to specific requests, for example in the HTTP response headers.

By default, Express sends the `X-Powered-By` response header that you can
disable using the `app.disable()` method:

### Note

Disabling the `X-Powered-By header` does not prevent a sophisticated attacker from determining
that an app is running Express. It may discourage a casual exploit, but there are other ways to
determine an app is running Express.

Express also sends its own formatted “404 Not Found” messages and formatter error response messages. These can be changed by adding your own not found handler and writing your own error handler:

## Use cookies securely

To ensure cookies don’t open your app to exploits, don’t use the default session cookie name and set cookie security options appropriately.

There are two main middleware cookie session modules:

- express-session that replaces `express.session`middleware built-in to Express 3.x.
- cookie-session that replaces `express.cookieSession`middleware built-in to Express 3.x.

The main difference between these two modules is how they save cookie session data. The express-session middleware stores session data on the server; it only saves the session ID in the cookie itself, not session data. By default, it uses in-memory storage and is not designed for a production environment. In production, you’ll need to set up a scalable session-store; see the list of compatible session stores.

In contrast, cookie-session middleware implements cookie-backed storage: it serializes the entire session to the cookie, rather than just a session key. Only use it when session data is relatively small and easily encoded as primitive values (rather than objects). Although browsers are supposed to support at least 4096 bytes per cookie, to ensure you don’t exceed the limit, don’t exceed a size of 4093 bytes per domain. Also, be aware that the cookie data will be visible to the client, so if there is any reason to keep it secure or obscure, then `express-session` may be a better choice.

### Don’t use the default session cookie name

Using the default session cookie name can open your app to attacks. The security issue posed is similar to `X-Powered-By`: a potential attacker can use it to fingerprint the server and target attacks accordingly.

To avoid this problem, use generic cookie names; for example using express-session middleware:

### Set cookie security options

Set the following cookie options to enhance security:

- `secure`- Ensures the browser only sends the cookie over HTTPS.
- `httpOnly`- Ensures the cookie is sent only over HTTP(S), not client JavaScript, helping to protect against cross-site scripting attacks.
- `domain`- indicates the domain of the cookie; use it to compare against the domain of the server in which the URL is being requested. If they match, then check the path attribute next.
- `path`- indicates the path of the cookie; use it to compare against the request path. If this and domain match, then send the cookie in the request.
- `expires`- use to set expiration date for persistent cookies.

Here is an example using cookie-session middleware:

## Prevent brute-force attacks against authorization

Make sure login endpoints are protected to make private data more secure.

A simple and powerful technique is to block authorization attempts using two metrics:

- The number of consecutive failed attempts by the same user name and IP address.
- The number of failed attempts from an IP address over some long period of time. For example, block an IP address if it makes 100 failed attempts in one day.

rate-limiter-flexible package provides tools to make this technique easy and fast. You can find an example of brute-force protection in the documentation

## Ensure your dependencies are secure

Using npm to manage your application’s dependencies is powerful and convenient. But the packages that you use may contain critical security vulnerabilities that could also affect your application. The security of your app is only as strong as the “weakest link” in your dependencies.

Since npm@6, npm automatically reviews every install request. Also, you can use `npm audit` to analyze your dependency tree.

If you want to stay more secure, consider Snyk.

Snyk offers both a command-line tool and a GitHub integration that checks your application against Snyk’s open source vulnerability database for any known vulnerabilities in your dependencies. Install the CLI as follows:

Use this command to test your application for vulnerabilities:

### Avoid other known vulnerabilities

Keep an eye out for GitHub Advisory Database or Snyk advisories that may affect Express or other modules that your app uses. In general, these databases are excellent resources for knowledge and tools about Node security.

Finally, Express apps—like any other web apps—can be vulnerable to a variety of web-based attacks. Familiarize yourself with known web vulnerabilities and take precautions to avoid them.

## Additional considerations

Here are some further recommendations from the excellent Node.js Security Checklist. Refer to that blog post for all the details on these recommendations:

- Always filter and sanitize user input to protect against cross-site scripting (XSS) and command injection attacks.
- Defend against SQL injection attacks by using parameterized queries or prepared statements.
- Use the open-source sqlmap tool to detect SQL injection vulnerabilities in your app.
- Use the nmap and sslyze tools to test the configuration of your SSL ciphers, keys, and renegotiation as well as the validity of your certificate.
- Use safe-regex to ensure your regular expressions are not susceptible to regular expression denial of service attacks.
