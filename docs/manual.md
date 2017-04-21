# LOGO


### Doc Version

## Visual Architecture



       +-------------------+                              +-------------------+
       |      BACSRV       |                              |       BACLNT      |
       |-------------------|                              |-------------------|
       |                   |                              |                   |
       |   +-----------+   |                              |   +-----------+   |
       |   | Management|+------------------------------------>| Management|   |
       |   | Srv/Clnt  |   |                              |   | Srv/Clnt  |   |
       |   |           |   |                              |   |           |   |
       |   |           |<------------------------------------+|           |   |
       |   +-----------+   |                              |   +-----------+   |
       |                   |                              |                   |
       |                   |                              |                   |
       |                   |                              |                   |
       |   +-----------+   |                              |   +-----------+   |
       |   | Data      |+------------------------------------>| Data      |   |
       |   | Srv/Clnt  |   |                              |   | Srv/Clnt  |   |
       |   |           |   |                              |   |           |   |
       |   |           |<------------------------------------+|           |   |
       |   +-----------+   |                              |   +-----------+   |
       |                   |                              |                   |
       |                   |                              |                   |
       |                   |                              |                   |
       +-------------------+                              +-------------------+


# BACKER MANUAL

## TABLE OF CONTENTS

<!-- TOC -->

- [LOGO](#logo)
        - [Doc Version](#doc-version)
    - [Visual Architecture](#visual-architecture)
- [BACKER MANUAL](#backer-manual)
    - [TABLE OF CONTENTS](#table-of-contents)
- [1. GENERAL INFORMATION](#1-general-information)
- [2. CONFIGURATION FILES](#2-configuration-files)
    - [2.1. Server config file](#21-server-config-file)
        - [2.1.1. Main part](#211-main-part)
        - [2.1.2. Repository part](#212-repository-part)
        - [2.1.3. Clients part](#213-clients-part)
        - [2.1.4. Example](#214-example)
    - [2.2. Jobs config file](#22-jobs-config-file)

<!-- /TOC -->

# 1. GENERAL INFORMATION

# 2. CONFIGURATION FILES
## 2.1. Server config file

### 2.1.1. Main part
| Configuration key name  |  Type of data   |  Default Value  |  Remark  |
|---|---|---|---|
| MgmtPort  |  string  | 8090 |    |
| DataPort  |  string  | 8000  |   |
| RestApiPort  |  string  | 8080   |   |
| LogOutput |  output  | STDOUT |
| Debug |  bool  | false |
| ExternalName | string | 
| DataTransferInterface | address | 

### 2.1.2. Repository part
| Configuration key name  |  Type of data   |  Default Value  |  Remark  |
|---|---|---|---|---|
| Localization | path |   | required |


### 2.1.3. Clients part
| Configuration key name  |  Type of data   |  Default Value  |  Remark  |
|---|---|---|---|---|
| ConfigFile | path | /etc/backer/backer.d/clients/ |

### 2.1.4. Example
```toml
[server]
MgmtPort = 8090
DataPort = 8000
RestApiPort = 8080
LogOutput = "STDOUT"
Debug = true
ExternalName = "127.0.0.1" # External name of server, can be hosts/ip/DNS
DataTransferInterface = "127.0.0.1" # IP addres on which data server will listen

[repository]
Localization = "/home/dixi/repository"

[clients]
ConfigFile = "/mnt/c/Users/Damian Rajca/dev/go/src/github.com/damekr/backer/config/clients"

```

## 2.2. Jobs config file
| Configuration key name  |  Type of data   |  Default Value  |  Remark  |
|---|---|---|---|---|
