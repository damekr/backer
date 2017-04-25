# Logo


# Doc Version
0.0.1 - DRAFT

# Visual Architecture



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

- [Logo](#logo)
- [Doc Version](#doc-version)
- [Visual Architecture](#visual-architecture)
- [BACKER MANUAL](#backer-manual)
    - [TABLE OF CONTENTS](#table-of-contents)
- [1. General information](#1-general-information)
- [2. Abbreviations](#2-abbreviations)
- [3. Use cases](#3-use-cases)
- [4. Feature Component specifications](#4-feature-component-specifications)
    - [4.1. FC_0 - Client server authentication](#41-fc_0---client-server-authentication)
    - [4.2. FC_1 - Data transfer encryption](#42-fc_1---data-transfer-encryption)
    - [4.3. FC_2 - Client integration](#43-fc_2---client-integration)
    - [4.4. FC_3 - Backup specific paths from client](#44-fc_3---backup-specific-paths-from-client)
    - [4.5. FC_4 - Restore data to given client](#45-fc_4---restore-data-to-given-client)
    - [4.6. FC_5 - List available backups from client side](#46-fc_5---list-available-backups-from-client-side)
- [5. Configuration files](#5-configuration-files)
    - [5.1. Server config file](#51-server-config-file)
        - [5.1.1. Main part](#511-main-part)
        - [5.1.2. Repository part](#512-repository-part)
        - [5.1.3. Clients part](#513-clients-part)
        - [5.1.4. Example](#514-example)
    - [5.2. Jobs config file](#52-jobs-config-file)

<!-- /TOC -->

# 1. General information

The application BACKER is mainly written in golang. Uses new libraries checked and implemented with successes on many environments. Some of them:

- Logrus
- Cobra
- gRPC
- etc.

# 2. Abbreviations

# 3. Use cases

# 4. Feature Component specifications
In this chapter are stored information about features that are being implemented. Interfaces for this features will be described in part of this document dedicated to this. All features components are under development and can be later on change. 

## 4.1. FC_0 - Client server authentication

Client server authentication needs to be provided since in some unknown environments always we have to be sure that we are integrating client which was in our mind or starting data transfer to proper client or server. This feature component can be divided in to two parts: gRPC authentication and data transfer authentication. 

gRPC Authentication:

Since gRPC is well standarized protocol crated by google the authentication model will be implemented according to their suggestions. A snipped of code can be find below:

[gRPC Simple authentication](https://github.com/grpc/grpc-go/issues/106)

- In case of server after installation, configuration needs to be done. As a part of this configuration username and password to gRPC interface must be configured. This user and password will be used only by clients to provide management from client side.

- In case of client username and password must by configured in configuration file. 



## 4.2. FC_1 - Data transfer encryption

The application uses two types of connections: gRPC and raw TCP sockets. Either gRPC or TCP does not provide any encryption mechanizm. Data encryption is also very important topic in case of privacy. For this aim TLS encryption will be used. TLS encryption requires certificate generation for each client which is being integrated to the system. This mechanizm will be implemented as part of client integration procedure. Also for development purposes there will be a possibility to disable encryption completely and work on raw data. A procedure for generating CA for server will be described in the docs. This mechanizm and the same certificates will be used in both types of connections.



## 4.3. FC_2 - Client integration

Client integration is a procedure where we want to add a new node(client/host) to our B&R System. The procedure will be executed on server side which means that an administrator must be able to login over ssh into console of server where backer server is working and update clients configuration file with client specific values. These values are specified in part of clients configuration file of this doc. When client has been added you need to reload server application with command described in Administration part of this document. After that server will try to connect to the client and exchanged certificate for next connections. Also external name of server needs to be in client configuration file, without this server will not be able to proceed the operation and integration will fail. 


## 4.4. FC_3 - Backup specific paths from client

Backup specific paths is a basic operation of all Backup and Restore application. In case of backer backup can be triggered from server side. Later on there is an idea to provide also a possibility to execute this procedure from client side. This is operation called 'backup on demand' it means when an administrator wants to get data from specific client then can trigger the backup and get the data immediately. Second option of backup is scheduled backup it means that during the integration you can specify when automatic backup will be triggered and at this time backup will run. How to create schedule please follow part of this doc which is about this under client configuration file. 

## 4.5. FC_4 - Restore data to given client

Restore data is second after backup basic operation. In this feature component two main things needs to be working. First is a restore data on client which is already integrated to the server. It covers use case when an adminitrator wants to get data from server in case of wrong configuration of service for example. Another option is when host which was integrated to the server fails. All data are lost and restore from full backup needs to be done. In this case integration procedure must be performed. After that regular restore can be done because server will consider this client as an integrated and known client. Restore procedure will be executed from client side. 


## 4.6. FC_5 - List available backups from client side

Feature component FC_4 requires to provide an interface for an administrator where the administator will be able to see which type of backups are stored on the server. This FC can be threated as more generic, because ths system needs to have command line interface (CLI). So in case of listen backups or another operation this CLI will be provided as sepearate tool to cover this and another related features. The CLI will communicate with server or client service over Mgmt interface and perform operations which is needed. 


# 5. Configuration files
## 5.1. Server config file

### 5.1.1. Main part
| Configuration key name  |  Type of data   |  Default Value  |  Remark  |
|---|---|---|---|
| MgmtPort  |  string  | 8090 |    |
| DataPort  |  string  | 8000  |   |
| RestApiPort  |  string  | 8080   |   |
| LogOutput |  output  | STDOUT |
| Debug |  bool  | false |
| ExternalName | string | 
| DataTransferInterface | address | 
| Login | string | | required|
| Password | string |  | required |

### 5.1.2. Repository part
| Configuration key name  |  Type of data   |  Default Value  |  Remark  |
|---|---|---|---|
| Localization | path |   | required |


### 5.1.3. Clients part
| Configuration key name  |  Type of data   |  Default Value  |  Remark  |
|---|---|---|---|
| ConfigFile | path | /etc/backer/backer.d/clients/ |

### 5.1.4. Example
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

## 5.2. Jobs config file
| Configuration key name  |  Type of data   |  Default Value  |  Remark  |
|---|---|---|---|
