### Example LookUp Output

#### virustotal
| Ratio | Link                                                                                                                          | API    | Scanned                |
| ----- | ----------------------------------------------------------------------------------------------------------------------------- | ------ | ---------------------- |
| 95%   | [link](https://www.virustotal.com/file/371d99fc5514f5a9816b4ec844cb816c52460a41b8e5d14bac1cb7bee57e0b1f/analysis/1312464222/) | Public | Thu 2011Aug04 13:23:42 |
#### shadow-server
 - Not found
##### AntiVirus
 - FirstSeen: 6/15/2010 3:09AM
 - LastSeen: 6/15/2010 3:09AM

| Vendor           | Signature                  |
| ---------------- | -------------------------- |
| G-Data           | Trojan.Generic.2609117     |
| NOD32            | Win32/AutoRun.VB.JP        |
| Norman           | Suspicious_Gen2.SKLJ       |
| QuickHeal        | Worm.VB.at                 |
| Sophos           | Troj/DwnLdr-HQY            |
| VBA32            | Trojan.VBO.011858          |
| AVG7             | Downloader.Generic9.URM    |
| DrWeb            | Win32.HLLW.Autoruner.6014  |
| F-Secure         | Worm:W32/Revois.gen!A      |
| Ikarus           | Trojan-Downloader.Win32.VB |
| Kaspersky        | Trojan.Win32.Cosmu.nyl     |
| McAfee           | Generic                    |
| TrendMicro       | TROJ_DLOADR.SMM            |
| AntiVir          | WORM/VB.NVA                |
| Clam             | Trojan.Downloader-50691    |
| Vexira           | Trojan.DL.VB.EEDT          |
| VirusBuster      | Worm.VB.FMYJ               |
| Avast-Commercial | Win32:Zbot-LRA             |
| F-Prot6          | W32/Worm.BAOX              |
| Panda            | W32/OverDoom.A             |

### Example Scan Output

#### File
| Field  | Value                                                                    |
| ------ | ------------------------------------------------------------------------ |
| Name   | befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408         |
| Path   | samples/befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408 |
| Size   | 40.96 kB                                                                 |
| MD5    | 669f87f2ec48dce3a76386eec94d7e3b                                         |
| SHA1   | 6b82f126555e7644816df5d4e4614677ee0bda5c                                 |
| SHA256 | befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408         |
| Mime   | application/x-dosexec                                                    |
#### SSDeep
768:15jQ4nVHQaeO379u4XckKVCsknBN9A4hUnDxDiNZ957ZpK0IUUiM95Zdz:15jQ4nVHQaeO9uwckKuBN9A4UnDxcbFi

#### TRiD
 -  31.0% (.EXE) Win32 Executable MS Visual C++ (generic) (31206/45/13)
 -  27.5% (.EXE) Win64 Executable (generic) (27644/34/4)
 -  26.4% (.EXE) Win32 EXE Yoda's Crypter (26569/9/4)
 -  6.5% (.DLL) Win32 Dynamic Link Library (generic) (6578/25/2)
 -  4.4% (.EXE) Win32 Executable (generic) (4508/7/1)

#### Exiftool
| Field                   | Value                                           |
| ----------------------- | ----------------------------------------------- |
| Legal Copyright         | (C) Microsoft Corporation. All rights reserved. |
| Product Version Number  | 6.0.2930.2180                                   |
| File Flags Mask         | 0x003f                                          |
| Object File Type        | Unknown                                         |
| Internal Name           | iexplore                                        |
| Subsystem Version       | 4.0                                             |
| Comments                |                                                 |
| Company Name            | Microsoft Corporation                           |
| Legal Trademarks        |                                                 |
| ExifTool Version Number | 10.09                                           |
| File Size               | 40 kB                                           |
| File Type Extension     | exe                                             |
| Machine Type            | Intel 386 or later, and compatibles             |
| MIME Type               | application/octet-stream                        |
| Code Size               | 20480                                           |
| Uninitialized Data Size | 0                                               |
| Language Code           | Neutral                                         |
| PE Type                 | PE32                                            |
| Linker Version          | 6.0                                             |
| File Description        | Internet Explorer                               |
| File Version            | 6.00.2900.2180 (xpsp_sp2_rtm.040803-2158)       |
| Image Version           | 0.0                                             |
| File Flags              | Private build                                   |
| Private Build           |                                                 |
| Product Name            | Microsoft(R) Windows(R) Operating System        |
| Character Set           | Unicode                                         |
| Product Version         | 6.00.2900.2180                                  |
| File Type               | Win32 EXE                                       |
| Initialized Data Size   | 20480                                           |
| File Version Number     | 6.0.2930.2180                                   |
| File OS                 | Unknown (0)                                     |
| Original File Name      | IEXPLORE.EXE                                    |
| Special Build           |                                                 |
| Entry Point             | 0x5a46                                          |
| OS Version              | 4.0                                             |
| Subsystem               | Windows GUI                                     |
| File Subtype            | 0                                               |
#### yara
| Rule                                   | Description                                 | Offset | Data                                 | Tags |
| -------------------------------------- | ------------------------------------------- | ------ | ------------------------------------ | ---- |
| _Microsoft_Visual_Cpp_v50v60_MFC_      | Microsoft Visual C++ v5.0/v6.0 (MFC)        | 5204   | U���                                 |      |
| _Borland_Delphi_v60__v70_              | Borland Delphi v6.0 - v7.0                  | 5204   | U��                                  |      |
| _dUP_v2x_Patcher__wwwdiablo2oo2cjbnet_ | dUP v2.x Patcher --> www.diablo2oo2.cjb.net | 78     | This program cannot be run in DOS mo |      |
| _Free_Pascal_v106_                     | Free Pascal v1.06                           | 14866  | ���@O�k                            |      |
| _Armadillo_v171_                       | Armadillo v1.71                             | 23110  | U��j�h b@h�[@d�                      |      |
#### AVG
| Infected | Result                | Engine    | Updated |
| -------- | --------------------- | --------- | ------- |
| true     | Found Win32/DH{YQMT?} | 13.0.3114 | 2016213 |
#### ClamAV
| Infected | Result               | Engine | Updated  |
| -------- | -------------------- | ------ | -------- |
| true     | Win.Trojan.Backspace | 0.99   | 20160214 |
#### F-PROT
| Infected | Result | Engine    | Updated |
| -------- | ------ | --------- | ------- |
| false    |        | 4.6.5.141 | 2016213 |
