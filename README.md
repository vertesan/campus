# campus

Background cronjob for a certain repository.

> [!IMPORTANT]  
> This program is not designated for any evil purpose however according to the way it works, it can be, just like any other technologies, altered for malicious intention. We choose to make it open source because we believe the innovation and creativity of the whole great open source community can create something amazing based on it. The developers of the project hereby declare that we strongly oppose any malicious uses against the original purpose. Please always remember it's not the technology do the harm, it's the people wielding it.

## Usage

You can choose to build from source on your own or utilize our pre-built docker image to use.

### Building

```bash
go mod download
env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=gcc go build -ldflags '-s -w' .
```

Some dependencies requires [CGO](https://pkg.go.dev/cmd/cgo) enabled. According to your developing environment, you may need to change `GOOS=windows CC=x86_64-w64-mingw32-gcc` or other related environment variables, please refer to the [documentation](https://pkg.go.dev/cmd/go#hdr-Environment_variables) for detailed information.
Required dependencies are listed in [go.mod](./go.mod).

### Docker

Get our pre-built docker image from https://github.com/vertesan/campus/pkgs/container/campus.

You must create an empty file named `config.yaml` and an empty directory named `cache` in advance if they don't exist.

```bash
touch config.yaml
mkdir cache
```

After that run `docker compose up` to see the results. By default the container run with an option `-h` to print the CLI instructions and returns, you can change this behaviour to fit your needs by editing [`compose.yaml`](./compose.yaml).

## CLI

After building the executable, run `./campus -h` to get instructions. Here is a glimpse.

```
Usage of ./campus:
  -ab
        Download and deobfuscate assetbundles if true.
        Deobfuscated files are saved in 'cache/assets' directory.
  -analyze
        Analyze dump.cs to retrieve proto schema.
        Generated codes are saved in 'cache/GeneratedProto' directory.
  -db
        Download and decrypt master database if true.
        Generated yaml files are saved in 'cache/masterYaml' directory.
  -forceab
        Download assetbundles without checking version.
        Takes no effect if 'ab' is absent.
        It's safe to set this flag to true if you only want to download a part of additional assets instead of the entire bulky thing because MD5 check will still be carried out before downloading.
  -forcedb
        Download and decrypt master database without checking local version.
        Take no effect if 'db' flag is absent.
  -keepab
        Do not delete obfuscated assetbundle files after deobfuscating.
        Take no effect if 'ab' flag is absent.
  -keepdb
        Do not delete encrypted master database files after decrypting.
        Take no effect if 'db' flag is absent.
  -token string
        The refresh token used to retrieve login idToken from firebase.
        If refreshToken field set in 'config.yaml' is not empty, the value in the config file will take precedence.
```

## Token

To make campus work properly, you must hand over your firebase refresh token as an option at the first run. After that your token will be saved in `./config.yaml` for subsequent use so that you can omit adding token from then on.

P.S. To get the firebase refresh token you will need to find a way to intercept the HTTPS traffic and that is not the content of this README.

## Configurations

After the first successful running, a `config.yaml` will be automatically created in the root directory. You can safely edit it if you want to tweak something manually.

## Scripts

- [`cronjob.sh`](./cronjob.sh): A cronjob script for a certain repository
- [`push_master.sh`](./push_master.sh): A post process script running after `cronjob.sh`
- [`analyze.sh`](./analyze.sh): Analyze the latest proto schema from `dump.cs` and automatically update the files inside `./proto`
- [`unpack.py`](./unpack.py): Unpack images from Unity assetbundle format to PNG

## License

AGPL-3.0 license
