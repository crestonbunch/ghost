Ghost: A Lightweight REST Framework
===================================

Ghost is a totally invisible framework for writing *testable* REST APIs. The
idea is that each resource is broken up into several components:

* An input model (which gets parsed from the request)
* A validator which validates the input model
* A processor which processes the input and possibly returns a different model,
  it may also write to the database / etc.
* A writer which decides how to serialize the processed data for the user.

See the examples/ directory for a straightforward how-to.

By separating resources into components, each resource can be tested by
component without needing a monolithic testing suite for every single resource.

Ghost is probably simple enough that you can write something like it yourself,
but if you want to contribute back then please feel free to fork and submit a
pull request.
