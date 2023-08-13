# go_app_template
This is a basic template to start a new golang application based on clean architecture 

## Setup of the component:

Must have golang installed version >= 1.18

make file consists of 4 steps: generate, test, build, run you can run all of them

```make all```

Run test cases run

```make test```
or

```go test -v ./tests```

Start the http server on port 9090:

```make run```

## Run By Docker

```
 docker build -t go_app .
 
 docker run -p 9090:9090 go_app
 
```
<p align="center"><img src="https://user-images.githubusercontent.com/1092882/86512217-bfd5a480-be1d-11ea-976c-a7c0ac0cd1f1.png" alt="goapp gopher" width="256px"/></p>


P.S: This guideline is not directly applicable for an independent package, as their primary use is to be consumed in other applications. In such cases, having most or all of the package code in the root is probably the best way of doing it. And that is where Go's recommendation of "no unnecessary sub packages" shines.

In my effort to try and make things easier to understand, the structure is explained based on a note taking web application (with hardly any features implemented 🤭).

## Table of contents

1. [Directory structure](#directory-structure)
2. [Configs package](#internalconfigs)
3. [API package](#internalapi)
4. [Users](#internalusers) (would be common for all such business logic units, 'notes' being similar to users) package.
5. [Testing](#internalusers_test)
6. [pkg package](#internalpkg)
   - 6.1. [datastore](#internalpkgdatastore)
   - 6.2. [logger](#internalpkglogger)
7. [HTTP server](#internalhttp)
8. [lib](#lib)
10. [docker](#docker)
11. [schemas](#schemas)
12. [main.go](#maingo)
13. [Error handling](#error-handling)
14. [Dependency flow](#dependency-flow)
16. [Note](#note)

## Directory structure

```bash
|
|____internal
|    |
|    |____configs
|    |    |____configs.go
|    |
|    |____api
|    |    |____notes.go
|    |    |____users.go
|    |
|    |____users
|    |    |____store.go
|    |    |____cache.go
|    |    |____users.go
|    |    |____users_test.go
|    |
|    |____notes
|    |    |____notes.go
|    |
|    |____pkg
|    |    |____stringutils
|    |    |____datastore
|    |    |     |____datastore.go
|    |    |____logger
|    |         |____logger.go
|    |
|    |____server
|         |____http
|         |    |____web
|         |    |    |____templates
|         |    |         |____index.html
|         |    |____handlers_notes.go
|         |    |____handlers_users.go
|         |    |____http.go
|         |
|
|____lib
|    |____notes
|         |____notes.go
|
|
|____docker
|    |____Dockerfile # your 'default' dockerfile
|
|____go.mod
|____go.sum
|
|____ciconfig.yml # depends on the CI/CD system you're using. e.g. .travis.yml
|____README.md
|____main.go
|
```

## internal

["internal" is a special directory name in Go](https://go.dev/doc/go1.4#internalpackages), wherein any exported name/entity can only be consumed within its immediate parent.

## internal/configs

Creating a dedicated configs package might seem like an overkill, but it makes a lot of things easier. In the example app provided, you see the HTTP configs are hardcoded and returned. Later you decide to change to consume from env variables. All you do is update the configs package. And further down the line, maybe you decide to introduce something like [etcd](https://github.com/etcd-io/etcd), then you define the dependency in `Configs` and update the functions accordingly. This is yet another separation of concern package, to try and keep `main` tidy.

## internal/api

The API package is supposed to have all the APIs exposed by the application. A dedicated API package is created to standardize the functionality, when there are different kinds of servers running. e.g. an HTTP & a gRPC server. In such cases, the respective "handler" functions would inturn call `api.<Method name>`. This gives a guarantee that all your APIs behave exactly the same without any accidental inconsistencies across different I/O methods.

But remember, middleware handling is still at the internal/server layer. e.g. access log, authentication etc. Even though this can be brought to the `api` package, it doesn't make much sense because middleware are mostly dependent on the server/handler implementation.

## internal/users

Users package is where all your actual user related business logic is implemented. e.g. Create a user after cleaning up the input, validation, and then store it inside a persistent datastore.

There's a `store.go` in this package which is where you write all the direct interactions with the datastore. There's an interface which is unique to the `users` package. It is introduced to handle dependency injection as well as dependency inversion elegantly. File naming convention for store files is `store_<logical group>.go`. e.g. `store_aggregations.go`. Or simply `store.go` if there's not much code.

`NewService/New` function is created in each package, which initializes and returns the respective package's handler. In case of users package, there's a `Users` struct. The name 'NewService' makes sense in most cases, and just reduces the burden of thinking of a good name for such scenarios. The Users struct here holds all the dependencies required for implementing features provided by users package.

### conclusion

At this point where you're testing individual package's datastore interaction, I'd rather you directly start testing the API. APIs would cover all the layers, API, business logic, datastore interaction etc. These tests can be built and deployed using external API testing frameworks (i.e. independent of your code). So my approach is a hybrid one, unit tests for all possible pure functions, and API test for the rest. And when it comes to API testing, your aim should be to try and "break the application". i.e. don't just cover happy paths. The lazier you are, more pure functions you will have(rather write unit tests than create API tests on yet another tool)!

P.S: I use [VSCode](https://code.visualstudio.com/) and it lets you auto [generate unit tests](https://code.visualstudio.com/docs/languages/go#_test). I'm positive other IDEs also have similar functionality. You could just right-click on the function and choose `Go: Generate unit tests for function`.

<p align="center">
<img src="https://user-images.githubusercontent.com/1092882/87034896-2be14f80-c206-11ea-8771-edabfffab39d.png" alt="generate unit test" width="384px" height="256px" style="margin-right: 16px" />
</p>

## internal/notes

Similar to the users package, 'notes' handles all business logic related to 'notes'.

## internal/pkg

pkg package contains all the packages which are to be consumed across multiple packages within the project. For instance the datastore package will be consumed by both users and notes package. I'm not really particular about the name _pkg_. This might as well be _utils_ or some other generic name of your choice.

### internal/pkg/datastore

The datastore package initializes `pgxpool.Pool` and returns a new instance. I'm using Postgres as the datastore in this sample app. Why create such a package? I for instance had to because the packages we are using for Postgres did not have readymade APM integration. So started off by writing methods which we use in the app (and not 100% mirroring of the library), with APM integration. Did the same for cachestore as well. And it gets us beautiful insights like the following:

<p align="center">
<img src="https://user-images.githubusercontent.com/1092882/86710556-baa07180-c038-11ea-8924-3b4d61db1476.png" alt="APM overall" width="384px" height="256px" style="margin-right: 16px" />
<img src="https://user-images.githubusercontent.com/1092882/86710547-b83e1780-c038-11ea-9829-b5585b3d599b.png" alt="APM 1 API" width="384px" height="256px" />
</p>

P.S: Similar to logger, we made these independent private packages hosted in our [VCS](https://en.wikipedia.org/wiki/Version_control). Shoutout to [Gitlab](https://gitlab.com/)!

### internal/pkg/logger

I usually define the logging interface as well as the package, in a private repository (internal to your company e.g. vcs.yourcompany.io/gopkgs/logger), and is used across all services. Logging interface helps you to easily switch between different logging libraries, as all your apps would be using the interface **you** defined (interface segregation principle from SOLID). But here I'm making it part of the application itself as it has fewer chances of going wrong when trying to cater to a larger audience.

**Logging might sound trivial but there are a few questions around it:**

1. Should it be made a dependency of all packages, or can it be global?

Logging just like any other dependency, is a dependency. And in most cases it's better to write packages (code in general) which have as few dependencies as practically possible. This is a general principle, fewer dependencies make a lot of things easier like maintainability, testing, porting, refactoring, etc. And creating singleton Globals bring in restrictions, also it's a dependency nevertheless. Global instances have another issue, it doesn't give you flexibility when you need varying functionality across different packages (since it's global, it's common for all consumers). E.g. in one package you'd like to have debug logs, and in the other you'd only want errors. So in my opinion, it's better not to use a global instance, but have global functions which implement the default behaviour for all your packages which do not have any custom requirements.

2. Where would you do it? Should you bubble up errors and log at the parent level, or write where the error occurs?

Keeping it at the root/outermost layer helps make things easier because you need to worry about injecting logging dependency only in this package. And easier to controls it in general. i.e. One less thing to worry about in majority of the code.

For developers, while troubleshooting (which is one of the foremost need for logging), the line number along with filename helps a lot. Then it's obvious, log where the error occurs, right?

Over the course of time, I found it's not really obvious. The more nested function calls you have, higher the chances of redundant logging. And setting up guidelines for your developers to only log at the origin of error is also not easy. A lot of developers get confused which level should be considered the origin (especially when there's deep nesting fn1 -> fn2 -> fn3 -> fn4). Thus I prefer logging at the Handlers layer, [with annotated errors](https://pkg.go.dev/errors)(using the '%w' verb in `fmt.Errorf`) to trace its origin. Recently I introduced a [minimal error handling package](https://github.com/bnkamalesh/errors/) which gives long file path, line number of the origin of error, stacktrace etc. as well as help set user friendly messages for API response. My earlier recommendation was to use API package for logging, but in the past 2+ years (> 2019), figured out that it's better/easier to handle in the handler layer. Now all the HTTP handlers return an error, and there's a wrapper to handle the logging (this is updated in the app as well) as well as responding to the HTTP request.

Though there are some exceptions to logging at the outer most layer alone, consider the case of `internal/users` package. I'm making use of cache, but it's a read-through cache. So even if there's a miss in cache or cache store is down altogether, the system should still work (a specific business logic). But then how do you find out if your cache is down when there are no logs? Hence you see the logger being made a dependency of users package. This would apply to any asynchronous behaviours as well, e.g. a queue subscriber

## internal/server/http

All HTTP related configurations and functionalities are kept inside this package. The naming convention followed for filenames, is also straightforward. i.e. all the HTTP handlers of a specific package/domain are grouped under `handlers_<business logic unit name>.go`. The special mention of naming handlers is because, often for decently large web applications (especially when building REST-ful services) you end up with a lot of handlers. I have services with 100+ handlers for individual APIs, so keeping them organized helps.

e.g. handlers_users.go. The advantage of naming this way is, it's easier for developers to look at and identify from a list of filenames. e.g. on VS code it looks like this

<p align="center"><img src="https://user-images.githubusercontent.com/1092882/86526182-24d8db00-beae-11ea-9681-0a31b2d67e1b.png" alt="handlers_users.go" width="512px"/></p>

### internal/server/http/web/templates

All HTML templates required for the application are to be put here. Sub directories based on the main business logic unit, e.g. we/templates/users, can be created if required. It is highly unlikely that HTML templates used for HTTP responses are reused elsewhere in the application. Hence it justifies its location within 'server/http'. Other static files shall also be made part of the `web` directory like `web/static/images`, `web/static/js` etc. Feel free to [embed](https://pkg.go.dev/embed) templates, static files etc.

## lib

This name is quite explicit and if you notice, it's outside of the special 'internal' directory. So any exported name or entity within this directory, is meant to be used in external projects.

It might seem redundant to add a sub-directory called 'goapp', the import path would be `github.com/bnkamalesh/goapp/lib/goapp`. Though this is not a mistake, while importing this package, you'd like to use it as follows `goapp.<something>`. So if you directly put it under lib, it'd be `lib.` and that's obviously too generic and you'd have to manually setup aliases every time. Or if you try solving it by having the package name which differ from the direcory name, it's going to be a tussle with your [IDE](https://en.wikipedia.org/wiki/Integrated_development_environment).

Another advantage is, if you have more than one package which you'd like to be made available for external consumption, you create `lib/<other>`. In this case, you reduce the dependencies which are imported to external functions. On the contrary if you put everything inside `lib` or in a single package, you'd be forcing import of all dependencies even when you'd need only a small part of it.

## docker

I've been a fan of Docker since a few years now. I like keeping a dedicated folder for Dockerfile, in anticipation of introducing multiple Docker files or maintaining other files required for Docker image build.

e.g. [Dockerfiles for Alpine & Debian based images](https://github.com/bnkamalesh/golang-dockerfile)

You can create the Docker image for the sample app provided:

```bash
$ git clone https://github.com/bnkamalesh/goapp.git
$ cd goapp
# Update the internal/configs/configs.go with valid datastore configuration. Or pass nil while calling user service. This would cause the app to panic when calling any API with database interaction
# Build the Docker image
$ docker build -t goapp -f docker/Dockerfile .
# and you can run the image with the following command
$ docker run -p 8080:8080 --rm -ti goapp
```

## schemas

All the SQL schemas required by the project in this directory. This is not nested inside individual package because it's not consumed by the application at all. Also the fact that, actual consumers of the schema (developers, DB maintainers etc.) are varied. It's better to make it easier for all the audience rather than just developers. Even if you use NoSQL databases, your application would need some sort of schema to function, which can still be maintained inside this.

I've recently started using [sqlc](https://sqlc.dev/) for code generation for all SQL interactions (and love it!). I use [Squirrel](https://github.com/Masterminds/squirrel) whenever I need to dynamically build queries. E.g. when updating a table, you want to update only certain columns based on the input. Also, this is a recommendation from a friend for maintaining SQL migrations, though I've never used it myself, [Goose](https://github.com/pressly/goose).

Even migrations can be maintained in a directory in the root, but it's best to keep the application never be responsible for database setup. i.e. let migrations, index creation etc. be handled outside the scope of the application itself. For instance, it's very easy to create deadlocks with databases if it's part of the application, when you deploy the application in a _horizontally_ scaled model.

## main.go

Finally the `main package`. I prefer putting the `main.go` file outside as shown here. No non-sense, straight up `go run main.go` would start the application (provided the required configurations are available). 'main' is probably going to be the ugliest package where all conventions and separation of concerns are broken, but this is acceptable. The responsibility of main package is one and only one, **get things started**.

`cmd` directory can be added in the root for adding multiple commands. This is usually required _when there are multiple modes of interacting with the application_. i.e. HTTP server, CLI etc. In which case each usecase can be initialized and started with subpackages under `cmd`. Even though Go advocates fewer use of packages, I would give higher precedence for separation of concerns at a package level to keep things tidy. And even the main.go can be in `cmd/main.go`.


## Dependency flow

<p align="center">
<img src="https://user-images.githubusercontent.com/1092882/104085767-f5999100-5277-11eb-808a-5fd9b6776ad6.png" alt="Dependency flow between the layers" width="768px"/>
</p>

# Note

You can clone this repository and actually run the application, it'd start an HTTP server listening on port 8080 with the following routes available.

- `/` GET, the root just returns "Hello world" text response
- `/-/health` GET, returns a JSON with some basic info. I like using this path to give out the status of the app, its dependencies etc
- `/users` POST, to create new user
- `/users/:ID` GET, reads a user from the database given the email id. e.g. http://localhost:8080/users/1
