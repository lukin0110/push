# Releasing

## 1. Bump version
* change the version in [version.go](../version/version.go)
* change the version urls in [install.sh](../install.sh)

## 2. Create binaries
```
$ docker-compose run app release
```

## 3. Create tag

https://git-scm.com/book/en/v2/Git-Basics-Tagging

```
$ git tag 0.0.2beta
$ git push origin 0.0.2beta
```

## 4. Upload binaries

Uploaded the generated binaries on the 
[Github Release](https://github.com/lukin0110/push/releases) page.
