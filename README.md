# rhsm-server

RHSM systemd service providing Varlink API

Build instructions
------------------

If you want to build the project, you can use the following command:

```console
$ make
```

If you want to create rpm package, then it is possible to build it using:

```console
$ make rpm
```

Integration tests
-----------------

To run integration tests, you have to have [behave](https://github.com/behave/behave) installed and you
have to install generated rhsm-server RPM packe. Then you can run integration tests with:

```console
$ sudo behave
```