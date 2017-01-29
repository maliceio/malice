# Scan  

```bash
$ malice scan befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408
```

=OR=

```
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
                -v `pwd`/data/samples:/malice/samples \
                -e MALICE_VT_API=$MALICE_VT_API \
                malice/engine scan befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408
```

#### File

| Field  | Value                                                                         |
|--------|-------------------------------------------------------------------------------|
| Name   | befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408              |
| Path   | data/samples/befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408 |
| Size   | 40.96 kB                                                                      |
| MD5    | 669f87f2ec48dce3a76386eec94d7e3b                                              |
| SHA1   | 6b82f126555e7644816df5d4e4614677ee0bda5c                                      |
| SHA256 | befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408              |
| Mime   | application/x-dosexec                                                         |

#### VirusTotal

| Ratio | Link                                                                                                                          | API    | Scanned             |
|-------|-------------------------------------------------------------------------------------------------------------------------------|--------|---------------------|
| 85%   | [link](https://www.virustotal.com/file/befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408/analysis/1455536823/) | Public | 2016-02-15 11:47:03 |

#### #totalhash

| Found              | URL                                                                                    |
| ------------------ | -------------------------------------------------------------------------------------- |
| :white_check_mark: | [link](https://totalhash.cymru.com/analysis/?6b82f126555e7644816df5d4e4614677ee0bda5c) |

#### NSRL Database

 - Not Found :grey_question:

#### ShadowServer

##### AntiVirus
 - FirstSeen: 11/03/2016 9:09AM
 - LastSeen: 11/08/2016 2:10AM

| Vendor      | Signature                              |
| ----------- | -------------------------------------- |
| QuickHeal   | TrojanAPT.LecnaCBack.MUE.Z3            |
| MicroWorld  | Backdoor.Lecna.AB                      |
| Fortinet    | W32/Generic.AC.1FAF5A!tr               |
| Ikarus      | Backdoor.Win32.Lecna                   |
| Eset        | Win32/Lecna.W                          |
| Clam        | PUA.Win.Packer.AcprotectUltraprotect-1 |
| DrWeb       | BackDoor.Dizhi                         |
| Symantec    | W32.Lecna.D                            |
| Avast       | Win32:Lecna-I                          |
| McAfee      | BackDoor-CSB                           |
| Sophos      | Troj/Lecna-Q                           |
| Comodo      | Backdoor.Win32.Lecna.AB                |
| AhnLab      | Win-Trojan/Lecna.61440                 |
| K7GW        | Trojan ( 00013a571 )                   |
| GData       | Backdoor.Lecna.AB                      |
| TrendMicro  | BKDR_LECNA.SM                          |
| Microsoft   | Backdoor:Win32/Lecna!dha               |
| K7          | Trojan ( 00013a571 )                   |
| Avira       | WORM/Rbot.Gen                          |
| Authentium  | W32/Lecnac.A.gen!Eldorado              |
| AVG         | Win32/DH{IIEXEx4XA2EICXwkIiU}          |
| BitDefender | Backdoor.Lecna.AB                      |
| FProt       | W32/Trojan.AAWD                        |

#### Yara

| Rule                                    | Description                                 | Offset | Data                                 | Tags |
|-----------------------------------------|---------------------------------------------|--------|--------------------------------------|------|
| *Microsoft_Visual_Cpp_v50v60_MFC*       | Microsoft Visual C++ v5.0/v6.0 (MFC)        | 5204   | U���                                 |      |
| *Borland_Delphi_v60\__v70*              | Borland Delphi v6.0 - v7.0                  | 5204   | U��                                  |      |
| *dUP_v2x_Patcher\__wwwdiablo2oo2cjbnet* | dUP v2.x Patcher --> www.diablo2oo2.cjb.net | 78     | This program cannot be run in DOS mo |      |
| *Free_Pascal_v106*                      | Free Pascal v1.06                           | 14866  | ���@O�k                              |      |
| *Armadillo_v171*                        | Armadillo v1.71                             | 23110  | U��j�h b@h�[@d�                      |      |

#### SSDeep

768:15jQ4nVHQaeO379u4XckKVCsknBN9A4hUnDxDiNZ957ZpK0IUUiM95Zdz:15jQ4nVHQaeO9uwckKuBN9A4UnDxcbFi

#### TRiD

-	31.0% (.EXE) Win32 Executable MS Visual C++ (generic) (31206/45/13)
-	27.5% (.EXE) Win64 Executable (generic) (27638/28/4)
-	26.4% (.EXE) Win32 EXE Yoda's Crypter (26569/9/4)
-	6.5% (.DLL) Win32 Dynamic Link Library (generic) (6578/25/2)
-	4.4% (.EXE) Win32 Executable (generic) (4508/7/1)

#### Exiftool

| Field                   | Value                                           |
|-------------------------|-------------------------------------------------|
| Special Build           |                                                 |
| Code Size               | 20480                                           |
| File Version            | 6.00.2900.2180 (xpsp_sp2_rtm.040803-2158)       |
| Legal Trademarks        |                                                 |
| Product Name            | Microsoft(R) Windows(R) Operating System        |
| Machine Type            | Intel 386 or later, and compatibles             |
| PE Type                 | PE32                                            |
| File Version Number     | 6.0.2930.2180                                   |
| Character Set           | Unicode                                         |
| Comments                |                                                 |
| MIME Type               | application/octet-stream                        |
| Linker Version          | 6.0                                             |
| Product Version Number  | 6.0.2930.2180                                   |
| File Flags              | Private build                                   |
| File OS                 | Unknown (0)                                     |
| File Description        | Internet Explorer                               |
| File Size               | 40 kB                                           |
| Object File Type        | Unknown                                         |
| Legal Copyright         | (C) Microsoft Corporation. All rights reserved. |
| Original File Name      | IEXPLORE.EXE                                    |
| Uninitialized Data Size | 0                                               |
| Image Version           | 0.0                                             |
| Subsystem               | Windows GUI                                     |
| File Flags Mask         | 0x003f                                          |
| Company Name            | Microsoft Corporation                           |
| Product Version         | 6.00.2900.2180                                  |
| Initialized Data Size   | 20480                                           |
| Entry Point             | 0x5a46                                          |
| OS Version              | 4.0                                             |
| File Subtype            | 0                                               |
| Language Code           | Neutral                                         |
| Internal Name           | iexplore                                        |
| File Type Extension     | exe                                             |
| File Type               | Win32 EXE                                       |
| Subsystem Version       | 4.0                                             |
| Private Build           |                                                 |
| ExifTool Version Number | 10.23                                           |

#### ClamAV

| Infected | Result                 | Engine | Updated  |
|----------|------------------------|--------|----------|
| true     | Win.Trojan.Backspace-1 | 0.99.2 | 20160919 |

#### Comodo

| Infected | Result                  | Engine | Updated |
|----------|-------------------------|--------|---------|
| true     | Backdoor.Win32.Lecna.AB | 1.1    |         |

#### F-Secure

| Infected | Result            | Engine         | Updated  |
|----------|-------------------|----------------|----------|
| true     | Backdoor.Lecna.AB | 11.00 build 79 | 20160919 |

#### F-PROT

| Infected | Result | Engine    | Updated  |
|----------|--------|-----------|----------|
| false    |        | 4.6.5.141 | 20160919 |

#### Avast

| Infected | Result              | Engine | Updated  |
| -------- | ------------------- | ------ | -------- |
| true     | Win32:Lecna-I [Trj] | 2.1.2  | 20170129 |

#### AVG

| Infected | Result                | Engine    | Updated  |
|----------|-----------------------|-----------|----------|
| true     | Found Win32/DH{YQMT?} | 13.0.3114 | 20160918 |

#### Bitdefender

| Infected | Result            | Engine  | Updated  |
|----------|-------------------|---------|----------|
| true     | Backdoor.Lecna.AB | 7.90123 | 20160919 |

#### Sophos

| Infected | Result       | Engine | Updated  |
| -------- | ------------ | ------ | -------- |
| true     | Troj/Lecna-Q | 5.27.0 | 20160920 |

#### Floss

##### Decoded Strings

Location: `0x402830`
 - `################################################################################################################################################################################################################################################################################################################################`

Location: `0x401059`
 - `*lecnaC*`
 - `Software\Microsoft\CurrentNetInf`
 - `SYSTEM\CurrentControlSet\Control\Lsa`
 - `Software\Microsoft\Windows\CurrentVersion\Policies\Explorer\Run`
 - `MicrosoftZj`
 - `LhbqnrnesDwhs`
 - `MicrosoftHaveExit`
 - `LhbqnrnesG`ud@bj`
 - `IEXPLORE.EXE`
 - `/ver.htm`
 - `/exe.htm`
 - `/app.htm`
 - `/myapp.htm`
 - `/hostlist.htm`
 - `.a`j-gsl`
 - `/SomeUpList.htm`
 - `/SomeUpVer.htm`
 - `www.flyeagles.com`
 - `www.km-nyc.com`
 - `/restore`
 - `/dizhi.gif`
 - `/connect.gif`
 - `\$NtUninstallKB900727$`
 - `\netsvc.exe`
 - `\netscv.exe`
 - `\netsvcs.exe`
 - `System Idle Process`
 - `Program Files`
 - `\Internet Exp1orer`
 - `forceguest`
 - `AudioPort`
 - `AudioPort.sys`
 - `SYSTEM\CurrentControlSet\Services`
 - `SYSTEM\ControlSet001\Services`
 - `SYSTEM\ControlSet002\Services`
 - `\drivers\`
 - `\DriverNum.dat`

Location: `0x40511A`
 - `\A|{@`

Location: `0x404DDE`
 - `SMBs`
 - `NTLMSSP`
 - `Windows 2000 2195`
 - `Windows 2000 5.0`
 - `SMBr`
 - `PC NETWORK PROGRAM 1.0`
 - `LANMAN1.0`
 - `Windows for Workgroups 3.1a`
 - `LM1.2X002`
 - `LANMAN2.1`
 - `NT LM 0.12`

Location: `0x401047`
 - `Ie_nkokbpAtep`
 - `+^]g*dpi`
 - `Ie_nkokbpD]ra=_g`

##### Stack Strings

 - `\A|{@`
 - `CAAA\`
 - `cmd.exe`
