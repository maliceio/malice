### Plugin List

```bash
$ malice plugin list --all --detail
```

| Name             | Description                               | Enabled | Image                   | Category | Mime                   |
|:-----------------|:------------------------------------------|:--------|:------------------------|:---------|:-----------------------|
| nsrl             | NSRL Database Hash Search                 | true    | malice/nsrl:sha1        | intel    | hash                   |
| virustotal       | VirusTotal - files scan and hash lookup   | true    | malice/virustotal       | intel    | hash                   |
| totalhash        | #totalhash - hash lookup                  | false   | malice/totalhash        | intel    | hash                   |
| shadow-server    | ShadowServer - hash lookup                | true    | malice/shadow-server    | intel    | hash                   |
| team-cymru       | TeamCymru - hash lookup                   | false   | malice/team-cymru       | intel    | hash                   |
| fileinfo         | ssdeep/TRiD/exiftool                      | true    | malice/fileinfo         | metadata | *                      |
| yara             | YARA Scan                                 | true    | malice/yara             | av       | *                      |
| avast            | Avast AntiVirus                           | true    | malice/avast            | av       | *                      |
| avg              | AVG AntiVirus                             | true    | malice/avg              | av       | *                      |
| avira            | Avira AntiVirus                           | false   | malice/avira            | av       | *                      |
| bitdefender      | Bitdefender AntiVirus                     | true    | malice/bitdefender      | av       | *                      |
| clamav           | ClamAV                                    | true    | malice/clamav           | av       | *                      |
| comodo           | Comodo AntiVirus                          | true    | malice/comodo           | av       | *                      |
| escan            | eScan AntiVirus                           | true    | malice/escan            | av       | *                      |
| fprot            | F-PROT AntiVirus                          | true    | malice/fprot            | av       | *                      |
| fsecure          | F-Secure AntiVirus                        | true    | malice/fsecure          | av       | *                      |
| sophos           | Sophos AntiVirus                          | true    | malice/sophos           | av       | *                      |
| windows-defender | Windows Defender AntiVirus                | false   | malice/windows-defender | av       | *                      |
| zoner            | ZonerAntiVirus                            | true    | malice/zoner            | av       | *                      |
| exe              | PE - tool to triage portable executables  | false   | malice/exe              | exe      | application/x-dosexec  |
| floss            | FireEye Labs Obfuscated String Solver     | true    | malice/floss            | exe      | application/x-dosexec  |
| office           | Office - tool to triage OLE/RTF documents | false   | malice/office           | document | *                      |
| pdf              | PDF - tool to triage PDF documents        | false   | malice/pdf              | document | application/pdf        |
| javascript       | Javascript - tool to triage JS scripts    | false   | malice/javascript       | document | application/javascript |
| archive          | Archive - tool to unarchive archives      | false   | malice/archive          | archive  | archive                |
