# `findthese`

Check URL for files from given source folder. Like `dirb` or `gobuster` but we know what we looking for.
```bash
findthese --src ../framework --url https://framework.xx/
```

## Installation
```bash
go get -v -u github.com/briiC/findthese
```

## Usage

```bash
# Example: Check endpoint for phpmyadmin files - which are accessible from internet
git clone https://github.com/phpmyadmin/phpmyadmin
findthese --src ./phpmyadmin --url https://some-site.xx/pma/
```
_NOTE_: You can clone different version of phpmyadmin if you know endpoint uses that version.



```
Flags:
     --version  Displays the program version string.
  -h --help  Displays help with available flag, subcommand, and positional value parameters.
  -s --src  Source path of directory -- REQUIRED
  -u --url  URL endpoint to hit -- REQUIRED
  -m --method  HTTP Method to use (default: HEAD) (default: HEAD)
  -o --output  Output report to file (default: ./findthese.report) (default: ./findthese.report)
     --depth  How deep go in folders. '0' no limit  (default: 0) (default: 0)
  -z --delay  Delay every request for N milliseconds (default: 150) (default: 150)
     --timeout  Timeout (seconds) to wait for response  (default: 10) (default: 10)
     --mutations  Mutations of checked file (default: [~ .swp .swo .tmp .dmp .bkp .backup .bak .zip .tar .old _* ~*]) (default: ~,.swp,.swo,.tmp,.dmp,.bkp,.backup,.bak,.zip,.tar,.old,_*,~*)
     --skip  Skip files with these extensions (default: [jquery css img images i18n po]) (default: jquery,css,img,images,i18n,po)
     --skip-ext  Skip files with these extensions (default: [.png .jpeg jpg Gif .CSS .less .sass]) (default: .png,.jpeg,jpg,Gif,.CSS,.less,.sass)
     --skip-code  Skip responses with this response HTTP code (default: [404]) (default: 404)
     --skip-size  Skip responses with this body size (default: [])
     --skip-content  Skip responses if given content found
     --dir-only  Scan directories only
     --user-agent  User-Agent used (default: random)
  -C --cookie  Cookie string sent with requests
  -H --headers  Custom Headers sent with requests

```


### TODO
- On key `p` pause scan. Run same command with additional params (fine-tuning) and scan will resume from previous with new settings. (Detects same src and url)
- Mark placeholder for file to put in URL: https://example.com?f=^FILE^&auth=john
- [--mode=info|download] (default: info)
- file download path where to download all files
- On first success pause show content and press Y to continue
- read robots and .htaccess for paths
- Generate wordlist based on source directory
- Use `tor` by default (golang sockets transport)
-
