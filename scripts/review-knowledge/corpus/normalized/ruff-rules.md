# Rules

**Ruff supports over 900 lint rules**, many of which are inspired by popular tools like Flake8,
isort, pyupgrade, and others. Regardless of the rule's origin, Ruff re-implements every rule in
Rust as a first-party feature.

By default, Ruff enables rules from the `F`, `E`, `B`, `UP`, and `RUF` categories,
as well as many more, omitting any stylistic rules that overlap with the use of a formatter, like
`ruff format` or Black.

If you're just getting started with Ruff, **the default rule set is a great place to start**: it
catches a wide variety of common errors (like unused imports) with zero configuration. See
*Default Rules* for the complete list.

## Legend

    🧪     The rule is unstable and is in "preview".

    ⚠️     The rule has been deprecated and will be removed in a future release.

    ❌     The rule has been removed only the documentation is available.

    🛠️     The rule is automatically fixable by the `--fix` command-line option.

All rules not marked as preview, deprecated or removed are stable.

## Airflow (AIR)

For more, see Airflow on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| AIR001 | airflow-variable-name-task-id-mismatch | Task variable name should match the `task_id`: "{task_id}" | Rule has been stable since v0.0.271 |
| AIR002 | airflow-dag-no-schedule-argument | `DAG`or`@dag`should have an explicit`schedule`argument | Rule has been stable since 0.13.0 |
| AIR003 | airflow-variable-get-outside-task | `Variable.get()`outside of a task | Rule has been in preview since 0.15.6 |
| AIR004 | airflow-task-branch-as-short-circuit | `@task.branch`can be replaced with`@task.short_circuit` | Rule has been in preview since 0.15.12 |
| AIR201 | airflow-xcom-pull-in-template-string | Use the `.output`attribute on the task object for "{task_id}" instead of`xcom_pull`in a template string | Rule has been in preview since 0.15.11Automatic fix available |
| AIR202 | airflow-task-implicit-multiple-outputs | `@task`-decorated function relies on`multiple_outputs`inference | Rule has been in preview since 0.15.14Automatic fix available |
| AIR301 | airflow3-removal | `{deprecated}`is removed in Airflow 3.0 | Rule has been stable since 0.13.0Automatic fix available |
| AIR302 | airflow3-moved-to-provider | `{deprecated}`is moved into`{provider}`provider in Airflow 3.0; | Rule has been stable since 0.13.0Automatic fix available |
| AIR303 | airflow3-incompatible-function-signature | `{function_name}`signature is changed in Airflow 3.0 | Rule has been stable since 0.16.0 |
| AIR304 | airflow3-dag-dynamic-value | `{function_name}()`produces a value that changes at runtime; using it in a Dag or task argument causes infinite Dag version creation | Rule has been in preview since 0.15.6 |
| AIR311 | airflow3-suggested-update | `{deprecated}`is removed in Airflow 3.0; It still works in Airflow 3.0 but is expected to be removed in a future version. | Rule has been stable since 0.13.0Automatic fix available |
| AIR312 | airflow3-suggested-to-move-to-provider | `{deprecated}`is deprecated and moved into`{provider}`provider in Airflow 3.0; It still works in Airflow 3.0 but is expected to be removed in a future version. | Rule has been stable since 0.13.0Automatic fix available |
| AIR321 | airflow31-moved | `{deprecated}`is moved in Airflow 3.1 | Rule has been in preview since 0.15.1Automatic fix available |

## eradicate (ERA)

For more, see eradicate on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| ERA001 | commented-out-code | Found commented-out code | Rule has been stable since v0.0.145 |

## FastAPI (FAST)

For more, see FastAPI on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| FAST001 | fast-api-redundant-response-model | FastAPI route with redundant `response_model`argument | Rule has been stable since 0.8.0Automatic fix available |
| FAST002 | fast-api-non-annotated-dependency | FastAPI dependency without `Annotated` | Rule has been stable since 0.8.0Automatic fix available |
| FAST003 | fast-api-unused-path-parameter | Parameter `{arg_name}`appears in route path, but not in`{function_name}`signature | Rule has been stable since 0.10.0Automatic fix available |

## flake8-2020 (YTT)

For more, see flake8-2020 on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| YTT101 | sys-version-slice3 | `sys.version[:3]`referenced (python3.10), use`sys.version_info` | Rule has been stable since v0.0.113 |
| YTT102 | sys-version2 | `sys.version[2]`referenced (python3.10), use`sys.version_info` | Rule has been stable since v0.0.113 |
| YTT103 | sys-version-cmp-str3 | `sys.version`compared to string (python3.10), use`sys.version_info` | Rule has been stable since v0.0.113 |
| YTT201 | sys-version-info0-eq3 | `sys.version_info[0] == 3`referenced (python4), use`>=` | Rule has been stable since v0.0.113 |
| YTT202 | six-py3 | `six.PY3`referenced (python4), use`not six.PY2` | Rule has been stable since v0.0.113 |
| YTT203 | sys-version-info1-cmp-int | `sys.version_info[1]`compared to integer (python4), compare`sys.version_info`to tuple | Rule has been stable since v0.0.113 |
| YTT204 | sys-version-info-minor-cmp-int | `sys.version_info.minor`compared to integer (python4), compare`sys.version_info`to tuple | Rule has been stable since v0.0.113 |
| YTT301 | sys-version0 | `sys.version[0]`referenced (python10), use`sys.version_info` | Rule has been stable since v0.0.113 |
| YTT302 | sys-version-cmp-str10 | `sys.version`compared to string (python10), use`sys.version_info` | Rule has been stable since v0.0.113 |
| YTT303 | sys-version-slice1 | `sys.version[:1]`referenced (python10), use`sys.version_info` | Rule has been stable since v0.0.113 |

## flake8-annotations (ANN)

For more, see flake8-annotations on PyPI.

For related settings, see flake8-annotations.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| ANN001 | missing-type-function-argument | Missing type annotation for function argument `{name}` | Rule has been stable since v0.0.105 |
| ANN002 | missing-type-args | Missing type annotation for `*{name}` | Rule has been stable since v0.0.105 |
| ANN003 | missing-type-kwargs | Missing type annotation for `**{name}` | Rule has been stable since v0.0.105 |
| ANN101 | missing-type-self | Missing type annotation for `{name}`in method | Rule was removed in 0.8.0 |
| ANN102 | missing-type-cls | Missing type annotation for `{name}`in classmethod | Rule was removed in 0.8.0 |
| ANN201 | missing-return-type-undocumented-public-function | Missing return type annotation for public function `{name}` | Rule has been stable since v0.0.105Automatic fix available |
| ANN202 | missing-return-type-private-function | Missing return type annotation for private function `{name}` | Rule has been stable since v0.0.105Automatic fix available |
| ANN204 | missing-return-type-special-method | Missing return type annotation for special method `{name}` | Rule has been stable since v0.0.105Automatic fix available |
| ANN205 | missing-return-type-static-method | Missing return type annotation for staticmethod `{name}` | Rule has been stable since v0.0.105Automatic fix available |
| ANN206 | missing-return-type-class-method | Missing return type annotation for classmethod `{name}` | Rule has been stable since v0.0.105Automatic fix available |
| ANN401 | any-type | Dynamically typed expressions (typing.Any) are disallowed in `{name}` | Rule has been stable since v0.0.108 |

## flake8-async (ASYNC)

For more, see flake8-async on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| ASYNC100 | cancel-scope-no-checkpoint | A `with {method_name}(...):`context does not contain any`await`statements. This makes it pointless, as the timeout can only be triggered by a checkpoint. | Rule has been stable since v0.0.269 |
| ASYNC105 | trio-sync-call | Call to `{method_name}`is not immediately awaited | Rule has been stable since 0.5.0Automatic fix available |
| ASYNC109 | async-function-with-timeout | Async function definition with a `timeout`parameter | Rule has been stable since 0.5.0 |
| ASYNC110 | async-busy-wait | Use `{module}.Event`instead of awaiting`{module}.sleep`in a`while`loop | Rule has been stable since 0.5.0 |
| ASYNC115 | async-zero-sleep | Use `{module}.lowlevel.checkpoint()`instead of`{module}.sleep(0)` | Rule has been stable since 0.5.0Automatic fix available |
| ASYNC116 | long-sleep-not-forever | `{module}.sleep()`with >24 hour interval should usually be`{module}.sleep_forever()` | Rule has been stable since 0.13.0Automatic fix available |
| ASYNC119 | yield-in-context-manager-in-async-generator | Yield in context manager in async generator may not trigger cleanup | Rule has been in preview since 0.15.16 |
| ASYNC210 | blocking-http-call-in-async-function | Async functions should not call blocking HTTP methods | Rule has been stable since 0.5.0 |
| ASYNC212 | blocking-http-call-httpx-in-async-function | Blocking httpx method {name}.{call}() in async context, use httpx.AsyncClient | Rule has been stable since 0.15.0 |
| ASYNC220 | create-subprocess-in-async-function | Async functions should not create subprocesses with blocking methods | Rule has been stable since 0.5.0 |
| ASYNC221 | run-process-in-async-function | Async functions should not run processes with blocking methods | Rule has been stable since 0.5.0 |
| ASYNC222 | wait-for-process-in-async-function | Async functions should not wait on processes with blocking methods | Rule has been stable since 0.5.0 |
| ASYNC230 | blocking-open-call-in-async-function | Async functions should not open files with blocking methods like `open` | Rule has been stable since 0.5.0 |
| ASYNC240 | blocking-path-method-in-async-function | Async functions should not use {path_library} methods, use trio.Path or anyio.path | Rule has been stable since 0.15.0 |
| ASYNC250 | blocking-input-in-async-function | Blocking call to `input()`in async context | Rule has been stable since 0.15.0 |
| ASYNC251 | blocking-sleep-in-async-function | Async functions should not call `time.sleep` | Rule has been stable since 0.5.0 |

## flake8-bandit (S)

For more, see flake8-bandit on PyPI.

For related settings, see flake8-bandit.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| S101 | assert | Use of `assert`detected | Rule has been stable since v0.0.116 |
| S102 | exec-builtin | Use of `exec`detected | Rule has been stable since v0.0.116 |
| S103 | bad-file-permissions | `os.chmod`setting a permissive mask`{mask:#o}`on file or directory | Rule has been stable since v0.0.211 |
| S104 | hardcoded-bind-all-interfaces | Possible binding to all interfaces | Rule has been stable since v0.0.116 |
| S105 | hardcoded-password-string | Possible hardcoded password assigned to: "{}" | Rule has been stable since v0.0.116 |
| S106 | hardcoded-password-func-arg | Possible hardcoded password assigned to argument: "{}" | Rule has been stable since v0.0.116 |
| S107 | hardcoded-password-default | Possible hardcoded password assigned to function default: "{}" | Rule has been stable since v0.0.116 |
| S108 | hardcoded-temp-file | Probable insecure usage of temporary file or directory: "{}" | Rule has been stable since v0.0.211 |
| S110 | try-except-pass | `try`-`except`-`pass`detected, consider logging the exception | Rule has been stable since v0.0.237 |
| S112 | try-except-continue | `try`-`except`-`continue`detected, consider logging the exception | Rule has been stable since v0.0.245 |
| S113 | request-without-timeout | Probable use of `{module}`call without timeout | Rule has been stable since v0.0.213 |
| S201 | flask-debug-true | Use of `debug=True`in Flask app detected | Rule has been stable since v0.2.0 |
| S202 | tarfile-unsafe-members | Uses of `tarfile.extractall()` | Rule has been stable since v0.2.0 |
| S301 | suspicious-pickle-usage | `pickle`and modules that wrap it can be unsafe when used to deserialize untrusted data, possible security issue | Rule has been stable since v0.0.258 |
| S302 | suspicious-marshal-usage | Deserialization with the `marshal`module is possibly dangerous | Rule has been stable since v0.0.258 |
| S303 | suspicious-insecure-hash-usage | Use of insecure MD2, MD4, MD5, or SHA1 hash function | Rule has been stable since v0.0.258 |
| S304 | suspicious-insecure-cipher-usage | Use of insecure cipher, replace with a known secure cipher such as AES | Rule has been stable since v0.0.258 |
| S305 | suspicious-insecure-cipher-mode-usage | Use of insecure block cipher mode, replace with a known secure mode such as CBC or CTR | Rule has been stable since v0.0.258 |
| S306 | suspicious-mktemp-usage | Use of insecure and deprecated function ( `mktemp`) | Rule has been stable since v0.0.258 |
| S307 | suspicious-eval-usage | Use of possibly insecure function; consider using `ast.literal_eval` | Rule has been stable since v0.0.258 |
| S308 | suspicious-mark-safe-usage | Use of `mark_safe`may expose cross-site scripting vulnerabilities | Rule has been stable since v0.0.258 |
| S310 | suspicious-url-open-usage | Audit URL open for permitted schemes. Allowing use of `file:`or custom schemes is often unexpected. | Rule has been stable since v0.0.258 |
| S311 | suspicious-non-cryptographic-random-usage | Standard pseudo-random generators are not suitable for cryptographic purposes | Rule has been stable since v0.0.258 |
| S312 | suspicious-telnet-usage | Telnet is considered insecure. Use SSH or some other encrypted protocol. | Rule has been stable since v0.0.258 |
| S313 | suspicious-xmlc-element-tree-usage | Using `xml`to parse untrusted data is known to be vulnerable to XML attacks; use`defusedxml`equivalents | Rule has been stable since v0.0.258 |
| S314 | suspicious-xml-element-tree-usage | Using `xml`to parse untrusted data is known to be vulnerable to XML attacks; use`defusedxml`equivalents | Rule has been stable since v0.0.258 |
| S315 | suspicious-xml-expat-reader-usage | Using `xml`to parse untrusted data is known to be vulnerable to XML attacks; use`defusedxml`equivalents | Rule has been stable since v0.0.258 |
| S316 | suspicious-xml-expat-builder-usage | Using `xml`to parse untrusted data is known to be vulnerable to XML attacks; use`defusedxml`equivalents | Rule has been stable since v0.0.258 |
| S317 | suspicious-xml-sax-usage | Using `xml`to parse untrusted data is known to be vulnerable to XML attacks; use`defusedxml`equivalents | Rule has been stable since v0.0.258 |
| S318 | suspicious-xml-mini-dom-usage | Using `xml`to parse untrusted data is known to be vulnerable to XML attacks; use`defusedxml`equivalents | Rule has been stable since v0.0.258 |
| S319 | suspicious-xml-pull-dom-usage | Using `xml`to parse untrusted data is known to be vulnerable to XML attacks; use`defusedxml`equivalents | Rule has been stable since v0.0.258 |
| S320 | suspicious-xmle-tree-usage | Using `lxml`to parse untrusted data is known to be vulnerable to XML attacks | Rule was removed in 0.12.0 |
| S321 | suspicious-ftp-lib-usage | FTP-related functions are being called. FTP is considered insecure. Use SSH/SFTP/SCP or some other encrypted protocol. | Rule has been stable since v0.0.258 |
| S323 | suspicious-unverified-context-usage | Python allows using an insecure context via the `_create_unverified_context`that reverts to the previous behavior that does not validate certificates or perform hostname checks. | Rule has been stable since v0.0.258 |
| S324 | hashlib-insecure-hash-function | Probable use of insecure hash functions in `{library}`:`{string}` | Rule has been stable since v0.0.212 |
| S401 | suspicious-telnetlib-import | `telnetlib`and related modules are considered insecure. Use SSH or another encrypted protocol. | Rule has been in preview since v0.1.12 |
| S402 | suspicious-ftplib-import | `ftplib`and related modules are considered insecure. Use SSH, SFTP, SCP, or another encrypted protocol. | Rule has been in preview since v0.1.12 |
| S403 | suspicious-pickle-import | `pickle`,`cPickle`,`dill`, and`shelve`modules are possibly insecure | Rule has been in preview since v0.1.12 |
| S404 | suspicious-subprocess-import | `subprocess`module is possibly insecure | Rule has been in preview since v0.1.12 |
| S405 | suspicious-xml-etree-import | `xml.etree`methods are vulnerable to XML attacks | Rule has been in preview since v0.1.12 |
| S406 | suspicious-xml-sax-import | `xml.sax`methods are vulnerable to XML attacks | Rule has been in preview since v0.1.12 |
| S407 | suspicious-xml-expat-import | `xml.dom.expatbuilder`is vulnerable to XML attacks | Rule has been in preview since v0.1.12 |
| S408 | suspicious-xml-minidom-import | `xml.dom.minidom`is vulnerable to XML attacks | Rule has been in preview since v0.1.12 |
| S409 | suspicious-xml-pulldom-import | `xml.dom.pulldom`is vulnerable to XML attacks | Rule has been in preview since v0.1.12 |
| S410 | suspicious-lxml-import | `lxml`is vulnerable to XML attacks | Rule was removed in v0.3.0 |
| S411 | suspicious-xmlrpc-import | XMLRPC is vulnerable to remote XML attacks | Rule has been in preview since v0.1.12 |
| S412 | suspicious-httpoxy-import | `httpoxy`is a set of vulnerabilities that affect application code running inCGI, or CGI-like environments. The use of CGI for web applications should be avoided | Rule has been in preview since v0.1.12 |
| S413 | suspicious-pycrypto-import | `pycrypto`library is known to have publicly disclosed buffer overflow vulnerability | Rule has been in preview since v0.1.12 |
| S415 | suspicious-pyghmi-import | An IPMI-related module is being imported. Prefer an encrypted protocol over IPMI. | Rule has been in preview since v0.1.12 |
| S501 | request-with-no-cert-validation | Probable use of `{string}`call with`verify=False`disabling SSL certificate checks | Rule has been stable since v0.0.213 |
| S502 | ssl-insecure-version | Call made with insecure SSL protocol: `{protocol}` | Rule has been stable since v0.2.0 |
| S503 | ssl-with-bad-defaults | Argument default set to insecure SSL protocol: `{protocol}` | Rule has been stable since v0.2.0 |
| S504 | ssl-with-no-version | `ssl.wrap_socket`called without an `ssl_version`` | Rule has been stable since v0.2.0 |
| S505 | weak-cryptographic-key | {cryptographic_key} key sizes below {minimum_key_size} bits are considered breakable | Rule has been stable since v0.2.0 |
| S506 | unsafe-yaml-load | Probable use of unsafe loader `{name}`with`yaml.load`. Allows instantiation of arbitrary objects. Consider`yaml.safe_load`. | Rule has been stable since v0.0.212 |
| S507 | ssh-no-host-key-verification | Paramiko call with policy set to automatically trust the unknown host key | Rule has been stable since v0.2.0 |
| S508 | snmp-insecure-version | The use of SNMPv1 and SNMPv2 is insecure. Use SNMPv3 if able. | Rule has been stable since v0.0.218 |
| S509 | snmp-weak-cryptography | You should not use SNMPv3 without encryption. `noAuthNoPriv`&`authNoPriv`is insecure. | Rule has been stable since v0.0.218 |
| S601 | paramiko-call | Possible shell injection via Paramiko call; check inputs are properly sanitized | Rule has been stable since v0.0.270 |
| S602 | subprocess-popen-with-shell-equals-true | `subprocess`call with`shell=True`seems safe, but may be changed in the future; consider rewriting without`shell` | Rule has been stable since v0.0.262 |
| S603 | subprocess-without-shell-equals-true | `subprocess`call: check for execution of untrusted input | Rule has been stable since v0.0.262 |
| S604 | call-with-shell-equals-true | Function call with `shell=True`parameter identified, security issue | Rule has been stable since v0.0.262 |
| S605 | start-process-with-a-shell | Starting a process with a shell: seems safe, but may be changed in the future; consider rewriting without `shell` | Rule has been stable since v0.0.262 |
| S606 | start-process-with-no-shell | Starting a process without a shell | Rule has been stable since v0.0.262 |
| S607 | start-process-with-partial-path | Starting a process with a partial executable path | Rule has been stable since v0.0.262 |
| S608 | hardcoded-sql-expression | Possible SQL injection vector through string-based query construction | Rule has been stable since v0.0.245 |
| S609 | unix-command-wildcard-injection | Possible wildcard injection in call due to `*`usage | Rule has been stable since v0.0.271 |
| S610 | django-extra | Use of Django `extra`can lead to SQL injection vulnerabilities | Rule has been stable since 0.5.0 |
| S611 | django-raw-sql | Use of `RawSQL`can lead to SQL injection vulnerabilities | Rule has been stable since v0.2.0 |
| S612 | logging-config-insecure-listen | Use of insecure `logging.config.listen`detected | Rule has been stable since v0.0.231 |
| S701 | jinja2-autoescape-false | Using jinja2 templates with `autoescape=False`is dangerous and can lead to XSS. Ensure`autoescape=True`or use the`select_autoescape`function. | Rule has been stable since v0.0.220 |
| S702 | mako-templates | Mako templates allow HTML and JavaScript rendering by default and are inherently open to XSS attacks | Rule has been stable since v0.2.0 |
| S704 | unsafe-markup-use | Unsafe use of `{name}`detected | Rule has been stable since 0.10.0 |

## flake8-blind-except (BLE)

For more, see flake8-blind-except on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| BLE001 | blind-except | Do not catch blind exception: `{name}` | Rule has been stable since v0.0.127 |

## flake8-boolean-trap (FBT)

For more, see flake8-boolean-trap on PyPI.

For related settings, see flake8-boolean-trap.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| FBT001 | boolean-type-hint-positional-argument | Boolean-typed positional argument in function definition | Rule has been stable since v0.0.127 |
| FBT002 | boolean-default-value-positional-argument | Boolean default positional argument in function definition | Rule has been stable since v0.0.127 |
| FBT003 | boolean-positional-value-in-call | Boolean positional value in function call | Rule has been stable since v0.0.127 |

## flake8-bugbear (B)

For more, see flake8-bugbear on PyPI.

For related settings, see flake8-bugbear.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| B002 | unary-prefix-increment-decrement | Python does not support the unary prefix increment operator ( `++`) | Rule has been stable since v0.0.83 |
| B003 | assignment-to-os-environ | Assigning to `os.environ`doesn't clear the environment | Rule has been stable since v0.0.102 |
| B004 | unreliable-callable-check | Using `hasattr(x, "__call__")`to test if x is callable is unreliable. Use`callable(x)`for consistent results. | Rule has been stable since v0.0.106Automatic fix available |
| B005 | strip-with-multi-characters | Using `.strip()`with multi-character strings is misleading | Rule has been stable since v0.0.106 |
| B006 | mutable-argument-default | Do not use mutable data structures for argument defaults | Rule has been stable since v0.0.92Automatic fix available |
| B007 | unused-loop-control-variable | Loop control variable `{name}`not used within loop body | Rule has been stable since v0.0.84Automatic fix available |
| B008 | function-call-in-default-argument | Do not perform function call `{name}`in argument defaults; instead, perform the call within the function, or read the default from a module-level singleton variable | Rule has been stable since v0.0.102 |
| B009 | get-attr-with-constant | Do not call `getattr`with a constant attribute value. It is not any safer than normal property access. | Rule has been stable since v0.0.110Automatic fix available |
| B010 | set-attr-with-constant | Do not call `setattr`with a constant attribute value. It is not any safer than normal property access. | Rule has been stable since v0.0.111Automatic fix available |
| B011 | assert-false | Do not `assert False`(`python -O`removes these calls), raise`AssertionError()` | Rule has been stable since v0.0.67Automatic fix available |
| B012 | jump-statement-in-finally | `{name}`inside`finally`blocks cause exceptions to be silenced | Rule has been stable since v0.0.116 |
| B013 | redundant-tuple-in-exception-handler | A length-one tuple literal is redundant in exception handlers | Rule has been stable since v0.0.89Automatic fix available |
| B014 | duplicate-handler-exception | Exception handler with duplicate exception: `{name}` | Rule has been stable since v0.0.67Automatic fix available |
| B015 | useless-comparison | Pointless comparison. Did you mean to assign a value? Otherwise, prepend `assert`or remove it. | Rule has been stable since v0.0.102 |
| B016 | raise-literal | Cannot raise a literal. Did you intend to return it or raise an Exception? | Rule has been stable since v0.0.102 |
| B017 | assert-raises-exception | Do not assert blind exception: `{exception}` | Rule has been stable since v0.0.83 |
| B018 | useless-expression | Found useless expression. Either assign it to a variable or remove it. | Rule has been stable since v0.0.100 |
| B019 | cached-instance-method | Use of `functools.lru_cache`or`functools.cache`on methods can lead to memory leaks | Rule has been stable since v0.0.114 |
| B020 | loop-variable-overrides-iterator | Loop control variable `{name}`overrides iterable it iterates | Rule has been stable since v0.0.121 |
| B021 | f-string-docstring | f-string used as docstring. Python will interpret this as a joined string, rather than a docstring. | Rule has been stable since v0.0.116 |
| B022 | useless-contextlib-suppress | No arguments passed to `contextlib.suppress`. No exceptions will be suppressed and therefore this context manager is redundant | Rule has been stable since v0.0.118 |
| B023 | function-uses-loop-variable | Function definition does not bind loop variable `{name}` | Rule has been stable since v0.0.139 |
| B024 | abstract-base-class-without-abstract-method | `{name}`is an abstract base class, but it has no abstract methods or properties | Rule has been stable since v0.0.118 |
| B025 | duplicate-try-block-exception | try-except* block with duplicate exception `{name}` | Rule has been stable since v0.0.67 |
| B026 | star-arg-unpacking-after-keyword-arg | Star-arg unpacking after a keyword argument is strongly discouraged | Rule has been stable since v0.0.109 |
| B027 | empty-method-without-abstract-decorator | `{name}`is an empty method in an abstract base class, but has no abstract decorator | Rule has been stable since v0.0.118 |
| B028 | no-explicit-stacklevel | No explicit `stacklevel`keyword argument found | Rule has been stable since v0.0.257Automatic fix available |
| B029 | except-with-empty-tuple | Using `except* ():`with an empty tuple does not catch anything; add exceptions to handle | Rule has been stable since v0.0.250 |
| B030 | except-with-non-exception-classes | `except*`handlers should only be exception classes or tuples of exception classes | Rule has been stable since v0.0.255 |
| B031 | reuse-of-groupby-generator | Using the generator returned from `itertools.groupby()`more than once will do nothing on the second usage | Rule has been stable since v0.0.260 |
| B032 | unintentional-type-annotation | Possible unintentional type annotation (using `:`). Did you mean to assign (using`=`)? | Rule has been stable since v0.0.250 |
| B033 | duplicate-value | Sets should not contain duplicate item `{value}` | Rule has been stable since v0.0.271Automatic fix available |
| B034 | re-sub-positional-args | `{method}`should pass`{param_name}`and`flags`as keyword arguments to avoid confusion due to unintuitive argument positions | Rule has been stable since v0.0.278 |
| B035 | static-key-dict-comprehension | Dictionary comprehension uses static key: `{key}` | Rule has been stable since v0.2.0 |
| B039 | mutable-contextvar-default | Do not use mutable data structures for `ContextVar`defaults | Rule has been stable since 0.8.0 |
| B043 | del-attr-with-constant | Do not call `delattr`with a constant attribute value. It is not any safer than normal property deletion. | Rule has been in preview since 0.15.6Automatic fix available |
| B901 | return-in-generator | Using `yield`and`return {value}`in a generator function can lead to confusing behavior | Rule has been in preview since v0.4.8 |
| B903 | class-as-data-structure | Class could be dataclass or namedtuple | Rule has been in preview since 0.9.0 |
| B904 | raise-without-from-inside-except | Within an `except*`clause, raise exceptions with`raise ... from err`or`raise ... from None`to distinguish them from errors in exception handling | Rule has been stable since v0.0.138 |
| B905 | zip-without-explicit-strict | `zip()`without an explicit`strict=`parameter | Rule has been stable since v0.0.167Automatic fix available |
| B909 | loop-iterator-mutation | Mutation to loop iterable `{name}`during iteration | Rule has been in preview since v0.3.7 |
| B911 | batched-without-explicit-strict | `itertools.batched()`without an explicit`strict`parameter | Rule has been stable since 0.10.0 |
| B912 | map-without-explicit-strict | `map()`without an explicit`strict=`parameter | Rule has been stable since 0.15.0Automatic fix available |

## flake8-builtins (A)

For more, see flake8-builtins on PyPI.

For related settings, see flake8-builtins.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| A001 | builtin-variable-shadowing | Variable `{name}`is shadowing a Python builtin | Rule has been stable since v0.0.48 |
| A002 | builtin-argument-shadowing | Function argument `{name}`is shadowing a Python builtin | Rule has been stable since v0.0.48 |
| A003 | builtin-attribute-shadowing | Python builtin is shadowed by class attribute `{name}`from {row} | Rule has been stable since v0.0.48 |
| A004 | builtin-import-shadowing | Import `{name}`is shadowing a Python builtin | Rule has been stable since 0.8.0 |
| A005 | stdlib-module-shadowing | Module `{name}`shadows a Python standard-library module | Rule has been stable since 0.9.0 |
| A006 | builtin-lambda-argument-shadowing | Lambda argument `{name}`is shadowing a Python builtin | Rule has been stable since 0.9.0 |

## flake8-commas (COM)

For more, see flake8-commas on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| COM812 | missing-trailing-comma | Trailing comma missing | Rule has been stable since v0.0.223Automatic fix available |
| COM818 | trailing-comma-on-bare-tuple | Trailing comma on bare tuple prohibited | Rule has been stable since v0.0.223 |
| COM819 | prohibited-trailing-comma | Trailing comma prohibited | Rule has been stable since v0.0.223Automatic fix available |

## flake8-comprehensions (C4)

For more, see flake8-comprehensions on PyPI.

For related settings, see flake8-comprehensions.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| C400 | unnecessary-generator-list | Unnecessary generator (rewrite using `list()`) | Rule has been stable since v0.0.61Automatic fix available |
| C401 | unnecessary-generator-set | Unnecessary generator (rewrite using `set()`) | Rule has been stable since v0.0.61Automatic fix available |
| C402 | unnecessary-generator-dict | Unnecessary generator (rewrite as a dict comprehension) | Rule has been stable since v0.0.61Automatic fix available |
| C403 | unnecessary-list-comprehension-set | Unnecessary list comprehension (rewrite as a set comprehension) | Rule has been stable since v0.0.58Automatic fix available |
| C404 | unnecessary-list-comprehension-dict | Unnecessary list comprehension (rewrite as a dict comprehension) | Rule has been stable since v0.0.58Automatic fix available |
| C405 | unnecessary-literal-set | Unnecessary {kind} literal (rewrite as a set literal) | Rule has been stable since v0.0.61Automatic fix available |
| C406 | unnecessary-literal-dict | Unnecessary {obj_type} literal (rewrite as a dict literal) | Rule has been stable since v0.0.61Automatic fix available |
| C408 | unnecessary-collection-call | Unnecessary `{kind}()`call (rewrite as a literal) | Rule has been stable since v0.0.61Automatic fix available |
| C409 | unnecessary-literal-within-tuple-call | Unnecessary list literal passed to `tuple()`(rewrite as a tuple literal) | Rule has been stable since v0.0.66Automatic fix available |
| C410 | unnecessary-literal-within-list-call | Unnecessary list literal passed to `list()`(remove the outer call to`list()`) | Rule has been stable since v0.0.66Automatic fix available |
| C411 | unnecessary-list-call | Unnecessary `list()`call (remove the outer call to`list()`) | Rule has been stable since v0.0.73Automatic fix available |
| C413 | unnecessary-call-around-sorted | Unnecessary `{func}()`call around`sorted()` | Rule has been stable since v0.0.73Automatic fix available |
| C414 | unnecessary-double-cast-or-process | Unnecessary `{inner}()`call within`{outer}()` | Rule has been stable since v0.0.70Automatic fix available |
| C415 | unnecessary-subscript-reversal | Unnecessary subscript reversal of iterable within `{func}()` | Rule has been stable since v0.0.64 |
| C416 | unnecessary-comprehension | Unnecessary {kind} comprehension (rewrite using `{kind}()`) | Rule has been stable since v0.0.73Automatic fix available |
| C417 | unnecessary-map | Unnecessary `map()`usage (rewrite using a {object_type}) | Rule has been stable since v0.0.74Automatic fix available |
| C418 | unnecessary-literal-within-dict-call | Unnecessary dict {kind} passed to `dict()`(remove the outer call to`dict()`) | Rule has been stable since v0.0.262Automatic fix available |
| C419 | unnecessary-comprehension-in-call | Unnecessary list comprehension | Rule has been stable since v0.0.262Automatic fix available |
| C420 | unnecessary-dict-comprehension-for-iterable | Unnecessary dict comprehension for iterable; use `dict.fromkeys`instead | Rule has been stable since 0.10.0Automatic fix available |

## flake8-copyright (CPY)

For more, see flake8-copyright on PyPI.

For related settings, see flake8-copyright.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| CPY001 | missing-copyright-notice | Missing copyright notice at top of file | Rule has been stable since 0.16.0 |

## flake8-datetimez (DTZ)

For more, see flake8-datetimez on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| DTZ001 | call-datetime-without-tzinfo | `datetime.datetime()`called without a`tzinfo`argument | Rule has been stable since v0.0.188 |
| DTZ002 | call-datetime-today | `datetime.datetime.today()`used | Rule has been stable since v0.0.188 |
| DTZ003 | call-datetime-utcnow | `datetime.datetime.utcnow()`used | Rule has been stable since v0.0.188 |
| DTZ004 | call-datetime-utcfromtimestamp | `datetime.datetime.utcfromtimestamp()`used | Rule has been stable since v0.0.188 |
| DTZ005 | call-datetime-now-without-tzinfo | `datetime.datetime.now()`called without a`tz`argument | Rule has been stable since v0.0.188 |
| DTZ006 | call-datetime-fromtimestamp | `datetime.datetime.fromtimestamp()`called without a`tz`argument | Rule has been stable since v0.0.188 |
| DTZ007 | call-datetime-strptime-without-zone | Naive datetime constructed using `datetime.datetime.strptime()`without %z | Rule has been stable since v0.0.188 |
| DTZ011 | call-date-today | `datetime.date.today()`used | Rule has been stable since v0.0.188 |
| DTZ012 | call-date-fromtimestamp | `datetime.date.fromtimestamp()`used | Rule has been stable since v0.0.188 |
| DTZ901 | datetime-min-max | Use of `datetime.datetime.{min_max}`without timezone information | Rule has been stable since 0.10.0 |

## flake8-debugger (T10)

For more, see flake8-debugger on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| T100 | debugger | Trace found: `{name}`used | Rule has been stable since v0.0.141 |

## flake8-django (DJ)

For more, see flake8-django on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| DJ001 | django-nullable-model-string-field | Avoid using `null=True`on string-based fields such as`{field_name}` | Rule has been stable since v0.0.246 |
| DJ003 | django-locals-in-render-function | Avoid passing `locals()`as context to a`render`function | Rule has been stable since v0.0.253 |
| DJ006 | django-exclude-with-model-form | Do not use `exclude`with`ModelForm`, use`fields`instead | Rule has been stable since v0.0.253 |
| DJ007 | django-all-with-model-form | Do not use `__all__`with`ModelForm`, use`fields`instead | Rule has been stable since v0.0.253 |
| DJ008 | django-model-without-dunder-str | Model does not define `__str__`method | Rule has been stable since v0.0.246 |
| DJ012 | django-unordered-body-content-in-model | Order of model's inner classes, methods, and fields does not follow the Django Style Guide: {element_type} should come before {prev_element_type} | Rule has been stable since v0.0.258 |
| DJ013 | django-non-leading-receiver-decorator | `@receiver`decorator must be on top of all the other decorators | Rule has been stable since v0.0.246 |

## flake8-errmsg (EM)

For more, see flake8-errmsg on PyPI.

For related settings, see flake8-errmsg.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| EM101 | raw-string-in-exception | Exception must not use a string literal, assign to variable first | Rule has been stable since v0.0.183Automatic fix available |
| EM102 | f-string-in-exception | Exception must not use an f-string literal, assign to variable first | Rule has been stable since v0.0.183Automatic fix available |
| EM103 | dot-format-in-exception | Exception must not use a `.format()`string directly, assign to variable first | Rule has been stable since v0.0.183Automatic fix available |

## flake8-executable (EXE)

For more, see flake8-executable on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| EXE001 | shebang-not-executable | Shebang is present but file is not executable | Rule has been stable since v0.0.233 |
| EXE002 | shebang-missing-executable-file | The file is executable but no shebang is present | Rule has been stable since v0.0.233 |
| EXE003 | shebang-missing-python | Shebang should contain `python`,`pytest`, or`uv run` | Rule has been stable since v0.0.229 |
| EXE004 | shebang-leading-whitespace | Avoid whitespace before shebang | Rule has been stable since v0.0.229Automatic fix available |
| EXE005 | shebang-not-first-line | Shebang should be at the beginning of the file | Rule has been stable since v0.0.229 |

## flake8-fixme (FIX)

For more, see flake8-fixme on GitHub.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| FIX001 | line-contains-fixme | Line contains FIXME, consider resolving the issue | Rule has been stable since v0.0.272 |
| FIX002 | line-contains-todo | Line contains TODO, consider resolving the issue | Rule has been stable since v0.0.272 |
| FIX003 | line-contains-xxx | Line contains XXX, consider resolving the issue | Rule has been stable since v0.0.272 |
| FIX004 | line-contains-hack | Line contains HACK, consider resolving the issue | Rule has been stable since v0.0.272 |

## flake8-future-annotations (FA)

For more, see flake8-future-annotations on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| FA100 | future-rewritable-type-annotation | Add `from __future__ import annotations`to simplify`{name}` | Rule has been stable since v0.0.269Automatic fix available |
| FA102 | future-required-type-annotation | Missing `from __future__ import annotations`, but uses {reason} | Rule has been stable since v0.0.271Automatic fix available |

## flake8-gettext (INT)

For more, see flake8-gettext on PyPI.

For related settings, see flake8-gettext.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| INT001 | f-string-in-get-text-func-call | f-string in plural argument is resolved before function call | Rule has been stable since v0.0.260 |
| INT002 | format-in-get-text-func-call | `format`method in plural argument is resolved before function call | Rule has been stable since v0.0.260 |
| INT003 | printf-in-get-text-func-call | printf-style format in plural argument is resolved before function call | Rule has been stable since v0.0.260 |

## flake8-implicit-str-concat (ISC)

For more, see flake8-implicit-str-concat on PyPI.

For related settings, see flake8-implicit-str-concat.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| ISC001 | single-line-implicit-string-concatenation | Implicitly concatenated string literals on one line | Rule has been stable since v0.0.201Automatic fix available |
| ISC002 | multi-line-implicit-string-concatenation | Implicitly concatenated string literals over multiple lines | Rule has been stable since v0.0.201 |
| ISC003 | explicit-string-concatenation | Explicitly concatenated string should be implicitly concatenated | Rule has been stable since v0.0.201Automatic fix available |
| ISC004 | implicit-string-concatenation-in-collection-literal | Unparenthesized implicit string concatenation in collection | Rule has been stable since 0.16.0Automatic fix available |

## flake8-import-conventions (ICN)

For more, see flake8-import-conventions on GitHub.

For related settings, see flake8-import-conventions.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| ICN001 | unconventional-import-alias | `{name}`should be imported as`{asname}` | Rule has been stable since v0.0.166Automatic fix available |
| ICN002 | banned-import-alias | `{name}`should not be imported as`{asname}` | Rule has been stable since v0.0.262 |
| ICN003 | banned-import-from | Members of `{name}`should not be imported explicitly | Rule has been stable since v0.0.263 |

## flake8-logging (LOG)

For more, see flake8-logging on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| LOG001 | direct-logger-instantiation | Use `logging.getLogger()`to instantiate loggers | Rule has been stable since v0.2.0Automatic fix available |
| LOG002 | invalid-get-logger-argument | Use `__name__`with`logging.getLogger()` | Rule has been stable since v0.2.0Automatic fix available |
| LOG004 | log-exception-outside-except-handler | `.exception()`call outside exception handlers | Rule has been stable since 0.16.0Automatic fix available |
| LOG007 | exception-without-exc-info | Use of `logging.exception`with falsy`exc_info` | Rule has been stable since v0.2.0 |
| LOG009 | undocumented-warn | Use of undocumented `logging.WARN`constant | Rule has been stable since v0.2.0Automatic fix available |
| LOG014 | exc-info-outside-except-handler | `exc_info=`outside exception handlers | Rule has been stable since 0.12.0Automatic fix available |
| LOG015 | root-logger-call | `{}()`call on root logger | Rule has been stable since 0.10.0 |

## flake8-logging-format (G)

For more, see flake8-logging-format on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| G001 | logging-string-format | Logging statement uses `str.format` | Rule has been stable since v0.0.236 |
| G002 | logging-percent-format | Logging statement uses `%` | Rule has been stable since v0.0.236 |
| G003 | logging-string-concat | Logging statement uses `+` | Rule has been stable since v0.0.236 |
| G004 | logging-f-string | Logging statement uses f-string | Rule has been stable since v0.0.236Automatic fix available |
| G010 | logging-warn | Logging statement uses `warn`instead of`warning` | Rule has been stable since v0.0.236Automatic fix available |
| G101 | logging-extra-attr-clash | Logging statement uses an `extra`field that clashes with a`LogRecord`field:`{key}` | Rule has been stable since v0.0.236 |
| G201 | logging-exc-info | Logging `.exception(...)`should be used instead of`.error(..., exc_info=True)` | Rule has been stable since v0.0.236 |
| G202 | logging-redundant-exc-info | Logging statement has redundant `exc_info` | Rule has been stable since v0.0.236 |

## flake8-no-pep420 (INP)

For more, see flake8-no-pep420 on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| INP001 | implicit-namespace-package | File `{filename}`is part of an implicit namespace package. Add an`__init__.py`. | Rule has been stable since v0.0.225 |

## flake8-pie (PIE)

For more, see flake8-pie on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| PIE790 | unnecessary-placeholder | Unnecessary `pass`statement | Rule has been stable since v0.0.208Automatic fix available |
| PIE794 | duplicate-class-field-definition | Class field `{name}`is defined multiple times | Rule has been stable since v0.0.208Automatic fix available |
| PIE796 | non-unique-enums | Enum contains duplicate value: `{value}` | Rule has been stable since v0.0.224 |
| PIE800 | unnecessary-spread | Unnecessary spread `**` | Rule has been stable since v0.0.231Automatic fix available |
| PIE804 | unnecessary-dict-kwargs | Unnecessary `dict`kwargs | Rule has been stable since v0.0.231Automatic fix available |
| PIE807 | reimplemented-container-builtin | Prefer `{container}`over useless lambda | Rule has been stable since v0.0.208Automatic fix available |
| PIE808 | unnecessary-range-start | Unnecessary `start`argument in`range` | Rule has been stable since v0.0.286Automatic fix available |
| PIE810 | multiple-starts-ends-with | Call `{attr}`once with a`tuple` | Rule has been stable since v0.0.243Automatic fix available |

## flake8-print (T20)

For more, see flake8-print on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| T201 | `print`found | Rule has been stable since v0.0.57Automatic fix available | |
| T203 | p-print | `pprint`found | Rule has been stable since v0.0.57Automatic fix available |

## flake8-pyi (PYI)

For more, see flake8-pyi on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| PYI001 | unprefixed-type-param | Name of private `{kind}`must start with`_` | Rule has been stable since v0.0.245 |
| PYI002 | complex-if-statement-in-stub | `if`test must be a simple comparison against`sys.platform`or`sys.version_info` | Rule has been stable since v0.0.276 |
| PYI003 | unrecognized-version-info-check | Unrecognized `sys.version_info`check | Rule has been stable since v0.0.276 |
| PYI004 | patch-version-comparison | Version comparison must use only major and minor version | Rule has been stable since v0.0.276 |
| PYI005 | wrong-tuple-length-version-comparison | Version comparison must be against a length-{expected_length} tuple | Rule has been stable since v0.0.276 |
| PYI006 | bad-version-info-comparison | Use `<`or`>=`for`sys.version_info`comparisons | Rule has been stable since v0.0.254 |
| PYI007 | unrecognized-platform-check | Unrecognized `sys.platform`check | Rule has been stable since v0.0.246 |
| PYI008 | unrecognized-platform-name | Unrecognized platform `{platform}` | Rule has been stable since v0.0.246 |
| PYI009 | pass-statement-stub-body | Empty body should contain `...`, not`pass` | Rule has been stable since v0.0.253Automatic fix available |
| PYI010 | non-empty-stub-body | Function body must contain only `...` | Rule has been stable since v0.0.253Automatic fix available |
| PYI011 | typed-argument-default-in-stub | Only simple default values allowed for typed arguments | Rule has been stable since v0.0.253Automatic fix available |
| PYI012 | pass-in-class-body | Class body must not contain `pass` | Rule has been stable since v0.0.260Automatic fix available |
| PYI013 | ellipsis-in-non-empty-class-body | Non-empty class body must not contain `...` | Rule has been stable since v0.0.270Automatic fix available |
| PYI014 | argument-default-in-stub | Only simple default values allowed for arguments | Rule has been stable since v0.0.253Automatic fix available |
| PYI015 | assignment-default-in-stub | Only simple default values allowed for assignments | Rule has been stable since v0.0.260Automatic fix available |
| PYI016 | duplicate-union-member | Duplicate union member `{}` | Rule has been stable since v0.0.262Automatic fix available |
| PYI017 | complex-assignment-in-stub | Stubs should not contain assignments to attributes or multiple targets | Rule has been stable since v0.0.279 |
| PYI018 | unused-private-type-var | Private {type_var_like_kind} `{type_var_like_name}`is never used | Rule has been stable since v0.0.281Automatic fix available |
| PYI019 | custom-type-var-for-self | Use `Self`instead of custom TypeVar`{}` | Rule has been stable since v0.0.283Automatic fix available |
| PYI020 | quoted-annotation-in-stub | Quoted annotations should not be included in stubs | Rule has been stable since v0.0.265Automatic fix available |
| PYI021 | docstring-in-stub | Docstrings should not be included in stubs | Rule has been stable since v0.0.253Automatic fix available |
| PYI024 | collections-named-tuple | Use `typing.NamedTuple`instead of`collections.namedtuple` | Rule has been stable since v0.0.271 |
| PYI025 | unaliased-collections-abc-set-import | Use `from collections.abc import Set as AbstractSet`to avoid confusion with the`set`builtin | Rule has been stable since v0.0.271Automatic fix available |
| PYI026 | type-alias-without-annotation | Use `{module}.TypeAlias`for type alias, e.g.,`{name}: TypeAlias = {value}` | Rule has been stable since v0.0.279Automatic fix available |
| PYI029 | str-or-repr-defined-in-stub | Defining `{name}`in a stub is almost always redundant | Rule has been stable since v0.0.271Automatic fix available |
| PYI030 | unnecessary-literal-union | Multiple literal members in a union. Use a single literal, e.g. `Literal[{}]` | Rule has been stable since v0.0.278Automatic fix available |
| PYI032 | any-eq-ne-annotation | Prefer `object`to`Any`for the second parameter to`{method_name}` | Rule has been stable since v0.0.271Automatic fix available |
| PYI033 | legacy-type-comment | Don't use type comments | Rule has been stable since v0.0.254 |
| PYI034 | non-self-return-type | `__new__`methods usually return`self`at runtime | Rule has been stable since v0.0.271Automatic fix available |
| PYI035 | unassigned-special-variable-in-stub | `{name}`in a stub file must have a value, as it has the same semantics as`{name}`at runtime | Rule has been stable since v0.0.271 |
| PYI036 | bad-exit-annotation | Star-args in `{method_name}`should be annotated with`object` | Rule has been stable since v0.0.279Automatic fix available |
| PYI041 | redundant-numeric-union | Use `{supertype}`instead of`{subtype} | {supertype}` | Rule has been stable since v0.0.279Automatic fix available |
| PYI042 | snake-case-type-alias | Type alias `{name}`should be CamelCase | Rule has been stable since v0.0.265 |
| PYI043 | t-suffixed-type-alias | Private type alias `{name}`should not be suffixed with`T`(the`T`suffix implies that an object is a`TypeVar`) | Rule has been stable since v0.0.265 |
| PYI044 | future-annotations-in-stub | `from __future__ import annotations`has no effect in stub files, since type checkers automatically treat stubs as having those semantics | Rule has been stable since v0.0.273Automatic fix available |
| PYI045 | iter-method-return-iterable | `__aiter__`methods should return an`AsyncIterator`, not an`AsyncIterable` | Rule has been stable since v0.0.271 |
| PYI046 | unused-private-protocol | Private protocol `{name}`is never used | Rule has been stable since v0.0.281 |
| PYI047 | unused-private-type-alias | Private TypeAlias `{name}`is never used | Rule has been stable since v0.0.281 |
| PYI048 | stub-body-multiple-statements | Function body must contain exactly one statement | Rule has been stable since v0.0.271 |
| PYI049 | unused-private-typed-dict | Private TypedDict `{name}`is never used | Rule has been stable since v0.0.281 |
| PYI050 | no-return-argument-annotation-in-stub | Prefer `{module}.Never`over`NoReturn`for argument annotations | Rule has been stable since v0.0.272 |
| PYI051 | redundant-literal-union | `Literal[{literal}]`is redundant in a union with`{builtin_type}` | Rule has been stable since v0.0.283 |
| PYI052 | unannotated-assignment-in-stub | Need type annotation for `{name}` | Rule has been stable since v0.0.269 |
| PYI053 | string-or-bytes-too-long | String and bytes literals longer than 50 characters are not permitted | Rule has been stable since v0.0.271Automatic fix available |
| PYI054 | numeric-literal-too-long | Numeric literals with a string representation longer than ten characters are not permitted | Rule has been stable since v0.0.271Automatic fix available |
| PYI055 | unnecessary-type-union | Multiple `type`members in a union. Combine them into one, e.g.,`type[{union_str}]`. | Rule has been stable since v0.0.283Automatic fix available |
| PYI056 | unsupported-method-call-on-all | Calling `.{name}()`on`__all__`may not be supported by all type checkers (use`+=`instead) | Rule has been stable since v0.0.281 |
| PYI057 | byte-string-usage | Do not use `{origin}.ByteString`, which has unclear semantics and is deprecated | Rule has been stable since 0.6.0 |
| PYI058 | generator-return-from-iter-method | Use `{return_type}`as the return value for simple`{method}`methods | Rule has been stable since v0.2.0Automatic fix available |
| PYI059 | generic-not-last-base-class | `Generic[]`should always be the last base class | Rule has been stable since 0.13.0Automatic fix available |
| PYI061 | redundant-none-literal | Use `None`rather than`Literal[None]` | Rule has been stable since 0.13.0Automatic fix available |
| PYI062 | duplicate-literal-member | Duplicate literal member `{}` | Rule has been stable since 0.6.0Automatic fix available |
| PYI063 | pep484-style-positional-only-parameter | Use PEP 570 syntax for positional-only parameters | Rule has been stable since 0.8.0 |
| PYI064 | redundant-final-literal | `Final[Literal[{literal}]]`can be replaced with a bare`Final` | Rule has been stable since 0.8.0Automatic fix available |
| PYI066 | bad-version-info-order | Put branches for newer Python versions first when branching on `sys.version_info`comparisons | Rule has been stable since 0.8.0 |

## flake8-pytest-style (PT)

For more, see flake8-pytest-style on PyPI.

For related settings, see flake8-pytest-style.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| PT001 | pytest-fixture-incorrect-parentheses-style | Use `@pytest.fixture{expected}`over`@pytest.fixture{actual}` | Rule has been stable since v0.0.208Automatic fix available |
| PT002 | pytest-fixture-positional-args | Configuration for fixture `{function}`specified via positional args, use kwargs | Rule has been stable since v0.0.208 |
| PT003 | pytest-extraneous-scope-function | `scope='function'`is implied in`@pytest.fixture()` | Rule has been stable since v0.0.208Automatic fix available |
| PT004 | pytest-missing-fixture-name-underscore | Fixture `{function}`does not return anything, add leading underscore | Rule was removed in 0.8.0 |
| PT005 | pytest-incorrect-fixture-name-underscore | Fixture `{function}`returns a value, remove leading underscore | Rule was removed in 0.8.0 |
| PT006 | pytest-parametrize-names-wrong-type | Wrong type passed to first argument of `pytest.mark.parametrize`; expected {expected_string} | Rule has been stable since v0.0.208Automatic fix available |
| PT007 | pytest-parametrize-values-wrong-type | Wrong values type in `pytest.mark.parametrize`expected`{values}` | Rule has been stable since v0.0.208Automatic fix available |
| PT008 | pytest-patch-with-lambda | Use `return_value=`instead of patching with`lambda` | Rule has been stable since v0.0.208 |
| PT009 | pytest-unittest-assertion | Use a regular `assert`instead of unittest-style`{assertion}` | Rule has been stable since v0.0.208Automatic fix available |
| PT010 | pytest-raises-without-exception | Set the expected exception in `pytest.raises()` | Rule has been stable since v0.0.208 |
| PT011 | pytest-raises-too-broad | `pytest.raises({exception})`is too broad, set the`match`parameter or use a more specific exception | Rule has been stable since v0.0.208 |
| PT012 | pytest-raises-with-multiple-statements | `pytest.raises()`block should contain a single simple statement | Rule has been stable since v0.0.208 |
| PT013 | pytest-incorrect-pytest-import | Incorrect import of `pytest`; use`import pytest`instead | Rule has been stable since v0.0.208 |
| PT014 | pytest-duplicate-parametrize-test-cases | Duplicate of test case at index {index} in `pytest.mark.parametrize` | Rule has been stable since v0.0.285Automatic fix available |
| PT015 | pytest-assert-always-false | Assertion always fails, replace with `pytest.fail()` | Rule has been stable since v0.0.208 |
| PT016 | pytest-fail-without-message | No message passed to `pytest.fail()` | Rule has been stable since v0.0.208 |
| PT017 | pytest-assert-in-except | Found assertion on exception `{name}`in`except`block, use`pytest.raises()`instead | Rule has been stable since v0.0.208 |
| PT018 | pytest-composite-assertion | Assertion should be broken down into multiple parts | Rule has been stable since v0.0.208Automatic fix available |
| PT019 | pytest-fixture-param-without-value | Fixture `{name}`without value is injected as parameter, use`@pytest.mark.usefixtures`instead | Rule has been stable since v0.0.208 |
| PT020 | pytest-deprecated-yield-fixture | `@pytest.yield_fixture`is deprecated, use`@pytest.fixture` | Rule has been stable since v0.0.208 |
| PT021 | pytest-fixture-finalizer-callback | Use `yield`instead of`request.addfinalizer` | Rule has been stable since v0.0.208 |
| PT022 | pytest-useless-yield-fixture | No teardown in fixture `{name}`, use`return`instead of`yield` | Rule has been stable since v0.0.208Automatic fix available |
| PT023 | pytest-incorrect-mark-parentheses-style | Use `@pytest.mark.{mark_name}{expected_parens}`over`@pytest.mark.{mark_name}{actual_parens}` | Rule has been stable since v0.0.208Automatic fix available |
| PT024 | pytest-unnecessary-asyncio-mark-on-fixture | `pytest.mark.asyncio`is unnecessary for fixtures | Rule has been stable since v0.0.208Automatic fix available |
| PT025 | pytest-erroneous-use-fixtures-on-fixture | `pytest.mark.usefixtures`has no effect on fixtures | Rule has been stable since v0.0.208Automatic fix available |
| PT026 | pytest-use-fixtures-without-parameters | Useless `pytest.mark.usefixtures`without parameters | Rule has been stable since v0.0.208Automatic fix available |
| PT027 | pytest-unittest-raises-assertion | Use `pytest.raises`instead of unittest-style`{assertion}` | Rule has been stable since v0.0.285Automatic fix available |
| PT028 | pytest-parameter-with-default-argument | Test function parameter `{}`has default argument | Rule has been stable since 0.12.0 |
| PT029 | pytest-warns-without-warning | Set the expected warning in `pytest.warns()` | Rule has been in preview since 0.9.2 |
| PT030 | pytest-warns-too-broad | `pytest.warns({warning})`is too broad, set the`match`parameter or use a more specific warning | Rule has been stable since 0.12.0 |
| PT031 | pytest-warns-with-multiple-statements | `pytest.warns()`block should contain a single simple statement | Rule has been stable since 0.12.0 |

## flake8-quotes (Q)

For more, see flake8-quotes on PyPI.

For related settings, see flake8-quotes.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| Q000 | bad-quotes-inline-string | Single quotes found but double quotes preferred | Rule has been stable since v0.0.88Automatic fix available |
| Q001 | bad-quotes-multiline-string | Single quote multiline found but double quotes preferred | Rule has been stable since v0.0.88Automatic fix available |
| Q002 | bad-quotes-docstring | Single quote docstring found but double quotes preferred | Rule has been stable since v0.0.88Automatic fix available |
| Q003 | avoidable-escaped-quote | Change outer quotes to avoid escaping inner quotes | Rule has been stable since v0.0.88Automatic fix available |
| Q004 | unnecessary-escaped-quote | Unnecessary escape on inner quote character | Rule has been stable since v0.2.0Automatic fix available |

## flake8-raise (RSE)

For more, see flake8-raise on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| RSE102 | unnecessary-paren-on-raise-exception | Unnecessary parentheses on raised exception | Rule has been stable since v0.0.239Automatic fix available |

## flake8-return (RET)

For more, see flake8-return on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| RET501 | unnecessary-return-none | Do not explicitly `return None`in function if it is the only possible return value | Rule has been stable since v0.0.154Automatic fix available |
| RET502 | implicit-return-value | Do not implicitly `return None`in function able to return non-`None`value | Rule has been stable since v0.0.154Automatic fix available |
| RET503 | implicit-return | Missing explicit `return`at the end of function able to return non-`None`value | Rule has been stable since v0.0.154Automatic fix available |
| RET504 | unnecessary-assign | Unnecessary assignment to `{name}`before`return`statement | Rule has been stable since v0.0.154Automatic fix available |
| RET505 | superfluous-else-return | Unnecessary `{branch}`after`return`statement | Rule has been stable since v0.0.154Automatic fix available |
| RET506 | superfluous-else-raise | Unnecessary `{branch}`after`raise`statement | Rule has been stable since v0.0.154Automatic fix available |
| RET507 | superfluous-else-continue | Unnecessary `{branch}`after`continue`statement | Rule has been stable since v0.0.154Automatic fix available |
| RET508 | superfluous-else-break | Unnecessary `{branch}`after`break`statement | Rule has been stable since v0.0.154Automatic fix available |

## flake8-self (SLF)

For more, see flake8-self on PyPI.

For related settings, see flake8-self.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| SLF001 | private-member-access | Private member accessed: `{access}` | Rule has been stable since v0.0.240 |

## flake8-simplify (SIM)

For more, see flake8-simplify on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| SIM101 | duplicate-isinstance-call | Multiple `isinstance`calls for`{name}`, merge into a single call | Rule has been stable since v0.0.212Automatic fix available |
| SIM102 | collapsible-if | Use a single `if`statement instead of nested`if`statements | Rule has been stable since v0.0.211Automatic fix available |
| SIM103 | needless-bool | Return the condition `{condition}`directly | Rule has been stable since v0.0.214Automatic fix available |
| SIM105 | suppressible-exception | Use `contextlib.suppress({exception})`instead of`try`-`except`-`pass` | Rule has been stable since v0.0.211Automatic fix available |
| SIM107 | return-in-try-except-finally | Don't use `return`in`try`-`except`and`finally` | Rule has been stable since v0.0.211 |
| SIM108 | if-else-block-instead-of-if-exp | Use ternary operator `{contents}`instead of`if`-`else`-block | Rule has been stable since v0.0.213Automatic fix available |
| SIM109 | compare-with-tuple | Use `{replacement}`instead of multiple equality comparisons | Rule has been stable since v0.0.213Automatic fix available |
| SIM110 | reimplemented-builtin | Use `{replacement}`instead of`for`loop | Rule has been stable since v0.0.211Automatic fix available |
| SIM112 | uncapitalized-environment-variables | Use capitalized environment variable `{expected}`instead of`{actual}` | Rule has been stable since v0.0.218Automatic fix available |
| SIM113 | enumerate-for-loop | Use `enumerate()`for index variable`{index}`in`for`loop | Rule has been stable since v0.2.0 |
| SIM114 | if-with-same-arms | Combine `if`branches using logical`or`operator | Rule has been stable since v0.0.246Automatic fix available |
| SIM115 | open-file-with-context-handler | Use a context manager for opening files | Rule has been stable since v0.0.219 |
| SIM116 | if-else-block-instead-of-dict-lookup | Use a dictionary instead of consecutive `if`statements | Rule has been stable since v0.0.250 |
| SIM117 | multiple-with-statements | Use a single `with`statement with multiple contexts instead of nested`with`statements | Rule has been stable since v0.0.211Automatic fix available |
| SIM118 | in-dict-keys | Use `key {operator} dict`instead of`key {operator} dict.keys()` | Rule has been stable since v0.0.176Automatic fix available |
| SIM201 | negate-equal-op | Use `{left} != {right}`instead of`not {left} == {right}` | Rule has been stable since v0.0.213Automatic fix available |
| SIM202 | negate-not-equal-op | Use `{left} == {right}`instead of`not {left} != {right}` | Rule has been stable since v0.0.213Automatic fix available |
| SIM208 | double-negation | Use `{expr}`instead of`not (not {expr})` | Rule has been stable since v0.0.213Automatic fix available |
| SIM210 | if-expr-with-true-false | Remove unnecessary `True if ... else False` | Rule has been stable since v0.0.214Automatic fix available |
| SIM211 | if-expr-with-false-true | Use `not ...`instead of`False if ... else True` | Rule has been stable since v0.0.214Automatic fix available |
| SIM212 | if-expr-with-twisted-arms | Use `{expr_else} if {expr_else} else {expr_body}`instead of`{expr_body} if not {expr_else} else {expr_else}` | Rule has been stable since v0.0.214Automatic fix available |
| SIM220 | expr-and-not-expr | Use `False`instead of`{name} and not {name}` | Rule has been stable since v0.0.211Automatic fix available |
| SIM221 | expr-or-not-expr | Use `True`instead of`{name} or not {name}` | Rule has been stable since v0.0.211Automatic fix available |
| SIM222 | expr-or-true | Use `{expr}`instead of`{replaced}` | Rule has been stable since v0.0.208Automatic fix available |
| SIM223 | expr-and-false | Use `{expr}`instead of`{replaced}` | Rule has been stable since v0.0.208Automatic fix available |
| SIM300 | yoda-conditions | Yoda condition detected | Rule has been stable since v0.0.207Automatic fix available |
| SIM401 | if-else-block-instead-of-dict-get | Use `{contents}`instead of an`if`block | Rule has been stable since v0.0.219Automatic fix available |
| SIM905 | split-static-string | Consider using a list literal instead of `str.{}` | Rule has been stable since 0.10.0Automatic fix available |
| SIM910 | dict-get-with-none-default | Use `{expected}`instead of`{actual}` | Rule has been stable since v0.0.261Automatic fix available |
| SIM911 | zip-dict-keys-and-values | Use `{expected}`instead of`{actual}` | Rule has been stable since v0.2.0Automatic fix available |

## flake8-slots (SLOT)

For more, see flake8-slots on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| SLOT000 | no-slots-in-str-subclass | Subclasses of `str`should define`__slots__` | Rule has been stable since v0.0.273 |
| SLOT001 | no-slots-in-tuple-subclass | Subclasses of `tuple`should define`__slots__` | Rule has been stable since v0.0.273 |
| SLOT002 | no-slots-in-namedtuple-subclass | Subclasses of {namedtuple_kind} should define `__slots__` | Rule has been stable since v0.0.273 |

## flake8-tidy-imports (TID)

For more, see flake8-tidy-imports on PyPI.

For related settings, see flake8-tidy-imports.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| TID251 | banned-api | `{name}`is banned: {message} | Rule has been stable since v0.0.201 |
| TID252 | relative-imports | Prefer absolute imports over relative imports from parent modules | Rule has been stable since v0.0.169Automatic fix available |
| TID253 | banned-module-level-imports | `{name}`is banned at the module level | Rule has been stable since v0.0.285 |
| TID254 | lazy-import-mismatch | `{name}`should be imported lazily | Rule has been in preview since 0.15.6Automatic fix available |
| TID255 | lazy-import-immediately-resolved | Lazy import `{name}`is resolved immediately | Rule has been in preview since 0.15.13Automatic fix available |

## flake8-todos (TD)

For more, see flake8-todos on GitHub.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| TD001 | invalid-todo-tag | Invalid TODO tag: `{tag}` | Rule has been stable since v0.0.269 |
| TD002 | missing-todo-author | Missing author in TODO; try: `# TODO(<author_name>): ...`or`# TODO @<author_name>: ...` | Rule has been stable since v0.0.269 |
| TD003 | missing-todo-link | Missing issue link for this TODO | Rule has been stable since v0.0.269 |
| TD004 | missing-todo-colon | Missing colon in TODO | Rule has been stable since v0.0.269 |
| TD005 | missing-todo-description | Missing issue description after `TODO` | Rule has been stable since v0.0.269 |
| TD006 | invalid-todo-capitalization | Invalid TODO capitalization: `{tag}`should be`TODO` | Rule has been stable since v0.0.269Automatic fix available |
| TD007 | missing-space-after-todo-colon | Missing space after colon in TODO | Rule has been stable since v0.0.269 |

## flake8-type-checking (TC)

For more, see flake8-type-checking on PyPI.

For related settings, see flake8-type-checking.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| TC001 | typing-only-first-party-import | Move application import `{}`into a type-checking block | Rule has been stable since 0.8.0Automatic fix available |
| TC002 | typing-only-third-party-import | Move third-party import `{}`into a type-checking block | Rule has been stable since 0.8.0Automatic fix available |
| TC003 | typing-only-standard-library-import | Move standard library import `{}`into a type-checking block | Rule has been stable since 0.8.0Automatic fix available |
| TC004 | runtime-import-in-type-checking-block | Move import `{qualified_name}`out of type-checking block. Import is used for more than type hinting. | Rule has been stable since 0.8.0Automatic fix available |
| TC005 | empty-type-checking-block | Found empty type-checking block | Rule has been stable since 0.8.0Automatic fix available |
| TC006 | runtime-cast-value | Add quotes to type expression in `typing.cast()` | Rule has been stable since 0.10.0Automatic fix available |
| TC007 | unquoted-type-alias | Add quotes to type alias | Rule has been stable since 0.10.0Automatic fix available |
| TC008 | quoted-type-alias | Remove quotes from type alias | Rule has been in preview since 0.8.1Automatic fix available |
| TC010 | runtime-string-union | Invalid string member in `X | Y`-style union type | Rule has been stable since 0.8.0 |

## flake8-unused-arguments (ARG)

For more, see flake8-unused-arguments on PyPI.

For related settings, see flake8-unused-arguments.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| ARG001 | unused-function-argument | Unused function argument: `{name}` | Rule has been stable since v0.0.168 |
| ARG002 | unused-method-argument | Unused method argument: `{name}` | Rule has been stable since v0.0.168 |
| ARG003 | unused-class-method-argument | Unused class method argument: `{name}` | Rule has been stable since v0.0.168 |
| ARG004 | unused-static-method-argument | Unused static method argument: `{name}` | Rule has been stable since v0.0.168 |
| ARG005 | unused-lambda-argument | Unused lambda argument: `{name}` | Rule has been stable since v0.0.168 |

## flake8-use-pathlib (PTH)

For more, see flake8-use-pathlib on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| PTH100 | os-path-abspath | `os.path.abspath()`should be replaced by`Path.resolve()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH101 | os-chmod | `os.chmod()`should be replaced by`Path.chmod()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH102 | os-mkdir | `os.mkdir()`should be replaced by`Path.mkdir()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH103 | os-makedirs | `os.makedirs()`should be replaced by`Path.mkdir(parents=True)` | Rule has been stable since v0.0.231Automatic fix available |
| PTH104 | os-rename | `os.rename()`should be replaced by`Path.rename()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH105 | os-replace | `os.replace()`should be replaced by`Path.replace()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH106 | os-rmdir | `os.rmdir()`should be replaced by`Path.rmdir()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH107 | os-remove | `os.remove()`should be replaced by`Path.unlink()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH108 | os-unlink | `os.unlink()`should be replaced by`Path.unlink()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH109 | os-getcwd | `os.getcwd()`should be replaced by`Path.cwd()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH110 | os-path-exists | `os.path.exists()`should be replaced by`Path.exists()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH111 | os-path-expanduser | `os.path.expanduser()`should be replaced by`Path.expanduser()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH112 | os-path-isdir | `os.path.isdir()`should be replaced by`Path.is_dir()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH113 | os-path-isfile | `os.path.isfile()`should be replaced by`Path.is_file()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH114 | os-path-islink | `os.path.islink()`should be replaced by`Path.is_symlink()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH115 | os-readlink | `os.readlink()`should be replaced by`Path.readlink()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH116 | os-stat | `os.stat()`should be replaced by`Path.stat()`,`Path.owner()`, or`Path.group()` | Rule has been stable since v0.0.231 |
| PTH117 | os-path-isabs | `os.path.isabs()`should be replaced by`Path.is_absolute()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH118 | os-path-join | `os.{module}.join()`should be replaced by`Path`with`/`operator | Rule has been stable since v0.0.231 |
| PTH119 | os-path-basename | `os.path.basename()`should be replaced by`Path.name` | Rule has been stable since v0.0.231Automatic fix available |
| PTH120 | os-path-dirname | `os.path.dirname()`should be replaced by`Path.parent` | Rule has been stable since v0.0.231Automatic fix available |
| PTH121 | os-path-samefile | `os.path.samefile()`should be replaced by`Path.samefile()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH122 | os-path-splitext | `os.path.splitext()`should be replaced by`Path.suffix`,`Path.stem`, and`Path.parent` | Rule has been stable since v0.0.231 |
| PTH123 | builtin-open | `open()`should be replaced by`Path.open()` | Rule has been stable since v0.0.231Automatic fix available |
| PTH124 | py-path | `py.path`is in maintenance mode, use`pathlib`instead | Rule has been stable since v0.0.231 |
| PTH201 | path-constructor-current-directory | Do not pass the current directory explicitly to `Path` | Rule has been stable since v0.0.279Automatic fix available |
| PTH202 | os-path-getsize | `os.path.getsize`should be replaced by`Path.stat().st_size` | Rule has been stable since v0.0.279Automatic fix available |
| PTH203 | os-path-getatime | `os.path.getatime`should be replaced by`Path.stat().st_atime` | Rule has been stable since v0.0.279Automatic fix available |
| PTH204 | os-path-getmtime | `os.path.getmtime`should be replaced by`Path.stat().st_mtime` | Rule has been stable since v0.0.279Automatic fix available |
| PTH205 | os-path-getctime | `os.path.getctime`should be replaced by`Path.stat().st_ctime` | Rule has been stable since v0.0.279Automatic fix available |
| PTH206 | os-sep-split | Replace `.split(os.sep)`with`Path.parts` | Rule has been stable since v0.0.281 |
| PTH207 | glob | Replace `{function}`with`Path.glob`or`Path.rglob` | Rule has been stable since v0.0.281 |
| PTH208 | os-listdir | Use `pathlib.Path.iterdir()`instead. | Rule has been stable since 0.10.0 |
| PTH210 | invalid-pathlib-with-suffix | Invalid suffix passed to `.with_suffix()` | Rule has been stable since 0.10.0Automatic fix available |
| PTH211 | os-symlink | `os.symlink`should be replaced by`Path.symlink_to` | Rule has been stable since 0.13.0Automatic fix available |

## flynt (FLY)

For more, see flynt on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| FLY002 | static-join-to-f-string | Consider `{expression}`instead of string join | Rule has been stable since v0.0.266Automatic fix available |

## isort (I)

For more, see isort on PyPI.

For related settings, see isort.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| I001 | unsorted-imports | Import block is un-sorted or un-formatted | Rule has been stable since v0.0.110Automatic fix available |
| I002 | missing-required-import | Missing required import: `{name}` | Rule has been stable since v0.0.218Automatic fix available |

## mccabe (C90)

For more, see mccabe on PyPI.

For related settings, see mccabe.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| C901 | complex-structure | `{name}`is too complex ({complexity} > {max_complexity}) | Rule has been stable since v0.0.127 |

## NumPy-specific rules (NPY)

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| NPY001 | numpy-deprecated-type-alias | Type alias `np.{type_name}`is deprecated, replace with builtin type | Rule has been stable since v0.0.247Automatic fix available |
| NPY002 | numpy-legacy-random | Replace legacy `np.random.{method_name}`call with`np.random.Generator` | Rule has been stable since v0.0.248 |
| NPY003 | numpy-deprecated-function | `np.{existing}`is deprecated; use`np.{replacement}`instead | Rule has been stable since v0.0.276Automatic fix available |
| NPY201 | numpy2-deprecation | `np.{existing}`will be removed in NumPy 2.0. {migration_guide} | Rule has been stable since v0.2.0Automatic fix available |

## pandas-vet (PD)

For more, see pandas-vet on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| PD002 | pandas-use-of-inplace-argument | `inplace=True`should be avoided; it has inconsistent behavior | Rule has been stable since v0.0.188Automatic fix available |
| PD003 | pandas-use-of-dot-is-null | `.isna`is preferred to`.isnull`; functionality is equivalent | Rule has been stable since v0.0.188 |
| PD004 | pandas-use-of-dot-not-null | `.notna`is preferred to`.notnull`; functionality is equivalent | Rule has been stable since v0.0.188 |
| PD007 | pandas-use-of-dot-ix | `.ix`is deprecated; use more explicit`.loc`or`.iloc` | Rule has been stable since v0.0.188 |
| PD008 | pandas-use-of-dot-at | Use `.loc`instead of`.at`. If speed is important, use NumPy. | Rule has been stable since v0.0.188 |
| PD009 | pandas-use-of-dot-iat | Use `.iloc`instead of`.iat`. If speed is important, use NumPy. | Rule has been stable since v0.0.188 |
| PD010 | pandas-use-of-dot-pivot-or-unstack | `.pivot_table`is preferred to`.pivot`or`.unstack`; provides same functionality | Rule has been stable since v0.0.188 |
| PD011 | pandas-use-of-dot-values | Use `.to_numpy()`or`.array`instead of`.values` | Rule has been stable since v0.0.188 |
| PD012 | pandas-use-of-dot-read-table | Use `.read_csv`instead of`.read_table`to read CSV files | Rule has been stable since v0.0.188 |
| PD013 | pandas-use-of-dot-stack | `.melt`is preferred to`.stack`; provides same functionality | Rule has been stable since v0.0.188 |
| PD015 | pandas-use-of-pd-merge | Use `.merge`method instead of`pd.merge`function. They have equivalent functionality. | Rule has been stable since v0.0.188 |
| PD101 | pandas-nunique-constant-series-check | Using `series.nunique()`for checking that a series is constant is inefficient | Rule has been stable since v0.0.279 |
| PD901 | pandas-df-variable-name | Avoid using the generic variable name `df`for DataFrames | Rule was removed in 0.13.0 |

## pep8-naming (N)

For more, see pep8-naming on PyPI.

For related settings, see pep8-naming.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| N801 | invalid-class-name | Class name `{name}`should use CapWords convention | Rule has been stable since v0.0.77 |
| N802 | invalid-function-name | Function name `{name}`should be lowercase | Rule has been stable since v0.0.77 |
| N803 | invalid-argument-name | Argument name `{name}`should be lowercase | Rule has been stable since v0.0.77 |
| N804 | invalid-first-argument-name-for-class-method | First argument of a class method should be named `cls` | Rule has been stable since v0.0.77Automatic fix available |
| N805 | invalid-first-argument-name-for-method | First argument of a method should be named `self` | Rule has been stable since v0.0.77Automatic fix available |
| N806 | non-lowercase-variable-in-function | Variable `{name}`in function should be lowercase | Rule has been stable since v0.0.89 |
| N807 | dunder-function-name | Function name should not start and end with `__` | Rule has been stable since v0.0.82 |
| N811 | constant-imported-as-non-constant | Constant `{name}`imported as non-constant`{asname}` | Rule has been stable since v0.0.82 |
| N812 | lowercase-imported-as-non-lowercase | Lowercase `{name}`imported as non-lowercase`{asname}` | Rule has been stable since v0.0.82 |
| N813 | camelcase-imported-as-lowercase | Camelcase `{name}`imported as lowercase`{asname}` | Rule has been stable since v0.0.82 |
| N814 | camelcase-imported-as-constant | Camelcase `{name}`imported as constant`{asname}` | Rule has been stable since v0.0.82 |
| N815 | mixed-case-variable-in-class-scope | Variable `{name}`in class scope should not be mixedCase | Rule has been stable since v0.0.89 |
| N816 | mixed-case-variable-in-global-scope | Variable `{name}`in global scope should not be mixedCase | Rule has been stable since v0.0.89 |
| N817 | camelcase-imported-as-acronym | CamelCase `{name}`imported as acronym`{asname}` | Rule has been stable since v0.0.82 |
| N818 | error-suffix-on-exception-name | Exception name `{name}`should be named with an Error suffix | Rule has been stable since v0.0.89 |
| N999 | invalid-module-name | Invalid module name: '{name}' | Rule has been stable since v0.0.248 |

## Perflint (PERF)

For more, see Perflint on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| PERF101 | unnecessary-list-cast | Do not cast an iterable to `list`before iterating over it | Rule has been stable since v0.0.276Automatic fix available |
| PERF102 | incorrect-dict-iterator | When using only the {subset} of a dict use the `{subset}()`method | Rule has been stable since v0.0.273Automatic fix available |
| PERF203 | try-except-in-loop | `try`-`except`within a loop incurs performance overhead | Rule has been stable since v0.0.276 |
| PERF401 | manual-list-comprehension | Use {message_str} to create a transformed list | Rule has been stable since v0.0.276Automatic fix available |
| PERF402 | manual-list-copy | Use `list`or`list.copy`to create a copy of a list | Rule has been stable since v0.0.276 |
| PERF403 | manual-dict-comprehension | Use a dictionary comprehension instead of {modifier} for-loop | Rule has been stable since 0.5.0Automatic fix available |

## pycodestyle (E, W)

For more, see pycodestyle on PyPI.

For related settings, see pycodestyle.

### Error (E)

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| E101 | mixed-spaces-and-tabs | Indentation contains mixed spaces and tabs | Rule has been stable since v0.0.229 |
| E111 | indentation-with-invalid-multiple | Indentation is not a multiple of {indent_width} | Rule has been in preview since v0.0.269 |
| E112 | no-indented-block | Expected an indented block | Rule has been in preview since v0.0.269 |
| E113 | unexpected-indentation | Unexpected indentation | Rule has been in preview since v0.0.269 |
| E114 | indentation-with-invalid-multiple-comment | Indentation is not a multiple of {indent_width} (comment) | Rule has been in preview since v0.0.269 |
| E115 | no-indented-block-comment | Expected an indented block (comment) | Rule has been in preview since v0.0.269 |
| E116 | unexpected-indentation-comment | Unexpected indentation (comment) | Rule has been in preview since v0.0.269 |
| E117 | over-indented | Over-indented (comment) | Rule has been in preview since v0.0.269 |
| E201 | whitespace-after-open-bracket | Whitespace after '{symbol}' | Rule has been in preview since v0.0.269Automatic fix available |
| E202 | whitespace-before-close-bracket | Whitespace before '{symbol}' | Rule has been in preview since v0.0.269Automatic fix available |
| E203 | whitespace-before-punctuation | Whitespace before '{symbol}' | Rule has been in preview since v0.0.269Automatic fix available |
| E204 | whitespace-after-decorator | Whitespace after decorator | Rule has been in preview since 0.5.1Automatic fix available |
| E211 | whitespace-before-parameters | Whitespace before '{bracket}' | Rule has been in preview since v0.0.269Automatic fix available |
| E221 | multiple-spaces-before-operator | Multiple spaces before operator | Rule has been in preview since v0.0.269Automatic fix available |
| E222 | multiple-spaces-after-operator | Multiple spaces after operator | Rule has been in preview since v0.0.269Automatic fix available |
| E223 | tab-before-operator | Tab before operator | Rule has been in preview since v0.0.269Automatic fix available |
| E224 | tab-after-operator | Tab after operator | Rule has been in preview since v0.0.269Automatic fix available |
| E225 | missing-whitespace-around-operator | Missing whitespace around operator | Rule has been in preview since v0.0.269Automatic fix available |
| E226 | missing-whitespace-around-arithmetic-operator | Missing whitespace around arithmetic operator | Rule has been in preview since v0.0.269Automatic fix available |
| E227 | missing-whitespace-around-bitwise-or-shift-operator | Missing whitespace around bitwise or shift operator | Rule has been in preview since v0.0.269Automatic fix available |
| E228 | missing-whitespace-around-modulo-operator | Missing whitespace around modulo operator | Rule has been in preview since v0.0.269Automatic fix available |
| E231 | missing-whitespace | Missing whitespace after {} | Rule has been in preview since v0.0.269Automatic fix available |
| E241 | multiple-spaces-after-comma | Multiple spaces after comma | Rule has been in preview since v0.0.281Automatic fix available |
| E242 | tab-after-comma | Tab after comma | Rule has been in preview since v0.0.281Automatic fix available |
| E251 | unexpected-spaces-around-keyword-parameter-equals | Unexpected spaces around keyword / parameter equals | Rule has been in preview since v0.0.269Automatic fix available |
| E252 | missing-whitespace-around-parameter-equals | Missing whitespace around parameter equals | Rule has been in preview since v0.0.269Automatic fix available |
| E261 | too-few-spaces-before-inline-comment | Insert at least two spaces before an inline comment | Rule has been in preview since v0.0.269Automatic fix available |
| E262 | no-space-after-inline-comment | Inline comment should start with `#` | Rule has been in preview since v0.0.269Automatic fix available |
| E265 | no-space-after-block-comment | Block comment should start with `#` | Rule has been in preview since v0.0.269Automatic fix available |
| E266 | multiple-leading-hashes-for-block-comment | Too many leading `#`before block comment | Rule has been in preview since v0.0.269Automatic fix available |
| E271 | multiple-spaces-after-keyword | Multiple spaces after keyword | Rule has been in preview since v0.0.269Automatic fix available |
| E272 | multiple-spaces-before-keyword | Multiple spaces before keyword | Rule has been in preview since v0.0.269Automatic fix available |
| E273 | tab-after-keyword | Tab after keyword | Rule has been in preview since v0.0.269Automatic fix available |
| E274 | tab-before-keyword | Tab before keyword | Rule has been in preview since v0.0.269Automatic fix available |
| E275 | missing-whitespace-after-keyword | Missing whitespace after keyword | Rule has been in preview since v0.0.269Automatic fix available |
| E301 | blank-line-between-methods | Expected {BLANK_LINES_NESTED_LEVEL:?} blank line, found 0 | Rule has been in preview since v0.2.2Automatic fix available |
| E302 | blank-lines-top-level | Expected {expected_blank_lines:?} blank lines, found {actual_blank_lines} | Rule has been in preview since v0.2.2Automatic fix available |
| E303 | too-many-blank-lines | Too many blank lines ({actual_blank_lines}) | Rule has been in preview since v0.2.2Automatic fix available |
| E304 | blank-line-after-decorator | Blank lines found after function decorator ({lines}) | Rule has been in preview since v0.2.2Automatic fix available |
| E305 | blank-lines-after-function-or-class | Expected 2 blank lines after class or function definition, found ({blank_lines}) | Rule has been in preview since v0.2.2Automatic fix available |
| E306 | blank-lines-before-nested-definition | Expected 1 blank line before a nested definition, found 0 | Rule has been in preview since v0.2.2Automatic fix available |
| E401 | multiple-imports-on-one-line | Multiple imports on one line | Rule has been stable since v0.0.191Automatic fix available |
| E402 | module-import-not-at-top-of-file | Module level import not at top of cell | Rule has been stable since v0.0.28Automatic fix available |
| E501 | line-too-long | Line too long ({width} > {limit}) | Rule has been stable since v0.0.18 |
| E502 | redundant-backslash | Redundant backslash | Rule has been in preview since v0.3.3Automatic fix available |
| E701 | multiple-statements-on-one-line-colon | Multiple statements on one line (colon) | Rule has been stable since v0.0.245 |
| E702 | multiple-statements-on-one-line-semicolon | Multiple statements on one line (semicolon) | Rule has been stable since v0.0.245 |
| E703 | useless-semicolon | Statement ends with an unnecessary semicolon | Rule has been stable since v0.0.245Automatic fix available |
| E711 | none-comparison | Comparison to `None`should be`cond is None` | Rule has been stable since v0.0.28Automatic fix available |
| E712 | true-false-comparison | Avoid equality comparisons to `True`; use`{cond}:`for truth checks | Rule has been stable since v0.0.28Automatic fix available |
| E713 | not-in-test | Test for membership should be `not in` | Rule has been stable since v0.0.28Automatic fix available |
| E714 | not-is-test | Test for object identity should be `is not` | Rule has been stable since v0.0.28Automatic fix available |
| E721 | type-comparison | Use `is`and`is not`for type comparisons, or`isinstance()`for isinstance checks | Rule has been stable since v0.0.39 |
| E722 | bare-except | Do not use bare `except` | Rule has been stable since v0.0.36 |
| E731 | lambda-assignment | Do not assign a `lambda`expression, use a`def` | Rule has been stable since v0.0.28Automatic fix available |
| E741 | ambiguous-variable-name | Ambiguous variable name: `{name}` | Rule has been stable since v0.0.34 |
| E742 | ambiguous-class-name | Ambiguous class name: `{name}` | Rule has been stable since v0.0.35 |
| E743 | ambiguous-function-name | Ambiguous function name: `{name}` | Rule has been stable since v0.0.35 |
| E902 | io-error | {message} | Rule has been stable since v0.0.28 |
| E999 | syntax-error | SyntaxError | Rule was removed in 0.8.0 |

### Warning (W)

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| W191 | tab-indentation | Indentation contains tabs | Rule has been stable since v0.0.254 |
| W291 | trailing-whitespace | Trailing whitespace | Rule has been stable since v0.0.253Automatic fix available |
| W292 | missing-newline-at-end-of-file | No newline at end of file | Rule has been stable since v0.0.61Automatic fix available |
| W293 | blank-line-with-whitespace | Blank line contains whitespace | Rule has been stable since v0.0.253Automatic fix available |
| W391 | too-many-newlines-at-end-of-file | Too many newlines at end of {domain} | Rule has been in preview since v0.3.3Automatic fix available |
| W505 | doc-line-too-long | Doc line too long ({width} > {limit}) | Rule has been stable since v0.0.219 |
| W605 | invalid-escape-sequence | Invalid escape sequence: `\{ch}` | Rule has been stable since v0.0.85Automatic fix available |

## pydoclint (DOC)

For more, see pydoclint on PyPI.

For related settings, see pydoclint.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| DOC102 | docstring-extraneous-parameter | Documented parameter `{id}`is not in the function's signature | Rule has been in preview since 0.14.1 |
| DOC201 | docstring-missing-returns | `return`is not documented in docstring | Rule has been in preview since 0.5.6 |
| DOC202 | docstring-extraneous-returns | Docstring should not have a returns section because the function doesn't return anything | Rule has been in preview since 0.5.6 |
| DOC402 | docstring-missing-yields | `yield`is not documented in docstring | Rule has been in preview since 0.5.7 |
| DOC403 | docstring-extraneous-yields | Docstring has a "Yields" section but the function doesn't yield anything | Rule has been in preview since 0.5.7 |
| DOC501 | docstring-missing-exception | Raised exception `{id}`missing from docstring | Rule has been in preview since 0.5.5 |
| DOC502 | docstring-extraneous-exception | Raised exception is not explicitly raised: `{id}` | Rule has been in preview since 0.5.5 |

## pydocstyle (D)

For more, see pydocstyle on PyPI.

For related settings, see pydocstyle.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| D100 | undocumented-public-module | Missing docstring in public module | Rule has been stable since v0.0.70 |
| D101 | undocumented-public-class | Missing docstring in public class | Rule has been stable since v0.0.70 |
| D102 | undocumented-public-method | Missing docstring in public method | Rule has been stable since v0.0.70 |
| D103 | undocumented-public-function | Missing docstring in public function | Rule has been stable since v0.0.70 |
| D104 | undocumented-public-package | Missing docstring in public package | Rule has been stable since v0.0.70 |
| D105 | undocumented-magic-method | Missing docstring in magic method | Rule has been stable since v0.0.70 |
| D106 | undocumented-public-nested-class | Missing docstring in public nested class | Rule has been stable since v0.0.70 |
| D107 | undocumented-public-init | Missing docstring in `__init__` | Rule has been stable since v0.0.70 |
| D200 | unnecessary-multiline-docstring | One-line docstring should fit on one line | Rule has been stable since v0.0.68Automatic fix available |
| D201 | blank-line-before-function | No blank lines allowed before function docstring (found {num_lines}) | Rule has been stable since v0.0.70Automatic fix available |
| D202 | blank-line-after-function | No blank lines allowed after function docstring (found {num_lines}) | Rule has been stable since v0.0.70Automatic fix available |
| D203 | incorrect-blank-line-before-class | 1 blank line required before class docstring | Rule has been stable since v0.0.70Automatic fix available |
| D204 | incorrect-blank-line-after-class | 1 blank line required after class docstring | Rule has been stable since v0.0.70Automatic fix available |
| D205 | missing-blank-line-after-summary | 1 blank line required between summary line and description | Rule has been stable since v0.0.68Automatic fix available |
| D206 | docstring-tab-indentation | Docstring should be indented with spaces, not tabs | Rule has been stable since v0.0.75 |
| D207 | under-indentation | Docstring is under-indented | Rule has been stable since v0.0.75Automatic fix available |
| D208 | over-indentation | Docstring is over-indented | Rule has been stable since v0.0.75Automatic fix available |
| D209 | new-line-after-last-paragraph | Multi-line docstring closing quotes should be on a separate line | Rule has been stable since v0.0.68Automatic fix available |
| D210 | surrounding-whitespace | No whitespaces allowed surrounding docstring text | Rule has been stable since v0.0.68Automatic fix available |
| D211 | blank-line-before-class | No blank lines allowed before class docstring | Rule has been stable since v0.0.70Automatic fix available |
| D212 | multi-line-summary-first-line | Multi-line docstring summary should start at the first line | Rule has been stable since v0.0.69Automatic fix available |
| D213 | multi-line-summary-second-line | Multi-line docstring summary should start at the second line | Rule has been stable since v0.0.69Automatic fix available |
| D214 | overindented-section | Section is over-indented ("{name}") | Rule has been stable since v0.0.73Automatic fix available |
| D215 | overindented-section-underline | Section underline is over-indented ("{name}") | Rule has been stable since v0.0.73Automatic fix available |
| D300 | triple-single-quotes | Use triple double quotes `"""` | Rule has been stable since v0.0.69Automatic fix available |
| D301 | escape-sequence-in-docstring | Use `r"""`if any backslashes in a docstring | Rule has been stable since v0.0.172Automatic fix available |
| D400 | missing-trailing-period | First line should end with a period | Rule has been stable since v0.0.68Automatic fix available |
| D401 | non-imperative-mood | First line of docstring should be in imperative mood: "{first_line}" | Rule has been stable since v0.0.228 |
| D402 | signature-in-docstring | First line should not be the function's signature | Rule has been stable since v0.0.70 |
| D403 | first-word-uncapitalized | First word of the docstring should be capitalized: `{}`->`{}` | Rule has been stable since v0.0.69Automatic fix available |
| D404 | docstring-starts-with-this | First word of the docstring should not be "This" | Rule has been stable since v0.0.71 |
| D405 | non-capitalized-section-name | Section name should be properly capitalized ("{name}") | Rule has been stable since v0.0.71Automatic fix available |
| D406 | missing-new-line-after-section-name | Section name should end with a newline ("{name}") | Rule has been stable since v0.0.71Automatic fix available |
| D407 | missing-dashed-underline-after-section | Missing dashed underline after section ("{name}") | Rule has been stable since v0.0.71Automatic fix available |
| D408 | missing-section-underline-after-name | Section underline should be in the line following the section's name ("{name}") | Rule has been stable since v0.0.71Automatic fix available |
| D409 | mismatched-section-underline-length | Section underline should match the length of its name ("{name}") | Rule has been stable since v0.0.71Automatic fix available |
| D410 | no-blank-line-after-section | Missing blank line after section ("{name}") | Rule has been stable since v0.0.71Automatic fix available |
| D411 | no-blank-line-before-section | Missing blank line before section ("{name}") | Rule has been stable since v0.0.71Automatic fix available |
| D412 | blank-lines-between-header-and-content | No blank lines allowed between a section header and its content ("{name}") | Rule has been stable since v0.0.71Automatic fix available |
| D413 | missing-blank-line-after-last-section | Missing blank line after last section ("{name}") | Rule has been stable since v0.0.71Automatic fix available |
| D414 | empty-docstring-section | Section has no content ("{name}") | Rule has been stable since v0.0.71 |
| D415 | missing-terminal-punctuation | First line should end with a period, question mark, or exclamation point | Rule has been stable since v0.0.69Automatic fix available |
| D416 | missing-section-name-colon | Section name should end with a colon ("{name}") | Rule has been stable since v0.0.74Automatic fix available |
| D417 | undocumented-param | Missing argument description in the docstring for `{definition}`:`{name}` | Rule has been stable since v0.0.73 |
| D418 | overload-with-docstring | Function decorated with `@overload`shouldn't contain a docstring | Rule has been stable since v0.0.71 |
| D419 | empty-docstring | Docstring is empty | Rule has been stable since v0.0.68 |
| D420 | incorrect-section-order | Section "{current}" appears after section "{previous}" but should be before it | Rule has been in preview since 0.15.3 |
| D421 | property-docstring-starts-with-verb | Property docstring should not start with a verb ("{first_word}") | Rule has been in preview since 0.15.18 |

## Pyflakes (F)

For more, see Pyflakes on PyPI.

For related settings, see Pyflakes.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| F401 | unused-import | `{name}`imported but unused; consider using`importlib.util.find_spec`to test for availability | Rule has been stable since v0.0.18Automatic fix available |
| F402 | import-shadowed-by-loop-var | Import `{name}`from {row} shadowed by loop variable | Rule has been stable since v0.0.44 |
| F403 | undefined-local-with-import-star | `from {name} import *`used; unable to detect undefined names | Rule has been stable since v0.0.18 |
| F404 | late-future-import | `from __future__`imports must occur at the beginning of the file | Rule has been stable since v0.0.34 |
| F405 | undefined-local-with-import-star-usage | `{name}`may be undefined, or defined from star imports | Rule has been stable since v0.0.44 |
| F406 | undefined-local-with-nested-import-star-usage | `from {name} import *`only allowed at module level | Rule has been stable since v0.0.37 |
| F407 | future-feature-not-defined | Future feature `{name}`is not defined | Rule has been stable since v0.0.34 |
| F501 | percent-format-invalid-format | `%`-format string has invalid format string: {message} | Rule has been stable since v0.0.142 |
| F502 | percent-format-expected-mapping | `%`-format string expected mapping but got sequence | Rule has been stable since v0.0.142 |
| F503 | percent-format-expected-sequence | `%`-format string expected sequence but got mapping | Rule has been stable since v0.0.142 |
| F504 | percent-format-extra-named-arguments | `%`-format string has unused named argument(s): {message} | Rule has been stable since v0.0.142Automatic fix available |
| F505 | percent-format-missing-argument | `%`-format string is missing argument(s) for placeholder(s): {message} | Rule has been stable since v0.0.142 |
| F506 | percent-format-mixed-positional-and-named | `%`-format string has mixed positional and named placeholders | Rule has been stable since v0.0.142 |
| F507 | percent-format-positional-count-mismatch | `%`-format string has {wanted} placeholder(s) but {got} substitution(s) | Rule has been stable since v0.0.142 |
| F508 | percent-format-star-requires-sequence | `%`-format string`*`specifier requires sequence | Rule has been stable since v0.0.142 |
| F509 | percent-format-unsupported-format-character | `%`-format string has unsupported format character`{char}` | Rule has been stable since v0.0.142 |
| F521 | string-dot-format-invalid-format | `.format`call has invalid format string: {message} | Rule has been stable since v0.0.138 |
| F522 | string-dot-format-extra-named-arguments | `.format`call has unused named argument(s): {message} | Rule has been stable since v0.0.139Automatic fix available |
| F523 | string-dot-format-extra-positional-arguments | `.format`call has unused arguments at position(s): {message} | Rule has been stable since v0.0.139Automatic fix available |
| F524 | string-dot-format-missing-arguments | `.format`call is missing argument(s) for placeholder(s): {message} | Rule has been stable since v0.0.139 |
| F525 | string-dot-format-mixing-automatic | `.format`string mixes automatic and manual numbering | Rule has been stable since v0.0.139 |
| F541 | f-string-missing-placeholders | f-string without any placeholders | Rule has been stable since v0.0.18Automatic fix available |
| F601 | multi-value-repeated-key-literal | Dictionary key literal `{name}`repeated | Rule has been stable since v0.0.30Automatic fix available |
| F602 | multi-value-repeated-key-variable | Dictionary key `{name}`repeated | Rule has been stable since v0.0.30Automatic fix available |
| F621 | expressions-in-star-assignment | Too many expressions in star-unpacking assignment | Rule has been stable since v0.0.32 |
| F622 | multiple-starred-expressions | Two starred expressions in assignment | Rule has been stable since v0.0.32 |
| F631 | assert-tuple | Assert test is a non-empty tuple, which is always `True` | Rule has been stable since v0.0.28 |
| F632 | is-literal | Use `==`to compare constant literals | Rule has been stable since v0.0.39Automatic fix available |
| F633 | invalid-print-syntax | Use of `>>`is invalid with`print`function | Rule has been stable since v0.0.39 |
| F634 | if-tuple | If test is a tuple, which is always `True` | Rule has been stable since v0.0.18 |
| F701 | break-outside-loop | `break`outside loop | Rule has been stable since v0.0.36 |
| F702 | continue-outside-loop | `continue`not properly in loop | Rule has been stable since v0.0.36 |
| F704 | yield-outside-function | `{keyword}`statement outside of a function | Rule has been stable since v0.0.22 |
| F706 | return-outside-function | `return`statement outside of a function/method | Rule has been stable since v0.0.18 |
| F707 | default-except-not-last | An `except`block as not the last exception handler | Rule has been stable since v0.0.28 |
| F722 | forward-annotation-syntax-error | Syntax error in forward annotation: {parse_error} | Rule has been stable since v0.0.39 |
| F811 | redefined-while-unused | Redefinition of unused `{name}`from {row} | Rule has been stable since v0.0.171Automatic fix available |
| F821 | undefined-name | Undefined name `{name}`. {tip} | Rule has been stable since v0.0.20 |
| F822 | undefined-export | Undefined name `{name}`in`__all__` | Rule has been stable since v0.0.25 |
| F823 | undefined-local | Local variable `{name}`referenced before assignment | Rule has been stable since v0.0.24 |
| F841 | unused-variable | Local variable `{name}`is assigned to but never used | Rule has been stable since v0.0.22Automatic fix available |
| F842 | unused-annotation | Local variable `{name}`is annotated but never used | Rule has been stable since v0.0.172 |
| F901 | raise-not-implemented | `raise NotImplemented`should be`raise NotImplementedError` | Rule has been stable since v0.0.18Automatic fix available |

## pygrep-hooks (PGH)

For more, see pygrep-hooks on GitHub.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| PGH001 | eval | No builtin `eval()`allowed | Rule was removed in v0.2.0 |
| PGH002 | deprecated-log-warn | `warn`is deprecated in favor of`warning` | Rule was removed in v0.2.0Automatic fix available |
| PGH003 | blanket-type-ignore | Use specific rule codes when ignoring type issues | Rule has been stable since v0.0.187 |
| PGH004 | blanket-noqa | Use specific rule codes when using `noqa` | Rule has been stable since v0.0.200Automatic fix available |
| PGH005 | invalid-mock-access | Mock method should be called: `{name}` | Rule has been stable since v0.0.266 |

## Pylint (PL)

For more, see Pylint on PyPI.

For related settings, see Pylint.

### Convention (PLC)

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| PLC0105 | type-name-incorrect-variance | `{kind}`name "{param_name}" does not reflect its {variance}; consider renaming it to "{replacement_name}" | Rule has been stable since v0.0.278 |
| PLC0131 | type-bivariance | `{kind}`cannot be both covariant and contravariant | Rule has been stable since v0.0.278 |
| PLC0132 | type-param-name-mismatch | `{kind}`name`{param_name}`does not match assigned variable name`{var_name}` | Rule has been stable since v0.0.277 |
| PLC0205 | single-string-slots | Class `__slots__`should be a non-string iterable | Rule has been stable since v0.0.276 |
| PLC0206 | dict-index-missing-items | Extracting value from dictionary without calling `.items()` | Rule has been stable since 0.8.0 |
| PLC0207 | missing-maxsplit-arg | String is split more times than necessary | Rule has been stable since 0.15.0Automatic fix available |
| PLC0208 | iteration-over-set | Use a sequence type instead of a `set`when iterating over values | Rule has been stable since v0.0.271Automatic fix available |
| PLC0414 | useless-import-alias | Import alias does not rename original package | Rule has been stable since v0.0.156Automatic fix available |
| PLC0415 | import-outside-top-level | `import`should be at the top-level of a file | Rule has been stable since 0.12.0 |
| PLC1802 | len-test | `len({expression})`used as condition without comparison | Rule has been stable since 0.10.0Automatic fix available |
| PLC1901 | compare-to-empty-string | `{existing}`can be simplified to`{replacement}`as an empty string is falsey | Rule has been in preview since v0.0.255 |
| PLC2401 | non-ascii-name | {kind} name `{name}`contains a non-ASCII character | Rule has been stable since 0.5.0 |
| PLC2403 | non-ascii-import-name | Module alias `{name}`contains a non-ASCII character | Rule has been stable since 0.5.0 |
| PLC2701 | import-private-name | Private name import `{name}`from external module`{module}` | Rule has been in preview since v0.1.14 |
| PLC2801 | unnecessary-dunder-call | Unnecessary dunder call to `{method}`. {replacement}. | Rule has been in preview since v0.1.12Automatic fix available |
| PLC3002 | unnecessary-direct-lambda-call | Lambda expression called directly. Execute the expression inline instead. | Rule has been stable since v0.0.153 |

### Error (PLE)

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| PLE0100 | yield-in-init | `__init__`method is a generator | Rule has been stable since v0.0.245 |
| PLE0101 | return-in-init | Explicit return in `__init__` | Rule has been stable since v0.0.248 |
| PLE0115 | nonlocal-and-global | Name `{name}`is both`nonlocal`and`global` | Rule has been stable since 0.5.0 |
| PLE0116 | continue-in-finally | `continue`not supported inside`finally`clause | Rule has been stable since v0.0.257 |
| PLE0117 | nonlocal-without-binding | Nonlocal name `{name}`found without binding | Rule has been stable since v0.0.174 |
| PLE0118 | load-before-global-declaration | Name `{name}`is used prior to global declaration on {row} | Rule has been stable since v0.0.174 |
| PLE0237 | non-slot-assignment | Attribute `{name}`is not defined in class's`__slots__` | Rule has been stable since v0.1.15 |
| PLE0241 | duplicate-bases | Duplicate base `{base}`for class`{class}` | Rule has been stable since v0.0.269Automatic fix available |
| PLE0302 | unexpected-special-method-signature | The special method `{}`expects {}, {} {} given | Rule has been stable since v0.0.263 |
| PLE0303 | invalid-length-return-type | `__len__`does not return a non-negative integer | Rule has been stable since 0.6.0 |
| PLE0304 | invalid-bool-return-type | `__bool__`does not return`bool` | Rule has been stable since 0.16.0 |
| PLE0305 | invalid-index-return-type | `__index__`does not return an integer | Rule has been stable since 0.6.0 |
| PLE0307 | invalid-str-return-type | `__str__`does not return`str` | Rule has been stable since v0.0.271 |
| PLE0308 | invalid-bytes-return-type | `__bytes__`does not return`bytes` | Rule has been stable since 0.6.0 |
| PLE0309 | invalid-hash-return-type | `__hash__`does not return an integer | Rule has been stable since 0.6.0 |
| PLE0604 | invalid-all-object | Invalid object in `__all__`, must contain only strings | Rule has been stable since v0.0.237 |
| PLE0605 | invalid-all-format | Invalid format for `__all__`, must be`tuple`or`list` | Rule has been stable since v0.0.237 |
| PLE0643 | potential-index-error | Expression is likely to raise `IndexError` | Rule has been stable since 0.5.0 |
| PLE0704 | misplaced-bare-raise | Bare `raise`statement is not inside an exception handler | Rule has been stable since 0.5.0 |
| PLE1132 | repeated-keyword-argument | Repeated keyword argument: `{duplicate_keyword}` | Rule has been stable since 0.5.0 |
| PLE1141 | dict-iter-missing-items | Unpacking a dictionary in iteration without calling `.items()` | Rule has been in preview since v0.3.0Automatic fix available |
| PLE1142 | await-outside-async | `await`should be used within an async function | Rule has been stable since v0.0.150 |
| PLE1205 | logging-too-many-args | Too many arguments for `logging`format string | Rule has been stable since v0.0.252 |
| PLE1206 | logging-too-few-args | Not enough arguments for `logging`format string | Rule has been stable since v0.0.252 |
| PLE1300 | bad-string-format-character | Unsupported format character '{format_char}' | Rule has been stable since v0.0.283 |
| PLE1307 | bad-string-format-type | Format type does not match argument type | Rule has been stable since v0.0.245 |
| PLE1310 | bad-str-strip-call | String `{strip}`call contains duplicate characters (did you mean`{removal}`?) | Rule has been stable since v0.0.242 |
| PLE1507 | invalid-envvar-value | Invalid type for initial `os.getenv`argument; expected`str` | Rule has been stable since v0.0.255 |
| PLE1519 | singledispatch-method | `@singledispatch`decorator should not be used on methods | Rule has been stable since 0.6.0Automatic fix available |
| PLE1520 | singledispatchmethod-function | `@singledispatchmethod`decorator should not be used on non-method functions | Rule has been stable since 0.6.0Automatic fix available |
| PLE1700 | yield-from-in-async-function | `yield from`statement in async function; use`async for`instead | Rule has been stable since v0.0.271 |
| PLE2502 | bidirectional-unicode | Contains control characters that can permit obfuscated code | Rule has been stable since v0.0.244 |
| PLE2510 | invalid-character-backspace | Invalid unescaped character backspace, use "\b" instead | Rule has been stable since v0.0.257Automatic fix available |
| PLE2512 | invalid-character-sub | Invalid unescaped character SUB, use "\x1a" instead | Rule has been stable since v0.0.257Automatic fix available |
| PLE2513 | invalid-character-esc | Invalid unescaped character ESC, use "\x1b" instead | Rule has been stable since v0.0.257Automatic fix available |
| PLE2514 | invalid-character-nul | Invalid unescaped character NUL, use "\0" instead | Rule has been stable since v0.0.257Automatic fix available |
| PLE2515 | invalid-character-zero-width-space | Invalid unescaped character zero-width-space, use "\u200B" instead | Rule has been stable since v0.0.257Automatic fix available |
| PLE4703 | modified-iterating-set | Iterated set `{name}`is modified within the`for`loop | Rule has been in preview since v0.3.5Automatic fix available |

### Refactor (PLR)

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| PLR0124 | comparison-with-itself | Name compared with itself, consider replacing `{actual}` | Rule has been stable since v0.0.273 |
| PLR0133 | comparison-of-constant | Two constants compared in a comparison, consider replacing `{left_constant} {op} {right_constant}` | Rule has been stable since v0.0.221 |
| PLR0202 | no-classmethod-decorator | Class method defined without decorator | Rule has been in preview since v0.1.7Automatic fix available |
| PLR0203 | no-staticmethod-decorator | Static method defined without decorator | Rule has been in preview since v0.1.7Automatic fix available |
| PLR0206 | property-with-parameters | Cannot have defined parameters for properties | Rule has been stable since v0.0.153 |
| PLR0402 | manual-from-import | Use `from {module} import {name}`in lieu of alias | Rule has been stable since v0.0.155Automatic fix available |
| PLR0904 | too-many-public-methods | Too many public methods ({methods} > {max_methods}) | Rule has been in preview since v0.0.290 |
| PLR0911 | too-many-return-statements | Too many return statements ({returns} > {max_returns}) | Rule has been stable since v0.0.242 |
| PLR0912 | too-many-branches | Too many branches ({branches} > {max_branches}) | Rule has been stable since v0.0.242 |
| PLR0913 | too-many-arguments | Too many arguments in function definition ({c_args} > {max_args}) | Rule has been stable since v0.0.238 |
| PLR0914 | too-many-locals | Too many local variables ({current_amount} > {max_amount}) | Rule has been in preview since v0.1.9 |
| PLR0915 | too-many-statements | Too many statements ({statements} > {max_statements}) | Rule has been stable since v0.0.240 |
| PLR0916 | too-many-boolean-expressions | Too many Boolean expressions ({expressions} > {max_expressions}) | Rule has been in preview since v0.1.1 |
| PLR0917 | too-many-positional-arguments | Too many positional arguments ({c_pos} > {max_pos}) | Rule has been stable since 0.16.0 |
| PLR1701 | repeated-isinstance-calls | Merge `isinstance`calls:`{expression}` | Rule was removed in 0.5.0Automatic fix available |
| PLR1702 | too-many-nested-blocks | Too many nested blocks ({nested_blocks} > {max_nested_blocks}) | Rule has been in preview since v0.1.15 |
| PLR1704 | redefined-argument-from-local | Redefining argument with the local name `{name}` | Rule has been stable since 0.5.0 |
| PLR1706 | and-or-ternary | Consider using if-else expression | Rule was removed in v0.2.0 |
| PLR1708 | stop-iteration-return | Explicit `raise StopIteration`in generator | Rule has been stable since 0.16.0 |
| PLR1711 | useless-return | Useless `return`statement at end of function | Rule has been stable since v0.0.257Automatic fix available |
| PLR1712 | swap-with-temporary-variable | Unnecessary temporary variable | Rule has been in preview since 0.15.3Automatic fix available |
| PLR1714 | repeated-equality-comparison | Consider merging multiple comparisons: `{expression}`. Use a`set`if the elements are hashable. | Rule has been stable since v0.0.279Automatic fix available |
| PLR1716 | boolean-chained-comparison | Contains chained boolean comparison that can be simplified | Rule has been stable since 0.9.0Automatic fix available |
| PLR1722 | sys-exit-alias | Use `sys.exit()`instead of`{name}` | Rule has been stable since v0.0.156Automatic fix available |
| PLR1730 | if-stmt-min-max | Replace `if`statement with`{replacement}` | Rule has been stable since 0.6.0Automatic fix available |
| PLR1733 | unnecessary-dict-index-lookup | Unnecessary lookup of dictionary value by key | Rule has been stable since 0.12.0Automatic fix available |
| PLR1736 | unnecessary-list-index-lookup | List index lookup in `enumerate()`loop | Rule has been stable since 0.5.0Automatic fix available |
| PLR2004 | magic-value-comparison | Magic value used in comparison, consider replacing `{value}`with a constant variable | Rule has been stable since v0.0.221 |
| PLR2044 | empty-comment | Line with empty comment | Rule has been stable since 0.5.0Automatic fix available |
| PLR5501 | collapsible-else-if | Use `elif`instead of`else`then`if`, to reduce indentation | Rule has been stable since v0.0.253Automatic fix available |
| PLR6104 | non-augmented-assignment | Use `{operator}`to perform an augmented assignment directly | Rule has been in preview since v0.3.7Automatic fix available |
| PLR6201 | literal-membership | Use a set literal when testing for membership | Rule has been in preview since v0.1.1Automatic fix available |
| PLR6301 | no-self-use | Method `{method_name}`could be a function, class method, or static method | Rule has been in preview since v0.0.286 |

### Warning (PLW)

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| PLW0108 | unnecessary-lambda | Lambda may be unnecessary; consider inlining inner function | Rule has been stable since 0.15.0Automatic fix available |
| PLW0120 | useless-else-on-loop | `else`clause on loop without a`break`statement; remove the`else`and dedent its contents | Rule has been stable since v0.0.156Automatic fix available |
| PLW0127 | self-assigning-variable | Self-assignment of variable `{name}` | Rule has been stable since v0.0.281 |
| PLW0128 | redeclared-assigned-name | Redeclared variable `{name}`in assignment | Rule has been stable since 0.5.0 |
| PLW0129 | assert-on-string-literal | Asserting on an empty string literal will never pass | Rule has been stable since v0.0.258 |
| PLW0131 | named-expr-without-context | Named expression used without context | Rule has been stable since v0.0.270 |
| PLW0133 | useless-exception-statement | Missing `raise`statement on exception | Rule has been stable since 0.5.0Automatic fix available |
| PLW0177 | nan-comparison | Comparing against a NaN value; use `math.isnan`instead | Rule has been stable since 0.12.0 |
| PLW0211 | bad-staticmethod-argument | First argument of a static method should not be named `{argument_name}` | Rule has been stable since 0.6.0 |
| PLW0244 | redefined-slots-in-subclass | Slot `{slot_name}`redefined from base class`{base}` | Rule has been in preview since 0.9.3 |
| PLW0245 | super-without-brackets | `super`call is missing parentheses | Rule has been stable since 0.5.0Automatic fix available |
| PLW0406 | import-self | Module `{name}`imports itself | Rule has been stable since v0.0.265 |
| PLW0602 | global-variable-not-assigned | Using global for `{name}`but no assignment is done | Rule has been stable since v0.0.174 |
| PLW0603 | global-statement | Using the global statement to update `{name}`is discouraged | Rule has been stable since v0.0.253 |
| PLW0604 | global-at-module-level | `global`at module level is redundant | Rule has been stable since 0.5.0 |
| PLW0642 | self-or-cls-assignment | Reassigned `{}`variable in {method_type} method | Rule has been stable since 0.6.0 |
| PLW0711 | binary-op-exception | Exception to catch is the result of a binary `and`operation | Rule has been stable since v0.0.258 |
| PLW0717 | too-many-statements-in-try-clause | Try clause contains too many statements ({statements} > {max_statements}) | Rule has been in preview since 0.15.14 |
| PLW1501 | bad-open-mode | `{mode}`is not a valid mode for`open` | Rule has been stable since 0.5.0 |
| PLW1507 | shallow-copy-environ | Shallow copy of `os.environ`via`copy.copy(os.environ)` | Rule has been stable since 0.10.0Automatic fix available |
| PLW1508 | invalid-envvar-default | Invalid type for environment variable default; expected `str`or`None` | Rule has been stable since v0.0.255 |
| PLW1509 | subprocess-popen-preexec-fn | `preexec_fn`argument is unsafe when using threads | Rule has been stable since v0.0.281 |
| PLW1510 | subprocess-run-without-check | `subprocess.run`without explicit`check`argument | Rule has been stable since v0.0.285Automatic fix available |
| PLW1514 | unspecified-encoding | `{function_name}`in text mode without explicit`encoding`argument | Rule has been in preview since v0.1.1Automatic fix available |
| PLW1641 | eq-without-hash | Object does not implement `__hash__`method | Rule has been stable since 0.12.0 |
| PLW2101 | useless-with-lock | Threading lock directly created in `with`statement has no effect | Rule has been stable since 0.5.0 |
| PLW2901 | redefined-loop-name | Outer {outer_kind} variable `{name}`overwritten by inner {inner_kind} target | Rule has been stable since v0.0.252 |
| PLW3201 | bad-dunder-method-name | Dunder method `{name}`has no special meaning in Python 3 | Rule has been in preview since v0.0.285 |
| PLW3301 | nested-min-max | Nested `{func}`calls can be flattened | Rule has been stable since v0.0.266Automatic fix available |

## pyupgrade (UP)

For more, see pyupgrade on PyPI.

For related settings, see pyupgrade.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| UP001 | useless-metaclass-type | `__metaclass__ = type`is implied | Rule has been stable since v0.0.155Automatic fix available |
| UP003 | type-of-primitive | Use `{}`instead of`type(...)` | Rule has been stable since v0.0.155Automatic fix available |
| UP004 | useless-object-inheritance | Class `{name}`inherits from`object` | Rule has been stable since v0.0.155Automatic fix available |
| UP005 | deprecated-unittest-alias | `{alias}`is deprecated, use`{target}` | Rule has been stable since v0.0.155Automatic fix available |
| UP006 | non-pep585-annotation | Use `{to}`instead of`{from}`for type annotation | Rule has been stable since v0.0.155Automatic fix available |
| UP007 | non-pep604-annotation-union | Use `X | Y`for type annotations | Rule has been stable since v0.0.155Automatic fix available |
| UP008 | super-call-with-parameters | Use `super()`instead of`super(__class__, self)` | Rule has been stable since v0.0.155Automatic fix available |
| UP009 | utf8-encoding-declaration | UTF-8 encoding declaration is unnecessary | Rule has been stable since v0.0.155Automatic fix available |
| UP010 | unnecessary-future-import | Unnecessary `__future__`import`{import}`for target Python version | Rule has been stable since v0.0.155Automatic fix available |
| UP011 | lru-cache-without-parameters | Unnecessary parentheses to `functools.lru_cache` | Rule has been stable since v0.0.155Automatic fix available |
| UP012 | unnecessary-encode-utf8 | Unnecessary call to `encode`as UTF-8 | Rule has been stable since v0.0.155Automatic fix available |
| UP013 | convert-typed-dict-functional-to-class | Convert `{name}`from`TypedDict`functional to class syntax | Rule has been stable since v0.0.155Automatic fix available |
| UP014 | convert-named-tuple-functional-to-class | Convert `{name}`from`NamedTuple`functional to class syntax | Rule has been stable since v0.0.155Automatic fix available |
| UP015 | redundant-open-modes | Unnecessary mode argument | Rule has been stable since v0.0.155Automatic fix available |
| UP017 | datetime-timezone-utc | Use `datetime.UTC`alias | Rule has been stable since v0.0.192Automatic fix available |
| UP018 | native-literals | Unnecessary `{literal_type}`call (rewrite as a literal) | Rule has been stable since v0.0.193Automatic fix available |
| UP019 | typing-text-str-alias | `{}.Text`is deprecated, use`str` | Rule has been stable since v0.0.195Automatic fix available |
| UP020 | open-alias | Use builtin `open` | Rule has been stable since v0.0.196Automatic fix available |
| UP021 | replace-universal-newlines | `universal_newlines`is deprecated, use`text` | Rule has been stable since v0.0.196Automatic fix available |
| UP022 | replace-stdout-stderr | Prefer `capture_output`over sending`stdout`and`stderr`to`PIPE` | Rule has been stable since v0.0.199Automatic fix available |
| UP023 | deprecated-c-element-tree | `cElementTree`is deprecated, use`ElementTree` | Rule has been stable since v0.0.199Automatic fix available |
| UP024 | os-error-alias | Replace aliased errors with `OSError` | Rule has been stable since v0.0.206Automatic fix available |
| UP025 | unicode-kind-prefix | Remove unicode literals from strings | Rule has been stable since v0.0.201Automatic fix available |
| UP026 | deprecated-mock-import | `mock`is deprecated, use`unittest.mock` | Rule has been stable since v0.0.206Automatic fix available |
| UP027 | unpacked-list-comprehension | Replace unpacked list comprehension with a generator expression | Rule was removed in 0.8.0 |
| UP028 | yield-in-for-loop | Replace `yield`over`for`loop with`yield from` | Rule has been stable since v0.0.210Automatic fix available |
| UP029 | unnecessary-builtin-import | Unnecessary builtin import: `{import}` | Rule has been stable since v0.0.211Automatic fix available |
| UP030 | format-literals | Use implicit references for positional format fields | Rule has been stable since v0.0.218Automatic fix available |
| UP031 | printf-string-formatting | Use format specifiers instead of percent format | Rule has been stable since v0.0.229Automatic fix available |
| UP032 | f-string | Use f-string instead of `format`call | Rule has been stable since v0.0.224Automatic fix available |
| UP033 | lru-cache-with-maxsize-none | Use `@functools.cache`instead of`@functools.lru_cache(maxsize=None)` | Rule has been stable since v0.0.225Automatic fix available |
| UP034 | extraneous-parentheses | Avoid extraneous parentheses | Rule has been stable since v0.0.228Automatic fix available |
| UP035 | deprecated-import | Import from `{target}`instead: {names} | Rule has been stable since v0.0.239Automatic fix available |
| UP036 | outdated-version-block | Version block is outdated for minimum Python version | Rule has been stable since v0.0.240Automatic fix available |
| UP037 | quoted-annotation | Remove quotes from type annotation | Rule has been stable since v0.0.242Automatic fix available |
| UP038 | non-pep604-isinstance | Use `X | Y`in`{}`call instead of`(X, Y)` | Rule was removed in 0.13.0Automatic fix available |
| UP039 | unnecessary-class-parentheses | Unnecessary parentheses after class definition | Rule has been stable since v0.0.273Automatic fix available |
| UP040 | non-pep695-type-alias | Type alias `{name}`uses {type_alias_method} instead of the`type`keyword | Rule has been stable since v0.0.283Automatic fix available |
| UP041 | timeout-error-alias | Replace aliased errors with `TimeoutError` | Rule has been stable since v0.2.0Automatic fix available |
| UP042 | replace-str-enum | Class {name} inherits from both `str`and`enum.Enum` | Rule has been stable since 0.15.0Automatic fix available |
| UP043 | unnecessary-default-type-args | Unnecessary default type arguments | Rule has been stable since 0.8.0Automatic fix available |
| UP044 | non-pep646-unpack | Use `*`for unpacking | Rule has been stable since 0.10.0Automatic fix available |
| UP045 | non-pep604-annotation-optional | Use `X | None`for type annotations | Rule has been stable since 0.12.0Automatic fix available |
| UP046 | non-pep695-generic-class | Generic class `{name}`uses`Generic`subclass instead of type parameters | Rule has been stable since 0.12.0Automatic fix available |
| UP047 | non-pep695-generic-function | Generic function `{name}`should use type parameters | Rule has been stable since 0.12.0Automatic fix available |
| UP049 | private-type-parameter | Generic {} uses private type parameters | Rule has been stable since 0.12.0Automatic fix available |
| UP050 | useless-class-metaclass-type | Class `{name}`uses`metaclass=type`, which is redundant | Rule has been stable since 0.13.0Automatic fix available |
| UP051 | deprecated-abc-decorator | Use `@{to}`and`@abstractmethod`instead of`{from}` | Rule has been in preview since 0.15.21Automatic fix available |

## refurb (FURB)

For more, see refurb on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| FURB101 | read-whole-file | `Path.open()`followed by`read()`can be replaced by`{filename}.{suggestion}` | Rule has been in preview since v0.1.2Automatic fix available |
| FURB103 | write-whole-file | `Path.open()`followed by`write()`can be replaced by`{filename}.{suggestion}` | Rule has been in preview since v0.3.6Automatic fix available |
| FURB105 | print-empty-string | Unnecessary empty string passed to `print` | Rule has been stable since 0.5.0Automatic fix available |
| FURB110 | if-exp-instead-of-or-operator | Replace ternary `if`expression with`or`operator | Rule has been stable since 0.15.0Automatic fix available |
| FURB113 | repeated-append | Use `{suggestion}`instead of repeatedly calling`{name}.append()` | Rule has been in preview since v0.0.287Automatic fix available |
| FURB116 | f-string-number-format | Replace `{function_name}`call with`{display}` | Rule has been stable since 0.13.0Automatic fix available |
| FURB118 | reimplemented-operator | Use `operator.{operator}`instead of defining a {target} | Rule has been in preview since v0.1.9Automatic fix available |
| FURB122 | for-loop-writes | Use of `{}.write`in a for loop | Rule has been stable since 0.12.0Automatic fix available |
| FURB129 | readlines-in-for | Instead of calling `readlines()`, iterate over file object directly | Rule has been stable since 0.5.0Automatic fix available |
| FURB131 | delete-full-slice | Prefer `clear`over deleting a full slice | Rule has been in preview since v0.0.287Automatic fix available |
| FURB132 | check-and-remove-from-set | Use `{suggestion}`instead of check and`remove` | Rule has been stable since 0.12.0Automatic fix available |
| FURB136 | if-expr-min-max | Replace `if`expression with`{min_max}`call | Rule has been stable since 0.5.0Automatic fix available |
| FURB140 | reimplemented-starmap | Use `itertools.starmap`instead of the generator | Rule has been in preview since v0.0.291Automatic fix available |
| FURB142 | for-loop-set-mutations | Use of `set.{}()`in a for loop | Rule has been in preview since v0.3.5Automatic fix available |
| FURB145 | slice-copy | Prefer `copy`method over slicing | Rule has been in preview since v0.0.290Automatic fix available |
| FURB148 | unnecessary-enumerate | `enumerate`value is unused, use`for x in range(len(y))`instead | Rule has been in preview since v0.0.291Automatic fix available |
| FURB152 | math-constant | Replace `{literal}`with`math.{constant}` | Rule has been in preview since v0.1.6Automatic fix available |
| FURB154 | repeated-global | Use of repeated consecutive `{}` | Rule has been in preview since v0.4.9Automatic fix available |
| FURB156 | hardcoded-string-charset | Use of hardcoded string charset | Rule has been in preview since 0.7.0Automatic fix available |
| FURB157 | verbose-decimal-constructor | Verbose expression in `Decimal`constructor | Rule has been stable since 0.12.0Automatic fix available |
| FURB161 | bit-count | Use of `bin({existing}).count('1')` | Rule has been stable since 0.5.0Automatic fix available |
| FURB162 | fromisoformat-replace-z | Unnecessary timezone replacement with zero offset | Rule has been stable since 0.12.0Automatic fix available |
| FURB163 | redundant-log-base | Prefer `math.{log_function}({arg})`over`math.log`with a redundant base | Rule has been stable since 0.5.0Automatic fix available |
| FURB164 | unnecessary-from-float | Verbose method `{method_name}`in`{constructor}`construction | Rule has been stable since 0.16.0Automatic fix available |
| FURB166 | int-on-sliced-str | Use of `int`with explicit`base={base}`after removing prefix | Rule has been stable since 0.12.0Automatic fix available |
| FURB167 | regex-flag-alias | Use of regular expression alias `re.{}` | Rule has been stable since 0.5.0Automatic fix available |
| FURB168 | isinstance-type-none | Prefer `is`operator over`isinstance`to check if an object is`None` | Rule has been stable since 0.5.0Automatic fix available |
| FURB169 | type-none-comparison | When checking against `None`, use`{}`instead of comparison with`type(None)` | Rule has been stable since 0.5.0Automatic fix available |
| FURB171 | single-item-membership-test | Membership test against single-item container | Rule has been stable since 0.15.0Automatic fix available |
| FURB177 | implicit-cwd | Prefer `Path.cwd()`over`Path().resolve()`for current-directory lookups | Rule has been stable since 0.5.0Automatic fix available |
| FURB180 | meta-class-abc-meta | Use of `metaclass=abc.ABCMeta`to define abstract base class | Rule has been in preview since v0.2.0Automatic fix available |
| FURB181 | hashlib-digest-hex | Use of hashlib's `.digest().hex()` | Rule has been stable since 0.5.0Automatic fix available |
| FURB187 | list-reverse-copy | Use of assignment of `reversed`on list`{name}` | Rule has been stable since 0.5.0Automatic fix available |
| FURB188 | slice-to-remove-prefix-or-suffix | Prefer `str.removeprefix()`over conditionally replacing with slice. | Rule has been stable since 0.9.0Automatic fix available |
| FURB189 | subclass-builtin | Subclassing `{subclass}`can be error prone, use`collections.{replacement}`instead | Rule has been in preview since 0.7.3Automatic fix available |
| FURB192 | sorted-min-max | Prefer `min`over`sorted()`to compute the minimum value in a sequence | Rule has been stable since 0.16.0Automatic fix available |

## Ruff-specific rules (RUF)

For related settings, see Ruff.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| RUF001 | ambiguous-unicode-character-string | String contains ambiguous {}. Did you mean {}? | Rule has been stable since v0.0.102 |
| RUF002 | ambiguous-unicode-character-docstring | Docstring contains ambiguous {}. Did you mean {}? | Rule has been stable since v0.0.102 |
| RUF003 | ambiguous-unicode-character-comment | Comment contains ambiguous {}. Did you mean {}? | Rule has been stable since v0.0.108 |
| RUF005 | collection-literal-concatenation | Consider `{expression}`instead of concatenation | Rule has been stable since v0.0.227Automatic fix available |
| RUF006 | asyncio-dangling-task | Store a reference to the return value of `{expr}.{method}` | Rule has been stable since v0.0.247 |
| RUF007 | zip-instead-of-pairwise | Prefer `itertools.pairwise()`over`zip()`when iterating over successive pairs | Rule has been stable since v0.0.257Automatic fix available |
| RUF008 | mutable-dataclass-default | Do not use mutable default values for dataclass attributes | Rule has been stable since v0.0.262 |
| RUF009 | function-call-in-dataclass-default-argument | Do not perform function call `{name}`in dataclass defaults | Rule has been stable since v0.0.262 |
| RUF010 | explicit-f-string-type-conversion | Use explicit conversion flag | Rule has been stable since v0.0.267Automatic fix available |
| RUF011 | ruff-static-key-dict-comprehension | Dictionary comprehension uses static key | Rule was removed in v0.2.0 |
| RUF012 | mutable-class-default | Mutable default value for class attribute | Rule has been stable since v0.0.273 |
| RUF013 | implicit-optional | PEP 484 prohibits implicit `Optional` | Rule has been stable since v0.0.273Automatic fix available |
| RUF015 | unnecessary-iterable-allocation-for-first-element | Prefer `next({iterable})`over single element slice | Rule has been stable since v0.0.278Automatic fix available |
| RUF016 | invalid-index-type | Slice in indexed access to type `{value_type}`uses type`{index_type}`instead of an integer | Rule has been stable since v0.0.278 |
| RUF017 | quadratic-list-summation | Avoid quadratic list summation | Rule has been stable since v0.0.285Automatic fix available |
| RUF018 | assignment-in-assert | Avoid assignment expressions in `assert`statements | Rule has been stable since v0.2.0 |
| RUF019 | unnecessary-key-check | Unnecessary key check before dictionary access | Rule has been stable since v0.2.0Automatic fix available |
| RUF020 | never-union | `{never_like} | T`is equivalent to`T` | Rule has been stable since v0.2.0Automatic fix available |
| RUF021 | parenthesize-chained-operators | Parenthesize `a and b`expressions when chaining`and`and`or`together, to make the precedence clear | Rule has been stable since 0.8.0Automatic fix available |
| RUF022 | unsorted-dunder-all | `__all__`is not sorted | Rule has been stable since 0.8.0Automatic fix available |
| RUF023 | unsorted-dunder-slots | `{}.__slots__`is not sorted | Rule has been stable since 0.8.0Automatic fix available |
| RUF024 | mutable-fromkeys-value | Do not pass mutable objects as values to `dict.fromkeys` | Rule has been stable since 0.5.0Automatic fix available |
| RUF026 | default-factory-kwarg | `default_factory`is a positional-only argument to`defaultdict` | Rule has been stable since 0.5.0Automatic fix available |
| RUF027 | missing-f-string-syntax | Possible f-string without an `f`prefix | Rule has been in preview since v0.2.1Automatic fix available |
| RUF028 | invalid-formatter-suppression-comment | This suppression comment is invalid because {} | Rule has been stable since 0.12.0Automatic fix available |
| RUF029 | unused-async | Function `{name}`is declared`async`, but doesn't`await`or use`async`features. | Rule has been in preview since v0.4.0 |
| RUF030 | assert-with-print-message | `print()`call in`assert`statement is likely unintentional | Rule has been stable since 0.8.0Automatic fix available |
| RUF031 | incorrectly-parenthesized-tuple-in-subscript | Use parentheses for tuples in subscripts | Rule has been in preview since 0.5.7Automatic fix available |
| RUF032 | decimal-from-float-literal | `Decimal()`called with float literal argument | Rule has been stable since 0.9.0Automatic fix available |
| RUF033 | post-init-default | `__post_init__`method with argument defaults | Rule has been stable since 0.9.0Automatic fix available |
| RUF034 | useless-if-else | Useless `if`-`else`condition | Rule has been stable since 0.9.0 |
| RUF035 | ruff-unsafe-markup-use | Unsafe use of `{name}`detected | Rule was removed in 0.10.0 |
| RUF036 | none-not-at-end-of-union | `None`not at the end of the type union. | Rule has been stable since 0.16.0Automatic fix available |
| RUF037 | unnecessary-empty-iterable-within-deque-call | Unnecessary empty iterable within a deque call | Rule has been stable since 0.15.0Automatic fix available |
| RUF038 | redundant-bool-literal | `Literal[True, False, ...]`can be replaced with`Literal[...] | bool` | Rule has been in preview since 0.8.0Automatic fix available |
| RUF039 | unraw-re-pattern | First argument to {call} is not raw string | Rule has been in preview since 0.8.0Automatic fix available |
| RUF040 | invalid-assert-message-literal-argument | Non-string literal used as assert message | Rule has been stable since 0.10.0 |
| RUF041 | unnecessary-nested-literal | Unnecessary nested `Literal` | Rule has been stable since 0.10.0Automatic fix available |
| RUF043 | pytest-raises-ambiguous-pattern | Pattern passed to `match=`contains metacharacters but is neither escaped nor raw | Rule has been stable since 0.13.0 |
| RUF045 | implicit-class-var-in-dataclass | Assignment without annotation found in dataclass body | Rule has been in preview since 0.9.7 |
| RUF046 | unnecessary-cast-to-int | Value being cast to `int`is already an integer | Rule has been stable since 0.10.0Automatic fix available |
| RUF047 | needless-else | Empty `else`clause | Rule has been in preview since 0.9.3Automatic fix available |
| RUF048 | map-int-version-parsing | `__version__`may contain non-integral-like elements | Rule has been stable since 0.10.0 |
| RUF049 | dataclass-enum | An enum class should not be decorated with `@dataclass` | Rule has been stable since 0.12.0 |
| RUF050 | unnecessary-if | Empty `if`statement | Rule has been in preview since 0.15.8Automatic fix available |
| RUF051 | if-key-in-dict-del | Use `pop`instead of`key in dict`followed by`del dict[key]` | Rule has been stable since 0.10.0Automatic fix available |
| RUF052 | used-dummy-variable | Local dummy variable `{}`is accessed | Rule has been in preview since 0.8.2Automatic fix available |
| RUF053 | class-with-mixed-type-vars | Class with type parameter list inherits from `Generic` | Rule has been stable since 0.12.0Automatic fix available |
| RUF054 | indented-form-feed | Indented form feed | Rule has been in preview since 0.9.6 |
| RUF055 | unnecessary-regular-expression | Plain string pattern passed to `re`function | Rule has been in preview since 0.8.1Automatic fix available |
| RUF056 | falsy-dict-get-fallback | Avoid providing a falsy fallback to `dict.get()`in boolean test positions. The default fallback`None`is already falsy. | Rule has been in preview since 0.8.5Automatic fix available |
| RUF057 | unnecessary-round | Value being rounded is already an integer | Rule has been stable since 0.12.0Automatic fix available |
| RUF058 | starmap-zip | `itertools.starmap`called on`zip`iterable | Rule has been stable since 0.12.0Automatic fix available |
| RUF059 | unused-unpacked-variable | Unpacked variable `{name}`is never used | Rule has been stable since 0.13.0Automatic fix available |
| RUF060 | in-empty-collection | Unnecessary membership test on empty collection | Rule has been stable since 0.15.0 |
| RUF061 | legacy-form-pytest-raises | Use context-manager form of `pytest.{}()` | Rule has been stable since 0.15.0Automatic fix available |
| RUF063 | access-annotations-from-class-dict | Use `{suggestion}`instead of`__dict__`access | Rule has been stable since 0.16.0 |
| RUF064 | non-octal-permissions | Non-octal mode | Rule has been stable since 0.15.0Automatic fix available |
| RUF065 | logging-eager-conversion | Unnecessary `oct()`conversion when formatting with`%s`. Use`%#o`instead of`%s` | Rule has been in preview since 0.13.2 |
| RUF066 | property-without-return | `{name}`is a property without a`return`statement | Rule has been in preview since 0.14.7 |
| RUF067 | non-empty-init-module | `__init__`module should not contain any code | Rule has been in preview since 0.14.11 |
| RUF068 | duplicate-entry-in-dunder-all | `__all__`contains duplicate entries | Rule has been stable since 0.16.0Automatic fix available |
| RUF069 | float-equality-comparison | Unreliable floating point equality comparison `{left} {operator} {right}` | Rule has been in preview since 0.15.1 |
| RUF070 | unnecessary-assign-before-yield | Unnecessary assignment to `{name}`before`yield from`statement | Rule has been in preview since 0.15.3Automatic fix available |
| RUF071 | os-path-commonprefix | `os.path.commonprefix()`compares strings character-by-character | Rule has been in preview since 0.15.6Automatic fix available |
| RUF072 | useless-finally | Empty `finally`clause | Rule has been in preview since 0.15.8Automatic fix available |
| RUF073 | f-string-percent-format | `%`operator used on an f-string | Rule has been in preview since 0.15.8 |
| RUF074 | incorrect-decorator-order | `@{outer_decorator}`should be placed after`@{inner_decorator}` | Rule has been in preview since 0.15.14 |
| RUF075 | fallible-context-manager | Context manager does not catch exceptions | Rule has been in preview since 0.15.14 |
| RUF076 | pytest-fixture-autouse | Avoid using `autouse=True`in`pytest.fixture`decorators | Rule was removed in 0.15.20 |
| RUF100 | unused-noqa | Unused {} | Rule has been stable since v0.0.155Automatic fix available |
| RUF101 | redirected-noqa | `{original}`is a redirect to`{target}` | Rule has been stable since 0.6.0Automatic fix available |
| RUF102 | invalid-rule-code | Invalid rule code in {}: {} | Rule has been stable since 0.15.0Automatic fix available |
| RUF103 | invalid-suppression-comment | Invalid suppression comment: {msg} | Rule has been stable since 0.15.0Automatic fix available |
| RUF104 | unmatched-suppression-comment | Suppression comment without matching `#ruff:enable`comment | Rule has been stable since 0.15.0 |
| RUF105 | noqa-comments | `noqa`comment used instead of`ruff: ignore` | Rule has been in preview since 0.15.22Automatic fix available |
| RUF106 | rule-codes-in-suppression-comments | Rule code used instead of name in suppression comment | Rule has been in preview since 0.15.22Automatic fix available |
| RUF200 | invalid-pyproject-toml | Failed to parse pyproject.toml: {message} | Rule has been stable since v0.0.271 |
| RUF201 | rule-codes-in-selectors | Rule code used instead of name in `lint.{selector}` | Rule has been in preview since 0.15.22Automatic fix available |

## tryceratops (TRY)

For more, see tryceratops on PyPI.

| Code | Name | Message | Fix/Status |
|---|---|---|---|
| TRY002 | raise-vanilla-class | Create your own exception | Rule has been stable since v0.0.236 |
| TRY003 | raise-vanilla-args | Avoid specifying long messages outside the exception class | Rule has been stable since v0.0.236 |
| TRY004 | type-check-without-type-error | Prefer `TypeError`exception for invalid type | Rule has been stable since v0.0.230 |
| TRY200 | reraise-no-cause | Use `raise from`to specify exception cause | Rule was removed in v0.2.0 |
| TRY201 | verbose-raise | Use `raise`without specifying exception name | Rule has been stable since v0.0.231Automatic fix available |
| TRY203 | useless-try-except | Remove exception handler; error is immediately re-raised | Rule has been stable since 0.7.0 |
| TRY300 | try-consider-else | Consider moving this statement to an `else`block | Rule has been stable since v0.0.229 |
| TRY301 | raise-within-try | Abstract `raise`to an inner function | Rule has been stable since v0.0.233 |
| TRY400 | error-instead-of-exception | Use `logging.exception`instead of`logging.error` | Rule has been stable since v0.0.236Automatic fix available |
| TRY401 | verbose-log-message | Redundant exception object included in `logging.exception`call | Rule has been stable since v0.0.250 |
