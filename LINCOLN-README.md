# Building click code

First we need to create the protobuf.  The protobuf for our click container is from https://pulwar.isi.edu/sabres/moa.git, which is a submodule with path pfcpiface/click\_pb/moa. We then do a submodule update to pull the latest code, then we need to generate the protobuf, which we do in a container and copy over the results.

```
make lincoln-test
```

Now the generated protobuf has been built, we can build our code

```
make test-click-integration-native
```
